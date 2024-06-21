package helper

import (
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

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
