/*
The [Szdiagnostic] implementation of the [senzing.SzDiagnostic] interface
communicates with the Senzing native C binary, libSz.so.
*/
package szdiagnostic

/*
#include <stdlib.h>
#include "libSzDiagnostic.h"
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

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

/*
Type Szdiagnostic struct implements the [senzing.SzDiagnostic] interface
for communicating with the Senzing C binaries.
*/
type Szdiagnostic struct {
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
// sz-sdk-go.SzDiagnostic interface methods
// ----------------------------------------------------------------------------

/*
Method CheckDatastorePerformance runs performance tests on the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - secondsToRun: Duration of the test in seconds.

Output

  - A JSON document containing performance results.
    Example: `{"numRecordsInserted":0,"insertTime":0}`
*/
func (client *Szdiagnostic) CheckDatastorePerformance(ctx context.Context, secondsToRun int) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(1, secondsToRun)

		entryTime := time.Now()
		defer func() { client.traceExit(2, secondsToRun, result, err, time.Since(entryTime)) }()
	}

	result, err = client.checkDatastorePerformance(ctx, secondsToRun)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szdiagnostic.CheckDatastorePerformance error: %w", err)
}

/*
Method GetDatastoreInfo returns information about the Senzing datastore.

Input
  - ctx: A context to control lifecycle.

Output

  - A JSON document containing Senzing datastore metadata.
*/
func (client *Szdiagnostic) GetDatastoreInfo(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(7)

		entryTime := time.Now()
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getDatastoreInfo(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szdiagnostic.GetDatastoreInfo error: %w", err)
}

/*
Method GetFeature is an experimental method that returns diagnostic information of a feature.
Not recommended for use.

Input
  - ctx: A context to control lifecycle.
  - featureID: The identifier of the feature to describe.

Output

  - A JSON document containing feature metadata.
*/
func (client *Szdiagnostic) GetFeature(ctx context.Context, featureID int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(9, featureID)

		entryTime := time.Now()
		defer func() { client.traceExit(10, featureID, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getFeature(ctx, featureID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"featureID": strconv.FormatInt(featureID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szdiagnostic.GetFeature error: %w", err)
}

/*
WARNING: Method PurgeRepository removes every record in the Senzing datastore.
This is a destructive method that cannot be undone.
Before calling purgeRepository(), all programs using Senzing MUST be terminated.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szdiagnostic) PurgeRepository(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(17)

		entryTime := time.Now()
		defer func() { client.traceExit(18, err, time.Since(entryTime)) }()
	}

	err = client.purgeRepository(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
		}()
	}

	return wraperror.Errorf(err, "szdiagnostic.PurgeRepository error: %w", err)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
Method Destroy will destroy and perform cleanup for the Senzing SzDiagnostic object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szdiagnostic) Destroy(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(5)

		entryTime := time.Now()
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}

	err = client.destroy(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}

	return wraperror.Errorf(err, "szdiagnostic.Destroy error: %w", err)
}

/*
Method GetObserverOrigin returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szdiagnostic) GetObserverOrigin(ctx context.Context) string {
	_ = ctx

	return client.observerOrigin
}

/*
Method Initialize initializes the SzDiagnostic object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.  0 for current default configuration.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szdiagnostic) Initialize(
	ctx context.Context,
	instanceName string,
	settings string,
	configID int64,
	verboseLogging int64,
) error {
	var err error

	if client.isTrace {
		client.traceEntry(15, instanceName, settings, configID, verboseLogging)

		entryTime := time.Now()
		defer func() {
			client.traceExit(16, instanceName, settings, configID, verboseLogging, err, time.Since(entryTime))
		}()
	}

	if configID == senzing.SzInitializeWithDefaultConfiguration {
		err = client.init(ctx, instanceName, settings, verboseLogging)
	} else {
		err = client.initWithConfigID(ctx, instanceName, settings, configID, verboseLogging)
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID":       strconv.FormatInt(configID, baseTen),
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}

	return wraperror.Errorf(err, "szdiagnostic.Initialize error: %w", err)
}

/*
Method RegisterObserver adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szdiagnostic) RegisterObserver(ctx context.Context, observer observer.Observer) error {
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

	return wraperror.Errorf(err, "szdiagnostic.RegisterObserver error: %w", err)
}

/*
Method Reinitialize re-initializes the Senzing SzDiagnostic object.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (client *Szdiagnostic) Reinitialize(ctx context.Context, configID int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(19, configID)

		entryTime := time.Now()
		defer func() { client.traceExit(20, configID, err, time.Since(entryTime)) }()
	}

	err = client.reinit(ctx, configID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}

	return wraperror.Errorf(err, "szdiagnostic.Reinitialize error: %w", err)
}

/*
Method SetLogLevel sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szdiagnostic) SetLogLevel(ctx context.Context, logLevelName string) error {
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

	return wraperror.Errorf(err, "szdiagnostic.SetLogLevel error: %w", err)
}

/*
Method SetObserverOrigin sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szdiagnostic) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
Method UnregisterObserver removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szdiagnostic) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
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

	return wraperror.Errorf(err, "szdiagnostic.UnregisterObserver error: %w", err)
}

// ----------------------------------------------------------------------------
// Private methods for calling the Senzing C API
// ----------------------------------------------------------------------------

func (client *Szdiagnostic) checkDatastorePerformance(ctx context.Context, secondsToRun int) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.SzDiagnostic_checkDatastorePerformance_helper(C.longlong(secondsToRun))
	if result.returnCode != noError {
		err = client.newError(ctx, 4001, secondsToRun, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szdiagnostic) destroy(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.SzDiagnostic_destroy()
	if result != noError {
		err = client.newError(ctx, 4002, result)
	}

	return err
}

func (client *Szdiagnostic) getDatastoreInfo(ctx context.Context) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.SzDiagnostic_getDatastoreInfo_helper()
	if result.returnCode != noError {
		err = client.newError(ctx, 4003, result.returnCode)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

func (client *Szdiagnostic) getFeature(ctx context.Context, featureID int64) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		err            error
		resultResponse string
	)

	result := C.SzDiagnostic_getFeature_helper(C.longlong(featureID))
	if result.returnCode != noError {
		err = client.newError(ctx, 4004, featureID, result)
	}

	resultResponse = C.GoString(result.response)

	C.SzHelper_free(unsafe.Pointer(result.response))

	return resultResponse, err
}

/*
Method init method the Senzing SzDiagnostic object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szdiagnostic) init(
	ctx context.Context,
	instanceName string,
	settings string,
	verboseLogging int64,
) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	instanceNameForC := C.CString(instanceName)

	defer C.free(unsafe.Pointer(instanceNameForC))

	settingsForC := C.CString(settings)

	defer C.free(unsafe.Pointer(settingsForC))

	result := C.SzDiagnostic_init(instanceNameForC, settingsForC, C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4005, instanceName, settings, verboseLogging, result)
	}

	return err
}

/*
Method initWithConfigID initializes the Senzing SzDiagnostic object with a non-default configuration ID.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szdiagnostic) initWithConfigID(
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

	result := C.SzDiagnostic_initWithConfigID(
		instanceNameForC,
		settingsForC,
		C.longlong(configID),
		C.longlong(verboseLogging),
	)
	if result != noError {
		err = client.newError(ctx, 4006, instanceName, settings, configID, verboseLogging, result)
	}

	return err
}

func (client *Szdiagnostic) purgeRepository(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.SzDiagnostic_purgeRepository()
	if result != noError {
		err = client.newError(ctx, 4007, result)
	}

	return err
}

func (client *Szdiagnostic) reinit(ctx context.Context, configID int64) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var err error

	result := C.SzDiagnostic_reinit(C.longlong(configID))
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
func (client *Szdiagnostic) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szdiagnostic.IDMessages, baseCallerSkip)
	}

	return client.logger
}

// Get the Messenger singleton.
func (client *Szdiagnostic) getMessenger() messenger.Messenger {
	if client.messenger == nil {
		client.messenger = helper.GetMessenger(ComponentID, szdiagnostic.IDMessages, baseCallerSkip)
	}

	return client.messenger
}

// Trace method entry.
func (client *Szdiagnostic) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szdiagnostic) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szdiagnostic) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
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
func (client *Szdiagnostic) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// --- Sz exception handling --------------------------------------------------

/*
Method clearLastException erases the last exception message held by the Senzing SzDiagnostic object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szdiagnostic) clearLastException(ctx context.Context) error {
	var err error

	_ = ctx

	if client.isTrace {
		client.traceEntry(3)

		entryTime := time.Now()
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}

	C.SzDiagnostic_clearLastException()

	return err
}

/*
Method getLastException retrieves the last exception thrown in Senzing's SzDiagnostic.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's SzDiagnostic.
*/
func (client *Szdiagnostic) getLastException(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	_ = ctx

	if client.isTrace {
		client.traceEntry(11)

		entryTime := time.Now()
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}

	stringBuffer := client.getByteArray(initialByteArraySize)
	C.SzDiagnostic_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	result = string(bytes.Trim(stringBuffer, "\x00"))

	return result, err
}

/*
Method getLastExceptionCode retrieves the code of the last exception thrown in Senzing's SzDiagnostic.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's SzDiagnostic.
*/
func (client *Szdiagnostic) getLastExceptionCode(ctx context.Context) (int, error) {
	var (
		err    error
		result int
	)

	_ = ctx

	if client.isTrace {
		client.traceEntry(13)

		entryTime := time.Now()
		defer func() { client.traceExit(14, result, err, time.Since(entryTime)) }()
	}

	result = int(C.SzDiagnostic_getLastExceptionCode())

	return result, err
}

// --- Misc -------------------------------------------------------------------

// Get space for an array of bytes of a given size.
// func (client *Szdiagnostic) getByteArrayC(size int) *C.char {
// 	bytes := C.malloc(C.size_t(size))
// 	return (*C.char)(bytes)
// }

// Make a byte array.
func (client *Szdiagnostic) getByteArray(size int) []byte {
	return make([]byte, size)
}

// A hack: Only needed to import the "senzing" package for the godoc comments.
// func junk() {
// 	fmt.Printf(senzing.SzNoAttributes)
// }
