package emanage

import (
	"fmt"

	"github.com/go-errors/errors"

	"helputils"
	"rest"
)

const emanageVMsUri = "/api/emanage_vms"
const monitorUri = "monitor"

type EmanageVms struct {
	ID               int    `json:"id,omitempty"`
	VMID             int    `json:"vm_id,omitempty"`
	HostDataNicName  string `json:"host_data_nic_name,omitempty"`
	HostDataNicName2 string `json:"host_data_nic_name2,omitempty"`
	HostID           int    `json:"host_id,omitempty"`
	VIP              string `json:"vip,omitempty"`
	VIPMacAddress    string `json:"vip_mac_address,omitempty"`
	VIPNetmask       string `json:"vip_netmask,omitempty"`
	HasVip           bool   `json:"has_vip"`
	Name             string `json:"name"`
	IP               string `json:"ip,omitempty"`
	User             string `json:"user,omitempty"`
	Password         string `json:"password,omitempty"`
	Status           string `json:"status,omitempty"`
	IsActive         bool   `json:"is_active,omitempty"`
	IsLocal          bool   `json:"is_local,omitempty"`
	DataIP           string `json:"data_ip,omitempty"`
	DataIP2          string `json:"data_ip2,omitempty"`
	URL              string `json:"url,omitempty"`
}

type VmStatus string

const (
	VmOnline  VmStatus = "online"
	VmOffline VmStatus = "offline"
	VmInit    VmStatus = "init"
)

func (v VmStatus) String() string {
	switch v {
	case VmOnline:
		return "online"
	case VmOffline:
		return "offline"
	case VmInit:
		return "init"
	default:
		return "Illegal vm status"
	}
}

/*
GET /api/emanage_vms
200
[
  {
    "id": 1,
    "vm_id": null,
    "host_data_nic_name": null,
    "host_data_nic_name2": null,
    "host_id": null,
    "vip": null,
    "vip_mac_address": null,
    "vip_netmask": null,
    "ip": null,
    "user": null,
    "status": "init",
    "is_active": false,
    "is_local": true,
    "data_ip": null,
    "data_ip2": null,
    "url": "http://test.host/api/emanage_vms/1"
  }
]
*/

/*
[
  {
    "id": 1,
    "vm_id": 1,
    "host_data_nic_name": "vmnic2",
    "host_data_nic_name2": "vmnic3",
    "host_id": 2,
    "vip": "10.11.209.11",
    "vip_mac_address": "00:50:56:a9:fe:5e",
    "vip_netmask": "255.255.255.0",
    "ip": null,
    "user": null,
    "status": "online",
    "is_active": true,
    "is_local": true,
    "data_ip": "11.0.0.100",
    "data_ip2": "11.10.0.227",
    "url": "http://10.11.209.11/api/emanage_vms/1"
  },
  {
    "id": 2,
    "vm_id": 6,
    "host_data_nic_name": null,
    "host_data_nic_name2": null,
    "host_id": 3,
    "vip": "10.11.209.11",
    "vip_mac_address": "00:50:56:a9:e5:cb",
    "vip_netmask": "255.255.255.0",
    "ip": "10.11.196.184",
    "user": "admin",
    "status": "online",
    "is_active": false,
    "is_local": false,
    "data_ip": "11.0.0.105",
    "data_ip2": "11.10.0.229",
    "url": "http://10.11.209.11/api/emanage_vms/2"
  }
]
*/

type EmanageVmsList []EmanageVms

func (evl *EmanageVmsList) GetVIP() (string, error) {
	for _, emanageVm := range *evl {
		if emanageVm.VIP != "" {
			return emanageVm.VIP, nil
		}
	}
	return "", fmt.Errorf("VIP not found")
}

func (evl *EmanageVmsList) GetPassiveEmanageIP() (string, error) {
	for _, emanageVm := range *evl {
		if !emanageVm.IsActive {
			return emanageVm.IP, nil
		}
	}
	return "", fmt.Errorf("Passive emanage IP was not found")
}

type emanageVms struct {
	conn *rest.Session
}

func (ev *emanageVms) Get() (EmanageVmsList, error) {
	var emsList EmanageVmsList
	err := ev.conn.Request(rest.MethodGet, emanageVMsUri, nil, &emsList)
	if err != nil {
		return emsList, err
	}
	return emsList, nil
}

func (evl *EmanageVmsList) ActiveEmanageIP() (string, error) {
	ips := helputils.StringSet{}
	for _, emanageVm := range *evl {
		ips.Add(emanageVm.IP)
		if emanageVm.IsActive {
			return emanageVm.IP, nil
		}
	}
	return "", errors.Errorf("No active emanaged found, ips: %v", ips)
}

func (ev *emanageVms) GetById(ID int) (*EmanageVms, error) {
	var ems EmanageVms
	fullUri := fmt.Sprintf("%s/%d", emanageVMsUri, ID)
	err := ev.conn.Request(rest.MethodGet, fullUri, nil, &ems)
	if err != nil {
		return nil, err
	}
	return &ems, nil
}

func (ev *emanageVms) Update(ems *EmanageVms) error {
	if ems == nil {
		return fmt.Errorf("Cannot update emanage VM with no data")
	}
	fullUri := fmt.Sprintf("%s/%d", emanageVMsUri, ems.ID)
	return ev.conn.Request(rest.MethodPut, fullUri, ems, nil)
}

func (ev *emanageVms) Create(ems *EmanageVms) (*EmanageVms, error) {
	if ems == nil {
		return nil, fmt.Errorf("Cannot create emanage VM with no data")
	}
	result := &EmanageVms{}
	err := ev.conn.Request(rest.MethodPost, emanageVMsUri, ems, result)
	return result, err
}

type EmanageMonitor struct {
	Active        bool   `json:"active"`
	Connected     bool   `json:"connected"`
	DataIP        string `json:"data_ip"`
	DataIP2       string `json:"data_ip2"`
	RemoteDataIP  string `json:"remote_data_ip"`
	RemoteDataIP2 string `json:"remote_data_ip2"`
}

/*
Active Emanage result:
{
  "active": true,
  "connected": true,
  "data_ip": "11.0.0.105",
  "data_ip2": "11.10.128.101",
  "remote_data_ip": "11.0.0.100",
  "remote_data_ip2": "11.10.128.99"
}
Passive Emanage result:
{
  "active": false,
  "connected": false,
  "data_ip": "11.0.0.100",
  "data_ip2": "11.10.128.99",
  "remote_data_ip": "11.0.0.105",
  "remote_data_ip2": "11.10.128.101"
}
*/
