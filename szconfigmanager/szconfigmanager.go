/*
The szconfigmanager implementation is a wrapper over the Senzing libg2configmgr library.
*/
package szconfigmanager

/*
#include <stdlib.h>
#include "libg2configmgr.h"
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
	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

type Szconfigmanager struct {
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
)

// ----------------------------------------------------------------------------
// sz-sdk-go.SzConfigManager interface methods
// ----------------------------------------------------------------------------

/*
The AddConfig method adds a Senzing configuration JSON document to the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: The Senzing configuration JSON document.
  - configComment: A free-form string describing the configuration document.

Output
  - A configuration identifier.
*/
func (client *Szconfigmanager) AddConfig(ctx context.Context, configDefinition string, configComment string) (int64, error) {
	var err error
	var result int64
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, configDefinition, configComment)
		defer func() {
			client.traceExit(2, configDefinition, configComment, result, err, time.Since(entryTime))
		}()
	}
	result, err = client.addConfig(ctx, configDefinition, configComment)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configComment": configComment,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}
	return result, err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2ConfigMgr object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfigmanager) Destroy(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}
	err = client.destroy(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
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
func (client *Szconfigmanager) GetConfig(ctx context.Context, configID int64) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7, configID)
		defer func() { client.traceExit(8, configID, result, err, time.Since(entryTime)) }()
	}
	result, err = client.getConfig(ctx, configID)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}
	return result, err
}

/*
The GetConfigs method retrieves a list of Senzing configurations from the Senzing database.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing configurations.
    See the example output.
*/
func (client *Szconfigmanager) GetConfigs(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}
	result, err = client.getConfigList(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}
	return result, err
}

/*
The GetDefaultConfigID method retrieves from the Senzing database the configuration identifier of the default Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A configuration identifier which identifies the current configuration in use.
*/
func (client *Szconfigmanager) GetDefaultConfigID(ctx context.Context) (int64, error) {
	var err error
	var result int64
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	result, err = client.getDefaultConfigID(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}
	return result, err
}

/*
The ReplaceDefaultConfigID method replaces the old configuration identifier with a new configuration identifier in the Senzing database.
It is like a "compare-and-swap" instruction to serialize concurrent editing of configuration.
If currentDefaultConfigID is no longer the "old configuration identifier", the operation will fail.
To simply set the default configuration ID, use SetDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - currentDefaultConfigID: The configuration identifier to replace.
  - newDefaultConfigID: The configuration identifier to use as the default.
*/
func (client *Szconfigmanager) ReplaceDefaultConfigID(ctx context.Context, currentDefaultConfigID int64, newDefaultConfigID int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, currentDefaultConfigID, newDefaultConfigID)
		defer func() { client.traceExit(20, currentDefaultConfigID, newDefaultConfigID, err, time.Since(entryTime)) }()
	}
	err = client.replaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"newDefaultConfigID": strconv.FormatInt(newDefaultConfigID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
		}()
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
func (client *Szconfigmanager) SetDefaultConfigID(ctx context.Context, configID int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, configID)
		defer func() { client.traceExit(22, configID, err, time.Since(entryTime)) }()
	}
	err = client.setDefaultConfigID(ctx, configID)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}
	return err
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
func (client *Szconfigmanager) GetObserverOrigin(ctx context.Context) string {
	_ = ctx
	return client.observerOrigin
}

