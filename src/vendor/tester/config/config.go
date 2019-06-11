/*
Ginkgo accepts a number of configuration options.

These are documented [here](http://onsi.github.io/ginkgo/#the_ginkgo_cli)

You can also learn more via

	ginkgo help

or (I kid you not):

	go test -asdf
*/
package config

import (
	"flag"
	"fmt"
	"time"
)

const VERSION = "2.1.0"
const DefaultTimeout = 0 // in seconds

type TesterConfigType struct {
	RandomSeed        int64
	RandomizeAllSpecs bool
	FocusString       string
	IncludeFile       string
	SkipString        string
	SpecTimeout       string
	SuiteTimeout      string
	FailOnPending     bool
	FailFast          bool
	EmitSpecProgress  bool
	DryRun            bool
}

var TesterConfig = TesterConfigType{}

type DefaultReporterConfigType struct {
	NoColor           bool
	SlowSpecThreshold float64
	NoisyPendings     bool
	Succinct          bool
	Verbose           bool
	FullTrace         bool
}

var DefaultReporterConfig = DefaultReporterConfigType{}

func processPrefix(prefix string) string {
	if prefix != "" {
		prefix = prefix + "."
	}
	return prefix
}

func Flags(flagSet *flag.FlagSet, prefix string) {
	prefix = processPrefix(prefix)
	flagSet.Int64Var(&(TesterConfig.RandomSeed), prefix+"seed", time.Now().Unix(), "The seed used to randomize the spec suite.")
	flagSet.BoolVar(&(TesterConfig.RandomizeAllSpecs), prefix+"randomizeAllSpecs", false, "If set, ginkgo will randomize all specs together.  By default, ginkgo only randomizes the top level Describe/Context groups.")
	flagSet.BoolVar(&(TesterConfig.FailOnPending), prefix+"failOnPending", false, "If set, ginkgo will mark the test suite as failed if any specs are pending.")
	flagSet.BoolVar(&(TesterConfig.FailFast), prefix+"failFast", false, "If set, ginkgo will stop running a test suite after a failure occurs.")
	flagSet.BoolVar(&(TesterConfig.DryRun), prefix+"dryRun", false, "If set, ginkgo will walk the test hierarchy without actually running anything.  Best paired with -v.")
	flagSet.StringVar(&(TesterConfig.FocusString), prefix+"focus", "", "If set, ginkgo will only run specs that match this regular expression.")
	flagSet.StringVar(&(TesterConfig.IncludeFile), prefix+"includeFile", "", "If set, ginkgo will only run specs that are included in the specified YAML file.")
	flagSet.StringVar(&(TesterConfig.SkipString), prefix+"skip", "", "If set, ginkgo will only run specs that do not match this regular expression.")
	flagSet.BoolVar(&(TesterConfig.EmitSpecProgress), prefix+"progress", false, "If set, ginkgo will emit progress information as each spec runs to the TesterWriter.")
	flagSet.StringVar(&(TesterConfig.SuiteTimeout), prefix+"suiteTimeout", "0s", "If elpased more than this since start of suite then abort.")
	flagSet.StringVar(&(TesterConfig.SpecTimeout), prefix+"specTimeout", "0s", "If running more than this then Timeout the spec and continue the suite.")

	flagSet.BoolVar(&(DefaultReporterConfig.NoColor), prefix+"noColor", false, "If set, suppress color output in default reporter.")
	flagSet.Float64Var(&(DefaultReporterConfig.SlowSpecThreshold), prefix+"slowSpecThreshold", 5.0, "(in seconds) Specs that take longer to run than this threshold are flagged as slow by the default reporter.")
	flagSet.BoolVar(&(DefaultReporterConfig.NoisyPendings), prefix+"noisyPendings", true, "If set, default reporter will shout about pending tests.")
	flagSet.BoolVar(&(DefaultReporterConfig.Verbose), prefix+"v", false, "If set, default reporter print out all specs as they begin.")
	flagSet.BoolVar(&(DefaultReporterConfig.Succinct), prefix+"succinct", false, "If set, default reporter prints out a very succinct report")
	flagSet.BoolVar(&(DefaultReporterConfig.FullTrace), prefix+"trace", false, "If set, default reporter prints out the full stack trace when a failure occurs")
}

func BuildFlagArgs(prefix string, ginkgo TesterConfigType, reporter DefaultReporterConfigType) []string {
	prefix = processPrefix(prefix)
	result := make([]string, 0)

	if ginkgo.RandomSeed > 0 {
		result = append(result, fmt.Sprintf("--%sseed=%d", prefix, ginkgo.RandomSeed))
	}

	if ginkgo.RandomizeAllSpecs {
		result = append(result, fmt.Sprintf("--%srandomizeAllSpecs", prefix))
	}

	if ginkgo.FailOnPending {
		result = append(result, fmt.Sprintf("--%sfailOnPending", prefix))
	}

	if ginkgo.FailFast {
		result = append(result, fmt.Sprintf("--%sfailFast", prefix))
	}

	if ginkgo.DryRun {
		result = append(result, fmt.Sprintf("--%sdryRun", prefix))
	}

	if ginkgo.SpecTimeout != "0s" {
		if _, err := time.ParseDuration(ginkgo.SpecTimeout); err != nil {
			panic("Malformatted specTimeout: " + err.Error())
		}
		result = append(result, fmt.Sprintf("--%sspecTimeout=%s", prefix, ginkgo.SpecTimeout))
	}

	if ginkgo.SuiteTimeout != "0s" {
		if _, err := time.ParseDuration(ginkgo.SuiteTimeout); err != nil {
			panic("Malformatted suiteTimeout: " + err.Error())
		}
		result = append(result, fmt.Sprintf("--%ssuiteTimeout=%s", prefix, ginkgo.SuiteTimeout))
	}

	if ginkgo.FocusString != "" {
		result = append(result, fmt.Sprintf("--%sfocus=%s", prefix, ginkgo.FocusString))
	}

	if ginkgo.IncludeFile != "" {
		result = append(result, fmt.Sprintf("--%include=%s", prefix, ginkgo.IncludeFile))
	}

	if ginkgo.SkipString != "" {
		result = append(result, fmt.Sprintf("--%sskip=%s", prefix, ginkgo.SkipString))
	}

	if ginkgo.EmitSpecProgress {
		result = append(result, fmt.Sprintf("--%sprogress", prefix))
	}

	if reporter.NoColor {
		result = append(result, fmt.Sprintf("--%snoColor", prefix))
	}

	if reporter.SlowSpecThreshold > 0 {
		result = append(result, fmt.Sprintf("--%sslowSpecThreshold=%.5f", prefix, reporter.SlowSpecThreshold))
	}

	if !reporter.NoisyPendings {
		result = append(result, fmt.Sprintf("--%snoisyPendings=false", prefix))
	}

	if reporter.Verbose {
		result = append(result, fmt.Sprintf("--%sv", prefix))
	}

	if reporter.Succinct {
		result = append(result, fmt.Sprintf("--%ssuccinct", prefix))
	}

	if reporter.FullTrace {
		result = append(result, fmt.Sprintf("--%strace", prefix))
	}

	return result
}
