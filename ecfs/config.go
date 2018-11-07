package main

import (
	"github.com/golang/glog"

	"github.com/elastifile/errors"
)

type config struct {
	NFSServer  string
	EmanageURL string
	Username   string
	Password   string
}

func newConfig() (conf *config, err error) {
	// TODO: Implement fetching cluster configuration and credentials from configmap+secret - these are fake values
	configMap, secret, err := GetProvisionerSettings()
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	// TODO: Check key availability
	conf = &config{
		NFSServer:  configMap[nfsAddress],
		EmanageURL: configMap[managementAddress],
		Username:   configMap[managementUserName],
		Password:   secret[managementPassword],
	}

	// TODO: Consider masking the password
	glog.Infof("Parsed config map and secrets %+v", conf)
	return
}
