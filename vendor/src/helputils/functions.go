package helputils

import "fmt"

type (
	Action  func() error
	Actions []Action
)

func (as Actions) Call() error {
	for i, a := range as {
		if e := a(); e != nil {
			fmt.Printf("failed on command: %d\n", i)
			return e
		}
	}
	return nil
}
