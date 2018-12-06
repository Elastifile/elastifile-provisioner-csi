/*
The stenographer is used by Ginkgo's reporters to generate output.

Move along, nothing to see here.
*/

package stenographer

import (
	"fmt"
	"runtime"
	"strings"

	tm "github.com/buger/goterm"

	"helputils"
	"tester/internal/logging"
	"tester/types"
)

const (
	defaultStyle = "\x1b[0m"
	boldStyle    = "\x1b[1m"
	grayColor    = "\x1b[37m"
)

var (
	redColor    = fmt.Sprintf("\x1b[3%dm", tm.RED)
	greenColor  = fmt.Sprintf("\x1b[3%dm", tm.GREEN)
	yellowColor = fmt.Sprintf("\x1b[3%dm", tm.YELLOW)
	cyanColor   = fmt.Sprintf("\x1b[3%dm", tm.CYAN)
	blueColor   = fmt.Sprintf("\x1b[3%dm", tm.BLUE)
)

type cursorStateType int

const (
	cursorStateTop cursorStateType = iota
	cursorStateStreaming
	cursorStateMidBlock
	cursorStateEndBlock
)

type Stenographer interface {
	AnnounceSuite(description string, randomSeed int64, randomizingAll bool, succinct bool)
	AnnounceAggregatedParallelRun(nodes int, succinct bool)
	AnnounceParallelRun(node int, nodes int, specsToRun int, totalSpecs int, succinct bool)
	AnnounceNumberOfSpecs(specsToRun int, total int, succinct bool)
	AnnounceSpecRunCompletion(summary *types.SuiteSummary, succinct bool)

	AnnounceSpecWillRun(spec *types.SpecSummary)
	AnnounceBeforeSuiteFailure(summary *types.SetupSummary, succinct bool, fullTrace bool)
	AnnounceAfterSuiteFailure(summary *types.SetupSummary, succinct bool, fullTrace bool)

	AnnounceCapturedOutput(output string)

	AnnounceSuccesfulSpec(spec *types.SpecSummary)
	AnnounceSuccesfulSlowSpec(spec *types.SpecSummary, succinct bool)

	AnnouncePendingSpec(spec *types.SpecSummary, noisy bool)
	AnnounceSkippedSpec(spec *types.SpecSummary, succinct bool, fullTrace bool)

	AnnounceSpecTimedOut(spec *types.SpecSummary, succinct bool, fullTrace bool)
	AnnounceSpecPanicked(spec *types.SpecSummary, succinct bool, fullTrace bool)
	AnnounceSpecFailed(spec *types.SpecSummary, succinct bool, fullTrace bool)

	SummarizeFailures(summaries []*types.SpecSummary)
	SummarizeAll(summaries []*types.SpecSummary)
}

func New(color bool) Stenographer {
	denoter := "•"
	if runtime.GOOS == "windows" {
		denoter = "+"
	}
	return &consoleStenographer{
		color:       color,
		denoter:     denoter,
		cursorState: cursorStateTop,
	}
}

type consoleStenographer struct {
	color       bool
	denoter     string
	cursorState cursorStateType
}

var alternatingColors = []string{defaultStyle, grayColor}

func (s *consoleStenographer) AnnounceSuite(description string, randomSeed int64, randomizingAll bool, succinct bool) {
	if succinct {
		s.print(0, "[%d] %s ", randomSeed, s.colorize(boldStyle, description))
		return
	}
	s.printBanner(fmt.Sprintf("Running Suite: %s", description), "=")
	s.print(0, "Random Seed: %s", s.colorize(boldStyle, "%d", randomSeed))
	if randomizingAll {
		s.print(0, " - Will randomize all specs")
	}
	s.flushLine()
}

func (s *consoleStenographer) AnnounceParallelRun(node int, nodes int, specsToRun int, totalSpecs int, succinct bool) {
	if succinct {
		s.print(0, "- node #%d ", node)
		return
	}
	s.println(0,
		"Parallel test node %s/%s. Assigned %s of %s specs.",
		s.colorize(boldStyle, "%d", node),
		s.colorize(boldStyle, "%d", nodes),
		s.colorize(boldStyle, "%d", specsToRun),
		s.colorize(boldStyle, "%d", totalSpecs),
	)
	s.flushLine()
}

