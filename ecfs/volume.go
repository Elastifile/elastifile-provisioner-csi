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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"path"
	"strings"

	"github.com/elastifile/emanage-go/pkg/emanage"
	"github.com/elastifile/errors"
	//"github.com/elastifile/emanage-go/vendor/github.com/elastifile/errors"
)

const (
	cephRootPrefix  = PluginFolder + "/controller/volumes/root-"
	cephVolumesRoot = "csi-volumes"

	namespacePrefix = "ns-"

	dcPolicy = 1 // TODO: Consider making the policy (e.g. compress/dedup) configurable
)

func getCephRootPath_local(volId volumeID) string {
	return cephRootPrefix + string(volId)
}

func getCephRootVolumePath_local(volId volumeID) string {
	return path.Join(getCephRootPath_local(volId), cephVolumesRoot, string(volId))
}

func getVolumeRootPath_ceph(volId volumeID) string {
	return path.Join("/", cephVolumesRoot, string(volId))
}

func getVolumeNamespace(volId volumeID) string {
	return namespacePrefix + string(volId)
}

func setVolumeAttribute(root, attrName, attrValue string) error {
	return execCommandAndValidate("setfattr", "-n", attrName, "-v", attrValue, root)
}

//func createVolume(volOptions *volumeOptions, adminCr *credentials, volId volumeID, bytesQuota int64) error {
//	cephRoot := getCephRootPath_local(volId)
//
//	if err := createMountPoint(cephRoot); err != nil {
//		return err
//	}
//
//	// RootPath is not set for a dynamically provisioned volume
//	// Access to ecfs's / is required
//	volOptions.RootPath = "/"
//
//	m, err := newMounter(volOptions)
//	if err != nil {
//		return fmt.Errorf("failed to create mounter: %v", err)
//	}
//
//	if err = m.mount(cephRoot, adminCr, volOptions, volId); err != nil {
//		return fmt.Errorf("error mounting ceph root: %v", err)
//	}
//
//	defer func() {
//		unmountVolume(cephRoot)
//		os.Remove(cephRoot)
//	}()
//
//	volOptions.RootPath = getVolumeRootPath_ceph(volId)
//	localVolRoot := getCephRootVolumePath_local(volId)
//
//	if err := createMountPoint(localVolRoot); err != nil {
//		return err
//	}
//
//	if err := setVolumeAttribute(localVolRoot, "ceph.quota.max_bytes", fmt.Sprintf("%d", bytesQuota)); err != nil {
//		return err
//	}
//
//	if err := setVolumeAttribute(localVolRoot, "ceph.dir.layout.pool", volOptions.Pool); err != nil {
//		return fmt.Errorf("%v\ncephfs: Does pool '%s' exist?", err, volOptions.Pool)
//	}
//
//	if err := setVolumeAttribute(localVolRoot, "ceph.dir.layout.pool_namespace", getVolumeNamespace(volId)); err != nil {
//		return err
//	}
//
//	return nil
//}

func dcExists(emsClient *emanage.Client, opt *volumeOptions) (found bool, err error) {
	dcList, err := emsClient.DataContainers.GetAll(&emanage.DcGetAllOpts{})
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to list Data Containers", 0)
		return
	}

	for _, d := range dcList {
		if d.Name == opt.Name {
			found = true
			break
		}
	}

	return
}

func exportExists(emsClient *emanage.Client, opt *volumeOptions) (found bool, export emanage.Export, err error) {
	exports, err := emsClient.Exports.GetAll(nil)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get exports", 0)
		return
	}

	for _, export = range exports {
		if export.Name == opt.Name && export.DataContainerId == opt.DataContainer.Id {
			glog.V(3).Infof("find export from data containers by %s", opt.Name)
			found = true
			break
		}
	}

	return
}

func alreadyCreatedError(err error) bool {
	return strings.Contains(err.Error(), "has already been taken") ||
		strings.Contains(err.Error(), "already exist")
}

func createDc(emsClient *emanage.Client, opt *volumeOptions) (*emanage.DataContainer, error) {
	// Create DataContainer for volume
	dc, err := emsClient.DataContainers.Create(opt.Name, dcPolicy, &emanage.DcCreateOpts{
		SoftQuota:      int(opt.Capacity), // TODO: Consider setting soft quota at 80% of hard quota
		HardQuota:      int(opt.Capacity),
		DirPermissions: opt.Permissions,
	})

	return &dc, err
}

func createExport(emsClient *emanage.Client, volOptions *volumeOptions) (export emanage.Export, err error) {
	found, export, err := exportExists(emsClient, volOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to check if export %v exists", volOptions.Name), 0)
		return
	}
	if found {
		glog.V(3).Infof("Export %v already exists - nothing to do", volOptions.Name)
		return
	}

	var exportOpt = &emanage.ExportCreateOpts{
		DcId:        int(volOptions.DataContainer.Id),
		Path:        "/",
		UserMapping: emanage.UserMappingNone,
		Uid:         volOptions.ExportUid,
		Gid:         volOptions.ExportGid,
	}

	export, err = emsClient.Exports.Create(volOptions.Name, exportOpt)
	if err != nil {
		if alreadyCreatedError(err) {
			glog.V(3).Infof("Export %v was recently created - nothing to do", volOptions.Name)
			err = nil
		} else {
			err = errors.Wrap(err, 0)
			return
		}
	}

	return
}

