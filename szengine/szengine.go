/*
The szengine is a wrapper over the Senzing libg2 library.
*/
package szengine

/*
#include <stdlib.h>
#include "libg2.h"
#include "gohelpers/golang_helpers.h"
#cgo CFLAGS: -g -I/opt/senzing/g2/sdk/c
#cgo windows CFLAGS: -g -I"C:/Program Files/Senzing/g2/sdk/c"
#cgo LDFLAGS: -L/opt/senzing/g2/lib -lG2
#cgo windows LDFLAGS: -L"C:/Program Files/Senzing/g2/lib" -lG2
*/
import "C"

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szengineapi "github.com/senzing-garage/sz-sdk-go/szengine"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

type Szengine struct {
	isTrace        bool
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

const initialByteArraySize = 65535

// ----------------------------------------------------------------------------
// sz-sdk-go.SzEngine interface methods
// ----------------------------------------------------------------------------

/*
The AddRecord method adds a record into the Senzing repository.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing repository.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) AddRecord(ctx context.Context, dataSourceCode string, recordId string, recordDefinition string, flags int64) (string, error) {
	if (flags & sz.SZ_WITH_INFO) > 0 {
		finalFlags := flags ^ sz.SZ_WITH_INFO
		return client.addRecordWithInfo(ctx, dataSourceCode, recordId, recordDefinition, finalFlags)
	}
	return client.addRecord(ctx, dataSourceCode, recordId, recordDefinition)
}

/*
The CloseExport method closes the exported document created by ExportJsonEntityReport().
It is part of the ExportJsonEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of sized entities.

Input
  - ctx: A context to control lifecycle.
  - exportHandle: A handle created by ExportJsonEntityReport() or ExportCsvEntityReport().
*/
func (client *Szengine) CloseExport(ctx context.Context, exportHandle uintptr) error {
	//  _DLEXPORT int G2_closeExport(ExportHandle responseHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(13, exportHandle)
		defer func() { client.traceExit(14, exportHandle, err, time.Since(entryTime)) }()
	}
	result := C.G2_closeExport_helper(C.uintptr_t(exportHandle))
	if result != 0 {
		err = client.newError(ctx, 4006, exportHandle, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8006, err, details)
		}()
	}
	return err
}

/*
The CountRedoRecords method returns the number of records in need of redo-ing.

Input
  - ctx: A context to control lifecycle.

Output
  - The number of redo records in Senzing's redo queue.
*/
func (client *Szengine) CountRedoRecords(ctx context.Context) (int64, error) {
	//  _DLEXPORT long long G2_countRedoRecords();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(15)
		defer func() { client.traceExit(16, result, err, time.Since(entryTime)) }()
	}
	result = int64(C.G2_countRedoRecords())
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8007, err, details)
		}()
	}
	return result, err
}

/*
The DeleteRecord method deletes a record from the Senzing repository.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) DeleteRecord(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	if (flags & sz.SZ_WITH_INFO) > 0 {
		finalFlags := flags ^ sz.SZ_WITH_INFO
		return client.deleteRecordWithInfo(ctx, dataSourceCode, recordId, finalFlags)
	}
	return client.deleteRecord(ctx, dataSourceCode, recordId)
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2 object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) Destroy(ctx context.Context) error {
	//  _DLEXPORT int G2_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(21)
		defer func() { client.traceExit(22, err, time.Since(entryTime)) }()
	}
	result := C.G2_destroy()
	if result != 0 {
		err = client.newError(ctx, 4009, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8010, err, details)
		}()
	}
	return err
}

/*
The ExportCsvEntityReport method initializes a cursor over a document of exported entities.
It is part of the ExportCsvEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - csvColumnList: A comma-separated list of column names for the CSV export.
  - flags: Flags used to control information returned.

Output
  - A handle that identifies the document to be scrolled through using FetchNext().
*/
func (client *Szengine) ExportCsvEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error) {
	//  _DLEXPORT int G2_exportCSVEntityReport(const char* csvColumnList, const long long flags, ExportHandle* responseHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultExportHandle uintptr
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(27, csvColumnList, flags)
		defer func() { client.traceExit(28, csvColumnList, flags, resultExportHandle, err, time.Since(entryTime)) }()
	}
	csvColumnListForC := C.CString(csvColumnList)
	defer C.free(unsafe.Pointer(csvColumnListForC))
	result := C.G2_exportCSVEntityReport_helper(csvColumnListForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4012, csvColumnList, flags, result.returnCode, result, time.Since(entryTime))
	}
	resultExportHandle = (uintptr)(result.exportHandle)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8013, err, details)
		}()
	}
	return resultExportHandle, err
}

/*
The ExportCsvEntityReportIterator method creates an Iterator that can be used in a for-loop
to scroll through a document of exported entities.
It is a convenience method for the ExportCsvEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - csvColumnList: A comma-separated list of column names for the CSV export.
  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportCsvEntityReportIterator(ctx context.Context, csvColumnList string, flags int64) chan sz.StringFragment {
	stringFragmentChannel := make(chan sz.StringFragment)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(stringFragmentChannel)
		var err error = nil
		if client.isTrace {
			entryTime := time.Now()
			client.traceEntry(163, csvColumnList, flags)
			defer func() { client.traceExit(164, csvColumnList, flags, err, time.Since(entryTime)) }()
		}
		reportHandle, err := client.ExportCsvEntityReport(ctx, csvColumnList, flags)
		if err != nil {
			result := sz.StringFragment{
				Error: err,
			}
			stringFragmentChannel <- result
			return
		}
		defer func() {
			client.CloseExport(ctx, reportHandle)
		}()
	forLoop:
		for {
			select {
			case <-ctx.Done():
				stringFragmentChannel <- sz.StringFragment{
					Error: ctx.Err(),
				}
				break forLoop
			default:
				entityReportFragment, err := client.FetchNext(ctx, reportHandle)
				if err != nil {
					stringFragmentChannel <- sz.StringFragment{
						Error: err,
					}
					break forLoop
				}
				if len(entityReportFragment) == 0 {
					break forLoop
				}
				stringFragmentChannel <- sz.StringFragment{
					Value: entityReportFragment,
				}
			}
		}
		if client.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8079, err, details)
			}()
		}
	}()
	return stringFragmentChannel
}

/*
The ExportJsonEntityReport method initializes a cursor over a document of exported entities.
It is part of the ExportJsonEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A handle that identifies the document to be scrolled through using FetchNext().
*/
func (client *Szengine) ExportJsonEntityReport(ctx context.Context, flags int64) (uintptr, error) {
	//  _DLEXPORT int G2_exportJSONEntityReport(const long long flags, ExportHandle* responseHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultExportHandle uintptr
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(29, flags)
		defer func() { client.traceExit(30, flags, resultExportHandle, err, time.Since(entryTime)) }()
	}
	result := C.G2_exportJSONEntityReport_helper(C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4013, flags, result.returnCode, result, time.Since(entryTime))
	}
	resultExportHandle = (uintptr)(result.exportHandle)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8014, err, details)
		}()
	}
	return resultExportHandle, err
}

