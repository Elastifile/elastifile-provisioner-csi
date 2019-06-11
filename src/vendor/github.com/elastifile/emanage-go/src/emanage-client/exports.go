package emanage

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"

	"eurl"
	"optional"
	"rest"
)

const exportsUri = "api/exports"

type exports struct {
	session *rest.Session
}

type UserMappingType string

const (
	UserMappingNone UserMappingType = "no_mapping"
	UserMappingRoot UserMappingType = "remap_root"
	UserMappingAll  UserMappingType = "remap_all"
)

var UserMappingValues = []UserMappingType{
	UserMappingNone,
	UserMappingRoot,
	UserMappingAll,
}

type ExportStateType string

const (
	ExportStateEnabled ExportStateType = "export_state_enabled"
	ExportStateDisabled ExportStateType = "export_state_disabled"
)

type ExportCreateForVolumeOpts struct {
	DcId        int                 `json:"data_container_id"`
	Path        string               `json:"path"`
	Uuid        optional.String      `json:"uuid,omitempty"` // Should type be uuid.UUID? How do we handle nilable?
	Access      ExportAccessModeType `json:"access_permission,omitempty"`
	UserMapping UserMappingType      `json:"user_mapping,omitempty"`
	Uid         int                  `json:"uid,omitempty"`
	Gid         int                  `json:"gid,omitempty"`
}

type ExportCreateForSnapshotOpts struct {
	SnapShotId  int                 `json:"snapshot_id"`
	Path        string               `json:"path"`
	Uuid        optional.String      `json:"uuid,omitempty"` // Should type be uuid.UUID? How do we handle nilable?
	Access      ExportAccessModeType `json:"access_permission,omitempty"`
	UserMapping UserMappingType      `json:"user_mapping,omitempty"`
	Uid         int                  `json:"uid,omitempty"`
	Gid         int                  `json:"gid,omitempty"`
}

type Export struct {
	Id              int                  `json:"id"`
	Name            string               `json:"name"`
	Path            string               `json:"path"` // TODO: should have a type?
	Uuid            uuid.UUID            `json:"uuid"`
	Access          ExportAccessModeType `json:"access_permission"`
	UserMapping     UserMappingType      `json:"user_mapping"`
	DataContainerId int                  `json:"data_container_id"`
	SnapshotId      int                  `json:"snapshot_id"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Url             eurl.URL             `json:"url"`
	Uid             int                  `json:"uid"`
	Gid             int                  `json:"uid"`
	ClientRules     []ExportRule         `json:"client_rules"`
	SnapshotName    string               `json:"-"`
	ExportState		ExportStateType		 `json:"export_state"`
}

func (ex *exports) GetAll(opt *GetAllOpts) ([]Export, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []Export
	if err := ex.session.Request(rest.MethodGet, exportsUri, opt, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ex *exports) GetFull(exportId int) (Export, error) {
	uri := fmt.Sprintf("%s/%d", exportsUri, exportId)
	var result Export
	err := ex.session.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

type ExportUpdateOpts struct {
	Path        string               `json:"path"`
	Uuid        optional.String      `json:"uuid,omitempty"` // Should type be uuid.UUID? How do we handle nilable?
	Access      ExportAccessModeType `json:"access_permission,omitempty"`
	UserMapping UserMappingType      `json:"user_mapping,omitempty"`
	Uid         optional.Int         `json:"uid,omitempty"` //changed to optional to allow o value pass
	Gid         optional.Int         `json:"gid,omitempty"` //changed to optional to allow o value pass
}

type ExportCreateOpts struct {
	DcId         int                  `json:"data_container_id,omitempty"`
	SnapshotId   int                  `json:"snapshot_id,omitempty"`
	Path         string               `json:"path"`
	Uuid         optional.String      `json:"uuid,omitempty"` // Should type be uuid.UUID? How do we handle nilable?
	Access       ExportAccessModeType `json:"access_permission,omitempty"`
	UserMapping  UserMappingType      `json:"user_mapping,omitempty"`
	Uid          optional.Int         `json:"uid,omitempty"` //changed to optional to allow o value pass
	Gid          optional.Int         `json:"gid,omitempty"` //changed to optional to allow o value pass
	SnapshotName string               `json:"-"`
}

type ExportAccessModeType string

const (
	ExportAccessRW   ExportAccessModeType = "read_write"
	ExportAccessRO   ExportAccessModeType = "read_only"
	ExportAccessList ExportAccessModeType = "list_only"
	ExportAccessNone ExportAccessModeType = "no_access"
)

var ExportAccessModeValues = []ExportAccessModeType{
	ExportAccessRW,
	ExportAccessRO,
	ExportAccessList,
	ExportAccessNone,
}

func (ex *exports) CreateForVolume(name string, opt *ExportCreateForVolumeOpts) (Export, error) {
	if opt == nil {
		opt = &ExportCreateForVolumeOpts{}
	}
	params := struct {
		Name string `json:"name"`
		ExportCreateForVolumeOpts
	}{name, *opt}

	var result Export
	err := ex.session.Request(rest.MethodPost, exportsUri, params, &result)
	return result, err
}

//TODO: need refactor
func (ex *exports) CreateForSnapshot(name string, opt *ExportCreateForSnapshotOpts) (Export, error) {
	if opt == nil {
		opt = &ExportCreateForSnapshotOpts{}
	}
	params := struct {
		Name string `json:"name"`
		ExportCreateForSnapshotOpts
	}{name, *opt}

	var result Export
	err := ex.session.Request(rest.MethodPost, exportsUri, params, &result)
	return result, err
}

func (ex *exports) Update(export *Export, opt *ExportUpdateOpts) (Export, error) {
	if opt == nil {
		panic(fmt.Errorf("missing export %s update options", export.Name))
	}

	if opt.Path == "" {
		opt.Path = export.Path
		logger.Warn("copied path from export to options, due to emanage requirement", "Path", opt.Path)
	}

	params := struct {
		Name string `json:"name"`
		ExportUpdateOpts
	}{export.Name, *opt}

	uri := fmt.Sprintf("%s/%d", exportsUri, export.Id)
	var result Export
	err := ex.session.Request(rest.MethodPut, uri, params, &result)
	return result, err
}

func (ex *exports) Delete(export *Export) (Export, error) {
	uri := fmt.Sprintf("%s/%d", exportsUri, export.Id)
	var result Export
	err := ex.session.Request(rest.MethodDelete, uri, nil, &result)
	return result, err
}
