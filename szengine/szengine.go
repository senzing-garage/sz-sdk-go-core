/*
The [Szengine] implementation of the [senzing.SzEngine] interface
communicates with the Senzing native C binary, libSz.so.
*/
package szengine

/*
#include <stdlib.h>
#include "libSz.h"
#include "szhelpers/SzLang_helpers.h"
#cgo linux CFLAGS: -g -I/opt/senzing/er/sdk/c
#cgo linux LDFLAGS: -L/opt/senzing/er/lib -lSz
*/
import "C"

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szengine"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

/*
Type Szengine struct implements the [senzing.SzEngine] interface
for communicating with the Senzing C binaries.
*/
type Szengine struct {
	isTrace        bool
	logger         logging.Logging
	messenger      messenger.Messenger
	observerOrigin string
	observers      subject.Subject
}

const (
	baseCallerSkip       = 4
	baseTen              = 10
	initialByteArraySize = 65535
	noError              = 0
	withoutInfo          = ""
)

// ----------------------------------------------------------------------------
// sz-sdk-go.SzEngine interface methods
// ----------------------------------------------------------------------------

/*
Method AddRecord adds a record into the Senzing datastore.
The unique identifier of a record is the [dataSourceCode, recordID] compound key.
If the unique identifier does not exist in the Senzing datastore, a new record definition is created in the
Senzing datastore.
If the unique identifier already exists, the new record definition will replace the old record definition.
If the record definition contains JSON keys of `DATA_SOURCE` and/or `RECORD_ID`, they must match the values of `
dataSourceCode` and `recordID`.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) AddRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	recordDefinition string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(1, dataSourceCode, recordID, recordDefinition, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(2, dataSourceCode, recordID, recordDefinition, flags, result, err, time.Since(entryTime))
		}()
	}

	if (flags & senzing.SzWithInfo) == senzing.SzNoFlags {
		result, err = client.addRecord(ctx, dataSourceCode, recordID, recordDefinition)
	} else {
		finalFlags := flags & ^senzing.SzWithInfo
		result, err = client.addRecordWithInfo(ctx, dataSourceCode, recordID, recordDefinition, finalFlags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.AddRecord error: %w", err)
	}

	return result, nil
}

/*
Method CloseExport closes the exported document created by [Szengine.ExportJSONEntityReport] or
[Szengine.ExportCsvEntityReport].
It is part of the ExportXxxEntityReport(), [Szengine.FetchNext], CloseExport lifecycle of a list of entities to export.
CloseExport is idempotent; an exportHandle may be closed multiple times.

Input

  - ctx: A context to control lifecycle.

  - exportHandle: A handle created by [Szengine.ExportJSONEntityReport] or [Szengine.ExportCsvEntityReport]
    that is to be closed.
*/
func (client *Szengine) CloseExport(ctx context.Context, exportHandle uintptr) error {
	var err error

	if client.isTrace {
		client.traceEntry(5, exportHandle)

		entryTime := time.Now()
		defer func() { client.traceExit(6, exportHandle, err, time.Since(entryTime)) }()
	}

	err = client.closeExport(ctx, exportHandle)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.CloseExport error: %w", err)
	}

	return nil
}

/*
Method CountRedoRecords returns the number of records needing re-evaluation.
These are often called "redo records".

Input
  - ctx: A context to control lifecycle.

Output
  - The number of redo records in Senzing's redo queue.
*/
func (client *Szengine) CountRedoRecords(ctx context.Context) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(7)

		entryTime := time.Now()
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}

	result, err = client.countRedoRecords(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}

	if err != nil {
		return 0, fmt.Errorf("szengine.CountRedoRecords error: %w", err)
	}

	return result, nil
}

