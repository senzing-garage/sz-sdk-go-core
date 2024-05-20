/*
The Szproduct implementation is a wrapper over the Senzing libg2product library.
*/
package szproduct

/*
#include <stdlib.h>
#include "libg2product.h"
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
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szproductapi "github.com/senzing-garage/sz-sdk-go/szproduct"
)

type Szproduct struct {
	isTrace        bool
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

const initialByteArraySize = 65535

// ----------------------------------------------------------------------------
// sz-sdk-go.SzProduct interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Product object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szproduct) Destroy(ctx context.Context) error {
	// _DLEXPORT int G2Config_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	result := C.G2Product_destroy()
	if result != 0 {
		err = client.newError(ctx, 4001, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return err
}

/*
The GetLicense method retrieves information about the currently used license by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
    See the example output.
*/
func (client *Szproduct) GetLicense(ctx context.Context) (string, error) {
	// _DLEXPORT char* G2Product_license();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2Product_license()
	resultResponse = C.GoString(result)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8002, err, details)
		}()
	}
	return resultResponse, err
}

/*
The GetVersion method returns the version of the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
    See the example output.
*/
func (client *Szproduct) GetVersion(ctx context.Context) (string, error) {
	// _DLEXPORT char* G2Product_license();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	var resultResponse string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, resultResponse, err, time.Since(entryTime)) }()
	}
	result := C.G2Product_version()
	resultResponse = C.GoString(result)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8003, err, details)
		}()
	}
	return resultResponse, err
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
	return client.observerOrigin
}

/*
The Initialize method initializes the Senzing SzProduct object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szproduct) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	// _DLEXPORT int G2Config_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
	entryTime := time.Now()
	if client.isTrace {
		client.traceEntry(13, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(14, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	moduleNameForC := C.CString(instanceName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(settings)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2Product_init(moduleNameForC, iniParamsForC, C.longlong(verboseLogging))
	if result != 0 {
		err = client.newError(ctx, 4003, instanceName, settings, verboseLogging, result, time.Since(entryTime))
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8004, err, details)
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
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(703, observer.GetObserverId(ctx))
		defer func() { client.traceExit(704, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8702, err, details)
		}()
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szproduct) SetLogLevel(ctx context.Context, logLevelName string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error = nil
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8703, err, details)
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
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szproduct) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(707, observer.GetObserverId(ctx))
		defer func() { client.traceExit(708, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerId": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8704, err, details)
		err = client.observers.UnregisterObserver(ctx, observer)
		if !client.observers.HasObservers(ctx) {
			client.observers = nil
		}
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szproduct) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, szproductapi.IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
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
	lastException, err := client.getLastException(ctx)
	defer client.clearLastException(ctx)
	message := lastException
	if err != nil {
		message = err.Error()
	}
	details = append(details, errors.New(message))
	errorMessage := client.getLogger().Json(errorNumber, details...)

	// TODO: Remove hack

	code := szerror.Code(message) // hack
	if code > 30000 {             // hack
		code = code - 27000 // hack
	} // hack

	return szerror.New(code, (errorMessage))
}

// --- Sz exception handling --------------------------------------------------

/*
The clearLastException method erases the last exception message held by the Senzing G2Product object.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szproduct) clearLastException(ctx context.Context) error {
	_ = ctx
	// _DLEXPORT void G2Config_clearLastException();
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1)
		defer func() { client.traceExit(2, err, time.Since(entryTime)) }()
	}
	C.G2Product_clearLastException()
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's client.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2Product.
*/
func (client *Szproduct) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	_ = ctx
	var err error = nil
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, result, err, time.Since(entryTime)) }()
	}
	stringBuffer := client.getByteArray(initialByteArraySize)
	C.G2Product_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.size_t(len(stringBuffer)))
	// if result == 0 { // "result" is length of exception message.
	// 	err = client.getLogger().Error(4002, result, time.Since(entryTime))
	// }
	result = string(bytes.Trim(stringBuffer, "\x00"))
	return result, err
}

/*
The GetLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2Product.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2Product.
*/
func (client *Szproduct) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Config_getLastExceptionCode();
	_ = ctx
	var err error = nil
	var result int
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}
	result = int(C.G2Product_getLastExceptionCode())
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
