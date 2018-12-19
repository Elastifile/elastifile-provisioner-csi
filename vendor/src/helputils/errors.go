package helputils

import (
	"fmt"
	"runtime/debug"

	"github.com/go-errors/errors"
)

func RecoverError(err error) {
	if r := recover(); r != nil {
		fmt.Println("Recovered in deferRecover", r)
		// find out exactly what the error was and set err
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case errors.Error:
			err = &x
		default:
			err = errors.New("Unknown panic")
		}
	}
}

func PrintErrorWithStack(err error) {
	switch e := err.(type) {
	case *errors.Error:
		fmt.Println(e.ErrorStack())
	default:
		debug.PrintStack()
	}
}