func createVolume(emsClient *emanage.Client, volOptions *volumeOptions) (err error) {
	glog.V(2).Infof("AAAAA createVolume - volOptions: %+v client: %+v", volOptions, emsClient) // TODO: DELME
	found, err := dcExists(emsClient, volOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to check if volume %v exists", volOptions.Name), 0)
		err = status.Error(codes.Internal, err.Error())
		return
	}

	glog.V(2).Infof("AAAAA createVolume - DC found: %v", found) // TODO: DELME

	if found {
		// TODO: Improve idempotency by returning success only if volume's settings (quota, user mapping etc.) match
		glog.V(3).Infof("Volume (data container) %v already exists - nothing to do", volOptions.Name)
		return status.Error(codes.AlreadyExists, err.Error())
	}

	dc, err := createDc(emsClient, volOptions)
	glog.V(2).Infof("AAAAA createVolume - createDc() err: %v, result: %+v", err, volOptions.DataContainer) // TODO: DELME
	if err != nil {
		if alreadyCreatedError(err) {
			glog.V(3).Infof("Volume %v was recently created - nothing to do", volOptions.Name)
			err = nil
			// TODO: fetch the dc anyway. Currently, volOptions.DataContainer will be assigned nil in this case
		} else {
			err = errors.Wrap(err, 0)
			return status.Error(codes.Internal, err.Error())
		}
	}
	volOptions.DataContainer = dc
	glog.V(2).Infof("AAAAA createVolume - DC created: %+v", volOptions.DataContainer) // TODO: DELME

	//cephRoot := getCephRootPath_local(volumeID(volOptions.Name))
	//
	//if err := createMountPoint(cephRoot); err != nil {
	//	return err
	//}

	export, err := createExport(emsClient, volOptions)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	} else {
		volOptions.Export = &export
	}
	glog.V(2).Infof("AAAAA createVolume - export created: %+v", volOptions.DataContainer) // TODO: DELME

	// RootPath is not set for a dynamically provisioned volume
	// Access to ecfs's / is required
	//volOptions.RootPath = "/"
	//
	//m, err := newMounter(volOptions)
	//if err != nil {
	//	return fmt.Errorf("failed to create mounter: %v", err)
	//}
	//
	//if err = m.mount(cephRoot, adminCr, volOptions, volId); err != nil {
	//	return fmt.Errorf("error mounting ceph root: %v", err)
	//}
	//
	//defer func() {
	//	unmountVolume(cephRoot)
	//	os.Remove(cephRoot)
	//}()
	//
	//volOptions.RootPath = getVolumeRootPath_ceph(volId)
	//localVolRoot := getCephRootVolumePath_local(volId)
	//
	//if err := createMountPoint(localVolRoot); err != nil {
	//	return err
	//}
	//
	//if err := setVolumeAttribute(localVolRoot, "ceph.quota.max_bytes", fmt.Sprintf("%d", bytesQuota)); err != nil {
	//	return err
	//}
	//
	//if err := setVolumeAttribute(localVolRoot, "ceph.dir.layout.pool", volOptions.Pool); err != nil {
	//	return fmt.Errorf("%v\ncephfs: Does pool '%s' exist?", err, volOptions.Pool)
	//}
	//
	//if err := setVolumeAttribute(localVolRoot, "ceph.dir.layout.pool_namespace", getVolumeNamespace(volOptions.Name)); err != nil {
	//	return err
	//}

	return nil
}

func purgeVolume(volId volumeID, adminCr *credentials, volOptions *volumeOptions) error {
	var (
		cephRoot        = getCephRootPath_local(volId)
		volRoot         = getCephRootVolumePath_local(volId)
		volRootDeleting = volRoot + "-deleting"
	)

	if err := createMountPoint(cephRoot); err != nil {
		return err
	}

	// Root path is not set for dynamically provisioned volumes
	// Access to ecfs's / is required
	volOptions.RootPath = "/"

	m, err := newMounter(volOptions)
	if err != nil {
		return fmt.Errorf("failed to create mounter: %v", err)
	}

	if err = m.mount(cephRoot, adminCr, volOptions, volId); err != nil {
		return fmt.Errorf("error mounting ceph root: %v", err)
	}

	defer func() {
		unmountVolume(volRoot)
		os.Remove(volRoot)
	}()

	if err := os.Rename(volRoot, volRootDeleting); err != nil {
		return fmt.Errorf("coudln't mark volume %s for deletion: %v", volId, err)
	}

	if err := os.RemoveAll(volRootDeleting); err != nil {
		return fmt.Errorf("failed to delete volume %s: %v", volId, err)
	}

	return nil
}
