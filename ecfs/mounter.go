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
	"os"

	"github.com/golang/glog"

	"github.com/elastifile/errors"
)

func mountNfs(args ...string) error {
	out, err := execCommand("mount", args[:]...)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Mount failed. Output: %v", string(out)), 0)
	}
	return nil
}

func mountEcfs(mountPoint string, volOptions *volumeOptions, volId volumeIdType) error {
	if err := createMountPoint(mountPoint); err != nil {
		return err
	}

	glog.Infof("ECFS: Mounting volume %v on %v", volId, mountPoint)
	// TODO: Don't create eManage client for each action (will need relogin support)
	var emsClient emanageClient
	dc, export, err := emsClient.GetDcDefaultExportByVolumeId(volId)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get DC/export for Volume Id %s", volId), 0)
	}

	// TODO: Add support for mount options once mountOptions and SupportsMountOption() are supported in K8s
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options
	args := []string{
		"-vvv",
		"-t", "nfs",
		"-o", "nolock,vers=3", // TODO: Remove these defaults once mount works
		fmt.Sprintf("%v:%v/%v", volOptions.NfsAddress, dc.Name, export.Name),
		mountPoint,
	}

	err = mountNfs(args...)
	if err != nil {
		return errors.WrapPrefix(err, "Failed to mount ECFS export", 0)
	}

	return nil
}

func mountEcfsSnapshotExport(mountPoint string, volOptions *volumeOptions, volId volumeIdType) error {
	if err := createMountPoint(mountPoint); err != nil {
		return err
	}

	glog.Infof("ECFS: Mounting volume from snapshot %v on %v", volId, mountPoint)

	volDesc, err := parseVolumeId(volId)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	var emsClient emanageClient
	snapshotExportPath, err := getSnapshotExportPath(emsClient.GetClient(), volDesc.SnapshotId)
	if err != nil {
		return errors.WrapPrefix(err,
			fmt.Sprintf("Failed to get snapshot export path by snapshot id: %v", volDesc.SnapshotId), 0)
	}

	//snapshot, err := emsClient.GetClient().Snapshots.GetById(volDesc.SnapshotId)
	//if err != nil {
	//	return errors.Wrap(err, 0)
	//}

	//dc, export, err := emsClient.GetDcSnapshotExportByVolumeId(volId)
	//if err != nil {
	//	return errors.WrapPrefix(err, fmt.Sprintf("Failed to get DC/export for Volume Id %s", volId), 0)
	//}

	// TODO: Add support for mount options once mountOptions and SupportsMountOption() are supported in K8s
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options
	args := []string{
		"-vvv",
		"-t", "nfs",
		"-o", "ro,nolock,vers=3", // TODO: Remove these defaults once mount works
		fmt.Sprintf("%v:%v", volOptions.NfsAddress, snapshotExportPath),
		mountPoint,
	}

	err = mountNfs(args...)
	if err != nil {
		return errors.WrapPrefix(err, "Failed to mount ECFS snapshot export", 0)
	}

	return nil
}

func bindMount(from, to string, readOnly bool) error {
	if err := execCommandAndValidate("mount", "--bind", from, to); err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("failed to bind-mount %s to %s", from, to), 0)
	}

	if readOnly {
		if err := execCommandAndValidate("mount", "-o", "remount,ro,bind", to); err != nil {
			return errors.WrapPrefix(err, fmt.Sprintf("failed read-only remount of %s", to), 0)
		}
	}

	return nil
}

func createMountPoint(mountPoint string) error {
	return os.MkdirAll(mountPoint, 0750)
}

func unmountVolume(mountPoint string) error {
	return execCommandAndValidate("umount", mountPoint)
}
