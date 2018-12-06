package emanage

import (
	"fmt"
	"path"
	"time"

	"github.com/go-errors/errors"
	"github.com/pborman/uuid"

	"eurl"
	"optional"
	"rest"
)

const dcUri = "api/data_containers"

type dataContainers struct {
	conn *rest.Session
}

type DataContainer struct {
	parent         dataContainers
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Uuid           uuid.UUID `json:"uuid"`
	Used           Bytes     `json:"used_capacity"`
	Scope          string    `json:"namespace_scope"`
	DataType       string    `json:"data_type"`
	Policy         Policy    `json:"policy"`
	PolicyId       uint      `json:"policy_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Url            eurl.URL  `json:"url"`
	Exports        []Export  `json:"exports"`
	ExportsCount   int       `json:"exports_count"`
	SoftQuota      Bytes     `json:"soft_quota"`
	HardQuota      Bytes     `json:"hard_quota"`
	DirPermissions int       `json:"dir_permissions,omitempty"`
	Dedup          int       `json:"dedup,omitempty"`
	Compression    int       `json:"compression,omitempty"`
	DirUID         int       `json:"dir_uid"`
	DirGid         int       `json:"dir_gid"`
	TenantID       int       `json:"tenant_id"`
	TotalUsedBytes Bytes `json:"total_used_bytes"`
}

type DcGetAllOpts struct {
	GetAllOpts

	FilterByVm     optional.Int `json:"filter_by_vm,omitempty"`     // Represents a vm_id
	FilterByPolicy optional.Int `json:"filter_by_policy,omitempty"` // Represents a policy_id
}

func populateDCConnections(dcList []DataContainer, parent *dataContainers) []DataContainer {
	for i, _ := range dcList {
		dcList[i].parent = *parent
	}
	return dcList
}

func (dcs *dataContainers) GetAll(opt *DcGetAllOpts) (result []DataContainer, err error) {
	if opt == nil {
		opt = &DcGetAllOpts{}
	}
	if err = dcs.conn.Request(rest.MethodGet, dcUri, opt, &result); err != nil {
		return nil, err
	}
	populateDCConnections(result, dcs)
	return result, nil
}

func (dcs *dataContainers) GetFull(dcId int) (result DataContainer, err error) {
	uri := fmt.Sprintf("%s/%d", dcUri, dcId)
	err = dcs.conn.Request(rest.MethodGet, uri, nil, &result)
	result.parent = *dcs
	return result, err
}

type DcCreateOpts struct {
	Name           string          `json:"name"`
	NamespaseScope string          `json:"namespace_scope,omitempty"`
	DataType       string          `json:"data_type,omitempty"`
	SoftQuota      int             `json:"soft_quota"`
	HardQuota      int             `json:"hard_quota"`
	DirPermissions int             `json:"dir_permissions,omitempty"`
	Share          optional.String `json:"share,omitempty"`
	VmIds          []int           `json:"vm_ids,omitempty"`
	Dedup          int             `json:"dedup,omitempty"`
	Compression    int             `json:"compression,omitempty"`
}

type DcUpdateOpts struct {
	Name        string          `json:"name,omitempty"`
	SoftQuota   int             `json:"soft_quota"`
	HardQuota   int             `json:"hard_quota"`
	Share       optional.String `json:"share,omitempty"`
	Dedup       int             `json:"dedup"`
	Compression int             `json:"compression"`
	PolicyId    int             `json:"policy_id,omitempty"`
}

type DcDirCreateOpts struct {
	Path        string       `json:"path"`
	Uid         optional.Int `json:"uid,omitempty"`
	Gid         optional.Int `json:"gid,omitempty"`
	Permissions optional.Int `json:"permissions,omitempty"`
}

func (dcs *dataContainers) Create(name string, policyId int, opt *DcCreateOpts) (DataContainer, error) {
	if opt == nil {
		opt = &DcCreateOpts{}
	}

	params := struct {
		Name     string `json:"name"`
		PolicyId int    `json:"policy_id"`
		DcCreateOpts
	}{name, policyId, *opt}
	var result DataContainer
	err := dcs.conn.Request(rest.MethodPost, dcUri, params, &result)
	result.parent = *dcs
	return result, err
}

func (dcs *dataContainers) DirCreate(dc *DataContainer, opt *DcDirCreateOpts) (*DataContainer, error) {
	if opt == nil {
		return nil, errors.New("Must specify directory opts")
	}

	logger.Info("Creating directory on data container", "path", opt.Path)

	var result DataContainer
	dcDirUri := fmt.Sprintf("%s/%d/create_dir", dcUri, dc.Id)
	err := dcs.conn.Request(rest.MethodPost, dcDirUri, opt, &result)
	result.parent = *dcs
	return &result, err
}

func (dcs *dataContainers) Update(dc *DataContainer, opt *DcUpdateOpts) (DataContainer, error) {
	if opt == nil {
		panic(fmt.Errorf("requireing update opts"))
	}

	params := struct {
		Name string `json:"name"`
		DcUpdateOpts
	}{dc.Name, *opt}
	var result DataContainer
	uri := fmt.Sprintf("%s/%d", dcUri, dc.Id)
	err := dcs.conn.Request(rest.MethodPut, uri, params, &result)
	result.parent = *dcs
	return result, err
}

func (dcs *dataContainers) Delete(dc *DataContainer) (result DataContainer, err error) {
	uri := path.Join(dcUri, fmt.Sprintf("%v", dc.Id))
	result = DataContainer{}
	if err = dcs.conn.Request(rest.MethodDelete, uri, nil, &result); err != nil {
		return result, err
	}
	result.parent = *dcs
	return result, nil
}

func (dcs *dataContainers) GetSnapshots(dc *DataContainer) (result []Snapshot, err error) {
	uri := path.Join(dcUri, fmt.Sprintf("%v", dc.Id), "snapshots")
	result = []Snapshot{}
	if err = dcs.conn.Request(rest.MethodGet, uri, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (dcs *dataContainers) GetExports(dc *DataContainer) (result []Export, err error) {
	uri := path.Join(dcUri, fmt.Sprintf("%v", dc.Id), "exports")
	result = []Export{}
	if err = dcs.conn.Request(rest.MethodGet, uri, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}