func (s *consoleStenographer) AnnounceAggregatedParallelRun(nodes int, succinct bool) {
	if succinct {
		s.print(0, "- %d nodes ", nodes)
		return
	}
	s.println(0,
		"Running in parallel across %s nodes",
		s.colorize(boldStyle, "%d", nodes),
	)
	s.flushLine()
}

func (s *consoleStenographer) AnnounceNumberOfSpecs(specsToRun int, total int, succinct bool) {
	if succinct {
		s.print(0, "- %d/%d specs ", specsToRun, total)
		s.stream()
		return
	}

	logging.Log().Info(fmt.Sprintf("Will run %d of %d specs", specsToRun, total))
}

func (s *consoleStenographer) AnnounceSpecRunCompletion(summary *types.SuiteSummary, succinct bool) {
	if succinct && summary.SuiteSucceeded {
		s.print(0, " %s %s ", s.colorize(greenColor, "SUCCESS!"), summary.RunTime)
		return
	}
	s.flushLine()
	color := greenColor
	if !summary.SuiteSucceeded {
		color = redColor
	}

	status, log := "PASS!", logging.Log().Info
	if !summary.SuiteSucceeded {
		status, log = "FAIL!", logging.Log().Warn
	}

	log(s.colorize(boldStyle+color, "%s Ran %d tests in %v", status, summary.NumberOfSpecsThatWillBeRun, helputils.RoundSeconds(summary.RunTime)))

	colorOrGray := func(color string, count int) string {
		if count == 0 {
			return grayColor
		} else {
			return color
		}
	}

	logging.Log().Info(fmt.Sprintf(
		"Summary: %s | %s | %s | %s | %s - %s=%s",
		s.colorize(colorOrGray(greenColor, summary.NumberOfPassedSpecs), "%d Passed", summary.NumberOfPassedSpecs),
		s.colorize(colorOrGray(redColor, summary.NumberOfFailedSpecs), "%d Failed", summary.NumberOfFailedSpecs),
		s.colorize(colorOrGray(yellowColor, summary.NumberOfAbortedSpecs), "%d Aborted", summary.NumberOfAbortedSpecs),
		s.colorize(colorOrGray(cyanColor, summary.NumberOfSkippedSpecs), "%d Skipped", summary.NumberOfSkippedSpecs),
		s.colorize(boldStyle, "%d Total", summary.NumberOfTotalSpecs),
		s.colorize(colorOrGray(blueColor, summary.NumberOfSkippedSpecs), "runtime"),
		helputils.RoundSeconds(summary.RunTime),
	))
}

func (s *consoleStenographer) AnnounceSpecWillRun(spec *types.SpecSummary) {
	if spec.Aborted() {
		s.announceAbort(spec, "[-] Abort")
	} else {
		s.announceSpec(spec, ">>> Run")
	}
	last := len(spec.ComponentTexts) - 1
	s.println(1, spec.ComponentCodeLocations[last].String())
}

func (s *consoleStenographer) announceSpec(spec *types.SpecSummary, msg string, args ...interface{}) {
	logging.Log().Info(s.announceMessage(spec, msg), args...)
}

func (s *consoleStenographer) announceAbort(spec *types.SpecSummary, msg string) {
	logging.Log().Warn(s.announceMessage(spec, msg))
}

func (s *consoleStenographer) announceMessage(spec *types.SpecSummary, msg string) string {
	last := len(spec.ComponentTexts) - 1
	return fmt.Sprintf("%s %s %s",
		s.colorize(boldStyle, msg),
		strings.Join(spec.ComponentTexts[1:last], " "),
		s.colorize(boldStyle, spec.ComponentTexts[last]),
	)
}

func (s *consoleStenographer) AnnounceBeforeSuiteFailure(summary *types.SetupSummary, succinct bool, fullTrace bool) {
	s.announceSetupFailure("BeforeSuite", summary, succinct, fullTrace)
}

func (s *consoleStenographer) AnnounceAfterSuiteFailure(summary *types.SetupSummary, succinct bool, fullTrace bool) {
	s.announceSetupFailure("AfterSuite", summary, succinct, fullTrace)
}

