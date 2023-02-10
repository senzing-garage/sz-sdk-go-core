/*
The G2diagnostic implementation is a wrapper over the Senzing libg2diagnostic library.
*/
package g2diagnostic

/*
#include "g2diagnostic.h"
#cgo CFLAGS: -g -I/opt/senzing/g2/sdk/c
#cgo LDFLAGS: -L/opt/senzing/g2/lib -lG2
*/
import "C"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	g2diagnosticapi "github.com/senzing/g2-sdk-go/g2diagnostic"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// G2diagnostic is the default implementation of the G2diagnostic interface.
type G2diagnostic struct {
	isTrace   bool
	logger    messagelogger.MessageLoggerInterface
	observers subject.Subject
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const initialByteArraySize = 65535

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (client *G2diagnostic) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *G2diagnostic) getByteArray(size int) []byte {
	return make([]byte, size)
}

// Create a new error.
func (client *G2diagnostic) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
	lastException, err := client.getLastException(ctx)
	defer client.clearLastException(ctx)
	message := lastException
	if err != nil {
		message = err.Error()
	}

	var newDetails []interface{}
	newDetails = append(newDetails, details...)
	newDetails = append(newDetails, errors.New(message))
	errorMessage, err := client.getLogger().Message(errorNumber, newDetails...)
	if err != nil {
		errorMessage = err.Error()
	}

	return errors.New(errorMessage)
}

// Get the Logger singleton.
func (client *G2diagnostic) getLogger() messagelogger.MessageLoggerInterface {
	if client.logger == nil {
		client.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, g2diagnosticapi.IdMessages, g2diagnosticapi.IdStatuses, messagelogger.LevelInfo)
	}
	return client.logger
}

func (client *G2diagnostic) notify(ctx context.Context, messageId int, err error, details map[string]string) {
	now := time.Now()
	details["subjectId"] = strconv.Itoa(ProductId)
	details["messageId"] = strconv.Itoa(messageId)
	details["messageTime"] = strconv.FormatInt(now.UnixNano(), 10)
	if err != nil {
		details["error"] = err.Error()
	}
	message, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		client.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (client *G2diagnostic) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2diagnostic) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

/*
The clearLastException method erases the last exception message held by the Senzing G2Config object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnostic) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2Diagnostic_clearLastException();
	if client.isTrace {
		client.traceEntry(3)
	}
	entryTime := time.Now()
	var err error = nil
	C.G2Diagnostic_clearLastException()
	if client.isTrace {
		defer client.traceExit(4, err, time.Since(entryTime))
	}
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's G2Diagnostic.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2Config.
*/
func (client *G2diagnostic) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	if client.isTrace {
		client.traceEntry(31)
	}
	entryTime := time.Now()
	var err error = nil
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2Diagnostic_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.ulong(len(stringBuffer)))
	// if result == 0 { // "result" is length of exception message.
	// 	err = client.getLogger().Error(4014, result, time.Since(entryTime))
	// }
	stringBuffer = bytes.Trim(stringBuffer, "\x00")
	if client.isTrace {
		defer client.traceExit(32, string(stringBuffer), err, time.Since(entryTime))
	}
	return string(stringBuffer), err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2Diagnostic.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2Config.
*/
func (client *G2diagnostic) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Diagnostic_getLastExceptionCode();
	if client.isTrace {
		client.traceEntry(33)
	}
	entryTime := time.Now()
	var err error = nil
	result := int(C.G2Diagnostic_getLastExceptionCode())
	if client.isTrace {
		defer client.traceExit(34, result, err, time.Since(entryTime))
	}
	return result, err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The CheckDBPerf method performs inserts to determine rate of insertion.

Input
  - ctx: A context to control lifecycle.
  - secondsToRun: Duration of the test in seconds.

Output

  - A string containing a JSON document.
    Example: `{"numRecordsInserted":0,"insertTime":0}`