/*
Method DeleteRecord deletes a record from the Senzing datastore.
The unique identifier of a record is the [dataSourceCode, recordID] compound key.
DeleteRecord() is idempotent.
Multiple calls to delete the same unique identifier will all succeed,
even if the unique identifier is not present in the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) DeleteRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(9, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(10, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	if (flags & senzing.SzWithInfo) == senzing.SzNoFlags {
		result, err = client.deleteRecord(ctx, dataSourceCode, recordID)
	} else {
		finalFlags := flags & ^senzing.SzWithInfo
		result, err = client.deleteRecordWithInfo(ctx, dataSourceCode, recordID, finalFlags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.DeleteRecord error: %w", err)
	}

	return result, nil
}

/*
Method ExportCsvEntityReport initializes a cursor over a CSV document of exported entities.
It is part of the ExportCsvEntityReport, [Szengine.FetchNext], [Szengine.CloseExport] lifecycle
of a list of entities to export.
The first exported line is the CSV header.
Each subsequent line contains metadata for a single entity.

Input

  - ctx: A context to control lifecycle.

  - csvColumnList: Use `*` to request all columns, an empty string to request "standard" columns,
    or a comma-separated list of column names for customized columns.

  - flags: Flags used to control information returned.

Output
  - exportHandle: A handle that identifies the document to be scrolled through using [Szengine.FetchNext].
*/
func (client *Szengine) ExportCsvEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error) {
	var (
		err    error
		result uintptr
	)

	if client.isTrace {
		client.traceEntry(13, csvColumnList, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(14, csvColumnList, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.exportCsvEntityReport(ctx, csvColumnList, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}

	if err != nil {
		return 0, fmt.Errorf("szengine.ExportCsvEntityReport error: %w", err)
	}

	return result, nil
}

/*
Method ExportCsvEntityReportIterator creates an Iterator that can be used in a for-loop
to scroll through a CSV document of exported entities.
It is a convenience method for the [Szenzine.ExportCsvEntityReport], [Szengine.FetchNext], [Szengine.CloseExport]
lifecycle of a list of entities to export.

Input

  - ctx: A context to control lifecycle.

  - csvColumnList: Use `*` to request all columns, an empty string to request "standard" columns,
    or a comma-separated list of column names for customized columns.

  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportCsvEntityReportIterator(
	ctx context.Context,
	csvColumnList string,
	flags int64,
) chan senzing.StringFragment {
	stringFragmentChannel := make(chan senzing.StringFragment)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(stringFragmentChannel)

		var err error

		if client.isTrace {
			client.traceEntry(15, csvColumnList, flags)

			entryTime := time.Now()
			defer func() { client.traceExit(16, csvColumnList, flags, err, time.Since(entryTime)) }()
		}

		reportHandle, err := client.ExportCsvEntityReport(ctx, csvColumnList, flags)
		if err != nil {
			result := senzing.StringFragment{
				Error: err,
			}
			stringFragmentChannel <- result

			return
		}

		defer func() {
			err = client.CloseExport(ctx, reportHandle)
			if err != nil {
				panic(err) // TODO:  Something better than panic(err)?
			}
		}()

		client.fetchNextIntoChannel(ctx, reportHandle, stringFragmentChannel)

		if client.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
			}()
		}
	}()

	return stringFragmentChannel
}

/*
Method ExportJSONEntityReport initializes a cursor over a JSON document of exported entities.
It is part of the ExportJSONEntityReport, [Szengine.FetchNext], [Szengine.CloseExport] lifecycle
of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A handle that identifies the document to be scrolled through using [Szengine.FetchNext].
*/
func (client *Szengine) ExportJSONEntityReport(ctx context.Context, flags int64) (uintptr, error) {
	var (
		err    error
		result uintptr
	)

	if client.isTrace {
		client.traceEntry(17, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(18, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.exportJSONEntityReport(ctx, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}

	if err != nil {
		return 0, fmt.Errorf("szengine.ExportJSONEntityReport error: %w", err)
	}

	return result, nil
}

/*
Method ExportJSONEntityReportIterator creates an Iterator that can be used in a for-loop
to scroll through a JSON document of exported entities.
It is a convenience method for the [Szengine.ExportJSONEntityReport], [Szengine.FetchNext], [Szengine.CloseExport]
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportJSONEntityReportIterator(ctx context.Context, flags int64) chan senzing.StringFragment {
	stringFragmentChannel := make(chan senzing.StringFragment)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(stringFragmentChannel)

		var err error

		if client.isTrace {
			client.traceEntry(19, flags)

			entryTime := time.Now()
			defer func() { client.traceExit(20, flags, err, time.Since(entryTime)) }()
		}

		reportHandle, err := client.ExportJSONEntityReport(ctx, flags)
		if err != nil {
			result := senzing.StringFragment{
				Error: err,
			}
			stringFragmentChannel <- result

			return
		}

		defer func() {
			err = client.CloseExport(ctx, reportHandle)
			if err != nil {
				panic(err) // TODO:  Something better than panic(err)?
			}
		}()

		client.fetchNextIntoChannel(ctx, reportHandle, stringFragmentChannel)

		if client.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
			}()
		}
	}()

	return stringFragmentChannel
}

/*
Method FetchNext is used to scroll through an exported JSON or CSV document.
It is part of the [Szengine.ExportJSONEntityReport] or [Szengine.ExportCsvEntityReport], FetchNext,
[Szengine.CloseExport] lifecycle of a list of exported entities.

Input
  - ctx: A context to control lifecycle.
  - exportHandle: A handle created by [Szengine.ExportJSONEntityReport] or [Szengine.ExportCsvEntityReport].

Output
  - The next chunk of exported data. An empty string signifies end of data.
*/
func (client *Szengine) FetchNext(ctx context.Context, exportHandle uintptr) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(21, exportHandle)

		entryTime := time.Now()
		defer func() { client.traceExit(22, exportHandle, result, err, time.Since(entryTime)) }()
	}

	result, err = client.fetchNext(ctx, exportHandle)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8010, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FetchNext error: %w", err)
	}

	return result, nil
}

/*
Method FindInterestingEntitiesByEntityID is an experimental method.
Not recommended for use.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindInterestingEntitiesByEntityID(
	ctx context.Context,
	entityID int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(23, entityID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(24, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.findInterestingEntitiesByEntityID(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8011, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindInterestingEntitiesByEntityID error: %w", err)
	}

	return result, nil
}

/*
Method FindInterestingEntitiesByRecordID is an experimental method.
Not recommended for use.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindInterestingEntitiesByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(25, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(26, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.findInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8012, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindInterestingEntitiesByRecordID error: %w", err)
	}

	return result, nil
}

/*
Method FindNetworkByEntityID finds a network of entities surrounding a requested set of entities.
This includes the requested entities, paths between them, and relations to other nearby entities.
The size and character of the returned network can be modified by input parameters.

Input

  - ctx: A context to control lifecycle.

  - entityIDs: A JSON document listing entities.
    Example: `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}, {"ENTITY_ID": 3}]}`

  - maxDegrees: The maximum number of degrees in paths between entityIDs.

  - buildOutDegrees: The number of degrees of relationships to show around each search entity. Zero (0)
    prevents buildout.

  - buildOutMaxEntities: The maximum number of entities to build out in the returned network.

  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindNetworkByEntityID(
	ctx context.Context,
	entityIDs string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(27, entityIDs, maxDegrees, buildOutDegrees, buildOutMaxEntities, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(
				28,
				entityIDs,
				maxDegrees,
				buildOutDegrees,
				buildOutMaxEntities,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.findNetworkByEntityIDV2(
		ctx,
		entityIDs,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityIDs": entityIDs,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8013, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindNetworkByEntityID error: %w", err)
	}

	return result, nil
}

/*
Method FindNetworkByRecordID finds a network of entities surrounding a requested set of entities identified
by record keys.
This includes the requested entities, paths between them, and relations to other nearby entities.
The size and character of the returned network can be modified by input parameters.

Input

  - ctx: A context to control lifecycle.

  - recordKeys: A JSON document listing records.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`

  - maxDegrees: The maximum number of degrees in paths between entities identified by the recordKeys.

  - buildOutDegrees: The number of degrees of relationships to show around each search entity.
    Zero (0) prevents buildout.

  - buildOutMaxEntities: The maximum number of entities to build out in the returned network.

  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindNetworkByRecordID(
	ctx context.Context,
	recordKeys string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(39, recordKeys, maxDegrees, buildOutDegrees, buildOutMaxEntities, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(
				40,
				recordKeys,
				maxDegrees,
				buildOutDegrees,
				buildOutMaxEntities,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.findNetworkByRecordIDV2(
		ctx,
		recordKeys,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordKeys": recordKeys,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8014, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindNetworkByRecordID error: %w", err)
	}

	return result, nil
}

/*
Method FindPathByEntityID finds a relationship path between two entities.
Paths are found using known relationships with other entities.
The path can be modified by input parameters.

Input
  - ctx: A context to control lifecycle.
  - startEntityID: The entity ID for the starting entity of the search path.
  - endEntityID: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidEntityIDs: A JSON document listing entities that should be avoided on the path.
    An empty string disables this capability.
    Example: `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}, {"ENTITY_ID": 3}]}`
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
    An empty string disables this capability.
    Example: `{"DATA_SOURCES": ["MY_DATASOURCE_1", "MY_DATASOURCE_2", "MY_DATASOURCE_3"]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindPathByEntityID(
	ctx context.Context, startEntityID int64, endEntityID int64, maxDegrees int64, avoidEntityIDs string,
	requiredDataSources string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(31, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(32, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources,
				flags, result, err, time.Since(entryTime))
		}()
	}

	switch {
	case len(requiredDataSources) > 0:
		result, err = client.findPathByEntityIDIncludingSourceV2(
			ctx, startEntityID, endEntityID, maxDegrees, avoidEntityIDs,
			requiredDataSources, flags)
	case len(avoidEntityIDs) > 0:
		result, err = client.findPathByEntityIDWithAvoidsV2(
			ctx, startEntityID, endEntityID, maxDegrees, avoidEntityIDs,
			flags)
	default:
		result, err = client.findPathByEntityIDV2(ctx, startEntityID, endEntityID, maxDegrees, flags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityID":       formatEntityID(startEntityID),
				"endEntityID":         formatEntityID(endEntityID),
				"avoidEntityIDs":      avoidEntityIDs,
				"requiredDataSources": requiredDataSources,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8015, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindPathByEntityID error: %w", err)
	}

	return result, nil
}

/*
Method FindPathByRecordID finds a relationship path between two entities identified by record keys.
Paths are found using known relationships with other entities.
The path can be modified by input parameters.

Input

  - ctx: A context to control lifecycle.

  - startDataSourceCode: Identifies the provenance of the record for the starting
    entity of the search path.

  - startRecordID: The unique identifier within the records of the same data source
    for the starting entity of the search path.

  - endDataSourceCode: Identifies the provenance of the record for the ending entity
    of the search path.

  - endRecordID: The unique identifier within the records of the same data source for
    the ending entity of the search path.

  - maxDegrees: The maximum number of degrees in paths between search entities.

  - avoidRecordKeys: A JSON document listing entities that should be avoided on the path.
    An empty string disables this capability.

    Example: `{"RECORDS": [
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "1"},
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "2"},
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "3"}
    ]}`

  - requiredDataSources: A JSON document listing data sources that should be included on the path.
    An empty string disables this capability.
    Example: `{"DATA_SOURCES": ["MY_DATASOURCE_1", "MY_DATASOURCE_2", "MY_DATASOURCE_3"]}`

  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindPathByRecordID(ctx context.Context, startDataSourceCode string, startRecordID string,
	endDataSourceCode string, endRecordID string, maxDegrees int64, avoidRecordKeys string,
	requiredDataSources string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(33, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees,
			avoidRecordKeys, requiredDataSources, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(34, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees,
				avoidRecordKeys, requiredDataSources, flags, result, err, time.Since(entryTime))
		}()
	}

	switch {
	case len(requiredDataSources) > 0:
		result, err = client.findPathByRecordIDIncludingSourceV2(
			ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees, avoidRecordKeys,
			requiredDataSources, flags)
	case len(avoidRecordKeys) > 0:
		result, err = client.findPathByRecordIDWithAvoidsV2(
			ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees, avoidRecordKeys,
			flags)
	default:
		result, err = client.findPathByRecordIDV2(
			ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID,
			maxDegrees, flags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordID":       startRecordID,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordID":         endRecordID,
				"avoidRecordKeys":     avoidRecordKeys,
				"requiredDataSources": requiredDataSources,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8016, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.FindPathByRecordID error: %w", err)
	}

	return result, nil
}

/*
Method GetActiveConfigID returns the Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.

Output
  - configID: The Senzing configuration JSON document identifier that is currently in use by the Senzing engine.
*/
func (client *Szengine) GetActiveConfigID(ctx context.Context) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(35)

		entryTime := time.Now()
		defer func() { client.traceExit(36, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getActiveConfigID(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8017, err, details)
		}()
	}

	if err != nil {
		return 0, fmt.Errorf("szengine.GetActiveConfigID error: %w", err)
	}

	return result, nil
}

/*
Method GetEntityByEntityID returns information about a resolved identity.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
*/
func (client *Szengine) GetEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(37, entityID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(38, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getEntityByEntityIDV2(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8018, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetEntityByEntityID error: %w", err)
	}

	return result, nil
}

/*
Method GetEntityByRecordID returns information about a resolved entity identified by a record
which is a member of the entity.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetEntityByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(39, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(40, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.getEntityByRecordIDV2(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8019, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetEntityByRecordID error: %w", err)
	}

	return result, nil
}

/*
Method GetRecord returns a JSON document containing a single record from the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(45, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(46, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.getRecordV2(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8020, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetRecord error: %w", err)
	}

	return result, nil
}

/*
Method GetRedoRecord returns the next maintenance record from the Senzing datastore.
Usually, [Szengine.ProcessRedoRecord] is called to process the maintenance record retrieved by GetRedoRecord.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document. If no redo records exist, an empty string is returned.
*/
func (client *Szengine) GetRedoRecord(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(47)

		entryTime := time.Now()
		defer func() { client.traceExit(48, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getRedoRecord(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8021, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetRedoRecord error: %w", err)
	}

	return result, nil
}

/*
Method GetStats retrieves workload statistics for the current process.
These statistics are automatically reset after each call.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) GetStats(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(49)

		entryTime := time.Now()
		defer func() { client.traceExit(50, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getStats(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8022, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetStats error: %w", err)
	}

	return result, nil
}

/*
Method GetVirtualEntityByRecordID describes a hypothetical entity based on a list of records.

Input
  - ctx: A context to control lifecycle.
  - recordKeys: A JSON document listing records to include in the hypothetical entity.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetVirtualEntityByRecordID(
	ctx context.Context,
	recordKeys string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(51, recordKeys, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(52, recordKeys, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getVirtualEntityByRecordIDV2(ctx, recordKeys, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordKeys": recordKeys}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8023, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.GetVirtualEntityByRecordID error: %w", err)
	}

	return result, nil
}

/*
Method HowEntityByEntityID explains how an entity was constructed from its constituent records.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) HowEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(53, entityID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(54, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.howEntityByEntityIDV2(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8024, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.HowEntityByEntityID error: %w", err)
	}

	return result, nil
}

/*
Method PreprocessRecord tests adding a record into the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - recordDefinition: A JSON document containing the record to be tested against the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) PreprocessRecord(ctx context.Context, recordDefinition string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(77, recordDefinition, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(78, recordDefinition, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.preprocessRecord(ctx, recordDefinition, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8035, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.PreprocessRecord error: %w", err)
	}

	return result, nil
}

/*
Method PrimeEngine pre-initializes some of the heavier weight internal resources of the Senzing engine.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) PrimeEngine(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(57)

		entryTime := time.Now()
		defer func() { client.traceExit(58, err, time.Since(entryTime)) }()
	}

	err = client.primeEngine(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8026, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.PrimeEngine error: %w", err)
	}

	return nil
}

/*
Method ProcessRedoRecord processes a redo record retrieved by [Szengine.GetRedoRecord].
Calling ProcessRedoRecord has the potential to create more redo records in certain situations.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) ProcessRedoRecord(ctx context.Context, redoRecord string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(59, redoRecord, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(60, redoRecord, flags, result, err, time.Since(entryTime)) }()
	}

	if (flags & senzing.SzWithInfo) == senzing.SzNoFlags {
		result, err = client.processRedoRecord(ctx, redoRecord)
	} else {
		result, err = client.processRedoRecordWithInfo(ctx, redoRecord)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8027, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.ProcessRedoRecord error: %w", err)
	}

	return result, nil
}

/*
Method ReevaluateEntity verifies that the entity is consistent with its records.
If inconsistent, ReevaluateEntity() adjusts the entity definition, splits entities, and/or merges entities.
Usually, the ReevaluateEntity method is called after a Senzing configuration change to impact
entities immediately.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateEntity(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(61, entityID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(62, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	if (flags & senzing.SzWithInfo) == senzing.SzNoFlags {
		result, err = client.reevaluateEntity(ctx, entityID, flags)
	} else {
		finalFlags := flags & ^senzing.SzWithInfo
		result, err = client.reevaluateEntityWithInfo(ctx, entityID, finalFlags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8028, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.ReevaluateEntity error: %w", err)
	}

	return result, nil
}

/*
Method ReevaluateRecord verifies that a record is consistent with the entity to which it belongs.
If inconsistent, ReevaluateRecord() adjusts the entity definition, splits entities, and/or merges entities.
Usually, the ReevaluateRecord method is called after a Senzing configuration change to impact
the record immediately.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(63, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(64, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	if (flags & senzing.SzWithInfo) == senzing.SzNoFlags {
		result, err = client.reevaluateRecord(ctx, dataSourceCode, recordID, flags)
	} else {
		finalFlags := flags & ^senzing.SzWithInfo
		result, err = client.reevaluateRecordWithInfo(ctx, dataSourceCode, recordID, finalFlags)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8029, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.ReevaluateRecord error: %w", err)
	}

	return result, nil
}

/*
Method SearchByAttributes retrieves entity data based on entity attributes
and an optional search profile.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the attributes desired in the result set.
    Example: `{"NAME_FULL": "BOB SMITH", "EMAIL_ADDRESS": "bsmith@work.com"}`
  - searchProfile: The name of the search profile to use in the search.
    An empty string will use the default search profile.
    Example: "SEARCH"
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) SearchByAttributes(
	ctx context.Context,
	attributes string,
	searchProfile string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(69, attributes, searchProfile, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(70, attributes, searchProfile, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.searchByAttributesV3(ctx, attributes, searchProfile, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"attributes":    attributes,
				"searchProfile": searchProfile,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8031, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.SearchByAttributes error: %w", err)
	}

	return result, nil
}

/*
Method WhyEntities explains the ways in which two entities are related to each other.

Input
  - ctx: A context to control lifecycle.
  - entityID1: The first of two entity IDs.
  - entityID2: The second of two entity IDs.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyEntities(
	ctx context.Context,
	entityID1 int64,
	entityID2 int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(71, entityID1, entityID2, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(72, entityID1, entityID2, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.whyEntitiesV2(ctx, entityID1, entityID2, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID1": formatEntityID(entityID1),
				"entityID2": formatEntityID(entityID2),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8032, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.WhyEntities error: %w", err)
	}

	return result, nil
}

/*
Method WhyRecordInEntity explains why a record belongs to its resolved entitiy.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyRecordInEntity(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(73, dataSourceCode, recordID, flags)

		entryTime := time.Now()
		defer func() { client.traceExit(74, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.whyRecordInEntityV2(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8033, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.WhyRecordInEntity error: %w", err)
	}

	return result, nil
}

/*
Method WhyRecords describes ways in which two records are related to each other.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordID1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordID2: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyRecords(
	ctx context.Context,
	dataSourceCode1 string,
	recordID1 string,
	dataSourceCode2 string,
	recordID2 string,
	flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(75, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)

		entryTime := time.Now()
		defer func() {
			client.traceExit(
				76,
				dataSourceCode1,
				recordID1,
				dataSourceCode2,
				recordID2,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.whyRecordsV2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode1": dataSourceCode1,
				"recordID1":       recordID1,
				"dataSourceCode2": dataSourceCode2,
				"recordID2":       recordID2,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8034, err, details)
		}()
	}

	if err != nil {
		return "", fmt.Errorf("szengine.WhyRecords error: %w", err)
	}

	return result, nil
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
Method Destroy will destroy and perform cleanup for the Senzing Sz object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) Destroy(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(11)

		entryTime := time.Now()
		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}

	err = client.destroy(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.Destroy error: %w", err)
	}

	return nil
}

/*
Method GetObserverOrigin returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) GetObserverOrigin(ctx context.Context) string {
	_ = ctx

	return client.observerOrigin
}

/*
Method Initialize initializes the SzEngine object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.  0 for current default configuration.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) Initialize(
	ctx context.Context,
	instanceName string,
	settings string,
	configID int64, verboseLogging int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(55, instanceName, settings, configID, verboseLogging)

		entryTime := time.Now()
		defer func() {
			client.traceExit(56, instanceName, settings, configID, verboseLogging, err, time.Since(entryTime))
		}()
	}

	if configID > 0 {
		err = client.initWithConfigID(ctx, instanceName, settings, configID, verboseLogging)
	} else {
		err = client.init(ctx, instanceName, settings, verboseLogging)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID":       strconv.FormatInt(configID, baseTen),
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8025, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.Initialize error: %w", err)
	}

	return nil
}

/*
Method RegisterObserver adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(703, observer.GetObserverID(ctx))

		entryTime := time.Now()
		defer func() { client.traceExit(704, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	if client.observers == nil {
		client.observers = &subject.SimpleSubject{}
	}

	err = client.observers.RegisterObserver(ctx, observer)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverID(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8702, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.RegisterObserver error: %w", err)
	}

	return nil
}

/*
Method Reinitialize re-initializes the Senzing engine with a specific Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (client *Szengine) Reinitialize(ctx context.Context, configID int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(65, configID)

		entryTime := time.Now()
		defer func() { client.traceExit(66, configID, err, time.Since(entryTime)) }()
	}

	err = client.reinit(ctx, configID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8030, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.Reinitialize error: %w", err)
	}

	return nil
}

/*
Method SetLogLevel sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szengine) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	if client.isTrace {
		client.traceEntry(705, logLevelName)

		entryTime := time.Now()
		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
	}

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s; %w", logLevelName, szerror.ErrSzSdk)
	}

	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8703, err, details)
		}()
	}

	if err != nil {
		return fmt.Errorf("szengine.SetLogLevel error: %w", err)
	}

	return nil
}

/*
Method SetObserverOrigin sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
Method UnregisterObserver removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(707, observer.GetObserverID(ctx))

		entryTime := time.Now()
		defer func() { client.traceExit(708, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8704, err, details)
		err = client.observers.UnregisterObserver(ctx, observer)

		if !client.observers.HasObservers(ctx) {
			client.observers = nil
		}
	}

	if err != nil {
		return fmt.Errorf("szengine.UnregisterObserver error: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------------
// Private methods for calling the Senzing C API
// ----------------------------------------------------------------------------

/*
Method addRecord adds a record into the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing datastore.
*/
func (client *Szengine) addRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	recordDefinition string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	recordDefinitionForC := C.CString(recordDefinition)

	defer C.free(unsafe.Pointer(recordDefinitionForC))

	result := C.Sz_addRecord(dataSourceCodeForC, recordIDForC, recordDefinitionForC)
	if result != noError {
		err = client.newError(ctx, 4001, dataSourceCode, recordID, recordDefinition, result)
	}

	return withoutInfo, err
}

/*
Method addRecordWithInfo adds a record into the Senzing datastore and returns information on the affected entities.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) addRecordWithInfo(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	recordDefinition string, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	recordDefinitionForC := C.CString(recordDefinition)

	defer C.free(unsafe.Pointer(recordDefinitionForC))

	result := C.Sz_addRecordWithInfo_helper(dataSourceCodeForC, recordIDForC, recordDefinitionForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4002, dataSourceCode, recordID, recordDefinition, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) closeExport(ctx context.Context, exportHandle uintptr) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.Sz_closeExport_helper(C.uintptr_t(exportHandle))
	if result != noError {
		err = client.newError(ctx, 4003, exportHandle, result)
	}

	return err
}

func (client *Szengine) countRedoRecords(ctx context.Context) (int64, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	_ = ctx

	result := int64(C.Sz_countRedoRecords())
	if result < 0 {
		err = client.newError(ctx, 4061, result)
	}

	return result, err
}

/*
Method deleteRecord deletes a record from the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
*/
func (client *Szengine) deleteRecord(ctx context.Context, dataSourceCode string, recordID string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_deleteRecord(dataSourceCodeForC, recordIDForC)
	if result != noError {
		err = client.newError(ctx, 4004, dataSourceCode, recordID, result)
	}

	return withoutInfo, err
}

/*
Method deleteRecordWithInfo deletes a record from the Senzing datastore and returns information
on the affected entities.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) deleteRecordWithInfo(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_deleteRecordWithInfo_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4005, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) destroy(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.Sz_destroy()
	if result != noError {
		err = client.newError(ctx, 4006, result)
	}

	return err
}

func (client *Szengine) exportCsvEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err                error
		resultExportHandle uintptr
	)

	csvColumnListForC := C.CString(csvColumnList)

	defer C.free(unsafe.Pointer(csvColumnListForC))

	result := C.Sz_exportCSVEntityReport_helper(csvColumnListForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4007, csvColumnList, flags, result.returnCode, result)
	}

	resultExportHandle = (uintptr)(result.exportHandle)

	return resultExportHandle, err
}

func (client *Szengine) exportJSONEntityReport(ctx context.Context, flags int64) (uintptr, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err                error
		resultExportHandle uintptr
	)

	result := C.Sz_exportJSONEntityReport_helper(C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4008, flags, result.returnCode, result)
	}

	resultExportHandle = (uintptr)(result.exportHandle)

	return resultExportHandle, err
}

func (client *Szengine) fetchNext(ctx context.Context, exportHandle uintptr) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_fetchNext_helper(C.uintptr_t(exportHandle))
	if result.returnCode < 0 {
		err = client.newError(ctx, 4009, exportHandle, result.returnCode, result)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) findInterestingEntitiesByEntityID(
	ctx context.Context,
	entityID int64,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_findInterestingEntitiesByEntityID_helper(C.longlong(entityID), C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4010, entityID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) findInterestingEntitiesByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_findInterestingEntitiesByRecordID_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))

	if result.returnCode != noError {
		err = client.newError(ctx, 4011, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) findNetworkByEntityIDV2(
	ctx context.Context,
	entityIDs string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	entityListForC := C.CString(entityIDs)

	defer C.free(unsafe.Pointer(entityListForC))

	result := C.Sz_findNetworkByEntityID_V2_helper(
		entityListForC,
		C.longlong(maxDegrees),
		C.longlong(buildOutDegrees),
		C.longlong(buildOutMaxEntities),
		C.longlong(flags),
	)
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4013,
			entityIDs,
			maxDegrees,
			buildOutDegrees,
			buildOutMaxEntities,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) findNetworkByRecordIDV2(
	ctx context.Context,
	recordKeys string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	recordListForC := C.CString(recordKeys)

	defer C.free(unsafe.Pointer(recordListForC))

	result := C.Sz_findNetworkByRecordID_V2_helper(
		recordListForC,
		C.longlong(maxDegrees),
		C.longlong(buildOutDegrees),
		C.longlong(buildOutMaxEntities),
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4015,
			recordKeys,
			maxDegrees,
			buildOutDegrees,
			buildOutMaxEntities,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByEntityIDV2 finds single relationship paths between two entities.
Paths are found using known relationships with other entities.

Input
  - ctx: A context to control lifecycle.
  - startEntityID: The entity ID for the starting entity of the search path.
  - endEntityID: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
*/
func (client *Szengine) findPathByEntityIDV2(
	ctx context.Context,
	startEntityID int64,
	endEntityID int64,
	maxDegrees int64,
	flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_findPathByEntityID_V2_helper(
		C.longlong(startEntityID),
		C.longlong(endEntityID),
		C.longlong(maxDegrees),
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4017, startEntityID, endEntityID, maxDegrees, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByRecordIDV2 finds single relationship paths between two entities.
The entities are identified by starting and ending records.
Paths are found using known relationships with other entities.
It extends FindPathByRecordID() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordID: The unique identifier within the records of the same data source for the
    starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordID: The unique identifier within the records of the same data source for the
    ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) findPathByRecordIDV2(
	ctx context.Context,
	startDataSourceCode string,
	startRecordID string,
	endDataSourceCode string,
	endRecordID string,
	maxDegrees int64,
	flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	startDataSourceCodeForC := C.CString(startDataSourceCode)

	defer C.free(unsafe.Pointer(startDataSourceCodeForC))

	startRecordIDForC := C.CString(startRecordID)

	defer C.free(unsafe.Pointer(startRecordIDForC))

	endDataSourceCodeForC := C.CString(endDataSourceCode)

	defer C.free(unsafe.Pointer(endDataSourceCodeForC))

	endRecordIDForC := C.CString(endRecordID)

	defer C.free(unsafe.Pointer(endRecordIDForC))

	result := C.Sz_findPathByRecordID_V2_helper(
		startDataSourceCodeForC,
		startRecordIDForC,
		endDataSourceCodeForC,
		endRecordIDForC,
		C.longlong(maxDegrees),
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4019,
			startDataSourceCode,
			startRecordID,
			endDataSourceCode,
			endRecordID,
			maxDegrees,
			flags,
			result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByEntityIDWithAvoidsV2 finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that avoid certain entities from being on the path.

When avoiding entities, the user may choose to either strictly exclude the entities,
or prefer to avoid the entities but still include them if no other path is found.
By default, entities will be strictly avoided.
A "preferred avoidance" may be done by specifying the SzFindPathStrictAvoid control flag.

Input
  - ctx: A context to control lifecycle.
  - startEntityID: The entity ID for the starting entity of the search path.
  - endEntityID: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidedEntities: A JSON document listing entities that should be avoided on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) findPathByEntityIDWithAvoidsV2(
	ctx context.Context,
	startEntityID int64,
	endEntityID int64,
	maxDegrees int64,
	avoidedEntities string,
	flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	avoidedEntitiesForC := C.CString(avoidedEntities)

	defer C.free(unsafe.Pointer(avoidedEntitiesForC))

	result := C.Sz_findPathByEntityIDWithAvoids_V2_helper(
		C.longlong(startEntityID),
		C.longlong(endEntityID),
		C.longlong(maxDegrees),
		avoidedEntitiesForC,
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4021,
			startEntityID,
			endEntityID,
			maxDegrees,
			avoidedEntities,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByRecordIDWithAvoidsV2 finds single relationship paths between two entities.
Paths are found using known relationships with other entities.
In addition, it will find paths that avoid certain entities from being on the path.
It extends FindPathExcludingByRecordID() by adding output control flags.

When avoiding entities, the user may choose to either strictly exclude the entities,
or prefer to avoid the entities but still include them if no other path is found.
By default, entities will be strictly avoided.
A "preferred avoidance" may be done by specifying the SzFindPathStrictAvoid control flag.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting entity of the search path.
  - startRecordID: The unique identifier within the records of the same data source for the starting
    entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity of the search path.
  - endRecordID: The unique identifier within the records of the same data source for the ending
    entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidedRecords: A JSON document listing records that should be avoided on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) findPathByRecordIDWithAvoidsV2(
	ctx context.Context,
	startDataSourceCode string,
	startRecordID string,
	endDataSourceCode string,
	endRecordID string,
	maxDegrees int64,
	avoidedRecords string,
	flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	startDataSourceCodeForC := C.CString(startDataSourceCode)

	defer C.free(unsafe.Pointer(startDataSourceCodeForC))

	startRecordIDForC := C.CString(startRecordID)

	defer C.free(unsafe.Pointer(startRecordIDForC))

	endDataSourceCodeForC := C.CString(endDataSourceCode)

	defer C.free(unsafe.Pointer(endDataSourceCodeForC))

	endRecordIDForC := C.CString(endRecordID)

	defer C.free(unsafe.Pointer(endRecordIDForC))

	avoidedRecordsForC := C.CString(avoidedRecords)

	defer C.free(unsafe.Pointer(avoidedRecordsForC))

	result := C.Sz_findPathByRecordIDWithAvoids_V2_helper(
		startDataSourceCodeForC,
		startRecordIDForC,
		endDataSourceCodeForC,
		endRecordIDForC,
		C.longlong(maxDegrees),
		avoidedRecordsForC,
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4023,
			startDataSourceCode,
			startRecordID,
			endDataSourceCode,
			endRecordID,
			maxDegrees,
			avoidedRecords,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByEntityIDIncludingSourceV2 finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityID_V2() and FindPathExcludingByRecordID_V2().
It extends FindPathIncludingSourceByEntityID() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - startEntityID: The entity ID for the starting entity of the search path.
  - endEntityID: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidedEntities: A JSON document listing entities that should be avoided on the path.
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) findPathByEntityIDIncludingSourceV2(
	ctx context.Context,
	startEntityID int64,
	endEntityID int64,
	maxDegrees int64,
	avoidedEntities string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	avoidedEntitiesForC := C.CString(avoidedEntities)

	defer C.free(unsafe.Pointer(avoidedEntitiesForC))

	requiredDataSourcesForC := C.CString(requiredDataSources)

	defer C.free(unsafe.Pointer(requiredDataSourcesForC))

	result := C.Sz_findPathByEntityIDIncludingSource_V2_helper(
		C.longlong(startEntityID),
		C.longlong(endEntityID),
		C.longlong(maxDegrees),
		avoidedEntitiesForC,
		requiredDataSourcesForC,
		C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4025,
			startEntityID,
			endEntityID,
			maxDegrees,
			avoidedEntities,
			requiredDataSources,
			flags,
			result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method findPathByRecordIDIncludingSourceV2 finds single relationship paths between two entities.
In addition, one of the enties along the path must include a specified data source.
Specific entities may also be excluded,
using the same methodology as the FindPathExcludingByEntityID_V2() and FindPathExcludingByRecordID_V2().
It extends FindPathIncludingSourceByRecordID() by adding output control flags.

Input

  - ctx: A context to control lifecycle.

  - startDataSourceCode: Identifies the provenance of the record for the starting entity
    of the search path.

  - startRecordID: The unique identifier within the records of the same data source for
    the starting entity of the search path.

  - endDataSourceCode: Identifies the provenance of the record for the ending entity of
    the search path.

  - endRecordID: The unique identifier within the records of the same data source for the
    ending entity of the search path.

  - maxDegrees: The maximum number of degrees in paths between search entities.

  - avoidedRecords: A JSON document listing records that should be avoided on the path.

  - requiredDataSources: A JSON document listing data sources that should be included on the path.

  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) findPathByRecordIDIncludingSourceV2(
	ctx context.Context,
	startDataSourceCode string,
	startRecordID string,
	endDataSourceCode string,
	endRecordID string,
	maxDegrees int64,
	avoidedRecords string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	startDataSourceCodeForC := C.CString(startDataSourceCode)

	defer C.free(unsafe.Pointer(startDataSourceCodeForC))

	startRecordIDForC := C.CString(startRecordID)

	defer C.free(unsafe.Pointer(startRecordIDForC))

	endDataSourceCodeForC := C.CString(endDataSourceCode)

	defer C.free(unsafe.Pointer(endDataSourceCodeForC))

	endRecordIDForC := C.CString(endRecordID)

	defer C.free(unsafe.Pointer(endRecordIDForC))

	avoidedRecordsForC := C.CString(avoidedRecords)

	defer C.free(unsafe.Pointer(avoidedRecordsForC))

	requiredDataSourcesForC := C.CString(requiredDataSources)

	defer C.free(unsafe.Pointer(requiredDataSourcesForC))

	result := C.Sz_findPathByRecordIDIncludingSource_V2_helper(
		startDataSourceCodeForC,
		startRecordIDForC,
		endDataSourceCodeForC,
		endRecordIDForC,
		C.longlong(maxDegrees),
		avoidedRecordsForC,
		requiredDataSourcesForC,
		C.longlong(flags),
	)
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4027,
			startDataSourceCode,
			startRecordID,
			endDataSourceCode,
			endRecordID,
			maxDegrees,
			avoidedRecords,
			requiredDataSources,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getActiveConfigID(ctx context.Context) (int64, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultConfigID int64
	)

	result := C.Sz_getActiveConfigID_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4028, result.returnCode, result)
	}

	resultConfigID = int64(result.configID)

	return resultConfigID, err
}

func (client *Szengine) getEntityByEntityIDV2(ctx context.Context, entityID int64, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_getEntityByEntityID_V2_helper(C.longlong(entityID), C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4030, entityID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getEntityByRecordIDV2(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_getEntityByRecordID_V2_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4032, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getRecordV2(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_getRecord_V2_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4035, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getRedoRecord(ctx context.Context) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_getRedoRecord_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4036, result.returnCode, result)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getStats(ctx context.Context) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_stats_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4054, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) getVirtualEntityByRecordIDV2(
	ctx context.Context,
	recordKeys string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	recordListForC := C.CString(recordKeys)

	defer C.free(unsafe.Pointer(recordListForC))

	result := C.Sz_getVirtualEntityByRecordID_V2_helper(recordListForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4038, recordKeys, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) howEntityByEntityIDV2(ctx context.Context, entityID int64, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_howEntityByEntityID_V2_helper(C.longlong(entityID), C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4040, entityID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method init initializes the SzEngine object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) init(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	instanceNameForC := C.CString(instanceName)

	defer C.free(unsafe.Pointer(instanceNameForC))

	settingsForC := C.CString(settings)

	defer C.free(unsafe.Pointer(settingsForC))

	result := C.Sz_init(instanceNameForC, settingsForC, C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4041, instanceName, settings, verboseLogging, result)
	}

	return err
}

/*
Method initWithConfigID initializes the Senzing Sz object with a non-default configuration ID.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) initWithConfigID(
	ctx context.Context,
	instanceName string,
	settings string,
	configID int64,
	verboseLogging int64,
) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	instanceNameForC := C.CString(instanceName)

	defer C.free(unsafe.Pointer(instanceNameForC))

	settingsForC := C.CString(settings)

	defer C.free(unsafe.Pointer(settingsForC))

	result := C.Sz_initWithConfigID(instanceNameForC, settingsForC, C.longlong(configID), C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4042, instanceName, settings, configID, verboseLogging, result)
	}

	return err
}

func (client *Szengine) primeEngine(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.Sz_primeEngine()
	if result != noError {
		err = client.newError(ctx, 4043, result)
	}

	return err
}

/*
Method preprocessRecord tests adding a record into the Senzing datastore and returns information
on the affected entities.

Input
  - ctx: A context to control lifecycle.
  - recordDefinition: A JSON document containing the record to be added to the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) preprocessRecord(ctx context.Context, recordDefinition string, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	recordDefinitionForC := C.CString(recordDefinition)

	defer C.free(unsafe.Pointer(recordDefinitionForC))

	result := C.Sz_preprocessRecord_helper(recordDefinitionForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4061, recordDefinition, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method processRedoRecord processes given redo record.
Calling processRedoRecord() has the potential to create more redo records in certain situations.

Input
  - ctx: A context to control lifecycle.
  - redoRecord: The redo record to be processed.

Output
  - An empty JSON document.
*/
func (client *Szengine) processRedoRecord(ctx context.Context, redoRecord string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	redoRecordForC := C.CString(redoRecord)

	defer C.free(unsafe.Pointer(redoRecordForC))

	result := C.Sz_processRedoRecord(redoRecordForC)
	if result != noError {
		err = client.newError(ctx, 4044, redoRecord, result)
	}

	return withoutInfo, err
}

/*
Method processRedoRecordWithInfo processes the next redo record and returns it and affected entities.
Calling processRedoRecordWithInfo() has the potential to create more redo records in certain situations.

Input
  - ctx: A context to control lifecycle.
  - redoRecord: The redo record to be processed.

Output
  - A JSON document with affected entities.
*/
func (client *Szengine) processRedoRecordWithInfo(ctx context.Context, redoRecord string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	redoRecordForC := C.CString(redoRecord)

	defer C.free(unsafe.Pointer(redoRecordForC))

	result := C.Sz_processRedoRecordWithInfo_helper(redoRecordForC)
	if result.returnCode != noError {
		err = client.newError(ctx, 4045, redoRecord, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method reevaluateEntity reevaluates the specified entity.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) reevaluateEntity(ctx context.Context, entityID int64, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.Sz_reevaluateEntity(C.longlong(entityID), C.longlong(flags))
	if result != noError {
		err = client.newError(ctx, 4046, entityID, flags, result)
	}

	return withoutInfo, err
}

/*
Method reevaluateEntityWithInfo reevaluates the specified entity and returns information.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
*/
func (client *Szengine) reevaluateEntityWithInfo(ctx context.Context, entityID int64, flags int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_reevaluateEntityWithInfo_helper(C.longlong(entityID), C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4047, entityID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method reevaluateRecord reevaluates a specific record.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) reevaluateRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_reevaluateRecord(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result != noError {
		err = client.newError(ctx, 4048, dataSourceCode, recordID, flags, result)
	}

	return withoutInfo, err
}

/*
Method reevaluateRecordWithInfo reevaluates a specific record and returns information.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
*/
func (client *Szengine) reevaluateRecordWithInfo(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_reevaluateRecordWithInfo_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4049, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szengine) reinit(ctx context.Context, configID int64) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.Sz_reinit(C.longlong(configID))
	if result != noError {
		err = client.newError(ctx, 4050, configID, result)
	}

	return err
}

/*
Method SearchByAttributes_V2 retrieves entity data based on a user-specified set of entity attributes.
It extends SearchByAttributes() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the record to be added to the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
// func (client *Szengine) searchByAttributesV2(ctx context.Context, attributes string, flags int64) (string, error) {
// 	runtime.LockOSThread()
// 	defer runtime.UnlockOSThread()
// 	var err error
// 	var resultResponse string
// 	attributesForC := C.CString(attributes)
// 	defer C.free(unsafe.Pointer(attributesForC))
// 	result := C.Sz_searchByAttributes_V2_helper(attributesForC, C.longlong(flags))
// 	if result.returnCode != noError {
// 		err = client.newError(ctx, 4052, attributes, flags, result.returnCode)
// 	}
// 	resultResponse = C.GoString(result.response)
// 	C.SzHelper_free(unsafe.Pointer(result.response))
// 	return resultResponse, err
// }

/*
Method SearchByAttributes_V3 retrieves entity data based on a user-specified set of entity attributes.
It extends searchByAttributesV2() by adding a search profile parameter.

Input
  - ctx: A context to control lifecycle.
  - jsonData: A JSON document containing the record to be added to the Senzing datastore.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) searchByAttributesV3(
	ctx context.Context,
	attributes string,
	searchProfile string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	attributesForC := C.CString(attributes)

	defer C.free(unsafe.Pointer(attributesForC))

	searchProfileForC := C.CString(searchProfile)

	defer C.free(unsafe.Pointer(searchProfileForC))

	result := C.Sz_searchByAttributes_V3_helper(attributesForC, searchProfileForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4053, attributes, searchProfile, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method whyEntitiesV2 explains why records belong to their resolved entities.
whyEntitiesV2() will compare the record data within an entity
against the rest of the entity data and show why they are connected.
This is calculated based on the features that record data represents.
It extends whyEntities() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - entityID1: The entity ID for the starting entity of the search path.
  - entityID2: The entity ID for the ending entity of the search path.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) whyEntitiesV2(
	ctx context.Context,
	entityID1 int64,
	entityID2 int64,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.Sz_whyEntities_V2_helper(C.longlong(entityID1), C.longlong(entityID2), C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4056, entityID1, entityID2, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method whyRecordInEntityV2 explains why a record belongs to its resolved entitiy.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) whyRecordInEntityV2(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSourceCodeForC := C.CString(dataSourceCode)

	defer C.free(unsafe.Pointer(dataSourceCodeForC))

	recordIDForC := C.CString(recordID)

	defer C.free(unsafe.Pointer(recordIDForC))

	result := C.Sz_whyRecordInEntity_V2_helper(dataSourceCodeForC, recordIDForC, C.longlong(flags))
	if result.returnCode != noError {
		err = client.newError(ctx, 4058, dataSourceCode, recordID, flags, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method whyRecordsV2 explains why records belong to their resolved entities.
It extends WhyRecords() by adding output control flags.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordID1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordID2: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) whyRecordsV2(
	ctx context.Context,
	dataSourceCode1 string,
	recordID1 string,
	dataSourceCode2 string,
	recordID2 string,
	flags int64,
) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	dataSource1CodeForC := C.CString(dataSourceCode1)

	defer C.free(unsafe.Pointer(dataSource1CodeForC))

	recordID1ForC := C.CString(recordID1)

	defer C.free(unsafe.Pointer(recordID1ForC))

	dataSource2CodeForC := C.CString(dataSourceCode2)

	defer C.free(unsafe.Pointer(dataSource2CodeForC))

	recordID2ForC := C.CString(recordID2)

	defer C.free(unsafe.Pointer(recordID2ForC))

	result := C.Sz_whyRecords_V2_helper(
		dataSource1CodeForC,
		recordID1ForC,
		dataSource2CodeForC,
		recordID2ForC,
		C.longlong(flags),
	)
	if result.returnCode != noError {
		err = client.newError(
			ctx,
			4060,
			dataSourceCode1,
			recordID1,
			dataSourceCode2,
			recordID2,
			flags,
			result.returnCode,
		)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (client *Szengine) fetchNextIntoChannel(
	ctx context.Context,
	reportHandle uintptr,
	stringFragmentChannel chan senzing.StringFragment,
) {
	for {
		select {
		case <-ctx.Done():
			stringFragmentChannel <- senzing.StringFragment{
				Error: ctx.Err(),
			}

			break
		default:
			entityReportFragment, err := client.FetchNext(ctx, reportHandle)
			if err != nil {
				stringFragmentChannel <- senzing.StringFragment{
					Error: err,
				}

				break
			}

			if len(entityReportFragment) == 0 {
				break
			}
			stringFragmentChannel <- senzing.StringFragment{
				Value: entityReportFragment,
			}
		}
	}
}

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szengine) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szengine.IDMessages, baseCallerSkip)
	}

	return client.logger
}

// Get the Messenger singleton.
func (client *Szengine) getMessenger() messenger.Messenger {
	if client.messenger == nil {
		client.messenger = helper.GetMessenger(ComponentID, szengine.IDMessages, baseCallerSkip)
	}

	return client.messenger
}

// Trace method entry.
func (client *Szengine) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szengine) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

func formatEntityID(entityID int64) string {
	return strconv.FormatInt(entityID, baseTen)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szengine) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
	defer func() { client.panicOnError(client.clearLastException(ctx)) }()

	lastExceptionCode, _ := client.getLastExceptionCode(ctx)

	lastException, err := client.getLastException(ctx)
	if err != nil {
		lastException = err.Error()
	}

	details = append(details, messenger.MessageCode{Value: fmt.Sprintf(ExceptionCodeTemplate, lastExceptionCode)})
	details = append(details, messenger.MessageReason{Value: lastException})
	details = append(details, fmt.Errorf("%s; %w", lastException, szerror.ErrSz))
	errorMessage := client.getMessenger().NewJSON(errorNumber, details...)

	return szerror.New(lastExceptionCode, errorMessage) //nolint
}

/*
Method panicOnError calls panic() when an error is not nil.

Input:
  - err: nil or an actual error
*/
func (client *Szengine) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// --- Sz exception handling --------------------------------------------------

/*
Method clearLastException erases the last exception message held by the Senzing Sz object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) clearLastException(ctx context.Context) error {
	var err error

	_ = ctx

	if client.isTrace {
		client.traceEntry(3)

		entryTime := time.Now()
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}

	C.Sz_clearLastException()

	return err
}

/*
Method getLastException retrieves the last exception thrown in Senzing's Sz.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's Sz.
*/
func (client *Szengine) getLastException(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	_ = ctx

	if client.isTrace {
		client.traceEntry(41)

		entryTime := time.Now()
		defer func() { client.traceExit(42, result, err, time.Since(entryTime)) }()
	}

	stringBuffer := client.getByteArray(initialByteArraySize)
	C.Sz_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	result = string(bytes.Trim(stringBuffer, "\x00"))

	return result, err
}

/*
Method getLastExceptionCode retrieves the code of the last exception thrown in Senzing's Sz.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's Sz.
*/
func (client *Szengine) getLastExceptionCode(ctx context.Context) (int, error) {
	var (
		err    error
		result int
	)

	_ = ctx

	if client.isTrace {
		client.traceEntry(43)

		entryTime := time.Now()
		defer func() { client.traceExit(44, result, err, time.Since(entryTime)) }()
	}

	result = int(C.Sz_getLastExceptionCode())

	return result, err
}

// --- Misc -------------------------------------------------------------------

// Make a byte array.
func (client *Szengine) getByteArray(size int) []byte {
	return make([]byte, size)
}

// A hack: Only needed to import the "senzing" package for the godoc comments.
// func junk() {
// 	fmt.Printf(senzing.SzNoAttributes)
// }
