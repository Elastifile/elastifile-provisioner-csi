package efaasclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"github.com/elastifile/efaasclient/efaasapi"
	"github.com/elastifile/efaasclient/log"
	"github.com/elastifile/efaasclient/size"
)

const (
	EfaasOperationStatusPending = "PENDING"
	EfaasOperationStatusRunning = "RUNNING"
	EfaasOperationStatusDone    = "DONE"
)

const (
	QuotaTypeFixed = "fixed"
	QuotaTypeAuto  = "auto"
)

const (
	CapacityUnitTypeSteps = "Steps"
	CapacityUnitTypeBytes = "Bytes"
)

const minQuota = int64(10 * size.GiB)

type Client struct {
	efaasapi.APIClient
	Context context.Context
}

func getEfaasConf(jsonData []byte) (efaasConf *efaasapi.Configuration, err error) {
	efaasToken, err := GetEfaasToken(jsonData)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get eFaaS client", 0)
		return
	}

	efaasConf = efaasapi.NewConfiguration()
	efaasConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %v", efaasToken))

	return
}

func NewClient(saKeyJson []byte, apiBaseURL string) (client *Client, err error) {
	conf, err := getEfaasConf(saKeyJson)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get eFaaS configuration", 0)
		return
	}

	if apiBaseURL != "" {
		conf.BasePath = apiBaseURL
	} else {
		conf.BasePath = EfaasApiUrl()
	}

	apiClient := efaasapi.NewAPIClient(conf)

	client = &Client{
		APIClient: *apiClient,
		Context:   context.Background(),
	}

	return
}

func CheckApiCall(err error, resp *http.Response, op *efaasapi.Operation) error {
	var summary string

	if op != nil && op.Error_ != nil {
		if len(op.Error_.Errors) > 0 {
			summary = fmt.Sprintf("Operation %v (%v) failed - %#v", op.Name, op.Id, op.Error_.Errors)
		}
	}

	if resp != nil {
		if resp.StatusCode >= http.StatusAccepted {
			if summary == "" {
				summary = "API call failed" // Generic error
			}
			summary = fmt.Sprintf("%v - HTTP code %v (%v)", summary, resp.StatusCode, resp.Status)
		}
	}

	if err != nil {
		swaggerErr, ok := err.(efaasapi.GenericSwaggerError)
		if ok {
			summary = fmt.Sprintf("%v - %v", summary, string(swaggerErr.Body()))
		}
		return errors.WrapPrefix(err, summary, 0)
	} else if summary != "" {
		return errors.New(summary)
	}

	return nil
}

func (c *Client) GetOperation(id string) (operation efaasapi.Operation, err error) {
	op, resp, err := c.ProjectsprojectoperationApi.GetOperation(c.Context, id, ProjectNumber())
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get operation by id %v project %v",
			id, ProjectNumber()), 0)
		return
	}

	if resp.StatusCode >= http.StatusAccepted {
		err = errors.Errorf("HTTP request failed - %v", resp.Status)
		return
	}

	if len(op.Error_.Errors) > 0 {
		err = errors.Errorf("Operation %v (%v) failed - %#v", op.Name, op.Id, op.Error_.Errors)
		return
	}

	return op, nil
}

func (c *Client) WaitForOperationStatus(id string, status string, timeout time.Duration) (err error) {
	glog.V(log.DEBUG).Infof("Waiting for operation %v to reach status %v", id, status)
	var (
		e  error
		op efaasapi.Operation
	)
	for startTime := time.Now(); time.Since(startTime) <= timeout; time.Sleep(time.Second) {
		op, e = c.GetOperation(id)
		if e != nil {
			glog.V(log.VERBOSE_DEBUG).Infof("GetOperation failed - retrying... %v", e.Error())
		} else if op.Status == status {
			glog.V(log.DEBUG).Infof("Operation %v reached status %v", id, op.Status)
			return nil
		}
	}

	message := fmt.Sprintf("Timed out waiting for operation %v to reach status %v after %v.", id, status, timeout)
	if e != nil {
		message += " Last error: " + e.Error()
	}

	return errors.New(message)
}

