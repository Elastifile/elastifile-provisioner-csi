package emanage_test

import (
	"strings"
	"testing"

	"emanage-client"
	"logging"
	logging_config "logging/config"
)

func init() {
	logging.Setup(logging_config.ConfigForUnitTest())
}

func TestGetRemoteSites(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	remoteSites, err := mgmt.RemoteSites.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Remote sites info:\n%#v", remoteSites)
}
func TestGetRemoteSite(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	remoteSite, err := mgmt.RemoteSites.GetById(1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Remote site (id 1) info:\n%#v", remoteSite)
}

func TestCreateRemoteSite(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	createOpts := emanage.RemoteSiteOpts{
		IpAddress:      "10.11.209.222",
		Login:          "admin",
		Password:       "changeme",
		LocalLogin:     "admin",
		LocalPassword:  "changeme",
		LocalIpAddress: "10.11.209.216",
	}
	remoteSite, err := mgmt.RemoteSites.Create(createOpts)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Remote site (id 1) info:\n%#v", remoteSite)
}

func TestConnectRemoteSite(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	remoteSites, err := mgmt.RemoteSites.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	tested := false
	for _, site := range remoteSites {
		if strings.Contains(site.ConnectionStatus, "disconnected") {
			remote, err := mgmt.RemoteSites.Connect(site.ID)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(remote.ConnectionStatus, "connected") {
				t.Fatalf("Remote site not connected, status: %v", remote.ConnectionStatus)
			}
			tested = true
			break
		}
	}

	if !tested {
		t.Fatal("Can't find a disconnected remote site to connect")
	}
}

//func TestDisconnectRemoteSite(t *testing.T) {
//	t.Skip("Integration test")
//
//	mgmt := startEManage(t)
//
//	remoteSites, err := mgmt.RemoteSites.GetAll()
//	if err != nil {
//		t.Fatal(err)
//	}
//	tested := false
//	for _, site := range remoteSites {
//		if strings.Contains(site.ConnectionStatus, "connected") {
//			remote, err := mgmt.RemoteSites.Disconnect(site.ID, emanage.RemoteDisconnectOpts{})
//			if err != nil {
//				t.Fatal(err)
//			}
//			if !strings.Contains(remote.ConnectionStatus, "disconnected") {
//				t.Fatalf("Remote site not disconnected, status: %v", remote.ConnectionStatus)
//			}
//			tested = true
//			break
//		}
//	}
//
//	if !tested {
//		t.Fatal("Can't find a connected remote site to disconnect")
//	}
//}

func TestDeleteRemoteSite(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	remoteSites, err := mgmt.RemoteSites.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	remoteSite, err := mgmt.RemoteSites.Delete(remoteSites[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Remote site (id 1) info:\n%#v", remoteSite)

}

func TestDataContinure(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	remoteSites, err := mgmt.RemoteSites.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	allDcs := 0
	for _, site := range remoteSites {
		dcs, err := mgmt.RemoteSites.DataContainers(site.ID)
		if err != nil {
			t.Fatal(err)
		}
		allDcs += len(dcs)
		t.Logf("Found %d data continures on site %v", len(dcs), site.ID)
	}
}
