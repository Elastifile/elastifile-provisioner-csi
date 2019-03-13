/*
Ginkgo is a BDD-style testing framework for Golang

The godoc documentation describes Ginkgo's API.  More comprehensive documentation (with examples!) is available at http://onsi.github.io/ginkgo/

Ginkgo's preferred matcher library is [Gomega](http://github.com/onsi/gomega)

Ginkgo on Github: http://github.com/onsi/ginkgo

Ginkgo is MIT-Licensed
*/
package tester

import (
	"flag"
	"io"
	"os"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"tester/config"
	"tester/internal/codelocation"
	"tester/internal/failer"
	"tester/internal/logging"
	"tester/internal/suite"
	"tester/internal/writer"
	"tester/reporters"
	"tester/reporters/stenographer"
	"tester/types"
)

const GINKGO_VERSION = config.VERSION
const GINKGO_PANIC = `
Your test failed.
Ginkgo panics to prevent subsequent assertions from running.
Normally Ginkgo rescues this panic so you shouldn't see it.

But, if you make an assertion in a goroutine, Ginkgo can't capture the panic.
To circumvent this, you should call

	defer TesterRecover()

at the top of the goroutine that caused this panic.
`

var globalSuite *suite.Suite
var globalFailer *failer.Failer

func init() {
	config.Flags(flag.CommandLine, "")
	config.Flags(flag.CommandLine, "ginkgo")
	flag.Parse()

	TesterWriter = writer.New(os.Stdout)
	globalFailer = failer.New()
	globalSuite = suite.New(globalFailer)
}

func GlobalSuite() *suite.Suite {
	return globalSuite
}

func GlobalFailer() *failer.Failer {
	return globalFailer
}

func TesterSetLogger(logger log.Logger, locationRoot string) {
	logging.SetLogger(logger)
	TesterWriter = writer.New(logger)
	types.SetSubpath(locationRoot)
}

//TesterWriter implements an io.Writer
//When running in verbose mode any writes to TesterWriter will be immediately printed
//to stdout.  Otherwise, TesterWriter will buffer any writes produced during the current test and flush them to screen
//only if the current test fails.
var TesterWriter io.Writer

//The interface by which Ginkgo receives *testing.T
type GinkgoTestingT interface {
	Fail()
}

//Custom Ginkgo test reporters must implement the Reporter interface.
//
//The custom reporter is passed in a SuiteSummary when the suite begins and ends,
//and a SpecSummary just before a spec begins and just after a spec ends
type Reporter reporters.Reporter

//Asynchronous specs are given a channel of the Done type.  You must close or write to the channel
//to tell Ginkgo that your async test is done.
type Done chan<- interface{}

//TestDescription represents the information about the current running test returned by CurrentTestDescription
//	FullTestText: a concatenation of ComponentTexts and the TestText
//	ComponentTexts: a list of all texts for the Describes & Contexts leading up to the current test
//	TestText: the text in the actual It or Measure node
//	FileName: the name of the file containing the current test
//	LineNumber: the line number for the current test
//	Failed: if the current test has failed, this will be true (useful in an AfterEach)
type TestDescription struct {
	FullTestText   string
	ComponentTexts []string
	TestText       string

	FileName   string
	LineNumber int

	Failed bool
}

//CurrentGinkgoTestDescripton returns information about the current running test.
func CurrentTestDescription() TestDescription {
	summary, ok := globalSuite.CurrentRunningSpecSummary()
	if !ok {
		return TestDescription{}
	}

	subjectCodeLocation := summary.ComponentCodeLocations[len(summary.ComponentCodeLocations)-1]

	return TestDescription{
		ComponentTexts: summary.ComponentTexts[1:],
		FullTestText:   strings.Join(summary.ComponentTexts[1:], " "),
		TestText:       summary.ComponentTexts[len(summary.ComponentTexts)-1],
		FileName:       subjectCodeLocation.FileName,
		LineNumber:     subjectCodeLocation.LineNumber,
		Failed:         summary.HasFailureState(),
	}
}

//RunSpecs is the entry point for the Ginkgo test runner.
//You must call this within a Golang testing TestX(t *testing.T) function.
//
//To bootstrap a test suite you can use the Ginkgo CLI:
//
//	ginkgo bootstrap
func RunSpecs(t GinkgoTestingT, description string) bool {
	specReporters := []Reporter{buildDefaultReporter()}
	return RunSpecsWithCustomReporters(t, description, specReporters)
}