*/
func (client *G2diagnostic) CheckDBPerf(ctx context.Context, secondsToRun int) (string, error) {
	// _DLEXPORT int G2Diagnostic_checkDBPerf(int secondsToRun, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(1, secondsToRun)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_checkDBPerf_helper(C.int(secondsToRun))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4001, secondsToRun, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8001, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(2, secondsToRun, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The CloseEntityListBySize method closes the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is created by the GetEntityListBySize() method.

Input
  - ctx: A context to control lifecycle.
  - entityListBySizeHandle: A handle created by GetEntityListBySize().
*/
func (client *G2diagnostic) CloseEntityListBySize(ctx context.Context, entityListBySizeHandle uintptr) error {
	//  _DLEXPORT int G2Diagnostic_closeEntityListBySize(EntityListBySizeHandle entityListBySizeHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(5)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_closeEntityListBySize_helper(C.uintptr_t(entityListBySizeHandle))
	if result != 0 {
		err = client.newError(ctx, 4002, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8002, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(6, err, time.Since(entryTime))
	}
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Diagnostic object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnostic) Destroy(ctx context.Context) error {
	//  _DLEXPORT int G2Diagnostic_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(7)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_destroy()
	if result != 0 {
		err = client.newError(ctx, 4003, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8003, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(8, err, time.Since(entryTime))
	}
	return err
}

/*
The FetchNextEntityBySize method gets the next section of the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is created by the GetEntityListBySize() method.

Input
  - ctx: A context to control lifecycle.
  - entityListBySizeHandle: A handle created by GetEntityListBySize().

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle uintptr) (string, error) {
	//  _DLEXPORT int G2Diagnostic_fetchNextEntityBySize(EntityListBySizeHandle entityListBySizeHandle, char *responseBuf, const size_t bufSize);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(9)
	}
	entryTime := time.Now()
	var err error = nil
	stringBuffer := client.getByteArray(initialByteArraySize)
	result := C.G2Diagnostic_fetchNextEntityBySize_helper(C.uintptr_t(entityListBySizeHandle), (*C.char)(unsafe.Pointer(&stringBuffer[0])), C.ulong(len(stringBuffer)))
	if result < 0 {
		err = client.newError(ctx, 4004, result, time.Since(entryTime))
	}
	stringBuffer = bytes.Trim(stringBuffer, "\x00")
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8004, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(10, string(stringBuffer), err, time.Since(entryTime))
	}
	return string(stringBuffer), err
}

/*
The FindEntitiesByFeatureIDs method finds entities having any of the lib feat id specified in the "features" JSON document.
The "features" also contains an entity id.
This entity is ignored in the returned values.

Input
  - ctx: A context to control lifecycle.
  - features: A JSON document having the format: `{"ENTITY_ID":<entity id>,"LIB_FEAT_IDS":[<id1>,<id2>,...<idn>]}` where ENTITY_ID specifies the entity to ignore in the returns and <id#> are the lib feat ids used to query for entities.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error) {
	//  _DLEXPORT int G2Diagnostic_findEntitiesByFeatureIDs(const char *features, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(11, features)
	}
	entryTime := time.Now()
	var err error = nil
	featuresForC := C.CString(features)
	defer C.free(unsafe.Pointer(featuresForC))
	result := C.G2Diagnostic_findEntitiesByFeatureIDs_helper(featuresForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4005, features, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8005, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(12, features, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetAvailableMemory method returns the available memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of available memory.
*/
func (client *G2diagnostic) GetAvailableMemory(ctx context.Context) (int64, error) {
	// _DLEXPORT long long G2Diagnostic_getAvailableMemory();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(13)
	}
	entryTime := time.Now()
	var err error = nil
	result := int64(C.G2Diagnostic_getAvailableMemory())
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8006, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(14, result, err, time.Since(entryTime))
	}
	return result, err
}

/*
The GetDataSourceCounts method returns information about data sources.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    See the example output.
*/
func (client *G2diagnostic) GetDataSourceCounts(ctx context.Context) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getDataSourceCounts(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(15)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getDataSourceCounts_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4006, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8007, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(16, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetDBInfo method returns information about the database connection.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    Example: `{"Hybrid Mode":false,"Database Details":[{"Name":"0.0.0.0","Type":"postgresql"}]}`
*/
func (client *G2diagnostic) GetDBInfo(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Diagnostic_getDBInfo(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(17)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getDBInfo_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4007, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8008, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(18, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetEntityDetails method returns information about the database connection.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - includeInternalFeatures: FIXME:

Output
  - A JSON document enumerating FIXME:.
    See the example output.
*/
func (client *G2diagnostic) GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getEntityDetails(const long long entityID, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(19, entityID, includeInternalFeatures)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getEntityDetails_helper(C.longlong(entityID), C.int(includeInternalFeatures))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4008, entityID, includeInternalFeatures, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8009, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(20, entityID, includeInternalFeatures, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetEntityListBySize method gets the next section of the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is used by the FetchNextEntityBySize() and CloseEntityListBySize() methods.

Input
  - ctx: A context to control lifecycle.
  - entitySize: FIXME:

Output
  - A handle to an entity list to be used with FetchNextEntityBySize() and CloseEntityListBySize().
*/
func (client *G2diagnostic) GetEntityListBySize(ctx context.Context, entitySize int) (uintptr, error) {
	//  _DLEXPORT int G2Diagnostic_getEntityListBySize(const size_t entitySize, EntityListBySizeHandle* entityListBySizeHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(21, entitySize)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getEntityListBySize_helper(C.size_t(entitySize))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4009, entitySize, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8010, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(22, entitySize, (uintptr)(result.response), err, time.Since(entryTime))
	}
	return (uintptr)(result.response), err
}

/*
The GetEntityResume method FIXME:

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetEntityResume(ctx context.Context, entityID int64) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getEntityResume(const long long entityID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(23, entityID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getEntityResume_helper(C.longlong(entityID))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4010, entityID, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8011, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(24, entityID, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetEntitySizeBreakdown method FIXME:

Input
  - ctx: A context to control lifecycle.
  - minimumEntitySize: FIXME:
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getEntitySizeBreakdown(const size_t minimumEntitySize, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(25, minimumEntitySize, includeInternalFeatures)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getEntitySizeBreakdown_helper(C.size_t(minimumEntitySize), C.int(includeInternalFeatures))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4011, minimumEntitySize, includeInternalFeatures, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8012, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(26, minimumEntitySize, includeInternalFeatures, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetFeature method retrieves a stored feature.

Input
  - ctx: A context to control lifecycle.
  - libFeatID: The identifier of the feature requested in the search.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetFeature(ctx context.Context, libFeatID int64) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getFeature(const long long libFeatID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(27, libFeatID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getFeature_helper(C.longlong(libFeatID))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4012, libFeatID, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8013, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(28, libFeatID, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetGenericFeatures method retrieves a stored feature.

Input
  - ctx: A context to control lifecycle.
  - featureType: FIXME:
  - maximumEstimatedCount: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getGenericFeatures(const char* featureType, const size_t maximumEstimatedCount, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(29, featureType, maximumEstimatedCount)
	}
	entryTime := time.Now()
	var err error = nil
	featureTypeForC := C.CString(featureType)
	defer C.free(unsafe.Pointer(featureTypeForC))
	result := C.G2Diagnostic_getGenericFeatures_helper(featureTypeForC, C.size_t(maximumEstimatedCount))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4013, featureType, maximumEstimatedCount, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8014, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(30, featureType, maximumEstimatedCount, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetLogicalCores method returns the number of logical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of logical cores.
*/
func (client *G2diagnostic) GetLogicalCores(ctx context.Context) (int, error) {
	// _DLEXPORT int G2Diagnostic_getLogicalCores();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(35)
	}
	entryTime := time.Now()
	var err error = nil
	result := int(C.G2Diagnostic_getLogicalCores())
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8015, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(36, result, err, time.Since(entryTime))
	}
	return result, err
}

/*
The GetMappingStatistics method FIXME:

Input
  - ctx: A context to control lifecycle.
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getMappingStatistics(const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(37, includeInternalFeatures)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getMappingStatistics_helper(C.int(includeInternalFeatures))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4015, includeInternalFeatures, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8016, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(38, includeInternalFeatures, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetPhysicalCores method returns the number of physical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of physical cores.
*/
func (client *G2diagnostic) GetPhysicalCores(ctx context.Context) (int, error) {
	// _DLEXPORT int G2Diagnostic_getPhysicalCores();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(39)
	}
	entryTime := time.Now()
	var err error = nil
	result := int(C.G2Diagnostic_getPhysicalCores())
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8017, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(40, result, err, time.Since(entryTime))
	}
	return result, err
}

/*
The GetRelationshipDetails method FIXME:

Input
  - ctx: A context to control lifecycle.
  - relationshipID: FIXME:
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getRelationshipDetails(const long long relationshipID, const int includeInternalFeatures, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(41, relationshipID, includeInternalFeatures)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getRelationshipDetails_helper(C.longlong(relationshipID), C.int(includeInternalFeatures))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4016, relationshipID, includeInternalFeatures, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8018, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(42, relationshipID, includeInternalFeatures, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetResolutionStatistics method FIXME:

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetResolutionStatistics(ctx context.Context) (string, error) {
	//  _DLEXPORT int G2Diagnostic_getResolutionStatistics(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(43)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_getResolutionStatistics_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4017, result.returnCode, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8019, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(44, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2diagnosticInterface.
For this implementation, "base" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnostic) GetSdkId(ctx context.Context) (string, error) {
	return "base", nil
}

/*
The GetTotalSystemMemory method returns the total memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of memory.
*/
func (client *G2diagnostic) GetTotalSystemMemory(ctx context.Context) (int64, error) {
	// _DLEXPORT long long G2Diagnostic_getTotalSystemMemory();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(57)
	}
	entryTime := time.Now()
	var err error = nil
	result := int64(C.G2Diagnostic_getTotalSystemMemory())
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8020, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(46, result, err, time.Since(entryTime))
	}
	return result, err
}

/*
The Init method initializes the Senzing G2Diagnosis object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2diagnostic) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	// _DLEXPORT int G2Diagnostic_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(47, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	var err error = nil
	moduleNameForC := C.CString(moduleName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(iniParams)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2Diagnostic_init(moduleNameForC, iniParamsForC, C.int(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4018, moduleName, iniParams, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8021, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(48, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The InitWithConfigID method initializes the Senzing G2Diagnosis object with a non-default configuration ID.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - initConfigID: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2diagnostic) InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error {
	//  _DLEXPORT int G2Diagnostic_initWithConfigID(const char *moduleName, const char *iniParams, const long long initConfigID, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(49, moduleName, iniParams, initConfigID, verboseLogging)
	}
	entryTime := time.Now()
	var err error = nil
	moduleNameForC := C.CString(moduleName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(iniParams)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2Diagnostic_initWithConfigID(moduleNameForC, iniParamsForC, C.longlong(initConfigID), C.int(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4019, moduleName, iniParams, initConfigID, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"initConfigID":   strconv.FormatInt(initConfigID, 10),
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8022, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(50, moduleName, iniParams, initConfigID, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

// Null shows how to report a BUG inline.
func (client *G2diagnostic) Null(ctx context.Context) (int64, error) {
	// BUG(mjd): Just an example of how to show bugs in GoDoc.
	var err error = nil
	return 1, err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnostic) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	return client.observers.RegisterObserver(ctx, observer)
}

/*
The Reinit method re-initializes the Senzing G2Diagnosis object.

Input
  - ctx: A context to control lifecycle.
  - initConfigID: The configuration ID used for the initialization.
*/
func (client *G2diagnostic) Reinit(ctx context.Context, initConfigID int64) error {
	//  _DLEXPORT int G2Diagnostic_reinit(const long long initConfigID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(51, initConfigID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Diagnostic_reinit(C.longlong(initConfigID))
	if result != 0 {
		err = client.newError(ctx, 4020, initConfigID, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"initConfigID": strconv.FormatInt(initConfigID, 10),
			}
			client.notify(ctx, 8023, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(52, initConfigID, err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2diagnostic) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(53, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	client.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	client.isTrace = (client.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if client.isTrace {
		defer client.traceExit(54, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnostic) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	err := client.observers.UnregisterObserver(ctx, observer)
	if err != nil {
		return err
	}
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
