package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/elastifile/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
	"csi-provisioner-elastifile/ecfs/log"
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

	// K8s service names' suffixes, appended to APP_NAME, e.g. "elastifile-app"
	k8sServiceNfsVipSuffix     = "-elastifile-svc"
	k8sServiceEmanageVipSuffix = "-emanage-svc"
)

type config struct {
	NFSServer  string `parameter:"nfsServer"`
	EmanageURL string `parameter:"emanageURL"`
	Username   string `parameter:"username"`
	Password   string
}

func (conf *config) String() string {
	return fmt.Sprintf("NFS Server: %v, Management URL: %v, Management username: %v, Management password: %v",
		conf.NFSServer, conf.EmanageURL, conf.Username, strings.Repeat("*", len(conf.Password)))
}

func GetPluginSettings() (configMap map[string]string, secrets map[string][]byte, err error) {
	glog.V(log.DETAILED_INFO).Infof("ecfs: Loading configuration from config map '%v'", configMapName)
	configMap, err = co.GetConfigMap(Namespace(), configMapName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get config map", 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Loading configuration from secrets '%v'", secretsName)
	secrets, err = co.GetSecret(Namespace(), secretsName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get secrets", 0)
	}

	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Provisioner settings - config map: %+v", configMap)

	//TODO: convert secrets map[string][]byte to type Secrets, and add String() with masked password
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Provisioner settings - secrets: %+v", secrets)

	return
}

func pluginConfig() (conf *config, err error) {
	pluginSettings, secret, err := GetPluginSettings()
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	conf = &config{
		NFSServer:  pluginSettings[nfsAddress],
		EmanageURL: pluginSettings[managementAddress],
		Username:   pluginSettings[managementUserName],
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

func updateNfsAddress() (err error) {
	serviceName := fmt.Sprintf("%v%v", AppName(), k8sServiceNfsVipSuffix)
	service, err := co.GetServiceWithRetries(Namespace(), serviceName, 5*time.Minute)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	addr := service.Spec.ClusterIP
	err = co.UpdateConfigMap(Namespace(), configMapName, map[string]string{nfsAddress: addr})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(log.DEBUG).Infof("ecfs: Updated NFS address in config map %v to %v", configMapName, addr)
	return
}

func updateEmanageAddress() (err error) {
	serviceName := fmt.Sprintf("%v%v", AppName(), k8sServiceEmanageVipSuffix)
	service, err := co.GetServiceWithRetries(Namespace(), serviceName, 5*time.Minute)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	addr := service.Spec.ClusterIP
	err = co.UpdateConfigMap(Namespace(), configMapName, map[string]string{managementAddress: addr})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(log.DEBUG).Infof("ecfs: Updated Management address in config map %v to %v", configMapName, addr)
	return
}

func updateConfigEkfs() (err error) {
	if !IsEKFS() {
		glog.V(log.DEBUG).Infof("ecfs: Running outside EKFS - skipping service-based config map updates")
		return
	}

	err = updateNfsAddress()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = updateEmanageAddress()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return
}
