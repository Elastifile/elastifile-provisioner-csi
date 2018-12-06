package emanage

import (
	"fmt"
	"time"

	"rest"
)

const replicationAgentsUri = "/api/replication_agents"

type replicationAgents struct {
	session *rest.Session
}

type ReplicationAgent struct {
	ID                       int                `json:"id"`
	Name                     string             `json:"name"`
	SystemID                 int                `json:"system_id"`
	VMID                     int                `json:"vm_id"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	HostName                 string             `json:"host_name"`
	Host                     replicationHost    `json:"host"`
	SetupStatus              string             `json:"setup_status"`
	Datastore                interface{}        `json:"datastore"`
	ParamsValidForDeployment bool               `json:"params_valid_for_deployment"`
	ClientNetworkNic         clientNetworkNic   `json:"client_network_nic"`
	ClientNetworkStaticIP    string             `json:"client_network_static_ip"`
	ClientNetworkActualIP    string             `json:"client_network_actual_ip"`
	ExternalNetworkNic       externalNetworkNic `json:"external_network_nic"`
	ExternalNetworkVlan      int                `json:"external_network_vlan"`
	ExternalNetworkStaticIP  string             `json:"external_network_static_ip"`
	ExternalNetworkIPRange   int                `json:"external_network_ip_range"`
	ExternalNetworkGatewayIP string             `json:"external_network_gateway_ip"`
	ExternalNetwork          string             `json:"external_network"`
	ExternalNetworkID        int                `json:"external_network_id"`
	ExternalNetworkIsDhcp    bool               `json:"external_network_is_dhcp"`
	ExternalNetworkActualIP  string             `json:"external_network_actual_ip"`
	Status                   string             `json:"status"`
}

func (repl *ReplicationAgent) IsRunning() bool {
	return repl.Status == "running"
}

type replicationHost struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Status           string      `json:"status"`
	HostType         string      `json:"host_type"`
	Role             string      `json:"role"`
	EnodeSetupStatus string      `json:"enode_setup_status"`
	PowerState       string      `json:"power_state"`
	Path             string      `json:"path"`
	Cluster          interface{} `json:"cluster"`
	DataCenter       interface{} `json:"data_center"`
	Model            interface{} `json:"model"`
	CPUType          string      `json:"cpu_type"`
	Vendor           string      `json:"vendor"`
	Cores            int         `json:"cores"`
	Memory           int64       `json:"memory"`
	Maintenance      bool        `json:"maintenance"`
	VMManagerID      int         `json:"vm_manager_id"`
	Software         string      `json:"software"`
	User             string      `json:"user"`
	EnableSriov      bool        `json:"enable_sriov"`
	IsManagement     bool        `json:"is_management"`
	CPUUsage         struct {
		Percent float64 `json:"percent"`
	} `json:"cpu_usage"`
	MemUsage struct {
		Percent float64 `json:"percent"`
	} `json:"mem_usage"`
	IsDeployed                      bool `json:"is_deployed"`
	IsReplicationAgentDeployed      bool `json:"is_replication_agent_deployed"`
	DevicesCount                    int  `json:"devices_count"`
	HddCount                        int  `json:"hdd_count"`
	NetworkInterfacesCount          int  `json:"network_interfaces_count"`
	NetworkBandwidth                int  `json:"network_bandwidth"`
	HighSpeedNetworkInterfacesCount int  `json:"high_speed_network_interfaces_count"`
	EnodeID                         int  `json:"enode_id"`
	ReplicationAgentID              int  `json:"replication_agent_id"`
	MoID                            int  `json:"mo_id"`
	HyperThreadingActive            bool `json:"hyper_threading_active"`
}

type clientNetworkNic struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Speed          interface{}   `json:"speed"`
	Status         string        `json:"status"`
	HostID         int           `json:"host_id"`
	EnodeID        int           `json:"enode_id"`
	Lladdress      string        `json:"lladdress"`
	VfuncID        int           `json:"vfunc_id"`
	Role           string        `json:"role"`
	ClientNetwork  bool          `json:"client_network"`
	DetectedByNic1 string        `json:"detected_by_nic1"`
	DetectedByNic2 string        `json:"detected_by_nic2"`
	VswitchID      int           `json:"vswitch_id"`
	VswitchName    string        `json:"vswitch_name"`
	SriovEnabled   bool          `json:"sriov_enabled"`
	SriovCapable   bool          `json:"sriov_capable"`
	SriovActive    bool          `json:"sriov_active"`
	Networks       []interface{} `json:"networks"`
}

type externalNetworkNic struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Speed          interface{}   `json:"speed"`
	Status         string        `json:"status"`
	HostID         int           `json:"host_id"`
	EnodeID        int           `json:"enode_id"`
	Lladdress      string        `json:"lladdress"`
	VfuncID        int           `json:"vfunc_id"`
	Role           string        `json:"role"`
	ClientNetwork  interface{}   `json:"client_network"`
	DetectedByNic1 interface{}   `json:"detected_by_nic1"`
	DetectedByNic2 interface{}   `json:"detected_by_nic2"`
	VswitchID      int           `json:"vswitch_id"`
	VswitchName    string        `json:"vswitch_name"`
	SriovEnabled   bool          `json:"sriov_enabled"`
	SriovCapable   bool          `json:"sriov_capable"`
	SriovActive    bool          `json:"sriov_active"`
	Networks       []interface{} `json:"networks"`
}

type ReplicationAgentsCreateOpts struct {
	HostID                   int    `json:"host_id,omitempty"`
	ClientNetworkNicId       int    `json:"client_network_nic_id,omitempty"`
	ClientNetworkStaticIp    string `json:"client_network_static_ip,omitempty"`
	ExternalNetworkNicId     int    `json:"external_network_nic_id,omitempty"`
	ExternalNetworkVlan      int    `json:"external_network_vlan,omitempty"`
	ExternalNetworkStaticIp  string `json:"external_network_static_ip,omitempty"`
	ExternalNetworkIpRange   int    `json:"external_network_ip_range,omitempty"`
	ExternalNetworkGatewayIp string `json:"external_network_gateway_ip,omitempty"`
	ExternalNetworkId        int    `json:"external_network_id,omitempty"`
	ExternalNetworkIsDhcp    bool   `json:"external_network_is_dhcp,omitempty"`
	DatastoreId              int    `json:"datastore_id,omitempty"`
}

func (ra *replicationAgents) GetAll(opt *GetAllOpts) (result []ReplicationAgent, err error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}
	err = ra.session.Request(rest.MethodGet, replicationAgentsUri, opt, &result)
	return result, err
}

func (ra *replicationAgents) GetById(agentID int) (result ReplicationAgent, err error) {
	uri := fmt.Sprintf("%s/%d", replicationAgentsUri, agentID)
	err = ra.session.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

func (ra *replicationAgents) Create(createOpts ReplicationAgentsCreateOpts) (result ReplicationAgent, err error) {
	err = ra.session.Request(rest.MethodPost, replicationAgentsUri, createOpts, &result)
	return result, err
}

func (ra *replicationAgents) Update(agentID int, createOpts ReplicationAgentsCreateOpts) (result ReplicationAgent, err error) {
	uri := fmt.Sprintf("%s/%d", replicationAgentsUri, agentID)
	err = ra.session.Request(rest.MethodPost, uri, createOpts, &result)
	return result, err
}

func (ra *replicationAgents) Delete(agentId int) (result ReplicationAgent, err error) {
	uri := fmt.Sprintf("%s/%d", replicationAgentsUri, agentId)
	err = ra.session.Request(rest.MethodDelete, uri, nil, &result)
	return result, err
}
