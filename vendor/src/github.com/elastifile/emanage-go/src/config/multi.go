package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/koding/multiconfig"

	"helputils"
)

const (
	teslaPrefix = "TESLA"
	camelCase   = true
)

func newLoader(configFiles []string, options []string) (*multiconfig.DefaultLoader, error) {
	loaders := []multiconfig.Loader{}

	// 1. from struct tags
	loaders = append(loaders, &multiconfig.TagLoader{})

	// 2. from files (toml for tesla, json for elab related stuff)
	for _, configFile := range configFiles {
		ext := filepath.Ext(configFile)
		switch strings.ToLower(ext) {
		case ".json":
			ensuredConfigFile, err := ensureSystemElabParents(configFile)
			if err != nil {
				return nil, err
			}
			loaders = append(loaders, &multiconfig.JSONLoader{Path: ensuredConfigFile})
		case ".toml":
			loaders = append(loaders, &tomlLoader{Path: configFile, Required: false})
		default:
			return nil, errors.Errorf("unsupported config file extension: %s, expecting: '.json' or '.toml'", ext)
		}
	}

	// 3. from environment variables
	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    teslaPrefix,
		CamelCase: camelCase,
	}
	loaders = append(loaders, envLoader)

	// 4. from cli flags
	optLoader := &OptionLoader{
		Options: options,
	}
	loaders = append(loaders, optLoader)

	result := &multiconfig.DefaultLoader{
		Loader: multiconfig.MultiLoader(loaders...),
		Validator: multiconfig.MultiValidator(
			&multiconfig.RequiredValidator{},
		),
	}
	return result, nil
}

func ensureSystemElabParents(configFile string) (string, error) {
	body, err := helputils.ReadAll(configFile)
	if err != nil {
		return configFile, errors.Wrap(err, 0)
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return configFile, errors.Wrap(err, 0)
	}

	if data, hasData := m["data"]; hasData {
		if _, hasEmanage := data.(map[string]interface{})["emanage"]; hasEmanage {
			newBody := []byte(fmt.Sprintf(`
{"system":
{"elab": 
%s
}
}
`, body))
			tmpConfigJson, err := helputils.TmpRandomFileName("ensureSystemElabParents")
			if err != nil {
				return configFile, errors.New(err)
			}
			tmpConfigJson += ".json"

			err = helputils.WriteFile(tmpConfigJson, newBody, os.ModePerm)
			return tmpConfigJson, err
		}
	}

	return configFile, nil
}
