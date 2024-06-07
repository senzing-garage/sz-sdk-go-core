package helpers

import (
	"github.com/senzing-garage/go-logging/logging"
)

func getLogger(componentID int, idMessages map[int]string, callerSkip int) logging.Logging {
	options := []interface{}{
		&logging.OptionCallerSkip{Value: callerSkip},
	}
	result, err := logging.NewSenzingSdkLogger(componentID, idMessages, options...)
	if err != nil {
		panic(err)
	}
	return result
}
