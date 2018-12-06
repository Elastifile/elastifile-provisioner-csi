// Package config manages Tesla's configuration. Sources for configuration
// can be a TOML config file, command line options and environment variables.
// Configuration is shared via the environment with all Tesla containers.
package config

import (
	"encoding/json"
	"fmt"

	"github.com/koding/multiconfig"
	log "gopkg.in/inconshreveable/log15.v2"

	"helputils"
	"issues"
	"logging"
	"runtimeutil"
	"types"
)

const JiraBaseUrl = "http://jira.il.elastifile.com"

var logger = logging.NewLogger("config")

// TODO: Add extraction of "doc" tags to show description for options in help.
// See vendor/github.com/koding/multiconfig/tag.go for the idea.

// Config holds the configuration. It is filled in with externally supplied
// configuration data, and can be accessed by all Tesla components.

func IsKnownIssue(conf *types.Config, issue string, summary string, err error) bool {
	known := false
	for _, knownIssue := range conf.Tesla.Issues.Known {
		if knownIssue == issue {
			known = true
			break
		}
	}

	if !known {
		logger.Info("Ignoring unknown issue", "issue", issue)
		return false
	}

	// TODO: Look up whether this is actually a known open issue in JIRA

	logger.Warn("Known issue", "issue", issue)
	issues.ReportIssue(issue, summary, err)
	return true
}

func ShouldWorkaround(conf *types.Config, issue string, logger log.Logger, summary string, ctx ...interface{}) bool {
	if conf != nil && helputils.ContainsStr(conf.Tesla.Issues.SkipWorkarounds, issue) {
		logger.Info("Ignore workaround: "+summary, "issue", issue)
		return false
	}
	ctx = append(helputils.StringSliceToInterfaces("issue", JiraBaseUrl+"/browse/"+issue), ctx...)
	caller := runtimeutil.CallerString(1)
	if caller != "" {
		ctx = append(ctx, "caller", caller)
	}
	logger.Warn("Workaround: "+summary, ctx...)
	return true
}

// Environment returns the config as a set of environment variables
// that can be passed e.g. to a Docker container.
func Environment(conf types.Config) []string {
	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    teslaPrefix,
		CamelCase: camelCase,
	}
	return envLoader.GetEnvironment(conf)
}

func LoadStartupEnvVars(conf types.StartupEnvVars) []string {
	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    types.EnvPrefix,
		CamelCase: camelCase,
	}
	return envLoader.GetEnvironment(conf)
}

// FromAllSources loads the config from its various sources: Config file, command line options, environment variables.
// This is used by the Tesla client which accepts configuration from all sources.
func FromAllSources(conf *types.Config, configFiles []string, options []string) error {
	loader, err := newLoader(configFiles, options)
	if err != nil {
		return err
	}
	return loader.Load(conf)
}

// ConfigFromEnvironment loads the config just from the appropriate environment variables.
// This is used by various Tesla agents that get the environment via Docker.
func FromEnvironment() (types.Config, error) {
	conf := types.NewConfig()

	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    teslaPrefix,
		CamelCase: camelCase,
	}
	err := envLoader.Load(conf)

	conf.Logging.PrintDebug("Config", "conf", fmt.Sprintf("%+v", conf))
	return *conf, err
}

func FromTagsAndEnvironment() (types.Config, error) {
	conf := types.NewConfig()

	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    teslaPrefix,
		CamelCase: camelCase,
	}
	if err := envLoader.Load(conf); err != nil {
		conf.Logging.PrintError("Failed env loader", "err", err.Error())
		return *conf, err
	}

	conf.Logging.PrintDebug("Config", "conf", fmt.Sprintf("%+v", conf))
	return *conf, nil
}

var DeployConfig string

const (
	SanityImage       = "agent/sanity"
	SlaveImage        = "agent/slave"
	ToolyImage        = "agent/tooly"
	ECSImage          = "system/ecs"
	ELFSImage         = "system/elfs"
	EmanageImage      = "system/emanage"
	EmanageDBImage    = "system/emanagedb"
	NewAPIImage       = "test/newapi"
	GNATSDImage       = "thirdparty/gnatsd"
	MinIOImage        = "thirdparty/minio"
	CthonImage        = "tool/cthon"
	ErunImage         = "tool/erun"
	FillCapacityImage = "tool/fstool"
	FIOImage          = "tool/fio"
	SFS2008Image      = "tool/sfs2008"
	SFS2014Image      = "tool/sfs2014"
	VDBenchImage      = "tool/vdbench"
)

func GetImage(name string) types.Image {
	parsed := map[string]types.Image{}
	if err := json.Unmarshal([]byte(DeployConfig), &parsed); err != nil {
		panic("Failed unmarshal of Deploy Config:\n" + DeployConfig + "\nerr: " + err.Error())
	}
	return parsed[name]
}

func GetToolImage(tool types.ToolName) types.Image {
	logger.Info("Asking for tool image", "tool", tool)
	return GetImage("tool/" + string(tool))
}
