/*
The G2configmgr implementation is a wrapper over the Senzing libg2configmgr library.
*/
package g2configmgr

/*
#include "g2configmgr.h"
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

	g2configmgrapi "github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2error"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// G2configmgr is the default implementation of the G2configmgr interface.
type G2configmgr struct {
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
func (client *G2configmgr) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *G2configmgr) getByteArray(size int) []byte {
	return make([]byte, size)
}

// Create a new error.
func (client *G2configmgr) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
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

	return g2error.G2Error(g2error.G2ErrorCode(message), (errorMessage))
}

// Get the Logger singleton.
func (client *G2configmgr) getLogger() messagelogger.MessageLoggerInterface {
	if client.logger == nil {
		client.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, g2configmgrapi.IdMessages, g2configmgrapi.IdStatuses, messagelogger.LevelInfo)
	}
	return client.logger
}

// Notify registered observers.
func (client *G2configmgr) notify(ctx context.Context, messageId int, err error, details map[string]string) {
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
func (client *G2configmgr) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2configmgr) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

/*
The clearLastException method erases the last exception message held by the Senzing G2ConfigMgr object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgr) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2Config_clearLastException()
	if client.isTrace {
		client.traceEntry(3)
	}
	entryTime := time.Now()
	var err error = nil
	C.G2ConfigMgr_clearLastException()
	if client.isTrace {
		defer client.traceExit(4, err, time.Since(entryTime))
	}
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's G2ConfigMgr.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2ConfigMgr.
*/
func (client *G2configmgr) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	if client.isTrace {
		client.traceEntry(13)
	}
	entryTime := time.Now()
	var err error = nil
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2ConfigMgr_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.ulong(len(stringBuffer)))
	// if result == 0 { // "result" is length of exception message.
	// 	err = client.getLogger().Error(4006, result, time.Since(entryTime))
	// }
	stringBuffer = bytes.Trim(stringBuffer, "\x00")
	if client.isTrace {
		defer client.traceExit(14, string(stringBuffer), err, time.Since(entryTime))
	}
	return string(stringBuffer), err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2ConfigMgr.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2ConfigMgr.
*/
func (client *G2configmgr) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Config_getLastExceptionCode();
	if client.isTrace {
		client.traceEntry(15)
	}
	entryTime := time.Now()
	var err error = nil
	result := int(C.G2ConfigMgr_getLastExceptionCode())
	if client.isTrace {
		defer client.traceExit(16, result, err, time.Since(entryTime))
	}
	return result, err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The AddConfig method adds a Senzing configuration JSON document to the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configStr: The Senzing configuration JSON document.
  - configComments: A free-form string of comments describing the configuration document.

Output
  - A configuration identifier.
