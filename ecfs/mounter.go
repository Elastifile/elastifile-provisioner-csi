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
	"github.com/golang/glog"
	"os"
)

func mountEcfs(mountPoint string, volOptions *volumeOptions, volId volumeID) error {
	if err := createMountPoint(mountPoint); err != nil {
		return err
	}

	// export = volOptions.Export.Name
	// ip = volOptions.NfsAddress

	glog.Infof("ECFS: Mounting volume %v on %v", volId, mountPoint)
	//glog.V(2).Infof("AAAAA mountEcfs. volId: %v, mountPoint: %v, volOptions: %+v", volId, mountPoint, volOptions) // TODO: DELME
	glog.V(2).Infof("AAAAA mountEcfs. Export: %+v", volOptions.Export) // TODO: DELME
	if volOptions.Export == nil {
		panic("Export is not initialized")
	}

	// TODO: Add support for mount options
	args := [...]string{
		"-vvv",
		"-t", "nfs",
		"-o", "nolock,vers=3", // TODO: Remove these defaults once mount works
		fmt.Sprintf("%v:%v/%v", volOptions.NfsAddress, volOptions.DataContainer.Name, volOptions.Export.Name),
		mountPoint,
	}

	out, err := execCommand("mount", args[:]...)
	if err != nil {
		return fmt.Errorf("ecfs: mount failed with following error: %s\necfs: mount output: %s", err, out)
	}

	return nil
}

func bindMount(from, to string, readOnly bool) error {
	if err := execCommandAndValidate("mount", "--bind", from, to); err != nil {
		return fmt.Errorf("failed to bind-mount %s to %s: %v", from, to, err)
	}

	if readOnly {
		if err := execCommandAndValidate("mount", "-o", "remount,ro,bind", to); err != nil {
			return fmt.Errorf("failed read-only remount of %s: %v", to, err)
		}
	}

	return nil
}

func unmountVolume(mountPoint string) error {
	return execCommandAndValidate("umount", mountPoint)
}

func createMountPoint(root string) error {
	return os.MkdirAll(root, 0750)
}
