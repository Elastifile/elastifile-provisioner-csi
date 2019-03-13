package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"helputils"
	config "logging/config"
	"types"
)

const logFileMode = 0777

type logTransport struct {
	lvl              log.Lvl
	locker           sync.Mutex
	localPath        string
	filerPath        string
	transferredBytes int64
}

func (trans *logTransport) String() string {
	return fmt.Sprintf("logTransport{locker:%+v localPath:%s filerPath:%s transferredBytes:%+v}",
		trans.locker, trans.localPath, trans.filerPath, trans.transferredBytes)
}

func (trans *logTransport) transfer() (size int, err error) {
	trans.locker.Lock()
	defer trans.locker.Unlock()

	var confTransLvl log.Lvl
	if confTransLvl, err = log.LvlFromString(global.conf.Filer.TransferLvl); err != nil {
		return
	} else if trans.lvl > confTransLvl {
		return
	}

	if stat, err := os.Stat(trans.localPath); os.IsNotExist(err) || stat.Size() == trans.transferredBytes {
		return 0, nil
	} else if global.conf.Filer.Verbose {
		logger.Debug("transferable log",
			"local", trans.localPath,
			"size", stat.Size(),
			"transferred", trans.transferredBytes,
			"filer", trans.filerPath,
		)
	}

	for startTime := time.Now(); time.Since(startTime) <= 5*time.Minute; time.Sleep(5 * time.Second) {
		if size, err = helputils.AppendFile(trans.filerPath, trans.localPath, trans.transferredBytes, logFileMode); err != nil {
			err = fmt.Errorf("Failed append to filer, err: %s", err)
		} else {
			trans.transferredBytes += int64(size)
			if global.conf.Filer.Verbose {
				logger.Debug("log transferred",
					"local", trans.localPath,
					"new", size,
					"total", trans.transferredBytes,
				)
			}
			break
		}
	}

	return
}

func (trans *logTransport) mustTransfer() int {
	size, err := trans.transfer()
	if err != nil {
		panic(fmt.Sprintf("%s - %+v", err, *trans))
	}
	return size
}

type logTransporters map[log.Lvl][]*logTransport

var logLvlTransports logTransporters = make(logTransporters)

func lockHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(rec *log.Record) error {
		for _, trans := range logLvlTransports[rec.Lvl] {
			locker := &trans.locker
			locker.Lock()
			defer locker.Unlock()
		}
		return h.Log(rec)
	})
}

func addTransfer(lvl log.Lvl, localPath string, filerPath string) {
	newTrans := logTransport{
		lvl:       lvl,
		localPath: localPath,
		filerPath: filerPath,
	}

	if _, err := os.Stat(filerPath); os.IsNotExist(err) {
		helputils.TouchFile(filerPath, logFileMode)
		global.conf.PrintDebug("initial touch",
			"filer", filerPath,
			"mode", fmt.Sprintf("%o", logFileMode),
		)
	}

	if _, ok := logLvlTransports[lvl]; !ok {
		logLvlTransports[lvl] = []*logTransport{&newTrans}
		go mustTransferToFilerRoutine(lvl, time.NewTicker(global.conf.Filer.Interval))
	} else {
		logLvlTransports[lvl] = append(logLvlTransports[lvl], &newTrans)
	}
}

func mustTransferToFilerRoutine(lvl log.Lvl, ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			for _, trans := range logLvlTransports[lvl] {
				trans.mustTransfer()
			}
		}
	}()
}

func finiTransfers() error {
	for _, transports := range logLvlTransports {
		for _, trans := range transports {
			if size := trans.mustTransfer(); size > 0 && global.conf.Filer.Verbose {
				logger.Debug("final log transfer",
					"size", size,
					"info", trans,
				)
			}
			if !global.conf.Filer.PersistLocal {
				if err := os.RemoveAll(trans.localPath); err != nil {
					return err
				} else if global.conf.Filer.Verbose {
					logger.Debug("removed log file",
						"local", trans.localPath,
					)
				}
			}
		}
	}

	if global.conf != nil {
		if global.conf.Filer.Enabled && global.conf.SessionId != "" {
			global.conf.PrintInfo("log transfers finalized",
				"session", global.conf.SessionPath(),
				"program", os.Args[0],
			)
		}
	}

	return nil
}

func StreamLogsToFiler() {
	for logFile, lvl := range global.logFilesMap {
		localPath, _ := types.FilePathInDatastore(logFile)
		logDir := global.conf.SessionPath()
		if global.conf.Filer.ProgramName != "" {
			logDir = extendDir(logDir, global.conf.Filer.ProgramName)
		}
		filerPath := filepath.Join(logDir, logFile)
		addTransfer(lvl, localPath, filerPath)
	}
}

func extendDir(logDir string, subdir string) string {
	logDir = filepath.Join(logDir, subdir)
	if err := helputils.MkdirAll(logDir, config.FileMode); err != nil {
		panic(err)
	}

	logger.Debug("created", "subdir", logDir)
	return logDir
}
