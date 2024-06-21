package helper

import (
	"fmt"

	"github.com/senzing-garage/go-messaging/messenger"
)

func GetMessenger(componentID int, idMessages map[int]string, callerSkip int, options ...interface{}) messenger.Messenger {
	optionMessageIDTemplate := fmt.Sprintf("%s%04d", MessageIDPrefix, componentID) + "%04d"
	messengerOptions := []interface{}{
		messenger.OptionCallerSkip{Value: callerSkip},
		messenger.OptionIDMessages{Value: idMessages},
		messenger.OptionMessageFields{Value: []string{"id", "reason"}},
		messenger.OptionMessageIDTemplate{Value: optionMessageIDTemplate},
	}
	messengerOptions = append(messengerOptions, options...)
	result, err := messenger.New(messengerOptions...)
	if err != nil {
		panic(err)
	}
	return result
}
