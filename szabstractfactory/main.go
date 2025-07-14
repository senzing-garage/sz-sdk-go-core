package szabstractfactory

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szabstractfactory package.
Package abstractfactory messages will have the format "SZSDK6000eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6000

var errForPackage = errors.New("szabstractfactory")
