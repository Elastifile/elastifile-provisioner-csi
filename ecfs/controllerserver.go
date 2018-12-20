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

	"github.com/elastifile/errors"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.V(2).Infof("ecfs: Creating volume: %v", req.GetName())
	glog.V(6).Infof("ecfs: Received CreateVolumeRequest: %+v", *req)

	// TODO: Convert Errorf() calls into WrapPrefix() ones
	if err := cs.validateCreateVolumeRequest(req); err != nil {
		glog.Errorf("CreateVolumeRequest validation failed: %v", err)
		err = status.Error(codes.InvalidArgument, err.Error())
		return nil, err
	}

	pluginConfig, err := pluginConfig()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: Don't create eManage client for each action (will need relogin support)
	// This has to wait until the new unified emanage client is available, since that one has generic relogin handler support
	var ems emanageClient
	volOptions, err := newVolumeOptions(req.GetParameters())
	if err != nil {
		glog.Errorf("Validation of volume options failed: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	glog.Warningf("AAAAA CreateVolume - volOptions: %+v", volOptions) // TODO: DELME

	capacity := req.GetCapacityRange().GetRequiredBytes()
	if capacity > 0 {
		volOptions.Capacity = capacity
	}
	volOptions.NfsAddress = pluginConfig.NFSServer

	var volumeId volumeIdType
	if req.VolumeContentSource == nil { // Create a regular volume, i.e. new Data Container
		glog.V(6).Infof("ecfs: Creating regular volume %v", req.GetName())
		volumeId, err = createVolume(ems.GetClient(), volOptions)
		if err != nil {
			glog.Errorf("Failed to create volume %v - %v", req.GetName(), err)
			err = errors.Wrap(err, 0)
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else { // Volume from snapshot
		source := req.GetVolumeContentSource()
		glog.V(6).Infof("ecfs: Creating volume %v from snapshot %v",
			req.GetName(), source.GetSnapshot().GetId())
		volumeId, err = createVolumeFromSnapshot(ems.GetClient(), volOptions, source)
		if err != nil {
			glog.Errorf("Failed to create volume %v from snapshot %v - %v",
				req.GetName(), req.VolumeContentSource.GetSnapshot().GetId(), err)
			err = errors.Wrap(err, 0)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	volOptions.VolumeId = volumeId

	glog.V(3).Infof("ecfs: Created volume %v", volOptions.VolumeId)

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
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	glog.V(2).Infof("ecfs: Deleting volume %v", req.GetVolumeId())
	if err := cs.validateDeleteVolumeRequest(req); err != nil {
		glog.Errorf("DeleteVolumeRequest validation failed: %v", err)
		return nil, err
	}

	var (
		volId = volumeIdType(req.GetVolumeId())
		err   error
		ems   emanageClient
	)

	volDesc, err := parseVolumeId(volId)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	//if volDesc.SnapshotId == 0 { // Regular volume
	//	err = deleteVolume(ems.GetClient(), volDesc) // Error is handled below
	//} else { // Volume from snapshot
	//	err = deleteVolumeFromSnapshot(ems.GetClient(), volDesc) // Error is handled below
	//}
	deleteVolumeFunc := deleteVolume
	if volDesc.SnapshotId != 0 { // Regular volume
		deleteVolumeFunc = deleteVolumeFromSnapshot
	}
	err = deleteVolumeFunc(ems.GetClient(), volDesc)
	if err != nil {
		if isErrorDoesNotExist(err) { // Operation MUST be idempotent
			glog.V(5).Infof("ecfs: Volume id %v not found - assuming already deleted", volId)
			return &csi.DeleteVolumeResponse{}, nil
		}
		glog.Warningf("Failed to delete volume %v: %v", volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.V(3).Infof("ecfs: Deleted volume %s", volId)
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, capability := range req.VolumeCapabilities {
		accessMode := int32(capability.GetAccessMode().GetMode())
		glog.V(3).Infof("Checking volume capability %v (%v)",
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

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	glog.V(2).Infof("ecfs: Publishing volume %v on node %v", req.GetVolumeId(), req.GetNodeId())
	glog.V(6).Infof("ecfs: ControllerPublishVolume - enter. req: %+v", *req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	glog.V(2).Infof("ecfs: Unpublishing volume %v on node", req.GetVolumeId(), req.GetNodeId())
	glog.V(6).Infof("ecfs: ControllerUnpublishVolume - enter. req: %+v", *req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

func (cs *controllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (response *csi.CreateSnapshotResponse, err error) {
	var ems emanageClient

	if isWorkaround("80 chars limit") {
		glog.V(6).Infof("ecfs: Received snapshot create request - %+v", req.GetName())
		const maxSnapshotNameLen = 36
		var snapshotName = req.GetName()
		k8sSnapshotPrefix := "snapshot-"
		snapshotName = strings.TrimPrefix(snapshotName, k8sSnapshotPrefix)
		if len(snapshotName) > maxSnapshotNameLen {
			err = errors.Errorf("Snapshot name exceeds max allowed length of %v - %v (short version: %v)",
				maxSnapshotNameLen, req.GetName(), snapshotName)
			//snapshotName = truncateStr(req.Name, maxSnapshotNameLen)
			return
		}
		req.Name = snapshotName
	}

	volumeId := volumeIdType(req.GetSourceVolumeId())
	glog.V(2).Infof("ecfs: Creating snapshot %v on volume %v", req.GetName(), volumeId)
	glog.V(6).Infof("ecfs: CreateSnapshot - enter. req: %+v", *req)
	ecfsSnapshot, err := createSnapshot(ems.GetClient(), req.GetName(), volumeId, req.GetParameters())
	if err != nil {
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for volume %v with name %v", volumeId, req.GetName()), 0)
		return
	}

	glog.V(6).Infof("ecfs: Parsing snapshot's CreatedAt timestamp: %v", ecfsSnapshot.CreatedAt)
	glog.V(10).Infof("ecfs: CreateSnapshot - ecfsSnapshot.CreatedAt: %+v", ecfsSnapshot.CreatedAt) // TODO: DELME
	creationTimestamp, err := parseTimestampRFC3339(ecfsSnapshot.CreatedAt)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	csiSnapshotStatus := snapshotStatusEcfsToCsi(ecfsSnapshot.Status)
	response = &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			Id:             ecfsSnapshot.Name,
			SourceVolumeId: string(volumeId),
			CreatedAt:      creationTimestamp,
			Status: &csi.SnapshotStatus{
				Type: csiSnapshotStatus,
			},
		},
	}

	glog.V(3).Infof("ecfs: Created snapshot %v", req.Name)
	return
}

func (cs *controllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (response *csi.DeleteSnapshotResponse, err error) {
	glog.V(2).Infof("ecfs: Deleting snapshot %v", req.GetSnapshotId())
	glog.V(6).Infof("ecfs: DeleteSnapshot - enter. req: %+v", *req)
	var ems emanageClient
	err = deleteSnapshot(ems.GetClient(), req.GetSnapshotId())
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete snapshot by id %v", req.GetSnapshotId()), 0)
		return
	}
	response = &csi.DeleteSnapshotResponse{}
	glog.V(3).Infof("ecfs: Deleted snapshot %v", req.SnapshotId)
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
	glog.V(2).Infof("ecfs: Listing snapshots %v", description)
	glog.V(6).Infof("ecfs: ListSnapshots - enter. req: %+v", *req)

	var ems emanageClient
	ecfsSnapshots, nextToken, err := listSnapshots(ems.GetClient(), req.GetSnapshotId(), req.GetSourceVolumeId(), req.GetMaxEntries(), req.GetStartingToken())
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots. Request: %+v", *req), 0)
		return
	}

	glog.Warningf("AAAAA ListSnapshots - ecfsSnapshots: %+v", ecfsSnapshots) // TODO: DELME

	var listEntries []*csi.ListSnapshotsResponse_Entry
	for _, ecfsSnapshot := range ecfsSnapshots {
		var csiSnapshot *csi.Snapshot
		glog.V(6).Infof("ecfs: Converting ECFS snapshot struct to CSI: %+v", ecfsSnapshot)
		csiSnapshot, err = snapshotEcfsToCsi(ems.GetClient(), ecfsSnapshot)
		if err != nil {
			err = errors.Wrap(err, 0)
			return
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

	glog.V(3).Infof("ecfs: Listed %v snapshots", len(listEntries))
	return
}

// TODO: Implement
// Found in master of https://github.com/container-storage-interface/spec/blob/master/spec.md#rpc-interface, but not in 1.0.0
//func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
//	glog.V(6).Infof("ControllerExpandVolume - enter. req: %+v", *req)
//	// Set VolumeExpansion = ONLINE
//}
