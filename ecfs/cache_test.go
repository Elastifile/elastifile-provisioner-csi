package main

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Logf("Comparing %v and %v", a, b)

	if a == nil && b == nil { // Can't compare nils directly - this will fail in case they're pointers to different types
		return
	}

	if a == b {
		return
	}

	t.Errorf("Values mismatch - %v (type %v) != %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	debug.PrintStack()
}

const (
	testVolName string       = "testvolname"
	testVolId   volumeIdType = "testvolid"
	testVolId2  volumeIdType = "testvolid2"
	volReady                 = true
	volNotReady              = false
)

var emptyVolume *CachedVolume

func TestCacheVolumeGet(t *testing.T) {
	val, hit := cacheVolumeGet(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	cacheVolumeSet(testVolName, testVolId, volReady)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val.ID, testVolId)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)
}

func TestCacheVolumeSet(t *testing.T) {
	cacheVolumeSet(testVolName, testVolId, volReady)
	val, hit := cacheVolumeGet(testVolName)
	AssertEqual(t, val.ID, testVolId)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)

	cacheVolumeSet(testVolName, testVolId2, volReady)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val.ID, testVolId2)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)

	cacheVolumeSet(testVolName, testVolId, volNotReady) // Reset Exists to false
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val.ID, testVolId)
	AssertEqual(t, val.IsReady, volNotReady)
	AssertEqual(t, hit, true)
}

func TestCacheVolumeRemove(t *testing.T) {
	// Remove non-existent volume
	cacheVolumeRemove(volumeIdType(testVolId))
	cacheVolumeRemove(volumeIdType(testVolId2))
	val, hit := cacheVolumeGet(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	// Remove existing volume
	cacheVolumeSet(testVolName, testVolId, volReady)
	cacheVolumeRemove(testVolId)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	// Remove previously existing and removed volume
	cacheVolumeRemove(testVolId)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)
}
