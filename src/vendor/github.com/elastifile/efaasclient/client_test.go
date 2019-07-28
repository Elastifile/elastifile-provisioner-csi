package efaasclient

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-errors/errors"

	"github.com/elastifile/efaasclient/efaasapi"
	"github.com/elastifile/efaasclient/size"
)

const (
	testInstName      = "demo-instance1"
	testFsName        = "test-fs"
	testSnapName      = "test-snap1"
	testShareName     = "e"
	testProjectNumber = "276859139519" // c934
	//testProjectNumber = "602010805072" // golden-eagle-dev-consumer10
	//testProjectNumber    = "507926947502" // elastifile-show
	testEfaasEnvironment = "https://silver-eagle.gcp.elastifile.com"
	//testEfaasEnvironment = "https://bronze-eagle.gcp.elastifile.com"
	testServiceAccountKeyFile = "/tmp/sa-key.json"
)

func testSaKey() (data []byte) {
	data, err := ioutil.ReadFile(testServiceAccountKeyFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read service account key file %v - %v", testServiceAccountKeyFile, err))
	}
	return data
}

func testEfaasApiClient() (client *Client) {
	err := os.Setenv(EnvProjectNumber, testProjectNumber)
	if err != nil {
		panic(fmt.Sprintf("Failed to set env %v to %v. err: %v", EnvProjectNumber, testProjectNumber, err))
	}

	err = os.Setenv(EnvEfaasUrl, testEfaasEnvironment)
	if err != nil {
		panic(fmt.Sprintf("Failed to set env %v to %v. err: %v", EnvEfaasUrl, testProjectNumber, err))
	}

	client, err = NewClient(testSaKey(), EfaasApiUrl())
	if err != nil {
		panic(fmt.Sprintf("Failed to create eFaaS API client %v", err))
	}

	return client
}

func TestClient_CreateInstance(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Creating instance %v", testInstName)
	err := createDefaultInstance(client, testInstName)
	if err != nil {
		t.Fatal("CreateDefaultInstance failed", "err", err)
	}

	t.Logf("Getting instance %v", testInstName)
	inst, err := client.GetInstance(testInstName)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}
	if inst.Name != testInstName {
		t.Fatalf("Instance name (%v) doesn't match the requested one ('%v')", inst.Name, testInstName)
	}
	t.Logf("Instance: %#v", inst)
}

func TestClient_GetInstance(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting instance %v", testInstName)
	inst, err := client.GetInstance(testInstName)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	t.Logf("Instance: %#v", inst)
}

func TestClient_AddFilesystem(t *testing.T) {
	client := testEfaasApiClient()

	fsName := testFsName
	snapshot := &efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 2.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "all",
		AccessRights: "readWrite",
	}

	accessors := &efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainerAdd{
		Name:        fsName,
		HardQuota:   int64(10 * size.GiB),
		QuotaType:   QuotaTypeFixed,
		Description: fmt.Sprintf("Filesystem %v", fsName),
		Accessors:   accessors,
		Snapshot:    snapshot,
	}

	t.Logf("Adding filesystem %v", filesystem.Name)
	err := client.AddFilesystem(testInstName, filesystem, 30*time.Minute)
	if err != nil {
		t.Fatalf("AddFilesystem failed: %v", err)
	}
}

func TestClient_GetFilesystemByName(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting filesystem %v", testFsName)
	fs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	if fs.Name != testFsName {
		t.Fatalf("GetFilesystemByName returned a wrong name: %v vs. %v", fs.Name, testFsName)
	}

}

func TestClient_GetFilesystemById(t *testing.T) {
	client := testEfaasApiClient()

	tmpFs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	t.Logf("Getting filesystem %v", tmpFs.Id)
	fs, err := client.GetFilesystemById(testInstName, tmpFs.Id)
	if err != nil {
		t.Fatalf("GetFilesystemById failed: %v", err)
	}

	if fs.Name != testFsName {
		t.Fatalf("GetFilesystemById returned a wrong name: %v vs. %v", fs.Name, testFsName)
	}
}

func TestClient_UpdateFilesystemQuotaById(t *testing.T) {
	client := testEfaasApiClient()

	fs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName %v failed: %v", testFsName, err)
	}

	quota := 5 * size.GiB
	t.Logf("Updating filesystem %v with quota %v", fs.Id, quota)
	err = client.UpdateFilesystemQuotaById(testInstName, fs.Id, quota, 5*time.Minute)
	if err != nil {
		t.Fatalf("UpdateFilesystemQuotaByName failed. fs: %v quota: %v err: %v", testFsName, quota, err)
	}
}

