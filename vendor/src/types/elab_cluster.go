package types

type ElabClusterData struct {
	Hosts    []ElabClusterHost `json:"hosts"`
	Networks struct {
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
		DataNetwork []struct {
			Name        string `json:"name"`
			NetworkID   string `json:"network_id"`
			NetworkMask int    `json:"network_mask"`
			VlanID      int    `json:"vlan_id"`
		} `json:"data_network"`
	} `json:"networks"`
	Size          int `json:"size"`
	SupportedRepl struct {
		Level2 bool `json:"level_2"`
		Level3 bool `json:"level_3"`
	} `json:"supported_repl"`
	SupportedRule struct {
		HCI bool `json:"HCI"`
		TOR bool `json:"TOR"`
	} `json:"supported_rule"`
}

type ElabClusterHost struct {
	DataNics  []string `json:"data_nics"`
	Datastore []struct {
		Accessible    bool   `json:"accessible"`
		BackingDevice string `json:"backing_device"`
		Capacity      int64  `json:"capacity"`
		FreeSpace     int64  `json:"free_space"`
		Local         bool   `json:"local,omitempty"`
		Name          string `json:"name"`
		Ssd           bool   `json:"ssd,omitempty"`
		Type          string `json:"type"`
	} `json:"datastore"`
	Hardware struct {
		CPU struct {
			Cores   int    `json:"cores"`
			Freq    int    `json:"freq"`
			Model   string `json:"model"`
			Sockets int    `json:"sockets"`
		} `json:"cpu"`
		RAM int64 `json:"ram"`
	} `json:"hardware"`
	IPAddress string `json:"ip_address"`
	Name      string `json:"name"`
	Nested    bool   `json:"nested"`
	Password  string `json:"password"`
	Storage   []struct {
		Capacity int64  `json:"capacity"`
		Name     string `json:"name"`
	} `json:"storage"`
	Type string `json:"type"`
	User string `json:"user"`
}
