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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"

	"ecfs/log"
)

const (
	Version = "0.6.5"
)

type ecfsDriver struct {
	driver *csicommon.CSIDriver

	is *identityServer
	ns *nodeServer
	cs *controllerServer

	caps   []*csi.VolumeCapability_AccessMode
	cscaps []*csi.ControllerServiceCapability
}

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
	// TODO: Consider checking EMS/NFS availability
	// Pro: Early failures are easier to debug
	// Con: The system may become available later, and this would result in unnecessary failures
	// Maybe add a warning instead of failing here

	// Initialize default library driver
	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Starting driver: %v version: %v", driverName, Version)
	fs.driver = csicommon.NewCSIDriver(driverName, Version, nodeId)
	if fs.driver == nil {
		glog.Fatalln("Failed to initialize CSI driver")
	}

	fs.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
		csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
		csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
		//csi.ControllerServiceCapability_RPC_GET_CAPACITY, // TODO: Add support for GetCapacity API (what's the use case?)
		//csi.ControllerServiceCapability_RPC_LIST_VOLUMES, // TODO: Add support for ListVolumes API (what's the use case?)
		//csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME, // Not needed at the controller level
	})

	fs.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		// TODO: Check if we need to advertise other (more limited) access modes, e.g. MULTI_NODE_SINGLE_WRITER
	})

	// Create gRPC servers
	fs.is = NewIdentityServer(fs.driver)
	fs.ns = NewNodeServer(fs.driver)
	fs.cs = NewControllerServer(fs.driver)

	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(endpoint, fs.is, fs.cs, fs.ns)
	server.Wait()
}
