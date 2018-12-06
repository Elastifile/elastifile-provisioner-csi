package emanage

import (
	"strconv"
	"strings"

	"helputils"
	"optional"
	"rest"
)

const clientsUri = "api/clients"

type clients struct {
	session *rest.Session
}

type EmanageClient struct {
	Id             int         `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Path           string      `json:"path,omitempty"`
	IsNas          bool        `json:"is_nas,omitempty"`
	IsDatastore    bool        `json:"is_datastore,omitempty"`
	NasConnected   string      `json:"nas_connected,omitempty"`
	PowerState     string      `json:"power_state,omitempty"`
	Status         string      `json:"status,omitempty"`
	Cores          int         `json:"cores,omitempty"`
	Memory         int         `json:"memory,omitempty"`
	Ip             string      `json:"ip,omitempty"`
	HostName       string      `json:"host_name,omitempty"`
	HostId         int         `json:"host_id,omitempty"`
	GuestOs        string      `json:"guest_os,omitempty"`
	MoId           interface{} `json:"mo_id,omitempty"`
	DeviceCapacity Bytes       `json:"device_capacity,omitempty"`
	DeviceUsage    Bytes       `json:"device_usage,omitempty"`
	ClientIp       string      `json:"client_ip,omitempty"`
	Url            string      `json:"url,omitempty"`
	Devices        interface{} `json:"devices"`
	Networks       []struct {
		Mac  string `json:"mac"`
		Ip   string `json:"ip"`
		Name string `json:"name"`
	} `json:"networks"`
}

func (cl *clients) GetAll(opt *GetAllOpts) ([]EmanageClient, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []EmanageClient
	err := cl.session.Request(rest.MethodGet, clientsUri, opt, &result)
	return result, err
}

func (cl *clients) GetById(id int) (EmanageClient, error) {
	var result EmanageClient
	fullUri := clientsUri + "/" + strconv.Itoa(id)
	err := cl.session.Request(rest.MethodGet, fullUri, nil, &result)
	return result, err
}

type EmanageClientsStats []struct {
	ClientIp string             `json:"client_ip"`
	Stats    EmanageClientStats `json:"stats"`
}

type EmanageClientStats struct {
	ID               int             `json:"id"`
	VmID             optional.Int    `json:"vm_id"`
	ClientMac        optional.String `json:"client_mac"`
	ReadNumEvents    optional.Int    `json:"read_num_events"`
	WriteNumEvents   optional.Int    `json:"write_num_events"`
	MdReadNumEvents  optional.Int    `json:"md_read_num_events"`
	MdWriteNumEvents optional.Int    `json:"md_write_num_events"`
	ReadIo           Bytes           `json:"read_io"`
	WriteIo          Bytes           `json:"write_io"`
	ReadLatency      Nanos           `json:"read_latency"`
	WriteLatency     Nanos           `json:"write_latency"`
	MdReadLatency    Nanos           `json:"md_read_latency"`
	MdWriteLatency   Nanos           `json:"md_write_latency"`
	Timestamp        string          `json:"timestamp"`
	EnodeID          optional.Int    `json:"enode_id"`
	EnodeName        optional.String `json:"enode_name"`
}

func (stats *EmanageClientStats) String() string {
	return "{" + strings.Join(helputils.StructToNamedStrings(*stats), " ") + "}"
}

func (cl *clients) GetStats() (EmanageClientsStats, error) {
	var result EmanageClientsStats
	fullUri := clientsUri + "/statistics"
	err := cl.session.Request(rest.MethodGet, fullUri, nil, &result)
	return result, err
}
