package efaas

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"size"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	efaasapi "csi-provisioner-elastifile/ecfs/efaas-api"
	"csi-provisioner-elastifile/ecfs/log"
)

type EfaasClient struct {
	*http.Client
	GoogleIdToken string
}

const (
	EfaasOperationStatusPending = "PENDING"
	EfaasOperationStatusRunning = "RUNNING"
	EfaasOperationStatusDone    = "DONE"
)

const (
	QuotaTypeFixed = "fixed"
	QuotaTypeAuto  = "auto"
)

func GetHttpClient() (client *http.Client, err error) {
	defaultTransport := http.DefaultTransport.(*http.Transport)

	// Create new Transport that ignores self-signed SSL
	httpTransportWithSelfSignedTLS := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	// TODO: IMPORTANT: Do not deploy in production while InsecureSkipVerify is in use
	client = &http.Client{Transport: httpTransportWithSelfSignedTLS}

	return client, nil
}

// Deprecated
func apiCallGet(client *EfaasClient, reqURL string) (body []byte, err error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create request: %v", req), 0)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", client.GoogleIdToken))

	glog.Infof("Req Header: %#v", req.Header)

	res, err := client.Do(req)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to send request: %v", req), 0)
		return
	}

	// Read the response body
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to read response body"), 0)
		return
	}

	return
}

func GetEfaasClient(data []byte) (client *EfaasClient, err error) {
	httpClient, err := GetHttpClient()
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get HTTP client"), 0)
		return
	}

	token, err := GetEfaasToken(data)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS client"), 0)
		return
	}

	client = &EfaasClient{
		Client:        httpClient,
		GoogleIdToken: token,
	}
	return
}

func NewEfaasConf(jsonData []byte) (efaasConf *efaasapi.Configuration, err error) {
	client, err := GetEfaasClient(jsonData)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get eFaaS client", 0)
		return
	}

	efaasConf = efaasapi.NewConfiguration()
	efaasConf.BasePath = BaseURL
	efaasConf.AccessToken = client.GoogleIdToken
	efaasConf.Debug = true
	efaasConf.DebugFile = "/tmp/api-debug.log"

	// Insecure transport
	defaultTransport := http.DefaultTransport.(*http.Transport)
	efaasConf.Transport = &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // TODO: FIXME before deploying to production
	}
	efaasConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %v", client.GoogleIdToken))
	return
}

func CheckApiCall(err error, resp *efaasapi.APIResponse, op *efaasapi.Operation) error {
	var summary string

	if op != nil {
		if len(op.Error_.Errors) > 0 {
			summary = fmt.Sprintf("Operation %v (%v) failed - %#v", op.Name, op.Id, op.Error_.Errors)
		}
	}

	if resp != nil {
		if resp.StatusCode >= http.StatusAccepted {
			if summary == "" {
				summary = "API call failed" // Generic error
			}
			summary = fmt.Sprintf("%v - HTTP code %v (%v). Details: %v",
				summary, resp.StatusCode, resp.Status, string(resp.Payload))
		}
	}

	if err != nil {
		return errors.WrapPrefix(err, summary, 0)
	} else if summary != "" {
		return errors.New(summary)
	}

	return nil
}

