package emanage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

type Version struct {
	Major    int
	Minor    int
	Build    int
	Revision int
}

type OptionalVersion *Version

func NewVersion(major, minor, build, revision int) OptionalVersion {
	return &Version{Major: major, Minor: minor, Build: build, Revision: revision}
}

func (v Version) String() string {
	return fmt.Sprintf("%v.%v.%v.%v",
		v.Major, v.Minor, v.Build, v.Revision)
}

func (v *Version) Compare(other *Version) int {
	var result int
	result = compare(v.Major, other.Major)
	if result != 0 {
		return result
	}
	result = compare(v.Minor, other.Minor)
	if result != 0 {
		return result
	}
	result = compare(v.Build, other.Build)
	if result != 0 {
		return result
	}
	result = compare(v.Revision, other.Revision)
	if result != 0 {
		return result
	}
	return 0
}

func compare(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

// JSON example:
// [
// {"type":"Elastifile vHead","name":"vHead-VM1","version":"2.0.0.0-34689.12d799bcdb6f.el7.centos"}
// ,
// {"type":"Elastifile vHead","name":"vHead-VM2","version":"2.0.0.0-34689.12d799bcdb6f.el7.centos"}
// ,
// {"type":"Elastifile vHead","name":"vHead-VM3","version":"2.0.0.0-34689.12d799bcdb6f.el7.centos"}
// ,
// {"type":"Elastifile vHead","name":"vHead-VM4","version":"2.0.0.0-34689.12d799bcdb6f.el7.centos"}
// ,
// {"type":"Elastifile Control Server","name":"deploy by tesla","version":"2.0.0.0-34689.12d799bcdb6f.el7.centos"}
// ]

type VersionDetails struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type VersionList []VersionDetails

func getECSVersion(versionList *VersionList) (*VersionDetails, error) {
	ecsTypeStr := "Elastifile Control Server"
	for _, versionDetails := range *versionList {
		if versionDetails.Type == ecsTypeStr {
			return &versionDetails, nil
		}
	}

	return nil, errors.Errorf("Type %v version was not found", ecsTypeStr)
}

func parseVersion(rawVersion string) (*Version, error) {
	versionStr := strings.Split(rawVersion, "-")
	if versionStr[0] == "" {
		return nil, errors.Errorf("Illegel version, found='%v'", versionStr[0])
	}

	strings := strings.Split(versionStr[0], ".")
	legelVersionCount := 4
	if len(strings) != legelVersionCount {
		return nil, errors.Errorf("Illegel version, expecting='%v', got='%v'",
			legelVersionCount, len(strings))
	}

	var ints []int
	for _, str := range strings {
		i, err := strconv.Atoi(str)
		if err != nil {
			return nil,
				errors.Errorf("Failed to parse string='%v' to int, error=%v", str, err)
		}
		ints = append(ints, i)
	}

	version := Version{
		Major:    ints[0],
		Minor:    ints[1],
		Build:    ints[2],
		Revision: ints[3],
	}

	return &version, nil
}
