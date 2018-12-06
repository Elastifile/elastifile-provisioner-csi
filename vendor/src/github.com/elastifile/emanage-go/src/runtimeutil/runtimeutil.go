package runtimeutil

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	teslaTopDir = "tesla/src/elastifile"
)

// Returns the caller current package name
func PackageName() string {
	pc, _, _, _ := runtime.Caller(1)
	callerFullPath := runtime.FuncForPC(pc).Name() // e.g. 'elastifile/tesla/pkgName.pkgFuncCaller'
	// Need to get the 'pkgName' so ...
	basePath := filepath.Base(callerFullPath)       // e.g. 'pkgName.pkgFuncCaller'
	suffix := filepath.Ext(basePath)                // e.g. suffix -> 'pkgFuncCaller'
	pkgName := strings.TrimSuffix(basePath, suffix) // e.g. after trimming suffix -> 'pkgName'
	return pkgName
}

func CallerString(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1) // Caller(0) is runtime.Caller() itself, so we add 1 for 'here'
	if !ok {
		return ""
	}
	if strings.Contains(file, teslaTopDir) {
		file = strings.Split(file, teslaTopDir+"/")[1]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