/*
The ExportJsonEntityReportIterator method creates an Iterator that can be used in a for-loop
to scroll through a document of exported entities.
It is a convenience method for the ExportJsonEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportJsonEntityReportIterator(ctx context.Context, flags int64) chan sz.StringFragment {
	stringFragmentChannel := make(chan sz.StringFragment)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(stringFragmentChannel)
		var err error = nil
		if client.isTrace {
			entryTime := time.Now()
			client.traceEntry(165, flags)
			defer func() { client.traceExit(166, flags, err, time.Since(entryTime)) }()
		}
		reportHandle, err := client.ExportJsonEntityReport(ctx, flags)
		if err != nil {
			result := sz.StringFragment{
				Error: err,
			}
			stringFragmentChannel <- result
			return
		}
		defer func() {
			client.CloseExport(ctx, reportHandle)
		}()
	forLoop:
		for {
			select {
			case <-ctx.Done():
				stringFragmentChannel <- sz.StringFragment{
					Error: ctx.Err(),
				}
				break forLoop
			default:
				entityReportFragment, err := client.FetchNext(ctx, reportHandle)
				if err != nil {
					stringFragmentChannel <- sz.StringFragment{
						Error: err,
					}
					break forLoop
				}
				if len(entityReportFragment) == 0 {
					break forLoop
				}
				stringFragmentChannel <- sz.StringFragment{
					Value: entityReportFragment,
				}
			}
		}
		if client.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8080, err, details)
			}()
		}
	}()
	return stringFragmentChannel
}

/*
The FetchNext method is used to scroll through an exported document.
It is part of the ExportJsonEntityReport() or ExportCsvEntityReport(), FetchNext(), CloseExport()
lifecycle of a list of exported entities.

Input
  - ctx: A context to control lifecycle.
  - exportHandle: A handle created by ExportJsonEntityReport() or ExportCsvEntityReport().

Output
  - TODO: Document output for FetchNext
*/
func (client *Szengine) FetchNext(ctx context.Context, exportHandle uintptr) (string, error) {
	//  _DLEXPORT int G2_fetchNext(ExportHandle responseHandle, char *responseBuf, const size_t bufSize);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(31, exportHandle)
		defer func() { client.traceExit(32, exportHandle, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_fetchNext_helper(C.uintptr_t(exportHandle))
	if result.returnCode < 0 {
		err = client.newError(ctx, 4014, exportHandle, result.returnCode, result, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8015, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Document FindInterestingEntitiesByEntityId
The FindInterestingEntitiesByEntityId method...

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) FindInterestingEntitiesByEntityId(ctx context.Context, entityId int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_findInterestingEntitiesByEntityID(const long long entityID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(33, entityId, flags)
		defer func() { client.traceExit(34, entityId, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_findInterestingEntitiesByEntityID_helper(C.longlong(entityId), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4015, entityId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": strconv.FormatInt(entityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8016, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Document FindInterestingEntitiesByRecordId
The FindInterestingEntitiesByRecordId method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) FindInterestingEntitiesByRecordId(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_findInterestingEntitiesByRecordID(const char* dataSourceCode, const char* recordID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(35, dataSourceCode, recordId, flags)
		defer func() {
			client.traceExit(36, dataSourceCode, recordId, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIDForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIDForC))
	result := C.G2_findInterestingEntitiesByRecordID_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4016, dataSourceCode, recordId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8017, err, details)
		}()
	}
	return resultResponse, err
}

/*
The FindNetworkByEntityId method finds all entities surrounding a requested set of entities.
This includes the requested entities, paths between them, and relations to other nearby entities.

Input
  - ctx: A context to control lifecycle.
  - entityList: A JSON document listing entities.
    Example: `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}, {"ENTITY_ID": 3}]}`
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - buildOutDegree: The number of degrees of relationships to show around each search entity.
  - maxEntities: The maximum number of entities to return in the discovered network.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"SEAMAN","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-11-29 22:25:18.997","LAST_SEEN_DT":"2022-11-29 22:25:19.005"}],"LAST_SEEN_DT":"2022-11-29 22:25:19.005"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-11-29 22:25:19.009","LAST_SEEN_DT":"2022-11-29 22:25:19.009"}],"LAST_SEEN_DT":"2022-11-29 22:25:19.009"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) FindNetworkByEntityId(ctx context.Context, entityList string, maxDegrees int64, buildOutDegree int64, maxEntities int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_findNetworkByEntityID_V2(const char* entityList, const int maxDegree, const int buildOutDegree, const int maxEntities, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(39, entityList, maxDegrees, buildOutDegree, maxDegrees, flags)
		defer func() {
			client.traceExit(40, entityList, maxDegrees, buildOutDegree, maxDegrees, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	entityListForC := C.CString(entityList)
	defer C.free(unsafe.Pointer(entityListForC))
	result := C.G2_findNetworkByEntityID_V2_helper(entityListForC, C.longlong(maxDegrees), C.longlong(buildOutDegree), C.longlong(maxEntities), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4018, entityList, maxDegrees, buildOutDegree, maxEntities, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityList": entityList,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8019, err, details)
		}()
	}
	return resultResponse, err
}

/*
The FindNetworkByRecordId method finds all entities surrounding a requested set of entities identified by record identifiers.
This includes the requested entities, paths between them, and relations to other nearby entities.

Input
  - ctx: A context to control lifecycle.
  - recordList: A JSON document listing entities.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`
  - maxDegree: The maximum number of degrees in paths between search entities.
  - buildOutDegree: The number of degrees of relationships to show around each search entity.
  - maxEntities: The maximum number of entities to return in the discovered network.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:40:34.285","LAST_SEEN_DT":"2022-12-06 14:40:34.420"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.420"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:40:34.359","LAST_SEEN_DT":"2022-12-06 14:40:34.359"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.359"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":3,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:40:34.424","LAST_SEEN_DT":"2022-12-06 14:40:34.424"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.424"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) FindNetworkByRecordId(ctx context.Context, recordList string, maxDegrees int64, buildOutDegree int64, maxEntities int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_findNetworkByRecordID_V2(const char* recordList, const int maxDegree, const int buildOutDegree, const int maxEntities, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(43, recordList, maxDegrees, buildOutDegree, maxDegrees, flags)
		defer func() {
			client.traceExit(44, recordList, maxDegrees, buildOutDegree, maxDegrees, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	recordListForC := C.CString(recordList)
	defer C.free(unsafe.Pointer(recordListForC))
	result := C.G2_findNetworkByRecordID_V2_helper(recordListForC, C.longlong(maxDegrees), C.longlong(buildOutDegree), C.longlong(maxEntities), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4020, recordList, maxDegrees, buildOutDegree, maxEntities, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordList": recordList,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8021, err, details)
		}()
	}
	return resultResponse, err
}

/*
The FindPathByEntityId method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:43:49.024","LAST_SEEN_DT":"2022-12-06 14:43:49.164"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.164"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:43:49.104","LAST_SEEN_DT":"2022-12-06 14:43:49.104"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.104"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) FindPathByEntityId(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, exclusions string, requiredDataSources string, flags int64) (string, error) {
	if len(requiredDataSources) > 0 {
		return client.findPathIncludingSourceByEntityId_V2(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	} else if len(exclusions) > 0 {
		return client.findPathExcludingByEntityId_V2(ctx, startEntityId, endEntityId, maxDegrees, exclusions, flags)
	}
	return client.findPathByEntityId_V2(ctx, startEntityId, endEntityId, maxDegrees, flags)
}

/*
The FindPathByRecordId method finds single relationship paths between two entities.
The entities are identified by starting and ending records.
Paths are found using known relationships with other entities.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:48:19.522","LAST_SEEN_DT":"2022-12-06 14:48:19.667"}],"LAST_SEEN_DT":"2022-12-06 14:48:19.667"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:48:19.593","LAST_SEEN_DT":"2022-12-06 14:48:19.593"}],"LAST_SEEN_DT":"2022-12-06 14:48:19.593"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) FindPathByRecordId(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, exclusions string, requiredDataSources string, flags int64) (string, error) {
	if len(requiredDataSources) > 0 {
		return client.findPathIncludingSourceByRecordId_V2(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags)
	} else if len(exclusions) > 0 {
		return client.findPathExcludingByRecordId_V2(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, flags)
	}
	return client.findPathByRecordId_V2(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, flags)
}

/*
The GetActiveConfigId method returns the identifier of the loaded Senzing engine configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - The identifier of the active Senzing Engine configuration.
*/
func (client *Szengine) GetActiveConfigId(ctx context.Context) (int64, error) {
	//  _DLEXPORT int G2_getActiveConfigID(long long* configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultConfigId int64
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(69)
		defer func() { client.traceExit(70, resultConfigId, err, time.Since(entryTime)) }()
	}
	result := C.G2_getActiveConfigID_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4033, result.returnCode, result, time.Since(entryTime))
	}
	resultConfigId = int64(C.longlong(result.configID))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8034, err, details)
		}()
	}
	return resultConfigId, err
}