func GetOperation(efaasConf *efaasapi.Configuration, id string) (operation *efaasapi.Operation, err error) {
	api := efaasapi.ProjectsprojectoperationApi{Configuration: efaasConf}
	op, resp, err := api.GetOperation(id, ProjectId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get operation by id %v project %v", id, ProjectId), 0)
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

func WaitForOperationStatus(efaasConf *efaasapi.Configuration, id string, status string, timeout time.Duration) (err error) {
	glog.V(log.DEBUG).Infof("Waiting for operation %v to reach status %v", id, status)
	var (
		e  error
		op *efaasapi.Operation
	)
	for startTime := time.Now(); time.Since(startTime) <= timeout; time.Sleep(time.Second) {
		op, e = GetOperation(efaasConf, id)
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

func WaitForOperationStatusComplete(efaasConf *efaasapi.Configuration, id string, timeout time.Duration) (err error) {
	err = WaitForOperationStatus(efaasConf, id, EfaasOperationStatusDone, timeout)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func GetFilesystemById(efaasConf *efaasapi.Configuration, instanceName string, fsId string) (filesystem efaasapi.DataContainer, err error) {
	inst, err := GetInstance(efaasConf, instanceName)
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

func GetFilesystemByName(efaasConf *efaasapi.Configuration, instanceName string, fsName string) (filesystem efaasapi.DataContainer, err error) {
	inst, err := GetInstance(efaasConf, instanceName)
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

func CreateInstance(efaasConf *efaasapi.Configuration, instanceName string, instance efaasapi.Instances) (err error) {
	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	op, resp, err := instancesAPI.CreateInstance(ProjectId, instance, "")
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create instance %v", instanceName), 0)
		return
	}

	glog.V(log.DETAILED_INFO).Infof("Waiting for instance %v to be created...", instanceName)
	err = WaitForOperationStatusComplete(efaasConf, op.Id, time.Hour)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to create instance %v", instanceName), 0)
		return
	}

	return
}

func CreateDefaultInstance(efaasConf *efaasapi.Configuration, instanceName string) (err error) {
	snapshot := efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 2.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "all",
		AccessRights: "readWrite",
	}

	accessors := efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainer{
		Name:        "dc1",                          // Filesystem name
		Description: fmt.Sprintf("Filesystem desc"), // Filesystem description
		QuotaType:   QuotaTypeFixed,                 // Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota.
		HardQuota:   10 * 1024 * 1024 * 1024,        // Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size.
		Snapshots:   snapshot,                       // Snapshot object
		Accessors:   accessors,                      // Defines the access rights to the File System. This is a listof access rights configured by the client for the file system.
	}

	instance := efaasapi.Instances{
		Name:        instanceName,
		Description: "eFaaS instance description",
		//ServiceClass:             "capacity-optimized-az",
		ServiceClass:             "capacity-optimized",
		Region:                   "us-east1",
		Zone:                     "us-east1-b",
		ProvisionedCapacityUnits: 3,
		Network:                  "default",
		Filesystems:              []efaasapi.DataContainer{filesystem},
	}

	err = CreateInstance(efaasConf, instanceName, instance)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func GetInstance(efaasConf *efaasapi.Configuration, instanceName string) (inst *efaasapi.Instances, err error) {
	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}
	inst, resp, err := instancesAPI.GetInstance(instanceName, ProjectId)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("GetInstance %v failed", instanceName), 0)
		return
	}

	return
}

const minQuota = int64(10 * size.GiB)

func GetFsByName(efaasConf *efaasapi.Configuration, instanceName string, fsName string) (filesystem *efaasapi.DataContainer, err error) {
	inst, err := GetInstance(efaasConf, instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for _, fs := range inst.Filesystems {
		if fs.Name == fsName {
			filesystem = &fs
			return
		}
	}

	err = errors.Errorf("Filesystem %v not found on instance %v", fsName, instanceName)
	return
}

func UpdateFilesystemQuotaById(efaasConf *efaasapi.Configuration, instanceName string, fsId string, quota size.Size) (err error) {
	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	payload := efaasapi.UpdateQuota{
		HardQuota: int64(quota),
		QuotaType: QuotaTypeFixed,
	}

	if payload.HardQuota < minQuota {
		glog.Warningf("Requested volume size (%v) is smaller than the minimal size - using %v", quota, minQuota)
		payload.HardQuota = minQuota
	}

	glog.V(log.INFO).Infof("Updating filesystem: %#v with %#v", fsId, payload)
	op, resp, err := instancesAPI.UpdateFilesystemQuota(instanceName, fsId, ProjectId, payload)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed updating filesystem %v quota %#v", fsId, payload), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 5*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v quota %#v update",
			fsId, payload), 0)
		return
	}

	return
}

func UpdateFilesystemQuotaByName(efaasConf *efaasapi.Configuration, instanceName string, fsName string, quota size.Size) (err error) {
	fs, err := GetFsByName(efaasConf, instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name: %v", fsName), 0)
		return
	}

	err = UpdateFilesystemQuotaById(efaasConf, instanceName, fs.Id, quota)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func AddFilesystem(efaasConf *efaasapi.Configuration, instanceName string, filesystem efaasapi.DataContainerAdd) (err error) {
	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}
	inst, err := GetInstance(efaasConf, instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	glog.V(log.HIGH_LEVEL_INFO).Infof("Adding filesystem: %#v: %#v", filesystem.Name, filesystem)
	op, resp, err := instancesAPI.AddFilesystem(inst.Name, ProjectId, filesystem)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to add filesystem %v - %#v", filesystem.Name, filesystem), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 5*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting to create filesystem %v", filesystem.Name), 0)
		return
	}

	return
}

func DeleteFilesystem(efaasConf *efaasapi.Configuration, instanceName string, fsName string) (err error) {
	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}
	inst, err := GetInstance(efaasConf, instanceName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	fs, err := GetFilesystemByName(efaasConf, inst.Name, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	glog.V(log.HIGH_LEVEL_INFO).Infof("Deleting filesystem: %v", fsName)

	op, resp, err := instancesAPI.DeleteFilesystem(inst.Name, fs.Id, ProjectId)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete filesystem %v", fsName), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 10*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed waiting for filesystem %v delete operation %v",
			fsName, op.Id), 0)
		return
	}

	return
}

