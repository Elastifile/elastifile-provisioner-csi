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

const (
	volumeExportName   = "root"
	snapshotExportName = "se"
)

var emsConfig *config

// Connect to eManage
func newEmanageClient() (client *emanageClient, err error) {
	if emsConfig == nil {
		emsConfig, err = pluginConfig()
		if err != nil {
			err = errors.WrapPrefix(err, "Failed to get plugin configuration", 0)
			return
		}
	}

	baseURL, err := url.Parse(strings.TrimSuffix(emsConfig.EmanageURL, "\n"))
	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	glog.V(5).Infof("ecfs: Connecting to ECFS management server on %v", emsConfig.EmanageURL)
	legacyClient := emanage.NewClient(baseURL)
	client = &emanageClient{legacyClient}
	glog.V(5).Infof("ecfs: Logging into ECFS management server as %v", emsConfig.Username)
	err = client.Sessions.RetriedLoginTimeout(emsConfig.Username, emsConfig.Password, 3*time.Minute)
	if err != nil {
		glog.Warningf("Failed to log into ECFS management (%v) - %v", emsConfig, err)
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to log into eManage on %v", emsConfig.EmanageURL), 0)
		err = status.Error(codes.Unauthenticated, err.Error())
		return
	}
	return
}

func (ems *emanageClient) GetClient() *emanageClient {
	if ems.Client == nil {
		glog.V(5).Infof("ecfs: Initializing eManage client")
		tmpClient, err := newEmanageClient()
		if err != nil {
			panic(fmt.Sprintf("Failed to create eManage client. err: %v", err))
		}
		ems.Client = tmpClient.Client
		glog.V(6).Infof("ecfs: Initialized new eManage client")
	}

	return ems
}

func (ems *emanageClient) GetDcByName(dcName string) (*emanage.DataContainer, error) {
	glog.V(6).Infof("ecfs: GetDcByName - getting DCs from ECFS management")
	dcs, err := ems.GetClient().DataContainers.GetAll(nil)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Failed to list Data Containers", 0)
	}
	for _, dc := range dcs {
		if dc.Name == dcName {
			return &dc, nil
		}
	}
	return nil, errors.Errorf("Data Container '%v' not found", dcName)
}

func (ems *emanageClient) GetDcDefaultExportByVolumeId(volId volumeIdType) (*emanage.DataContainer, *emanage.Export, error) {
	glog.V(6).Infof("ecfs: Looking for DC/export by Volume Id %v", volId)

	volDesc, err := parseVolumeId(volId)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	dc, err := ems.GetClient().DataContainers.GetFull(volDesc.DcId)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	exports, err := ems.GetClient().Exports.GetAll(nil)
	if err != nil {
		return nil, nil, errors.WrapPrefix(err, "Failed to get exports from eManage", 0)
	}
	for _, export := range exports {
		if dc.Id == export.DataContainerId && export.Name == volumeExportName {
			glog.V(10).Infof("ecfs: Found Dc and Export by Volume Id - DC: %+v EXPORT: %+v", dc, export)
			return &dc, &export, nil
		}
	}
	return nil, nil, errors.Errorf("Export not found by Volume Id %v", volId)
}

func (ems *emanageClient) GetDcSnapshotExportByVolumeId(volId volumeIdType) (*emanage.DataContainer, *emanage.Export, error) {
	glog.V(6).Infof("ecfs: Looking for DC/export by Volume Id %v", volId)

	volDesc, err := parseVolumeId(volId)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	dc, err := ems.GetClient().DataContainers.GetFull(volDesc.DcId)
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	exports, err := ems.GetClient().Exports.GetAll(nil)
	if err != nil {
		return nil, nil, errors.WrapPrefix(err, "Failed to get exports from eManage", 0)
	}
	for _, export := range exports {
		if dc.Id == export.DataContainerId && export.Name == volumeExportName {
			glog.V(10).Infof("ecfs: Found Snapshot Export by Volume Id - success. Returning DC: %+v EXPORT: %+v", dc, export)
			return &dc, &export, nil
		}
	}
	return nil, nil, errors.Errorf("Export not found by Volume Id %v", volId)
}

func (ems *emanageClient) GetSnapshotByName(snapshotName string) (snapshot *emanage.Snapshot, err error) {
	glog.V(6).Infof("ecfs: Looking for snapshot named: %v", snapshotName)
	snapshots, err := ems.GetClient().Snapshots.Get()
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, snap := range snapshots {
		if snap.Name == snapshotName {
			glog.V(6).Infof("ecfs: GetSnapshotByName - matched snapshot by name %v on DC %v", snap.Name, snap.DataContainerID)
			return snap, nil
		}
	}
	return nil, errors.Errorf("Snapshot not found by name %v", snapshotName)
}

