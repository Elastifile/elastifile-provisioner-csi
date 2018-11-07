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
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"os"
	"os/exec"
)

const (
	volumeMounter_fuse   = "fuse"
	volumeMounter_kernel = "kernel"
)

var (
	availableMounters []string
)

// Load available ceph mounters installed on system into availableMounters
// Called from driver.go's Run()
func loadAvailableMounters() error {
	fuseMounterProbe := exec.Command("ceph-fuse", "--version")
	kernelMounterProbe := exec.Command("mount.ceph")

	if fuseMounterProbe.Run() == nil {
		availableMounters = append(availableMounters, volumeMounter_fuse)
	}

	if kernelMounterProbe.Run() == nil {
		availableMounters = append(availableMounters, volumeMounter_kernel)
	}

	if len(availableMounters) == 0 {
		return fmt.Errorf("no ceph mounters found on system")
	}

	return nil
}

// TODO: Get rid of this interface - we only support one mounting method
type volumeMounter interface {
	mount(mountPoint string, cr *credentials, volOptions *volumeOptions, volId volumeID) error
	name() string
}

func newMounter(volOptions *volumeOptions) (volumeMounter, error) {
	// Get the mounter from the configuration

	wantMounter := volOptions.Mounter

	if wantMounter == "" {
		wantMounter = DefaultVolumeMounter
	}

	// Verify that it's available

	var chosenMounter string

	for _, availMounter := range availableMounters {
		if chosenMounter == "" {
			if availMounter == wantMounter {
				chosenMounter = wantMounter
			}
		}
	}

	if chosenMounter == "" {
		// Otherwise pick whatever is left
		chosenMounter = availableMounters[0]
	}

	// Create the mounter

	switch chosenMounter {
	case volumeMounter_fuse:
		return &fuseMounter{}, nil
	case volumeMounter_kernel:
		return &kernelMounter{}, nil
	}

	return nil, fmt.Errorf("unknown mounter '%s'", chosenMounter)
}

type fuseMounter struct{}

func mountFuse(mountPoint string, cr *credentials, volOptions *volumeOptions, volId volumeID) error {
	args := [...]string{
		mountPoint,
		"-c", getCephConfPath(volId),
		"-n", cephEntityClientPrefix + cr.id,
		"--keyring", getCephKeyringPath(volId, cr.id),
		"-r", volOptions.RootPath,
		"-o", "nonempty",
	}

	out, err := execCommand("ceph-fuse", args[:]...)
	if err != nil {
		return fmt.Errorf("ecfs: ceph-fuse failed with following error: %s\necfs: ceph-fuse output: %s", err, out)
	}

	if !bytes.Contains(out, []byte("starting fuse")) {
		return fmt.Errorf("ecfs: ceph-fuse failed:\necfs: ceph-fuse output: %s", out)
	}

	return nil
}

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

func (m *fuseMounter) mount(mountPoint string, cr *credentials, volOptions *volumeOptions, volId volumeID) error {
	if err := createMountPoint(mountPoint); err != nil {
		return err
	}

	return mountFuse(mountPoint, cr, volOptions, volId)
}

func (m *fuseMounter) name() string { return "Ceph FUSE driver" }

type kernelMounter struct{}

func mountKernel(mountPoint string, cr *credentials, volOptions *volumeOptions, volId volumeID) error {
	if err := execCommandAndValidate("modprobe", "ceph"); err != nil {
		return err
	}

	return execCommandAndValidate("mount",
		"-t", "ceph",
		fmt.Sprintf("%s", volOptions.RootPath),
		mountPoint,
		"-o",
		fmt.Sprintf("name=%s,secretfile=%s", cr.id, getCephSecretPath(volId, cr.id)),
	)
}

func (m *kernelMounter) mount(mountPoint string, cr *credentials, volOptions *volumeOptions, volId volumeID) error {
	if err := createMountPoint(mountPoint); err != nil {
		return err
	}

	return mountKernel(mountPoint, cr, volOptions, volId)
}

func (m *kernelMounter) name() string { return "Ceph kernel client" }

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