func TestClient_UpdateFilesystemQuotaByName(t *testing.T) {
	client := testEfaasApiClient()

	quota := 5 * size.GiB
	t.Logf("Updating filesystem %v with quota %v", testFsName, quota)
	err := client.UpdateFilesystemQuotaByName(testInstName, testFsName, quota)
	if err != nil {
		t.Fatalf("UpdateFilesystemQuotaByName failed. fs: %v quota: %v err: %v", testFsName, quota, err)
	}
}

func TestClient_UpdateSnapshotScheduler(t *testing.T) {
	client := testEfaasApiClient()

	fs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	newScheduleValue := "Daily"
	if fs.Snapshots.Schedule == newScheduleValue {
		newScheduleValue = "Monthly"
	}

	schedule := efaasapi.SnapshotSchedule{
		Enable:    !fs.Snapshots.Enable,
		Schedule:  newScheduleValue,
		Retention: fs.Snapshots.Retention + 1,
	}

	t.Logf("Updating snapshot scheduler for filesystem %v", fs.Id)
	err = client.UpdateSnapshotScheduler(testInstName, fs.Id, schedule, 5*time.Minute)
	if err != nil {
		t.Fatalf("UpdateSnapshotScheduler failed - %v", err)
	}

	fs2, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	if fs.Snapshots.Enable == fs2.Snapshots.Enable {
		t.Fatalf("Enable value not updated: %v", fs2.Snapshots.Enable)
	}

	if fs.Snapshots.Retention == fs2.Snapshots.Retention {
		t.Fatalf("Retention value not updated: %v", fs2.Snapshots.Retention)
	}

	if fs.Snapshots.Schedule == fs2.Snapshots.Schedule {
		t.Fatalf("Schedule value not updated: %v", fs2.Snapshots.Schedule)
	}
}

func TestClient_UpdateFilesystemAccessors(t *testing.T) {
	client := testEfaasApiClient()

	fs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "all",
		AccessRights: "readOnly",
	}

	accessors := efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	t.Logf("Updating accessors for filesystem  %v", fs.Id)
	err = client.UpdateFilesystemAccessors(testInstName, fs.Id, accessors, 5*time.Minute)
	if err != nil {
		t.Fatalf("UpdateFilesystemAccessors failed: %v", err)
	}

	fs, err = client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	if fs.Accessors.Items[0].AccessRights != accessor1.AccessRights {
		t.Fatalf("UpdateFilesystemAccessors failed to update AccessRights: %v vs. %v",
			fs.Accessors.Items[0].AccessRights, accessor1.AccessRights)
	}
}

func TestClient_UpdateInstanceCapacity(t *testing.T) {
	client := testEfaasApiClient()

	inst, err := client.GetInstance(testInstName)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	newCapacity := inst.ProvisionedCapacityUnits + 1
	capacity := efaasapi.SetCapacity{
		ProvisionedCapacityUnits: newCapacity,
		CapacityUnitType:         inst.CapacityUnitType,
	}

	t.Logf("Updating instance capacity to %v", capacity)
	err = client.UpdateInstanceCapacity(testInstName, capacity, 30*time.Minute)
	if err != nil {
		t.Fatalf("UpdateInstanceCapacity failed: %v", err)
	}

	inst2, err := client.GetInstance(testInstName)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	if inst2.ProvisionedCapacityUnits != newCapacity {
		t.Fatalf("UpdateInstanceCapacity failed to update capacity: %v vs. %v",
			inst2.ProvisionedCapacityUnits, newCapacity)
	}
}

func TestClient_ListSnapshotsByFsName(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Listing snapshots by filesystem name %v", testFsName)
	snapshots, err := client.ListSnapshotsByFsName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName %v failed: %v", testFsName, err)
	}
	for _, snap := range snapshots {
		t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
	}
}

func TestClient_ListInstanceSnapshots(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Listing instance snapshots")
	snapshots, err := client.ListInstanceSnapshots(testInstName)
	if err != nil {
		t.Fatalf("ListInstanceSnapshots %v failed: %v", testInstName, err)
	}
	for _, snap := range snapshots {
		t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
	}
}

