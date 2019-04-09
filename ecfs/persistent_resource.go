package main

// IMPORTANT: Current consistency mechanism relies on the clocks being (at least roughly) in sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"csi-provisioner-elastifile/ecfs/co"
	"csi-provisioner-elastifile/ecfs/log"
)

type resourceTypeEnum int

const (
	TTL = 3 * time.Minute

	resourceTypeIdVolume resourceTypeEnum = iota + 1
	resourceTypeIdSnapshot
)

var (
	resourceTypeName = map[resourceTypeEnum]string{
		resourceTypeIdVolume:   "volume",
		resourceTypeIdSnapshot: "snapshot",
	}
)

func (rte resourceTypeEnum) String() string {
	return resourceTypeName[rte]
}

type serializableTime struct { // Custom struct is needed to change the default stringer
	time.Time
}

func (st *serializableTime) String() string {
	return st.Format(time.RFC3339)
}

// PersistentResource is used to indicate the K8s node that owns specific resource
// and when was the last time the CSI plugin on that node reported that it was working on the resource
type PersistentResource struct {
	// Static values
	ResourceName string           `json:",omitempty"`
	ResourceType resourceTypeEnum `json:",omitempty"`

	// Dynamic values
	OwnedBy   string
	LastAlive serializableTime
}

func NewPersistentResource(resourceType resourceTypeEnum, resourceName string) (pr PersistentResource, err error) {
	pr = PersistentResource{
		ResourceType: resourceType,
		ResourceName: resourceName,
	}

	glog.V(log.DEBUG).Infof("ecfs: Creating new persistent resource %v", pr.resourceKey())
	err = co.CreateConfigMap(Namespace(), pr.resourceKey(), pr.toMap())
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Config map %v already exists, loading contents", pr.resourceKey())
			e := pr.loadFromPersistentStore()
			if e != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to load existing config map: %v", e), 0)
			}
		} else {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create config map %v", pr.resourceKey()), 0)
			return
		}
	}

	err = pr.KeepAlive()
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	glog.V(log.DETAILED_DEBUG).Infof("ecfs: Created new persistent resource %+v", pr)
	return
}

// resourceKey creates resource identifier to be used in persistent storage
func (pr *PersistentResource) resourceKey() string {
	return fmt.Sprintf("%v-%v", resourceTypeName[pr.ResourceType], pr.ResourceName)
}

// loadFromPersistentStore updates the PersistentResource from config map
func (pr *PersistentResource) loadFromPersistentStore() error {
	LastAliveWithPartialTTL := pr.LastAlive.Add(TTL / 3)
	if LastAliveWithPartialTTL.After(time.Now()) {
		glog.V(log.DETAILED_DEBUG).Infof("ecfs: Using cached persistent resource for %v %v",
			resourceTypeName[pr.ResourceType], pr.ResourceName)
		return nil
	}

	confMapName := pr.resourceKey()
	glog.V(log.DETAILED_INFO).Infof("ecfs: Loading persistent resource configuration from config map '%v'", confMapName)

	data, err := co.GetConfigMap(Namespace(), confMapName)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Config map %v not found - assuming new resource", confMapName)
		} else {
			return errors.WrapPrefix(err, fmt.Sprintf("Failed to get config map %v", confMapName), 0)
		}
	}

	err = pr.fromMap(data)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to load resource from map - %+v", data), 0)
	}

	glog.V(log.DETAILED_DEBUG).Infof("AAAAA ecfs: Loaded persistent resource from config map %v: %+v", pr.resourceKey(), pr) // TODO: DELME
	return nil
}

func (pr *PersistentResource) GetOwner() string {
	err := pr.loadFromPersistentStore()
	if err != nil {
		panic(err.Error())
	}
	return pr.OwnedBy
}

// isAlive returns true if the resource's alive timestamp was updated within the TTL window
func (pr *PersistentResource) isAlive() (isAlive bool) {
	err := pr.loadFromPersistentStore()
	if err != nil {
		panic(err.Error())
	}
	LastAliveWithTTL := pr.LastAlive.Add(TTL)
	return LastAliveWithTTL.After(time.Now())
}

// isOwnedByMe returns true if the current plugin instance is the official resource owner
// Note: Liveness is not checked
func (pr *PersistentResource) isOwnedByMe() (ownedByMe bool) {
	owner := pr.GetOwner()
	if owner == GetPluginNodeName() {
		// TODO: Rename DETAILED_DEBUG to VERBOSE_DEBUG
		glog.V(log.DETAILED_DEBUG).Infof("ecfs: Resource is owned by the current plugin instance")
		return true
	} else {
		if owner == "" {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Resource has no owner")
		} else {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Resource is not owned by the current plugin instance %v (owner %v)",
				GetPluginNodeName(), owner)
		}
		return false
	}
}

