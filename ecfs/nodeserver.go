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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	//"github.com/container-storage-interface/spec/lib/go/csi" // TODO: Uncomment when switching to CSI 1.0
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elastifile/errors"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
}

func (ns *nodeServer) nodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if err := validateNodePublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Configuration
	targetPath := req.GetTargetPath()
	volId := req.GetVolumeId()

	glog.V(10).Infof("ecfs: Creating mount point: %v", targetPath)
	if err := createMountPoint(targetPath); err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create mount point at %v", targetPath), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check if the volume is already mounted
	isMnt, err := isMountPoint(targetPath)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Checking path '%v' for being a mount point failed", targetPath), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if isMnt {
		glog.Infof("ecfs: volume %s is already bind-mounted to %s", volId, targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// Mount the volume
	if err = bindMount(req.GetStagingTargetPath(), req.GetTargetPath(), req.GetReadonly()); err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to bind-mount volume %v", volId), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.V(3).Infof("ecfs: Bind-mounted volume %v to %v", volId, targetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.V(5).Infof("ecfs: Publishing volume %v", req.VolumeId)
	glog.V(6).Infof("ecfs: NodePublishVolume - enter. req: %+v", *req)
	return ns.nodePublishVolume(ctx, req)
}

func (ns *nodeServer) nodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if err := validateNodeStageVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Configuration
	stagingTargetPath := req.GetStagingTargetPath()
	volId := volumeIdType(req.GetVolumeId())

	volDesc, err := parseVolumeId(volId)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	volOptions, err := newVolumeOptions(req.GetVolumeAttributes())
	//volOptions, err := newVolumeOptions(req.GetVolumeContext()) // TODO: Uncomment when switching to CSI 1.0
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Error reading volume options for volume %v", volId), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = createMountPoint(stagingTargetPath); err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create staging mount point at %v for volume %v",
			stagingTargetPath, volId), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check if the volume is already mounted
	isMount, err := isMountPoint(stagingTargetPath)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to check mount point %v", stagingTargetPath), 0)
		glog.Errorf(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if isMount {
		glog.Infof("ecfs: volume %s is already mounted on %s, skipping", volId, stagingTargetPath)
		return &csi.NodeStageVolumeResponse{}, nil
	}

	// Mount the volume
	if volDesc.SnapshotId == 0 { // Regular volume
		err = mountEcfs(stagingTargetPath, volOptions, volId)
		if err != nil {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to mount volume %v", volId), 0)
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else { // Volume from snapshot
		err = mountEcfsSnapshotExport(stagingTargetPath, volOptions, volId)
		if err != nil {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to mount volume from snapshot %v", volId), 0)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	glog.Infof("ecfs: successfully mounted volume %s to %s", volId, stagingTargetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	glog.V(6).Infof("NodeStageVolume - enter. req: %+v", *req)
	return ns.nodeStageVolume(ctx, req)
}

func (ns *nodeServer) nodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if err := validateNodeUnpublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	targetPath := req.GetTargetPath()

	// Unmount the bind-mount
	if err := unmount(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_ = os.Remove(targetPath)

	glog.V(3).Infof("ecfs: Unbinded volume %s from %s", req.GetVolumeId(), targetPath)

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	glog.V(5).Infof("ecfs: Unpublishing volume %v", req.VolumeId)
	glog.V(6).Infof("ecfs: NodeUnpublishVolume - enter. req: %+v", *req)
	return ns.nodeUnpublishVolume(ctx, req)
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	glog.V(5).Infof("ecfs: Unstaging volume %v", req.VolumeId)
	glog.V(6).Infof("ecfs: NodeUnstageVolume - enter. req: %+v", *req)
	if err := validateNodeUnstageVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stagingTargetPath := req.GetStagingTargetPath()

	// Unmount the volume
	if err := unmount(stagingTargetPath); err != nil {
		const errorNotMounted = "mountpoint not found"
		if strings.Contains(err.Error(), errorNotMounted) {
			glog.Warningf("ecfs: Unstaging failed as '%v' does not exist - for idempotency's sake assuming "+
				"it has already been removed. Error: %v", stagingTargetPath, err.Error())
			return &csi.NodeUnstageVolumeResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Delete the mount dir
	if err := os.Remove(stagingTargetPath); err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete staging mount dir %v", stagingTargetPath), 0)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.Infof("ecfs: successfully umounted volume %s from %s", req.GetVolumeId(), stagingTargetPath)

	return &csi.NodeUnstageVolumeResponse{}, nil
}

// TODO: Implement. What's the use case? Is it needed?
// Enabled via NodeServiceCapability_RPC_GET_VOLUME_STATS
// TODO: Uncomment when switching to CSI 1.0
//func (ns *nodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
//	glog.V(6).Infof("NodeGetVolumeStats - enter. req: %+v", *req)
//	return nil, status.Error(codes.Unimplemented, "QQQQQ - not yet supported")
//}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (capabilities *csi.NodeGetCapabilitiesResponse, err error) {
	glog.V(6).Infof("ecfs: NodeGetCapabilities - enter. req: %+v", *req)

	capabilities = &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
			//{
			//	Type: &csi.NodeServiceCapability_Rpc{
			//		Rpc: &csi.NodeServiceCapability_RPC{
			//			Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
			//		},
			//	},
			//},
		},
	}

	glog.V(5).Infof("ecfs: Returning node capabilities")
	glog.V(6).Infof("ecfs: Returning node capabilities: %+v", capabilities)
	return
}
