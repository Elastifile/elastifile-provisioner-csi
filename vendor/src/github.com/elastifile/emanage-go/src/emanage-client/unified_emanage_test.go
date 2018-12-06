package emanage

import (
	"fmt"
	"net/url"
	"testing"
)

/*
EManage Client Struct Members:
	ClusterReports      *clusterReports
	DataContainers      *dataContainers
	Enodes              *enodes
	ControlTasks        *controlTasks
	Exports             *exports
	Hosts               *hosts
	Events              *events
	NetworkInterfaces   *netInterfaces
	Policies            *policies
	Sessions            *rest.Session
	Devices             *devices
	Snapshots           *snapshots
	Statistics          *statistics
	Systems             *systems
	EmanageVMs          *emanageVms
	Tenants             *tenants
	VMManagers          *vmManagers
	VMs                 *vms
	Clients             *clients
	ClientNetworks      *clientNetworks
	ReplicationAgents   *replicationAgents
	RemoteSites         *remoteSites
	DcPairs             *dcPairs
	CloudProviders      *cloudProviders
	CloudConfigurations *cloudConfigurations
 */


func validateHosts(EMSClient *Client, opts *GetAllOpts) {
	hostArr, err := EMSClient.Hosts.GetAll(opts)

	if err != nil {
		fmt.Printf("Error getting all hosts: %s", err)
		return
	}
	// Validate hosts getters:
	for i, host := range hostArr {
	    fmt.Printf("Host# %v: %v, Role: %v\n", i, host.Name, host.Role)

	    // Retrieve single host:
	    hostGetSingle, err := EMSClient.Hosts.GetHost(i+1)
	    if (err != nil) { fmt.Printf("Error retrieving host information: %s\n", err) }
		if (hostGetSingle.ID != (i+1)) { fmt.Printf("Error: host ID mismatches?!\n") }

	    fmt.Printf("Host %s: cores=%d, device count=%d, memory=%d\n",
	    	hostGetSingle.Name, hostGetSingle.Cores, hostGetSingle.DevicesCount, hostGetSingle.Memory)
	}
}

func validateEnodes(EMSClient *Client, opts *GetAllOpts) {
	enodesArr, err := EMSClient.Enodes.GetAll()
	if err != nil {
		fmt.Printf("Error getting all enodes: %s", err)
		return
	}
	for i, enode := range enodesArr {
		fmt.Printf("ENode #%d: IP: %s, CPU Usage %.2f%%, Version: %s\n",
			i, enode.DataIP, enode.CpuUsage.Percent, enode.SoftwareVersion)
	}
}

func validateDataContainers(EMSClient *Client, opts *GetAllOpts) {
	dcArr, err := EMSClient.DataContainers.GetAll(nil);
	if err != nil {
		fmt.Printf("Error getting all data containers: %s", err)
		return
	}
	for i, dc := range dcArr {
		fmt.Printf("DC #%d: Name: %s, Used Bytes: %d Data Type: %s\n",
			i, dc.Name, dc.TotalUsedBytes.Bytes, dc.DataType)
	}
}

func TestCreateClient(t *testing.T) {

	fmt.Println("Starting Unified Emanage Client test")
	EMSClient := clientLogin("10.11.209.208")
	opts := &GetAllOpts{}

	validateHosts(EMSClient, opts)
	validateEnodes(EMSClient, opts)
	validateDataContainers(EMSClient, opts)
}

func clientLogin(IP string) *Client {
	eurl := &url.URL{ Scheme: "http",  Host: IP }
	EMSClient := NewClient(eurl)
	// login
	err := EMSClient.Sessions.Login("admin", "changeme")
	if err != nil {
		fmt.Printf("Error logging in: %s", err)
		return nil
	}
	fmt.Println("Logged in")
	return EMSClient
}