//To run your tests with Ginkgo's default reporter and your custom reporter(s), replace
//RunSpecs() with this method.
func RunSpecsWithDefaultAndCustomReporters(t GinkgoTestingT, description string, specReporters ...Reporter) bool {
	specReporters = append([]Reporter{buildDefaultReporter()}, specReporters...)
	return RunSpecsWithCustomReporters(t, description, specReporters)
}

//To run your tests with your custom reporter(s) (and *not* Ginkgo's default reporter), replace
//RunSpecs() with this method.  Note that parallel tests will not work correctly without the default reporter
func RunSpecsWithCustomReporters(t GinkgoTestingT, description string, specReporters []Reporter) bool {
	writer := TesterWriter.(*writer.Writer)
	writer.SetStream(config.DefaultReporterConfig.Verbose)
	reporters := make([]reporters.Reporter, len(specReporters))
	for i, reporter := range specReporters {
		reporters[i] = reporter
	}
	passed, hasFocusedTests := globalSuite.Run(t, description, reporters, writer, config.TesterConfig)
	if passed && hasFocusedTests {
		logging.Log().Info("PASS | FOCUSED")
		os.Exit(types.TESTER_FOCUS_EXIT_CODE)
	}
	return passed
}

func buildDefaultReporter() Reporter {
	stenographer := stenographer.New(!config.DefaultReporterConfig.NoColor)
	return reporters.NewDefaultReporter(config.DefaultReporterConfig, stenographer)
}

//Skip notifies Ginkgo that the current spec should be skipped.
func Skip(message string, callerSkip ...int) {
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}

	globalFailer.Skip(message, codelocation.New(skip+1))
	panic(GINKGO_PANIC)
}

func AbortSuite() {
	globalSuite.Abort()
}

//Fail notifies Ginkgo that the current spec has failed. (Gomega will call Fail for you automatically when an assertion fails.)
func Fail(message string, callerSkip ...int) {
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}

	globalFailer.Fail(message, codelocation.New(skip+1))
	panic(GINKGO_PANIC)
}

//TesterRecover should be deferred at the top of any spawned goroutine that (may) call `Fail`
//Since Gomega assertions call fail, you should throw a `defer TesterRecover()` at the top of any goroutine that
//calls out to Gomega
//
//Here's why: Ginkgo's `Fail` method records the failure and then panics to prevent
//further assertions from running.  This panic must be recovered.  Ginkgo does this for you
//if the panic originates in a Ginkgo node (an It, BeforeEach, etc...)
//
//Unfortunately, if a panic originates on a goroutine *launched* from one of these nodes there's no
//way for Ginkgo to rescue the panic.  To do this, you must remember to `defer TesterRecover()` at the top of such a goroutine.
func TesterRecover() {
	e := recover()
	if e != nil {
		globalFailer.Panic(codelocation.New(1), e)
	}
}

//Describe blocks allow you to organize your specs.  A Describe block can contain any number of
//BeforeEach, AfterEach, JustBeforeEach, It, and Measurement blocks.
//
//In addition you can nest Describe and Context blocks.  Describe and Context blocks are functionally
//equivalent.  The difference is purely semantic -- you typical Describe the behavior of an object
//or method and, within that Describe, outline a number of Contexts.
func Describe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeNone, codelocation.New(1))
	return true
}

//You can focus the tests within a describe block using FDescribe
func FDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeFocused, codelocation.New(1))
	return true
}

//You can mark the tests within a describe block as pending using PDescribe
func PDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}

//You can mark the tests within a describe block as pending using XDescribe
func XDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}

//Context blocks allow you to organize your specs.  A Context block can contain any number of
//BeforeEach, AfterEach, JustBeforeEach, It, and Measurement blocks.
//
//In addition you can nest Describe and Context blocks.  Describe and Context blocks are functionally
//equivalent.  The difference is purely semantic -- you typical Describe the behavior of an object
//or method and, within that Describe, outline a number of Contexts.
func Context(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeNone, codelocation.New(1))
	return true
}

//You can focus the tests within a describe block using FContext
func FContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeFocused, codelocation.New(1))
	return true
}

//You can mark the tests within a describe block as pending using PContext
func PContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}

//You can mark the tests within a describe block as pending using XContext
func XContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}

