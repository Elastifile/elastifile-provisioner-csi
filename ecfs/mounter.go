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

	"csi-provisioner-elastifile/ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

func mountNfs(args ...string) error {
	out, err := execCommand("mount", args[:]...)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Mount failed. Output: %v", string(out)), 0)
	}
	return nil
}

func getNfsAddress() (string, error) {
	settings, _, err := GetPluginSettings()
	if err != nil {
		return "", errors.WrapPrefix(err, "Failed to get plugin settings", 0)
	}
	return settings[nfsAddress], nil
}

func mountEcfs(mountPoint string, volId volumeHandleType) error {
	nfsAddr, err := getNfsAddress()
	if err != nil {
		return errors.WrapPrefix(err, "Failed to mount ECFS export", 0)
	}

	if err = createMountPoint(mountPoint); err != nil {
		return err
	}

	glog.V(log.INFO).Infof("ecfs: Mounting volume %v on %v", volId, mountPoint)

	dcName := string(volId)

	// TODO: Add support for mount options once mountOptions and SupportsMountOption() are supported in K8s
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options
	args := []string{
		"-vvv",
		"-t", "nfs",
		"-o", "nolock,vers=3", // TODO: Remove these defaults once mount works
		fmt.Sprintf("%v:%v/%v", nfsAddr, dcName, volumeExportName),
		mountPoint,
	}

	err = mountNfs(args...)
	if err != nil {
		if isErrorAlreadyMounted(err) && isWorkaround("b/153705643") {
			glog.V(log.DEBUG).Infof("ecfs: Mount point %v is already mounted", mountPoint)
		} else {
			return errors.WrapPrefix(err, "Failed to mount ECFS export", 0)
		}
	}

	return nil
}

// mountEcfsSnapshot creates a snapshot export and mounts it
func mountEcfsSnapshot(mountPoint string, snapshot *emanage.Snapshot) error {
	if err := createMountPoint(mountPoint); err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Mounting snapshot %v on %v", snapshot.Name, mountPoint)

	snapExportOptions := &volumeOptions{
		UserMapping: emanage.UserMappingNone,
	}

	var emsClient emanageClient
	snapshotExportPath, err := getSnapshotExportPath(emsClient.GetClient(), snapshot.ID)
	if err != nil {
		if isErrorDoesNotExist(err) { // Create export
			_, err = createExportOnSnapshot(emsClient.GetClient(), snapshot, snapExportOptions)
			if err != nil {
				return errors.WrapPrefix(err, fmt.Sprintf("Failed to create export for snapshot %v (%v)",
					snapshot.ID, snapshot.Name), 0)
			}
			snapshotExportPath, err = getSnapshotExportPath(emsClient.GetClient(), snapshot.ID)
			if err != nil {
				return errors.WrapPrefix(err, fmt.Sprintf("Failed to get newly created snapshot export path "+
					"by snapshot id %v (%v)", snapshot.ID, snapshot.Name), 0)
			}
		} else {
			return errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot export path by snapshot id %v (%v)",
				snapshot.ID, snapshot.Name), 0)
		}
	}

	nfsAddr, err := getNfsAddress()
	if err != nil {
		return errors.WrapPrefix(err, "Failed to mount ECFS export", 0)
	}

	args := []string{
		"-vvv",
		"-t", "nfs",
		"-o", "ro,nolock,vers=3",
		fmt.Sprintf("%v:%v", nfsAddr, snapshotExportPath),
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

func unmount(mountPoint string) error {
	return execCommandAndValidate("umount", mountPoint)
}

func unmountAndCleanup(mountPoint string) (err error) {
	err = unmount(mountPoint)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to unmount %v", mountPoint), 0)
	}

	err = os.Remove(mountPoint)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete mount dir %v", mountPoint), 0)
		return
	}

	return
}
