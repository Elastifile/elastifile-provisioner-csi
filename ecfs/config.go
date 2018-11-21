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

func newPluginConfig() (conf *config, err error) {
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

	// TODO: Mask the password in the log
	glog.Infof("Parsed config map and secrets %+v", conf)
	return
}