//It blocks contain your test code and assertions.  You cannot nest any other Ginkgo blocks
//within an It block.
//
//Ginkgo will normally run It blocks synchronously.  To perform asynchronous tests, pass a
//function that accepts a Done channel.  When you do this, you can also provide an optional timeout.
func It(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeNone, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//You can focus individual Its using FIt
func FIt(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeFocused, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//You can mark Its as pending using PIt
func PIt(text string, _ ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}

//You can mark Its as pending using XIt
func XIt(text string, _ ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}

//Specify blocks are aliases for It blocks and allow for more natural wording in situations
//which "It" does not fit into a natural sentence flow. All the same protocols apply for Specify blocks
//which apply to It blocks.
func Specify(text string, body interface{}, timeout ...float64) bool {
	return It(text, body, timeout...)
}

//You can focus individual Specifys using FSpecify
func FSpecify(text string, body interface{}, timeout ...float64) bool {
	return FIt(text, body, timeout...)
}

//You can mark Specifys as pending using PSpecify
func PSpecify(text string, is ...interface{}) bool {
	return PIt(text, is...)
}

//You can mark Specifys as pending using XSpecify
func XSpecify(text string, is ...interface{}) bool {
	return XIt(text, is...)
}

//By allows you to better document large Its.
//
//Generally you should try to keep your Its short and to the point.  This is not always possible, however,
//especially in the context of integration tests that capture a particular workflow.
//
//By allows you to document such flows.  By must be called within a runnable node (It, BeforeEach, Measure, etc...)
//By will simply log the passed in text to the TesterWriter.  If By is handed a function it will immediately run the function.
func By(text string, callbacks ...func()) {
	msg := "Step: " + text
	if !config.DefaultReporterConfig.NoColor {
		msg = "\x1b[1m" + msg + "\x1b[0m"
	}
	logging.Log().Info(msg)
	if len(callbacks) == 1 {
		callbacks[0]()
	}
	if len(callbacks) > 1 {
		panic("just one callback per By, please")
	}
}

//BeforeSuite blocks are run just once before any specs are run.  When running in parallel, each
//parallel node process will call BeforeSuite.
//
//BeforeSuite blocks can be made asynchronous by providing a body function that accepts a Done channel
//
//You may only register *one* BeforeSuite handler per test suite.  You typically do so in your bootstrap file at the top level.
func BeforeSuite(body interface{}, timeout ...float64) bool {
	globalSuite.SetBeforeSuiteNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//AfterSuite blocks are *always* run after all the specs regardless of whether specs have passed or failed.
//Moreover, if Ginkgo receives an interrupt signal (^C) it will attempt to run the AfterSuite before exiting.
//
//When running in parallel, each parallel node process will call AfterSuite.
//
//AfterSuite blocks can be made asynchronous by providing a body function that accepts a Done channel
//
//You may only register *one* AfterSuite handler per test suite.  You typically do so in your bootstrap file at the top level.
func AfterSuite(body interface{}, timeout ...float64) bool {
	globalSuite.SetAfterSuiteNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//BeforeEach blocks are run before It blocks.  When multiple BeforeEach blocks are defined in nested
//Describe and Context blocks the outermost BeforeEach blocks are run first.
//
//Like It blocks, BeforeEach blocks can be made asynchronous by providing a body function that accepts
//a Done channel
func BeforeEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushBeforeEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//JustBeforeEach blocks are run before It blocks but *after* all BeforeEach blocks.  For more details,
//read the [documentation](http://onsi.github.io/ginkgo/#separating_creation_and_configuration_)
//
//Like It blocks, BeforeEach blocks can be made asynchronous by providing a body function that accepts
//a Done channel
func JustBeforeEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushJustBeforeEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

//AfterEach blocks are run after It blocks.   When multiple AfterEach blocks are defined in nested
//Describe and Context blocks the innermost AfterEach blocks are run first.
//
//Like It blocks, AfterEach blocks can be made asynchronous by providing a body function that accepts
//a Done channel
func AfterEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushAfterEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

func parseTimeout(timeout ...float64) time.Duration {
	if len(timeout) == 0 {
		dur, err := time.ParseDuration(config.TesterConfig.SpecTimeout)
		if err != nil {
			panic("Malformatter specTimeout: " + config.TesterConfig.SpecTimeout + ", " + err.Error())
		}
		return dur
	} else {
		return time.Duration(timeout[0] * float64(time.Second))
	}
}
