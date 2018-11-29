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
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
}

func (ns *nodeServer) nodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.V(2).Infof("AAAAA NodePublishVolume - enter. ctx: %+v req: %+v", ctx, req) // TODO: DELME
	if err := validateNodePublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Configuration

	targetPath := req.GetTargetPath()
	volId := req.GetVolumeId()

	glog.V(2).Infof("AAAAA NodePublishVolume - createMountPoint: %v", targetPath) // TODO: DELME
	if err := createMountPoint(targetPath); err != nil {
		glog.Errorf("failed to create mount point at %s: %v", targetPath, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check if the volume is already mounted
	isMnt, err := isMountPoint(targetPath)
	if err != nil {
		glog.Errorf("stat failed: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if isMnt {
		glog.Infof("ecfs: volume %s is already bind-mounted to %s", volId, targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// Mount the volume
	if err = bindMount(req.GetStagingTargetPath(), req.GetTargetPath(), req.GetReadonly()); err != nil {
		glog.Errorf("failed to bind-mount volume %s: %v", volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.V(2).Infof("AAAAA NodePublishVolume - done. volId: %v, targetPath: %v", volId, targetPath) // TODO: DELME

	glog.Infof("ecfs: successfully bind-mounted volume %s to %s", volId, targetPath)

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.V(2).Infof("AAAAA NodePublishVolume - enter. ctx: %+v req: %+v", ctx, req) // TODO: DELME
	return ns.nodePublishVolume(ctx, req)
}

func (ns *nodeServer) nodeStageVolume1(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if err := validateNodeStageVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Configuration
	stagingTargetPath := req.GetStagingTargetPath()
	volId := volumeID(req.GetVolumeId())

	glog.V(2).Infof("AAAAA NodeStageVolume - calling newVolumeOptions(). volId: %+v", volId) // TODO: DELME
	volOptions, err := newVolumeOptions(req.VolumeId, req.GetVolumeAttributes())             // TODO: Here we rely on volume id being identical to its name. Check if the actual name is stored in its attributes
	if err != nil {
		glog.Errorf("Error reading volume options for volume %s: %v", volId, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	glog.V(2).Infof("AAAAA NodeStageVolume - calling newVolumeOptions(). volId: %+v volOptions: %+v", volId, volOptions) // TODO: DELME

	if err = createMountPoint(stagingTargetPath); err != nil {
		glog.Errorf("failed to create staging mount point at %s for volume %s: %v", stagingTargetPath, volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check if the volume is already mounted
	isMount, err := isMountPoint(stagingTargetPath)
	if err != nil {
		glog.Errorf("stat failed: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if isMount {
		glog.Infof("ecfs: volume %s is already mounted on %s, skipping", volId, stagingTargetPath)
		return &csi.NodeStageVolumeResponse{}, nil
	}

	// Mount the volume
	err = mountEcfs(stagingTargetPath, volOptions, volId)
	if err != nil {
		glog.Errorf("failed to mount volume %s: %v", volId, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	glog.Infof("ecfs: successfully mounted volume %s to %s", volId, stagingTargetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return ns.nodeStageVolume1(ctx, req)
}

func (ns *nodeServer) nodeUnpublishVolume1(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if err := validateNodeUnpublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	targetPath := req.GetTargetPath()

	// Unmount the bind-mount
	if err := unmountVolume(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	os.Remove(targetPath)

	glog.Infof("ecfs: successfully unbinded volume %s from %s", req.GetVolumeId(), targetPath)

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return ns.nodeUnpublishVolume1(ctx, req)
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	if err := validateNodeUnstageVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stagingTargetPath := req.GetStagingTargetPath()

	// Unmount the volume
	if err := unmountVolume(stagingTargetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	os.Remove(stagingTargetPath)

	glog.Infof("ecfs: successfully umounted volume %s from %s", req.GetVolumeId(), stagingTargetPath)

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
			Capabilities: []*csi.NodeServiceCapability{
				{
					Type: &csi.NodeServiceCapability_Rpc{
						Rpc: &csi.NodeServiceCapability_RPC{
							Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
						},
					},
				},
			},
		},
		nil
}
