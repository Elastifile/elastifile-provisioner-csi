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

	"github.com/elastifile/emanage-go/pkg/emanage"
	"github.com/elastifile/emanage-go/pkg/optional"
	"github.com/elastifile/errors"
)

type volumeOptions struct {
	Name           string
	NfsAddress     string
	Export         *emanage.Export
	DataContainer  *emanage.DataContainer
	UserMapping    emanage.UserMappingType
	UserMappingUid int
	UserMappingGid int
	ExportUid      optional.Int
	ExportGid      optional.Int
	Capacity       int64
	Permissions    int

	// TODO: Remove ceph-specific options
	//Monitors string `json:"monitors"`
	Pool            string `json:"pool"`
	RootPath        string `json:"rootPath"`
	Mounter         string `json:"mounter"`
	ProvisionVolume bool   `json:"provisionVolume"`
}

func validateNonEmptyField(field, fieldName string) error {
	if field == "" {
		return fmt.Errorf("parameter '%s' cannot be empty", fieldName)
	}

	return nil
}

func (o *volumeOptions) validate() error {
	// TODO: Validate Elastifile options

	//if err := validateNonEmptyField(o.Monitors, "monitors"); err != nil {
	//	return err
	//}
	//
	//if err := validateNonEmptyField(o.RootPath, "rootPath"); err != nil {
	//	if !o.ProvisionVolume {
	//		return err
	//	}
	//} else {
	//	if o.ProvisionVolume {
	//		return fmt.Errorf("Non-empty field rootPath is in conflict with provisionVolume=true")
	//	}
	//}
	//
	//if o.ProvisionVolume {
	//	if err := validateNonEmptyField(o.Pool, "pool"); err != nil {
	//		return err
	//	}
	//}
	//
	//if o.Mounter != "" {
	//	if err := validateMounter(o.Mounter); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func extractOption(dest *string, optionLabel string, options map[string]string) error {
	if opt, ok := options[optionLabel]; !ok {
		return errors.New("[IN SRC] Missing required field " + optionLabel)
	} else {
		*dest = opt
		return nil
	}
}

func validateMounter(m string) error {
	switch m {
	case volumeMounter_fuse:
	case volumeMounter_kernel:
	default:
		return fmt.Errorf("Unknown mounter '%s'. Valid options are 'fuse' and 'kernel'", m)
	}

	return nil
}

func newVolumeOptions(volumeName string, volOptions map[string]string) (opts *volumeOptions, err error) {
	var ems emanageClient
	opts = &volumeOptions{}

	glog.V(2).Infof("AAAAA newVolumeOptions - enter. volOptions: %+v", volOptions) // TODO: DELME

	configMap, _, err := GetProvisionerSettings()
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get provisioner settings", 0)
		return
	}

	opts.NfsAddress = configMap[nfsAddress]

	// Opportunistically fill out Dc and Export (useful when not creating a new volume
	opts.DataContainer, opts.Export, err = ems.GetDcExportByName(volumeName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("No Data Container & Export found for volume %v", volumeName), 0)
		glog.Infof(err.Error())
	}

	//if err = extractOption(&opts.Monitors, "monitors", volOptions); err != nil {
	//	return nil, err
	//}
	//
	//if err = extractOption(&provisionVolumeBool, "provisionVolume", volOptions); err != nil {
	//	return nil, err
	//}
	//
	//if opts.ProvisionVolume, err = strconv.ParseBool(provisionVolumeBool); err != nil {
	//	return nil, fmt.Errorf("Failed to parse provisionVolume: %v", err)
	//}
	//
	//if opts.ProvisionVolume {
	//	if err = extractOption(&opts.Pool, "pool", volOptions); err != nil {
	//		return nil, err
	//	}
	//} else {
	//	if err = extractOption(&opts.RootPath, "rootPath", volOptions); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//// This field is optional, don't check for its presence
	//extractOption(&opts.Mounter, "mounter", volOptions)

	glog.V(2).Infof("AAAAA newVolumeOptions - validating opts: %+v", opts) // TODO: DELME
	if err = opts.validate(); err != nil {
		err = errors.WrapPrefix(err, "Failed to validate new volume options", 0)
		return
	}

	glog.V(2).Infof("AAAAA newVolumeOptions - returning: %+v", opts) // TODO: DELME
	return
}
