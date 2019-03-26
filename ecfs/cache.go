package main

import (
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/log"
)

type CachedVolume struct {
	ID      volumeIdType
	IsReady bool
}

// Note that this is a naive cache implementation and is broken if ControllerServer crashes in-between requests
// Use something like groupcache (probably an overkill) to work around this limitation
var volumeCache map[string]*CachedVolume

func cacheVolumeGet(volumeName string) (cachedVolume *CachedVolume, cacheHit bool) {
	cachedVolume, cacheHit = volumeCache[volumeName]
	return
}

func cacheVolumeSet(volumeName string, volumeId volumeIdType, isReady bool) {
	if volumeCache == nil {
		volumeCache = make(map[string]*CachedVolume)
	}
	volumeCache[volumeName] = &CachedVolume{
		ID:      volumeId,
		IsReady: isReady,
	}
}

func cacheVolumeRemove(volumeId volumeIdType) {
	for volName, volume := range volumeCache {
		if volume.ID == volumeId {
			delete(volumeCache, volName)
			return
		}
	}
	glog.V(log.DEBUG).Infof("Tried to remove from cache Volume Id that wasn't there - %v", volumeId)
}

type CachedSnapshot struct {
	ID     int // ECFS snapshot ID
	Exists bool
}

var snapshotCache map[string]*CachedSnapshot

func cacheSnapshotGet(snapshotName string) (cachedSnapshot *CachedSnapshot, cacheHit bool) {
	cachedSnapshot, cacheHit = snapshotCache[snapshotName]
	return
}

func cacheSnapshotSet(snapshotName string, snapshotId int, exists bool) {
	if snapshotCache == nil {
		snapshotCache = make(map[string]*CachedSnapshot)
	}
	snapshotCache[snapshotName] = &CachedSnapshot{
		ID:     snapshotId,
		Exists: exists,
	}
}

func cacheSnapshotRemoveById(snapshotId int) {
	for volName, snapshot := range snapshotCache {
		if snapshot.ID == snapshotId {
			delete(snapshotCache, volName)
			return
		}
	}
	glog.V(log.DEBUG).Infof("Tried to remove from cache Snapshot Id that wasn't there - %v", snapshotId)
}

func cacheSnapshotRemoveByName(snapshotName string) {
	delete(snapshotCache, snapshotName)
}
