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

package ecfs

import (
	"errors"
	"fmt"
	"strconv"
)

type volumeOptions struct {
	Monitors string `json:"monitors"`
	Pool     string `json:"pool"`
	RootPath string `json:"rootPath"`

	Mounter         string `json:"mounter"`
	ProvisionVolume bool   `json:"provisionVolume"`
}

func validateNonEmptyField(field, fieldName string) error {
	if field == "" {
		return fmt.Errorf("Parameter '%s' cannot be empty", fieldName)
	}

	return nil
}

func (o *volumeOptions) validate() error {
	if err := validateNonEmptyField(o.Monitors, "monitors"); err != nil {
		return err
	}

	if err := validateNonEmptyField(o.RootPath, "rootPath"); err != nil {
		if !o.ProvisionVolume {
			return err
		}
	} else {
		if o.ProvisionVolume {
			return fmt.Errorf("Non-empty field rootPath is in conflict with provisionVolume=true")
		}
	}

	if o.ProvisionVolume {
		if err := validateNonEmptyField(o.Pool, "pool"); err != nil {
			return err
		}
	}

	if o.Mounter != "" {
		if err := validateMounter(o.Mounter); err != nil {
			return err
		}
	}

	return nil
}

func extractOption(dest *string, optionLabel string, options map[string]string) error {
	if opt, ok := options[optionLabel]; !ok {
		return errors.New("[IN PKG] Missing required field " + optionLabel)
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

func newVolumeOptions(volOptions map[string]string) (*volumeOptions, error) {
	var (
		opts                volumeOptions
		provisionVolumeBool string
		err                 error
	)

	if err = extractOption(&opts.Monitors, "monitors", volOptions); err != nil {
		return nil, err
	}

	if err = extractOption(&provisionVolumeBool, "provisionVolume", volOptions); err != nil {
		return nil, err
	}

	if opts.ProvisionVolume, err = strconv.ParseBool(provisionVolumeBool); err != nil {
		return nil, fmt.Errorf("Failed to parse provisionVolume: %v", err)
	}

	if opts.ProvisionVolume {
		if err = extractOption(&opts.Pool, "pool", volOptions); err != nil {
			return nil, err
		}
	} else {
		if err = extractOption(&opts.RootPath, "rootPath", volOptions); err != nil {
			return nil, err
		}
	}

	// This field is optional, don't check for its presence
	extractOption(&opts.Mounter, "mounter", volOptions)

	if err = opts.validate(); err != nil {
		return nil, err
	}

	return &opts, nil
}
