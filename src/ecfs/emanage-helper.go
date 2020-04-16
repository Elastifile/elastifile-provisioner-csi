package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

type emanageClient struct {
	*emanage.Client
}

const (
	volumeExportName   = "root"
	snapshotExportName = "e"
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

	glog.V(log.DETAILED_INFO).Infof("ecfs: Connecting to ECFS management server on %v", emsConfig.EmanageURL)
	legacyClient := emanage.NewClient(baseURL)
	client = &emanageClient{legacyClient}
	glog.V(log.DETAILED_INFO).Infof("ecfs: Logging into ECFS management server as %v", emsConfig.Username)
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
		glog.V(log.DETAILED_INFO).Infof("ecfs: Initializing eManage client")
		tmpClient, err := newEmanageClient()
		if err != nil {
			panic(fmt.Sprintf("Failed to create eManage client. err: %v", err))
		}
		ems.Client = tmpClient.Client
		glog.V(log.DEBUG).Infof("ecfs: Initialized new eManage client")
	}

	return ems
}

func (ems *emanageClient) GetDcByName(dcName string) (*emanage.DataContainer, error) {
	glog.V(log.DEBUG).Infof("ecfs: GetDcByName %v - getting DCs from ECFS management", dcName)
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

func (ems *emanageClient) GetDcDefaultExportByVolumeId(volId volumeHandleType) (*emanage.DataContainer, *emanage.Export, error) {
	glog.V(log.DEBUG).Infof("ecfs: Looking for DC/export by Volume Id %v", volId)

	dc, err := ems.GetClient().GetDcByName(string(volId))
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	exports, err := ems.GetClient().Exports.GetAll(nil)
	if err != nil {
		return nil, nil, errors.WrapPrefix(err, "Failed to get exports from eManage", 0)
	}
	for _, export := range exports {
		if dc.Id == export.DataContainerId && export.Name == volumeExportName {
			glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Found Dc and Export by Volume Id %v - DC: %+v EXPORT: %+v",
				volId, dc, export)
			return dc, &export, nil
		}
	}
	return nil, nil, errors.Errorf("Export not found by Volume Id %v", volId)
}

func (ems *emanageClient) GetDcSnapshotExportByVolumeId(volId volumeHandleType) (*emanage.DataContainer, *emanage.Export, error) {
	glog.V(log.DEBUG).Infof("ecfs: Looking for DC/export by Volume Id %v", volId)

	dc, err := ems.GetClient().GetDcByName(string(volId))
	if err != nil {
		return nil, nil, errors.Wrap(err, 0)
	}

	exports, err := ems.GetClient().Exports.GetAll(nil)
	if err != nil {
		return nil, nil, errors.WrapPrefix(err, "Failed to get exports from eManage", 0)
	}
	for _, export := range exports {
		if dc.Id == export.DataContainerId && export.Name == volumeExportName {
			glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Found Snapshot Export by Volume Id - success. "+
				"Returning DC: %+v EXPORT: %+v", dc, export)
			return dc, &export, nil
		}
	}
	return nil, nil, errors.Errorf("Export not found by Volume Id %v", volId)
}

func (ems *emanageClient) GetSnapshotByName(snapshotName string) (snapshot *emanage.Snapshot, err error) {
	glog.V(log.DEBUG).Infof("ecfs: Looking for snapshot named: %v", snapshotName)
	snapshots, err := ems.GetClient().Snapshots.Get()
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, snap := range snapshots {
		if snap.Name == snapshotName {
			glog.V(log.DEBUG).Infof("ecfs: GetSnapshotByName - matched snapshot by name %v on DC %v",
				snap.Name, snap.DataContainerID)
			return snap, nil
		}
	}
	return nil, errors.Errorf("Snapshot not found by name %v", snapshotName)
}

func (ems *emanageClient) GetSnapshotByStrId(snapshotID string) (snapshot *emanage.Snapshot, err error) {
	glog.V(log.DEBUG).Infof("ecfs: Looking for snapshot with ID: %v", snapshotID)
	snapID, err := strconv.Atoi(snapshotID)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	snapshot, err = ems.GetClient().Snapshots.GetById(snapID)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func parseTimestamp(dateTime string, format string) (ts *timestamp.Timestamp, err error) {
	parsedDateTime, err := time.Parse(format, dateTime)
	if err != nil {
		return ts, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to parse time/date '%v' - expected format '%v'", dateTime, format), 0)
	}

	ts = &timestamp.Timestamp{
		Seconds: parsedDateTime.Unix(),
		Nanos:   int32(parsedDateTime.Nanosecond()),
	}

	return ts, nil
}

func parseTimestampRFC3339(dateTime string) (ts *timestamp.Timestamp, err error) {
	return parseTimestamp(dateTime, time.RFC3339)
}

func snapshotEcfsToCsi(ems *emanageClient, ecfsSnapshot *emanage.Snapshot) (csiSnapshot *csi.Snapshot, err error) {
	glog.V(log.DEBUG).Infof("ecfs: Converting ECFS snapshot struct to CSI: %+v", *ecfsSnapshot)
	dcId := ecfsSnapshot.DataContainerID
	dc, err := ems.GetClient().DataContainers.GetFull(dcId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by id %v", dcId), 0)
		return
	}

	creationTimestamp, err := parseTimestampRFC3339(ecfsSnapshot.CreatedAt)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	csiSnapshot = &csi.Snapshot{
		SnapshotId:     ecfsSnapshot.Name,
		SourceVolumeId: dc.Name,
		CreationTime:   creationTimestamp,
		ReadyToUse:     isSnapshotUsable(ecfsSnapshot),
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

func isSnapshotUsable(snapshot *emanage.Snapshot) bool {
	return snapshot.Status == ecfsSnapshotStatus_VALID
}

func createExportOnSnapshot(emsClient *emanageClient, snapshot *emanage.Snapshot, volOptions *volumeOptions) (exportRef *emanage.Export, err error) {
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

	glog.V(log.DETAILED_INFO).Infof("Creating export %v on snapshot %v", exportName, snapshot.Name)
	export, err := emsClient.GetClient().Exports.CreateForSnapshot(exportName, &exportOpts)
	if err != nil {
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create export %v on snapshot %v", exportName, snapshot.Name), 0)
		return
	}
	exportRef = &export

	glog.V(log.DETAILED_INFO).Infof("Created export %v on snapshot %v (DD id: %v)",
		export.Name, snapshot.Name, snapshot.DataContainerID)
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

	glog.V(log.VERBOSE_DEBUG).Infof("Found snapshot export by snapshot ID %v - snapshot: %+v export: %+v",
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
	glog.V(log.DEBUG).Infof("ecfs: Using Snapshot Export Path: %v", snapshotExportPath)
	return
}