func (s *consoleStenographer) announceSetupFailure(name string, summary *types.SetupSummary, succinct bool, fullTrace bool) {
	s.startBlock()

	log := logging.Log().Info
	var message string
	switch summary.State {
	case types.SpecStateFailed:
		message = "Failure"
		log = logging.Log().Error
	case types.SpecStatePanicked:
		message = "Panic"
		log = logging.Log().Crit
	case types.SpecStateTimedOut:
		message = "Timeout"
		log = logging.Log().Error
	}

	log(s.colorize(redColor+boldStyle, "%s [%.3f seconds]", message, summary.RunTime.Seconds()))

	indentation := s.printCodeLocationBlock([]string{name}, []types.CodeLocation{summary.CodeLocation}, summary.ComponentType, 0, summary.State, true)

	s.flushLine()
	log(strings.Repeat(" ", indentation), summary.State, summary.Failure, fullTrace)

	s.endBlock()
}

func (s *consoleStenographer) AnnounceCapturedOutput(output string) {
	if output == "" {
		return
	}

	s.startBlock()
	s.println(0, output)
	s.midBlock()
}

func (s *consoleStenographer) AnnounceSuccesfulSpec(spec *types.SpecSummary) {
	s.println(0, s.colorize(greenColor, s.denoter))
	s.stream()
	s.announceSpec(spec, "[v] Passed", "runtime", helputils.RoundSeconds(spec.RunTime))
}

func (s *consoleStenographer) AnnounceSuccesfulSlowSpec(spec *types.SpecSummary, succinct bool) {
	s.printBlockWithMessage(
		s.colorize(greenColor, "%s [SLOW TEST:%.3f seconds]", s.denoter, spec.RunTime.Seconds()),
		"",
		spec,
		succinct,
	)
	s.announceSpec(spec, "[v] Passed", "runtime", helputils.RoundSeconds(spec.RunTime))
}

func (s *consoleStenographer) AnnouncePendingSpec(spec *types.SpecSummary, noisy bool) {
	if noisy {
		s.printBlockWithMessage(
			s.colorize(yellowColor, "P [PENDING]"),
			"",
			spec,
			false,
		)
	} else {
		s.print(0, s.colorize(yellowColor, "P"))
		s.stream()
	}
}

func (s *consoleStenographer) AnnounceSkippedSpec(spec *types.SpecSummary, succinct bool, fullTrace bool) {
	// Skips at runtime will have a non-empty spec.Failure. All others should be succinct.
	if succinct || spec.Failure == (types.SpecFailure{}) {
		s.print(0, s.colorize(cyanColor, "S"))
		s.stream()
	} else {
		s.startBlock()
		s.printSkip(0, spec)
		_ = s.printCodeLocationBlock(spec.ComponentTexts, spec.ComponentCodeLocations, spec.Failure.ComponentType, spec.Failure.ComponentIndex, spec.State, succinct)
		s.endBlock()
	}
}

func (s *consoleStenographer) AnnounceSpecTimedOut(spec *types.SpecSummary, succinct bool, fullTrace bool) {
	s.printSpecFailure(fmt.Sprintf("%s... Timeout", s.denoter), spec, succinct, fullTrace)
}

func (s *consoleStenographer) AnnounceSpecPanicked(spec *types.SpecSummary, succinct bool, fullTrace bool) {
	s.printSpecFailure(fmt.Sprintf("%s! Panic", s.denoter), spec, succinct, fullTrace)
}

func (s *consoleStenographer) AnnounceSpecFailed(spec *types.SpecSummary, succinct bool, fullTrace bool) {
	s.printSpecFailure(fmt.Sprintf("%s Failure", s.denoter), spec, succinct, fullTrace)
}

func (s *consoleStenographer) SummarizeAll(summaries []*types.SpecSummary) {
	s.flushLine()

	colorMap := map[rune]string{
		'v': greenColor,
		'x': redColor,
		'-': yellowColor,
		'~': cyanColor,
	}

	titleMap := map[rune]string{
		'v': "Passed",
		'x': "Failed",
		'-': "Aborted",
		'~': "Skipped",
	}

	results := make(map[rune][]string)

	s.println(0, s.colorize(boldStyle, "Session Results:"))

	for _, summary := range summaries {
		result := 'v'
		if summary.HasFailureState() {
			result = 'x'
		} else if summary.Aborted() {
			result = '-'
		} else if summary.Skipped() {
			result = '~'
		}

		s.print(0, s.colorize(colorMap[result], fmt.Sprintf("[%s] ", string(result))))
		testHierarchy := strings.Join(summary.ComponentTexts[1:len(summary.ComponentTexts)-1], " ")
		if testHierarchy != "" {
			testHierarchy += " "
		}
		s.print(0, s.colorize(defaultStyle, testHierarchy))
		testName := summary.ComponentTexts[len(summary.ComponentTexts)-1]
		s.print(0, s.colorize(boldStyle, testName))
		s.print(0, fmt.Sprintf(" <%v>", helputils.RoundSeconds(summary.RunTime)))

		results[result] = append(results[result], strings.Split(testName, " ")[0])

		fmt.Println(s.getPrintBuffer())
		s.printNewLine()
	}

	for _, result := range "vx~-" {
		if len(results[result]) > 0 {
			testNames := strings.Join(results[result], ", ")
			s.action(0, s.colorize(colorMap[result], "[%s] %s ")+testNames, string(result), titleMap[result])
		}
	}
}

