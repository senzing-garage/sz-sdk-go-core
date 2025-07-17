package szconfig

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szconfig package.
Package szconfig messages will have the format "SZSDK6001eeee" where "eeee" is the error identifier.

ExceptionCodeTemplate is a template for the error code returned by the Senzing C binary.
*/
const (
	ComponentID           = 6001
	ExceptionCodeTemplate = "SENZ%04d"
)
