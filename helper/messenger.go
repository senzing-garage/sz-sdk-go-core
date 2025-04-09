package helper

import (
	"fmt"

	"github.com/senzing-garage/go-messaging/messenger"
)

/*
The GetMessenger function returns a message generator that creates "SZSDKcccceeee" messages.

Input
  - componentID: The 4-digit identifier of the component used as "nnnn" in the "SZSDKcccceee" message identifier.
  - idMessages: A map of error identifiers to error strings.
  - callerSkip: Default number of "frames" to ascend. See [runtime.Caller].
  - options: Zero or more Option* types from [logging].

Output
  - A configured message generator.

[logging]: http://localhost:6060/pkg/github.com/senzing-garage/go-logging/logging/
[runtime.Caller]: https://pkg.go.dev/runtime#Caller
*/
func GetMessenger(
	componentID int,
	idMessages map[int]string,
	callerSkip int,
	options ...interface{}) messenger.Messenger {
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
