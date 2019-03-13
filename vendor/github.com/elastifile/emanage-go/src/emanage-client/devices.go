package emanage

import (
	"fmt"

	"rest"
)

const (
	devicesUri = "/api/devices"
)

type devices struct {
	conn *rest.Session
}

type DeviceStatus string

const (
	StatusOk     DeviceStatus = "OK"
	StatusFailed DeviceStatus = "Failed"
)

type EmanageDevices struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Vendor       string       `json:"vendor"`
	IsWritable   bool         `json:"is_writable"`
	Ssd          bool         `json:"ssd"`
	DevicePath   string       `json:"device_path"`
	Status       DeviceStatus `json:"status"`
	HostID       interface{}  `json:"host_id"`
	EnodeID      int          `json:"enode_id"`
	HostName     string       `json:"host_name"`
	Capacity     Bytes        `json:"capacity"`
	Usage        Bytes        `json:"usage"`
	Datastore    interface{}  `json:"datastore"`
	IsLocal      bool         `json:"is_local"`
	IsInternalFs bool         `json:"is_internal_fs"`
	IsLogFs      bool         `json:"is_log_fs"`
}

type EmanageDevicesList []*EmanageDevices

func (devs *devices) GetAll() (EmanageDevicesList, error) {
	var devList EmanageDevicesList
	err := devs.conn.Request(rest.MethodGet, devicesUri, nil, &devList)
	return devList, err
}

func (devs *devices) GetDeviceById(id int) (*EmanageDevices, error) {
	var device EmanageDevices
	fullURI := fmt.Sprintf("%s/%d", devicesUri, id)
	err := devs.conn.Request(rest.MethodGet, fullURI, nil, &device)

	return &device, err
}