func (c *Client) WaitForOperationStatusComplete(id string, timeout time.Duration) (err error) {
	err = c.WaitForOperationStatus(id, EfaasOperationStatusDone, timeout)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func (c *Client) GetFilesystemById(instanceName string, fsId string) (filesystem efaasapi.DataContainer, err error) {
	inst, err := c.GetInstance(instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, fs := range inst.Filesystems {
		if fs.Id == fsId {
			filesystem = fs
			return
		}
	}
	err = errors.Errorf("Filesystem %v not found in instance %v", fsId, instanceName)
	return
}

func (c *Client) GetFilesystemBySnapshotName(instanceName string, snapName string) (filesystem efaasapi.DataContainer, err error) {
	snap, err := c.GetSnapshotByName(instanceName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot %v", snapName), 0)
		return
	}

	filesystem, err = c.GetFilesystemById(instanceName, snap.FilesystemId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem %v", snap.FilesystemName), 0)
		return
	}

	return
}

func (c *Client) GetFilesystemByName(instanceName string, fsName string) (filesystem efaasapi.DataContainer, err error) {
	inst, err := c.GetInstance(instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, fs := range inst.Filesystems {
		if fs.Name == fsName {
			filesystem = fs
			return
		}
	}
	err = errors.Errorf("Filesystem %v not found in instance %v", fsName, instanceName)
	return
}

func (c *Client) CreateInstance(instance efaasapi.Instances, timeout time.Duration) (err error) {
	op, resp, err := c.ProjectsprojectinstancesApi.CreateInstance(c.Context, ProjectNumber(), instance, nil)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create instance %v", instance.Name), 0)
		return
	}

	glog.V(log.DETAILED_INFO).Infof("Waiting for instance %v to be created...", instance.Name)
	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to create instance %v", instance.Name), 0)
		return
	}

	return
}

func (c *Client) DeleteInstance(instanceName string, timeout time.Duration) (err error) {
	opts := &efaasapi.DeleteInstanceItemOpts{
		Force: optional.NewString("false"),
	}
	op, resp, err := c.ProjectsprojectinstancesApi.DeleteInstanceItem(c.Context, instanceName, ProjectNumber(), opts)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Delete Instance %v failed", instanceName), 0)
		return
	}

	glog.V(log.DETAILED_INFO).Infof("Waiting for instance %v to be created...", instanceName)
	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to create instance %v", instanceName), 0)
		return
	}

	return
}

func (c *Client) GetInstance(instanceName string) (inst efaasapi.Instances, err error) {
	inst, resp, err := c.ProjectsprojectinstancesApi.GetInstance(c.Context, instanceName, ProjectNumber())
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("GetInstance %v failed", instanceName), 0)
		return
	}

	return
}

func (c *Client) UpdateFilesystemQuotaById(instanceName string, fsId string, quota size.Size, timeout time.Duration) (err error) {
	payload := efaasapi.UpdateQuota{
		HardQuota: int64(quota),
		QuotaType: QuotaTypeFixed,
	}

	if payload.HardQuota < minQuota {
		glog.Warningf("Requested volume size (%v) is smaller than the minimal size - using %v", quota, minQuota)
		payload.HardQuota = minQuota
	}

	glog.V(log.INFO).Infof("Updating filesystem: %#v with %#v", fsId, payload)
	op, resp, err := c.ProjectsprojectinstancesApi.UpdateFilesystemQuota(c.Context, instanceName, fsId, ProjectNumber(), payload)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed updating filesystem %v quota %#v", fsId, payload), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v quota %#v update",
			fsId, payload), 0)
		return
	}

	return
}

func (c *Client) UpdateFilesystemQuotaByName(instanceName string, fsName string, quota size.Size) (err error) {
	fs, err := c.GetFilesystemByName(instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name: %v", fsName), 0)
		return
	}

	err = c.UpdateFilesystemQuotaById(instanceName, fs.Id, quota, 5*time.Minute)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func (c *Client) UpdateInstanceCapacity(instanceName string, capacity efaasapi.SetCapacity, timeout time.Duration) (err error) {
	op, resp, err := c.ProjectsprojectinstancesApi.PostInstanceSetCapacity(c.Context, instanceName, ProjectNumber(), capacity)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to add capacity to instance %v", instanceName), 0)
		return
	}
	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to add capacity to instace %v", instanceName), 0)
		return
	}

	return
}

func (c *Client) UpdateFilesystemAccessors(instanceName string, fsId string, accessors efaasapi.Accessors, timeout time.Duration) (err error) {
	glog.V(log.INFO).Infof("Updating filesystem: %#v with %#v", fsId, accessors)
	op, resp, err := c.ProjectsprojectinstancesApi.SetAccessorsToFilesystem(c.Context, instanceName, fsId, ProjectNumber(), accessors)

	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed updating filesystem %v accessors %#v", fsId, accessors), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v accessors %#v update",
			fsId, accessors), 0)
		return
	}
	return
}

