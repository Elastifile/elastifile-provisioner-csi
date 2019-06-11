package types

import (
	"fmt"
	"strings"
)

var subpath = ""

type CodeLocation struct {
	FileName       string
	LineNumber     int
	FullStackTrace string
}

func (codeLocation CodeLocation) String() string {
	path := codeLocation.FileName
	if subpath != "" && strings.Contains(codeLocation.FileName, subpath) {
		fmt.Println(codeLocation, subpath)
		splits := strings.Split(path, subpath)
		path = splits[len(splits)-1]
	}
	return fmt.Sprintf("%s:%d", path, codeLocation.LineNumber)
}

func SetSubpath(subpath string) {
	subpath = subpath
}
