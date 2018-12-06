package emanage

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/go-errors/errors"

	"rest"
)

var (
	hostsUri     = "/api/hosts"
	syncHostsUri = path.Join(hostsUri, "sync")
	instancesUri = path.Join(hostsUri, "create_instances")
)

type hosts struct {
	conn *rest.Session
}

type DetectHostOpts struct {
	Vlan              int    `json:"vlan"`
	HostIDs           []int  `json:"host_ids"`
	BroadcastNic      string `json:"broadcast_nic,omitempty"`
	DataNetworkNumber int    `json:"data_network_number,omitempty"`
}

func (h *hosts) Detect(opts *DetectHostOpts) error {
	return h.conn.Request(rest.MethodPost, path.Join(hostsUri, "detect"), opts, nil)
}

func (h *hosts) Sync() error {
	return h.conn.Request(rest.MethodPost, syncHostsUri, nil, nil)
}

func (h *hosts) GetAll(opt *GetAllOpts) ([]Host, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []Host
	if err := h.conn.Request(rest.MethodGet, hostsUri, opt, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *hosts) GetHost(id int) (*Host, error) {
	var result Host
	err := h.conn.Request(
		rest.MethodGet,
		path.Join(hostsUri, fmt.Sprintf("%v", id)),
		nil,
		&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type UpdateHostOpts struct {
	User     string
	Password string
}

func (h *hosts) Update(id int, opts *UpdateHostOpts) error {
	hosts, err := h.GetAll(nil)
	if err != nil {
		return err
	}

	var result Host

	for _, hst := range hosts {
		if hst.ID == id {

			hst.User = opts.User
			hst.Password = opts.Password
			err := h.conn.Request(
				rest.MethodPut,
				path.Join(hostsUri, fmt.Sprintf("%v", hst.ID)),
				&hostOpts{Host: &hst},
				&result)

			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.Errorf("Host update: didn't find any host with id: %v", id)
}

func (h *hosts) Create(opts *Host) (result Host, err error) {
	err = h.conn.Request(rest.MethodPost, hostsUri, opts, &result)
	return result, err
}

type Bytes struct {
	Bytes int `json:"bytes"`
}

type Nanos struct {
	Nanos uint64 `json:"nanos"`
}

type Device struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Model      string      `json:"model"`
	Vendor     string      `json:"vendor"`
	IsWritable bool        `json:"is_writable"`
	Ssd        bool        `json:"ssd"`
	DevicePath interface{} `json:"device_path"`
	Status     string      `json:"status"`
	HostID     int         `json:"host_id"`
	EnodeID    interface{} `json:"enode_id"`
	HostName   string      `json:"host_name"`
	Capacity   struct {
		Bytes int64 `json:"bytes"`
	} `json:"capacity"`
	Usage struct {
		Bytes interface{} `json:"bytes"`
	} `json:"usage"`
	Datastore         string `json:"datastore"`
	IsLocal           bool   `json:"is_local"`
	IsInternalFs      bool   `json:"is_internal_fs"`
	IsLogFs           bool   `json:"is_log_fs"`
	PassthroughActive bool   `json:"passthrough_active"`
	URL               string `json:"url"`
}

// type host struct {
// }

type Host struct {
	Type                   string             `json:"host_type"`
	SystemId               int                `json:"system_id"`
	Cores                  int                `json:"cores"`
	Datastores             []DataStore        `json:"datastores"`
	Devices                []Device           `json:"devices"`
	DevicesCount           int                `json:"devices_count"`
	EnableSriov            interface{}        `json:"enable_sriov"`
	IsManagement           bool               `json:"is_management"`
	ID                     int                `json:"id"`
	Maintenance            bool               `json:"maintenance"`
	Memory                 int                `json:"memory"`
	Model                  string             `json:"model"`
	Name                   string             `json:"name"`
	NetworkInterfaces      []NetworkInterface `json:"network_interfaces"`
	NetworkInterfacesCount int                `json:"network_interfaces_count"`
	Networks               []struct {
		CreatedAt string `json:"created_at"`
		HostID    int    `json:"host_id"`
		ID        int    `json:"id"`
		Name      string `json:"name"`
		UpdatedAt string `json:"updated_at"`
		Vlan      int    `json:"vlan"`
		VswitchID int    `json:"vswitch_id"`
	} `json:"networks"`
	Path        string      `json:"path"`
	PowerState  string      `json:"power_state"`
	Role        string      `json:"role"`
	Software    string      `json:"software"`
	Status      string      `json:"status"`
	User        interface{} `json:"user"`
	Password    interface{} `json:"password"`
	Vendor      string      `json:"vendor"`
	VMManagerID int         `json:"vm_manager_id"`
}

func (h *Host) GetDataStoreByPrefix(prefix string) (*DataStore, error) {
	if len(prefix) == 0 {
		return nil, errors.Errorf("DataStore name is empty")
	}

	for _, d := range h.Datastores {
		if strings.HasPrefix(d.Name, prefix) {
			return &d, nil
		}
	}
	return nil, errors.Errorf("Failed to find datastore with prefix '%v'.", prefix)
}

func (h *Host) DeviceIDsByPrefixes(prefix ...string) (result []int) {
	for _, dev := range h.Devices {
		for _, prefix := range prefix {
			if strings.HasPrefix(dev.Name, prefix) {
				logger.Debug("Setup device match", "disk name", dev.Name, "prefix", prefix)
				result = append(result, dev.ID)
			}
		}
	}
	return
}

func (h *Host) DeviceIDsByFS() (result []int) {
	for _, dev := range h.Devices {
		if !dev.IsInternalFs && !dev.IsLogFs {
			logger.Debug("Host device match", "id", dev.ID, "name", dev.Name)
			result = append(result, dev.ID)
		}
	}
	return
}

func (h *Host) MacByInterfaceNetworkName(name string) string {
	for _, nif := range h.NetworkInterfaces {
		for _, net := range nif.Networks {
			if net.Name == name {
				return nif.Lladdress
			}
		}
	}
	return ""
}

func (h *Host) IsPhysical() bool {
	return h.Type == "physical"
}

func (h *Host) NetInterfaceRoles() []string {
	roles := make([]string, 0)
	for _, nif := range h.NetworkInterfaces {
		roles = append(roles, nif.Role)
	}
	return roles
}

type hostOpts struct {
	Host *Host `json:"host"`
}

type InstancesCreateOpts struct {
	Instances int  `json:"instances,omitempty"`
	Async     bool `json:"async,omitempty"`
}

type Instances struct {
	ID           int         `json:"id,omitempty"`
	UUID         string      `json:"uuid,omitempty"`
	Priority     int         `json:"priority,omitempty"`
	Attempts     int         `json:"attempts,omitempty"`
	Status       string      `json:"status,omitempty"`
	Name         string      `json:"name,omitempty"`
	CurrentStep  string      `json:"current_step,omitempty"`
	StepProgress int         `json:"step_progress,omitempty"`
	StepTotal    int         `json:"step_total,omitempty"`
	CreatedAt    time.Time   `json:"created_at,omitempty"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
	Host         interface{} `json:"host,omitempty"`
	TaskType     string      `json:"task_type,omitempty"`
	URL          string      `json:"url,omitempty"`
}

func (en *hosts) CreateInstances(opts *InstancesCreateOpts) ([]Instances, error) {
	result := []Instances{}

	err := en.conn.Request(rest.MethodPost, instancesUri, opts, &result)
	return result, err
}