func parseTimestampRFC3339(timestamp string) (int64, error) {
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to parse timestamp '%v' - expected format '%v'", timestamp, time.RFC3339), 0)
	}
	return ts.Unix(), nil
}

func snapshotEcfsToCsi(ems *emanageClient, ecfsSnapshot *emanage.Snapshot) (csiSnapshot *csi.Snapshot, err error) {
	glog.V(6).Infof("ecfs: Converting ECFS snapshot struct to CSI: %+v", *ecfsSnapshot)
	dcId := ecfsSnapshot.DataContainerID
	dc, err := ems.GetClient().DataContainers.GetFull(dcId)
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

const (
	ecfsSnapshotStatus_ADDING    = "status_adding"
	ecfsSnapshotStatus_VALID     = "status_valid"
	ecfsSnapshotStatus_REMOVING  = "status_removing"
	ecfsSnapshotStatus_MODIFYING = "status_modifying"
	ecfsSnapshotStatus_REMOVED   = "status_removed"
)

func snapshotStatusEcfsToCsi(ecfsSnapshotStatus string) csi.SnapshotStatus_Type {
	var snapEcfs2CsiMap = map[string]csi.SnapshotStatus_Type{
		ecfsSnapshotStatus_ADDING:    csi.SnapshotStatus_UPLOADING,
		ecfsSnapshotStatus_VALID:     csi.SnapshotStatus_READY,
		ecfsSnapshotStatus_REMOVING:  csi.SnapshotStatus_UNKNOWN,
		ecfsSnapshotStatus_MODIFYING: csi.SnapshotStatus_UNKNOWN,
		ecfsSnapshotStatus_REMOVED:   csi.SnapshotStatus_UNKNOWN,
		"":                           csi.SnapshotStatus_UNKNOWN,
	}

	csiSnapshotStatus, ok := snapEcfs2CsiMap[ecfsSnapshotStatus]
	if !ok {
		glog.Warningf("ecfs: Unrecognized snapshot status %v - using csi.SnapshotStatus_UNKNOWN", ecfsSnapshotStatus)
		csiSnapshotStatus = csi.SnapshotStatus_UNKNOWN
	}

	return csiSnapshotStatus
}

func createExportOnSnapshot(emsClient *emanageClient, snapshot *emanage.Snapshot, volOptions *volumeOptions) (volumeDescriptor volumeDescriptorType, exportRef *emanage.Export, err error) {
	var (
		exportName = snapshotExportName
		exportOpts = emanage.ExportCreateForSnapshotOpts{
			Path:        "/",
			SnapShotId:  snapshot.ID,
			Access:      emanage.ExportAccessRO,
			UserMapping: volOptions.UserMapping,
			Uid:         volOptions.UserMappingUid,
			Gid:         volOptions.UserMappingGid,
		}
	)

	volumeDescriptor.DcId = snapshot.DataContainerID
	volumeDescriptor.SnapshotId = snapshot.ID

	glog.V(5).Infof("Creating export %v on snapshot %v", exportName, snapshot.Name)
	export, err := emsClient.GetClient().Exports.CreateForSnapshot(exportName, &exportOpts)
	if err != nil {
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create export %v on snapshot %v", exportName, snapshot.Name), 0)
		return
	}
	exportRef = &export

	glog.V(5).Infof("Created export %v on snapshot %v (DD id: %v)", export.Name, snapshot.Name, snapshot.DataContainerID)
	return
}

func getSnapshotExport(emsClient *emanageClient, snapshotId int) (snapshotRef *emanage.Snapshot, exportRef *emanage.Export, err error) {
	snapshotRef, err = emsClient.GetClient().Snapshots.GetById(snapshotId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by id: %v", snapshotId), 0)
		return
	}

	exports, err := emsClient.GetClient().Exports.GetAll(nil)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get exports", 0)
		return
	}

	var found bool
	for _, export := range exports {
		if snapshotRef.ID == export.SnapshotId {
			exportRef = &export
			found = true
			break
		}
	}

	if !found {
		err = errors.Errorf("Export for snapshot id %v not found", snapshotRef.ID)
		return
	}

	glog.V(10).Infof("Found snapshot export by snapshot ID %v - snapshot: %+v export: %+v",
		snapshotRef.ID, *snapshotRef, *exportRef)
	return
}

func getSnapshotExportPath(emsClient *emanageClient, snapshotId int) (snapshotExportPath string, err error) {
	snapshot, export, err := getSnapshotExport(emsClient, snapshotId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot export path by id: %v", snapshotId), 0)
		return
	}

	dc, err := emsClient.GetClient().DataContainers.GetFull(export.DataContainerId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by id: %v", export.DataContainerId), 0)
		return
	}

	snapshotExportPath = fmt.Sprintf("%v/%v_%v", dc.Name, snapshot.Name, export.Name)
	glog.V(6).Infof("ecfs: using Snapshot Export Path: %v", snapshotExportPath)
	return
}
