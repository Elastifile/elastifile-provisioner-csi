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
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	//"github.com/container-storage-interface/spec/lib/go/csi" // TODO: Uncomment when switching to CSI 1.0
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"csi-provisioner-elastifile/ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func getCreateVolumeResponse(volumeId volumeHandleType, volOptions *volumeOptions, req *csi.CreateVolumeRequest) (response *csi.CreateVolumeResponse) {
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            string(volumeId),
			CapacityBytes: int64(volOptions.Capacity),
			Attributes:    req.GetParameters(),
		},
		// TODO: Uncomment when switching to CSI 1.0
		//Volume: &csi.Volume{
		//	VolumeId:      string(volumeId),
		//	CapacityBytes: int64(volOptions.Capacity),
		//	VolumeContext: req.GetParameters(),
		//},
	}
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (response *csi.CreateVolumeResponse, err error) {
	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Creating volume: %v", req.GetName())
	glog.V(log.DEBUG).Infof("ecfs: Received CreateVolumeRequest: %+v", *req)

	if err = cs.validateCreateVolumeRequest(req); err != nil {
		err = errors.WrapPrefix(err, "CreateVolumeRequest validation failed", 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	volOptions, err := newVolumeOptions(req)
	if err != nil {
		err = errors.WrapPrefix(err, "Validation of volume options failed", 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	glog.V(log.DETAILED_DEBUG).Infof("ecfs: CreateVolume options: %+v", volOptions)

	// Check if volume with this name has already been requested - this operation MUST be idempotent
	var volumeId volumeHandleType
	volume, cacheHit := volumeCache.Get(req.GetName())
	if !cacheHit {
		err = volumeCache.Create(req.GetName())
		if err != nil {
			err = errors.WrapPrefix(err, "Failed to create volume cache entry", 0)
			glog.Errorf(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		if !volume.IsReady {
			glog.V(log.DETAILED_INFO).Infof("ecfs: Received repeat request to create volume %v, "+
				"but the volume is not ready yet", req.GetName())
			return nil, status.Error(codes.Aborted, "Volume creation is already in progress")
		} else {
			glog.V(log.DETAILED_INFO).Infof("ecfs: Received repeat request to create volume %v, "+
				"returning success", req.GetName())
			response = getCreateVolumeResponse(volumeHandleType(volume.ID), volOptions, req)
			return // Success
		}
	}

	// TODO: Don't create eManage client for each action (will need relogin support)
	// This has to wait until the new unified emanage client is available, since that one has generic relogin handler support
	var ems emanageClient

	if req.GetVolumeContentSource() == nil { // Create a regular volume, i.e. new empty Data Container
		glog.V(log.DEBUG).Infof("ecfs: Creating regular volume %v", req.GetName())
		volumeId, err = createEmptyVolume(ems.GetClient(), volOptions)
		if err != nil {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create volume %v", req.GetName()), 0)
			glog.Errorf(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else { // Create a pre-populated volume
		source := req.GetVolumeContentSource()
		sourceType := source.GetType()
		switch sourceType.(type) {
		case *csi.VolumeContentSource_Volume: // Clone volume
			srcVolume := source.GetVolume()
			glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Cloning volume %v to %v",
				srcVolume.GetVolumeId(), req.GetName())
			volumeId, err = cloneVolume(ems.GetClient(), srcVolume, volOptions)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to clone volume %v to %v",
					srcVolume.GetVolumeId(), req.GetName()), 0)
				glog.Errorf(err.Error())
				return nil, status.Error(codes.Internal, err.Error())
			}
		case *csi.VolumeContentSource_Snapshot: // Restore from snapshot
			snapshot := source.GetSnapshot()
			glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Creating volume %v from snapshot %v",
				req.GetName(), snapshot.GetSnapshotId())
			volumeId, err = restoreSnapshotToVolume(ems.GetClient(), snapshot, volOptions)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create volume %v from snapshot %v",
					req.GetName(), snapshot.GetSnapshotId()), 0)
				glog.Errorf(err.Error())
				return nil, status.Error(codes.Internal, err.Error())
			}
		default:
			err = errors.Errorf("Unsupported volume source type: %v (%v)", sourceType, source)
			glog.Errorf(err.Error())
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	err = volumeCache.Set(req.GetName(), true)
	if err != nil {
		err = errors.Errorf("ecfs: Failed to update cache entry for volume %v (%v), "+
			"but the volume was successfully created", req.GetName(), volumeId)
		glog.Warningf(err.Error())
	}
	glog.V(log.INFO).Infof("ecfs: Created volume %v", volumeId)

	response = getCreateVolumeResponse(volumeId, volOptions, req)
	return // Success
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	var (
		volName = volumeHandleType(req.GetVolumeId())
		err     error
		ems     emanageClient
	)

	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Deleting volume %v", req.GetVolumeId())
	if err = cs.validateDeleteVolumeRequest(req); err != nil {
		err = errors.WrapPrefix(err, "DeleteVolumeRequest validation failed", 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = deleteVolume(ems.GetClient(), volName)
	if err != nil {
		if isErrorDoesNotExist(err) { // Operation MUST be idempotent
			glog.V(log.DEBUG).Infof("ecfs: Volume id %v not found - assuming already deleted", volName)
			return &csi.DeleteVolumeResponse{}, nil // Success
		}
		if isWorkaround("EL-13618 - Failed read-dir for volume deletion") {
			const EL13618 = "Failed read-dir"
			if strings.Contains(err.Error(), EL13618) {
				glog.Warningf("ecfs: Data Container delete failed due to EL-13618 - returning success "+
					"to cleanup the pv. Actual error: %v", err)
				return &csi.DeleteVolumeResponse{}, nil // Success
			}
		}

		glog.Warningf("Failed to delete volume %v: %v", volName, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = volumeCache.Remove(req.GetVolumeId())
	if err != nil {
		err = errors.Errorf("ecfs: Failed to remove cache entry for volume %v, "+
			"but the volume was successfully deleted", req.GetVolumeId())
		glog.Warningf(err.Error())
	}

	glog.V(log.INFO).Infof("ecfs: Deleted volume %s", volName)
	return &csi.DeleteVolumeResponse{}, nil // Success
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, capability := range req.VolumeCapabilities {
		accessMode := int32(capability.GetAccessMode().GetMode())
		glog.V(log.DETAILED_INFO).Infof("Checking volume capability %v (%v)",
			csi.VolumeCapability_AccessMode_Mode_name[accessMode], accessMode)
		// TODO: Consider checking the actual requested AccessMode - not the most general one
		if capability.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER {
			return &csi.ValidateVolumeCapabilitiesResponse{
				Supported: false,
				// TODO: Uncomment when switching to CSI 1.0
				//Confirmed: nil,
				Message: ""}, nil
		}
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Supported: true,
		// TODO: Uncomment when switching to CSI 1.0
		//Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
		//	VolumeContext:      req.GetVolumeContext(),
		//	VolumeCapabilities: req.GetVolumeCapabilities(),
		//	Parameters:         req.GetParameters(),
		//},
		Message: ""}, nil
}

func getCreateSnapshotResponse(ecfsSnapshot *emanage.Snapshot, req *csi.CreateSnapshotRequest) (response *csi.CreateSnapshotResponse, err error) {
	creationTimestamp, err := parseTimestamp(ecfsSnapshot.CreatedAt)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	isReady := isSnapshotUsable(ecfsSnapshot)

	response = &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     ecfsSnapshot.Name,
			SourceVolumeId: req.GetSourceVolumeId(),
			CreationTime:   creationTimestamp,
			ReadyToUse:     isReady,
		},
	}

	return
}

func (cs *controllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (response *csi.CreateSnapshotResponse, err error) {
	var ems emanageClient

	if isWorkaround("80 chars limit") {
		glog.V(log.DEBUG).Infof("ecfs: Received snapshot create request - %+v", req.GetName())
		var snapshotName = req.GetName()
		k8sSnapshotPrefix := "snapshot-"
		snapshotName = strings.TrimPrefix(snapshotName, k8sSnapshotPrefix)
		if len(snapshotName) > maxSnapshotNameLen {
			//snapshotName = truncateStr(req.Name, maxSnapshotNameLen)
			err = errors.Errorf("Snapshot name exceeds max allowed length of %v characters - %v "+
				"(truncated version: %v)", maxSnapshotNameLen, req.GetName(), snapshotName)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		req.Name = snapshotName
	}

	var ecfsSnapshot *emanage.Snapshot
	snapshot, cacheHit := snapshotCache.Get(req.GetName())
	if cacheHit {
		glog.V(log.DEBUG).Infof("ecfs: Received repeat request to create snapshot %v", req.GetName())

		ecfsSnapshot, err = ems.GetClient().Snapshots.GetById(snapshot.ID)
		if err != nil {
			if isErrorDoesNotExist(err) {
				return nil, status.Error(codes.Aborted, "Snapshot creation is already in progress")
			}
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by id %v", snapshot.ID), 0)
			return nil, status.Error(codes.Internal, err.Error())
		}

		response, err = getCreateSnapshotResponse(ecfsSnapshot, req)
		if err != nil {
			err = errors.Wrap(err, 0)
			return nil, status.Error(codes.Internal, err.Error())
		}

		glog.V(log.INFO).Infof("ecfs: Created snapshot %v (repeat response). Ready: %v",
			req.GetName(), response.Snapshot.ReadyToUse)
		return // Success
	}

	volumeId := volumeHandleType(req.GetSourceVolumeId())
	glog.V(log.DETAILED_INFO).Infof("ecfs: Creating snapshot %v on volume %v", req.GetName(), volumeId)
	glog.V(log.DEBUG).Infof("ecfs: CreateSnapshot details. req: %+v", *req)
	ecfsSnapshot, err = createSnapshot(ems.GetClient(), req.GetName(), volumeId, req.GetParameters())
	if err != nil {
		if isErrorAlreadyExists(err) {
			// Fetching snapshot by name to prevent a race between snapshot creation and snapshot id update in cache
			// Assumption - snapshot name is unique in K8s cluster
			var e error
			ecfsSnapshot, e = ems.GetClient().GetSnapshotByName(req.GetName())
			if e != nil {
				e = errors.WrapPrefix(e,
					fmt.Sprintf("Failed to create snapshot for volume %v with name %v", volumeId, req.GetName()), 0)
			}
			response, err = getCreateSnapshotResponse(ecfsSnapshot, req)
			if err != nil {
				err = errors.WrapPrefix(err, "Snapshot has been created, "+
					"but there was a problem generating proper response", 0)
				return nil, status.Error(codes.Internal, err.Error())
			}
			return // Success
		}
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for volume %v with name %v", volumeId, req.GetName()), 0)
		return nil, status.Error(codes.Internal, err.Error())
	}

	snapshotCache.Set(req.GetName(), ecfsSnapshot.ID, true)

	response, err = getCreateSnapshotResponse(ecfsSnapshot, req)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Snapshot %v created successfully, but generating response "+
			"has failed", req.Name), 0)
		return nil, status.Error(codes.Unknown, err.Error())
	}

	glog.V(log.INFO).Infof("ecfs: Created snapshot %v. Ready: %v", req.GetName(), response.Snapshot.ReadyToUse)
	return // Success
}

func (cs *controllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (response *csi.DeleteSnapshotResponse, err error) {
	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Deleting snapshot %v", req.GetSnapshotId())
	glog.V(log.DEBUG).Infof("ecfs: DeleteSnapshot - enter. req: %+v", *req)
	var ems emanageClient
	err = deleteSnapshot(ems.GetClient(), req.GetSnapshotId())
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete snapshot by id %v", req.GetSnapshotId()), 0)
		return nil, status.Error(codes.Internal, err.Error())
	}

	snapshotCache.RemoveByName(req.GetSnapshotId())

	glog.V(log.INFO).Infof("ecfs: Deleted snapshot %v", req.SnapshotId)

	response = &csi.DeleteSnapshotResponse{}
	return
}

func (cs *controllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (response *csi.ListSnapshotsResponse, err error) {
	var (
		args        []string
		description string
	)

	// Print detailed description
	if req.GetSnapshotId() != "" {
		args = append(args, fmt.Sprintf("snapshot id: %v", req.GetSnapshotId()))
	}
	if req.GetSourceVolumeId() != "" {
		args = append(args, fmt.Sprintf("volume id: %v", req.GetSourceVolumeId()))
	}
	if req.GetMaxEntries() != 0 {
		args = append(args, fmt.Sprintf("max entries: %v", req.GetMaxEntries()))
	}
	if req.GetStartingToken() != "" {
		args = append(args, fmt.Sprintf("starting token: %v", req.GetStartingToken()))
	}
	if len(args) > 0 {
		description = " - " + strings.Join(args, ", ")
	}
	glog.V(log.INFO).Infof("ecfs: Listing snapshots %v", description)
	glog.V(log.DETAILED_DEBUG).Infof("ecfs: ListSnapshots - enter. req: %+v", *req)

	var ems emanageClient
	ecfsSnapshots, nextToken, err := listSnapshots(ems.GetClient(), req.GetSnapshotId(), req.GetSourceVolumeId(),
		req.GetMaxEntries(), req.GetStartingToken())
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots. Request: %+v", *req), 0)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Listing snapshots: %+v", ecfsSnapshots)

	var listEntries []*csi.ListSnapshotsResponse_Entry
	for _, ecfsSnapshot := range ecfsSnapshots {
		var csiSnapshot *csi.Snapshot
		glog.V(log.DEBUG).Infof("ecfs: Converting ECFS snapshot struct to CSI: %+v", ecfsSnapshot)
		csiSnapshot, err = snapshotEcfsToCsi(ems.GetClient(), ecfsSnapshot)
		if err != nil {
			err = errors.Wrap(err, 0)
			return nil, status.Error(codes.Internal, err.Error())
		}
		listEntry := &csi.ListSnapshotsResponse_Entry{
			Snapshot: csiSnapshot,
		}
		listEntries = append(listEntries, listEntry)
	}

	response = &csi.ListSnapshotsResponse{
		Entries:   listEntries,
		NextToken: nextToken,
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Listed %v snapshots", len(listEntries))
	return // Success
}

// TODO: Implement (CSI 1.1)
func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	glog.V(log.DEBUG).Infof("ControllerExpandVolume - enter. req: %+v", *req)
	// Set VolumeExpansion = ONLINE
	panic("ControllerExpandVolume not implemented")
}
