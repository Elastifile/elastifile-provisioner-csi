package logging

import (
	"fmt"
	"log/syslog"
	"strconv"
	"strings"

	colorable "github.com/mattn/go-colorable"
	log "gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2_ild/stack"

	"helputils"
	config "logging/config"
	"types"
)

const (
	syslogTag = "Tesla"

	entryFormatMono = iota
	entryFormatColor
	entryFormatHtml
)

// Terminal log handler
// Used by main.go:init(), Keep It Simple (Stupid)
func TerminalHandler() log.Handler {
	terminalHandler := log.StreamHandler(colorable.NewColorableStdout(), coloredEntryFormat(msgPadding))
	terminalHandler = callerHandler(terminalHandler)
	terminalHandler = log.FilterHandler(logLevelFilter, terminalHandler)
	return terminalHandler
}

// Network log handler
func netHandler(conf *config.NetLogger) log.Handler {
	if conf.Host == "" || conf.Port == 0 {
		return nil
	}

	protocol := conf.Protocol
	if protocol == "" {
		protocol = "udp"
	}

	connStr := conf.Host + ":" + strconv.Itoa(conf.Port)
	retNetHandler, err := log.NetHandler(protocol, connStr, log.JsonFormat())
	if err != nil {
		logger.Warn("cannot open net logger handler", "conn", connStr, "err", err)
		return nil
	}

	return callerHandler(retNetHandler)
}

// Remote syslog handler
func sysLogHandler(net, addr string, port int, frmt log.Format) (log.Handler, error) {
	dest := fmt.Sprintf("%s:%d", addr, port)
	return log.SyslogNetHandler(net, dest, syslog.LOG_INFO, syslogTag, frmt)
}

// getSysLogHandlers supports multiple syslog targets, sets a proper format + log lvl filter.
func getSysLogHandlers(lvl log.Lvl, port int, addrs ...string) ([]log.Handler, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("No address is given for remote syslog, len(addrs) == 0")
	}

	var (
		handlers = make([]log.Handler, 0)
		net      = "tcp"
	)

	log.Debug("Trying to set remote syslog handler", "remote hosts", addrs)

	for _, addr := range addrs {
		handler, err := sysLogHandler(
			net, addr, port, monochromeEntryFormat(noPadding))
		if err != nil {
			return handlers, fmt.Errorf(
				"Failed setting remote syslog logger on target: %v\nerr: %v", addr, err)
		}
		h := log.LvlFilterHandler(lvl, handler)
		handlers = append(handlers, h)
	}
	return handlers, nil
}

type fileOpts struct {
	lvl    log.Lvl
	format int
}

func fileHandler(opts fileOpts) log.Handler {
	fileExt := config.LogFileExt
	switch opts.format {
	case entryFormatColor:
		fileExt = "colored." + config.LogFileExt
	case entryFormatHtml:
		fileExt = config.LogFileExt + ".html"
	}
	logFile := fmt.Sprintf(config.LogFileNameFmt, strings.ToUpper(opts.lvl.String()), fileExt)

	global.logFilesMap[logFile] = opts.lvl

	localPath, _ := types.FilePathInDatastore(logFile)

	entryFormatter := monochromeEntryFormat
	switch opts.format {

	case entryFormatColor:
		entryFormatter = coloredEntryFormat

	case entryFormatHtml:
		entryFormatter = htmlEntryFormat
		_ = helputils.WriteFile(localPath, []byte(htmlBody), config.FileMode)
	}

	fh := log.Must.FileHandler(localPath, entryFormatter(msgPadding))
	if global.conf.Filer.Enabled {
		fh = lockHandler(fh)
	}
	fh = callerHandler(fh)
	return log.LvlFilterHandler(opts.lvl, fh)
}

// Caller log Handler
func callerHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(rec *log.Record) error {
		call := stack.Call(rec.Call.Frame().PC)
		caller := fmt.Sprintf("%+v", call)
		caller = strings.TrimPrefix(caller, "elastifile/tesla/")

		if ShowCaller(caller, rec.Lvl) {
			rec.Ctx = append(rec.Ctx, "caller", caller)
		}

		return h.Log(rec)
	})
}

func ShowCaller(caller string, lvl log.Lvl) bool {
	return !isNoCallerPkg(caller) && (isCallerLogLvl(lvl) || isCallerPkg(caller))
}

func isCallerLogLvl(lvl log.Lvl) bool {
	if global.conf == nil {
		return false
	}
	for _, lvlName := range strings.Split(global.conf.CallerLogLvlNames, "|") {
		if lvlName == "" {
			continue
		}
		if callerLvl, err := log.LvlFromString(lvlName); err != nil {
			panic(err.Error())
		} else if lvl == callerLvl {
			return true
		}
	}
	return false
}

func isCallerPkg(caller string) bool {
	if global.conf == nil {
		return false
	}
	return isContainingAnyPkg(caller, global.conf.CallerPkgs)
}

func isNoCallerPkg(caller string) bool {
	if global.conf == nil {
		return false
	}
	return isContainingAnyPkg(caller, global.conf.NoCallerPkgs)
}

func isContainingAnyPkg(caller string, pkgs string) bool {
	for _, pkg := range strings.Split(pkgs, "|") {
		if pkg != "" && strings.Contains(caller, pkg) {
			return true
		}
	}
	return false
}

//////

// Service functions

func getDefaultHandlers() []log.Handler {
	handlers := []log.Handler{TerminalHandler()}

	if global.conf.LogToFile {
		handlers = append(handlers,
			fileHandler(fileOpts{lvl: log.LvlError, format: entryFormatMono}),
			fileHandler(fileOpts{lvl: log.LvlError, format: entryFormatColor}),
			fileHandler(fileOpts{lvl: log.LvlError, format: entryFormatHtml}),
			fileHandler(fileOpts{lvl: log.LvlInfo, format: entryFormatMono}),
			fileHandler(fileOpts{lvl: log.LvlInfo, format: entryFormatColor}),
			fileHandler(fileOpts{lvl: log.LvlInfo, format: entryFormatHtml}),
			fileHandler(fileOpts{lvl: log.LvlDebug, format: entryFormatMono}),
		)
	}
	return handlers
}

// AddSysLogHandlers adds syslog handlers to the global used handlers
func AddSysLogHandlers(port int, addrs ...string) error {
	sysLogHandlers, err := getSysLogHandlers(global.level, port, addrs...)
	if err != nil {
		return err
	}
	AddHandlers(sysLogHandlers...)
	return nil
}

func logLevelFilter(rec *log.Record) bool {
	return rec.Lvl <= global.level
}
