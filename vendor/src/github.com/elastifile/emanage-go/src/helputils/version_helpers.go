package helputils

import (
	"github.com/hashicorp/go-version"
	"strings"
)

func GetVersion(ver string) *version.Version {
	versionObj, err := version.NewVersion(ver)
	if err != nil {
		panic(err)
	}
	return versionObj
}

func ParseVersion(version string) string {
	lastIndex := strings.LastIndex(version, ".")
	return version[:lastIndex]
}
