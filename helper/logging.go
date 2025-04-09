package helper

import (
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

/*
The GetLogger function returns a logger that logs "SZSDKcccceeee" messages.

Input
  - componentID: The 4-digit identifier of the component used as "nnnn" in the "SZSDKcccceeee" message identifier.
  - idMessages: A map of error identifiers to error strings.
  - callerSkip: Default number of "frames" to ascend. See [runtime.Caller].
  - options: Zero or more Option* types from [logging].

Output
  - A configured logger.

[logging]: http://localhost:6060/pkg/github.com/senzing-garage/go-logging/logging/
[runtime.Caller]: https://pkg.go.dev/runtime#Caller
*/
func GetLogger(componentID int, idMessages map[int]string, callerSkip int, options ...interface{}) logging.Logging {
	optionMessageID := fmt.Sprintf("%s%04d", MessageIDPrefix, componentID) + "%04d"
	loggerOptions := []interface{}{
		logging.OptionCallerSkip{Value: callerSkip},
		logging.OptionComponentID{Value: componentID},
		logging.OptionIDMessages{Value: idMessages},
		logging.OptionMessageFields{Value: []string{"id", "text"}},
		logging.OptionMessageIDTemplate{Value: optionMessageID},
	}
	loggerOptions = append(loggerOptions, options...)

	result, err := logging.New(loggerOptions...)
	if err != nil {
		panic(err)
	}

	return result
}