func (s *consoleStenographer) SummarizeFailures(summaries []*types.SpecSummary) {
	failingSpecs := []*types.SpecSummary{}

	for _, summary := range summaries {
		if summary.HasFailureState() {
			failingSpecs = append(failingSpecs, summary)
		}
	}

	if len(failingSpecs) == 0 {
		return
	}

	plural := "s"
	if len(failingSpecs) == 1 {
		plural = ""
	}

	logging.Log().Error(s.colorize(redColor+boldStyle, "Details of %d Failure%s:", len(failingSpecs), plural))
	for _, summary := range failingSpecs {
		if summary.HasFailureState() {
			if summary.TimedOut() {
				s.print(0, s.colorize(redColor+boldStyle, "[Timeout...] "))
			} else if summary.Panicked() {
				s.print(0, s.colorize(redColor+boldStyle, "[Panic!] "))
			} else if summary.Failed() {
				s.print(0, s.colorize(redColor+boldStyle, "[Fail] "))
			}
			indentation := s.printSpecContext(summary.ComponentTexts, summary.ComponentCodeLocations, summary.Failure.ComponentType, summary.Failure.ComponentIndex, summary.State, true)
			s.flushLine()
			s.println(indentation, s.colorize(grayColor, summary.Failure.Location.String()))
		}
	}
}

func (s *consoleStenographer) startBlock() {
	if s.cursorState == cursorStateStreaming {
		s.flushLine()
		s.printDelimiter()
	} else if s.cursorState == cursorStateMidBlock {
		s.flushLine()
	}
}

func (s *consoleStenographer) midBlock() {
	s.cursorState = cursorStateMidBlock
}

func (s *consoleStenographer) endBlock() {
	s.printDelimiter()
	s.cursorState = cursorStateEndBlock
}

func (s *consoleStenographer) stream() {
	s.cursorState = cursorStateStreaming
}

func (s *consoleStenographer) printBlockWithMessage(header string, message string, spec *types.SpecSummary, succinct bool) {
	s.startBlock()
	s.println(0, header)

	indentation := s.printCodeLocationBlock(spec.ComponentTexts, spec.ComponentCodeLocations, types.SpecComponentTypeInvalid, 0, spec.State, succinct)

	if message != "" {
		s.flushLine()
		s.println(indentation, message)
	}

	s.endBlock()
}

func (s *consoleStenographer) printSpecFailure(message string, spec *types.SpecSummary, succinct bool, fullTrace bool) {
	s.startBlock()
	logging.Log().Error(s.colorize(redColor+boldStyle, "%s%s [%v]", message, s.failureContext(spec.Failure.ComponentType), spec.RunTime))

	indentation := s.printCodeLocationBlock(spec.ComponentTexts, spec.ComponentCodeLocations, spec.Failure.ComponentType, spec.Failure.ComponentIndex, spec.State, succinct)

	s.printFailure(indentation, spec.State, spec.Failure, fullTrace)
	s.endBlock()
}

func (s *consoleStenographer) failureContext(failedComponentType types.SpecComponentType) string {
	switch failedComponentType {
	case types.SpecComponentTypeBeforeSuite:
		return " in Suite Setup (BeforeSuite)"
	case types.SpecComponentTypeAfterSuite:
		return " in Suite Teardown (AfterSuite)"
	case types.SpecComponentTypeBeforeEach:
		return " in Spec Setup (BeforeEach)"
	case types.SpecComponentTypeJustBeforeEach:
		return " in Spec Setup (JustBeforeEach)"
	case types.SpecComponentTypeAfterEach:
		return " in Spec Teardown (AfterEach)"
	}

	return ""
}

