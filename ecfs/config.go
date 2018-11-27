package main

import (
	"github.com/elastifile/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
)

const (
	//SecretNamespace = "csiProvisionerSecretNamespace"
	configMapName = "elastifile"
	secretsName   = "elastifile"
)

type config struct {
	NFSServer  string `parameter:"nfsServer"`
	EmanageURL string `parameter:"emanageURL"`
	Username   string `parameter:"username"`
	Password   string
	namespace  string
	//SecretName      string `parameter:"secretName"`
	//SecretNamespace string `parameter:"secretNamespace"`
}

func GetProvisionerSettings() (configMap map[string]string, secrets map[string][]byte, err error) {
	namespace := "default" // TODO: Support user-specified (or automatically detected) namespace

	configMap, err = co.GetConfigMap(namespace, configMapName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get config map", 0)
	}

	secrets, err = co.GetSecret(namespace, secretsName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get secrets", 0)
	}

	glog.Infof("CCCCC GetProvisionerSettings - config map: %+v", configMap) // TODO: DELME
	glog.Infof("CCCCC GetProvisionerSettings - secrets: %+v", secrets)      // TODO: DELME

	return
}

func pluginConfig() (conf *config, err error) {
	configMap, secret, err := GetProvisionerSettings()
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	// TODO: Check key availability
	conf = &config{
		NFSServer:  configMap[nfsAddress],
		EmanageURL: configMap[managementAddress],
		Username:   configMap[managementUserName],
		Password:   string(secret[managementPassword]),
	}

	// TODO: Make sure conf.EmanageURL starts with https

	// TODO: Mask the password in the log using Stringer
	glog.Infof("Parsed config map and secrets %+v", conf)
	return
}
