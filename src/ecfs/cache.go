package main

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"ecfs/log"
)

type CachedVolume struct {
	ID                 string
	Error              error
	IsReady            bool
	persistentResource PersistentResource
}

type VolumeCache map[string]*CachedVolume

var volumeCache VolumeCache // Global volume cache

func (c *VolumeCache) Create(volumeName string) (err error) {
	if *c == nil {
		*c = make(VolumeCache)
	}
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Creating persistent resource for volume %v", volumeName)
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

func (c *VolumeCache) Set(volumeName string, isReady bool, operationFailure error) (err error) {
	if *c == nil {
		glog.V(log.TRACE).Infof("ecfs: Creating persistent resource for volume %v in Set", volumeName)
		err = c.Create(volumeName)
		if err != nil {
			return errors.WrapPrefix(err, fmt.Sprintf(
				"Failed to create volume cache entry in Set. Volume=%v to isReady=%v operationFailure=%v",
				volumeName, isReady, operationFailure), 0)
		}
	}

	(*c)[volumeName].IsReady = isReady
	(*c)[volumeName].Error = operationFailure

	err = (*c)[volumeName].persistentResource.KeepAlive()
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to keep ownership of volume %v", volumeName), 0)
	}

	return
}

func (c *VolumeCache) Remove(volumeId string) (err error) {
	cacheEntry := (*c)[volumeId]
	if cacheEntry == nil {
		glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Cache entry for %v not found - assuming already deleted", volumeId)
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
	// ECFS/EFAAS snapshot ID
	ID      string // Needs to be a string due to the nature of EFAAS snapshot ids
	IsReady bool
}

type SnapshotCache map[string]*CachedSnapshot

var snapshotCache SnapshotCache

func (c *SnapshotCache) Get(snapshotName string) (cachedSnapshot *CachedSnapshot, cacheHit bool) {
	cachedSnapshot, cacheHit = (*c)[snapshotName]
	return
}

func (c *SnapshotCache) Set(snapshotName string, snapshotId string, exists bool) {
	if *c == nil {
		*c = make(map[string]*CachedSnapshot)
	}
	(*c)[snapshotName] = &CachedSnapshot{
		ID:      snapshotId,
		IsReady: exists,
	}
}

func (c *SnapshotCache) RemoveById(snapshotId string) {
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
