package types

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	//"golang.org/x/tools/go/gcimporter15/testdata"
)

type ElabSystem struct {
	Data ElabData `json:"data"`
	Name string
}

type ElabData struct {
	Id                string                       `json:"id"`
	Name              string                       `json:"name"`
	Diagnostics       []string                     `json:"_diagnostics"`
	Emanage           []Node                       `json:"emanage"`
	EmanageVip        string                       `json:"emanage_vip"`
	Loaders           []LoaderNode                 `json:"loaders"`
	Vheads            []VheadNode                  `json:"vheads"`
	Nested            bool                         `json:"nested"`
	Networks          Networks                     `json:"networks"`
	Site              string                       `json:"site"`
	Type              systemType                   `json:"type"` //e.g. "DSM" or some int
	ReplicationAgents []ElabSystemReplicationAgent `json:"replication_agents,omitempty"`
}

const (
	_ = iota
	sysTypeHCI
	sysTypeDSM
	sysTypeVDSM
	sysTypeCloud
)

const (
	InternalNetworkId = iota
	ExternalNetworkId
)

const (
	InternalNetworkName = "internal"
	ExternalNetworkName = "external"
)

var (
	NetworkId2Name = map[int]string{
		InternalNetworkId: InternalNetworkName,
		ExternalNetworkId: ExternalNetworkName,
	}
)

type systemType string

func (s *systemType) UnmarshalJSON(data []byte) error {
	// Workaround for a inconsistent api on elab. Field 'Type' might return string or int.
	// the following will handle the case of incoming int and will convert it to string.
	sysTypeStrId := string(data)

	unquoted, err := strconv.Unquote(sysTypeStrId)
	if err == nil {
		sysTypeStrId = unquoted
	}

	typeId2Str := map[int]string{sysTypeHCI: "HCI", sysTypeDSM: "DSM", sysTypeVDSM: "vDSM", sysTypeCloud: "Cloud"}
	typeStr2Id := map[string]int{"HCI": sysTypeHCI, "DSM": sysTypeDSM, "vDSM": sysTypeVDSM, "Cloud": sysTypeCloud}

	// Get sysTypeStrId's numeric value
	typ, err := strconv.Atoi(sysTypeStrId)
	if err != nil { // TODO: Check where the code is that calls the function with its own return value
		// Handle the case where sysTypeStrId is a string
		var ok bool
		typ, ok = typeStr2Id[sysTypeStrId]
		if !ok {
			return err
		}
	}

	sysType, ok := typeId2Str[typ]
	if !ok {
		fmt.Printf("UnmarshalJSON - got unsupported value '%+v' - returning Unknown\n", typ)
		sysType = "Unknown"
	}

	*s = systemType(sysType)

	return nil
}

func (s *systemType) SetValue(prefix string, ctx interface{}) error {
	// Workaround for a inconsistent api on elab. Field 'Type' might return string or int.
	// It will set the systemType value.
	// This is called by both TagLoader and EnvLoader. this fix a failure of multiconfig Env loader which
	// tries to set the 'ElabData.Type' with string as 'systemType'.
	sysType := os.Getenv("TESLA_ELAB_SYSTEM_DATA_TYPE")
	*s = systemType(sysType)
	return nil
}

type NodeNetwork struct {
	Interface  string `json:"interface"`
	IpAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	Name       string `json:"name"`
	Prefix     int    `json:"prefix"`
	Version    string `json:"version"`
}

type Node struct {
	Host       string        `json:"host"`
	Hostname   string        `json:"hostname"`
	IpAddress  string        `json:"ip_address"`
	State      string        `json:"state"`
	VmName     string        `json:"vm_name"`
	MacAddress string        `json:"mac_address"`
	Networks   []NodeNetwork `json:"networks"`

	DataNics []string
}

func (n Node) Usable() bool {
	return n.State == "on"
}

// getNetworkByName returns node network by name
func (n Node) getNetworkByName(name string) NodeNetwork {
	for _, network := range n.Networks {
		if network.Name == name {
			return network
		}
	}
	return NodeNetwork{}
}

// getIpByNetworkName returns node's IP on network with the requested name
func (n Node) getIpByNetworkName(name string) string {
	network := n.getNetworkByName(name)
	return network.IpAddress
}

// getNetworkByType returns NodeNetwork by type
func (n Node) getNetworkByType(networkType int) (network NodeNetwork) {
	switch networkType {
	case InternalNetworkId:
		network = n.getNetworkByName(InternalNetworkName)
	case ExternalNetworkId:
		network = n.getNetworkByName(ExternalNetworkName)
	default:
		panic(fmt.Sprintf("Unsupported network type id: %v", networkType))
	}

	return
}

