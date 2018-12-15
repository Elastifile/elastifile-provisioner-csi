package main

import (
	"fmt"

	"github.com/golang/glog"

	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

func createSnapshot(emsClient *emanageClient, name string, dcName string, params map[string]string) (snapshot *emanage.Snapshot, err error) {
	glog.V(2).Infof("ecfs: Creating snapshot %v for volume %v", name, dcName)
	glog.V(6).Infof("ecfs: Creating snapshot %v - parameters: %v", name, params)

	dc, err := emsClient.GetDcByName(dcName)
	if err != nil {
		return nil, errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by name %v", dcName), 0)
	}

	snap := &emanage.Snapshot{Name: name, DataContainerID: dc.Id}
	snapshot, err = emsClient.Snapshots.Create(snap)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(6).Infof("ecfs: Snapshot %v for volume %v already exists - assuming duplicate request", name, dcName)
			// TODO: Make sure snapshot name is unique
			snapshot, err = emsClient.GetSnapshotByName(name)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", name), 0)
				return
			}
		}
		return nil, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for Data Container %v with name %v", dcName, name), 0)
	}

	return
}

func deleteSnapshot(emsClient *emanageClient, name string) error {
	glog.V(2).Infof("ecfs: Deleting snapshot %v", name)
	snapshot, err := emsClient.GetSnapshotByName(name)
	if err != nil {
		if isErrorDoesNotExist(err) { // This operation has to be idempotent
			glog.V(6).Infof("ecfs: Snapshot %v not found - assuming already deleted", name)
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", name), 0)
	}

	glog.V(6).Infof("AAAAA snapshot: %+v", snapshot) // TODO: DELME
	glog.V(5).Infof("ecfs: Calling emanage snapshot.Delete on snapshot %v, dcId: %v, status: %v",
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
