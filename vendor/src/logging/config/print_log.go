package logging_config

import (
	"fmt"
	"strings"
	"time"

	tm "github.com/buger/goterm"

	"helputils"
)

const (
	EROR = "EROR"
	WARN = "WARN"
	INFO = "INFO"
	DBUG = "DBUG"
)

var printLogLvlNames = []string{
	EROR,
	WARN,
	INFO,
	DBUG,
}

func (conf *Config) PrintDebug(msg string, ctx ...interface{}) {
	conf.Print(DBUG, tm.CYAN, msg, ctx...)
}

func (conf *Config) PrintInfo(msg string, ctx ...interface{}) {
	conf.Print(INFO, tm.GREEN, msg, ctx...)
}

func (conf *Config) PrintWarn(msg string, ctx ...interface{}) {
	conf.Print(WARN, tm.YELLOW, msg, ctx...)
}

func (conf *Config) PrintError(msg string, ctx ...interface{}) {
	conf.Print(EROR, tm.RED, msg, ctx...)
}

func (conf *Config) Print(lvl string, color int, msg string, ctx ...interface{}) {
	if i := helputils.FindStr(printLogLvlNames, lvl); i < 0 {
		panic("invalid log lvl name: " + lvl)
	} else if i <= helputils.FindStr(printLogLvlNames, conf.PrintLvl) {
		conf.printLogEntry(tm.Color(lvl, color), msg, ctx...)
	}
}

func (conf *Config) printLogEntry(lvl string, msg string, kv ...interface{}) {
	kvStr := make([]string, len(kv)/2)

	for i := range kvStr {
		kvStr[i] = fmt.Sprintf("%s=%+v", kv[i*2], kv[i*2+1])
	}

	fmt.Println(time.Now().Format("20060102-150405.000"), lvl+":", msg, strings.Join(kvStr, " "))
}
