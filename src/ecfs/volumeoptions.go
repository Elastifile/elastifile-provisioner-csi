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
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"

	"ecfs/log"
	"github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/elastifile/errors"
	"optional"
	"size"
)

type volumeOptions struct {
	VolumeId volumeHandleType

	Export        *emanage.Export
	DataContainer *emanage.DataContainer

	NfsAddress string

	Capacity          int64
	UserMapping       emanage.UserMappingType
	UserMappingUid    int
	UserMappingGid    int
	ExportPermissions int
	ExportUid         optional.Int
	ExportGid         optional.Int
	Access            string
}

func extractOptionString(paramName StorageClassCustomParameter, options map[string]string) (value string, err error) {
	if opt, ok := options[string(paramName)]; !ok {
		err = errors.New("Missing volume parameter: " + paramName)
	} else {
		value = opt
	}
	return
}

// Strings used in storageclass configuration file
type StorageClassCustomParameter string

const (
	UserMapping       StorageClassCustomParameter = "userMapping"
	UserMappingUid    StorageClassCustomParameter = "userMappingUid"
	UserMappingGid    StorageClassCustomParameter = "userMappingGid"
	ExportUid         StorageClassCustomParameter = "exportUid"
	ExportGid         StorageClassCustomParameter = "exportGid"
	Permissions       StorageClassCustomParameter = "permissions"
	DefaultVolumeSize StorageClassCustomParameter = "defaultVolumeSize"
	Access            StorageClassCustomParameter = "access"
)

func newVolumeOptions(req *csi.CreateVolumeRequest) (*volumeOptions, error) {
	var (
		volParams = req.GetParameters()
		opts      = &volumeOptions{}
		paramStr  string
		paramInt  int
		paramSize size.Size
	)

	pluginSettings, _, err := GetPluginSettings()
	if err != nil {
		return nil, errors.WrapPrefix(err, "Failed to get provisioner settings", 0)
	}

	opts.VolumeId = volumeHandleType(req.GetName())
	opts.NfsAddress = pluginSettings[nfsAddress]

	// UserMapping
	if paramStr, err = extractOptionString(UserMapping, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	opts.UserMapping = emanage.UserMappingType(paramStr)

	// UserMappingUid
	if paramStr, err = extractOptionString(UserMappingUid, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if opts.UserMappingUid, err = strconv.Atoi(paramStr); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	// UserMappingGid
	if paramStr, err = extractOptionString(UserMappingGid, volParams); err != nil {
		return nil, err
	}
	if opts.UserMappingGid, err = strconv.Atoi(paramStr); err != nil {
		return nil, err
	}

	// ExportUid
	if paramStr, err = extractOptionString(ExportUid, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if paramInt, err = strconv.Atoi(paramStr); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	opts.ExportUid = optional.NewInt(paramInt)

	// ExportGid
	if paramStr, err = extractOptionString(ExportGid, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if paramInt, err = strconv.Atoi(paramStr); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	opts.ExportGid = optional.NewInt(paramInt)

	// ExportPermissions
	if paramStr, err = extractOptionString(Permissions, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if opts.ExportPermissions, err = strconv.Atoi(paramStr); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	// DefaultVolumeSize
	capacity := req.GetCapacityRange().GetRequiredBytes()
	if capacity > 0 {
		opts.Capacity = capacity
	} else {
		if paramStr, err = extractOptionString(DefaultVolumeSize, volParams); err != nil {
			return nil, errors.Wrap(err, 0)
		}
		if paramSize, err = size.Parse(paramStr); err != nil {
			return nil, errors.Wrap(err, 0)
		}
		if paramSize > 0 {
			opts.Capacity = int64(paramSize)
		} else {
			opts.Capacity = int64(1 * size.TiB)
		}
	}

	// Access
	if paramStr, err = extractOptionString(Access, volParams); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if paramStr == "" { // Default value
		paramStr = string(emanage.ExportAccessRW)
	}
	opts.Access = paramStr

	glog.V(log.DEBUG).Infof("ecfs: Current volume options: %+v", opts)
	return opts, nil
}