func (c *Client) UpdateSnapshotScheduler(instanceName string, fsId string, scheduler efaasapi.SnapshotSchedule, timeout time.Duration) (err error) {
	glog.V(log.INFO).Infof("Updating filesystem: %#v with %#v", fsId, scheduler)
	op, resp, err := c.ProjectsprojectinstancesApi.SetFilesystemSnapshotScheduling(
		c.Context, instanceName, fsId, ProjectNumber(), scheduler)

	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed updating filesystem %v scheduler %#v", fsId, scheduler), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v scheduler %#v update",
			fsId, scheduler), 0)
		return
	}
	return
}

func (c *Client) AddFilesystem(instanceName string, filesystem efaasapi.DataContainerAdd, timeout time.Duration) (err error) {
	inst, err := c.GetInstance(instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	glog.V(log.HIGH_LEVEL_INFO).Infof("Adding filesystem: %#v: %#v", filesystem.Name, filesystem)
	op, resp, err := c.ProjectsprojectinstancesApi.AddFilesystem(c.Context, inst.Name, ProjectNumber(), filesystem)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to add filesystem %v - %#v", filesystem.Name, filesystem), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to create filesystem %v", filesystem.Name), 0)
		return
	}

	return
}

func (c *Client) DeleteFilesystem(instanceName string, fsName string, timeout time.Duration) (err error) {
	inst, err := c.GetInstance(instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	fs, err := c.GetFilesystemByName(inst.Name, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	glog.V(log.HIGH_LEVEL_INFO).Infof("Deleting filesystem: %v", fsName)

	op, resp, err := c.ProjectsprojectinstancesApi.DeleteFilesystem(c.Context, inst.Name, fs.Id, ProjectNumber())
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete filesystem %v", fsName), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v delete operation %v",
			fsName, op.Id), 0)
		return
	}

	return
}

func (c *Client) GetSnapshotByFsAndName(instanceName string, fsName string, snapName string) (snapshot efaasapi.Snapshots, err error) {
	snapshots, err := c.ListSnapshotsByFsName(instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots in filesystem %v", fsName), 0)
		return
	}

	for _, snap := range snapshots {
		if snap.Name == snapName {
			snapshot = snap
			return
		}
	}

	err = errors.Errorf("Snapshot name %v not found in filesystem %v", snapName, fsName)
	return
}

func (c *Client) ListInstanceSnapshots(instanceName string) (snapshots []efaasapi.Snapshots, err error) {
	glog.V(log.DEBUG).Infof("Listing all snapshots")

	snapshots, resp, err := c.ProjectsprojectsnapshotsApi.ListInstanceSnapshots(c.Context, ProjectNumber(), instanceName, nil)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list instance snapshots"), 0)
		return
	}

	return
}

func (c *Client) GetSnapshotByName(instanceName string, snapName string) (snapshot efaasapi.Snapshots, err error) {
	snapshots, err := c.ListInstanceSnapshots(instanceName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list all snapshots"), 0)
		return
	}

	for _, snap := range snapshots {
		if snap.Name == snapName {
			snapshot = snap
			return
		}
	}

	err = errors.Errorf("Snapshot name %v not found", snapName)
	return
}

func (c *Client) GetSnapshotById(snapId string) (snapshot efaasapi.Snapshots, err error) {
	glog.V(log.DEBUG).Infof("Getting snapshot by Id %v", snapId)
	snapshot, resp, err := c.ProjectsprojectsnapshotsApi.GetSnapshot(c.Context, ProjectNumber(), snapId)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshots byt Id %v", snapId), 0)
		return
	}

	return
}

func (c *Client) ListSnapshotsByFsName(instanceName string, fsName string) (snapshots []efaasapi.Snapshots, err error) {
	fs, err := c.GetFilesystemByName(instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	glog.V(log.DEBUG).Infof("Listing snapshots for filesystem %v", fsName)
	listOpts := &efaasapi.ListSnapshotsOpts{
		Filesystem: optional.NewString(fs.Id),
		Instance:   optional.NewString(instanceName),
	}
	snapshots, resp, err := c.ProjectsprojectsnapshotsApi.ListSnapshots(c.Context, ProjectNumber(), listOpts)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots for filesystem %v", fsName), 0)
		return
	}

	return
}

