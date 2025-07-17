package szengine

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szengine package.
Package szengine messages will have the format "SZSDK6004eeee" where "eeee" is the error identifier.

ExceptionCodeTemplate is a template for the error code returned by the Senzing C binary.
*/
const (
	ComponentID           = 6004
	ExceptionCodeTemplate = "SENZ%04d"
)

var errForPackage = errors.New("szengine")
