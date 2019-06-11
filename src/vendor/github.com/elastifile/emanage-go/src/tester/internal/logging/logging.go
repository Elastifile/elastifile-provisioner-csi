package logging

import log15 "gopkg.in/inconshreveable/log15.v2"

var logger = log15.New()

func SetLogger(newLogger log15.Logger) {
	logger = newLogger
	logger.Info("tester", "logger", logger)
}

func Log() log15.Logger {
	return logger
}
