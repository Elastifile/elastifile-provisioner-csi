package emanage_test

import (
	"emanage-client"
	"testing"
)

func TestCreateAgent(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	createOpts := emanage.ReplicationAgentsCreateOpts{HostID: 5}
	replicationAgent, err := mgmt.ReplicationAgents.Create(createOpts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("replicationAgent info:\n%#v", replicationAgent)
}
func TestDeleteAgent(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	replicationAgent, err := mgmt.ReplicationAgents.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}
	replicationAgents, err := mgmt.ReplicationAgents.Delete(replicationAgent[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("replicationAgent info:\n%#v", replicationAgents)
}

func TestGetAll(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	replicationAgents, err := mgmt.ReplicationAgents.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("replicationAgent Amount:%v", len(replicationAgents))
	for i, agent := range replicationAgents {
		t.Logf("replicationAgent %v:\n%#v", i, agent)
	}
}

func TestConfigureReplicationAgent(t *testing.T) {
	t.Skip("Integration test")

	mgmt := startEManage(t)

	var createOpts = emanage.ReplicationAgentsCreateOpts{}
	createOpts.HostID = 5
	createOpts.ClientNetworkNicId = 19
	createOpts.ClientNetworkStaticIp = "172.16.208.150"
	createOpts.ExternalNetworkNicId = 18
	createOpts.ExternalNetworkStaticIp = "10.11.147.87"
	createOpts.ExternalNetworkIpRange = 16

	replicationAgent, err := mgmt.ReplicationAgents.Create(createOpts)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("replicationAgent info:\n%#v", replicationAgent)
}
