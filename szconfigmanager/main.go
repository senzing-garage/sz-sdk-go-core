package szconfigmanager

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szconfigmanager package.
szconfigmanager package messages will have the format "SZSDK6002eeee" where "eeee" is the error identifier.

ExceptionCodeTemplate is a template for the error code returned by the Senzing C binary
*/
const (
	ComponentID           = 6002
	ExceptionCodeTemplate = "SENZ%04d"
)
