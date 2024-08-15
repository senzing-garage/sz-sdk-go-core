package szdiagnostic

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szdiagnostic package.
szdiagnostic package messages will have the format "SZSDK6003eeee" where "eeee" is the error identifier.

ExceptionCodeTemplate is a template for the error code returned by the Senzing C binary
*/
const (
	ComponentID           = 6003
	ExceptionCodeTemplate = "SENZ%04d"
)