func TestClient_CreateSnapshot(t *testing.T) {
	client := testEfaasApiClient()

	// Create snapshot
	snapshot := efaasapi.Snapshot{
		Name:      testSnapName,
		Retention: 3.0,
	}

	t.Logf("Creating snapshot %v", snapshot.Name)
	err := client.CreateSnapshot(testInstName, testFsName, snapshot, 10*time.Minute)
	if err != nil {
		t.Fatalf("CreateSnapshot failed - %v", err)
	}

	// Verify snapshot creation
	snapshots, err := client.ListSnapshotsByFsName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName failed: %v", err)
	}

	var found bool
	for _, snap := range snapshots {
		if snap.Name == testSnapName {
			t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
			found = true
		}
	}

	if !found {
		t.Fatalf("Snapshot %v not found", testSnapName)
	}
}

func TestClient_GetFilesystemBySnapshotName(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting filesystem by snapshot %v", testSnapName)
	fs, err := client.GetFilesystemBySnapshotName(testInstName, testSnapName)
	if err != nil {
		t.Fatalf("GetFilesystemBySnapshotName %v failed - %v", testSnapName, err)
	}

	if fs.Name != testFsName {
		t.Fatalf("GetFilesystemBySnapshotName returned wrong filesystem: %v vs. %v", fs.Name, testFsName)
	}
}

func TestClient_GetSnapshotByName(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting snapshot %v", testSnapName)
	snap, err := client.GetSnapshotByName(testInstName, testSnapName)
	if err != nil {
		t.Fatalf("GetSnapshotByName failed: %v", err)
	}

	if snap.Name != testSnapName {
		t.Fatalf("GetSnapshotByName returned wrong snapshot: %v vs. %v", snap.Name, testSnapName)
	}
}

func TestClient_GetSnapshotByFsAndName(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting snapshot by filesystem %v and name %v", testFsName, testSnapName)
	snap, err := client.GetSnapshotByFsAndName(testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatalf("GetSnapshotByFsAndName failed: %v", err)
	}
	t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
}

func TestClient_GetSnapshotById(t *testing.T) {
	client := testEfaasApiClient()

	snap, err := client.GetSnapshotByName(testInstName, testSnapName)
	if err != nil {
		t.Fatalf("GetSnapshotByName failed: %v", err)
	}

	t.Logf("Getting snapshot %v", snap.Id)
	snap, err = client.GetSnapshotById(snap.Id)
	if err != nil {
		t.Fatalf("GetSnapshot failed: %v", err)
	}
	t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
}

func TestClient_CreateShareWithFs(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Creating share %v for snapshot %v on filesystem %v", testShareName, testSnapName, testFsName)
	err := client.CreateShareWithFs(testInstName, testFsName, testSnapName, testShareName, 5*time.Minute)
	if err != nil {
		t.Fatalf("CreateShareWithFs failed: %v", err)
	}

	share, err := client.GetShareWithFs(testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatalf("Failed to get share on snapshot %v", err)
	}

	if share.Name != testShareName {
		t.Fatalf("Share %v not found on snapshot %v", testShareName, testSnapName)
	}
}

func TestClient_CreateShare(t *testing.T) {
	client := testEfaasApiClient()

	share, err := client.GetShare(testInstName, testSnapName)
	if err == nil {
		err := client.DeleteShare(testInstName, testFsName, testSnapName, 5*time.Minute)
		if err != nil {
			t.Fatalf("DeleteShare from snapshot %v failed: %v", testSnapName, err)
		}
	}

	t.Logf("Creating share %v on snapshot %v", testShareName, testSnapName)
	err = client.CreateShare(testInstName, testSnapName, testShareName, 5*time.Minute)
	if err != nil {
		t.Fatalf("CreateShareWithFs failed: %v", err)
	}

	share, err = client.GetShare(testInstName, testSnapName)
	if err != nil {
		t.Fatalf("Failed to get share on snapshot %v", err)
	}

	if share.Name != testShareName {
		t.Fatalf("Share %v not found on snapshot %v", testShareName, testSnapName)
	}
}

func TestClient_GetShareWithFs(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting share on snapshot %v, filesystem %v", testSnapName, testFsName)
	share, err := client.GetShareWithFs(testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatalf("GetShareWithFs failed for snapshot %v", testSnapName)
	}

	t.Logf("Share on snapshot %v: %+v", testSnapName, *share)
}

