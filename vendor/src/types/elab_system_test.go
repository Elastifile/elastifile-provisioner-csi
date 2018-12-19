package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestElabSysJsonUnmarshal(t *testing.T) {
	jsonContent := `{
		"data": {
		  "_diagnostics": [],
		  "diagnostics": [],
		  "emanage": [
			{
			  "host": "eManage-1",
			  "hostname": "eManage-1",
			  "ip_address": "10.11.209.214",
			  "mac_address": "52:54:00:94:cc:e4",
			  "role": null,
			  "state": "on",
			  "type": "virtual",
			  "vm_name": "eManage-1"
			}
		  ],
		  "emanage_vip": "10.11.209.214",
		  "id": "214",
		  "loaders": [
			{
			  "host": "Loader-1",
			  "hostname": "Loader-1",
			  "ip_address": "10.11.197.216",
			  "role": null,
			  "state": "on",
			  "type": "virtual",
			  "vm_name": "Loader-1"
			}
		  ],
		  "name": "C214",
		  "nested": false,
		  "networks": {
			"client_external_network": {
			  "name": "Elastifile-EXT-CLN-sw",
			  "network_id": "172.16.0.0",
			  "network_mask": 16,
			  "vheads_ip_range": [
				"172.16.214.1",
				"172.16.214.2",
				"172.16.214.3",
				"172.16.214.4"
			  ],
			  "vlan_id": 2214
			},
			"client_internal_network": {
			  "name": "Elastifile-CLN-sw",
			  "network_id": "192.168.214.0",
			  "network_mask": 24,
			  "vlan_id": 0
			},
			"data_network": [
			  {
				"name": "Elastifile-DATA-A-sw",
				"network_id": "10.214.10.0",
				"network_mask": 24,
				"vlan_id": 214
			  },
			  {
				"name": "Elastifile-DATA-B-sw",
				"network_id": "10.214.11.0",
				"network_mask": 24,
				"vlan_id": 1214
			  }
			]
		  },
		  "site": "il-lab",
		  "size": 1,
		  "supported_repl": {
			"level_2": false,
			"level_3": false
		  },
		  "supported_rule": {
			"HCI": false,
			"TOR": false
		  },
		  "type": 2,
		  "vCenter": {
			"host": "vc214.lab.il.elastifile.com",
			"password": "vmware",
			"user": "root"
		  },
		  "vheads": [
			{
			  "host": "vhead-1",
			  "hostname": "vhead-1",
			  "ip_address": "10.11.198.222",
			  "role": null,
			  "state": "on",
			  "type": "virtual",
			  "vm_name": "vhead-1"
			}
		  ]
		},
		"meta": {
		  "api_version": 1,
		  "elab_version": "1.0.154 (250fe7277b9b)",
		  "error": false,
		  "exec_time": 117.923,
		  "group": "system",
		  "resource": "cluster",
		  "timestamp": 1504603246
		}
	  }`
	var elabSys ElabSystem

	if err := json.Unmarshal([]byte(jsonContent), &elabSys); err != nil {
		t.Fatalf("error returned %+v", err)
	}

	fmt.Printf("%+v\n", elabSys)
	fmt.Printf("Type: %+v\n", elabSys.Data.Type)
}