const jsonRoot = "data"

// toMap converts the receiver's contents to a map
func (pr *PersistentResource) toMap() map[string]string {
	bytes, err := json.Marshal(pr)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal struct %+v", pr))
	}

	return map[string]string{
		jsonRoot: string(bytes),
	}
}

// fromMap populates the receiver with data from the map
func (pr *PersistentResource) fromMap(data map[string]string) error {

	// Convert map to json
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to marshal map %+v", data), 0)
	}

	// Convert json to struct
	err = json.Unmarshal(bytes, pr)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to unmarshal json %v", string(bytes)), 0)
	}
	return nil
}

func (pr *PersistentResource) updatePersistentConf() (err error) {
	confMapName := pr.resourceKey()
	data := pr.toMap()
	glog.V(log.DETAILED_DEBUG).Infof("AAAAA ecfs: Updating config map. id: %v data: %+v", pr.resourceKey(), data) // TODO: DELME

	err = co.UpdateConfigMap(Namespace(), confMapName, data)
	if err != nil {
		if isErrorDoesNotExist(err) {
			err = co.CreateConfigMap(Namespace(), confMapName, data)
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create config map %v", confMapName), 0)
			} else {
				glog.V(log.DEBUG).Infof("ecfs: Created config map %v with %+v", confMapName, data)
				return
			}
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to update config map %v with %+v", confMapName, data), 0)
	}

	glog.V(log.DEBUG).Infof("ecfs: Updated config map %v with %+v", confMapName, data)
	return
}

// TakeOwnership changes the owner to the current plugin instance
func (pr *PersistentResource) takeOwnership() (err error) {
	originalOwner := pr.GetOwner()
	if pr.isAlive() && !pr.isOwnedByMe() {
		err = errors.Errorf("Can't take ownership - resource is owned by another active plugin instance - %v",
			originalOwner)
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
		return
	}

	// TODO: Protect by a lock - consider using a lock file on ELFS
	// Take ownership
	pr.OwnedBy = GetPluginNodeName()
	pr.LastAlive = serializableTime{time.Now()}
	err = pr.updatePersistentConf()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	// Re-check
	// It's not bulletproof, but in practical terms it reduces the chances of contention virtually to zero,
	// since K8s doesn't flood the plugin with duplicate requests
	if pr.isAlive() && pr.isOwnedByMe() {
		glog.V(log.DETAILED_INFO).Infof("ecfs: Transferred resource ownership to %v (from %v)",
			GetPluginNodeName(), originalOwner)
	} else {
		err = errors.Errorf("Failed to transfer resource ownership to %v (owned by %v)",
			GetPluginNodeName(), pr.GetOwner())
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
		return
	}

	return
}

// KeepAlive updates the resource's LastAlive timestamp
func (pr *PersistentResource) KeepAlive() (err error) {
	glog.V(log.TRACE).Infof("KeepAlive %v", pr.resourceKey())
	err = pr.takeOwnership()
	if err != nil {
		err = errors.WrapPrefix(err, "Can't keep alive", 0)
		glog.V(log.DETAILED_DEBUG).Infof(err.Error())
		return
	}
	return
}

func (pr *PersistentResource) KeepAliveRoutine(errChan chan error, stopChan chan struct{}, timeout time.Duration) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			glog.V(log.DETAILED_DEBUG).Infof("AAAAA KeepAliveRoutine %v", pr.resourceKey()) // TODO: DELME
			err := pr.KeepAlive()
			if err != nil {
				err = errors.WrapPrefix(err, fmt.Sprintf("ecfs: Failed to send KeepAlive for %v", pr.resourceKey()), 0)
				glog.Warning(err.Error())
			}
		case <-timer.C:
			errChan <- errors.Errorf("Timed out sending KeepAlive for %v - aborting routine", pr.resourceKey())
			return
		case <-stopChan:
			errChan <- nil
			return
		}
	}
}

// Delete removes persistent data used by the resource
func (pr *PersistentResource) Delete() (err error) {
	err = pr.loadFromPersistentStore()
	if err != nil {
		glog.V(log.DEBUG).Infof("Failed to load persistent resource %v - assuming already deleted", pr.resourceKey())
		return nil
	}

	err = co.DeleteConfigMap(Namespace(), pr.resourceKey())
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DETAILED_DEBUG).Infof("ecfs: Config map %v not found - assuming success", pr.resourceKey())
			err = nil
		} else {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete persistent resource: %v", pr.resourceKey()), 0)
		}
	}
	return
}

// TODO: Add background keepalive thread per resource being actively worked on

// TODO: Decide what to do with existing DCs whose volumes were deleted (e.g. due to existing data/snapshots)
// Deleting all snapshots and their exports makes sense, but existing data is something that needs PM's decision
