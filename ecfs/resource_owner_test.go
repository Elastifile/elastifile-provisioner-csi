package main

import (
	"os"
	"testing"
)

var ownedResource = NewResourceOwner("blah")

const (
	nodeID1 = "MY_NODE_ID1"
	nodeID2 = "MY_NODE_ID2"
)

func TestResourceOwner_IsAlive(t *testing.T) {
	isAlive := ownedResource.IsAlive()
	AssertEqual(t, isAlive, false)

	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err := ownedResource.TakeOwnership()
	AssertEqual(t, err, nil)

	isAlive = ownedResource.IsAlive()
	AssertEqual(t, isAlive, true)
}

func TestResourceOwner_KeepAlive(t *testing.T) {
	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err := ownedResource.TakeOwnership()
	AssertEqual(t, err, nil)

	// Successful keepalive
	err = ownedResource.KeepAlive()
	AssertEqual(t, err, nil)

	// Different node fails to send keepalive
	_ = os.Setenv(envVarK8sNodeID, nodeID2)
	err = ownedResource.KeepAlive()
	AssertEqual(t, err != nil, true)
}

func TestResourceOwner_TakeOwnership(t *testing.T) {
	// Successful ownership change
	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err := ownedResource.TakeOwnership()
	AssertEqual(t, err, nil)

	// Successful ownership change to the same owner
	_ = os.Setenv(envVarK8sNodeID, nodeID1)
	err = ownedResource.TakeOwnership()
	AssertEqual(t, err, nil)

	// Different node fails to take ownership
	_ = os.Setenv(envVarK8sNodeID, nodeID2)
	err = ownedResource.TakeOwnership()
	AssertEqual(t, err != nil, true)
}
