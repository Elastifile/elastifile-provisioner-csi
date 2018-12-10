package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"github.com/container-storage-interface/spec/lib/go/csi" // TODO: Uncomment when switching to CSI 1.0

	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

type emanageClient struct {
	*emanage.Client
}

const exportName = "root"

var emsConfig *config

// Connect to eManage
func newEmanageClient() (client *emanageClient, err error) {
	if emsConfig == nil {
		glog.V(2).Infof("AAAAA GetClient - initializing new eManage client") // TODO: DELME
		emsConfig, err = pluginConfig()
		if err != nil {
			err = errors.WrapPrefix(err, "Failed to get plugin configuration", 0)
			return
		}
	}

	// TODO: Add retries

	glog.V(2).Infof("AAAAA newEmanageClient - config: %+v", emsConfig) // TODO: DELME
	glog.Infof("newEmanageClient on %s", emsConfig.EmanageURL)
	baseURL, err := url.Parse(strings.TrimSuffix(emsConfig.EmanageURL, "\n"))
	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	glog.V(2).Info("AAAAA newEmanageClient - calling NewClient()") // TODO: DELME
	legacyClient := emanage.NewClient(baseURL)
	client = &emanageClient{legacyClient}
	glog.V(2).Infof("AAAAA newEmanageClient - got new client: %+v", client) // TODO: DELME
	err = client.Sessions.Login(emsConfig.Username, emsConfig.Password)
	glog.V(2).Infof("AAAAA newEmanageClient - login err: %v", err) // TODO: DELME
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to log into eManage on %v", emsConfig.EmanageURL), 0)
		err = status.Error(codes.Unauthenticated, err.Error())
		return
	}
	return
}

func (ems *emanageClient) GetClient() *emanageClient {
	if ems.Client == nil {
		glog.V(5).Infof("Initializing eManage client")
		tmpClient, err := newEmanageClient()
		if err != nil {
			panic(fmt.Sprintf("Failed to create eManage client. err: %v", err))
		}
		ems.Client = tmpClient.Client
		glog.Infof("AAAAA GetClient - initialized new eManage client - %+v", ems.Client) // TODO: DELME
	}

	return ems
}

func (ems *emanageClient) GetDcByName(dcName string) (*emanage.DataContainer, error) {
	glog.V(2).Infof("AAAAA GetDcByName - getting DCs from ems: %+v", ems) // TODO: DELME
	dcs, err := ems.GetClient().DataContainers.GetAll(nil)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Failed to list Data Containers", 0)
	}
	for _, dc := range dcs {
		if dc.Name == dcName {
			return &dc, nil
		}
	}
	return nil, errors.Errorf("Container '%v' not found", dcName)
}

func (ems *emanageClient) GetDcExportByName(dcName string) (*emanage.DataContainer, *emanage.Export, error) {
	// Here we assume the Dc and the Export have the same name
	glog.V(2).Infof("AAAAA GetDcExportByName - enter. Looking for Dc & export named: %v", dcName) // TODO: DELME
	dc, err := ems.GetDcByName(dcName)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	exports, err := ems.GetClient().Exports.GetAll(nil)
	if err != nil {
		return nil, nil, errors.WrapPrefix(err, "Failed to get exports from eManage", 0)
	}

	for _, export := range exports {
		if dc.Id == export.DataContainerId && export.Name == exportName {
			glog.V(2).Infof("AAAAA GetDcExportByName - success. Returning DC: %+v EXPORT: %+v", dc, export) // TODO: DELME
			return dc, &export, nil
		}
	}
	return nil, nil, errors.Errorf("Export not found by DataContainer&Export name", dcName)
}

func (ems *emanageClient) GetSnapshotByName(snapshotName string) (snapshot *emanage.Snapshot, err error) {
	// Here we assume the Dc and the Export have the same name
	glog.V(2).Infof("AAAAA GetSnapshotByName - enter. Looking for snapshot named: %v", snapshotName) // TODO: DELME
	snapshots, err := ems.GetClient().Snapshots.Get()
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, snap := range snapshots {
		// TODO: Fix the potential issue with snapshot names being unique per-DC, while in K8s they *might* be cluster-wide
		// Find a way to make sure snapshot belongs to a specific volume (e.g. by prepending the volume name)
		if snap.Name == snapshotName {
			glog.V(6).Infof("ecfs: GetSnapshotByName - matched snapshot by name %v on DC %v", snap.Name, snap.DataContainerID)
			return
		}
	}
	return nil, errors.Errorf("Snapshot not found by name %v", snapshotName)
}

func parseTimestampRFC3339(timestamp string) (int64, error) {
	const RFC3339 = "2016-01-01T10:02:00.000Z"
	ts, err := time.Parse(RFC3339, timestamp)
	if err != nil {
		return 0, errors.WrapPrefix(err, "Failed to parse timestamp '%v' - expected format '%v'", 0)
	}
	return ts.Unix(), nil
}

func snapshotEcfsToCsi(ems *emanageClient, ecfsSnapshot *emanage.Snapshot) (csiSnapshot *csi.Snapshot, err error) {
	dcId := ecfsSnapshot.DataContainerID
	dc, err := ems.DataContainers.GetFull(dcId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by id %v", dcId), 0)
		return
	}

	csiCreatedAt, err := parseTimestampRFC3339(ecfsSnapshot.CreatedAt)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	csiSnapshot = &csi.Snapshot{
		Id:             ecfsSnapshot.Name,
		SourceVolumeId: dc.Name,
		CreatedAt:      csiCreatedAt,
		Status: &csi.SnapshotStatus{
			Type: csi.SnapshotStatus_READY,
		},
	}
	return
}

func snapshotStatusEcfsToCsi(ecfsSnapshotStatus string) csi.SnapshotStatus_Type {
	var snapEcfs2CsiMap = map[string]csi.SnapshotStatus_Type{
		"status_adding":    csi.SnapshotStatus_UPLOADING,
		"status_valid":     csi.SnapshotStatus_READY,
		"status_removing":  csi.SnapshotStatus_UNKNOWN,
		"status_modifying": csi.SnapshotStatus_UNKNOWN,
		"status_removed":   csi.SnapshotStatus_UNKNOWN,
		"":                 csi.SnapshotStatus_UNKNOWN,
	}
	csiSnapshotStatus, ok := snapEcfs2CsiMap[ecfsSnapshotStatus]
	if !ok {
		glog.Warningf("ecfs: Unrecognized snapshot status %v - using csi.SnapshotStatus_UNKNOWN")
		csiSnapshotStatus = csi.SnapshotStatus_UNKNOWN
	}

	return csiSnapshotStatus
}
