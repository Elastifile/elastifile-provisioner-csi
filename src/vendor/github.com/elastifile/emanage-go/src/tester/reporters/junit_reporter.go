/*

JUnit XML Reporter for Ginkgo

For usage instructions: http://onsi.github.io/ginkgo/#generating_junit_xml_output

*/

package reporters

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"tester/config"
	"tester/internal/logging"
	"tester/types"
)

type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	TestCases []JUnitTestCase `xml:"testcase"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      float64         `xml:"time,attr"`

	Properties []JUnitProperty `xml:"properties>property,omitempty"`
}

type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JUnitTestCase struct {
	Name           string               `xml:"name,attr"`
	ClassName      string               `xml:"classname,attr"`
	TestCustomData string               `xml:"testCustomData,attr"`
	FailureMessage *JUnitFailureMessage `xml:"failure,omitempty"`
	Skipped        *JUnitSkipped        `xml:"skipped,omitempty"`
	Time           float64              `xml:"time,attr"`
}

type JUnitFailureMessage struct {
	Type    string `xml:"type,attr"`
	Message string `xml:",chardata"`
}

type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`
}

type JUnitReporter struct {
	suite         JUnitTestSuite
	filename      string
	testSuiteName string
}

//NewJUnitReporter creates a new JUnit XML reporter.  The XML will be stored in the passed in filename.
func NewJUnitReporter(filename string) *JUnitReporter {
	return &JUnitReporter{
		filename: filename,
	}
}

func (reporter *JUnitReporter) AddProperty(p JUnitProperty) {
	reporter.suite.Properties = append(reporter.suite.Properties, p)
}

func (reporter *JUnitReporter) SpecSuiteWillBegin(config config.TesterConfigType, summary *types.SuiteSummary) {
	reporter.suite = JUnitTestSuite{
		Tests:      summary.NumberOfSpecsThatWillBeRun,
		TestCases:  []JUnitTestCase{},
		Properties: reporter.suite.Properties,
	}
	reporter.testSuiteName = summary.SuiteDescription
}

func (reporter *JUnitReporter) SpecWillRun(specSummary *types.SpecSummary) {
}

func (reporter *JUnitReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {
	reporter.handleSetupSummary("BeforeSuite", setupSummary)
}

func (reporter *JUnitReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {
	reporter.handleSetupSummary("AfterSuite", setupSummary)
}

var customDataMap = make(map[string]string)

func (reporter *JUnitReporter) SetTestCustomData(name string, data string) {
	customDataMap[name] = data
}
func (reporter *JUnitReporter) getTestCustomData(name string) string {
	data, ok := customDataMap[name]
	if !ok {
		return ""
	}
	return data
}

func (reporter *JUnitReporter) handleSetupSummary(name string, setupSummary *types.SetupSummary) {
	if setupSummary.State != types.SpecStatePassed {
		testCase := JUnitTestCase{
			Name:      name,
			ClassName: reporter.testSuiteName,
		}

		testCase.FailureMessage = &JUnitFailureMessage{
			Type:    reporter.failureTypeForState(setupSummary.State),
			Message: fmt.Sprintf("%s\n%s", setupSummary.Failure.ComponentCodeLocation.String(), setupSummary.Failure.Message),
		}
		testCase.Time = setupSummary.RunTime.Seconds()
		reporter.suite.TestCases = append(reporter.suite.TestCases, testCase)
	}
}

func (reporter *JUnitReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	testName := strings.Join(specSummary.ComponentTexts[1:], " ")
	testCase := JUnitTestCase{
		Name:      testName,
		ClassName: reporter.testSuiteName,
	}
	if specSummary.State == types.SpecStateFailed || specSummary.State == types.SpecStateTimedOut || specSummary.State == types.SpecStatePanicked {
		testCase.FailureMessage = &JUnitFailureMessage{
			Type:    reporter.failureTypeForState(specSummary.State),
			Message: fmt.Sprintf("%s\n%s", specSummary.Failure.ComponentCodeLocation.String(), specSummary.Failure.Message),
		}
	}
	if specSummary.State == types.SpecStateSkipped || specSummary.State == types.SpecStatePending {
		testCase.Skipped = &JUnitSkipped{}
	}
	testCase.Time = specSummary.RunTime.Seconds()
	testCase.TestCustomData = reporter.getTestCustomData(testName)

	reporter.suite.TestCases = append(reporter.suite.TestCases, testCase)
}

func (reporter *JUnitReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	reporter.suite.Time = summary.RunTime.Seconds()
	reporter.suite.Failures = summary.NumberOfFailedSpecs
	file, err := os.Create(reporter.filename)
	if err != nil {
		logging.Log().Error("Failed to create JUnit report file", "name", reporter.filename, "err", err)
	}
	defer file.Close()
	file.WriteString(xml.Header)
	encoder := xml.NewEncoder(file)
	encoder.Indent("  ", "    ")
	err = encoder.Encode(reporter.suite)
	if err != nil {
		logging.Log().Error("Failed to generate JUnit report", "err", err)
	}
}

func (reporter *JUnitReporter) failureTypeForState(state types.SpecState) string {
	switch state {
	case types.SpecStateFailed:
		return "Failure"
	case types.SpecStateTimedOut:
		return "Timeout"
	case types.SpecStatePanicked:
		return "Panic"
	default:
		return ""
	}
}
