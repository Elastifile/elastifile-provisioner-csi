/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

type volumeHandleType string

const dcPolicy = 1 // TODO: Consider making the policy (e.g. compress/dedup) configurable

func dcExists(emsClient *emanageClient, opt *volumeOptions) (bool, error) {
	_, err := emsClient.GetClient().GetDcByName(string(opt.VolumeId))
	if err != nil {
		if isErrorDoesNotExist(err) {
			return false, nil
		}
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Containers by name: %v", opt.VolumeId), 0)
		return false, err
	}

	return true, nil
}

func exportExists(emsClient *emanageClient, exportName string, opt *volumeOptions) (found bool, export emanage.Export, err error) {
	exports, err := emsClient.Exports.GetAll(nil)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get exports", 0)
		return
	}

	for _, export = range exports {
		if export.Name == exportName && export.DataContainerId == opt.DataContainer.Id {
			glog.V(log.DEBUG).Infof("find export from data containers by id %v", opt.VolumeId)
			found = true
			break
		}
	}

	return
}

func createDc(emsClient *emanageClient, opt *volumeOptions) (*emanage.DataContainer, error) {
	dc, err := emsClient.DataContainers.Create(string(opt.VolumeId), dcPolicy, &emanage.DcCreateOpts{
		// TODO: Consider setting soft quota at a %% of capacity (to be set via storageclass)
		SoftQuota:      int(opt.Capacity),
		HardQuota:      int(opt.Capacity),
		DirPermissions: opt.ExportPermissions,
	})
	return &dc, err
}

func createExportForVolume(emsClient *emanageClient, volOptions *volumeOptions) (export emanage.Export, err error) {
	exportOpt := &emanage.ExportCreateForVolumeOpts{
		DcId:        int(volOptions.DataContainer.Id),
		Path:        "/",
		UserMapping: volOptions.UserMapping,
		Uid:         volOptions.UserMappingUid,
		Gid:         volOptions.UserMappingGid,
		Access:      emanage.ExportAccessModeType(volOptions.Access),
	}

	export, err = emsClient.Exports.CreateForVolume(volumeExportName, exportOpt)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Export for volume %v was recently created - nothing to do", volOptions.VolumeId)
			err = nil
		} else {
			err = errors.Wrap(err, 0)
			return
		}
	}

	return
}

func createEmptyVolume(emsClient *emanageClient, volOptions *volumeOptions) (volumeId volumeHandleType, err error) {
	glog.V(log.DETAILED_INFO).Infof("ecfs: Creating Volume - settings: %+v", volOptions)
	volumeId = volOptions.VolumeId

	// Create Data Container
	var dc *emanage.DataContainer
	dc, err = createDc(emsClient, volOptions)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Volume %v was already created - assuming it was created "+
				"during previous, failed, attempt", volumeId)
			var e error
			dc, e = emsClient.GetDcByName(string(volumeId))
			if e != nil {
				logSecondaryError(err, e)
				return
			}
		} else {
			err = errors.Wrap(err, 0)
			return "", errors.Wrap(err, 0)
		}
	}
	volOptions.DataContainer = dc
	glog.V(log.DEBUG).Infof("ecfs: Data Container created: %+v", volOptions.DataContainer.Name)

	// Create Export
	export, err := createExportForVolume(emsClient, volOptions)
	if err != nil {
		return "", errors.Wrap(err, 0)
	} else {
		volOptions.Export = &export
	}
	glog.V(log.DEBUG).Infof("ecfs: Export %v created on Data Container %v",
		volOptions.Export.Name, volOptions.DataContainer.Name)

	glog.V(log.DEBUG).Infof("ecfs: Created volume with id %v", volumeId)

	return
}

