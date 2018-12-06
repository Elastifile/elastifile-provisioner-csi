package runtimeutil_test

import (
	"testing"

	"runtimeutil"
)

// "runtimeutil"

func TestGetPkgName(t *testing.T) {
	myPkg := runtimeutil.PackageName()
	t.Log("Package name: ", myPkg)
	if myPkg != "runtimeutil_test" {
		t.Fatal(myPkg)
	}
}
