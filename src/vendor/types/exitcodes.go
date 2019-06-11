package types

import tester_types "tester/types"

// Tesla exit codes

//go:generate stringer -type=ExitCode

type ExitCode int

const (
	ExitCodeSuccess ExitCode = iota
	ExitCodeTeslaFailed
	ExitCodeToolFailed
	ExitCodeVerifyFailed
	ExitCodeFailed
	ExitCodeUnknown

	ExitCodeTesterFocus ExitCode = tester_types.TESTER_FOCUS_EXIT_CODE
)
