package emanage_test

import (
	"emanage-client"
	"testing"

	"github.com/go-errors/errors"
)

// func TestVersionGet(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	version, err := system.GetVersion()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(version)
// }

func TestVersionCompareHigher(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	baseVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 2,
	}

	otherVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 1,
	}

	if baseVersion.Compare(&otherVersion) != 1 {
		t.Fatal(errors.Errorf("Expected '%#v' > '%#v'", baseVersion, otherVersion))
	}
}

func TestVersionCompareEqual(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	baseVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 2,
	}

	otherVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 2,
	}

	if baseVersion.Compare(&otherVersion) != 0 {
		t.Fatal(errors.Errorf("Expected '%#v' == '%#v'", baseVersion, otherVersion))
	}
}

func TestVersionCompareLower(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	baseVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 2,
	}

	otherVersion := emanage.Version{
		Major:    5,
		Minor:    4,
		Build:    3,
		Revision: 3,
	}

	if baseVersion.Compare(&otherVersion) != -1 {
		t.Fatal(errors.Errorf("Expected '%#v' < '%#v'", baseVersion, otherVersion))
	}
}
