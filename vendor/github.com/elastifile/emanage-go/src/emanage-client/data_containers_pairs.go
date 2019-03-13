package emanage

import (
	"fmt"
	"time"

	"rest"
)

type dcPairs struct {
	conn *rest.Session
}

type Pair struct {
	ID                   int    `json:"id,omitempty"`
	UUID                 string `json:"uuid,omitempty"`
	DataContainerID      int    `json:"data_container_id,omitempty"`
	RemoteSiteID         int    `json:"remote_site_id,omitempty"`
	RemoteSystemName     string `json:"remote_system_name,omitempty"`
	RemoteDcID           int    `json:"remote_dc_id,omitempty"`
	RemoteDcUUID         string `json:"remote_dc_uuid,omitempty"`
	RemoteDcPairID       int    `json:"remote_dc_pair_id,omitempty"`
	RemoteDcPairUUID     string `json:"remote_dc_pair_uuid,omitempty"`
	DrRole               string `json:"dr_role,omitempty"`
	Rpo                  int    `json:"rpo,omitempty"`
	ConnectionStatus     string `json:"connection_status,omitempty"`
	Reason               string `json:"reason,omitempty"`
	ReplicateAcls        bool   `json:"replicate_acls,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
	LastStartTime        string `json:"last_start_time,omitempty"`
	LastEndTime          string `json:"last_end_time,omitempty"`
	LastDrStatus         string `json:"last_dr_status,omitempty"`
	ReplicationGatewayIP string `json:"replication_gateway_ip,omitempty"`
	ActiveSyncOp         bool   `json:"active_sync_op,omitempty"`
	OverrideScheduler    bool   `json:"override_scheduler,omitempty"`
	URL                  string `json:"url,omitempty"`
}

func (dc *dcPairs) GetPairs(dcId int) (result []Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs", dcUri, dcId)
	err = dc.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

func (dc *dcPairs) GetById(dcId int, pairId int) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%d", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

func (dc *dcPairs) GetBySiteId(dcId int, siteId int) (result []Pair, err error) {
	pairs, err := dc.GetPairs(dcId)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		if pair.RemoteSiteID == siteId {
			result = append(result, pair)
		}
	}

	return result, err
}

type DrStatus string

const (
	Sync              DrStatus = "sync"
	ExceededRpo       DrStatus = "exceeded_rpo"
	ReplicationFailed DrStatus = "replication_failed"
	Promoting         DrStatus = "promoting"
	Collision         DrStatus = "collision"
)

type DrRoleOpts string

const (
	DrRoleNone    DrRoleOpts = "role_dc_none"
	DrRoleActive  DrRoleOpts = "role_dc_active"
	DrRolePassive DrRoleOpts = "role_dc_passive"
)

type PairCreateOpts struct {
	RemoteSiteId      string     `json:"remote_site_id,omitempty"`
	Rpo               int        `json:"rpo,omitempty"`
	DrRole            DrRoleOpts `json:"dr_role"`
	OverrideScheduler bool       `json:"override_scheduler,omitempty"`
	ReplicateAcls     bool       `json:"replicate_acls,omitempty"`
}

func (dc *dcPairs) CreatePair(dcId int, pairCreateOpts PairCreateOpts) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs", dcUri, dcId)
	err = dc.conn.Request(rest.MethodPost, uri, pairCreateOpts, &result)
	return result, err
}

func (dc *dcPairs) Update(dcId int, pairId int, pairCreateOpts PairCreateOpts) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodPut, uri, pairCreateOpts, &result)
	return result, err
}

func (dc *dcPairs) Connect(dcId int, pairId int) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/connect", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodPost, uri, nil, &result)
	return result, err
}

type PairsDisconnectOpts struct {
	Force bool `json:"force" default:"false"`
}

func (dc *dcPairs) Disconnect(dcId int, pairId int, pairsDisconnectOpts PairsDisconnectOpts) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/disconnect", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodPost, uri, pairsDisconnectOpts, &result)
	return result, err
}

type ReplicationLogs struct {
	DcPairID  int       `json:"dc_pair_id"`
	DrStatus  string    `json:"dr_status"`
	DrRole    string    `json:"dr_role"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (dc *dcPairs) ReplicationLogs(dcId int, pairId int) (result []ReplicationLogs, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/replication_logs", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

func (dc *dcPairs) Delete(dcId int, pairId int) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodDelete, uri, nil, &result)
	return result, err
}

type TestDataContainereOpts struct {
	DataContainerName string `json:"data_container_name"`
	SnapshotId        int    `json:"snapshot_id,omitempty"`
	Async             bool   `json:"async,omitempty"`
	Dedup             int    `json:"dedup,omitempty"`
	Compression       int    `json:"compression,omitempty"`
}

type AsyncDataContainere []struct {
	ID                    int         `json:"id"`
	UUID                  string      `json:"uuid"`
	LastError             interface{} `json:"last_error"`
	Priority              int         `json:"priority"`
	Attempts              int         `json:"attempts"`
	Queue                 interface{} `json:"queue"`
	Status                string      `json:"status"`
	Name                  string      `json:"name"`
	CurrentStep           interface{} `json:"current_step"`
	StepProgress          interface{} `json:"step_progress"`
	StepTotal             interface{} `json:"step_total"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	Host                  string      `json:"host"`
	TaskType              string      `json:"task_type"`
	TargetDataContainerID int         `json:"target_data_container_id"`
	DcPairID              int         `json:"dc_pair_id"`
	URL                   string      `json:"url"`
}

func (dc *dcPairs) TestImage(dcId int, pairId int, testDataContainereOpts TestDataContainereOpts) (result AsyncDataContainere, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/test_image", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodPost, uri, testDataContainereOpts, &result)
	return result, err
}

type ForcePromoteOpts struct {
	SnapshotId int  `json:"snapshot_id,omitempty"`
	Async      bool `json:"async,omitempty"`
}

func (dc *dcPairs) ForcePromote(dcId int, pairId int, forcePromoteOpts ForcePromoteOpts) (result Pair, err error) {
	uri := fmt.Sprintf("%s/%d/dc_pairs/%v/force_promote", dcUri, dcId, pairId)
	err = dc.conn.Request(rest.MethodPost, uri, forcePromoteOpts, &result)
	return result, err
}
