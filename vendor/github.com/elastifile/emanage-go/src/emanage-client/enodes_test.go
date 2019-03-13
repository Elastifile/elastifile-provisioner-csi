package emanage_test

import (
	"net/url"
	"os"
	"testing"
)

func getEmanageAddress(t *testing.T) *url.URL {
	const envName = "TESLA_EMANAGE_SERVER"
	host := os.Getenv(envName)
	if host == "" {
		t.Fatalf("Environment variable %v not set", envName)
	}

	baseURL := &url.URL{
		Scheme: "http",
		Host:   host,
	}

	return baseURL
}

// func TestEnodesGetAllNames(t *testing.T) {
// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	enodes, err := mgmt.Enodes.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	eNames := make([]string, len(enodes))
// 	for i, en := range enodes {
// 		eNames[i] = en.Name
// 	}
// 	t.Logf("%d enode, names:\n%s", len(enodes), strings.Join(eNames, "\n"))
// }

// func TestEnodesGetFirstInfo(t *testing.T) {
// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	enodes, err := mgmt.Enodes.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	enode := enodes[0]
// 	t.Logf("first enode info: %#v", enode)
// 	t.Logf("Is powered on: %v", enode.IsPoweredOn())
// 	t.Logf("Name: %s", enode.Name)
// 	t.Logf("Role: %s", enode.Role)
// 	t.Logf("Status: %s", enode.Status)
// 	t.Logf("Power: %s", enode.PowerState)
// 	t.Logf("Memory Util.: %.2f %%", enode.MemoryUsage.Percent)
// 	t.Logf("CPU Util.: %.2f %%", enode.CpuUsage.Percent)
// 	t.Logf("Data Networks: %s %s", enode.DataNicStatus, enode.DataNic2Status)
// 	t.Logf("Connections: %d", enode.ActiveConns)
// 	t.Logf("Data Devices: %d (returns EMPTY from eManage REST!)", len(enode.Devices))
// 	t.Logf("Effective Capacity: %d b (%d GB)", enode.Capacity.Bytes, enode.Capacity.Bytes>>30)
// }

// func TestEnodesCreate(t *testing.T) {
// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	enodes, err := mgmt.Enodes.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	enode := enodes[len(enodes)-1]

// 	opts := EnodesCreateOpts{
// 		Name:       enode.Name + "~",
// 		ExternalIP: enode.ExternalIP,
// 		DataMAC:    enode.DataMAC,
// 		DataMAC2:   enode.DataMAC2,
// 		DataIp:     enode.DataIP,
// 		DataIp2:    enode.DataIP2,
// 		HostID:     enode.Host.ID,
// 		// InternalMAC: enode.InternalMAC,
// 		// DatastoreID: enode.DatastoreID,
// 		Role: string(enode.Role),
// 	}
// 	opts.DeviceIDs = make([]int, len(enode.Devices))
// 	for i, dev := range enode.Devices {
// 		opts.DeviceIDs[i] = dev.ID
// 	}
// 	newEnode, err := mgmt.Enodes.Create(&opts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("%+v", *newEnode)
// }
