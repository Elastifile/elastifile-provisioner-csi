package logging

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"

	tm "github.com/buger/goterm"
	"github.com/go-errors/errors"
	"github.com/koding/multiconfig"
	log "gopkg.in/inconshreveable/log15.v2"

	config "logging/config"
	"types"
)

var logger = NewLogger("logging")

var global loggingGlobalData

func init() {
	global.level = log.LvlInfo
	log.Root().SetHandler(&global)
}

func Fini() error {
	return finiTransfers()
}

type loggingGlobalData struct {
	setupOnce sync.Once
	handler   log.Handler
	level     log.Lvl
	conf      *config.Config
	datastore string

	logFilesMap map[string]log.Lvl
}

func (l *loggingGlobalData) Log(rec *log.Record) error {
	if l.handler == nil {
		return nil
	}
	return l.handler.Log(rec)
}

func (l *loggingGlobalData) init() {
	l.logFilesMap = make(map[string]log.Lvl)

	var err error
	if l.datastore, err = types.FilePathInDatastore(""); err == nil && l.datastore != "" {
		logger.Info("Got file path", "datastore", l.datastore)
	}
}

func NewLogger(pkg string) log.Logger {
	return log.New("package", pkg)
}

// Setup is the proper way to setup your logger behaviour.
func Setup(conf *config.Config) {
	global.setupOnce.Do(func() {
		global.init()

		if conf == nil {
			conf = &config.Config{}
			tagLoader := &multiconfig.TagLoader{}
			if err := tagLoader.Load(conf); err != nil {
				panic(err)
			}
		}
		global.conf = conf

		AddHandlers(getDefaultHandlers()...)

		if nh := netHandler(&global.conf.LogServer); nh != nil {
			AddHandlers(nh)
		}

		setLevelString(global.conf.Level)
	})
}

// AddHandlers is the common way for adding handlers to your logger.
func AddHandlers(handlers ...log.Handler) {
	for _, h := range handlers {
		if global.handler == nil {
			global.handler = h
		} else {
			global.handler = log.MultiHandler(global.handler, h)
		}
	}
}

func setLevelString(level string) {
	if level == "" {
		return
	}
	if l, err := log.LvlFromString(level); err != nil {
		fmt.Println(tm.Color("WARN:", tm.YELLOW), "unexpected log level:", tm.Color(level, tm.YELLOW))
	} else {
		setLevel(l)
	}
}

func setLevel(level log.Lvl) {
	global.level = level
}

func LogErrorStack(logger log.Logger, err error, kv ...interface{}) {
	logger.Error(Framed("ERROR", '*'), "caller", "")
	kv = append(kv, "caller", "")
	logger.Error(err.Error(), kv...)

	switch e := err.(type) {
	case *errors.Error:
		logger.Error(Framed("Stacktrace", '-'), "caller", "")
		for _, errStk := range strings.Split(e.ErrorStack(), "\n")[1:] {
			if errStk == "" {
				continue
			}
			if errStk[0] == '/' {
				errStk = errStk[1:]
			}
			logger.Error(errStk, "caller", "")
		}
	default:
		debug.PrintStack()
	}

	logger.Error(Framed("", '*'), "caller", "")
}