// ipByNetworkType returns node's IP on network with the requested type
func (n Node) ipByNetworkType(networkType int) (addr string) {
	switch networkType {
	case InternalNetworkId:
		addr = n.getIpByNetworkName(InternalNetworkName)
	case ExternalNetworkId:
		addr = n.getIpByNetworkName(ExternalNetworkName)
	default:
		panic(fmt.Sprintf("Unsupported network type id: %v", networkType))
	}

	if addr == "" {
		addr = n.IpAddress
	}

	return
}

func (n Node) InternalAddr() string {
	return n.ipByNetworkType(InternalNetworkId)
}

func (n Node) ExternalAddr() string {
	return n.ipByNetworkType(ExternalNetworkId)
}

func (n Node) InternalMac() (mac string) {
	mac = n.getNetworkByType(InternalNetworkId).MacAddress
	if mac == "" {
		mac = n.MacAddress
	}
	return
}

func (n Node) ExternalMac() (mac string) {
	mac = n.getNetworkByType(ExternalNetworkId).MacAddress
	if mac == "" {
		mac = n.MacAddress
	}
	return
}

type LoaderNode struct {
	Node
	ClnNetworkMac string `json:"cln_network_mac"`
}

type Networks struct {
	ClientExternalNetwork Network   `json:"client_external_network"`
	ClientInternalNetwork Network   `json:"client_internal_network"`
	DataNetwork           []Network `json:"data_network"`
}

type Network struct {
	Name         string   `json:"name"`
	NetworkId    string   `json:"network_id"`
	NetworkMask  int      `json:"network_mask"`
	VlanId       int      `json:"vlan_id"`
	VheadIpRange []string `json:"vheads_ip_range"`
}

type VheadNode struct {
	Node
	Type string `json:"type"`
}

var ElabNodeTypes []string = []string{
	"physical",
	"virtual",
}

func (vhn *VheadNode) IsPhysical() bool {
	return vhn.Type == "physical"
}

func (sys *ElabSystem) IsCloud() bool {
	return string(sys.Data.Type) == "Cloud"
}

func (sys *ElabSystem) EmanageVip() string {
	vip := sys.Data.EmanageVip
	if strings.HasSuffix(vip, ".%d") {
		vip = fmt.Sprintf(vip, sys.Data.Id)
	}
	return vip
}

func (sys *ElabSystem) SetHost(name string) error {
	if strings.HasPrefix(name, "func") {
		return fmt.Errorf("unsupported cluster type 'func'")
	} else {
		sys.Name = name
	}
	return nil
}

func (sys *ElabSystem) EmanageIp() string {
	if sys.Data.EmanageVip != "" {
		return sys.Data.EmanageVip
	}
	emsIps := sys.EmanageIps()
	if len(emsIps) > 0 {
		return emsIps[0]
	}
	return ""
}
func (sys *ElabSystem) UsableEmanages() (emngs []Node) {
	for _, emng := range sys.Data.Emanage {
		if emng.Usable() {
			emngs = append(emngs, emng)
		}
	}
	return emngs
}

//func (sys *ElabSystem) DefaultVip(net *Network) string {
//	return sys.Vip(net, 1)
//}
//
//func (sys *ElabSystem) Vip(net *Network, n int) string {
//	if sys.IsCloud() { // Eduard - be used LoadBalancer
//		if len(sys.VheadIps()) > 0 {
//			return sys.VheadIps()[0] //TODO - in case we're in system deploy - it's empty!!!!! Need some IP!!!!!
//		} else {
//			return ""
//		}
//
//	} else {
//		if net.NetworkId == "" {
//			return ""
//		} else if strings.HasSuffix(net.NetworkId, ".0") {
//			return strings.Join(strings.Split(net.NetworkId, ".")[:3], ".") + "." + strconv.Itoa(n)
//		}
//		return net.NetworkId
//	}
//}

func (sys *ElabSystem) EmanageIps() (ips []string) {
	for _, emanage := range sys.Data.Emanage {
		ips = append(ips, emanage.IpAddress)
	}
	return ips
}

func (sys *ElabSystem) EmanageMacs() (macs []string) {
	for _, emanage := range sys.UsableEmanages() {
		if sys.IsCloud() {
			macs = append(macs, emanage.ExternalMac())
		} else {
			macs = append(macs, emanage.MacAddress)
		}
	}
	return macs
}

