package emanage_test

import (
	"emanage-client"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCreateInstanses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName, emanageIp)

	mgmt := startEManage(t)

	ins, err := mgmt.Hosts.CreateInstances(&emanage.InstancesCreateOpts{
		Instances: 1,
		Async:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(ins)
}

// func TestHostsUpdate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	err := mgmt.Hosts.Update(0, &emanage.UpdateHostOpts{
// 		User:     "foo",
// 		Password: "bar",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("Updated host:\n%#v", 3)
// }

// func TestHostsCreate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	hosts, err := mgmt.Hosts.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(hosts) < 0 {
// 		t.Fatal("must have at least one host configured")
// 	}
// 	host := hosts[0]
// 	t.Logf("got host: %+v", host)

// 	host.Name += fmt.Sprintf("_%d", time.Now().Unix())
// 	host, err = mgmt.Hosts.Create(&host)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("Created host:\n%#v", host)
// }

// func TestHostsDetect(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	err := mgmt.Hosts.Sync()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	vlan := 8
// 	err = mgmt.Hosts.Detect(&emanage.DetectHostOpts{Vlan: vlan, HostIDs: []int{3}})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
