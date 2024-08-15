package szproduct

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szproduct package.
szproduct package messages will have the format "SZSDK6006eeee" where "eeee" is the error identifier.

ExceptionCodeTemplate is a template for the error code returned by the Senzing C binary
*/
const (
	ComponentID           = 6006
	ExceptionCodeTemplate = "SENZ%04d"
)