func (sys *ElabSystem) EmanageVmNames() (vmNames []string) {
	for _, emanage := range sys.UsableEmanages() {
		vmNames = append(vmNames, emanage.VmName)
	}
	return vmNames
}

func (sys *ElabSystem) UsableLoaders() (lds []LoaderNode) {
	for _, loader := range sys.Data.Loaders {
		if loader.Usable() {
			lds = append(lds, loader)
		}
	}
	return lds
}

func (sys *ElabSystem) LoaderIps() (ips []string) {
	for _, loader := range sys.UsableLoaders() {
		ips = append(ips, loader.IpAddress)
	}
	return ips
}

func (sys *ElabSystem) LoaderIpsExternal() (ips []string) {
	for _, loader := range sys.UsableLoaders() {
		ips = append(ips, loader.ExternalAddr())
	}
	return ips
}

func (sys *ElabSystem) LoaderIpsInternal() (ips []string) {
	for _, loader := range sys.UsableLoaders() {
		ips = append(ips, loader.InternalAddr())
	}
	return ips
}

func (sys *ElabSystem) VheadIps() (ips []string) {
	for _, vhead := range sys.Data.Vheads {
		ips = append(ips, vhead.IpAddress)
	}
	return ips
}

func (sys *ElabSystem) VheadByIp(ip string) *VheadNode {
	var currentIp string
	for _, vhead := range sys.Data.Vheads {
		currentIp = vhead.IpAddress
		if currentIp == ip {
			return &vhead
		}
	}
	return nil
}

func (sys *ElabSystem) VheadByHostName(name string) *VheadNode {
	name = strings.Split(name, ".")[0]
	for _, vhead := range sys.Data.Vheads {
		if strings.Contains(vhead.Hostname, name) {
			return &vhead
		}
	}
	return nil
}
func (sys *ElabSystem) ActiveReplicationAgents() []ElabSystemReplicationAgent {
	var repAgents []ElabSystemReplicationAgent
	for _, repAgent := range sys.Data.ReplicationAgents {
		if strings.Contains(repAgent.State, "on") && net.ParseIP(repAgent.ExternalIPAddress) != nil {
			repAgents = append(repAgents, repAgent)
		}
	}
	return repAgents
}

func (sys *ElabSystem) IsDsm() bool {
	return sys.Data.Type == "DSM" || sys.Data.Type == "vDSM"
}

func (sys *ElabSystem) IsHci() bool {
	return sys.Data.Type == "HCI"
}

func (sys *ElabSystem) IsNested() bool {
	return sys.Data.Nested
}

func (sys *ElabSystem) IsFunctional() bool {
	return sys.IsDsm() && sys.IsNested()
}

func (sys *ElabSystem) Issues() (issues []string) {
	for _, issue := range sys.Data.Diagnostics {
		issues = append(issues, issue)
	}
	return issues
}

func (sys *ElabSystem) EmanageMacAddressMap() map[string]string {
	macAddrs := make(map[string]string)
	for _, emng := range sys.UsableEmanages() {
		if sys.IsCloud() {
			macAddrs[emng.IpAddress] = emng.Networks[ExternalNetworkId].MacAddress
		} else {
			macAddrs[emng.IpAddress] = emng.MacAddress
		}
	}
	return macAddrs
}

func (sys *ElabSystem) NfsNetwork() *Network {
	if sys.IsCloud() {
		// TODO: THIS IS A HACK!!! Provide proper solution
		// This function is N/A for cloud, but it gets called from quite a few places
		nets := &sys.Data.Networks
		return &nets.ClientInternalNetwork
	} else {
		nets := &sys.Data.Networks
		if sys.IsDsm() {
			return &nets.ClientExternalNetwork
		}
		return &nets.ClientInternalNetwork
	}
}

func (sys *ElabSystem) FileNamePrefix() string {
	return ConfElabSystemBucket
}

func (sys *ElabSystem) FileNameInList(i, total int) string {
	log10 := int(math.Log10(float64(total)))
	iFmt := fmt.Sprintf("%%0%dd", 1+log10)
	return fmt.Sprintf("%s-"+iFmt+"-%s", ConfElabSystemBucket, i, sys.Data.Id)
}

