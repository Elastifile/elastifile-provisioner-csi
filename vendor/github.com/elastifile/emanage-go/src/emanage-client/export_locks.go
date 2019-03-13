package emanage

import (
	"fmt"
	"time"

	"rest"
)

type DirLockedFilesReply struct {
	Name           string      `json:"name"`
	Type           int         `json:"type"`
	UID            int         `json:"uid"`
	Gid            int         `json:"gid"`
	Size           int         `json:"size"`
	Mode           string      `json:"mode"`
	IsDirectory    bool        `json:"isDirectory"`
	FullPath       string      `json:"fullPath"`
	ParentFullPath string      `json:"parentFullPath"`
	Cursor         interface{} `json:"cursor"`
}

type GetFileLocksReply struct {
	Type        string `json:"type"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
	Address     string `json:"address"`
	Access      int    `json:"access"`
	Mode        int    `json:"mode"`
}

type BreakLockReply struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Path             string `json:"path"`
	UUID             string `json:"uuid"`
	AccessPermission string `json:"access_permission"`
	UserMapping      string `json:"user_mapping"`
	UID              int    `json:"uid"`
	Gid              int    `json:"gid"`
	DataContainer    DataContainer `json:"data_container"`
	DataContainerID  int           `json:"data_container_id"`
	SnapshotID       interface{}   `json:"snapshot_id"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	ClientRulesCount int           `json:"client_rules_count"`
	NamespaceScope   string        `json:"namespace_scope"`
	DataType         string        `json:"data_type"`
	ClientRules      []interface{} `json:"client_rules"`
}

type LockArgs struct {
	Path string `json:"path"`
}

func (ex *exports) getFullURI(export *Export, op string) string {
	return fmt.Sprintf("%s/%d/%s", exportsUri, export.Id, op)
}

func (ex *exports) uriGetDirLockedFiles(export *Export) string {
	return ex.getFullURI(export, "get_dir_locked_files")
}

func (ex *exports) uriGetFileLocks(export *Export) string {
	return ex.getFullURI(export, "get_file_locks")
}

func (ex *exports) uriBreakLock(export *Export) string {
	return ex.getFullURI(export, "break_lock")
}

func (ex *exports) GetDirLockedFiles(export *Export, lockArgs *LockArgs) (result []DirLockedFilesReply, err error) {
	err = ex.session.Request(rest.MethodGet, ex.uriGetDirLockedFiles(export), lockArgs, &result)
	return result, err
}

func (ex *exports) GetFileLocks(export *Export, lockArgs *LockArgs) (result []GetFileLocksReply, err error) {
	err = ex.session.Request(rest.MethodGet, ex.uriGetFileLocks(export), lockArgs, &result)
	return result, err
}

func (ex *exports) BreakLock(export *Export, lockArgs *LockArgs) (result BreakLockReply, err error) {
	err = ex.session.Request(rest.MethodPost, ex.uriBreakLock(export), lockArgs, &result)
	return result, err
}
