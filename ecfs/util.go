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
	"os/exec"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-errors/errors"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/kubernetes/pkg/util/mount"

	"csi-provisioner-elastifile/ecfs/log"
)

type volumeDescriptorType struct {
	DcId       int
	SnapshotId int
}

type volumeIdType string

// TODO: Remove snapshot id from the volume name as it's N/A after snapshot-export-as-read-only-volume has been removed
// TODO: cont'd - consider reverting volumeId to the original K8s volume name
func newVolumeId(volumeDescriptor volumeDescriptorType) volumeIdType {
	const volumeIdTemplate = "csi-dc-%v-snap-%v"
	return volumeIdType(fmt.Sprintf(volumeIdTemplate, volumeDescriptor.DcId, volumeDescriptor.SnapshotId))
}

// parseVolumeId takes internal volumeIdType that's defined by newVolumeId() and
// returns volumeDescriptorType{} with DataContainer and Snapshot ids
func parseVolumeId(volumeId volumeIdType) (volDesc *volumeDescriptorType, err error) {
	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Parsing Volume Id %v", volumeId)
	parts := strings.Split(string(volumeId), "-")
	if len(parts) != 5 {
		err = errors.Errorf("Invalid volume id: %v", volumeId)
		return
	}

	dcId, err := strconv.Atoi(parts[2])
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Illegal Data Container Id %v in Volume Id %v", parts[2], volumeId), 0)
		return
	}

	snapshotId, err := strconv.Atoi(parts[4])
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Illegal Snapshot Id %v in Volume Id %v", parts[2], volumeId), 0)
		return
	}

	volDesc = &volumeDescriptorType{
		DcId:       dcId,
		SnapshotId: snapshotId,
	}

	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Parsed Volume Id %v into %+v", volumeId, volDesc)
	return
}

func execCommand(command string, args ...string) ([]byte, error) {
	glog.V(log.DEBUG).Infof("ecfs: Running command: %s %s", command, args)

	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}

func execCommandAndValidate(program string, args ...string) error {
	out, err := execCommand(program, args...)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Command %v failed with following output: %v",
			program, string(out)), 0)
	}

	return nil
}

var dummyMount = mount.New("")

func isMountPoint(path string) (bool, error) {
	notMnt, err := dummyMount.IsLikelyNotMountPoint(path)
	if err != nil {
		return false, status.Error(codes.Internal, err.Error())
	}

	return !notMnt, nil
}

//
// Controller service request validation
//

func (cs *controllerServer) validateCreateVolumeRequest(req *csi.CreateVolumeRequest) error {
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Invalid CreateVolumeRequest: %+v", req), 0)
		return err
	}

	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "Volume Name cannot be empty")
	}

	if req.GetVolumeCapabilities() == nil {
		return status.Error(codes.InvalidArgument, "Volume Capabilities cannot be empty")
	}

	return nil
}

func (cs *controllerServer) validateDeleteVolumeRequest(req *csi.DeleteVolumeRequest) error {
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Invalid DeleteVolumeRequest: %+v", req), 0)
	}

	return nil
}

//
// Node service request validation
//

func validateNodeStageVolumeRequest(req *csi.NodeStageVolumeRequest) error {
	if req.GetVolumeCapability() == nil {
		return fmt.Errorf("volume capability missing in request")
	}

	if req.GetVolumeId() == "" {
		return fmt.Errorf("volume ID missing in request")
	}

	if req.GetStagingTargetPath() == "" {
		return fmt.Errorf("staging target path missing in request")
	}

	//if req.GetNodeStageSecrets() == nil || len(req.GetNodeStageSecrets()) == 0 {
	//	return fmt.Errorf("stage secrets cannot be nil or empty")
	//}

	return nil
}

func validateNodeUnstageVolumeRequest(req *csi.NodeUnstageVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return fmt.Errorf("volume ID missing in request")
	}

	if req.GetStagingTargetPath() == "" {
		return fmt.Errorf("staging target path missing in request")
	}

	return nil
}

func validateNodePublishVolumeRequest(req *csi.NodePublishVolumeRequest) error {
	if req.GetVolumeCapability() == nil {
		return fmt.Errorf("volume capability missing in request")
	}

	if req.GetVolumeId() == "" {
		return fmt.Errorf("volume ID missing in request")
	}

	if req.GetTargetPath() == "" {
		return fmt.Errorf("varget path missing in request")
	}

	return nil
}

func validateNodeUnpublishVolumeRequest(req *csi.NodeUnpublishVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return fmt.Errorf("volume ID missing in request")
	}

	if req.GetTargetPath() == "" {
		return fmt.Errorf("target path missing in request")
	}

	return nil
}

func isErrorAlreadyExists(err error) bool {
	var errorAlreadyExists = []string{
		"has already been taken",
		"already exist",
	}

	for _, text := range errorAlreadyExists {
		if strings.Contains(err.Error(), text) {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Entity already exists. Error: %v", err)
			return true
		}
	}
	return false
}

func isErrorDoesNotExist(err error) bool {
	var errorDoesNotExist = []string{
		"not found",
		"not exist",
		"RecordNotFound",
	}

	for _, text := range errorDoesNotExist {
		if strings.Contains(err.Error(), text) {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Entity does not exist. Error: %v", err)
			return true
		}
	}
	return false
}

func isWorkaround(desc string) bool {
	glog.Warningf("USING WORKAROUND FOR %v", desc)
	return true
}

func truncateStr(str string, maxLen int) string {
	if len(str) > maxLen {
		return str[:maxLen]
	}
	return str
}

func copyDir(src, dst string) (err error) {
	// TODO: Add timeout

	glog.V(log.DETAILED_INFO).Infof("ecfs: Going to copy %v to %v", src, dst)
	cmd := exec.Command("cp", "-a", fmt.Sprintf("%v/.", src), dst)
	out, err := cmd.CombinedOutput()
	glog.V(log.DETAILED_INFO).Infof("ecfs: Done copying %v to %v", src, dst)
	if err != nil {
		glog.Warningf("ecfs: Copy failure output: %v", string(out))
	}
	return
}
