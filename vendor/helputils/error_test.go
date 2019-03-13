package helputils

import (
	"testing"

	"github.com/go-errors/errors"
)

// TestPanic is just a testing harness to test the RecoverError function (test without verification)
func DontTestPanic(t *testing.T) {
	if err := panicFunction(); err != nil {
		t.Error(err)
		// print stacktrace if its an errors.Error
		eerr, ok := err.(*errors.Error)
		if ok {
			t.Log(eerr.ErrorStack())
		}
	}
}

func panicFunction() (err error) {
	println("in panic fn")
	defer RecoverError(err)
	panic("someone please call 911")
	return errors.New("error 1")
}