func (c *Client) CreateSnapshot(instanceName string, fsName string, snapshot efaasapi.Snapshot, timeout time.Duration) (err error) {
	fs, err := c.GetFilesystemByName(instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	createOpts := &efaasapi.FilesystemManualCreateSnapshotOpts{}
	op, resp, err := c.ProjectsprojectinstancesApi.FilesystemManualCreateSnapshot(
		c.Context, instanceName, fs.Id, ProjectNumber(), snapshot, createOpts)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create snapshot %v on filesystem %v", snapshot.Name, fsName), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v creation failed", snapshot.Name), 0)
		return
	}

	return
}

func (c *Client) DeleteSnapshot(instanceName string, fsName string, snapName string, timeout time.Duration) (err error) {
	snap, err := c.GetSnapshotByFsAndName(instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	deleteOpts := &efaasapi.DeleteSnapshotOpts{}
	op, resp, err := c.ProjectsprojectsnapshotsApi.DeleteSnapshot(c.Context, ProjectNumber(), snap.Id, deleteOpts)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete snapshot %v (%v) from filesystem %v",
			snap.Name, snap.Id, fsName), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v (%v) deletion failed",
			snap.Name, snap.Id), 0)
		return
	}

	return
}

func (c *Client) CreateShareWithFs(instanceName string, fsName string, snapName string, shareName string, timeout time.Duration) (err error) {
	payload := efaasapi.SnapshotShareCreate{
		// Create Share
		ShareName: shareName,
	}

	snap, err := c.GetSnapshotByFsAndName(instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	createOpts := &efaasapi.CreateShareOpts{}
	op, resp, err := c.ProjectsprojectsnapshotsApi.CreateShare(c.Context, ProjectNumber(), snap.Id, payload, createOpts)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create share on snapshot %v (%v)", snapName, snap.Id), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for share %v creation on snapshot %v failed",
			shareName, snapName), 0)
		return
	}

	return
}

func (c *Client) CreateShare(instanceName string, snapName string, shareName string, timeout time.Duration) (err error) {
	payload := efaasapi.SnapshotShareCreate{
		// Create Share
		ShareName: shareName,
	}

	snap, err := c.GetSnapshotByName(instanceName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", snapName), 0)
		return
	}

	createOpts := &efaasapi.CreateShareOpts{}
	op, resp, err := c.ProjectsprojectsnapshotsApi.CreateShare(c.Context, ProjectNumber(), snap.Id, payload, createOpts)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create share on snapshot %v (%v)", snapName, snap.Id), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for share %v creation on snapshot %v failed",
			shareName, snapName), 0)
		return
	}

	return
}

func (c *Client) DeleteShare(instanceName string, fsName string, snapName string, timeout time.Duration) (err error) {
	snap, err := c.GetSnapshotByFsAndName(instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	deleteOpts := &efaasapi.DeleteShareOpts{}
	op, resp, err := c.ProjectsprojectsnapshotsApi.DeleteShare(c.Context, ProjectNumber(), snap.Id, deleteOpts)
	err = CheckApiCall(err, resp, &op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete share on snapshot %v (%v)", snapName, snap.Id), 0)
		return
	}

	err = c.WaitForOperationStatusComplete(op.Id, timeout)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v (%v) deletion failed",
			snap.Name, snap.Id), 0)
		return
	}

	return
}

func (c *Client) GetShareWithFs(instanceName string, fsName string, snapName string) (share *efaasapi.Share, err error) {
	snap, err := c.GetSnapshotByName(instanceName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	share = snap.Share
	if share == nil {
		err = errors.Errorf("Got nil share in snapshot %v filesystem %v", snapName, fsName)
		return
	}
	if share.Name == "" {
		err = errors.Errorf("No shares found on snapshot %v on filesystem %v", snapName, fsName)
		return
	}

	return
}

func (c *Client) GetShare(instanceName string, snapName string) (share *efaasapi.Share, err error) {
	snap, err := c.GetSnapshotByName(instanceName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v", snapName), 0)
		return
	}

	share = snap.Share
	if share == nil {
		err = errors.Errorf("Got nil share in snapshot %v", snapName)
		return
	}
	if share.Name == "" {
		err = errors.Errorf("No shares found on snapshot %v", snapName)
		return
	}

	return
}
