/*
The szconfig implementation is a wrapper over the Senzing libg2config library.
*/
package szconfig

/*
#include <stdlib.h>
#include "libg2config.h"
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
	"github.com/senzing-garage/sz-sdk-go-core/helpers"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

type Szconfig struct {
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
// sz-sdk-go.SzConfig interface methods
// ----------------------------------------------------------------------------

/*
The AddDataSource method adds a data source to an existing in-memory configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - dataSourceCode: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.

Output
  - A string containing a JSON document listing the newly created data source.
    See the example output.
*/
func (client *Szconfig) AddDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, configHandle, dataSourceCode)
		defer func() {
			client.traceExit(2, configHandle, dataSourceCode, result, err, time.Since(entryTime))
		}()
	}
	result, err = client.addDataSource(ctx, configHandle, dataSourceCode)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"return":         result,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}
	return result, err
}

/*
The CloseConfig method cleans up the Senzing G2Config object pointed to by the handle.
The handle was created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
*/
func (client *Szconfig) CloseConfig(ctx context.Context, configHandle uintptr) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5, configHandle)
		defer func() { client.traceExit(6, configHandle, err, time.Since(entryTime)) }()
	}
	err = client.close(ctx, configHandle)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}
	return err
}

/*
The CreateConfig method creates an in-memory Senzing configuration from the g2config.json
template configuration file located in the PIPELINE.RESOURCEPATH path.
A handle is returned to identify the in-memory configuration.
The handle is used by the AddDataSource(), ListDataSources(), DeleteDataSource(), and Save() methods.
The handle is terminated by the Close() method.

Input
  - ctx: A context to control lifecycle.

Output
  - A Pointer to an in-memory Senzing configuration.
*/
func (client *Szconfig) CreateConfig(ctx context.Context) (uintptr, error) {
	var err error
	var result uintptr
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}
	result, err = client.create(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}
	return result, err
}

/*
The DeleteDataSource method removes a data source from an existing configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - dataSourceCode: The datasource name (e.g. "TEST_DATASOURCE").
*/
func (client *Szconfig) DeleteDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9, configHandle, dataSourceCode)
		defer func() { client.traceExit(10, configHandle, dataSourceCode, err, time.Since(entryTime)) }()
	}
	err = client.deleteDataSource(ctx, configHandle, dataSourceCode)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing Szconfig object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfig) Destroy(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}
	err = client.destroy(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}
	return err
}

/*
The ExportConfig method creates a JSON string representation of the Senzing Szconfig object.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON Document representation of the Senzing Szconfig object.
    See the example output.
*/
func (client *Szconfig) ExportConfig(ctx context.Context, configHandle uintptr) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13, configHandle)
		defer func() { client.traceExit(14, configHandle, result, err, time.Since(entryTime)) }()
	}
	result, err = client.save(ctx, configHandle)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}
	return result, err
}

/*
The GetDataSources method returns a JSON document of data sources.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON document listing all of the data sources.
    See the example output.
*/
func (client *Szconfig) GetDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(15, configHandle)
		defer func() { client.traceExit(16, configHandle, result, err, time.Since(entryTime)) }()
	}
	result, err = client.listDataSources(ctx, configHandle)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}
	return result, err
}

/*
The ImportConfig method initializes the in-memory Senzing G2Config object from a JSON string.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: A JSON document containing the Senzing configuration.

Output
  - An identifier of an in-memory configuration.
*/
func (client *Szconfig) ImportConfig(ctx context.Context, configDefinition string) (uintptr, error) {
	var err error
	var result uintptr
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, configDefinition)
		defer func() { client.traceExit(22, configDefinition, result, err, time.Since(entryTime)) }()
	}
	result, err = client.load(ctx, configDefinition)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
		}()
	}
	return result, err
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
func (client *Szconfig) GetObserverOrigin(ctx context.Context) string {
	_ = ctx
	return client.observerOrigin
}

