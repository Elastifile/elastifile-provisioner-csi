package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

func createSnapshot(emsClient *emanageClient, name string, volumeId volumeIdType, params map[string]string) (snapshot *emanage.Snapshot, err error) {
	glog.V(2).Infof("ecfs: Creating snapshot %v for volume %v", name, volumeId)
	glog.V(6).Infof("ecfs: Creating snapshot %v - parameters: %v", name, params)

	volumeDescriptor, err := parseVolumeId(volumeId)
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	_, err = emsClient.GetClient().DataContainers.GetFull(volumeDescriptor.DcId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by Volume Id: %v", volumeDescriptor.DcId), 0)
		return
	}

	snap := &emanage.Snapshot{Name: name, DataContainerID: volumeDescriptor.DcId}
	snapshot, err = emsClient.Snapshots.Create(snap)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(6).Infof("ecfs: Snapshot %v for volume %v already exists - assuming duplicate request", name, volumeId)
			snapshot, err = emsClient.GetSnapshotByName(name)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", name), 0)
				return
			}
		}
		return nil, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for Data Container %v with name %v", volumeId, name), 0)
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
					glog.V(6).Infof("ecfs: Snapshot id %v not found - assuming already deleted", snapshotId)
					return nil
				}
				return errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by id %v", snapshot.ID), 0)
			}
			if snapshot.Status == ecfsSnapshotStatus_REMOVED {
				glog.V(5).Infof("ecfs: Snapshot delete operation reported completed by EMS - snapshot %v, dcId: %v, status: %v",
					snapshot.Name, snapshot.DataContainerID, snapshot.Status)
				return err
			}
		case <-timeoutExpired:
			return errors.Errorf("Timed out waiting for snapshot %v to be deleted after %v", snapshotId, timeout)
		}
	}
}

func deleteSnapshot(emsClient *emanageClient, name string) error {
	glog.V(2).Infof("ecfs: Deleting snapshot %v", name)
	snapshot, err := emsClient.GetSnapshotByName(name)
	if err != nil {
		if isErrorDoesNotExist(err) { // This operation has to be idempotent
			glog.V(6).Infof("ecfs: Snapshot %v not found - assuming already deleted", name)
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
		glog.V(5).Infof("ecfs: Requested to delete snapshot that's already being removed - snapshot %v, dcId: %v, status: %v",
			snapshot.Name, snapshot.DataContainerID, snapshot.Status)
		err = waitForSnapshotToBeDeleted(emsClient, snapshot.ID, 3*time.Minute)
		if err != nil {

		}
	}

	glog.V(5).Infof("ecfs: Calling emanage snapshot.Delete - snapshot %v, dcId: %v, status: %v",
		snapshot.Name, snapshot.DataContainerID, snapshot.Status)
	tasks := snapshot.Delete()
	if tasks.Error() != nil {
		return tasks.Error()
	}
	glog.V(4).Infof("ecfs: Waiting for snapshot %v to be deleted by the backend", name)
	return tasks.Wait()
}

func listSnapshots(emsClient *emanageClient, snapshotId, volumeId string, maxEntries int32, startToken string) (snapshots emanage.SnapshotList, nextToken string, err error) {
	// TODO: List pagination is not supported in eManage client (page, per_page) - see TESLA-3310

	glog.V(5).Info("Listing snapshots",
		"snapshotId", snapshotId, "volumeId", volumeId, "maxEntries", maxEntries, "startToken", startToken)
	if snapshotId != "" {
		glog.V(6).Infof("ecfs: Listing snapshots by snapshotId %v", snapshotId)
		var snapshot *emanage.Snapshot
		snapshot, err = emsClient.GetSnapshotByName(snapshotId)
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
		snapshots = append(snapshots, snapshot)
	} else if volumeId != "" {
		glog.V(6).Infof("ecfs: Listing snapshots by volumeId %v", volumeId)
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
		glog.V(6).Infof("ecfs: Listing all snapshots")
		snapshots, err = emsClient.Snapshots.Get()
		if err != nil {
			err = errors.Wrap(err, 0)
			return
		}
	}

	return
}
