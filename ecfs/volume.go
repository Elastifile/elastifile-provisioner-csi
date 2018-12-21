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

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

const dcPolicy = 1 // TODO: Consider making the policy (e.g. compress/dedup) configurable

func dcExists(emsClient *emanageClient, opt *volumeOptions) (bool, error) {
	volumeDescriptor, err := parseVolumeId(opt.VolumeId)
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	_, err = emsClient.GetClient().DataContainers.GetFull(volumeDescriptor.DcId)
	if err != nil {
		if isErrorDoesNotExist(err) {
			return false, nil
		}
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Data Containers by id: %v", volumeDescriptor.DcId), 0)
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
			glog.V(3).Infof("find export from data containers by id %v", opt.VolumeId)
			found = true
			break
		}
	}

	return
}

func createDc(emsClient *emanageClient, opt *volumeOptions) (*emanage.DataContainer, error) {
	dcName := fmt.Sprintf("csi-%v", uuid.NewUUID())

	dc, err := emsClient.DataContainers.Create(dcName, dcPolicy, &emanage.DcCreateOpts{
		SoftQuota:      int(opt.Capacity), // TODO: Consider setting soft quota at 80% of hard quota
		HardQuota:      int(opt.Capacity),
		DirPermissions: opt.ExportPermissions,
	})

	return &dc, err
}

func createExportForVolume(emsClient *emanageClient, volOptions *volumeOptions) (export emanage.Export, err error) {
	found, export, err := exportExists(emsClient, volumeExportName, volOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to check if export %v exists on DC %v (%v)",
			volumeExportName, volOptions.DataContainer.Id, volOptions.DataContainer.Name), 0)
		return
	}
	if found {
		glog.V(3).Infof("ecfs: Export %v for volume %v already exists - nothing to do", volumeExportName, volOptions.VolumeId)
		return
	}

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
			glog.V(3).Infof("ecfs: Export for volume %v was recently created - nothing to do", volOptions.VolumeId)
			err = nil
		} else {
			err = errors.Wrap(err, 0)
			return
		}
	}

	return
}

func createVolume(emsClient *emanageClient, volOptions *volumeOptions) (volumeId volumeIdType, err error) {
	var volumeDescriptor volumeDescriptorType

	glog.V(6).Infof("ecfs: Creating Volume - settings: %+v", volOptions)

	// TODO: IMPORTANT: For idempotency's sake, make sure duplicate requests with the same volume name are treated as such

	// Create Data Container
	//found, err := dcExists(emsClient, volOptions)
	//if err != nil {
	//	err = errors.WrapPrefix(err, fmt.Sprintf("Failed to check if volume %v exists", volOptions.VolumeId), 0)
	//	err = status.Error(codes.Internal, err.Error())
	//	return
	//}
	//glog.V(2).Infof("AAAAA createVolume - DC found: %v", found) // TODO: DELME
	//
	//if !found { // Create Data Container
	var dc *emanage.DataContainer
	dc, err = createDc(emsClient, volOptions)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(3).Infof("ecfs: Volume %v was recently created - nothing to do", volOptions.VolumeId)
			err = nil
			// TODO: fetch the dc anyway. Currently, volOptions.DataContainer will be assigned nil in this case
			panic("Fetching DC is not implemented")
		} else {
			err = errors.Wrap(err, 0)
			return volumeId, status.Error(codes.Internal, err.Error())
		}
	}
	volumeDescriptor.DcId = dc.Id
	volOptions.DataContainer = dc
	glog.V(6).Infof("ecfs: Data Container created: %+v", volOptions.DataContainer.Name)
	//} else {
	//	glog.V(3).Infof("ecfs: Volume (data container) %v already exists - nothing to do", volOptions.VolumeId)
	//
	//	//TODO: Find the dc and return its id
	//	panic("Fetching DC is not implemented")
	//
	//	//return status.Error(codes.AlreadyExists, err.Error())
	//}

	// Create Export
	export, err := createExportForVolume(emsClient, volOptions)
	if err != nil {
		return volumeId, status.Error(codes.Internal, err.Error())
	} else {
		volOptions.Export = &export
	}
	glog.V(6).Infof("ecfs: Export %v created on Data Container",
		volOptions.Export.Name, volOptions.DataContainer.Name)

	volumeId = newVolumeId(volumeDescriptor)
	glog.V(5).Infof("ecfs: Created volume with id %v", volumeId)

	return
}

