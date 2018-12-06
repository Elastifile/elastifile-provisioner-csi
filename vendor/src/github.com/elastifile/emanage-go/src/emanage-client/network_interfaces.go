package emanage

import (
	"fmt"
	"path"

	"rest"
)

var (
	netInterfacesUri = "/api/network_interfaces"
)

type netInterfaces struct {
	conn *rest.Session
}

func (n *netInterfaces) GetAll() ([]NetworkInterface, error) {
	var result []NetworkInterface
	if err := n.conn.Request(rest.MethodGet, netInterfacesUri, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (n *netInterfaces) GetById(id int) (*NetworkInterface, error) {
	var result NetworkInterface
	fullUri := path.Join(netInterfacesUri, fmt.Sprintf("%v", id))
	if err := n.conn.Request(rest.MethodGet, fullUri, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (n *netInterfaces) GetByHost(hostId int) ([]NetworkInterface, error) {
	if nifs, err := n.GetAll(); err != nil {
		return nifs, err
	} else {
		hostNifs := make([]NetworkInterface, 0)
		for _, nif := range nifs {
			if nif.HostID == hostId {
				hostNifs = append(hostNifs)
			}
		}
		return hostNifs, nil
	}
}

type UpdateNetInterfacesOpts struct {
	Role          string `json:"role"`
	ClientNetwork bool   `json:"client_network"`
}

// Updates network_interfaces matching given id with requested UpdateNetInterfacesOpts
func (n *netInterfaces) Update(id int, opts *UpdateNetInterfacesOpts) error {
	iface, err := n.GetById(id)
	if err != nil {
		return err
	}

	type netIfaceOpts struct {
		NetworkInterface *NetworkInterface `json:"network_interface"`
	}

	iface.Role = opts.Role
	iface.ClientNetwork = opts.ClientNetwork
	fullUri := path.Join(netInterfacesUri, fmt.Sprintf("%v", id))
	return n.conn.Request(rest.MethodPut, fullUri, netIfaceOpts{NetworkInterface: iface}, nil)
}

type NetworkInterface struct {
	Name           string `json:"name"`
	ID             int    `json:"id"`
	DetectedByNic1 string `json:"detected_by_nic1"`
	DetectedByNic2 string `json:"detected_by_nic2"`
	Dhcp           bool   `json:"dhcp"`
	EnodeID        int    `json:"enode_id"`
	HostID         int    `json:"host_id"`
	Lladdress      string `json:"lladdress"`
	MaxVfunc       int    `json:"max_vfunc"`
	NumVfunc       int    `json:"num_vfunc"`
	Role           string `json:"role"`
	ClientNetwork  bool   `json:"client_network"`
	Speed          int    `json:"speed"`
	SriovActive    bool   `json:"sriov_active"`
	SriovCapable   bool   `json:"sriov_capable"`
	SriovEnabled   bool   `json:"sriov_enabled"`
	Status         string `json:"status"`
	Subnet         string `json:"subnet"`
	URL            string `json:"url"`
	VfuncID        int    `json:"vfunc_id"`
	VswitchID      int    `json:"vswitch_id"`
	VswitchName    string `json:"vswitch_name"`
	Networks       []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Vlan int    `json:"vlan"`
	} `json:"networks"`
}
