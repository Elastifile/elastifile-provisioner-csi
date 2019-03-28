package main

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/log"
)

// TODO: Use configMap to make the data available to more than 1 instance

// ResourceOwner is used to indicate the K8s node that owns specific resource
// and when was the last time the CSI plugin on that node reported that it was working on the resource
type ResourceOwner struct {
	resourceID string
	owner      string
	lastAlive  time.Time
}

func NewResourceOwner(resourceID string) ResourceOwner {
	return ResourceOwner{
		resourceID: resourceID,
	}
}

func (ro *ResourceOwner) GetOwner() string {
	// TODO: Read from config map
	return ro.owner
}

// IsAlive returns true if the resource's alive timestamp was updated within the tolerance window
func (ro *ResourceOwner) IsAlive() (isAlive bool) {
	const tolerance = 3 * time.Minute
	LastAliveWithTolerance := ro.lastAlive.Add(tolerance)
	return LastAliveWithTolerance.After(time.Now())
}

// IsOwnedByMe returns true if the current plugin instance is the official resource owner
func (ro *ResourceOwner) IsOwnedByMe() (ownedByMe bool) {
	if ro.GetOwner() == GetPluginNodeName() {
		// TODO: Rename DETAILED_DEBUG to VERBOSE_DEBUG
		glog.V(log.DETAILED_DEBUG).Infof("Resource is owned by the current plugin instance")
		return true
	} else {
		glog.V(log.DETAILED_DEBUG).Infof("Resource is not owned by the current plugin instance %v (owner %v)",
			GetPluginNodeName(), ro.GetOwner())
		return false
	}
}

// KeepAlive updates the resource's lastAlive timestamp
func (ro *ResourceOwner) KeepAlive() (err error) {
	if ro.IsOwnedByMe() {
		ro.lastAlive = time.Now()
		// TODO: Write to config map
	} else {
		err = errors.Errorf("Can't keep alive - resource not owned by current plugin instance (K8s node). "+
			"Current plugin: '%v' Official owner: '%v'", GetPluginNodeName(), ro.GetOwner())
	}
	return
}

// TakeOwnership changes the owner to the current plugin instance
func (ro *ResourceOwner) TakeOwnership() (err error) {
	// TODO: Needs to be atomic
	originalOwner := ro.GetOwner()

	if ro.IsOwnedByMe() {
		ro.KeepAlive()
		return
	}

	if !ro.IsAlive() {
		ro.owner = GetPluginNodeName()
		// TODO: Write to config map

		err = ro.KeepAlive()
		if err != nil {
			return errors.Wrap(err, 0)
		}

		if ro.IsOwnedByMe() {
			glog.V(log.DETAILED_INFO).Infof("Transferred resource ownership to %v (from %v)",
				GetPluginNodeName(), originalOwner)
		} else {
			err = errors.Errorf("Failed to transfer resource ownership to %v (owned by %v)",
				GetPluginNodeName(), ro.GetOwner())
			glog.V(log.DETAILED_DEBUG).Infof(err.Error())
		}
	} else {
		err = errors.Errorf("Resource is NOT owned by the current plugin instance %v (owner %v)",
			GetPluginNodeName(), ro.GetOwner())
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
	}

	return
}