func createVolumeFromSnapshot(emsClient *emanageClient, volOptions *volumeOptions, volSource *csi.VolumeContentSource) (volumeId volumeIdType, err error) {
	glog.V(6).Infof("ecfs: Creating Volume from Snapshot - volOptions: %+v", volOptions)

	// Verify volSource type
	_, ok := volSource.GetType().(*csi.VolumeContentSource_Snapshot)
	if !ok {
		err = errors.Errorf("Received bad volume volSource type - %v", volSource.GetType())
		return
	}

	snapSource := volSource.GetSnapshot()

	// Get snapshot
	// Create export on snapshot
	ecfsSnapshot, err := emsClient.GetSnapshotByName(snapSource.GetId())
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get ECFS snapshot by name", 0)
		return
	}

	// Create Export
	volumeDescriptor, export, err := createExportOnSnapshot(emsClient, ecfsSnapshot, volOptions)
	if err != nil {
		if isErrorAlreadyExists(err) { // Snapshot volume creation MUST be idempotent
			volumeDescriptor = volumeDescriptorType{
				DcId:       ecfsSnapshot.DataContainerID,
				SnapshotId: ecfsSnapshot.ID,
			}
			volumeId = newVolumeId(volumeDescriptor)
			glog.V(5).Infof("Snapshot export already exists on volume %v", volumeId)
			_, export, err = getSnapshotExport(emsClient.GetClient(), ecfsSnapshot.ID)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Snapshot Export by snapshot Id %v",
					ecfsSnapshot.ID), 0)
			}
		}
		err = status.Error(codes.Internal, errors.Wrap(err, 0).Error())
		return
	}

	volOptions.Export = export // TODO: Bad design, check if this is used anywhere
	glog.V(5).Infof("Export %v from snapshot %v created", export.Name, ecfsSnapshot.Name)

	volumeId = newVolumeId(volumeDescriptor)
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
			_, err := emsClient.Exports.Delete(&export)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		glog.V(3).Infof("ecfs: Export %v for volume %v not found. Assuming already deleted",
			volumeExportName, dc.Name)
	}

	return nil
}

func deleteExportFromSnapshot(emsClient *emanageClient, dc *emanage.DataContainer, snapshotId int) error {
	exports, err := emsClient.Exports.GetAll(&emanage.GetAllOpts{})
	if err != nil {
		return errors.WrapPrefix(err, "Failed to get exports", 0)
	}

	var found bool
	for _, export := range exports {
		if export.DataContainerId == dc.Id && export.SnapshotId == snapshotId {
			found = true
			_, err := emsClient.Exports.Delete(&export)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		glog.V(3).Infof("ecfs: Export from DC %v Snapshot Id %v not found. Assuming already deleted",
			dc.Name, snapshotId)
	}

	return nil
}

func deleteDataContainer(emsClient *emanageClient, dc *emanage.DataContainer) (err error) {
	_, err = emsClient.DataContainers.Delete(dc)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(6).Infof("ecfs: Data Container not found - assuming already deleted")
			return nil
		}
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete Data Container %v", dc.Name), 0)
	}
	return
}

func deleteVolume(emsClient *emanageClient, volDesc *volumeDescriptorType) (err error) {
	var (
		//found bool
		dc emanage.DataContainer
	)

	dc, err = emsClient.DataContainers.GetFull(volDesc.DcId)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(6).Infof("ecfs: Data Container not found - assuming already deleted")
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf(
			"Failed to get Data Container by Id %v", volDesc.DcId), 0)
	}

	err = deleteExport(emsClient, &dc)
	if err != nil {
		return err
	}

	err = deleteDataContainer(emsClient, &dc)
	if err != nil {
		return err
	}

	glog.V(2).Infof("ecfs: Deleted Data Container %v (%v)", volDesc.DcId, dc.Name)
	return nil
}

// deleteVolumeFromSnapshot deletes volume that was created from a snapshot
func deleteVolumeFromSnapshot(emsClient *emanageClient, volDesc *volumeDescriptorType) (err error) {
	var dc emanage.DataContainer

	dc, err = emsClient.DataContainers.GetFull(volDesc.DcId)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(6).Infof("ecfs: Data Container not found - assuming already deleted")
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf(
			"Failed to get Data Container by Id %v", volDesc.DcId), 0)
	}

	err = deleteExportFromSnapshot(emsClient, &dc, volDesc.SnapshotId)
	if err != nil {
		return err
	}

	glog.Infof("ecfs: Deleted Export from Snapshot %v (%v)", volDesc.DcId, dc.Name)
	return nil
}
