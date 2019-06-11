package emanage

import (
	"fmt"
	"strings"
	"time"

	"github.com/koding/multiconfig"

	"rest"
)

const remoteSitesUri = "/api/remote_sites"

type remoteSites struct {
	session *rest.Session
}

type RemoteSite struct {
	ID                int       `json:"id"`
	UUID              string    `json:"uuid"`
	IPAddress         string    `json:"ip_address"`
	Login             string    `json:"login"`
	Password          string    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	SystemID          int       `json:"system_id"`
	RemoteSystemID    int       `json:"remote_system_id"`
	RemoteSystemUUID  string    `json:"remote_system_uuid"`
	RemoteSystemName  string    `json:"remote_system_name"`
	ConnectionStatus  string    `json:"connection_status"`
	Reason            string    `json:"reason"`
	LocalIPAddress    string    `json:"local_ip_address"`
	ReplicatedDcCount int       `json:"replicated_dc_count"`
	HostedDcCount     int       `json:"hosted_dc_count"`
	URL               string    `json:"url"`
}

func (rs *RemoteSite) Connected() bool {
	return strings.Contains(rs.ConnectionStatus, "connected")
}

func (rs *remoteSites) GetAll() (result []RemoteSite, err error) {
	err = rs.session.Request(rest.MethodGet, remoteSitesUri, nil, &result)
	return result, err
}

func (rs *remoteSites) GetById(siteId int) (result RemoteSite, err error) {
	uri := fmt.Sprintf("%s/%d", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

type RemoteSiteOpts struct {
	RemoteSystemName string `json:"remote_system_name,omitempty"`
	IpAddress        string `json:"ip_address"`
	Login            string `json:"login,omitempty" default:"admin"`
	Password         string `json:"password,omitempty" default:"changeme"`
	SystemId         int    `json:"system_id,omitempty"`
	LocalLogin       string `json:"local_login,omitempty" default:"admin"`
	LocalPassword    string `json:"local_password,omitempty" default:"changeme"`
	LocalIpAddress   string `json:"local_ip_address,omitempty"`
}

func NewRemoteSiteOpts() (opts RemoteSiteOpts, _ error) {
	tagLoader := &multiconfig.TagLoader{}
	return opts, tagLoader.Load(&opts)
}

func (rs *remoteSites) Create(remoteSiteOpts RemoteSiteOpts) (result RemoteSite, err error) {
	params := struct {
		RemoteSiteOpts `json:"remote_site"`
	}{remoteSiteOpts}
	err = rs.session.Request(rest.MethodPost, remoteSitesUri, params, &result)
	return result, err
}

func (rs *remoteSites) Update(siteId int, remoteSiteOpts RemoteSiteOpts) (result RemoteSite, err error) {
	params := struct {
		RemoteSiteOpts `json:"remote_site"`
	}{remoteSiteOpts}
	uri := fmt.Sprintf("%s/%d", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodPut, uri, params, &result)
	return result, err
}

func (rs *remoteSites) Connect(siteId int) (result RemoteSite, err error) {
	uri := fmt.Sprintf("%s/%d/connect", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodPost, uri, nil, &result)
	return result, err
}

type RemoteDisconnectOpts struct {
	SkipRemoteSite bool `json:"skip_remote_site" default:"true"`
}

func (rs *remoteSites) Disconnect(siteId int, remoteDisconnectOpts *RemoteDisconnectOpts) (result RemoteSite, err error) {
	uri := fmt.Sprintf("%s/%d/disconnect", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodPost, uri, remoteDisconnectOpts, &result)
	return result, err
}

func (rs *remoteSites) Delete(siteId int) (result RemoteSite, err error) {
	uri := fmt.Sprintf("%s/%d", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodDelete, uri, nil, &result)
	return result, err
}

type RemoteDataContainers struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	UsedCapacity struct {
		Bytes int64 `json:"bytes"`
	} `json:"used_capacity"`
	NamespaceScope string `json:"namespace_scope"`
	DataType       string `json:"data_type"`
	Policy         struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Dedup       int       `json:"dedup"`
		Compression int       `json:"compression"`
		Replication int       `json:"replication"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		SoftQuota   struct {
			Bytes int `json:"bytes"`
		} `json:"soft_quota"`
		HardQuota struct {
			Bytes int `json:"bytes"`
		} `json:"hard_quota"`
		IsTemplate bool `json:"is_template"`
		IsDefault  bool `json:"is_default"`
		TenantID   int  `json:"tenant_id"`
	} `json:"policy"`
	PolicyID    int         `json:"policy_id"`
	Dedup       interface{} `json:"dedup"`
	Compression interface{} `json:"compression"`
	SoftQuota   struct {
		Bytes int `json:"bytes"`
	} `json:"soft_quota"`
	HardQuota struct {
		Bytes int `json:"bytes"`
	} `json:"hard_quota"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ExportsCount   int       `json:"exports_count"`
	DirPermissions int       `json:"dir_permissions"`
	DirUID         int       `json:"dir_uid"`
	DirGid         int       `json:"dir_gid"`
	TenantID       int       `json:"tenant_id"`
	TotalUsedBytes struct {
		Bytes int64 `json:"bytes"`
	} `json:"total_used_bytes"`
	DcPairs []struct {
		ID                   int         `json:"id"`
		UUID                 string      `json:"uuid"`
		DataContainerID      int         `json:"data_container_id"`
		RemoteSiteID         int         `json:"remote_site_id"`
		RemoteSystemName     string      `json:"remote_system_name"`
		RemoteDcID           int         `json:"remote_dc_id"`
		RemoteDcUUID         string      `json:"remote_dc_uuid"`
		RemoteDcPairID       int         `json:"remote_dc_pair_id"`
		RemoteDcPairUUID     string      `json:"remote_dc_pair_uuid"`
		DrRole               string      `json:"dr_role"`
		Rpo                  int         `json:"rpo"`
		ConnectionStatus     string      `json:"connection_status"`
		Reason               interface{} `json:"reason"`
		ReplicateAcls        bool        `json:"replicate_acls"`
		CreatedAt            time.Time   `json:"created_at"`
		UpdatedAt            time.Time   `json:"updated_at"`
		LastStartTime        time.Time   `json:"last_start_time"`
		LastEndTime          interface{} `json:"last_end_time"`
		LastDrStatus         string      `json:"last_dr_status"`
		ReplicationGatewayIP string      `json:"replication_gateway_ip"`
	} `json:"dc_pairs"`
	URL string `json:"url"`
}

func (rs *remoteSites) DataContainers(siteId int) (result []RemoteDataContainers, err error) {
	uri := fmt.Sprintf("%s/%d/data_containers", remoteSitesUri, siteId)
	err = rs.session.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}
