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
	"os"

	"github.com/golang/glog"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	PluginFolder = "/var/lib/kubelet/plugins/csi-ecfsfsplugin"
	Version      = "0.1.0"
)

type ecfsDriver struct {
	driver *csicommon.CSIDriver

	is *identityServer
	ns *nodeServer
	cs *controllerServer

	caps   []*csi.VolumeCapability_AccessMode
	cscaps []*csi.ControllerServiceCapability
}

var (
	driver               *ecfsDriver
	DefaultVolumeMounter string
)

func NewECFSDriver() *ecfsDriver {
	return &ecfsDriver{}
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
	}
}

func NewNodeServer(d *csicommon.CSIDriver) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func (fs *ecfsDriver) Run(driverName, nodeId, endpoint, volumeMounter string) {
	glog.Infof("Starting driver: %v version: %v", driverName, Version)
	glog.Infof("Driver: %v version: %v", driverName, Version)

	// Configuration

	if err := os.MkdirAll(controllerCacheRoot, 0755); err != nil {
		glog.Fatalf("ecfs: failed to create %s: %v", controllerCacheRoot, err)
		return
	}

	if err := loadControllerCache(); err != nil {
		glog.Errorf("ecfs: failed to read volume cache: %v", err)
	}

	// TODO: Check NFS Address availability

	//if err := loadAvailableMounters(); err != nil {
	//	glog.Fatalf("ecfs: failed to load mounters: %v", err)
	//}
	//
	//if volumeMounter != "" {
	//	if err := validateMounter(volumeMounter); err != nil {
	//		glog.Fatalln(err)
	//	} else {
	//		DefaultVolumeMounter = volumeMounter
	//	}
	//} else {
	//	// Pick the first available mounter as the default one.
	//	// The choice is biased towards "fuse" in case both
	//	// ceph fuse and kernel mounters are available.
	//	DefaultVolumeMounter = availableMounters[0]
	//}
	//
	//glog.Infof("ecfs: setting default volume mounter to %s", DefaultVolumeMounter)
	//
	// Initialize default library driver

	fs.driver = csicommon.NewCSIDriver(driverName, Version, nodeId)
	if fs.driver == nil {
		glog.Fatalln("Failed to initialize CSI driver")
	}

	fs.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})

	fs.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})

	// Create gRPC servers

	fs.is = NewIdentityServer(fs.driver)
	fs.ns = NewNodeServer(fs.driver)
	fs.cs = NewControllerServer(fs.driver)

	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(endpoint, fs.is, fs.cs, fs.ns)
	server.Wait()
}