// createVolumeFromSnapshot is intended to be used by restore/clone functions
func createVolumeFromSnapshot(emsClient *emanageClient, srcSnapName string, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	var srcSnapMountPath = fmt.Sprintf("/mnt/%v", srcSnapName)

	glog.V(log.DETAILED_INFO).Infof("ecfs: Restoring snapshot %v - dstVolOptions: %+v", srcSnapName, dstVolOptions)

	srcSnapshot, err := emsClient.GetSnapshotByName(srcSnapName)
	if err != nil {
		err = status.Error(codes.Internal, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to to get source snapshot by name %v", srcSnapName), 0).Error())
		return
	}

	// Create destination volume
	dstVolumeId, err = createEmptyVolume(emsClient, dstVolOptions)
	if err != nil {
		if !isErrorAlreadyExists(err) {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create destination volume %v",
				dstVolOptions.VolumeId), 0)
			glog.Errorf(err.Error())
			err = status.Error(codes.Internal, err.Error())
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Destination volume %v exists. Assuming CSI retried operation "+
				"that was previous aborted", dstVolOptions.VolumeId)
		}
	}

	// Mount the source snapshot
	err = mountEmsSnapshot(srcSnapMountPath, srcSnapshot)
	if err != nil {
		isMount, e := isMountPoint(srcSnapMountPath) // TODO: Consider remounting or using unique paths
		if e != nil {
			logSecondaryError(err, e)
		}
		if !isMount {
			err = errors.WrapPrefix(err, "Failed to mount source snapshot's export", 0)
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Source snapshot is already mounted on %v - "+
				"assuming it was mounted during previous, aborted, attempt", srcSnapMountPath)
		}
	}

	defer func() { // Umount the source export
		e := unmountAndCleanup(srcSnapMountPath)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, "Failed to unmount source snapshot", 0)
				glog.Warning(err.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	// Mount the destination volume
	dstVolMountPath := fmt.Sprintf("/mnt/%v", dstVolumeId)
	err = mountEcfs(dstVolMountPath, dstVolumeId, []string{"vers=3"})
	if err != nil {
		isMount, e := isMountPoint(dstVolMountPath) // TODO: Consider remounting or using unique paths
		if e != nil {
			logSecondaryError(err, e)
		}
		if !isMount {
			err = errors.WrapPrefix(err, "Failed to mount destination volume", 0)
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Destination volume is already mounted on %v - "+
				"assuming it was mounted during previous, aborted, attempt", dstVolMountPath)
		}
	}

	defer func() { // Umount the destination volume
		e := unmountAndCleanup(dstVolMountPath)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, "Failed to unmount destination volume", 0)
				glog.Warning(err.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	// Copy the source snapshot's contents into the destination volume
	err = copyDirWithKeepalive(srcSnapMountPath, dstVolMountPath, string(dstVolumeId))
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to copy snapshot %v (%v) contents to volume %v (%v)",
			srcSnapName, srcSnapMountPath, dstVolumeId, dstVolMountPath), 0)
		return
	}

	// Testing instrumentation
	delaySec := getDebugValueInt(debugValueCloneDelaySec, nil)
	if delaySec > 0 {
		glog.V(log.VERBOSE_DEBUG).Infof("ecfs: DEBUG - delaying volume creation completion by %v sec", delaySec)
		time.Sleep(time.Duration(delaySec) * time.Second)
	}

	return
}

func copyDirWithKeepalive(srcSnapMountPath, dstVolMountPath, volumeName string) (err error) {
	var (
		errChan  = make(chan error)
		stopChan = make(chan struct{})
	)

	cachedVolume, cacheHit := volumeCache.Get(volumeName)
	if !cacheHit {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get cached volume by name %v", volumeName), 0)
	}

	go cachedVolume.persistentResource.KeepAliveRoutine(errChan, stopChan, 30*24*time.Hour)
	defer close(stopChan)

	err = copyDir(srcSnapMountPath, dstVolMountPath)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to copy %v to %v", srcSnapMountPath, dstVolMountPath), 0)
		return
	}

	select { // Non-blocking read from channel
	case e, ok := <-errChan:
		if ok { // There was some data to read
			if e != nil {
				glog.Warningf("Copying data finished, but there was an error in keep-alive: %v", err)
			}
		}
	default:
		glog.V(log.DEBUG).Infof("Copying data finished")
	}

	// Testing instrumentation
	delaySec := getDebugValueInt(debugValueCopyDelaySec, nil)
	if delaySec > 0 {
		glog.V(log.VERBOSE_DEBUG).Infof("ecfs: DEBUG - delaying copy completion by %v sec", delaySec)
		time.Sleep(time.Duration(delaySec) * time.Second)
	}

	return
}

