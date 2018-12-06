package suite

import (
	"math/rand"
	"strings"
	"time"

	"tester/config"
	"tester/internal/containernode"
	"tester/internal/failer"
	"tester/internal/leafnodes"
	"tester/internal/logging"
	"tester/internal/spec"
	"tester/internal/specrunner"
	"tester/internal/writer"
	"tester/reporters"
	"tester/types"
)

type ginkgoTestingT interface {
	Fail()
}

type Suite struct {
	topLevelContainer *containernode.ContainerNode
	currentContainer  *containernode.ContainerNode
	containerIndex    int
	beforeSuiteNode   leafnodes.SuiteNode
	afterSuiteNode    leafnodes.SuiteNode
	runner            *specrunner.SpecRunner
	failer            *failer.Failer
	running           bool
}

func New(failer *failer.Failer) *Suite {
	topLevelContainer := containernode.New("[Top Level]", types.FlagTypeNone, types.CodeLocation{})

	return &Suite{
		topLevelContainer: topLevelContainer,
		currentContainer:  topLevelContainer,
		failer:            failer,
		containerIndex:    1,
	}
}

func (suite *Suite) Run(t ginkgoTestingT, description string, reporters []reporters.Reporter, writer writer.WriterInterface, config config.TesterConfigType) (bool, bool) {
	r := rand.New(rand.NewSource(config.RandomSeed))
	suite.topLevelContainer.Shuffle(r)

	suite.runner = specrunner.New(description, suite.beforeSuiteNode, suite.afterSuiteNode, reporters, writer, config)

	hasProgFocus := false

	if config.FocusString != "" {
		hasProgFocus = suite.addSpecsByFocus(config.FocusString, description, config)
	}

	if config.IncludeFile != "" {
		for _, incFile := range strings.Split(config.IncludeFile, ";") {
			err := suite.walkIncludeFile(incFile,
				func(focus string) {
					hasProgFocus = hasProgFocus || suite.addSpecsByFocus(focus, description, config)
					logging.Log().Debug("Added specs", "focus", focus)
				})
			if err != nil {
				logging.Log().Crit(err.Error())
				panic(err.Error())
			}
			logging.Log().Info("Added all specs", "includFile", incFile)
		}
	}

	suite.running = true
	success := suite.runner.Run()
	if !success {
		t.Fail()
	}

	return success, hasProgFocus
}

func (suite *Suite) walkIncludeFile(includeFile string, do func(string)) (err error) {
	//yamlBody, err := ioutil.ReadFile(includeFile)
	//if err != nil {
	//	return err
	//}
	//includeDir := filepath.Dir(includeFile)
	//
	//err = elvish.Walk(elvish.WalkOpts{
	//	YamlBody: yamlBody,
	//	Do:       do,
	//	Recursion: func(i int, item string) error {
	//		path := helputils.FullPath(includeDir, item)
	//		if err := suite.walkIncludeFile(path, do); err != nil {
	//			return fmt.Errorf("%s (item #%d: '%s')", err, i+1, item)
	//		}
	//		return nil
	//	},
	//})

	return err
}

func (suite *Suite) addSpecsByFocus(focus string, description string, config config.TesterConfigType) (hasProgFocus bool) {
	for _, splitFocus := range strings.Split(focus, ";") {
		splitFocus = strings.TrimSpace(splitFocus)
		if splitFocus == "" {
			continue
		}
		specs := suite.generateSpecs(description, config, splitFocus)

		var matchingSpecs []*spec.Spec
		for _, s := range specs.Specs() {
			if !s.Skipped() {
				matchingSpecs = append(matchingSpecs, s)
			}
		}
		if len(matchingSpecs) == 0 {
			if config.SkipString == "" {
				panicStr := "No matching specs for focus: " + splitFocus
				logging.Log().Crit(panicStr)
				panic(panicStr)
			} else {
				logging.Log().Warn("Focus " + splitFocus + " not matched or skipped")
			}
		}

		addedSpecs := spec.NewSpecs(matchingSpecs)
		suite.runner.AddSpecs(addedSpecs)
		logging.Log().Debug("Added specs", "focus", splitFocus)

		hasProgFocus = hasProgFocus || addedSpecs.HasProgrammaticFocus()
	}

	return hasProgFocus
}

func (suite *Suite) Abort() {
	suite.runner.Abort()
}

func (suite *Suite) generateSpecs(description string, config config.TesterConfigType, focusString string) *spec.Specs {
	specsSlice := []*spec.Spec{}
	suite.topLevelContainer.BackPropagateProgrammaticFocus()
	for _, collatedNodes := range suite.topLevelContainer.Collate() {
		specsSlice = append(specsSlice, spec.New(collatedNodes.Subject, collatedNodes.Containers, config.EmitSpecProgress))
	}

	specs := spec.NewSpecs(specsSlice)

	if config.RandomizeAllSpecs {
		specs.Shuffle(rand.New(rand.NewSource(config.RandomSeed)))
	}

	specs.ApplyFocus(description, focusString, config.SkipString)

	return specs
}

func (suite *Suite) CurrentRunningSpecSummary() (*types.SpecSummary, bool) {
	return suite.runner.CurrentSpecSummary()
}

func (suite *Suite) SetBeforeSuiteNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.beforeSuiteNode != nil {
		panic("You may only call BeforeSuite once!")
	}
	suite.beforeSuiteNode = leafnodes.NewBeforeSuiteNode(body, codeLocation, timeout, suite.failer)
}

func (suite *Suite) SetAfterSuiteNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.afterSuiteNode != nil {
		panic("You may only call AfterSuite once!")
	}
	suite.afterSuiteNode = leafnodes.NewAfterSuiteNode(body, codeLocation, timeout, suite.failer)
}

func (suite *Suite) PushContainerNode(text string, body func(), flag types.FlagType, codeLocation types.CodeLocation) {
	container := containernode.New(text, flag, codeLocation)
	suite.currentContainer.PushContainerNode(container)

	previousContainer := suite.currentContainer
	suite.currentContainer = container
	suite.containerIndex++

	body()

	suite.containerIndex--
	suite.currentContainer = previousContainer
}

func (suite *Suite) PushItNode(text string, body interface{}, flag types.FlagType, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.running {
		suite.failer.Fail("You may only call It from within a Describe or Context", codeLocation)
	}
	suite.currentContainer.PushSubjectNode(leafnodes.NewItNode(text, body, flag, codeLocation, timeout, suite.failer, suite.containerIndex))
}

func (suite *Suite) PushBeforeEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.running {
		suite.failer.Fail("You may only call BeforeEach from within a Describe or Context", codeLocation)
	}
	suite.currentContainer.PushSetupNode(leafnodes.NewBeforeEachNode(body, codeLocation, timeout, suite.failer, suite.containerIndex))
}

func (suite *Suite) PushJustBeforeEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.running {
		suite.failer.Fail("You may only call JustBeforeEach from within a Describe or Context", codeLocation)
	}
	suite.currentContainer.PushSetupNode(leafnodes.NewJustBeforeEachNode(body, codeLocation, timeout, suite.failer, suite.containerIndex))
}

func (suite *Suite) PushAfterEachNode(body interface{}, codeLocation types.CodeLocation, timeout time.Duration) {
	if suite.running {
		suite.failer.Fail("You may only call AfterEach from within a Describe or Context", codeLocation)
	}
	suite.currentContainer.PushSetupNode(leafnodes.NewAfterEachNode(body, codeLocation, timeout, suite.failer, suite.containerIndex))
}
