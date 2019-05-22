package main

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/efaas"
	efaasapi "csi-provisioner-elastifile/ecfs/efaas-api"
	"csi-provisioner-elastifile/ecfs/log"
	"size"
)

//type EfaasWrapper struct {
//	*efaasapi.Configuration
//}

func newEfaasConf() (efaasConf *efaasapi.Configuration) {
	_, secret, err := GetPluginSettings()
	if err != nil {
		panic("Failed to get plugin settings - " + err.Error())
	}

	jsonData := secret[efaasSecretsKeySaJson]
	efaasConf, err = efaas.NewEfaasConf(jsonData)
	if err != nil {
		panic(fmt.Sprintf("Failed to get eFaaS client based on json %v", string(jsonData)))
	}

	return
}

//func (ew *EfaasWrapper) GetClient(jsonData []byte) *efaasapi.Configuration {
//	if ew.Configuration == nil {
//		glog.V(log.DETAILED_INFO).Infof("ecfs: Initializing eFaaS client")
//		conf, err := efaas.NewEfaasConf()
//		if err != nil {
//			panic(fmt.Sprintf("Failed to create eFaaS client. err: %v", err))
//		}
//		ew.Configuration = conf
//		glog.V(log.DEBUG).Infof("ecfs: Initialized new eFaaS client")
//	}
//
//	return ew.Configuration
//}

func efaasGetInstanceName() string {
	return "jean-instance1" // TODO: FIXME - get from os.Getenv or some such
}

func efaasCreateEmptyVolume(volOptions *volumeOptions) (volumeId volumeHandleType, err error) {
	efaasConf := newEfaasConf()
	glog.V(log.DETAILED_INFO).Infof("ecfs: Creating Volume - settings: %+v", volOptions)
	volumeId = volOptions.VolumeId

	snapshot := efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 2.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "all",
		AccessRights: "readWrite",
	}

	accessors := efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainerAdd{
		Name:        string(volumeId),
		HardQuota:   int64(10 * size.GiB),
		QuotaType:   efaas.QuotaTypeFixed,
		Description: fmt.Sprintf("Filesystem %v", volumeId),
		Accessors:   accessors,
		Snapshot:    snapshot,
	}

	// Create Filesystem
	err = efaas.AddFilesystem(efaasConf, efaasGetInstanceName(), filesystem)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Volume %v was already created - assuming it was created "+
				"during previous, failed, attempt", volumeId)
			var e error
			_, e = efaas.GetFilesystemByName(efaasConf, efaasGetInstanceName(), string(volumeId))
			if e != nil {
				logSecondaryError(err, e)
				return
			}
		} else {
			err = errors.Wrap(err, 0)
			return "", errors.Wrap(err, 0)
		}
	}
	//volOptions.DataContainer = fs

	glog.V(log.DEBUG).Infof("ecfs: Created volume with id %v", volumeId)

	return
}

func efaasDeleteVolume(volName volumeHandleType) (err error) {
	efaasConf := newEfaasConf()
	err = efaas.DeleteFilesystem(efaasConf, efaasGetInstanceName(), string(volName))
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DEBUG).Infof("ecfs: Filesystem %v not found - assuming already deleted", volName)
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to delete filesystem %v", volName), 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Deleted filesystem %v", volName)
	return nil
}
