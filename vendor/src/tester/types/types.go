package types

import (
	"fmt"
	"time"
)

const TESTER_FOCUS_EXIT_CODE = 197

type SuiteSummary struct {
	SuiteDescription string
	SuiteSucceeded   bool
	SuiteID          string

	NumberOfSpecsBeforeParallelization int
	NumberOfTotalSpecs                 int
	NumberOfSpecsThatWillBeRun         int
	NumberOfPendingSpecs               int
	NumberOfSkippedSpecs               int
	NumberOfPassedSpecs                int
	NumberOfFailedSpecs                int
	NumberOfAbortedSpecs               int
	RunTime                            time.Duration
}

type SpecSummary struct {
	ComponentTexts         []string
	ComponentCodeLocations []CodeLocation

	State           SpecState
	RunTime         time.Duration
	Failure         SpecFailure
	NumberOfSamples int

	CapturedOutput string
	SuiteID        string
}

func (s SpecSummary) HasFailureState() bool {
	return s.State.IsFailure()
}

func (s SpecSummary) TimedOut() bool {
	return s.State == SpecStateTimedOut
}

func (s SpecSummary) Panicked() bool {
	return s.State == SpecStatePanicked
}

func (s SpecSummary) Failed() bool {
	return s.State == SpecStateFailed
}

func (s SpecSummary) Passed() bool {
	return s.State == SpecStatePassed
}

func (s SpecSummary) Skipped() bool {
	return s.State == SpecStateSkipped
}

func (s SpecSummary) Aborted() bool {
	return s.State == SpecStateAborted
}

func (s SpecSummary) Pending() bool {
	return s.State == SpecStatePending
}

type SetupSummary struct {
	ComponentType SpecComponentType
	CodeLocation  CodeLocation

	State   SpecState
	RunTime time.Duration
	Failure SpecFailure

	CapturedOutput string
	SuiteID        string
}

type SpecFailure struct {
	Message        string
	Location       CodeLocation
	ForwardedPanic string

	ComponentIndex        int
	ComponentType         SpecComponentType
	ComponentCodeLocation CodeLocation
}

type SpecState uint

const (
	SpecStateInvalid SpecState = iota

	SpecStatePending
	SpecStateSkipped
	SpecStateAborted
	SpecStatePassed
	SpecStateFailed
	SpecStatePanicked
	SpecStateTimedOut
)

func (state SpecState) IsFailure() bool {
	return state == SpecStateTimedOut || state == SpecStatePanicked || state == SpecStateFailed
}

func (state SpecState) Panicked() bool {
	return state == SpecStatePanicked
}

type SpecComponentType uint

const (
	SpecComponentTypeInvalid SpecComponentType = iota

	SpecComponentTypeContainer
	SpecComponentTypeBeforeSuite
	SpecComponentTypeAfterSuite
	SpecComponentTypeBeforeEach
	SpecComponentTypeJustBeforeEach
	SpecComponentTypeAfterEach
	SpecComponentTypeIt
	SpecComponentTypeMeasure
)

func (styp SpecComponentType) String() string {
	blockType := ""
	switch styp {
	case SpecComponentTypeContainer:
		blockType = "Container"
	case SpecComponentTypeBeforeSuite:
		blockType = "BeforeSuite"
	case SpecComponentTypeAfterSuite:
		blockType = "AfterSuite"
	case SpecComponentTypeBeforeEach:
		blockType = "BeforeEach"
	case SpecComponentTypeJustBeforeEach:
		blockType = "JustBeforeEach"
	case SpecComponentTypeAfterEach:
		blockType = "AfterEach"
	case SpecComponentTypeIt:
		blockType = "It"
	case SpecComponentTypeMeasure:
		blockType = "Measure"
	default:
		panic(fmt.Errorf("invalid SpecComponentType: %d", styp))
	}
	return blockType
}

type FlagType uint

const (
	FlagTypeNone FlagType = iota
	FlagTypeFocused
	FlagTypePending
	FlagType_last = FlagTypePending
)

var flagTypeNames = []string{
	"None",
	"Focused",
	"Pending",
}

func (ft *FlagType) String() string {
	if *ft <= FlagType_last {
		return flagTypeNames[*ft]
	} else {
		return "illegal flag type code"
	}
}
