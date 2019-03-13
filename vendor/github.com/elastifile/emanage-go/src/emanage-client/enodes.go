package emanage

import (
	"fmt"
	"path/filepath"

	"rest"
	"types"
)

const (
	enodesUri = "/api/enodes"
)

type enodes struct {
	conn *rest.Session
}

// type EnodeUpgradeState string

// const (
// 	EnodeUpgradeStateIdle     EnodeUpgradeState = "idle"
// 	EnodeUpgradeStateUpgrade  EnodeUpgradeState = "upgrading"
// 	EnodeUpgradeStateFail     EnodeUpgradeState = "failed"
// 	EnodeUpgradeStateStop     EnodeUpgradeState = "stopped"
// 	EnodeUpgradeStateStopping EnodeUpgradeState = "stopping"
// )

type Enode struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Cores             int                `json:"cores"`
	Memory            int                `json:"memory"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces"`
	Role              types.VHeadRole    `json:"role"`
	SetupStatus       types.SetupStatus  `json:"setup_status"`
	SoftwareVersion   string             `json:"software_version"`
	Status            types.Status       `json:"status"`
	SystemID          int                `json:"system_id"`
	ActiveConns       int                `json:"active_connections"`
	Capacity          types.Capacity     `json:"capacity"`
	PowerState        string             `json:"power_state"`
	CpuUsage          types.Usage        `json:"cpu_usage"`
	MemoryUsage       types.Usage        `json:"mem_usage"`
	UpdatedAt         string             `json:"updated_at"`
	VMFolder          string             `json:"vm_folder"`
	ActiveCo          string             `json:"active_co"`
	IsActiveCo        bool               `json:"is_active_co"`
	IsArpOfficer      bool               `json:"is_arp_officer"`
	IsActiveSo        bool               `json:"is_active_so"`
	CreatedAt         string             `json:"created_at"`
	Datastore         string             `json:"datastore"`
	DataIP            string             `json:"data_ip"`
	DataIP2           string             `json:"data_ip2"`
	DataNicStatus     string             `json:"data_nic_status"`
	DataNic2Status    string             `json:"data_nic2_status"`
	DataMAC           string             `json:"data_mac"`
	DataMAC2          string             `json:"data_mac2"`
	IsEcdb            bool               `json:"is_ecdb"`
	IsOrc             bool               `json:"is_orc"`
	NumOrcs           int                `json:"num_orcs"`
	ExternalIP        string             `json:"external_ip"`
	FrontendCores     []int              `json:"frontend_cores"`
	FEAfterPanic      bool               `json:"frontend_after_panic"`
	BackendCores      []int              `json:"backend_cores"`
	Host              Host               `json:"host"`
	StatusName        types.StatusName   `json:"status_name"`
	Devices           []struct {
		CanonicalName string      `json:"canonical_name"`
		Capacity      Bytes       `json:"capacity"`
		CreatedAt     string      `json:"created_at"`
		DevicePath    string      `json:"device_path"`
		EnodeID       int         `json:"enode_id"`
		Format        interface{} `json:"format"`
		Free          int         `json:"free"`
		HostID        int         `json:"host_id"`
		ID            int         `json:"id"`
		IsWritable    bool        `json:"is_writable"`
		Model         string      `json:"model"`
		Name          string      `json:"name"`
		PciID         string      `json:"pci_id"`
		Ssd           interface{} `json:"ssd"`
		Status        string      `json:"status"`
		UpdatedAt     string      `json:"updated_at"`
		Usage         Bytes       `json:"usage"`
		UUID          string      `json:"uuid"`
		Vendor        string      `json:"vendor"`
		VMID          int         `json:"vm_id"`
	} `json:"devices"`
	UpgradeState types.EnodeUpgradeState `json:"upgrade_state"`
}

type EnodesCreateOpts struct {
	Name        string `json:"name"`
	ExternalIP  string `json:"external_ip"`
	DataMAC     string `json:"data_mac"`
	DataMAC2    string `json:"data_mac2,omitempty"`
	DataIp      string `json:"data_ip"`
	DataIp2     string `json:"data_ip2"`
	HostID      int    `json:"host_id"`
	InternalMAC string `json:"internal_mac"`
	DeviceIDs   []int  `json:"device_ids"`
	DatastoreID int    `json:"datastore_id"`
	Role        string `json:"role"`
}

func (en *enodes) Create(opts *EnodesCreateOpts) (*Enode, error) {
	var result Enode
	if len(opts.DeviceIDs) == 0 {
		logger.Debug("new enode's device list is empty", "opts", opts)
	}
	if err := en.conn.Request(rest.MethodPost, enodesUri, opts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (en *enodes) Delete(enode *Enode) (tasks AsyncTasks) {
	fullUri := filepath.Join(enodesUri, fmt.Sprintf("%v", enode.ID))
	var err error
	if tasks.taskIDs, err = en.conn.AsyncRequest(rest.MethodDelete, fullUri, nil); err != nil {
		tasks.err = err
	} else {
		tasks.conn = en.conn
	}
	return tasks
}

func (en *enodes) GetAll() ([]Enode, error) {
	var result []Enode
	if err := en.conn.Request(rest.MethodGet, enodesUri, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (en Enode) IsPoweredOn() bool {
	return en.PowerState == "poweredOn"
}
