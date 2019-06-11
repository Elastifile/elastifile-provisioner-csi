package types_test

import (
	"fmt"
	"testing"

	"types"
)

func TestFile(t *testing.T) {
	f := types.File{
		Name:    "myfile",
		Content: []byte("mycontent is very long and tiring and therefore I am truncating it a bit"),
	}
	t.Logf("%v", f)
}

// func TestVCenter(t *testing.T) {
// 	vcUrl := "https://root:vmware@vc7a.lab.il.elastifile.com/sdk"

// 	vc := types.VCenter{
// 		Host:     "vc7a.lab.il.elastifile.com",
// 		Username: "root",
// 		Password: "vmware",
// 	}

// 	if vc.URL().String() != vcUrl {
// 		t.Fatalf("Got %v, expected %v", vc.URL(), vcUrl)
// 	}
// }

func TestLoaders(t *testing.T) {
	loaders := []string{"l1", "l2", "l3"}
	conf := &types.Config{}
	conf.SetLoaders(loaders...)

	for i, l := range conf.Loaders() {
		fmt.Printf("%d) %s\n", i, l)
	}
}
