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
	"github.com/elastifile/emanage-go/pkg/optional"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/elastifile/emanage-go/pkg/emanage"
	"github.com/elastifile/errors"
	//"github.com/elastifile/emanage-go/vendor/github.com/elastifile/errors"
)

const dcPolicy = 1 // TODO: Consider making the policy (e.g. compress/dedup) configurable

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
	dc, err := emsClient.DataContainers.Create(opt.Name, dcPolicy, &emanage.DcCreateOpts{
		SoftQuota:      int(opt.Capacity), // TODO: Consider setting soft quota at 80% of hard quota
		HardQuota:      int(opt.Capacity),
		DirPermissions: opt.ExportPermissions,
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
		UserMapping: volOptions.UserMapping,
		Uid:         optional.NewInt(volOptions.UserMappingUid),
		Gid:         optional.NewInt(volOptions.UserMappingGid),
		Access:      emanage.ExportAccessModeType(volOptions.Access),
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

	// Create Data Container
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

	// Create Export
	export, err := createExport(emsClient, volOptions)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	} else {
		volOptions.Export = &export
	}
	glog.Infof("AAAAA createVolume - export created: %+v", volOptions.DataContainer) // TODO: DELME

	return nil
}

func deleteExport(emsClient *emanage.Client, dc *emanage.DataContainer) (err error) {
	exports, err := emsClient.Exports.GetAll(&emanage.GetAllOpts{})
	if err != nil {
		return errors.WrapPrefix(err, "Failed to get exports", 0)
	}

	var found bool
	for _, export := range exports {
		if export.DataContainerId == dc.Id && export.Name == dc.Name {
			found = true
			_, err := emsClient.Exports.Delete(&export)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		glog.Infof("ecfs: deleteVolume - Export for volume %v not found. Assuming already deleted", dc.Name)
	}

	return nil
}

func deleteDataContainer(emsClient *emanage.Client, dc *emanage.DataContainer) (err error) {
	_, err = emsClient.DataContainers.Delete(dc)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete Data Container %v", dc.Name), 0)
	}
	return
}

func deleteVolume(emsClient *emanage.Client, volId string) error {
	var (
		found bool
		dc    emanage.DataContainer
	)

	dcs, err := emsClient.DataContainers.GetAll(&emanage.DcGetAllOpts{})
	if err != nil {
		return errors.WrapPrefix(err, "Failed to get Data Containers", 0)
	}

	// Find the DC to be deleted
	for _, dc = range dcs {
		if dc.Name == volId {
			found = true
			break
		}
	}
	if !found {
		glog.Infof("deleteVolume - Data Container for volume %v not found. Assuming already deleted", volId)
		return nil
	}

	err = deleteExport(emsClient, &dc)
	if err != nil {
		return err
	}

	err = deleteDataContainer(emsClient, &dc)
	if err != nil {
		return err
	}

	glog.Infof("ecfs: Deleted Data Container '%v'", dc.Name)
	return nil
}
