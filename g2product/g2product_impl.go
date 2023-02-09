// The G2productImpl implementation is a wrapper over the Senzing libg2product library.
package g2product

/*
#include "g2product.h"
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

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// G2productImpl is the default implementation of the G2product interface.
type G2productImpl struct {
	isTrace   bool
	logger    messagelogger.MessageLoggerInterface
	observers subject.Subject
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const initialByteArraySize = 65535

// ----------------------------------------------------------------------------
// Internal methods - names begin with lower case
// ----------------------------------------------------------------------------

// Get space for an array of bytes of a given size.
func (g2product *G2productImpl) getByteArrayC(size int) *C.char {
	bytes := C.malloc(C.size_t(size))
	return (*C.char)(bytes)
}

// Make a byte array.
func (g2product *G2productImpl) getByteArray(size int) []byte {
	return make([]byte, size)
}

// Create a new error.
func (g2product *G2productImpl) newError(ctx context.Context, errorNumber int, details ...interface{}) error {
	lastException, err := g2product.getLastException(ctx)
	defer g2product.clearLastException(ctx)
	message := lastException
	if err != nil {
		message = err.Error()
	}

	var newDetails []interface{}
	newDetails = append(newDetails, details...)
	newDetails = append(newDetails, errors.New(message))
	errorMessage, err := g2product.getLogger().Message(errorNumber, newDetails...)
	if err != nil {
		errorMessage = err.Error()
	}

	return errors.New(errorMessage)
}

// Get the Logger singleton.
func (g2product *G2productImpl) getLogger() messagelogger.MessageLoggerInterface {
	if g2product.logger == nil {
		g2product.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return g2product.logger
}

func (g2product *G2productImpl) notify(ctx context.Context, messageId int, err error, details map[string]string) {
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
		g2product.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (g2product *G2productImpl) traceEntry(errorNumber int, details ...interface{}) {
	g2product.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (g2product *G2productImpl) traceExit(errorNumber int, details ...interface{}) {
	g2product.getLogger().Log(errorNumber, details...)
}

/*
The clearLastException method erases the last exception message held by the Senzing G2Config object.

Input
  - ctx: A context to control lifecycle.
*/
func (g2product *G2productImpl) clearLastException(ctx context.Context) error {
	// _DLEXPORT void G2Config_clearLastException();
	if g2product.isTrace {
		g2product.traceEntry(1)
	}
	entryTime := time.Now()
	var err error = nil
	C.G2Product_clearLastException()
	if g2product.isTrace {
		defer g2product.traceExit(2, err, time.Since(entryTime))
	}
	return err
}

/*
The getLastException method retrieves the last exception thrown in Senzing's G2Product.

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing the error received from Senzing's G2Product.
*/
func (g2product *G2productImpl) getLastException(ctx context.Context) (string, error) {
	// _DLEXPORT int G2Config_getLastException(char *buffer, const size_t bufSize);
	if g2product.isTrace {
		g2product.traceEntry(5)
	}
	entryTime := time.Now()
	var err error = nil
	stringBuffer := g2product.getByteArray(initialByteArraySize)
	C.G2Product_getLastException((*C.char)(unsafe.Pointer(&stringBuffer[0])), C.ulong(len(stringBuffer)))
	// if result == 0 { // "result" is length of exception message.
	// 	err = g2product.getLogger().Error(4002, result, time.Since(entryTime))
	// }
	stringBuffer = bytes.Trim(stringBuffer, "\x00")
	if g2product.isTrace {
		defer g2product.traceExit(6, string(stringBuffer), err, time.Since(entryTime))
	}
	return string(stringBuffer), err
}

