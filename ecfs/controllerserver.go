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
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elastifile/emanage-go/pkg/size"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Warningf("AAAAA CreateVolume - req: %+v", req) // TODO: DELME
	// From the log:
	// AAAAA CreateVolume - req: name:"pvc-2b1f2e59f27e11e8" capacity_range:<required_bytes:53687091200 > volume_capabilities:<mount:<> access_mode:<mode:SINGLE_NODE_WRITER > > parameters:<key:"csiProvisionerSecretName" value:"elastifile" > parameters:<key:"csiProvisionerSecretNamespace" value:"default" > parameters:<key:"username" value:"admin" > controller_create_secrets:<key:"password" value:"changeme\n" >
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

	volOptions.Capacity = req.GetCapacityRange().GetRequiredBytes()
	if volOptions.Capacity == 0 {
		// TODO: Make the default configurable
		volOptions.Capacity = int64(100 * size.GiB)
	}
	volOptions.NfsAddress = pluginConfig.NFSServer

	glog.V(2).Infof("AAAAA CreateVolume - calling createVolume()") // TODO: DELME
	err = createVolume(ems.GetClient(), volOptions)
	if err != nil {
		glog.Errorf("failed to create volume %s: %v", req.GetName(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.V(2).Info("AAAAA CreateVolume - after createVolume()") // TODO: DELME

	glog.Infof("ecfs: successfully created volume %s", volOptions.Name)

	//glog.Infof("BBBBB inserting volume into controller cache - %v", volOptions.Name) // TODO: DELME
	//if err = ctrCache.insert(&controllerCacheEntry{VolumeID: volumeID(volOptions.Name), VolOptions: *volOptions}); err != nil {
	//	glog.Errorf("Failed to store cache entry for volume %s: %v", volOptions.Name, err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//glog.Infof("BBBBB inserted volume into controller cache - %v (%+v)", volOptions.Name, ctrCache) // TODO: DELME

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            volOptions.Name,
			CapacityBytes: int64(volOptions.Capacity),
			Attributes:    req.GetParameters(),
		},
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

	// Load volume info from cache
	//glog.Infof("BBBBB popping volume from controller cache - %s", volId) // TODO: DELME
	//ent, err := ctrCache.pop(volId)
	//if err != nil {
	//	glog.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	//defer func() {
	//	if err != nil {
	//		// Reinsert cache entry for retry
	//		if insErr := ctrCache.insert(ent); insErr != nil {
	//			glog.Errorf("failed to reinsert volume cache entry in rollback procedure for volume %s: %v", volId, err)
	//		}
	//	}
	//}()

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
		if cap.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER {
			return &csi.ValidateVolumeCapabilitiesResponse{Supported: false, Message: ""}, nil
		}
	}
	return &csi.ValidateVolumeCapabilitiesResponse{Supported: true, Message: ""}, nil
}

//func (cs *controllerServer) ValidateVolumeCapabilities(
//	ctx context.Context,
//	req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
//	// Cephfs doesn't support Block volume
//	for _, cap := range req.VolumeCapabilities {
//		if cap.GetBlock() != nil {
//			return &csi.ValidateVolumeCapabilitiesResponse{Supported: false, Message: ""}, nil
//		}
//	}
//	return &csi.ValidateVolumeCapabilitiesResponse{Supported: true}, nil
//}
