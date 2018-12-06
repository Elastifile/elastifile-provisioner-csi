package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-errors/errors"
	log "gopkg.in/inconshreveable/log15.v2"
)

// tomlLoader satisifies the multiconfig.Loader interface.
type tomlLoader struct {
	Path     string
	Required bool
}

var defaultTomlPath string = os.ExpandEnv("$HOME/.tesla/config.toml")

// Load loads the source into the config defined by struct s.
func (t *tomlLoader) Load(s interface{}) error {
	filename := t.Path

	if filename == "" {
		filename = defaultTomlPath
	}

	log.Debug("Trying to load configuration file", "filename", filename)

	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		if t.Required {
			return errors.New("Configuration file does not exist")
		}
		logger := log.Warn
		if filename == defaultTomlPath {
			logger = log.Debug
		}
		logger("Configuration file does not exist", "path", filename)
		return nil
	}
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err := ioutil.ReadAll(f)

	err = toml.Unmarshal(data, s)
	if err != nil {
		return err
	}

	log.Debug("Loaded configuration file", "filename", filename)
	return nil
}
