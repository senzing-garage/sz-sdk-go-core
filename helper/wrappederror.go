package helper

import "fmt"

/*
The Errorf function returns a wrapped error, if err != nil.

Input
  - err: The unwrapped/raw error.
  - format: The format string (think fmt.Sprintf())
  - messages: values to be put into the format string.

Output
  - Either nil, or a wrapped error.
*/
func Errorf(err error, format string, messages ...any) error {
	// var result error
	if err != nil {
		return fmt.Errorf(format, messages...) //nolint:err113
	}

	return nil
}
