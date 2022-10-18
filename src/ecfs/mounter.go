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
	"strings"

	"github.com/golang/glog"

	"ecfs/log"
	"github.com/elastifile/efaasclient/efaasapi"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
)

func mountNfs(args ...string) error {
	out, err := execCommand("mount", args[:]...)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Mount failed. Output: %v", string(out)), 0)
	}
	glog.V(log.DEBUG).Infof("Command output: %v", string(out))
	return nil
}

func getNfsAddress() (string, error) {
	settings, _, err := GetPluginSettings()
	if err != nil {
		return "", errors.WrapPrefix(err, "Failed to get plugin settings", 0)
	}
	return settings[nfsAddress], nil
}

func getNfsExportEcfs(volId volumeHandleType) (nfsExport string, err error) {
	nfsAddr, err := getNfsAddress()
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get ECFS export for filesystem %v", volId), 0)
		return
	}

	return fmt.Sprintf("%v:%v/%v", nfsAddr, string(volId), volumeExportName), nil
}

func getNfsExportEfaas(volId volumeHandleType) (nfsExport string, err error) {
	client := newEfaasClient()
	fs, err := client.GetFilesystemByName(efaasGetInstanceName(), string(volId))
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS filesystem %v", volId), 0)
		return
	}
	if len(fs.Exports) == 0 {
		err = errors.Errorf("eFaaS filesystem %v has no exports", volId)
		return
	}
	nfsExport = fs.Exports[0].NfsMountPoint
	return
}

func mountEcfs(mountPoint string, volId volumeHandleType, mountFlags []string) (err error) {
	var nfsExport string

	glog.V(log.INFO).Infof("ecfs: Mounting volume %v on %v", volId, mountPoint)
	if IsEFAAS() {
		nfsExport, err = getNfsExportEfaas(volId)
	} else {
		nfsExport, err = getNfsExportEcfs(volId)
	}
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get NFS export for volume %v", volId), 0)
	}

	if err = createMountPoint(mountPoint); err != nil {
		return err
	}

	var args []string
	if len(mountFlags) > 0 {
		args = append(args, "-o", strings.Join(mountFlags, ","))
	}
	args = append(args,
		"-vvv",
		"-t", "nfs",
		nfsExport,
		mountPoint,
	)

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

// mountEmsSnapshot creates an EMS snapshot export and mounts it
func mountEmsSnapshot(mountPoint string, snapshot *emanage.Snapshot) error {
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

// mountEfaasSnapshot creates an eFaaS snapshot export and mounts it
func mountEfaasSnapshot(mountPoint string, snapName string) error {
	if err := createMountPoint(mountPoint); err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Mounting snapshot %v on %v", snapName, mountPoint)
	client := newEfaasClient()
	share, err := client.GetShare(efaasGetInstanceName(), snapName)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to get share for snapshot %v", snapName), 0)
	}

	exportPath := share.NfsMountPoint
	if isWorkaround("'default' in export path") {
		var fs efaasapi.DataContainer
		fs, err = client.GetFilesystemBySnapshotName(efaasGetInstanceName(), snapName)
		if err != nil {
			return errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by snapshot name %v", snapName), 0)
		}

		exportPath = strings.Replace(exportPath, "default", fs.Name, 1)
	}

	args := []string{
		"-vvv",
		"-t", "nfs",
		"-o", "ro,nolock,vers=3",
		exportPath,
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