func cloneVolume(emsClient *emanageClient, source *csi.VolumeContentSource_VolumeSource, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	var (
		reqParams   map[string]string
		srcVolumeId = volumeHandleType(source.GetVolumeId())
		srcSnapName = truncateStr(fmt.Sprintf("4-%v", dstVolOptions.VolumeId), maxSnapshotNameLen)
	)

	glog.V(log.DETAILED_INFO).Infof("ecfs: Cloning volume %v to %v via snapshot %v - dstVolOptions: %+v",
		srcVolumeId, dstVolOptions.VolumeId, srcSnapName, dstVolOptions)

	// Take source volume's snapshot
	_, err = createSnapshot(emsClient, srcSnapName, srcVolumeId, reqParams)
	if err != nil {
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for volume %v with name %v", srcVolumeId, srcSnapName), 0)
		return
	}

	defer func() { // Cleanup source snapshot
		e := deleteSnapshot(emsClient, srcSnapName)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, fmt.Sprintf("Failed to delete source snapshot %v", srcSnapName), 0)
				glog.Warning(e.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	dstVolumeId, err = createVolumeFromSnapshot(emsClient, srcSnapName, dstVolOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to clone volume %v", srcVolumeId), 0)
		return
	}

	return
}

func restoreSnapshotToVolume(emsClient *emanageClient, source *csi.VolumeContentSource_SnapshotSource, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	srcSnapName := source.GetSnapshotId()

	dstVolumeId, err = createVolumeFromSnapshot(emsClient, srcSnapName, dstVolOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to restore snapshot %v", srcSnapName), 0)
		return
	}

	return
}

func deleteExport(emsClient *emanageClient, dc *emanage.DataContainer) error {
	exports, err := emsClient.Exports.GetAll(&emanage.GetAllOpts{})
	if err != nil {
		return errors.WrapPrefix(err, "Failed to get exports", 0)
	}

	var found bool
	for _, export := range exports {
		if export.DataContainerId == dc.Id && export.Name == volumeExportName {
			found = true
			_, err = emsClient.Exports.Delete(&export)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		glog.V(log.DEBUG).Infof("ecfs: Export %v for volume %v not found. Assuming already deleted",
			volumeExportName, dc.Name)
	}

	return nil
}

func deleteExportFromSnapshot(emsClient *emanageClient, snapshotId int) error {
	exports, err := emsClient.Exports.GetAll(&emanage.GetAllOpts{})
	if err != nil {
		return errors.WrapPrefix(err, "Failed to get exports", 0)
	}

	var found bool
	for _, export := range exports {
		if export.SnapshotId == snapshotId {
			found = true
			_, err := emsClient.Exports.Delete(&export)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		glog.V(log.DEBUG).Infof("ecfs: Export from Snapshot Id %v not found. Assuming already deleted", snapshotId)
	}

	return nil
}

func deleteDataContainer(emsClient *emanageClient, dc *emanage.DataContainer) (err error) {
	_, err = emsClient.DataContainers.Delete(dc)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DEBUG).Infof("ecfs: Data Container not found - assuming already deleted")
			return nil
		}
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete Data Container %v", dc.Name), 0)
	}
	return
}

func deleteVolume(emsClient *emanageClient, volName volumeHandleType) (err error) {
	var dc *emanage.DataContainer

	dc, err = emsClient.GetDcByName(string(volName))
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DEBUG).Infof("ecfs: Data Container not found - assuming already deleted")
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Container by name %v", volName), 0)
	}

	err = deleteExport(emsClient, dc)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = deleteDataContainer(emsClient, dc)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Deleted Data Container %v (%v)", dc.Id, dc.Name)
	return nil
}
