package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"

	"ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

const maxSnapshotNameLen = 36

func createSnapshot(emsClient *emanageClient, name string, volumeId volumeHandleType, params map[string]string) (snapshot *emanage.Snapshot, err error) {
	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Creating snapshot %v for volume %v", name, volumeId)
	glog.V(log.DEBUG).Infof("ecfs: Creating snapshot %v - parameters: %v", name, params)

	dc, err := emsClient.GetClient().GetDcByName(string(volumeId))
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by name: %v", volumeId), 0)
		return
	}

	snap := &emanage.Snapshot{Name: name, DataContainerID: dc.Id}
	snapshot, err = emsClient.Snapshots.Create(snap)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Snapshot %v for volume %v already exists - assuming duplicate request", name, volumeId)
			snapshot, err = emsClient.GetSnapshotByName(name)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", name), 0)
				return
			}
		}
		return nil, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot %v on Data Container %v", name, volumeId), 0)
	}

	return
}

func waitForSnapshotToBeDeleted(emsClient *emanageClient, snapshotId int, timeout time.Duration) (err error) {
	timeoutExpired := time.After(timeout)
	tick := time.Tick(10 * time.Second)
	var snapshot *emanage.Snapshot
	for {
		select {
		case <-tick:
			snapshot, err = emsClient.GetClient().Snapshots.GetById(snapshotId)
			if err != nil {
				if isErrorDoesNotExist(err) {
					glog.V(log.DEBUG).Infof("ecfs: Snapshot id %v not found - assuming already deleted", snapshotId)
					return nil
				}
				return errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by id %v", snapshot.ID), 0)
			}
			if snapshot.Status == ecfsSnapshotStatus_REMOVED {
				glog.V(log.DETAILED_INFO).Infof("ecfs: Snapshot delete operation reported completed by EMS - "+
					"snapshot %v, dcId: %v, status: %v", snapshot.Name, snapshot.DataContainerID, snapshot.Status)
				return err
			}
		case <-timeoutExpired:
			return errors.Errorf("Timed out waiting for snapshot %v to be deleted after %v", snapshotId, timeout)
		}
	}
}

func deleteSnapshot(emsClient *emanageClient, name string) error {
	glog.V(log.INFO).Infof("ecfs: Deleting snapshot %v", name)
	snapshot, err := emsClient.GetSnapshotByName(name)
	if err != nil {
		if isErrorDoesNotExist(err) { // This operation has to be idempotent
			glog.V(log.DEBUG).Infof("ecfs: Snapshot %v not found - assuming already deleted", name)
			return nil
		}
		if isWorkaround("EL-13618 - Failed read-dir") {
			const EL13618 = "Failed read-dir"
			if strings.Contains(err.Error(), EL13618) {
				glog.Warningf("ecfs: Snapshot delete failed due to EL-13618 - returning success to cleanup the pv. Actual error: %v", err)
				return nil
			}
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", name), 0)
	}

	// Handle subsequent requests to remove snapshot
	if snapshot.Status == ecfsSnapshotStatus_REMOVING { // Operation MUST be idempotent
		glog.V(log.DEBUG).Infof("ecfs: Requested to delete snapshot that's already being removed - snapshot %v, dcId: %v, status: %v",
			snapshot.Name, snapshot.DataContainerID, snapshot.Status)
		err = waitForSnapshotToBeDeleted(emsClient, snapshot.ID, 3*time.Minute)
		if err != nil {
			return errors.Wrap(err, 0)

		}
	}

	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Deleting export from snapshot %v (%v)", snapshot.ID, snapshot.Name)
	err = deleteExportFromSnapshot(emsClient.GetClient(), snapshot.ID)
	if err != nil {
		if !isErrorDoesNotExist(err) {
			glog.Warningf("Failed to delete export from snapshot %v (%v)", snapshot.ID, snapshot.Name)
		}
	}

	glog.V(log.DEBUG).Infof("ecfs: Calling emanage snapshot.Delete - snapshot %v, dcId: %v, status: %v",
		snapshot.Name, snapshot.DataContainerID, snapshot.Status)

	tasks := snapshot.Delete()
	if tasks.Error() != nil {
		return tasks.Error()
	}
	glog.V(log.DEBUG).Infof("ecfs: Waiting for snapshot %v to be deleted by the backend", name)
	return tasks.Wait()
}

func listSnapshots(emsClient *emanageClient, snapshotId, volumeId string, maxEntries int32, startToken string) (snapshots emanage.SnapshotList, nextToken string, err error) {
	// TODO: List pagination is not supported in eManage client (page, per_page) - see TESLA-3310

	glog.V(log.DETAILED_INFO).Info("Listing snapshots",
		"snapshotId", snapshotId, "volumeId", volumeId, "maxEntries", maxEntries, "startToken", startToken)
	if snapshotId != "" {
		glog.V(log.DEBUG).Infof("ecfs: Listing snapshots by snapshotId %v", snapshotId)
		var snapshot *emanage.Snapshot
		snapshot, err = emsClient.GetSnapshotByName(snapshotId)
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
		snapshots = append(snapshots, snapshot)
	} else if volumeId != "" {
		glog.V(log.DEBUG).Infof("ecfs: Listing snapshots by volumeId %v", volumeId)
		var dc *emanage.DataContainer
		dc, err = emsClient.GetDcByName(volumeId)
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
		snapshots, err = emsClient.Snapshots.GetByDataContainerId(dc.Id)
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
	} else {
		glog.V(log.DEBUG).Infof("ecfs: Listing all snapshots")
		snapshots, err = emsClient.Snapshots.Get()
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
	}

	return
}