*/
func (client *G2configmgr) AddConfig(ctx context.Context, configStr string, configComments string) (int64, error) {
	// _DLEXPORT int G2ConfigMgr_addConfig(const char* configStr, const char* configComments, long long* configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(1, configStr, configComments)
	}
	entryTime := time.Now()
	var err error = nil
	configStrForC := C.CString(configStr)
	defer C.free(unsafe.Pointer(configStrForC))
	configCommentsForC := C.CString(configComments)
	defer C.free(unsafe.Pointer(configCommentsForC))
	result := C.G2ConfigMgr_addConfig_helper(configStrForC, configCommentsForC)
	if result.returnCode != 0 {
		err = client.newError(ctx, 4001, configStr, configComments, result.returnCode, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configComments": configComments,
			}
			client.notify(ctx, 8001, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(2, configStr, configComments, int64(C.longlong(result.configID)), err, time.Since(entryTime))
	}
	return int64(C.longlong(result.configID)), err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2ConfigMgr object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgr) Destroy(ctx context.Context) error {
	// _DLEXPORT int G2Config_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(5)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_destroy()
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
The GetConfig method retrieves a specific Senzing configuration JSON document from the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configID: The configuration identifier of the desired Senzing Engine configuration JSON document to retrieve.

Output
  - A JSON document containing the Senzing configuration.
    See the example output.
*/
func (client *G2configmgr) GetConfig(ctx context.Context, configID int64) (string, error) {
	// _DLEXPORT int G2ConfigMgr_getConfig(const long long configID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(7, configID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_getConfig_helper(C.longlong(configID))
	if result.returnCode != 0 {
		err = client.newError(ctx, 4003, configID, result.returnCode, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8003, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(8, configID, C.GoString(result.config), err, time.Since(entryTime))
	}
	return C.GoString(result.config), err
}

/*
The GetConfigList method retrieves a list of Senzing configurations from the Senzing database.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing configurations.
    See the example output.
*/
func (client *G2configmgr) GetConfigList(ctx context.Context) (string, error) {
	// _DLEXPORT int G2ConfigMgr_getConfigList(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(9)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_getConfigList_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4004, result.returnCode, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8004, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(10, C.GoString(result.configList), err, time.Since(entryTime))
	}
	return C.GoString(result.configList), err
}

/*
The GetDefaultConfigID method retrieves from the Senzing database the configuration identifier of the default Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A configuration identifier which identifies the current configuration in use.
*/
func (client *G2configmgr) GetDefaultConfigID(ctx context.Context) (int64, error) {
	//  _DLEXPORT int G2ConfigMgr_getDefaultConfigID(long long* configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(11)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_getDefaultConfigID_helper()
	if result.returnCode != 0 {
		err = client.newError(ctx, 4005, result.returnCode, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8005, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(12, int64(C.longlong(result.configID)), err, time.Since(entryTime))
	}
	return int64(C.longlong(result.configID)), err
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2configmgrInterface.
For this implementation, "base" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgr) GetSdkId(ctx context.Context) string {
	if client.isTrace {
		client.traceEntry(29)
	}
	entryTime := time.Now()
	var err error = nil
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8010, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(30, err, time.Since(entryTime))
	}
	return "base"
}

/*
The Init method initializes the Senzing G2ConfigMgr object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2configmgr) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	// _DLEXPORT int G2Config_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(17, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	var err error = nil
	moduleNameForC := C.CString(moduleName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(iniParams)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2ConfigMgr_init(moduleNameForC, iniParamsForC, C.int(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4007, moduleName, iniParams, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8006, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(18, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2configmgr) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(25, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	err := client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			client.notify(ctx, 8010, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(26, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The ReplaceDefaultConfigID method replaces the old configuration identifier with a new configuration identifier in the Senzing database.
It is like a "compare-and-swap" instruction to serialize concurrent editing of configuration.
If oldConfigID is no longer the "old configuration identifier", the operation will fail.
To simply set the default configuration ID, use SetDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - oldConfigID: The configuration identifier to replace.
  - newConfigID: The configuration identifier to use as the default.
*/
func (client *G2configmgr) ReplaceDefaultConfigID(ctx context.Context, oldConfigID int64, newConfigID int64) error {
	// _DLEXPORT int G2ConfigMgr_replaceDefaultConfigID(const long long oldConfigID, const long long newConfigID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(19, oldConfigID, newConfigID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_replaceDefaultConfigID(C.longlong(oldConfigID), C.longlong(newConfigID))
	if result != 0 {
		err = client.newError(ctx, 4008, oldConfigID, newConfigID, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"newConfigID": strconv.FormatInt(newConfigID, 10),
			}
			client.notify(ctx, 8007, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(20, oldConfigID, newConfigID, err, time.Since(entryTime))
	}
	return err
}

/*
The SetDefaultConfigID method replaces the sets a new configuration identifier in the Senzing database.
To serialize modifying of the configuration identifier, see ReplaceDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - configID: The configuration identifier of the Senzing Engine configuration to use as the default.
*/
func (client *G2configmgr) SetDefaultConfigID(ctx context.Context, configID int64) error {
	// _DLEXPORT int G2ConfigMgr_setDefaultConfigID(const long long configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(21, configID)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2ConfigMgr_setDefaultConfigID(C.longlong(configID))
	if result != 0 {
		err = client.newError(ctx, 4009, configID, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, 10),
			}
			client.notify(ctx, 8008, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(22, configID, err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2configmgr) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if client.isTrace {
		client.traceEntry(23, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	client.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	client.isTrace = (client.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			client.notify(ctx, 8011, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(24, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2configmgr) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(27, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		client.notify(ctx, 8012, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	if client.isTrace {
		defer client.traceExit(28, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