/*
The GetLastExceptionCode method retrieves the code of the last exception thrown in Senzing's G2Product.

Input:
  - ctx: A context to control lifecycle.

Output:
  - An int containing the error received from Senzing's G2Product.
*/
func (g2product *G2productImpl) getLastExceptionCode(ctx context.Context) (int, error) {
	//  _DLEXPORT int G2Config_getLastExceptionCode();
	if g2product.isTrace {
		g2product.traceEntry(7)
	}
	entryTime := time.Now()
	var err error = nil
	result := int(C.G2Product_getLastExceptionCode())
	if g2product.isTrace {
		defer g2product.traceExit(8, result, err, time.Since(entryTime))
	}
	return result, err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Product object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (g2product *G2productImpl) Destroy(ctx context.Context) error {
	// _DLEXPORT int G2Config_destroy();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(3)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Product_destroy()
	if result != 0 {
		err = g2product.newError(ctx, 4001, result, time.Since(entryTime))
	}
	if g2product.observers != nil {
		go func() {
			details := map[string]string{}
			g2product.notify(ctx, 8001, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(4, err, time.Since(entryTime))
	}
	return err
}

/*
The Init method initializes the Senzing G2Product object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (g2product *G2productImpl) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	// _DLEXPORT int G2Config_init(const char *moduleName, const char *iniParams, const int verboseLogging);
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(9, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	var err error = nil
	moduleNameForC := C.CString(moduleName)
	defer C.free(unsafe.Pointer(moduleNameForC))
	iniParamsForC := C.CString(iniParams)
	defer C.free(unsafe.Pointer(iniParamsForC))
	result := C.G2Product_init(moduleNameForC, iniParamsForC, C.int(verboseLogging))
	if result != 0 {
		err = g2product.newError(ctx, 4003, moduleName, iniParams, verboseLogging, result, time.Since(entryTime))
	}
	if g2product.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			g2product.notify(ctx, 8002, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(10, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The License method retrieves information about the currently used license by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
    See the example output.
*/
func (g2product *G2productImpl) License(ctx context.Context) (string, error) {
	// _DLEXPORT char* G2Product_license();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(11)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Product_license()
	if g2product.observers != nil {
		go func() {
			details := map[string]string{}
			g2product.notify(ctx, 8003, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(12, C.GoString(result), err, time.Since(entryTime))
	}
	return C.GoString(result), err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (g2product *G2productImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if g2product.observers == nil {
		g2product.observers = &subject.SubjectImpl{}
	}
	return g2product.observers.RegisterObserver(ctx, observer)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (g2product *G2productImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(13, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2product.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	g2product.isTrace = (g2product.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if g2product.isTrace {
		defer g2product.traceExit(14, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (g2product *G2productImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	err := g2product.observers.UnregisterObserver(ctx, observer)
	if err != nil {
		return err
	}
	if !g2product.observers.HasObservers(ctx) {
		g2product.observers = nil
	}
	return err
}

/*
The ValidateLicenseFile method validates the licence file has not expired.

Input
  - ctx: A context to control lifecycle.
  - licenseFilePath: A fully qualified path to the Senzing license file.

Output
  - if error is nil, license is valid.
  - If error not nil, license is not valid.
  - The returned string has additional information.
*/
func (g2product *G2productImpl) ValidateLicenseFile(ctx context.Context, licenseFilePath string) (string, error) {
	// _DLEXPORT int G2Product_validateLicenseFile(const char* licenseFilePath, char **errorBuf, size_t *errorBufSize, void *(*resizeFunc)(void *ptr,size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(15, licenseFilePath)
	}
	entryTime := time.Now()
	var err error = nil
	licenseFilePathForC := C.CString(licenseFilePath)
	defer C.free(unsafe.Pointer(licenseFilePathForC))
	result := C.G2Product_validateLicenseFile_helper(licenseFilePathForC)
	if result.returnCode != 0 {
		err = g2product.newError(ctx, 4004, licenseFilePath, result.returnCode, result, time.Since(entryTime))
	}
	if g2product.observers != nil {
		go func() {
			details := map[string]string{}
			g2product.notify(ctx, 8004, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(16, licenseFilePath, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The ValidateLicenseStringBase64 method validates the licence, represented by a Base-64 string, has not expired.

Input
  - ctx: A context to control lifecycle.
  - licenseString: A Senzing license represented by a Base-64 encoded string.

Output
  - if error is nil, license is valid.
  - If error not nil, license is not valid.
  - The returned string has additional information.
    See the example output.
*/
func (g2product *G2productImpl) ValidateLicenseStringBase64(ctx context.Context, licenseString string) (string, error) {
	// _DLEXPORT int G2Product_validateLicenseStringBase64(const char* licenseString, char **errorBuf, size_t *errorBufSize, void *(*resizeFunc)(void *ptr,size_t newSize));
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(17, licenseString)
	}
	entryTime := time.Now()
	var err error = nil
	licenseStringForC := C.CString(licenseString)
	defer C.free(unsafe.Pointer(licenseStringForC))
	result := C.G2Product_validateLicenseStringBase64_helper(licenseStringForC)
	if result.returnCode != 0 {
		err = g2product.newError(ctx, 4005, licenseString, result.returnCode, result, time.Since(entryTime))
	}
	if g2product.observers != nil {
		go func() {
			details := map[string]string{}
			g2product.notify(ctx, 8005, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(18, licenseString, C.GoString(result.response), err, time.Since(entryTime))
	}
	return C.GoString(result.response), err
}

/*
The Version method returns the version of the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
    See the example output.
*/
func (g2product *G2productImpl) Version(ctx context.Context) (string, error) {
	// _DLEXPORT char* G2Product_license();
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if g2product.isTrace {
		g2product.traceEntry(19)
	}
	entryTime := time.Now()
	var err error = nil
	result := C.G2Product_version()
	if g2product.observers != nil {
		go func() {
			details := map[string]string{}
			g2product.notify(ctx, 8006, err, details)
		}()
	}
	if g2product.isTrace {
		defer g2product.traceExit(20, C.GoString(result), err, time.Since(entryTime))
	}
	return C.GoString(result), err
}
