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

const (
	oneGB = 1073741824
)

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Warningf("AAAAA CreateVolume - req: %+v", req) // TODO: DELME
	if err := cs.validateCreateVolumeRequest(req); err != nil {
		glog.Errorf("CreateVolumeRequest validation failed: %v", err)
		err = status.Error(codes.InvalidArgument, err.Error())
		return nil, err
	}

	ecfsConfig, err := newConfig()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var ems emanageClient
	//glog.V(2).Infof("AAAAA CreateVolume - creating eManage client. ecfsConfig: %+v", ecfsConfig) // TODO: DELME
	//emsClient, err := newEmanageClient(*ecfsConfig)
	//if err != nil {
	//	glog.Errorf("failed to create emanage client %+v - err:%v", ecfsConfig, err)
	//	return nil, err
	//}

	// Configuration

	//volId := newVolumeId()
	//glog.Warningf("AAAAA CreateVolume - req.Name: %v", req.Name) // TODO: DELME

	// TODO: Here we can get params from the user regarding User Mapping, mount options etc.
	volOptions, err := newVolumeOptions(req.Name, req.GetParameters())
	if err != nil {
		glog.Errorf("Validation of volume options failed: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	glog.Warningf("AAAAA CreateVolume - volOptions: %+v", volOptions) // TODO: DELME

	//emsConfig := ecfsConfigData{VolumeID: volId}
	//if err = emsConfig.writeToFile(); err != nil {
	//	glog.Errorf("failed to write ceph config file to %s: %v", getCephConfPath(volId), err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	// Create a volume in case the user didn't provide one

	// Admin credentials are required

	//cr, err := getAdminCredentials(req.GetControllerCreateSecrets())
	//if err != nil {
	//	return nil, status.Error(codes.InvalidArgument, err.Error())
	//}

	//if err = storeCephCredentials(volId, cr); err != nil {
	//	glog.Errorf("failed to store admin credentials for '%s': %v", cr.id, err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	glog.V(2).Infof("AAAAA CreateVolume - before volOptions.Name. volOptions: %+v", volOptions) // TODO: DELME
	volOptions.Name = req.Name
	glog.V(2).Infof("AAAAA CreateVolume - after volOptions.Name") // TODO: DELME
	volOptions.Capacity = req.GetCapacityRange().GetRequiredBytes()
	if volOptions.Capacity == 0 {
		// TODO: Is this check even needed? If so, make the default configurable
		volOptions.Capacity = int64(100 * size.GiB)
	}
	volOptions.NfsAddress = ecfsConfig.NFSServer

	//XXXXXX
	// Add to ctx: nfsAddress, export.Name (volOptions are not persistent enough)
	//context.WithValue()

	glog.V(2).Infof("AAAAA CreateVolume - calling createVolume()") // TODO: DELME
	err = createVolume(ems.GetClient(), volOptions)
	if err != nil {
		glog.Errorf("failed to create volume %s: %v", req.GetName(), err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.V(2).Info("AAAAA CreateVolume - after createVolume()") // TODO: DELME

	//if _, err = createCephUser(volOptions, cr, volId); err != nil {
	//	glog.Errorf("failed to create ceph user for volume %s: %v", req.GetName(), err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	glog.Infof("ecfs: successfully created volume %s", volOptions.Name)

	if err = ctrCache.insert(&controllerCacheEntry{VolOptions: *volOptions, VolumeID: volumeID(req.Name)}); err != nil {
		glog.Errorf("failed to store a cache entry for volume %s: %v", req.Name, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            req.Name,
			CapacityBytes: int64(volOptions.Capacity),
			Attributes:    req.GetParameters(),
		},
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if err := cs.validateDeleteVolumeRequest(req); err != nil {
		glog.Errorf("DeleteVolumeRequest validation failed: %v", err)
		return nil, err
	}

	var (
		volId = volumeID(req.GetVolumeId())
		err   error
	)

	// Load volume info from cache

	ent, err := ctrCache.pop(volId)
	if err != nil {
		glog.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !ent.VolOptions.ProvisionVolume {
		// DeleteVolume() is forbidden for statically provisioned volumes!

		glog.Warningf("volume %s is provisioned statically, aborting delete", volId)
		return &csi.DeleteVolumeResponse{}, nil
	}

	defer func() {
		if err != nil {
			// Reinsert cache entry for retry
			if insErr := ctrCache.insert(ent); insErr != nil {
				glog.Errorf("failed to reinsert volume cache entry in rollback procedure for volume %s: %v", volId, err)
			}
		}
	}()

	// Deleting a volume requires admin credentials

	cr, err := getAdminCredentials(req.GetControllerDeleteSecrets())
	if err != nil {
		glog.Errorf("failed to retrieve admin credentials: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = purgeVolume(volId, cr, &ent.VolOptions); err != nil {
		glog.Errorf("failed to delete volume %s: %v", volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = deleteCephUser(cr, volId); err != nil {
		glog.Errorf("failed to delete ceph user for volume %s: %v", volId, err)
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