/*
The Initialize method initializes the Senzing Szconfig object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfig) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(24, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	err = client.init(ctx, instanceName, settings, verboseLogging)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
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
func (client *Szconfig) RegisterObserver(ctx context.Context, observer observer.Observer) error {
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
func (client *Szconfig) SetLogLevel(ctx context.Context, logLevelName string) error {
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
func (client *Szconfig) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfig) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
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

func (client *Szconfig) addDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) (string, error) {
	// _DLEXPORT int G2Config_addDataSource(ConfigHandle configHandle, const char *inputJson, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	dataSourceDefinition := `{"DSRC_CODE": "` + dataSourceCode + `"}`
	dataSourceDefinitionForC := C.CString(dataSourceDefinition)
	defer C.free(unsafe.Pointer(dataSourceDefinitionForC))
	result := C.G2Config_addDataSource_helper(C.uintptr_t(configHandle), dataSourceDefinitionForC)
	if result.returnCode != noError {
		err = client.newError(ctx, 4001, configHandle, dataSourceCode, result.returnCode, result)
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	return resultResponse, err
}

func (client *Szconfig) close(ctx context.Context, configHandle uintptr) error {
	// _DLEXPORT int G2Config_close(ConfigHandle configHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.G2Config_close_helper(C.uintptr_t(configHandle))
	if result != noError {
		err = client.newError(ctx, 4002, configHandle, result)
	}
	return err
}

func (client *Szconfig) create(ctx context.Context) (uintptr, error) {
	// _DLEXPORT int G2Config_create(ConfigHandle* configHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse uintptr
	result := C.G2Config_create_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4003, result.returnCode)
	}
	resultResponse = uintptr(result.response)
	return resultResponse, err
}

func (client *Szconfig) deleteDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) error {
	// _DLEXPORT int G2Config_deleteDataSource(ConfigHandle configHandle, const char *inputJson);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	dataSourceDefinition := `{"DSRC_CODE": "` + dataSourceCode + `"}`
	dataSourceDefinitionForC := C.CString(dataSourceDefinition)
	defer C.free(unsafe.Pointer(dataSourceDefinitionForC))
	result := C.G2Config_deleteDataSource_helper(C.uintptr_t(configHandle), dataSourceDefinitionForC)
	if result != noError {
		err = client.newError(ctx, 4004, configHandle, dataSourceCode, result)
	}
	return err
}

func (client *Szconfig) destroy(ctx context.Context) error {
	// _DLEXPORT int G2Config_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.G2Config_destroy()
	if result != noError {
		err = client.newError(ctx, 4005, result)
	}
	return err
}

func (client *Szconfig) save(ctx context.Context, configHandle uintptr) (string, error) {
	// _DLEXPORT int G2Config_save(ConfigHandle configHandle, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize) );
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.G2Config_save_helper(C.uintptr_t(configHandle))
	if result.returnCode != noError {
		err = client.newError(ctx, 4010, configHandle, result.returnCode, result)
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	return resultResponse, err
}

func (client *Szconfig) listDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	// _DLEXPORT int G2Config_listDataSources(ConfigHandle configHandle, char **responseBuf, size_t *bufSize, void *(*resizeFunc)(void *ptr, size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.G2Config_listDataSources_helper(C.uintptr_t(configHandle))
	if result.returnCode != noError {
		err = client.newError(ctx, 4008, result.returnCode, result)
	}
	resultResponse = C.GoString(result.response)
	C.G2GoHelper_free(unsafe.Pointer(result.response))
	return resultResponse, err
}

func (client *Szconfig) load(ctx context.Context, configDefinition string) (uintptr, error) {
	// _DLEXPORT int G2Config_load(const char *jsonConfig,ConfigHandle* configHandle);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse uintptr
	jsonConfigForC := C.CString(configDefinition)
	defer C.free(unsafe.Pointer(jsonConfigForC))
	result := C.G2Config_load_helper(jsonConfigForC)
	if result.returnCode != noError {
		err = client.newError(ctx, 4009, configDefinition, result.returnCode)
	}
	resultResponse = uintptr(result.response)
	return resultResponse, err
}

func (client *Szconfig) init(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	// _DLEXPORT int G2Config_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	instanceNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(instanceNameForC))
	settingsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(settingsForC))
	result := C.G2Config_init(instanceNameForC, settingsForC, C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4007, instanceName, settings, verboseLogging, result)
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfig) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helpers.GetLogger(ComponentID, szconfig.IDMessages, baseCallerSkip)
	}
	return client.logger
}

// Get the Messenger singleton.
func (client *Szconfig) getMessenger() messenger.Messenger {
	if client.messenger == nil {
		client.messenger = helpers.GetMessenger(ComponentID, szconfig.IDMessages, baseCallerSkip)
	}
	return client.messenger
}

// Trace method entry.
func (client *Szconfig) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szconfig) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szconfig) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
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
func (client *Szconfig) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// --- Sz exception handling --------------------------------------------------

/*
The clearLastException method erases the last exception message held by the Senzing Szconfig object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfig) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2Config_clearLastException();
	_ = ctx
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	C.G2Config_clearLastException()
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's Szconfig.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's Szconfig.
*/
func (client *Szconfig) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	_ = ctx
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17)
		defer func() { client.traceExit(18, result, err, time.Since(entryTime)) }()
	}
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2Config_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	result = string(bytes.Trim(stringBuffer, "\x00"))
	return result, err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's Szconfig.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2Config.
*/
func (client *Szconfig) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Config_getLastExceptionCode();
	_ = ctx
	var err error
	var result int
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19)
		defer func() { client.traceExit(20, result, err, time.Since(entryTime)) }()
	}
	result = int(C.G2Config_getLastExceptionCode())
	return result, err
}

// --- Misc -------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (client *Szconfig) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *Szconfig) getByteArray(size int) []byte {
	return make([]byte, size)
}
