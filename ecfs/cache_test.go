package main

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("Values mismatch - %v (type %v) != %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	debug.PrintStack()
}

func AssertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		return
	}
	t.Errorf("Values are equal - %v (type %v) == %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	debug.Stack()
}

//func TestGetSnapshotByName(t *testing.T) {
//	var snapshotName = "vs-111-222"
//	var ems emanageClient
//
//	err := fakeEmsConfig()
//	if err != nil {
//		t.Fatal("GetSnapshotByName failed: ", err)
//	}
//
//	snapshot, err := ems.GetSnapshotByName(snapshotName)
//	if err != nil {
//		t.Fatal("GetSnapshotByName failed: ", err)
//	}
//
//	t.Log("TestGetSnapshotByName", "snapshot.Name", snapshot.Name, "snapshot.ID", snapshot.ID, "snapshot", *snapshot)
//}

const (
	testVolName string       = "testvolname"
	testVolId   volumeIdType = "testvolid"
	testVolId2  volumeIdType = "testvolid2"
	emptyVolume volumeIdType = ""
)

func TestCacheVolumeGet(t *testing.T) {
	val, hit := cacheVolumeGet(testVolName)
	AssertEqual(t, val, emptyVolume)
	AssertEqual(t, hit, false)

	cacheVolumeAdd(testVolName, testVolId)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val, testVolId)
	AssertEqual(t, hit, true)
}

func TestCacheVolumeAdd(t *testing.T) {
	cacheVolumeAdd(testVolName, testVolId)
	val, hit := cacheVolumeGet(testVolName)
	AssertEqual(t, val, testVolId)
	AssertEqual(t, hit, true)

	cacheVolumeAdd(testVolName, testVolId2)
	val, hit = cacheVolumeGet(testVolName)
	AssertEqual(t, val, testVolId2)
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
	cacheVolumeAdd(testVolName, testVolId)
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