func (s *consoleStenographer) printSkip(indentation int, spec *types.SpecSummary) {
	s.warn(0, s.colorize(boldStyle+cyanColor,
		"S [SKIPPING]%s: %s [%s]", s.failureContext(spec.Failure.ComponentType), spec.Failure.Message, helputils.RoundSeconds(spec.RunTime),
	))
	s.println(indentation, spec.Failure.Location.String())
}

func (s *consoleStenographer) printFailure(indentation int, state types.SpecState, failure types.SpecFailure, fullTrace bool) {
	delimiter := strings.Repeat("\U0001f621 ", 40)
	indentStr := strings.Repeat("  ", indentation)

	if state == types.SpecStatePanicked {
		logging.Log().Error(indentStr + s.colorize(redColor+boldStyle, failure.Message))
		logging.Log().Error(indentStr + s.colorize(redColor+boldStyle, delimiter))
		logging.Log().Error(indentStr + s.colorize(redColor, failure.ForwardedPanic))
		logging.Log().Error(indentStr + s.colorize(redColor+boldStyle, delimiter))
		logging.Log().Error(indentStr + failure.Location.String())
		logging.Log().Error(indentStr + s.colorize(redColor, "Full Stack Trace"))
		for _, stkLine := range strings.Split(failure.Location.FullStackTrace, "\n") {
			logging.Log().Error(indentStr + stkLine)
		}
	} else {
		logging.Log().Error(indentStr + s.colorize(redColor+boldStyle, delimiter))
		logging.Log().Error(indentStr + s.colorize(redColor, failure.Message))
		logging.Log().Error(indentStr + s.colorize(redColor+boldStyle, delimiter))
		logging.Log().Error(indentStr + failure.Location.String())
		if fullTrace {
			logging.Log().Error(indentStr + s.colorize(redColor, "Full Stack Trace"))
			for _, stkLine := range strings.Split(failure.Location.FullStackTrace, "\n") {
				logging.Log().Error(indentStr + stkLine)
			}
		}
	}

	helputils.ReportFullStack(func(msg string) { logging.Log().Debug(msg) })
}

func (s *consoleStenographer) printSpecContext(componentTexts []string, componentCodeLocations []types.CodeLocation, failedComponentType types.SpecComponentType, failedComponentIndex int, state types.SpecState, succinct bool) int {
	startIndex := 1
	indentation := 0

	if len(componentTexts) == 1 {
		startIndex = 0
	}

	for i := startIndex; i < len(componentTexts); i++ {
		if (state.IsFailure() || state == types.SpecStateSkipped) && i == failedComponentIndex {
			if state == types.SpecStateSkipped {
				s.warn(indentation, s.colorize(cyanColor+boldStyle, "%s [%s]", componentTexts[i], failedComponentType))
			} else if succinct {
				s.error(indentation, s.colorize(redColor+boldStyle, "[%s] %s ", failedComponentType, componentTexts[i]))
			} else {
				s.error(indentation, s.colorize(redColor+boldStyle, "%s [%s]", componentTexts[i], failedComponentType))
				s.println(indentation, s.colorize(grayColor, "%s", componentCodeLocations[i]))
			}
		} else {
			if succinct {
				s.print(0, s.colorize(alternatingColors[i%2], "%s ", componentTexts[i]))
			} else {
				s.println(indentation, componentTexts[i])
				s.println(indentation, s.colorize(grayColor, "%s", componentCodeLocations[i]))
			}
		}
		indentation++
	}

	return indentation
}

func (s *consoleStenographer) printCodeLocationBlock(componentTexts []string, componentCodeLocations []types.CodeLocation, failedComponentType types.SpecComponentType, failedComponentIndex int, state types.SpecState, succinct bool) int {
	indentation := s.printSpecContext(componentTexts, componentCodeLocations, failedComponentType, failedComponentIndex, state, succinct)

	if succinct {
		if len(componentTexts) > 0 {
			s.flushLine()
			s.print(0, s.colorize(grayColor, "%s", componentCodeLocations[len(componentCodeLocations)-1]))
		}
		s.flushLine()
		indentation = 1
	} else {
		indentation--
	}

	return indentation
}