func TestClient_GetShare(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Getting share on snapshot %v", testSnapName)
	share, err := client.GetShare(testInstName, testSnapName)
	if err != nil {
		t.Fatalf("GetShareWithFs failed for snapshot %v", testSnapName)
	}

	t.Logf("Share on snapshot %v: %+v", testSnapName, *share)
}

func TestClient_DeleteShare(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Deleting share on snapshot %v, filesystem %v", testSnapName, testFsName)
	err := client.DeleteShare(testInstName, testFsName, testSnapName, 15*time.Minute)
	if err != nil {
		t.Fatalf("DeleteShare from snapshot %v failed: %v", testSnapName, err)
	}

	share, err := client.GetShareWithFs(testInstName, testFsName, testSnapName)
	if err == nil {
		t.Fatalf("Received snapshot share when it was supposed to have been deleted: %#v", *share)
	}
}

func TestClient_DeleteSnapshot(t *testing.T) {
	client := testEfaasApiClient()

	// Delete snapshot
	t.Logf("Deleting snapshot %v on filesystem %v", testSnapName, testFsName)
	err := client.DeleteSnapshot(testInstName, testFsName, testSnapName, 15*time.Minute)
	if err != nil {
		t.Fatalf("DeleteSnapshot %v failed: %v", testSnapName, err)
	}

	// Verify snapshot has been deleted
	snapshots, err := client.ListSnapshotsByFsName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName failed: %v", err)
	}

	for _, snap := range snapshots {
		if snap.Name == testSnapName {
			t.Fatalf("Snapshot %v found when it should've been deleted", testSnapName)
		}
	}
}

func TestClient_DeleteFilesystem(t *testing.T) {
	client := testEfaasApiClient()

	// Cleanup snapshots
	fs, err := client.GetFilesystemByName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("GetFilesystemByName failed: %v", err)
	}

	schedule := efaasapi.SnapshotSchedule{
		Enable:   false,
		Schedule: "Monthly",
	}

	err = client.UpdateSnapshotScheduler(testInstName, fs.Id, schedule, 5*time.Minute)
	if err != nil {
		t.Fatalf("UpdateSnapshotScheduler failed - %v", err)
	}

	snapshots, err := client.ListSnapshotsByFsName(testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName failed: %v", err)
	}

	for _, snap := range snapshots {
		t.Logf("Deleting snapshot %v", snap.Name)
		err = client.DeleteSnapshot(testInstName, testFsName, snap.Name, 15*time.Minute)
		if err != nil {
			t.Fatalf("DeleteSnapshot %v failed: %v", testSnapName, err)
		}
	}

	// Delete filesystem
	t.Logf("Deleting filesystem %v", fs.Name)
	err = client.DeleteFilesystem(testInstName, testFsName, 15*time.Minute)
	if err != nil {
		t.Fatalf("DeleteFilesystem %v failed: %v", testFsName, err)
	}
}

func TestClient_DeleteInstance(t *testing.T) {
	client := testEfaasApiClient()

	t.Logf("Deleting instance %v", testInstName)
	err := client.DeleteInstance(testInstName, time.Hour)
	if err != nil {
		t.Fatalf("DeleteInstance %v failed: %v", testInstName, err)
	}
}

func createDefaultInstance(client *Client, instanceName string) (err error) {
	snapshot := &efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 2.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "all",
		AccessRights: "readWrite",
	}

	accessors := &efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainer{
		Name:        "dc1",                          // Filesystem name
		Description: fmt.Sprintf("Filesystem desc"), // Filesystem description
		QuotaType:   QuotaTypeFixed,                 // Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota.
		HardQuota:   10 * 1024 * 1024 * 1024,        // Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size.
		Snapshots:   snapshot,                       // Snapshot object
		Accessors:   accessors,                      // Defines the access rights to the File System. This is a list of access rights configured by the client for the file system.
	}

	instance := efaasapi.Instances{
		Name:                     instanceName,
		Description:              "eFaaS instance description",
		ServiceClass:             "tiny-testing", // capacity-optimized
		Region:                   "us-east1",
		Zone:                     "us-east1-b",
		ProvisionedCapacityUnits: 1,
		Network:                  "default",
		Filesystems:              []efaasapi.DataContainer{filesystem},
	}

	err = client.CreateInstance(instance, time.Hour)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}
