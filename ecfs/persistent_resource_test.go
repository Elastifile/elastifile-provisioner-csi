package main

import (
	"os"
	"testing"
	"time"
)

//var ownedResource, _ = NewPersistentResource(resourceTypeIdVolume, "blah")
var ownedResource PersistentResource

const (
	nodeID1 = "MY_NODE_ID1"
	nodeID2 = "MY_NODE_ID2"
)

func TestPersistentResource_ToMapFromMap(t *testing.T) {
	pr := PersistentResource{
		ResourceType: resourceTypeIdVolume,
		ResourceName: "blah-name",
	}

	t.Logf("data = pr.toMap()")
	data := pr.toMap()
	t.Logf("pr: %+v", pr)
	t.Logf("data: %+v", data)

	pr2 := PersistentResource{}
	data["LastAlive"] = "2006-01-02T15:04:05Z" // RFC3339
	t.Logf("pr2.fromMap() with time")
	err := pr2.fromMap(data)
	AssertEqual(t, err, nil)
	t.Logf("data: %+v", data)
	t.Logf("pr2: %+v", pr2)

	t.Logf("data2 = pr2.toMap()")
	pr2.LastAlive = serializableTime{time.Now()}
	data2 := pr2.toMap()
	t.Logf("pr2: %+v", pr2)
	t.Logf("data2: %+v", data2)
}

func TestSerializableTime_String(t *testing.T) {
	tm1 := serializableTime{time.Now()}
	tmStr := tm1.String()
	t.Logf("tmStr: %v", tmStr)

	tmTmp, err := time.Parse(time.RFC3339, tmStr)
	AssertEqual(t, err, nil)
	tm2 := serializableTime{tmTmp}

	t.Logf("tm1: %v", tm1)
	t.Logf("tm1.String(): %v", tm1.String())
	t.Logf("tm2: %v", tm2)
	t.Logf("tm2.String(): %v", tm2.String())
}

func TestResourceOwner_IsAlive(t *testing.T) {
	isAlive := ownedResource.isAlive()
	AssertEqual(t, isAlive, false)

	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err := ownedResource.takeOwnership()
	AssertEqual(t, err, nil)

	isAlive = ownedResource.isAlive()
	AssertEqual(t, isAlive, true)
}

func TestResourceOwner_KeepAlive(t *testing.T) {
	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err := ownedResource.takeOwnership()
	AssertEqual(t, err, nil)

	// Successful keepalive
	err = ownedResource.KeepAlive()
	AssertEqual(t, err, nil)

	// Different node fails to send keepalive
	_ = os.Setenv(envVarK8sNodeID, nodeID2)
	err = ownedResource.KeepAlive()
	AssertEqual(t, err != nil, true)
}
