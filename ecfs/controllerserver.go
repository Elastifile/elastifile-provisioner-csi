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
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	//"github.com/container-storage-interface/spec/lib/go/csi" // TODO: Uncomment when switching to CSI 1.0
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Warningf("AAAAA CreateVolume - req: %+v", req) // TODO: DELME
	if err := cs.validateCreateVolumeRequest(req); err != nil {
		glog.Errorf("CreateVolumeRequest validation failed: %v", err)
		err = status.Error(codes.InvalidArgument, err.Error())
		return nil, err
	}

	// req.Parameters[SecretNamespace]
	pluginConfig, err := pluginConfig()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: Don't create eManage client for each action (will need relogin support)
	var ems emanageClient
	// TODO: Here we can get User Mapping, mount options and other user-specified params (How?)
	volOptions, err := newVolumeOptions(req.GetName(), req.GetParameters())
	if err != nil {
		glog.Errorf("Validation of volume options failed: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	glog.Warningf("AAAAA CreateVolume - volOptions: %+v", volOptions) // TODO: DELME

	capacity := req.GetCapacityRange().GetRequiredBytes()

	glog.Warningf("CCCCC CreateVolume - LimitBytes: %+v", req.GetCapacityRange().GetLimitBytes())       // TODO: DELME
	glog.Warningf("CCCCC CreateVolume - RequiredBytes: %+v", req.GetCapacityRange().GetRequiredBytes()) // TODO: DELME

	if capacity > 0 {
		volOptions.Capacity = capacity
	}
	volOptions.NfsAddress = pluginConfig.NFSServer

	glog.Infof("AAAAA CreateVolume - calling createVolume()") // TODO: DELME
	err = createVolume(ems.GetClient(), volOptions)
	if err != nil {
		glog.Errorf("failed to create volume %s: %v", req.GetName(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Info("AAAAA CreateVolume - after createVolume()") // TODO: DELME

	glog.Infof("ecfs: successfully created volume %s", volOptions.Name)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            volOptions.Name,
			CapacityBytes: int64(volOptions.Capacity),
			Attributes:    req.GetParameters(),
		},
		// TODO: Uncomment when switching to CSI 1.0
		//Volume: &csi.Volume{
		//	VolumeId:      volOptions.Name,
		//	CapacityBytes: int64(volOptions.Capacity),
		//	VolumeContext: req.GetParameters(),
		//},
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	glog.Infof("ecfs: Deleting volume %v", req.VolumeId)
	//glog.Infof("BBBBB current cache state - %+v", ctrCache) // TODO: DELME
	if err := cs.validateDeleteVolumeRequest(req); err != nil {
		glog.Errorf("DeleteVolumeRequest validation failed: %v", err)
		return nil, err
	}

	var (
		volId = volumeID(req.GetVolumeId())
		err   error
	)

	var ems emanageClient
	err = deleteVolume(ems.GetClient(), req.GetVolumeId())
	if err != nil {
		glog.Errorf("failed to delete volume %s: %v", req.GetVolumeId(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.Infof("ecfs: successfully deleted volume %s", volId)

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, cap := range req.VolumeCapabilities {
		accessMode := int32(cap.GetAccessMode().GetMode())
		glog.V(3).Infof("Checking volume capability %v (%v)",
			csi.VolumeCapability_AccessMode_Mode_name[accessMode], accessMode)
		// TODO: Consider checking the actual requested AccessMode - not the most general one
		if cap.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER {
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
	glog.V(6).Infof("ControllerPublishVolume - enter. req: %+v", req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	glog.V(6).Infof("ControllerUnpublishVolume - enter. req: %+v", req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

// TODO: Implement
func (cs *controllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	glog.V(6).Infof("CreateSnapshot - enter. req: %+v", req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

// TODO: Implement
func (cs *controllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	glog.V(6).Infof("DeleteSnapshot - enter. req: %+v", req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

// TODO: Implement
func (cs *controllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	glog.V(6).Infof("ListSnapshots - enter. req: %+v", req)
	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
}

// TODO: Implement
// Found in master of https://github.com/container-storage-interface/spec/blob/master/spec.md#rpc-interface, but not in 1.0.0
//func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
//	glog.V(6).Infof("ControllerExpandVolume - enter. req: %+v", req)
//	// Set VolumeExpansion = ONLINE
//}
