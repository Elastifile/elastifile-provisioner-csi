package emanage_test

import (
	"emanage-client"
	"testing"
)

func TestDcsPairs(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)
	dcs, err := mgmt.DataContainers.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Data containers: %v", len(dcs))

	for _, dc := range dcs {
		pairs, err := mgmt.DcPairs.GetPairs(dc.Id)
		if err != nil {
			t.Fatal(err)
		}
		if len(pairs) == 0 {
			t.Logf("DC - %v: Not paired", dc.Name)
			continue
		}

		for _, pair := range pairs {
			t.Logf("DC: %v, Pair ID: %v connected to Remote DC ID: %v", dc.Name, pair.ID, pair.RemoteDcID)
			data, err := mgmt.DcPairs.GetById(pair.DataContainerID, pair.ID)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Updated pair %v", pair.ID)
			t.Logf("DC: %v connected to Remote DC ID: %v", dc.Name, data.RemoteDcID)
		}
	}
}

func TestCreateConnectPair(t *testing.T) {
	t.Skip("Integration test")

	remoteSite := 2
	dcName := "AutoTest"

	mgmt := startEManage(t)
	dcs, err := mgmt.DataContainers.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Data containers: %v", len(dcs))

	var dc emanage.DataContainer
	for _, d := range dcs {
		if d.Name == dcName {
			exports, err := mgmt.Exports.GetAll(nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, exp := range exports {
				if exp.DataContainerId == d.Id {
					exp, err = mgmt.Exports.Delete(&exp)
					if err != nil {
						t.Fatal(err)
					}
					t.Logf("Export %v deleted", exp.Id)
				}
			}
			dc, err = mgmt.DataContainers.Delete(&d)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("DC %v deleted", dc.Id)
		}
	}

	t.Logf("Pairing DC %v to remote site %v", dc.Id, remoteSite)

	dc, err = mgmt.DataContainers.Create("AutoTest", 1, &emanage.DcCreateOpts{Name: dcName, Dedup: 0, Compression: 0})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created dc: %v ID: %v", dc.Name, dc.Id)
	pairCreateOpts := emanage.PairCreateOpts{RemoteSiteId: "2", Rpo: 30, DrRole: emanage.DrRoleActive}
	pair, err := mgmt.DcPairs.CreatePair(dc.Id, pairCreateOpts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Paired Dc: %v to remote site: %v, Connection status: %v", dc.Id, remoteSite, pair.ConnectionStatus)
	pair, err = mgmt.DcPairs.Connect(pair.DataContainerID, pair.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Connected Dc: %v to remote site: %v, Connection status: %v", dc.Id, remoteSite, pair.ConnectionStatus)
}

func TestDCTset(t *testing.T) {
	t.Skip("Integration test")

	remoteSite := 2
	dcName := "TestDC-a"

	mgmt := startEManage(t)
	dc, err := mgmt.DataContainers.Create(dcName, 1, &emanage.DcCreateOpts{Name: dcName, Dedup: 0, Compression: 0})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created dc: %v ID: %v", dc.Name, dc.Id)
	pairCreateOpts := emanage.PairCreateOpts{RemoteSiteId: "2", Rpo: 30, DrRole: emanage.DrRoleActive}
	pair, err := mgmt.DcPairs.CreatePair(dc.Id, pairCreateOpts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Paired Dc: %v to remote site: %v, Connection status: %v", dc.Id, remoteSite, pair.ConnectionStatus)
	pair, err = mgmt.DcPairs.Connect(dc.Id, pair.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Connected Dc: %v to remote site: %v, Connection status: %v", dc.Id, remoteSite, pair.ConnectionStatus)
	aDC, err := mgmt.DcPairs.TestImage(dc.Id, pair.ID, emanage.TestDataContainereOpts{DataContainerName: dcName})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Test image: %v", aDC)
}
