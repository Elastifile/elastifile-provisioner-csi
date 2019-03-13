package leafnodes

import (
	"fmt"
	"time"

	"tester/internal/failer"
	"tester/internal/logging"
	"tester/types"
)

type ItNode struct {
	runner *runner

	flag types.FlagType
	text string
}

func NewItNode(text string, body interface{}, flag types.FlagType, codeLocation types.CodeLocation, timeout time.Duration, failer *failer.Failer, componentIndex int) *ItNode {
	logging.Log().Debug(fmt.Sprintf("New It(\"%s\")", text), "flag", flag, "timeout", timeout, "level", componentIndex)

	funcBody := body

	if timeout > 0 {
		bodyType := mustFuncType(body, codeLocation)

		if bodyType.NumIn() == 0 {
			funcBody = func(done chan<- interface{}) {
				body.(func())()
			}
		}
	}

	return &ItNode{
		runner: newRunner(funcBody, codeLocation, timeout, failer, types.SpecComponentTypeIt, componentIndex),
		flag:   flag,
		text:   text,
	}
}

func (node *ItNode) Run() (outcome types.SpecState, failure types.SpecFailure) {
	return node.runner.run()
}

func (node *ItNode) Type() types.SpecComponentType {
	return types.SpecComponentTypeIt
}

func (node *ItNode) Text() string {
	return node.text
}

func (node *ItNode) Flag() types.FlagType {
	return node.flag
}

func (node *ItNode) CodeLocation() types.CodeLocation {
	return node.runner.codeLocation
}

func (node *ItNode) Samples() int {
	return 1
}