/*
The GetEntityByEntityId method returns entity data based on the ID of a resolved identity.

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
    Example: `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:09:48.577","LAST_SEEN_DT":"2022-12-06 15:09:48.705"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.705","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:09:48.577"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+EXACTLY_SAME","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:09:48.705"}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:09:48.647","LAST_SEEN_DT":"2022-12-06 15:09:48.647"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.647"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:09:48.709","LAST_SEEN_DT":"2022-12-06 15:09:48.709"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.709"}]}`
*/
func (client *Szengine) GetEntityByEntityId(ctx context.Context, entityId int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_getEntityByEntityID_V2(const long long entityID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(73, entityId, flags)
		defer func() { client.traceExit(74, entityId, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_getEntityByEntityID_V2_helper(C.longlong(entityId), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4035, entityId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId": strconv.FormatInt(entityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8036, err, details)
		}()
	}
	return resultResponse, err
}

/*
The GetEntityByRecordId method returns entity data based on the ID of a record which is a member of the entity.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:12:25.464","LAST_SEEN_DT":"2022-12-06 15:12:25.597"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.597","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:12:25.464"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+EXACTLY_SAME","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:12:25.597"}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:12:25.536","LAST_SEEN_DT":"2022-12-06 15:12:25.536"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.536"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:12:25.603","LAST_SEEN_DT":"2022-12-06 15:12:25.603"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.603"}]}`
*/
func (client *Szengine) GetEntityByRecordId(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_getEntityByRecordID_V2(const char* dataSourceCode, const char* recordID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(77, dataSourceCode, recordId, flags)
		defer func() {
			client.traceExit(78, dataSourceCode, recordId, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_getEntityByRecordID_V2_helper(dataSourceCodeForC, recordIdForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4037, dataSourceCode, recordId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8038, err, details)
		}()
	}
	return resultResponse, err
}

/*
The GetRecord method returns a JSON document of a single record from the Senzing repository.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) GetRecord(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_getRecord_V2(const char* dataSourceCode, const char* recordID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(85, dataSourceCode, recordId, flags)
		defer func() {
			client.traceExit(86, dataSourceCode, recordId, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_getRecord_V2_helper(dataSourceCodeForC, recordIdForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4040, dataSourceCode, recordId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8040, err, details)
		}()
	}
	return resultResponse, err
}

/*
The GetRedoRecord method returns the next internally queued maintenance record from the Senzing repository.
Usually, the ProcessRedoRecord() or ProcessRedoRecordWithInfo() method is called to process the maintenance record
retrieved by GetRedoRecord().

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) GetRedoRecord(ctx context.Context) (string, error) {
	//  _DLEXPORT int G2_getRedoRecord(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(87)
		defer func() { client.traceExit(88, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_getRedoRecord_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4041, result.returnCode, result, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8041, err, details)
		}()
	}
	return resultResponse, err
}

/*
The GetStats method retrieves workload statistics for the current process.
These statistics will automatically reset after retrieval.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
    Example: `{"workload":{"loadedRecords":5,"addedRecords":2,"deletedRecords":0,"reevaluations":0,"repairedEntities":0,"duration":56,"retries":0,"candidates":19,"actualAmbiguousTest":0,"cachedAmbiguousTest":0,"libFeatCacheHit":219,"libFeatCacheMiss":73,"unresolveTest":1,"abortedUnresolve":0,"gnrScorersUsed":1,"unresolveTriggers":{"normalResolve":0,"update":0,"relLink":0,"extensiveResolve":0,"ambiguousNoResolve":1,"ambiguousMultiResolve":0},"reresolveTriggers":{"abortRetry":0,"unresolveMovement":0,"multipleResolvableCandidates":0,"resolveNewFeatures":1,"newFeatureFTypes":[{"DOB":1}]},"reresolveSkipped":0,"filteredObsFeat":0,"expressedFeatureCalls":[{"EFCALL_ID":1,"EFUNC_CODE":"PHONE_HASHER","numCalls":1},{"EFCALL_ID":2,"EFUNC_CODE":"EXPRESS_ID","numCalls":1},{"EFCALL_ID":3,"EFUNC_CODE":"EXPRESS_ID","numCalls":1},{"EFCALL_ID":5,"EFUNC_CODE":"EXPRESS_BOM","numCalls":1},{"EFCALL_ID":7,"EFUNC_CODE":"NAME_HASHER","numCalls":4},{"EFCALL_ID":9,"EFUNC_CODE":"ADDR_HASHER","numCalls":1},{"EFCALL_ID":10,"EFUNC_CODE":"EXPRESS_BOM","numCalls":1},{"EFCALL_ID":14,"EFUNC_CODE":"EXPRESS_ID","numCalls":1},{"EFCALL_ID":16,"EFUNC_CODE":"EXPRESS_ID","numCalls":4}],"expressedFeaturesCreated":[{"ADDR_KEY":2},{"ID_KEY":7},{"NAME_KEY":14},{"PHONE_KEY":1},{"SEARCH_KEY":2}],"scoredPairs":[{"ACCT_NUM":16},{"ADDRESS":16},{"DOB":25},{"GENDER":16},{"LOGIN_ID":16},{"NAME":19},{"PHONE":16},{"SSN":19}],"cacheHit":[{"ADDRESS":12},{"DOB":18},{"NAME":13},{"PHONE":15}],"cacheMiss":[{"ADDRESS":4},{"DOB":7},{"NAME":6},{"PHONE":1}],"redoTriggers":[],"latchContention":[],"highContentionFeat":[],"highContentionResEnt":[],"genericDetect":[],"candidateBuilders":[{"ACCT_NUM":7},{"ADDR_KEY":7},{"DOB":7},{"ID_KEY":9},{"LOGIN_ID":7},{"NAME_KEY":9},{"PHONE":7},{"PHONE_KEY":7},{"SEARCH_KEY":7},{"SSN":9}],"suppressedCandidateBuilders":[],"suppressedScoredFeatureType":[],"reducedScoredFeatureType":[],"suppressedDisclosedRelationshipDomainCount":0,"CorruptEntityTestDiagnosis":{},"threadState":{"active":0,"idle":4,"sqlExecuting":0,"loader":0,"resolver":0,"scoring":0,"dataLatchContention":0,"obsEntContention":0,"resEntContention":0},"systemResources":{"initResources":[{"physicalCores":16},{"logicalCores":16},{"totalMemory":"62.6GB"},{"availableMemory":"49.5GB"}],"currResources":[{"availableMemory":"47.4GB"},{"activeThreads":0},{"workerThreads":4},{"systemLoad":[{"cpuUser":13.442277},{"cpuSystem":2.635741},{"cpuIdle":82.024246},{"cpuWait":1.634159},{"cpuSoftIrq":0.263574}]}]}}}`
*/
func (client *Szengine) GetStats(ctx context.Context) (string, error) {
	//  _DLEXPORT int G2_stats(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(139)
		defer func() { client.traceExit(140, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_stats_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4066, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8066, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Write description for GetVirtualEntityByRecordId
The GetVirtualEntityByRecordId method...

Input
  - ctx: A context to control lifecycle.
  - recordList: A JSON document.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:20:17.088","LAST_SEEN_DT":"2022-12-06 15:20:17.161"}],"LAST_SEEN_DT":"2022-12-06 15:20:17.161","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","LAST_SEEN_DT":"2022-12-06 15:20:17.088","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","LAST_SEEN_DT":"2022-12-06 15:20:17.161","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]}}`
*/
func (client *Szengine) GetVirtualEntityByRecordId(ctx context.Context, recordList string, flags int64) (string, error) {
	//  _DLEXPORT int G2_getVirtualEntityByRecordID_V2(const char* recordList, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(93, recordList, flags)
		defer func() { client.traceExit(94, recordList, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	recordListForC := C.CString(recordList)
	defer C.free(unsafe.Pointer(recordListForC))
	result := C.G2_getVirtualEntityByRecordID_V2_helper(recordListForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4044, recordList, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordList": recordList}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8043, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Write description for HowEntityByEntityId
The HowEntityByEntityId method...

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) HowEntityByEntityId(ctx context.Context, entityId int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_howEntityByEntityID_V2(const long long entityID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(97, entityId, flags)
		defer func() { client.traceExit(98, entityId, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_howEntityByEntityID_V2_helper(C.longlong(entityId), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4046, entityId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId": strconv.FormatInt(entityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8046, err, details)
		}()
	}
	return resultResponse, err
}

/*
The PrimeEngine method pre-initializes some of the heavier weight internal resources of the G2 engine.
The G2 Engine uses "lazy initialization".
PrimeEngine() forces initialization.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) PrimeEngine(ctx context.Context) error {
	//  _DLEXPORT int G2_primeEngine();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(103)
		defer func() { client.traceExit(104, err, time.Since(entryTime)) }()
	}
	result := C.G2_primeEngine()
	if result != 0 {
		err = client.newError(ctx, 4049, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8049, err, details)
		}()
	}
	return err
}

/*
The ProcessRedoRecord method processes the next redo record and returns it.
Calling ProcessRedoRecord() has the potential to create more redo records in certain situations.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) ProcessRedoRecord(ctx context.Context, redoRecord string, flags int64) (string, error) {

	if (flags & sz.SZ_WITH_INFO) > 0 {
		finalFlags := flags ^ sz.SZ_WITH_INFO
		return client.processRedoRecordWithInfo(ctx, finalFlags)
	}
	return client.processRedoRecord(ctx, flags)
}

/*
TODO: Write description for ReevaluateEntity
The ReevaluateEntity method...

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateEntity(ctx context.Context, entityId int64, flags int64) (string, error) {

	if (flags & sz.SZ_WITH_INFO) > 0 {
		finalFlags := flags ^ sz.SZ_WITH_INFO
		return client.reevaluateEntityWithInfo(ctx, entityId, finalFlags)
	}
	return client.reevaluateEntity(ctx, entityId, flags)
}

/*
TODO: Write description for ReevaluateRecord
The ReevaluateRecord method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateRecord(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	if (flags & sz.SZ_WITH_INFO) > 0 {
		finalFlags := flags ^ sz.SZ_WITH_INFO
		return client.reevaluateRecordWithInfo(ctx, dataSourceCode, recordId, finalFlags)
	}
	return client.reevaluateRecord(ctx, dataSourceCode, recordId, flags)
}

/*
The Reinitialize method re-initializes the Senzing G2Engine object using a specified configuration identifier.

Input
  - ctx: A context to control lifecycle.
  - configId: The configuration ID used for the initialization.
*/
func (client *Szengine) Reinitialize(ctx context.Context, configId int64) error {
	//  _DLEXPORT int G2_reinit(const long long initConfigID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(127, configId)
		defer func() { client.traceExit(128, configId, err, time.Since(entryTime)) }()
	}
	result := C.G2_reinit(C.longlong(configId))
	if result != 0 {
		err = client.newError(ctx, 4061, configId, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configId": strconv.FormatInt(configId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8061, err, details)
		}()
	}
	return err
}

/*
The SearchByAttributes method retrieves entity data based on a user-specified set of entity attributes.

TODO: Write description of Input for SearchByAttributes
Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the record to be added to the Senzing repository.
  - searchProfile: TODO:
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","MATCH_KEY":"+NAME+SSN","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"NAME":[{"INBOUND_FEAT":"JOHNSON","CANDIDATE_FEAT":"JOHNSON","GNR_FN":100,"GNR_SN":100,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1}],"SSN":[{"INBOUND_FEAT":"053-39-3251","CANDIDATE_FEAT":"053-39-3251","FULL_SCORE":100}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:38:06.175","LAST_SEEN_DT":"2022-12-06 15:38:06.957"}],"LAST_SEEN_DT":"2022-12-06 15:38:06.957"}}}]}`
*/
func (client *Szengine) SearchByAttributes(ctx context.Context, attributes string, searchProfile string, flags int64) (string, error) {
	if len(searchProfile) > 0 {
		return client.searchByAttributes_V3(ctx, attributes, searchProfile, flags)
	}
	return client.searchByAttributes_V2(ctx, attributes, flags)
}

/*
The WhyEntities method explains why records belong to their resolved entities.
WhyEntities() will compare the record data within an entity
against the rest of the entity data and show why they are connected.
This is calculated based on the features that record data represents.

Input
  - ctx: A context to control lifecycle.
  - entityId1: The entity ID for the starting entity of the search path.
  - entityId2: The entity ID for the ending entity of the search path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"},{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"4/8/1983","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":86,"SCORE_BUCKET":"PLAUSIBLE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.906","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.400","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.404","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.407","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.410","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.259","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.201","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]}]}`
*/
func (client *Szengine) WhyEntities(ctx context.Context, entityId1 int64, entityId2 int64, flags int64) (string, error) {
	return client.whyEntities_V2(ctx, entityId1, entityId2, flags)
}

/*
TODO: Write description for WhyRecordInEntity
TODO: Assign entry/exit trace IDs
The WhyRecordInEntity method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) WhyRecordInEntity(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	return client.whyRecordInEntity_V2(ctx, dataSourceCode, recordId, flags)
}

/*
The WhyRecords method explains why records belong to their resolved entities.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordId1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordId2: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":2,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"TEST","RECORD_ID":"222"}],"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-DOB-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.916","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.405","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.408","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.411","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.418","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.265","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.208","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]}]}`
*/
func (client *Szengine) WhyRecords(ctx context.Context, dataSourceCode1 string, recordId1 string, dataSourceCode2 string, recordId2 string, flags int64) (string, error) {
	return client.whyRecords_V2(ctx, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
The GetObserverOrigin method returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) GetObserverOrigin(ctx context.Context) string {
	return client.observerOrigin
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same SzEngine interface.
For this implementation, "base" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) GetSdkId(ctx context.Context) string {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(161)
		defer func() { client.traceExit(162, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8075, err, details)
		}()
	}
	return "base"
}

/*
The Initialize method initializes the SzEngine object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configId: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) Initialize(ctx context.Context, instanceName string, settings string, configId int64, verboseLogging int64) error {
	if configId > 0 {
		return client.initializeWithConfigId(ctx, instanceName, settings, configId, verboseLogging)
	}
	return client.initialize(ctx, instanceName, settings, verboseLogging)
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(157, observer.GetObserverId(ctx))
		defer func() { client.traceExit(158, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	err = client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerId": observer.GetObserverId(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8076, err, details)
		}()
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szengine) SetLogLevel(ctx context.Context, logLevelName string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(137, logLevelName)
		defer func() { client.traceExit(138, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8077, err, details)
		}()
	}
	return err
}

/*
The SetObserverOrigin method sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) SetObserverOrigin(ctx context.Context, origin string) {
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(159, observer.GetObserverId(ctx))
		defer func() { client.traceExit(160, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerId": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8078, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}

// ----------------------------------------------------------------------------
// Private, delegated methods for interface methods
// ----------------------------------------------------------------------------

/*
The addRecord method adds a record into the Senzing repository.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing repository.
*/
func (client *Szengine) addRecord(ctx context.Context, dataSourceCode string, recordId string, recordDefinition string) (string, error) {
	//  _DLEXPORT int G2_addRecord(const char* dataSourceCode, const char* recordID, const char* jsonData, const char *loadID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(1, dataSourceCode, recordId, recordDefinition)
		defer func() {
			client.traceExit(2, dataSourceCode, recordId, recordDefinition, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	jsonDataForC := C.CString(recordDefinition)
	defer C.free(unsafe.Pointer(jsonDataForC))
	result := C.G2_addRecord(dataSourceCodeForC, recordIdForC, jsonDataForC)
	if result != 0 {
		err = client.newError(ctx, 4001, dataSourceCode, recordId, recordDefinition, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return "{}", err
}

/*
The addRecordWithInfo method adds a record into the Senzing repository and returns information on the affected entities.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing repository.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) addRecordWithInfo(ctx context.Context, dataSourceCode string, recordId string, recordDefinition string, flags int64) (string, error) {
	//  _DLEXPORT int G2_addRecordWithInfo(const char* dataSourceCode, const char* recordID, const char* jsonData, const char *loadID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(3, dataSourceCode, recordId, recordDefinition, flags)
		defer func() {
			client.traceExit(4, dataSourceCode, recordId, recordDefinition, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	jsonDataForC := C.CString(recordDefinition)
	defer C.free(unsafe.Pointer(jsonDataForC))
	result := C.G2_addRecordWithInfo_helper(dataSourceCodeForC, recordIdForC, jsonDataForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4002, dataSourceCode, recordId, recordDefinition, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8002, err, details)
		}()
	}
	return resultResponse, err
}

/*
The deleteRecord method deletes a record from the Senzing repository.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
*/
func (client *Szengine) deleteRecord(ctx context.Context, dataSourceCode string, recordId string) (string, error) {
	//  _DLEXPORT int G2_deleteRecord(const char* dataSourceCode, const char* recordID, const char* loadID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(17, dataSourceCode, recordId)
		defer func() { client.traceExit(18, dataSourceCode, recordId, err, time.Since(entryTime)) }()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_deleteRecord(dataSourceCodeForC, recordIdForC)
	if result != 0 {
		err = client.newError(ctx, 4007, dataSourceCode, recordId, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8008, err, details)
		}()
	}
	return "{}", err
}

/*
The deleteRecordWithInfo method deletes a record from the Senzing repository and returns information on the affected entities.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) deleteRecordWithInfo(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_deleteRecordWithInfo(const char* dataSourceCode, const char* recordID, const char* loadID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(19, dataSourceCode, recordId, flags)
		defer func() {
			client.traceExit(20, dataSourceCode, recordId, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_deleteRecordWithInfo_helper(dataSourceCodeForC, recordIdForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4008, dataSourceCode, recordId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8009, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathByEntityId_V2 method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:43:49.024","LAST_SEEN_DT":"2022-12-06 14:43:49.164"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.164"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:43:49.104","LAST_SEEN_DT":"2022-12-06 14:43:49.104"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.104"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) findPathByEntityId_V2(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathByEntityID_V2(const long long entityID1, const long long entityID2, const int maxDegree, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(47, startEntityId, endEntityId, maxDegrees, flags)
		defer func() {
			client.traceExit(48, startEntityId, endEntityId, maxDegrees, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	result := C.G2_findPathByEntityID_V2_helper(C.longlong(startEntityId), C.longlong(endEntityId), C.longlong(maxDegrees), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4022, startEntityId, endEntityId, maxDegrees, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityId": strconv.FormatInt(startEntityId, 10),
				"endEntityId":   strconv.FormatInt(endEntityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8023, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathByRecordId_V2 method finds single relationship paths between two entities.
The entities are identified by starting and ending records.
Paths are found using known relationships with other entities.
It extends FindPathByRecordId() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) findPathByRecordId_V2(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathByRecordID_V2(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const int maxDegree, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(51, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, flags)
		defer func() {
			client.traceExit(52, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	startDataSourceCodeForC := C.CString(startDataSourceCode)
	defer C.free(unsafe.Pointer(startDataSourceCodeForC))
	startRecordIdForC := C.CString(startRecordId)
	defer C.free(unsafe.Pointer(startRecordIdForC))
	endDataSourceCodeForC := C.CString(endDataSourceCode)
	defer C.free(unsafe.Pointer(endDataSourceCodeForC))
	endRecordIdForC := C.CString(endRecordId)
	defer C.free(unsafe.Pointer(endRecordIdForC))
	result := C.G2_findPathByRecordID_V2_helper(startDataSourceCodeForC, startRecordIdForC, endDataSourceCodeForC, endRecordIdForC, C.longlong(maxDegrees), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4024, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordId":       startRecordId,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordId":         endRecordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8025, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathExcludingByEntityId method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that exclude certain entities from being on the path.
To control output, use FindPathExcludingByEntityId_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:50:49.222","LAST_SEEN_DT":"2022-12-06 14:50:49.356"}],"LAST_SEEN_DT":"2022-12-06 14:50:49.356"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:50:49.295","LAST_SEEN_DT":"2022-12-06 14:50:49.295"}],"LAST_SEEN_DT":"2022-12-06 14:50:49.295"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) findPathExcludingByEntityId(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, exclusions string) (string, error) {
	//  _DLEXPORT int G2_findPathExcludingByEntityID(const long long entityID1, const long long entityID2, const int maxDegree, const char* excludedEntities, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(53, startEntityId, endEntityId, maxDegrees, exclusions)
		defer func() {
			client.traceExit(54, startEntityId, endEntityId, maxDegrees, exclusions, resultResponse, err, time.Since(entryTime))
		}()
	}
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	result := C.G2_findPathExcludingByEntityID_helper(C.longlong(startEntityId), C.longlong(endEntityId), C.longlong(maxDegrees), exclusionsForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4025, startEntityId, endEntityId, maxDegrees, exclusions, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityId": strconv.FormatInt(startEntityId, 10),
				"endEntityId":   strconv.FormatInt(endEntityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8026, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathExcludingByEntityId_V2 method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that exclude certain entities from being on the path.
It extends FindPathExcludingByEntityId() by adding output control flags.

When excluding entities, the user may choose to either strictly exclude the entities,
or prefer to exclude the entities but still include them if no other path is found.
By default, entities will be strictly excluded.
A "preferred exclude" may be done by specifying the G2_FIND_PATH_PREFER_EXCLUDE control flag.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) findPathExcludingByEntityId_V2(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, exclusions string, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathExcludingByEntityID_V2(const long long entityID1, const long long entityID2, const int maxDegree, const char* excludedEntities, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(55, startEntityId, endEntityId, maxDegrees, exclusions, flags)
		defer func() {
			client.traceExit(56, startEntityId, endEntityId, maxDegrees, exclusions, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	result := C.G2_findPathExcludingByEntityID_V2_helper(C.longlong(startEntityId), C.longlong(endEntityId), C.longlong(maxDegrees), exclusionsForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4026, startEntityId, endEntityId, maxDegrees, exclusions, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityId": strconv.FormatInt(startEntityId, 10),
				"endEntityId":   strconv.FormatInt(endEntityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8027, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathExcludingByRecordId method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that exclude certain entities from being on the path.
To control output, use FindPathExcludingByRecordId_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:55:02.577","LAST_SEEN_DT":"2022-12-06 14:55:02.711"}],"LAST_SEEN_DT":"2022-12-06 14:55:02.711"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:55:02.649","LAST_SEEN_DT":"2022-12-06 14:55:02.649"}],"LAST_SEEN_DT":"2022-12-06 14:55:02.649"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) findPathExcludingByRecordId(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, exclusions string) (string, error) {
	//  _DLEXPORT int G2_findPathExcludingByRecordID(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const int maxDegree, const char* excludedRecords, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(57, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions)
		defer func() {
			client.traceExit(58, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, resultResponse, err, time.Since(entryTime))
		}()
	}
	startDataSourceCodeForC := C.CString(startDataSourceCode)
	defer C.free(unsafe.Pointer(startDataSourceCodeForC))
	startRecordIdForC := C.CString(startRecordId)
	defer C.free(unsafe.Pointer(startRecordIdForC))
	endDataSourceCodeForC := C.CString(endDataSourceCode)
	defer C.free(unsafe.Pointer(endDataSourceCodeForC))
	endRecordIdForC := C.CString(endRecordId)
	defer C.free(unsafe.Pointer(endRecordIdForC))
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	result := C.G2_findPathExcludingByRecordID_helper(startDataSourceCodeForC, startRecordIdForC, endDataSourceCodeForC, endRecordIdForC, C.longlong(maxDegrees), exclusionsForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4027, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordId":       startRecordId,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordId":         endRecordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8028, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathExcludingByRecordId_V2 method finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that exclude certain entities from being on the path.
It extends FindPathExcludingByRecordId() by adding output control flags.

When excluding entities, the user may choose to either strictly exclude the entities,
or prefer to exclude the entities but still include them if no other path is found.
By default, entities will be strictly excluded.
A "preferred exclude" may be done by specifying the G2_FIND_PATH_PREFER_EXCLUDE control flag.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) findPathExcludingByRecordId_V2(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, exclusions string, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathExcludingByRecordID_V2(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const int maxDegree, const char* excludedRecords, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(59, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, flags)
		defer func() {
			client.traceExit(60, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	startDataSourceCodeForC := C.CString(startDataSourceCode)
	defer C.free(unsafe.Pointer(startDataSourceCodeForC))
	startRecordIdForC := C.CString(startRecordId)
	defer C.free(unsafe.Pointer(startRecordIdForC))
	endDataSourceCodeForC := C.CString(endDataSourceCode)
	defer C.free(unsafe.Pointer(endDataSourceCodeForC))
	endRecordIdForC := C.CString(endRecordId)
	defer C.free(unsafe.Pointer(endRecordIdForC))
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	result := C.G2_findPathExcludingByRecordID_V2_helper(startDataSourceCodeForC, startRecordIdForC, endDataSourceCodeForC, endRecordIdForC, C.longlong(maxDegrees), exclusionsForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4028, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordId":       startRecordId,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordId":         endRecordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8029, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathIncludingSourceByEntityId method finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityIdD() and FindPathExcludingByRecordId().
To control output, use FindPathIncludingSourceByEntityId_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:00:30.268","LAST_SEEN_DT":"2022-12-06 15:00:30.429"}],"LAST_SEEN_DT":"2022-12-06 15:00:30.429"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:00:30.339","LAST_SEEN_DT":"2022-12-06 15:00:30.339"}],"LAST_SEEN_DT":"2022-12-06 15:00:30.339"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) findPathIncludingSourceByEntityId(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, exclusions string, requiredDataSources string) (string, error) {
	//  _DLEXPORT int G2_findPathIncludingSourceByEntityID(const long long entityID1, const long long entityID2, const int maxDegree, const char* excludedEntities, const char* requiredDsrcs, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(61, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources)
		defer func() {
			client.traceExit(62, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, resultResponse, err, time.Since(entryTime))
		}()
	}
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	requiredDataSourcesForC := C.CString(requiredDataSources)
	defer C.free(unsafe.Pointer(requiredDataSourcesForC))
	result := C.G2_findPathIncludingSourceByEntityID_helper(C.longlong(startEntityId), C.longlong(endEntityId), C.longlong(maxDegrees), exclusionsForC, requiredDataSourcesForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4029, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityId": strconv.FormatInt(startEntityId, 10),
				"endEntityId":   strconv.FormatInt(endEntityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8030, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathIncludingSourceByEntityId_V2 method finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityId_V2() and FindPathExcludingByRecordId_V2().
It extends FindPathIncludingSourceByEntityId() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - startEntityId: The entity ID for the starting entity of the search path.
  - endEntityId: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) findPathIncludingSourceByEntityId_V2(ctx context.Context, startEntityId int64, endEntityId int64, maxDegrees int64, exclusions string, requiredDataSources string, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathIncludingSourceByEntityID_V2(const long long entityID1, const long long entityID2, const int maxDegree, const char* excludedEntities, const char* requiredDsrcs, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(63, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
		defer func() {
			client.traceExit(64, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	requiredDataSourcesForC := C.CString(requiredDataSources)
	defer C.free(unsafe.Pointer(requiredDataSourcesForC))
	result := C.G2_findPathIncludingSourceByEntityID_V2_helper(C.longlong(startEntityId), C.longlong(endEntityId), C.longlong(maxDegrees), exclusionsForC, requiredDataSourcesForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4030, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityId": strconv.FormatInt(startEntityId, 10),
				"endEntityId":   strconv.FormatInt(endEntityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8031, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathIncludingSourceByRecordId method finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityId() and FindPathExcludingByRecordId().
To control output, use FindPathIncludingSourceByRecordId_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - endRecordId: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.

Output
  - A JSON document.
    Example: `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:03:52.805","LAST_SEEN_DT":"2022-12-06 15:03:52.947"}],"LAST_SEEN_DT":"2022-12-06 15:03:52.947"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:03:52.876","LAST_SEEN_DT":"2022-12-06 15:03:52.876"}],"LAST_SEEN_DT":"2022-12-06 15:03:52.876"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
*/
func (client *Szengine) findPathIncludingSourceByRecordId(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, exclusions string, requiredDataSources string) (string, error) {
	//  _DLEXPORT int G2_findPathIncludingSourceByRecordID(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const int maxDegree, const char* excludedRecords, const char* requiredDsrcs, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(65, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources)
		defer func() {
			client.traceExit(66, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, resultResponse, err, time.Since(entryTime))
		}()
	}
	startDataSourceCodeForC := C.CString(startDataSourceCode)
	defer C.free(unsafe.Pointer(startDataSourceCodeForC))
	startRecordIdForC := C.CString(startRecordId)
	defer C.free(unsafe.Pointer(startRecordIdForC))
	endDataSourceCodeForC := C.CString(endDataSourceCode)
	defer C.free(unsafe.Pointer(endDataSourceCodeForC))
	endRecordIdForC := C.CString(endRecordId)
	defer C.free(unsafe.Pointer(endRecordIdForC))
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	requiredDataSourcesForC := C.CString(requiredDataSources)
	defer C.free(unsafe.Pointer(requiredDataSourcesForC))
	result := C.G2_findPathIncludingSourceByRecordID_helper(startDataSourceCodeForC, startRecordIdForC, endDataSourceCodeForC, endRecordIdForC, C.longlong(maxDegrees), exclusionsForC, requiredDataSourcesForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4031, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordId":       startRecordId,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordId":         endRecordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8032, err, details)
		}()
	}
	return resultResponse, err
}

/*
The findPathIncludingSourceByRecordId_V2 method finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityId_V2() and FindPathExcludingByRecordId_V2().
It extends FindPathIncludingSourceByRecordId() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordId: The unique identifier within the records of the same data source for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordId: The unique identifier within the records of the same data source for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - exclusions: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) findPathIncludingSourceByRecordId_V2(ctx context.Context, startDataSourceCode string, startRecordId string, endDataSourceCode string, endRecordId string, maxDegrees int64, exclusions string, requiredDataSources string, flags int64) (string, error) {
	//  _DLEXPORT int G2_findPathIncludingSourceByRecordID_V2(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const int maxDegree, const char* excludedRecords, const char* requiredDsrcs, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(67, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags)
		defer func() {
			client.traceExit(68, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	startDataSourceCodeForC := C.CString(startDataSourceCode)
	defer C.free(unsafe.Pointer(startDataSourceCodeForC))
	startRecordIdForC := C.CString(startRecordId)
	defer C.free(unsafe.Pointer(startRecordIdForC))
	endDataSourceCodeForC := C.CString(endDataSourceCode)
	defer C.free(unsafe.Pointer(endDataSourceCodeForC))
	endRecordIdForC := C.CString(endRecordId)
	defer C.free(unsafe.Pointer(endRecordIdForC))
	exclusionsForC := C.CString(exclusions)
	defer C.free(unsafe.Pointer(exclusionsForC))
	requiredDataSourcesForC := C.CString(requiredDataSources)
	defer C.free(unsafe.Pointer(requiredDataSourcesForC))
	result := C.G2_findPathIncludingSourceByRecordID_V2_helper(startDataSourceCodeForC, startRecordIdForC, endDataSourceCodeForC, endRecordIdForC, C.longlong(maxDegrees), exclusionsForC, requiredDataSourcesForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4032, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordId":       startRecordId,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordId":         endRecordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8033, err, details)
		}()
	}
	return resultResponse, err
}

/*
The initialize method initializes the SzEngine object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	// _DLEXPORT int G2_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(99, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(100, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	instanceNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(instanceNameForC))
	settingsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(settingsForC))
	result := C.G2_init(instanceNameForC, settingsForC, C.longlong(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4047, instanceName, settings, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8047, err, details)
		}()
	}
	return err
}

/*
The initializeWithConfigId method initializes the Senzing G2 object with a non-default configuration ID.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configId: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) initializeWithConfigId(ctx context.Context, instanceName string, settings string, configId int64, verboseLogging int64) error {
	//  _DLEXPORT int G2_initWithConfigID(const char *moduleName, const char *iniParams, const long long initConfigID, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(101, instanceName, settings, configId, verboseLogging)
		defer func() {
			client.traceExit(102, instanceName, settings, configId, verboseLogging, err, time.Since(entryTime))
		}()
	}
	instanceNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(instanceNameForC))
	settingsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(settingsForC))
	result := C.G2_initWithConfigID(instanceNameForC, settingsForC, C.longlong(configId), C.longlong(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4048, instanceName, settings, configId, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"settings":       settings,
				"configId":       strconv.FormatInt(configId, 10),
				"instanceName":   instanceName,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8048, err, details)
		}()
	}
	return err
}

// TODO: https://senzing.atlassian.net/browse/GDEV-3771
// TODO: Implement processRedoRecord
func (client *Szengine) processRedoRecord(ctx context.Context, flags int64) (string, error) {
	_ = ctx
	_ = flags
	// //  _DLEXPORT int G2_processRedoRecord(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()
	// var err error = nil
	// var resultResponse string
	// entryTime := time.Now()
	//
	//	if client.isTrace {
	//		client.traceEntry(107)
	//		defer func() { client.traceExit(108, resultResponse, err, time.Since(entryTime)) }()
	//	}
	//
	// result := C.G2_processRedoRecord_helper()
	//
	//	if result.returnCode != 0 {
	//		err = client.newError(ctx, 4051, result.returnCode, result, time.Since(entryTime))
	//	}
	//
	// resultResponse = C.GoString(result.response)
	// C.G2GoHelper_free(unsafe.Pointer(result.response))
	//
	//	if client.observers != nil {
	//		go func() {
	//			details := map[string]string{}
	//			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8051, err, details)
	//		}()
	//	}
	//
	// return resultResponse, err
	return "{}", nil
}

/*
The processRedoRecordWithInfo method processes the next redo record and returns it and affected entities.
Calling processRedoRecordWithInfo() has the potential to create more redo records in certain situations.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A JSON document with the record that was re-done.
  - A JSON document with affected entities.
*/
func (client *Szengine) processRedoRecordWithInfo(ctx context.Context, flags int64) (string, error) {
	_ = ctx
	_ = flags

	// TODO: Implement processRedoRecordWithInfo

	return "", nil
}

/*
TODO: Write description for reevaluateEntity
The reevaluateEntity method...

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) reevaluateEntity(ctx context.Context, entityId int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_reevaluateEntity(const long long entityID, const long long flags);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(119, entityId, flags)
		defer func() { client.traceExit(120, entityId, flags, err, time.Since(entryTime)) }()
	}
	result := C.G2_reevaluateEntity(C.longlong(entityId), C.longlong(flags))
	if result != 0 {
		err = client.newError(ctx, 4057, entityId, flags, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId": strconv.FormatInt(entityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8057, err, details)
		}()
	}
	return "{}", err
}

/*
TODO: Write description for reevaluateEntityWithInfo
The reevaluateEntityWithInfo method...

Input
  - ctx: A context to control lifecycle.
  - entityId: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
    See the example output.
*/
func (client *Szengine) reevaluateEntityWithInfo(ctx context.Context, entityId int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_reevaluateEntityWithInfo(const long long entityID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(121, entityId, flags)
		defer func() { client.traceExit(122, entityId, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_reevaluateEntityWithInfo_helper(C.longlong(entityId), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4058, entityId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId": strconv.FormatInt(entityId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8058, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Write description for reevaluateRecord
The reevaluateRecord method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) reevaluateRecord(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_reevaluateRecord(const char* dataSourceCode, const char* recordID, const long long flags);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(123, dataSourceCode, recordId, flags)
		defer func() { client.traceExit(124, dataSourceCode, recordId, flags, err, time.Since(entryTime)) }()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_reevaluateRecord(dataSourceCodeForC, recordIdForC, C.longlong(flags))
	if result != 0 {
		err = client.newError(ctx, 4059, dataSourceCode, recordId, flags, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8059, err, details)
		}()
	}
	return "{}", err
}

/*
TODO: Write description for reevaluateRecordWithInfo
The reevaluateRecordWithInfo method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
    See the example output.
*/
func (client *Szengine) reevaluateRecordWithInfo(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	//  _DLEXPORT int G2_reevaluateRecordWithInfo(const char* dataSourceCode, const char* recordID, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(125, dataSourceCode, recordId, flags)
		defer func() {
			client.traceExit(126, dataSourceCode, recordId, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSourceCodeForC := C.CString(dataSourceCode)
	defer C.free(unsafe.Pointer(dataSourceCodeForC))
	recordIdForC := C.CString(recordId)
	defer C.free(unsafe.Pointer(recordIdForC))
	result := C.G2_reevaluateRecordWithInfo_helper(dataSourceCodeForC, recordIdForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4060, dataSourceCode, recordId, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordId":       recordId,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8060, err, details)
		}()
	}
	return resultResponse, err
}

/*
The searchByAttributes method retrieves entity data based on a user-specified set of entity attributes.
To control output, use SearchByAttributes_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the record to be added to the Senzing repository.

Output
  - A JSON document.
    Example: `{"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","MATCH_KEY":"+NAME+SSN","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"NAME":[{"INBOUND_FEAT":"JOHNSON","CANDIDATE_FEAT":"JOHNSON","GNR_FN":100,"GNR_SN":100,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1}],"SSN":[{"INBOUND_FEAT":"053-39-3251","CANDIDATE_FEAT":"053-39-3251","FULL_SCORE":100}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:38:06.175","LAST_SEEN_DT":"2022-12-06 15:38:06.957"}],"LAST_SEEN_DT":"2022-12-06 15:38:06.957"}}}]}`
*/
func (client *Szengine) searchByAttributes(ctx context.Context, attributes string) (string, error) {
	//  _DLEXPORT int G2_searchByAttributes(const char* jsonData, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(133, attributes)
		defer func() { client.traceExit(134, attributes, resultResponse, err, time.Since(entryTime)) }()
	}
	jsonDataForC := C.CString(attributes)
	defer C.free(unsafe.Pointer(jsonDataForC))
	result := C.G2_searchByAttributes_helper(jsonDataForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4064, attributes, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8064, err, details)
		}()
	}
	return resultResponse, err
}

/*
The SearchByAttributes_V2 method retrieves entity data based on a user-specified set of entity attributes.
It extends SearchByAttributes() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the record to be added to the Senzing repository.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) searchByAttributes_V2(ctx context.Context, attributes string, flags int64) (string, error) {
	//  _DLEXPORT int G2_searchByAttributes_V2(const char* jsonData, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(135, attributes, flags)
		defer func() { client.traceExit(136, attributes, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	jsonDataForC := C.CString(attributes)
	defer C.free(unsafe.Pointer(jsonDataForC))
	result := C.G2_searchByAttributes_V2_helper(jsonDataForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4065, attributes, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8065, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Write description for searchByAttributes_V3
The searchByAttributes_V3 method...

Input
  - ctx: A context to control lifecycle.
  - jsonData: A JSON document containing the record to be added to the Senzing repository.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) searchByAttributes_V3(ctx context.Context, attributes string, searchProfile string, flags int64) (string, error) {
	var err error = nil
	var resultResponse string
	_ = ctx
	_ = attributes
	_ = searchProfile
	_ = flags

	// TODO: implement searchByAttributes_V3

	return resultResponse, err
}

/*
The whyEntities method explains why records belong to their resolved entities.
whyEntities() will compare the record data within an entity
against the rest of the entity data and show why they are connected.
This is calculated based on the features that record data represents.
To control output, use WhyEntities_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - entityId1: The entity ID for the starting entity of the search path.
  - entityId2: The entity ID for the ending entity of the search path.

Output
  - A JSON document.
    Example: `{"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"},{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"4/8/1983","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":86,"SCORE_BUCKET":"PLAUSIBLE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.906","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.400","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.404","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.407","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.410","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.259","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.201","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]}]}`
*/
func (client *Szengine) whyEntities(ctx context.Context, entityId1 int64, entityId2 int64) (string, error) {
	//  _DLEXPORT int G2_whyEntities(const long long entityID1, const long long entityID2, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(141, entityId1, entityId2)
		defer func() { client.traceExit(142, entityId1, entityId2, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_whyEntities_helper(C.longlong(entityId1), C.longlong(entityId2))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4067, entityId1, entityId2, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId1": strconv.FormatInt(entityId1, 10),
				"entityId2": strconv.FormatInt(entityId2, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8067, err, details)
		}()
	}
	return resultResponse, err
}

/*
The whyEntities_V2 method explains why records belong to their resolved entities.
whyEntities_V2() will compare the record data within an entity
against the rest of the entity data and show why they are connected.
This is calculated based on the features that record data represents.
It extends whyEntities() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - entityId1: The entity ID for the starting entity of the search path.
  - entityId2: The entity ID for the ending entity of the search path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) whyEntities_V2(ctx context.Context, entityId1 int64, entityId2 int64, flags int64) (string, error) {
	//  _DLEXPORT int G2_whyEntities_V2(const long long entityID1, const long long entityID2, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(143, entityId1, entityId2, flags)
		defer func() { client.traceExit(144, entityId1, entityId2, flags, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2_whyEntities_V2_helper(C.longlong(entityId1), C.longlong(entityId2), C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4068, entityId1, entityId2, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityId1": strconv.FormatInt(entityId1, 10),
				"entityId2": strconv.FormatInt(entityId2, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8068, err, details)
		}()
	}
	return resultResponse, err
}

/*
TODO: Write description for whyRecordInEntity
The whyRecordInEntity method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) whyRecordInEntity(ctx context.Context, dataSourceCode string, recordId string) (string, error) {
	var err error = nil
	var resultResponse string

	_ = ctx
	_ = dataSourceCode
	_ = recordId

	// TODO: Implement whyRecordInEntity

	return resultResponse, err
}

/*
TODO: Write description for whyRecordInEntity_V2
The whyRecordInEntity_V2 method...

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordId: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) whyRecordInEntity_V2(ctx context.Context, dataSourceCode string, recordId string, flags int64) (string, error) {
	var err error = nil
	var resultResponse string

	_ = ctx
	_ = dataSourceCode
	_ = recordId
	_ = flags

	// TODO: Implement whyRecordInEntity_V2

	return resultResponse, err
}

/*
The whyRecords method explains why records belong to their resolved entities.
To control output, use whyRecords_V2() instead.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordId1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordId2: The unique identifier within the records of the same data source.

Output

  - A JSON document.
    Example: `{"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":2,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"TEST","RECORD_ID":"222"}],"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-DOB-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.916","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.405","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.408","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.411","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.418","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.265","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.208","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]}]}`
*/
func (client *Szengine) whyRecords(ctx context.Context, dataSourceCode1 string, recordId1 string, dataSourceCode2 string, recordId2 string) (string, error) {
	//  _DLEXPORT int G2_whyRecords(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(153, dataSourceCode1, recordId1, dataSourceCode2, recordId2)
		defer func() {
			client.traceExit(154, dataSourceCode1, recordId1, dataSourceCode2, recordId2, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSource1CodeForC := C.CString(dataSourceCode1)
	defer C.free(unsafe.Pointer(dataSource1CodeForC))
	recordId1ForC := C.CString(recordId1)
	defer C.free(unsafe.Pointer(recordId1ForC))
	dataSource2CodeForC := C.CString(dataSourceCode2)
	defer C.free(unsafe.Pointer(dataSource2CodeForC))
	recordId2ForC := C.CString(recordId2)
	defer C.free(unsafe.Pointer(recordId2ForC))
	result := C.G2_whyRecords_helper(dataSource1CodeForC, recordId1ForC, dataSource2CodeForC, recordId2ForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4073, dataSourceCode1, recordId1, dataSourceCode2, recordId2, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode1": dataSourceCode1,
				"recordId1":       recordId1,
				"dataSourceCode2": dataSourceCode2,
				"recordId2":       recordId2,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8073, err, details)
		}()
	}
	return resultResponse, err
}

/*
The whyRecords_V2 method explains why records belong to their resolved entities.
It extends WhyRecords() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordId1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordId2: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
    See the example output.
*/
func (client *Szengine) whyRecords_V2(ctx context.Context, dataSourceCode1 string, recordId1 string, dataSourceCode2 string, recordId2 string, flags int64) (string, error) {
	//  _DLEXPORT int G2_whyRecords_V2(const char* dataSourceCode1, const char* recordID1, const char* dataSourceCode2, const char* recordID2, const long long flags, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(155, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags)
		defer func() {
			client.traceExit(156, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags, resultResponse, err, time.Since(entryTime))
		}()
	}
	dataSource1CodeForC := C.CString(dataSourceCode1)
	defer C.free(unsafe.Pointer(dataSource1CodeForC))
	recordId1ForC := C.CString(recordId1)
	defer C.free(unsafe.Pointer(recordId1ForC))
	dataSource2CodeForC := C.CString(dataSourceCode2)
	defer C.free(unsafe.Pointer(dataSource2CodeForC))
	recordId2ForC := C.CString(recordId2)
	defer C.free(unsafe.Pointer(recordId2ForC))
	result := C.G2_whyRecords_V2_helper(dataSource1CodeForC, recordId1ForC, dataSource2CodeForC, recordId2ForC, C.longlong(flags))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4074, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags, result.returnCode, time.Since(entryTime))
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode1": dataSourceCode1,
				"recordId1":       recordId1,
				"dataSourceCode2": dataSourceCode2,
				"recordId2":       recordId2,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8074, err, details)
		}()
	}
	return resultResponse, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szengine) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, szengineapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *Szengine) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szengine) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

