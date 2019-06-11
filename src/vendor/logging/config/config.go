package logging_config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/koding/multiconfig"

	"helputils"
)

const (
	LogFileNameFmt = "tesla.%s.%s"
	LogFileExt     = "log"
)

type NetLogger struct {
	Protocol string
	Host     string
	Port     int
}

type LogFiler struct {
	Enabled      bool          `default:"true"`
	MountFlags   string        `default:"-o soft,nolock"`
	NFSShare     string        `default:"file5.il.elastifile.com"`
	NFSExport    string        `default:"/mnt/ssd-6T/tesla"`
	LocalMount   string        `default:"/mnt/tesla"`
	TransferLvl  string        `default:"debug"`
	Interval     time.Duration `default:"5s"`
	MountTimeout time.Duration `default:"5m"`
	PersistLocal bool          `default:"true"`
	Verbose      bool          `default:"false"`
	Remount      bool          `default:"false"`
	ProgramName  string
}

const (
	SessionSubpathFormat = "sessions/{session}"
	osModeSticky         = 01000
	FileMode             = osModeSticky + os.ModePerm
)

type Config struct {
	Level             string `default:"info"`
	CallerPkgs        string `default:""`
	NoCallerPkgs      string `default:"stenographer"`
	CallerLogLvlNames string `default:"crit|error|info|debug"`
	LogToFile         bool   `default:"true"`
	LogToSyslog       bool   `default:"true"`
	RemoteSysLogPort  int    `default:"514"`
	SessionId         string
	Filer             LogFiler
	LogServer         NetLogger
	PrintLvl          string `default:"INFO"`
}

func NewConfig() (conf Config, err error) {
	tagLoader := &multiconfig.TagLoader{}
	return conf, tagLoader.Load(&conf)
}

func ConfigForUnitTest() *Config {
	utConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}

	utConfig.Level = "debug"
	return &utConfig
}

func (conf *Config) SessionSubpath() string {
	if conf.SessionId == "" {
		panic("session ID not set")
	}
	return strings.Replace(SessionSubpathFormat, "{session}", conf.SessionId, 1)
}

func (conf *Config) SessionPath() string {
	sessSubpath := conf.SessionSubpath()
	return filepath.Join(conf.Filer.LocalMount, sessSubpath)
}

func (conf *Config) StartSession(id string, tag string) error {
	if id == "" {
		sessionID, err := newSessionId(tag)
		if err != nil {
			return err
		}
		conf.SessionId = sessionID
		for attempts := 10; ; attempts -= 1 {
			if _, err := os.Stat(conf.SessionPath()); err != nil {
				if os.IsNotExist(err) {
					break
				} else {
					return errors.Wrap(err, 0)
				}
			} else if attempts == 0 {
				return errors.New(conf.sessionExistsMsg())
			}
			println(conf.sessionExistsMsg())
			time.Sleep(time.Millisecond)
			sessionID, err := newSessionId(tag)
			if err != nil {
				return err
			}
			conf.SessionId = sessionID
		}
	} else {
		conf.SessionId = id
		if tag != "" {
			conf.SessionId += "_" + tag
		}
		if _, err := os.Stat(conf.SessionPath()); err == nil {
			return errors.New(conf.sessionExistsMsg())
		}
	}

	logFilePattern := fmt.Sprintf(LogFileNameFmt, "*", LogFileExt+".*")
	return helputils.RemoveMatches(logFilePattern)
}

func (conf *Config) sessionExistsMsg() string {
	return "session " + conf.SessionId + " already exists"
}

func newSessionId(tag string) (string, error) {
	currUser, err := helputils.CurrentUser()
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	newId := currUser.Username + "_" + time.Now().Format("20060102-150405.000")
	if tag != "" {
		newId += "_" + tag
	}
	return newId, nil
}

func (lf *LogFiler) Mount() error {
	isRootUser, err := helputils.IsRootUser()
	if err != nil {
		return err
	}
	script := lf.MountScriptStr(isRootUser)

	if out, err := helputils.ExecuteShellString(script); err != nil {
		println("EROR:", err.Error(), "\nscript:\n", script, "\nout:\n", string(out))
		return err
	}

	hostname, _ := os.Hostname()
	fmt.Println("mounted", *lf, "@"+hostname)
	return nil
}

func (lf *LogFiler) MountScriptStr(root bool) string {
	script := `
set -e
set -x
mkdir -p ` + lf.LocalMount + `
if [ -z "$(mount | grep ` + lf.LocalMount + `)" ] ; then
	mount ` + fmt.Sprintf("%s %s:%s %s", lf.MountFlags, lf.NFSShare, lf.NFSExport, lf.LocalMount) + `
fi
`
	if !root {
		script = lf.addSudo(script)
	}

	return script
}

func (lf *LogFiler) UmountScriptStr(root bool) string {
	script := `
set -e
set -x
if [ -n "$(mount | grep ` + lf.LocalMount + `)" ] ; then
    sudo umount ` + lf.LocalMount + `
fi
`
	if !root {
		script = lf.addSudo(script)
	}

	return script
}

func (lf *LogFiler) Umount(root bool) error {
	script := lf.UmountScriptStr(root)

	if out, err := helputils.ExecuteShellString(script); err != nil {
		println("EROR:", err.Error(), "\nout:", string(out))
		return err
	}

	return nil
}

func (lf *LogFiler) addSudo(script string) string {
	sudoRe := regexp.MustCompile("^(mount |umount |mkdir )")

	lines := helputils.MapStr(
		func(line string) string {
			if sudoRe.MatchString(strings.TrimSpace(line)) {
				line = "sudo " + line
			}
			return line
		},
		strings.Split(script, "\n"),
	)

	return strings.Join(lines, "\n")
}
