package main

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

const (
	// Environment variables are set at deployment time via plugin container's manifest
	envNamespace         = "CSI_NAMESPACE"
	envVarK8sNodeID      = "NODE_ID"
	envAppName           = "APP_NAME"
	envEKFS              = "EKFS"
	envEfaasInstance     = "CSI_EFAAS_INSTANCE"
	envProjectNumber     = "GCP_PROJECT_NUMBER"
	EnvProjectNumberPrev = "CSI_GCP_PROJECT_NUMBER" // Deprecated, kept for backwards compatibility
	envEfaasUrl          = "EFAAS_URL"

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

// IsEFAAS checks if eFaaS instance is configured
func IsEFAAS() bool {
	if os.Getenv(envEfaasInstance) == "" {
		return false
	}
	return true
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
