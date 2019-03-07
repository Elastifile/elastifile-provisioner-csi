package main

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

const (
	envNamespace = "CSI_NAMESPACE"
	envEKFS      = "EKFS"
	envAppName   = "APP_NAME"

	// Defaults
	defaultNamespace = "default"
	defaultAppName   = "elastifile-app"
)

// IsEKFS checks if we're running in EKFS environment
func IsEKFS() bool {
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

func Namespace() (namespace string) {
	namespace = os.Getenv(envNamespace)
	if namespace == "" {
		namespace = defaultNamespace
		glog.Warningf("Failed getting environment variable %v - falling back to '%v'",
			envNamespace, namespace)
	}
	return
}

func AppName() (name string) {
	name = os.Getenv(envAppName)
	if name == "" {
		name = defaultAppName
		glog.Warningf("Environment variable %v is not set - falling back to '%v'",
			envAppName, name)
	}
	return
}
