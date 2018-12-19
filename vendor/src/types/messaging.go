package types

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/go-errors/errors"
)

type RemoteError struct {
	Message string
	Stack   string
}

func NewRemoteError(err error) *RemoteError {
	var result *RemoteError
	if err != nil {
		var stack string
		if wrappedErr, ok := err.(*errors.Error); ok {
			stack = wrappedErr.ErrorStack()
		} else {
			stack = strings.Join(strings.Split(string(debug.Stack()), "\n")[1:], "\n")
		}

		result = &RemoteError{
			Message: err.Error(),
			Stack:   stack,
		}
	}
	return result
}

func (e *RemoteError) Error() string {
	if e.Stack != "" {
		return fmt.Sprintf("%s\n%s", e.Message, e.Stack)
	}
	return e.Message
}
