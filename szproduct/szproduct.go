/*
The [Szproduct] implementation of the [senzing.SzProduct] interface
communicates with the Senzing native C binary, libSz.so.
*/
package szproduct

/*
#include <stdlib.h>
#include "libSzProduct.h"
#include "gohelpers/SzLang_helpers.h"
#cgo CFLAGS: -g -I/opt/senzing/er/sdk/c
#cgo windows CFLAGS: -g -I"C:/Program Files/Senzing/er/sdk/c"
#cgo LDFLAGS: -L/opt/senzing/er/lib -lSz
#cgo windows LDFLAGS: -L"C:/Program Files/Senzing/er/lib" -lSz
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
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/senzing-garage/sz-sdk-go/szproduct"
)

/*
The Szproduct implementation of the [senzing.SzProduct] interface
communicates with the Senzing C binaries.
*/
type Szproduct struct {
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
// sz-sdk-go.SzProduct interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method will destroy and perform cleanup for the Senzing SzProduct object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szproduct) Destroy(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	err = client.destroy(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}
	return err
}

/*
The GetLicense method retrieves information about the license used by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
*/
func (client *Szproduct) GetLicense(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}
	result, err = client.license(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}
	return result, err
}

/*
The GetVersion method returns the Senzing API version information.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
*/
func (client *Szproduct) GetVersion(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	result, err = client.version(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
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
func (client *Szproduct) GetObserverOrigin(ctx context.Context) string {
	_ = ctx
	return client.observerOrigin
}

/*
The Initialize method initializes the Senzing SzProduct object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szproduct) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(14, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	err = client.init(ctx, instanceName, settings, verboseLogging)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
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
func (client *Szproduct) RegisterObserver(ctx context.Context, observer observer.Observer) error {
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
func (client *Szproduct) SetLogLevel(ctx context.Context, logLevelName string) error {
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
func (client *Szproduct) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szproduct) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
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

func (client *Szproduct) destroy(ctx context.Context) error {
	// _DLEXPORT int SzProduct_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	result := C.SzProduct_destroy()
	if result != noError {
		err = client.newError(ctx, 4001, result)
	}
	return err
}

func (client *Szproduct) license(ctx context.Context) (string, error) {
	// _DLEXPORT char* SzProduct_license();
	_ = ctx
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.SzProduct_license()
	resultResponse = C.GoString(result)
	return resultResponse, err
}

func (client *Szproduct) version(ctx context.Context) (string, error) {
	// _DLEXPORT char* SzProduct_license();
	_ = ctx
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	var resultResponse string
	result := C.SzProduct_version()
	resultResponse = C.GoString(result)
	return resultResponse, err
}

func (client *Szproduct) init(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	// _DLEXPORT int SzProduct_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	moduleNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.SzProduct_init(moduleNameForC, iniParamsForC, C.longlong(verboseLogging))
	if result != noError {
		err = client.newError(ctx, 4002, instanceName, settings, verboseLogging, result)
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szproduct) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szproduct.IDMessages, baseCallerSkip)
	}
	return client.logger
}

// Get the Messenger singleton.
func (client *Szproduct) getMessenger() messenger.Messenger {
	if client.messenger == nil {
		client.messenger = helper.GetMessenger(ComponentID, szproduct.IDMessages, baseCallerSkip)
	}
	return client.messenger
}

// Trace method entry.
func (client *Szproduct) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szproduct) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create a new error.
func (client *Szproduct) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
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
func (client *Szproduct) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// --- Sz exception handling --------------------------------------------------

/*
The clearLastException method erases the last exception message held by the Senzing SzProduct object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szproduct) clearLastException(ctx context.Context) error {
	// _DLEXPORT void SzProduct_clearLastException();
	_ = ctx
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1)
		defer func() { client.traceExit(2, err, time.Since(entryTime)) }()
	}
	C.SzProduct_clearLastException()
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's SzProduct.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's SzProduct.
*/
func (client *Szproduct) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int SzProduct_getLastException(char *buffer, const size_t bufSize);
	_ = ctx
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, result, err, time.Since(entryTime)) }()
	}
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.SzProduct_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	result = string(bytes.Trim(stringBuffer, "\x00"))
	return result, err
}

/*
The getLastExceptionCode method retrieves the code of the last exception thrown in Senzing's SzProduct.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's SzProduct.
*/
func (client *Szproduct) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int SzProduct_getLastExceptionCode();
	_ = ctx
	var err error
	var result int
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}
	result = int(C.SzProduct_getLastExceptionCode())
	return result, err
}

// --- Misc -------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (client *Szproduct) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (client *Szproduct) getByteArray(size int) []byte {
	return make([]byte, size)
}

// A hack: Only needed to import the "senzing" package for the godoc comments.
func junk() {
	fmt.Printf(senzing.SzNoAttributes)
}
