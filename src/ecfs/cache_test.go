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

// ============================================================================

const (
	testVolName string = "testvolname"
	volReady           = true
	volNotReady        = false
)

var emptyVolume *CachedVolume

func TestCacheVolumeGet(t *testing.T) {
	val, hit := volumeCache.Get(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	err := volumeCache.Set(testVolName, volReady, nil)
	AssertEqual(t, err, nil)

	val, hit = volumeCache.Get(testVolName)
	AssertEqual(t, val.ID, testVolName)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)
}

func TestCacheVolumeSet(t *testing.T) {
	err := volumeCache.Set(testVolName, volReady, nil)
	AssertEqual(t, err, nil)

	val, hit := volumeCache.Get(testVolName)
	AssertEqual(t, val.ID, testVolName)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)

	err = volumeCache.Set(testVolName, volReady, nil)
	AssertEqual(t, err, nil)

	val, hit = volumeCache.Get(testVolName)
	AssertEqual(t, val.ID, testVolName)
	AssertEqual(t, val.IsReady, volReady)
	AssertEqual(t, hit, true)

	err = volumeCache.Set(testVolName, volNotReady, nil) // Reset Exists to false
	AssertEqual(t, err, nil)

	val, hit = volumeCache.Get(testVolName)
	AssertEqual(t, val.ID, testVolName)
	AssertEqual(t, val.IsReady, volNotReady)
	AssertEqual(t, hit, true)
}

func TestCacheVolumeRemove(t *testing.T) {
	// Remove non-existent volume
	err := volumeCache.Remove(testVolName)
	AssertEqual(t, err, nil)
	err = volumeCache.Remove(testVolName)
	AssertEqual(t, err, nil)

	val, hit := volumeCache.Get(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	// Remove existing volume
	err = volumeCache.Set(testVolName, volReady, nil)
	AssertEqual(t, err, nil)
	err = volumeCache.Remove(testVolName)
	AssertEqual(t, err, nil)
	val, hit = volumeCache.Get(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	// Remove previously existing and removed volume
	err = volumeCache.Remove(testVolName)
	AssertEqual(t, err, nil)

	val, hit = volumeCache.Get(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)
}

// ============================================================================

const (
	testSnapName string = "testsnapname"
	testSnapId   string    = "111"
	testSnapId2  string    = "222"
	snapReady           = true
	snapNotReady        = false
)

var emptySnapshot *CachedSnapshot

func TestCacheSnapshotGet(t *testing.T) {
	val, hit := snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)

	snapshotCache.Set(testSnapName, testSnapId, snapReady)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val.ID, testSnapId)
	AssertEqual(t, val.IsReady, snapReady)
	AssertEqual(t, hit, true)
}

func TestCacheSnapshotSet(t *testing.T) {
	snapshotCache.Set(testSnapName, testSnapId, snapReady)
	val, hit := snapshotCache.Get(testSnapName)
	AssertEqual(t, val.ID, testSnapId)
	AssertEqual(t, val.IsReady, snapReady)
	AssertEqual(t, hit, true)

	snapshotCache.Set(testSnapName, testSnapId2, snapReady)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val.ID, testSnapId2)
	AssertEqual(t, val.IsReady, snapReady)
	AssertEqual(t, hit, true)

	snapshotCache.Set(testSnapName, testSnapId, snapNotReady) // Reset IsReady to false
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val.ID, testSnapId)
	AssertEqual(t, val.IsReady, snapNotReady)
	AssertEqual(t, hit, true)
}

func TestCacheSnapshotRemoveById(t *testing.T) {
	// Remove non-existent snapshot
	snapshotCache.RemoveById(testSnapId)
	snapshotCache.RemoveById(testSnapId2)
	val, hit := snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)

	// Remove existing snapshot
	snapshotCache.Set(testSnapName, testSnapId, snapReady)
	snapshotCache.RemoveById(testSnapId)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)

	// Remove previously existing and removed snapshot
	snapshotCache.RemoveById(testSnapId)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)
}

func TestCacheSnapshotRemoveByName(t *testing.T) {
	// Remove non-existent snapshot
	snapshotCache.RemoveByName(testSnapName)
	val, hit := snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)

	// Remove existing snapshot
	snapshotCache.Set(testSnapName, testSnapId, snapReady)
	snapshotCache.RemoveByName(testSnapName)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)

	// Remove previously existing and removed snapshot
	snapshotCache.RemoveByName(testSnapName)
	val, hit = snapshotCache.Get(testSnapName)
	AssertEqual(t, val, emptySnapshot)
	AssertEqual(t, hit, false)
}
