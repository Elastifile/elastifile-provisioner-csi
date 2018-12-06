package emanage

import (
	"path"
	"strings"

	"rest"
)

var (
	vmsUri     = "/api/vms"
	syncVMsUri = path.Join(vmsUri, "sync")
)

type VM struct {
	Name   string `json:"name"`
	Cores  int    `json:"cores"`
	HostID int    `json:"host_id"`

	Disks []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"devices"`

	Networks []struct {
		MAC  string `json:"mac"`
		IP   string `json:"ip,omitempty"`
		Name string `json:"name"`
	} `json:"networks"`
}

func (vm *VM) ExternalIpByName(name string) string {
	for _, n := range vm.Networks {
		if n.Name == name {
			return n.IP
		}
	}
	return ""
}

func (vm *VM) MacByName(name string) string {
	for _, net := range vm.Networks {
		if net.Name == name {
			return net.MAC
		}
	}
	return ""
}

func (vm *VM) DeviceIDsByPrefixes(prefixes ...string) []int {
	var result []int
	for _, disk := range vm.Disks {
		disk.Name = strings.TrimPrefix(disk.Name, "\"")
		for _, prefix := range prefixes {
			logger.Info("Setup Disks match", "disk name", disk.Name, "prefix", prefix)
			if strings.HasPrefix(disk.Name, prefix) {
				result = append(result, disk.Id)
			}
		}
	}
	return result
}

type vms struct {
	conn *rest.Session
}

func (vms *vms) Sync() error {
	if err := vms.conn.Request(rest.MethodPost, syncVMsUri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (vms *vms) GetAll() ([]VM, error) {
	var result []VM
	if err := vms.conn.Request(rest.MethodGet, vmsUri, nil, &result); err != nil {
		return result, err
	}
	return result, nil
}