func GetSnapshotById(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapId string) (snapshot *efaasapi.Snapshots, err error) {
	snapshotsAPI := efaasapi.ProjectsprojectsnapshotsApi{Configuration: efaasConf}
	snapshot, resp, err := snapshotsAPI.GetSnapshot(ProjectId, snapId)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by Id %v. Filesystem name: %v",
			snapId, fsName), 0)
		return
	}

	return
}

func GetSnapshotByName(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapName string) (snapshot *efaasapi.Snapshots, err error) {
	snapshots, err := ListSnapshots(efaasConf, instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots in filesystem %v", fsName), 0)
		return
	}

	for _, snap := range snapshots {
		if snap.Name == snapName {
			snapshot = &snap
			return
		}
	}

	err = errors.Errorf("Snapshot name %v not found in filesystem %v", snapName, fsName)
	return
}

func ListSnapshots(efaasConf *efaasapi.Configuration, instanceName string, fsName string) (snapshots []efaasapi.Snapshots, err error) {
	fs, err := GetFilesystemByName(efaasConf, instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	glog.V(log.HIGH_LEVEL_INFO).Infof("Listing snapshots for filesystem %v", fsName)
	snapshotsAPI := efaasapi.ProjectsprojectsnapshotsApi{Configuration: efaasConf}
	snapshots, resp, err := snapshotsAPI.ListSnapshots(ProjectId, fs.Id, instanceName)
	err = CheckApiCall(err, resp, nil)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to list snapshots for filesystem %v", fsName), 0)
		return
	}

	return
}

func CreateSnapshot(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapshot efaasapi.Snapshot) (err error) {
	fs, err := GetFilesystemByName(efaasConf, instanceName, fsName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get filesystem by name %v", fsName), 0)
		return
	}

	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}
	requestId := "" // Used for idempotency
	op, resp, err := instancesAPI.FilesystemManualCreateSnapshot(instanceName, fs.Id, ProjectId, snapshot, requestId)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create snapshot %v on filesystem %v", snapshot.Name, fsName), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 5*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v creation failed", snapshot.Name), 0)
		return
	}

	return
}

func DeleteSnapshot(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapName string) (err error) {
	snap, err := GetSnapshotByName(efaasConf, instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	snapshotsAPI := efaasapi.ProjectsprojectsnapshotsApi{Configuration: efaasConf}
	requestId := "" // Used for idempotency
	op, resp, err := snapshotsAPI.DeleteSnapshot(ProjectId, snap.Id, requestId)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete snapshot %v (%v) from filesystem %v",
			snap.Name, snap.Id, fsName), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 10*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v (%v) deletion failed",
			snap.Name, snap.Id), 0)
		return
	}

	return
}

func CreateShare(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapName string, shareName string) (err error) {
	payload := efaasapi.SnapshotShareCreate{
		// Create Share
		ShareName: shareName,
	}

	snap, err := GetSnapshotByName(efaasConf, instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	snapshotsAPI := efaasapi.ProjectsprojectsnapshotsApi{Configuration: efaasConf}
	requestId := "" // Used for idempotency
	op, resp, err := snapshotsAPI.CreateShare(ProjectId, snap.Id, payload, requestId)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create share on snapshot %v (%v)", snapName, snap.Id), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 5*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for share %v creation on snapshot %v failed",
			shareName, snapName), 0)
		return
	}

	return
}

func DeleteShare(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapName string) (err error) {
	snap, err := GetSnapshotByName(efaasConf, instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	snapshotsAPI := efaasapi.ProjectsprojectsnapshotsApi{Configuration: efaasConf}
	requestId := "" // Used for idempotency
	op, resp, err := snapshotsAPI.DeleteShare(ProjectId, snap.Id, requestId)
	err = CheckApiCall(err, resp, op)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete share on snapshot %v (%v)", snapName, snap.Id), 0)
		return
	}

	err = WaitForOperationStatusComplete(efaasConf, op.Id, 5*time.Minute)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Waiting for snapshot %v (%v) deletion failed",
			snap.Name, snap.Id), 0)
		return
	}

	return
}

func GetShare(efaasConf *efaasapi.Configuration, instanceName string, fsName string, snapName string) (share *efaasapi.Share, err error) {
	snap, err := GetSnapshotByName(efaasConf, instanceName, fsName, snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			snapName, fsName), 0)
		return
	}

	share = &snap.Share
	if share.Name == "" {
		err = errors.Errorf("No shares found on snapshot %v on filesystem %v", snapName, fsName)
		return
	}

	return
}
