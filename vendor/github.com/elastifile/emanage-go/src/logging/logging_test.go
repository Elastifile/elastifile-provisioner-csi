package logging

import (
	"fmt"
	"os"
	"testing"
	"time"

	tm "github.com/buger/goterm"
	log "gopkg.in/inconshreveable/log15.v2"

	logging_config "logging/config"
	"types"
)

func TestLogger(t *testing.T) {
	conf, _ := logging_config.NewConfig()
	conf.LogToSyslog = false
	conf.StartSession("", "")
	t.Logf("%+v", conf)
	Setup(&conf)

	for level := log.LvlCrit; level <= log.LvlDebug; level++ {
		setLevel(level)
		logAllLvls("min level: " + level.String())
		fmt.Println()
	}
}

func logAllLvls(msg string) {
	logger.Crit(msg)
	logger.Error(msg)
	logger.Warn(msg)
	logger.Info(msg)
	logger.Info(msg)
	logger.Debug(msg)
}

func TestNoColor(t *testing.T) {
	for c := tm.RED; c < tm.WHITE; c++ {
		cStr := fmt.Sprintf("%d", c)
		cColored := tm.Color(cStr, c)
		cBoldColored := tm.Bold(cColored)
		unColored := noColor(cBoldColored, 0)
		if unColored != cStr {
			t.Fatalf("bad uncoloring: \"%s\" != \"%s\"", unColored, cStr)
		}
		if lenNoColor(cBoldColored) != len(cStr) {
			t.Fatalf("lenNoColor(\"%s\") != len(\"%s\")", cColored, cStr)
		}
		t.Logf("str: %s (%d), bold-colored: %s (%d|%d), uncolored: %s (%d)",
			cStr, len(cStr), cBoldColored, len(cBoldColored), lenNoColor(cBoldColored), unColored, len(unColored))
	}
}

func TestHtmlColor(t *testing.T) {
	for color := tm.RED; color < tm.WHITE; color++ {
		logRec := log.Record{
			Lvl:  log.LvlInfo,
			Time: time.Now(),
			Msg:  "Hello, " + tm.Bold("world!"),
			Ctx:  []interface{}{"key1", "val1", "key2", "val2"},
		}
		buf := recordToBytes(&logRec, color, msgPadding)
		t.Log(htmlColor(string(buf)))
	}
}

func TestLoggingConfig(t *testing.T) {
	if conf, err := fromTags(); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("%+v", conf.Logging)
	}
}

func fromTags() (types.Config, error) {
	conf := types.NewConfig()
	Setup(&conf.Logging)
	return *conf, nil
}

func TestMount(t *testing.T) {
	t.Skip()
	if conf, err := fromTags(); err != nil {
		t.Fatal(err.Error())
	} else {
		if err := conf.Logging.Filer.Mount(); err != nil {
			t.Fatal(err.Error())
		} else {
			t.Logf("mounted %+v", conf.Logging.Filer)
		}
	}
}

func TestSysLogger(t *testing.T) {
	const envName = "TESLA_EMANAGE_SERVER"

	host := os.Getenv(envName)
	if host == "" {
		t.Skipf("Environment variable %v not set", envName)
	}

	logger := log.New()
	global.level = log.LvlInfo
	sysLogPort := 514

	err := AddSysLogHandlers(sysLogPort, host)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		logger.Info("Info msg", "key1", "A", "key2", "B")
		time.Sleep(time.Second * 1)
	}
}
