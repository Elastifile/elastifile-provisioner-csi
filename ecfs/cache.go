package main

import (
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/log"
)

type CachedVolume struct {
	Owner     ResourceOwner
	ID        volumeIdType
	IsReady   bool
	IsCopying bool
}

type VolumeCache map[string]*CachedVolume

// Note that this is a naive cache implementation and is broken if ControllerServer crashes in-between requests
// Use something like groupcache (probably an overkill) to work around this limitation
var volumeCache VolumeCache

func (c *VolumeCache) Get(volumeName string) (cachedVolume *CachedVolume, cacheHit bool) {
	cachedVolume, cacheHit = (*c)[volumeName]
	return
}

// TODO: Wrap Set with CreateStarted/CreateEnded, add timestamps
func (c *VolumeCache) Set(volumeName string, volumeId volumeIdType, isReady bool) {
	if *c == nil {
		*c = make(VolumeCache)
	}
	(*c)[volumeName] = &CachedVolume{
		ID:      volumeId,
		IsReady: isReady,
	}
}

func (c *VolumeCache) Remove(volumeId volumeIdType) {
	for volName, volume := range *c {
		if volume.ID == volumeId {
			delete(*c, volName)
			return
		}
	}
	glog.V(log.DEBUG).Infof("Tried to remove from cache Volume Id that wasn't there - %v", volumeId)
}

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