/*
The Initialize method initializes the Senzing G2ConfigMgr object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfigmanager) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(18, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	err = client.init(ctx, instanceName, settings, verboseLogging)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(703, observer.GetObserverID(ctx))
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
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szconfigmanager) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(705, logLevelName)
		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8703, err, details)
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
func (client *Szconfigmanager) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(707, observer.GetObserverID(ctx))
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
	return err
}

// ----------------------------------------------------------------------------
// Private methods for calling the Senzing C API
// ----------------------------------------------------------------------------

func (client *Szconfigmanager) addConfig(ctx context.Context, configDefinition string, configComment string) (int64, error) {
	// _DLEXPORT int G2ConfigMgr_addConfig(const char* configStr, const char* configComments, long long* configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultConfigID int64
	configDefinitionForC := C.CString(configDefinition)
	defer C.free(unsafe.Pointer(configDefinitionForC))
	configCommentForC := C.CString(configComment)
	defer C.free(unsafe.Pointer(configCommentForC))
	result := C.G2ConfigMgr_addConfig_helper(configDefinitionForC, configCommentForC)
	if result.returnCode != noError {
		err = client.newError(ctx, 4001, configDefinition, configComment, result.returnCode, result)
	}
	resultConfigID = int64(C.longlong(result.configID))
	return resultConfigID, err
}

func (client *Szconfigmanager) destroy(ctx context.Context) error {
	// _DLEXPORT int G2Config_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.G2ConfigMgr_destroy()
	if result != noError {
		err = client.newError(ctx, 4002, result)
	}
	return err
}

func (client *Szconfigmanager) getConfig(ctx context.Context, configID int64) (string, error) {
	// _DLEXPORT int G2ConfigMgr_getConfig(const long long configID, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.G2ConfigMgr_getConfig_helper(C.longlong(configID))
	if result.returnCode != noError {
		err = client.newError(ctx, 4003, configID, result.returnCode, result)
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	return resultResponse, err
}

func (client *Szconfigmanager) getConfigList(ctx context.Context) (string, error) {
	// _DLEXPORT int G2ConfigMgr_getConfigList(char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.G2ConfigMgr_getConfigList_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4004, result.returnCode, result)
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	return resultResponse, err
}

func (client *Szconfigmanager) getDefaultConfigID(ctx context.Context) (int64, error) {
	//  _DLEXPORT int G2ConfigMgr_getDefaultConfigID(long long* configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultConfigID int64
	result := C.G2ConfigMgr_getDefaultConfigID_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4005, result.returnCode, result)
	}
	resultConfigID = int64(C.longlong(result.configID))
	return resultConfigID, err
}

func (client *Szconfigmanager) init(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	// _DLEXPORT int G2Config_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	moduleNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2ConfigMgr_init(moduleNameForC, iniParamsForC, C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4006, instanceName, settings, verboseLogging, result)
	}
	return err
}

func (client *Szconfigmanager) replaceDefaultConfigID(ctx context.Context, currentDefaultConfigID int64, newDefaultConfigID int64) error {
	// _DLEXPORT int G2ConfigMgr_replaceDefaultConfigID(const long long oldConfigID, const long long newConfigID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.G2ConfigMgr_replaceDefaultConfigID(C.longlong(currentDefaultConfigID), C.longlong(newDefaultConfigID))
	if result != noError {
		err = client.newError(ctx, 4007, currentDefaultConfigID, newDefaultConfigID, result)
	}
	return err
}

func (client *Szconfigmanager) setDefaultConfigID(ctx context.Context, configID int64) error {
	// _DLEXPORT int G2ConfigMgr_setDefaultConfigID(const long long configID);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.G2ConfigMgr_setDefaultConfigID(C.longlong(configID))
	if result != noError {
		err = client.newError(ctx, 4008, configID, result)
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfigmanager) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szconfigmanager.IDMessages, baseCallerSkip)
	}
	return client.logger
}

// Get the Messenger singleton.
func (client *Szconfigmanager) getMessenger() messenger.Messenger {
	if client.messenger == nil {
		client.messenger = helper.GetMessenger(ComponentID, szconfigmanager.IDMessages, baseCallerSkip)
	}
	return client.messenger
}

// Trace method entry.
func (client *Szconfigmanager) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szconfigmanager) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szconfigmanager) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
	defer func() { client.panicOnError(client.clearLastException(ctx)) }()
	lastExceptionCode, _ := client.getLastExceptionCode(ctx)
	lastException, err := client.getLastException(ctx)
	if err != nil {
		lastException = err.Error()
	}
	details = append(details, messenger.MessageCode{Value: fmt.Sprintf(ExceptionCodeTemplate, lastExceptionCode)})
	details = append(details, messenger.MessageReason{Value: lastException})
	details = append(details, errors.New(lastException))
	errorMessage := client.getMessenger().NewJSON(errorNumber, details...)
	return szerror.New(lastExceptionCode, errorMessage)
}

/*
The panicOnError method calls panic() when an error is not nil.

Input:
  - err: nil or an actual error
*/
func (client *Szconfigmanager) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// --- Sz exception handling --------------------------------------------------

/*
The clearLastException method erases the last exception message held by the Senzing G2ConfigMgr object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfigmanager) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2Config_clearLastException();
	_ = ctx
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	C.G2ConfigMgr_clearLastException()
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's G2ConfigMgr.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2ConfigMgr.
*/
func (client *Szconfigmanager) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	_ = ctx
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13)
		defer func() { client.traceExit(14, result, err, time.Since(entryTime)) }()
	}
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2ConfigMgr_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	result = string(bytes.Trim(stringBuffer, "\x00"))
	return result, err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2ConfigMgr.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2ConfigMgr.
*/
func (client *Szconfigmanager) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Config_getLastExceptionCode();
	_ = ctx
	var err error
	var result int
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(15)
		defer func() { client.traceExit(16, result, err, time.Since(entryTime)) }()
	}
	result = int(C.G2ConfigMgr_getLastExceptionCode())
	return result, err
}

// --- Misc -------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (client *Szconfigmanager) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *Szconfigmanager) getByteArray(size int) []byte {
	return make([]byte, size)
}
