package main

import "github.com/golang/glog"

// Note that this is a naive cache implementation and is broken if ControllerServer crashes in-between requests
// Use something like groupcache (probably an overkill) to work around this limitation
var volumeCache map[string]volumeIdType

func cacheVolumeGet(volumeName string) (volumeId volumeIdType, cacheHit bool) {
	volumeId, cacheHit = volumeCache[volumeName]
	return
}

func cacheVolumeAdd(volumeName string, volumeId volumeIdType) {
	if volumeCache == nil {
		volumeCache = make(map[string]volumeIdType)
	}
	volumeCache[volumeName] = volumeId
}

func cacheVolumeRemove(volumeId volumeIdType) {
	for volName, volId := range volumeCache {
		if volId == volumeId {
			delete(volumeCache, volName)
			return
		}
	}
	glog.V(6).Infof("Tried to remove from cache Volume Id that wasn't there - %v", volumeId)
}
