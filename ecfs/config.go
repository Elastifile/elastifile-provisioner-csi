package main

import (
	"fmt"
	"strings"

	"github.com/elastifile/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
)

const (
	//SecretNamespace = "csiProvisionerSecretNamespace"
	configMapName = "elastifile"
	secretsName   = "elastifile"

	// Config map / secret keys
	managementAddress  = "managementAddress"
	managementUserName = "managementUserName"
	managementPassword = "password"
	nfsAddress         = "nfsAddress"
)

type config struct {
	NFSServer  string `parameter:"nfsServer"`
	EmanageURL string `parameter:"emanageURL"`
	Username   string `parameter:"username"`
	Password   string
	//SecretName      string `parameter:"secretName"`
	//SecretNamespace string `parameter:"secretNamespace"`
}

func (conf *config) String() string {
	return fmt.Sprintf("NFS Server: %v, Management URL: %v, Management username: %v, Management password: %v",
		conf.NFSServer, conf.EmanageURL, conf.Username, strings.Repeat("*", len(conf.Password)))
}

func GetProvisionerSettings() (configMap map[string]string, secrets map[string][]byte, err error) {
	namespace := "default"

	glog.V(5).Infof("ecfs: Loading configuration from config map '%v'", configMapName)
	configMap, err = co.GetConfigMap(namespace, configMapName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get config map", 0)
	}

	glog.V(5).Infof("ecfs: Loading configuration from secrets '%v'", secretsName)
	secrets, err = co.GetSecret(namespace, secretsName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get secrets", 0)
	}

	glog.V(10).Infof("Provisioner settings - config map: %+v", configMap)
	glog.V(10).Infof("Provisioner settings - secrets: %+v", secrets)

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

	const tlsPrefix = "https://"
	if !strings.HasPrefix(strings.ToLower(conf.EmanageURL), tlsPrefix) {
		glog.Warningf("ECFS management URL has to start with https:// - got %v", conf.EmanageURL)
		conf.EmanageURL = tlsPrefix + conf.EmanageURL
	}

	glog.Infof("ecfs: Parsed config map and secrets: %+v", conf)
	return
}
