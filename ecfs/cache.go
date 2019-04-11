package main

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/log"
)

type CachedVolume struct {
	ID      string
	IsReady bool
	//IsCopying          bool
	persistentResource PersistentResource
}

type VolumeCache map[string]*CachedVolume

var volumeCache VolumeCache // Global volume cache

func (c *VolumeCache) Create(volumeName string) (err error) {
	if *c == nil {
		*c = make(VolumeCache)
	}
	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Creating persistent resource for volume %v", volumeName)
	resource, err := NewPersistentResource(resourceTypeIdVolume, volumeName)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to create descriptor for volume %v", volumeName), 0)
	}

	(*c)[volumeName] = &CachedVolume{
		IsReady:            false,
		persistentResource: resource,
	}

	return
}

func (c *VolumeCache) Get(volumeName string) (cachedVolume *CachedVolume, cacheHit bool) {
	cachedVolume, cacheHit = (*c)[volumeName]
	return
}

func (c *VolumeCache) Set(volumeName string, isReady bool) (err error) {
	if *c == nil {
		*c = make(VolumeCache)
	}
	err = (*c)[volumeName].persistentResource.KeepAlive()
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to keep ownership of volume %v", volumeName), 0)
	}

	(*c)[volumeName] = &CachedVolume{
		ID:      volumeName,
		IsReady: isReady,
	}

	return
}

func (c *VolumeCache) Remove(volumeId string) (err error) {
	cacheEntry := (*c)[volumeId]
	if cacheEntry == nil {
		glog.V(log.DETAILED_DEBUG).Infof("ecfs: Cache entry for %v not found - assuming already deleted", volumeId)
		return nil
	}

	resource, err := NewPersistentResource(resourceTypeIdVolume, volumeId)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to create descriptor for volume %v", volumeId), 0)
	}

	err = resource.Delete()
	if err != nil {
		glog.Warningf("ecfs: Failed to free up resource ownership information for volume %v", volumeId)
	}
	delete(*c, volumeId)
	return
}

///////////////////////////////////////////////////////////////////////////////

type CachedSnapshot struct {
	ID      int // ECFS snapshot ID
	IsReady bool
}

type SnapshotCache map[string]*CachedSnapshot

var snapshotCache SnapshotCache

func (c *SnapshotCache) Get(snapshotName string) (cachedSnapshot *CachedSnapshot, cacheHit bool) {
	cachedSnapshot, cacheHit = (*c)[snapshotName]
	return
}

func (c *SnapshotCache) Set(snapshotName string, snapshotId int, exists bool) {
	if *c == nil {
		*c = make(map[string]*CachedSnapshot)
	}
	(*c)[snapshotName] = &CachedSnapshot{
		ID:      snapshotId,
		IsReady: exists,
	}
}

func (c *SnapshotCache) RemoveById(snapshotId int) {
	for volName, snapshot := range *c {
		if snapshot.ID == snapshotId {
			delete(*c, volName)
			return
		}
	}
	glog.V(log.DEBUG).Infof("Tried to remove from cache Snapshot Id that wasn't there - %v", snapshotId)
}

func (c *SnapshotCache) RemoveByName(snapshotName string) {
	delete(*c, snapshotName)
}