func formatFlags(flags int64) string {
	return strconv.FormatInt(flags, 10)
}

func formatEntityId(entityId int64) string {
	return strconv.FormatInt(entityId, 10)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szengine) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
	lastException, err := client.getLastException(ctx)
	defer client.clearLastException(ctx)
	message := lastException
	if err != nil {
		message = err.Error()
	}
	details = append(details, errors.New(message))
	errorMessage := client.getLogger().Json(errorNumber, details...)
	return szerror.SzError(szerror.SzErrorCode(message), (errorMessage))
}

// --- G2 exception handling --------------------------------------------------

/*
The clearLastException method erases the last exception message held by the Senzing G2 object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2_clearLastException();
	_ = ctx
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}
	C.G2_clearLastException()
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's G2.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2Product.
*/
func (client *Szengine) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2_getLastException(char *buffer, const size_t bufSize);
	_ = ctx
	var err error = nil
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(79)
		defer func() { client.traceExit(80, result, err, time.Since(entryTime)) }()
	}
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	// if result == 0 { // "result" is length of exception message.
	// 	err = client.getLogger().Error(4038, result, time.Since(entryTime))
	// }
	result = string(bytes.Trim(stringBuffer, "\x00"))
	return result, err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2Product.
*/
func (client *Szengine) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2_getLastExceptionCode();
	_ = ctx
	var err error = nil
	var result int
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(81)
		defer func() { client.traceExit(82, result, err, time.Since(entryTime)) }()
	}
	result = int(C.G2_getLastExceptionCode())
	return result, err
}

// --- Misc -------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (client *Szengine) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *Szengine) getByteArray(size int) []byte {
	return make([]byte, size)
}