// elab deployment info - only actually used fields are defined!
// Since Tesla uses this struct primarily for OVA deploy and elab is using it for internal services,
// we keep the field definition to the minimum.
type ElabDeployData struct {
	ClusterId string
	Data      struct {
		Cluster struct {
			Hosts    []ElabDeployHostData  `json:"hosts"`
			Networks ElabDeployNetworkData `json:"networks"`
			/*
				Size     int                   `json:"size"`
				SupportedRepl struct {
					Level2 bool `json:"level_2"`
					Level3 bool `json:"level_3"`
				} `json:"supported_repl"`
				SupportedRule struct {
					HCI bool `json:"HCI"`
					TOR bool `json:"TOR"`
				} `json:"supported_rule"`
			*/
			ReplicationAgent []ElabDeployReplicationAgents `json:"replication_agent",omitempty`
		} `json:"cluster"`
		/*
			VCenter struct {
				Host     string `json:"host"`
				Password string `json:"password"`
				User     string `json:"user"`
			} `json:"vCenter"`
		*/
	} `json:"data"`
	/*
		Meta struct {
			APIVersion  int     `json:"api_version"`
			ElabVersion string  `json:"elab_version"`
			Error       bool    `json:"error"`
			ExecTime    float64 `json:"exec_time"`
			Group       string  `json:"group"`
			Resource    string  `json:"resource"`
			Timestamp   int     `json:"timestamp"`
		} `json:"meta"`
	*/
}

type ActiveDataNic struct {
	Name string `json:"name"`
	Vlan int    `json:"vlan"`
}

