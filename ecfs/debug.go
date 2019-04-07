package main

import (
	"fmt"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
	"csi-provisioner-elastifile/ecfs/log"
)

const (
	debugConfigMapName      = "csi-debug"
	debugValueCloneDelaySec = "cloneDelaySec"
)

func getDebugValue(key string, defaultValue *string) string {
	value, err := co.GetConfigMapValue(Namespace(), debugConfigMapName, key)
	if err != nil {
		if defaultValue != nil {
			value = *defaultValue
		}

		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get debug value %v from config map %v - using default",
			key, debugConfigMapName), 0)
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
	}

	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Returning debug value for %v: %v", key, value)
	return value
}

func getDebugValueInt(key string, defaultValue *int) (value int) {
	valueStr := getDebugValue(key, nil)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		if defaultValue != nil {
			value = *defaultValue
		}

		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to convert debug value %v from %v to int - using default",
			key, valueStr), 0)
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
	}

	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Returning debug value for %v: %v", key, value)
	return value
}
