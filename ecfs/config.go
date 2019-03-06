package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastifile/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
)

const (
	// Environment variable names
	envNamespace = "CSI_NAMESPACE"
	envEKFS      = "EKFS"

	//SecretNamespace = "csiProvisionerSecretNamespace"
	configMapName = "elastifile"
	secretsName   = "elastifile"

	// Config map / secret keys
	managementAddress  = "managementAddress"
	managementUserName = "managementUserName"
	managementPassword = "password"
	nfsAddress         = "nfsAddress"

	// K8s service names
	k8sServiceNfsVip     = "elastifile-app-elastifile-svc"
	k8sServiceEmanageVip = "elastifile-app-emanage-svc"

	defaultNamespace = "default"
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

func Namespace() (namespace string) {
	namespace = os.Getenv(envNamespace)
	if namespace == "" {
		namespace = defaultNamespace
		glog.Warningf("Failed getting environment variable %v - falling back to the default value '%v'",
			envNamespace, namespace)
	}
	return
}

func GetProvisionerSettings() (configMap map[string]string, secrets map[string][]byte, err error) {
	glog.V(5).Infof("ecfs: Loading configuration from config map '%v'", configMapName)
	configMap, err = co.GetConfigMap(Namespace(), configMapName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get config map", 0)
	}

	glog.V(5).Infof("ecfs: Loading configuration from secrets '%v'", secretsName)
	secrets, err = co.GetSecret(Namespace(), secretsName)
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

// isEKFS checks if we're running in EKFS environment
func isEKFS() bool {
	isEkfsStr := os.Getenv(envEKFS)
	if isEkfsStr == "" {
		return false
	}

	isEkfs, err := strconv.ParseBool(isEkfsStr)
	if err != nil {
		glog.Warningf("Failed to parse environment variable %v's value (%v) as bool - assuming running in EKFS",
			envEKFS, isEkfsStr)
		return true
	}

	return isEkfs
}

func updateNfsAddress() (err error) {
	if !isEKFS() {
		glog.V(6).Infof("ecfs: Running outside EKFS - skipping service-based update of ECFS NFS address")
		return
	}
	service, err := co.GetServiceWithRetries(Namespace(), k8sServiceNfsVip, 5*time.Minute)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	addr := service.Spec.ClusterIP
	err = co.UpdateConfigMap(Namespace(), configMapName, map[string]string{nfsAddress: addr})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(6).Infof("ecfs: Updated NFS address in config map %v to %v", configMapName, addr)
	return
}

func updateEmanageAddress() (err error) {
	service, err := co.GetServiceWithRetries(Namespace(), k8sServiceEmanageVip, 5*time.Minute)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	addr := service.Spec.ClusterIP
	err = co.UpdateConfigMap(Namespace(), configMapName, map[string]string{managementAddress: addr})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	glog.V(6).Infof("ecfs: Updated Management address in config map %v to %v", configMapName, addr)
	return
}

func updateConfigEkfs() (err error) {
	if !isEKFS() {
		glog.V(6).Infof("ecfs: Running outside EKFS - skipping service-based config map updates")
		return
	}

	err = updateNfsAddress()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	_ = updateEmanageAddress()
	// TODO: Uncomment the following block once EKFS adds the corresponding service
	//err = updateEmanageAddress()
	//if err != nil {
	//	return errors.Wrap(err, 0)
	//}

	return
}