type ElabDeployHostData struct {
	ActiveDataNics []ActiveDataNic `json:"active_data_nics"`
	DataNics       []string        `json:"data_nics"`
	Datastore      []struct {
		Accessible    bool   `json:"accessible"`
		BackingDevice string `json:"backing_device"`
		Capacity      int64  `json:"capacity"`
		FreeSpace     int64  `json:"free_space"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Local         bool   `json:"local,omitempty"`
		Ssd           bool   `json:"ssd,omitempty"`
	} `json:"datastore"`
	Fqdn      string `json:"fqdn"`
	IPAddress string `json:"ip_address"`
	Name      string `json:"name"`
	/*	Hardware struct {
			CPU struct {
				Cores   int    `json:"cores"`
				Freq    int    `json:"freq"`
				Model   string `json:"model"`
				Sockets int    `json:"sockets"`
			} `json:"cpu"`
			RAM int64 `json:"ram"`
		} `json:"hardware"`
	*/
	/*
		Nested    bool   `json:"nested"`
		Password  string `json:"password"`
		Storage []struct {
			Capacity int64  `json:"capacity"`
			Name     string `json:"name"`
		} `json:"storage"`
		Type int    `json:"type"`
		User string `json:"user"`
	*/
}

type ElabDeployNetworkData struct {
	/*
		ClientExternalNetwork struct {
			Name          string   `json:"name"`
			NetworkID     string   `json:"network_id"`
			NetworkMask   int      `json:"network_mask"`
			VheadsIPRange []string `json:"vheads_ip_range"`
			VlanID        int      `json:"vlan_id"`
		} `json:"client_external_network"`
		ClientInternalNetwork struct {
			Name        string `json:"name"`
			NetworkID   string `json:"network_id"`
			NetworkMask int    `json:"network_mask"`
			VlanID      int    `json:"vlan_id"`
		} `json:"client_internal_network"`
	*/
	DataNetwork []struct {
		Name        string `json:"name"`
		NetworkID   string `json:"network_id"`
		NetworkMask int    `json:"network_mask"`
		VlanID      int    `json:"vlan_id"`
	} `json:"data_network"`
}

type ElabSystemReplicationAgent struct {
	ClientIPAddress    string `json:"client_ip_address,omitempty"`
	ClientNetworkMac   string `json:"client_network_mac,omitempty"`
	ExternalIPAddress  string `json:"external_ip_address"`
	ExternalNetworkMac string `json:"external_network_mac"`
	Host               string `json:"host"`
	Hostname           string `json:"hostname"`
	Role               int    `json:"role"`
	State              string `json:"state"`
	Type               string `json:"type"`
	VMName             string `json:"vm_name"`
}

type ElabDeployReplicationAgents struct {
	Client struct {
		Gateway         string `json:"gateway"`
		IPAddress       string `json:"ip_address"`
		Mtu             int    `json:"mtu"`
		NetworkMask     string `json:"network_mask"`
		NetworkMaskBits int    `json:"network_mask_bits"`
		VlanID          int    `json:"vlan_id"`
	} `json:"client"`
	External struct {
		Gateway         string `json:"gateway"`
		IPAddress       string `json:"ip_address"`
		Mtu             int    `json:"mtu"`
		NetworkMask     string `json:"network_mask"`
		NetworkMaskBits int    `json:"network_mask_bits"`
		VlanID          int    `json:"vlan_id"`
	} `json:"external"`
	VMName string `json:"vm_name"`
}

type ProvisionData struct {
	Type             int    `json:"type" "default:"4"`
	ElfsVersion      string `json:"elfs_version"`
	EloaderVersion   string `json:"eloader_version,omitempty"`
	EmsInstances     int    `json:"ems_instances" "default:"1"`
	EloaderInstances int    `json:"eloader_instances" "default:"1"`
}

type ProvisionInfo struct {
	Operation ProvisionOperation `json:"data"`
	Meta      struct {
		APIVersion  int     `json:"api_version"`
		ElabVersion string  `json:"elab_version"`
		Error       bool    `json:"error"`
		ExecTime    float64 `json:"exec_time"`
		Group       string  `json:"group"`
		Instance    string  `json:"instance"`
		Resource    string  `json:"resource"`
		Status      int     `json:"status"`
		Timestamp   int     `json:"timestamp"`
	} `json:"meta"`
}

type ProvisionMessages struct {
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type ProvisionOperation struct {
	CreatedAt   int                 `json:"created_at"`
	Description string              `json:"description"`
	Error       bool                `json:"error"`
	Errors      []string            `json:"errors"`
	ID          string              `json:"id"`
	Kind        string              `json:"kind"`
	Link        string              `json:"link"`
	Messages    []ProvisionMessages `json:"messages"`
	Name        string              `json:"name"`
	Progress    int                 `json:"progress"`
	Status      string              `json:"status"`
	UpdatedAt   int                 `json:"updated_at"`
	User        string              `json:"user"`
	Warning     bool                `json:"warning"`
	Warnings    []string            `json:"warnings"`
}

func (info *ProvisionInfo) IsDone() bool {
	return info.Operation.Status == "DONE"
}

type ClusterInfo struct {
	Data struct {
		Emanage []struct {
			Active     bool   `json:"active"`
			Host       string `json:"host"`
			Hostname   string `json:"hostname"`
			IPAddress  string `json:"ip_address"`
			MacAddress string `json:"mac_address"`
			Networks   []struct {
				Interface string `json:"interface"`
				IPAddress string `json:"ip_address"`
				Name      string `json:"name"`
			} `json:"networks"`
			Role   string `json:"role"`
			State  string `json:"state"`
			Type   string `json:"type"`
			VMName string `json:"vm_name"`
		} `json:"emanage"`
		EmanageVip string   `json:"emanage_vip"`
		Functional bool     `json:"functional"`
		Hosts      []string `json:"hosts"`
		Hypervisor int      `json:"hypervisor"`
		ID         string   `json:"id"`
		Loaders    []struct {
			Host       string `json:"host"`
			Hostname   string `json:"hostname"`
			IPAddress  string `json:"ip_address"`
			MacAddress string `json:"mac_address"`
			Networks   []struct {
				Interface string `json:"interface"`
				IPAddress string `json:"ip_address"`
				Name      string `json:"name"`
			} `json:"networks"`
			Role   string `json:"role"`
			State  string `json:"state"`
			Type   string `json:"type"`
			VMName string `json:"vm_name"`
		} `json:"loaders"`
		Name     string `json:"name"`
		Nested   bool   `json:"nested"`
		Networks struct {
		} `json:"networks"`
		ReplicationAgents []string `json:"replication_agents"`
		Site              string   `json:"site"`
		Size              int      `json:"size"`
		SupportedRepl     struct {
			Level2 bool `json:"level_2"`
			Level3 bool `json:"level_3"`
		} `json:"supported_repl"`
		SupportedRule struct {
			HCI bool `json:"HCI"`
			TOR bool `json:"TOR"`
		} `json:"supported_rule"`
		Type    int `json:"type"`
		VCenter struct {
			Host     string `json:"host"`
			Password string `json:"password"`
			User     string `json:"user"`
		} `json:"vCenter"`
		Vheads []string `json:"vheads"`
	} `json:"data"`
	Meta struct {
		APIVersion  int     `json:"api_version"`
		ElabVersion string  `json:"elab_version"`
		Error       bool    `json:"error"`
		ExecTime    float64 `json:"exec_time"`
		Group       string  `json:"group"`
		Instance    string  `json:"instance"`
		Resource    string  `json:"resource"`
		Status      int     `json:"status"`
		Timestamp   int     `json:"timestamp"`
	} `json:"meta"`
}
