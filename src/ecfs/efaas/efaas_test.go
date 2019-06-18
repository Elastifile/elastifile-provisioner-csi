package efaas

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	efaasapi "ecfs/efaas-api"
	"github.com/go-errors/errors"
	"size"
)

const (
	testInstName      = "demo-instance1"
	testFsName        = "test-fs"
	testSnapId        = "12316016938850064433"
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
		panic(fmt.Sprintf("Failed to read service account key file %v - %v", testServiceAccountKeyFile, err.Error()))
	}
	return data
}

func testEfaasConf() (efaasConf *efaasapi.Configuration) {
	err := os.Setenv(envProjectNumber, testProjectNumber)
	if err != nil {
		panic(fmt.Sprintf("Failed to set env %v to %v. err: %v", envProjectNumber, testProjectNumber, err.Error()))
	}

	err = os.Setenv(envEfaasUrl, testEfaasEnvironment)
	if err != nil {
		panic(fmt.Sprintf("Failed to set env %v to %v. err: %v", envEfaasUrl, testProjectNumber, err.Error()))
	}

	efaasConf, err = NewEfaasConf(testSaKey())
	if err != nil {
		panic(fmt.Sprintf("Failed to create NewEfaasConf %v", err.Error()))
	}

	return efaasConf
}

func TestDirectAPI_apiCallGet(t *testing.T) {
	_ = testEfaasConf()
	client, err := GetEfaasClient(testSaKey())
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to get eFaaS client - %v", err.Error()))
	}

	InstancesURL := testEfaasEnvironment + "/api/v2/projects/" + ProjectNumber() + "/instances"
	res, err := apiCallGet(client, InstancesURL)
	if err != nil {
		t.Fatal(fmt.Sprintf("apiCallGet failed - %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}

func TestOpenAPI_CallAPI(t *testing.T) {
	_ = testEfaasConf()
	client, err := GetEfaasClient(testSaKey())
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to get eFaaS client - %v", err.Error()))
	}

	apiConf := efaasapi.NewConfiguration()
	apiConf.BasePath = EfaasApiUrl()
	apiConf.AccessToken = client.GoogleIdToken
	apiConf.Debug = true
	apiConf.DebugFile = "/tmp/api-debug.log"

	defaultTransport := http.DefaultTransport.(*http.Transport)
	apiConf.Transport = &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
	}
	apiConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %v", client.GoogleIdToken))

	res, err := apiConf.APIClient.CallAPI(testEfaasEnvironment+"/api/v2/regions", "GET",
		nil, apiConf.DefaultHeader, nil, nil, "", nil)
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to call API - %v", err.Error()))
	}
	t.Logf("RES: %+v", res)
}

func TestOpenAPI_CreateInstance(t *testing.T) {
	efaasConf := testEfaasConf()

	err := createDefaultInstance(efaasConf, testInstName)
	if err != nil {
		t.Fatal("CreateDefaultInstance failed", "err", err)
	}

	inst, err := GetInstance(efaasConf, testInstName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetInstance failed: %v", err.Error()))
	}
	if inst.Name != testInstName {
		t.Fatal(fmt.Sprintf("Instance name (%v) doesn't match the requested one ('%v')", inst.Name, testInstName))
	}
	t.Logf("Instance: %#v", inst)
}

func TestOpenAPI_GetInstance(t *testing.T) {
	efaasConf := testEfaasConf()

	inst, err := GetInstance(efaasConf, testInstName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetInstance failed: %v", err.Error()))
	}

	t.Logf("Instance: %#v", inst)
}

func TestOpenAPI_UpdateFilesystemQuotaById(t *testing.T) {
	efaasConf := testEfaasConf()

	fs, err := GetFilesystemByName(efaasConf, testInstName, testFsName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetFilesystemByName %v failed: %v", testFsName, err.Error()))
	}

	quota := 5 * size.GiB
	err = UpdateFilesystemQuotaById(efaasConf, testInstName, fs.Id, quota)
	if err != nil {
		t.Fatal(fmt.Sprintf("UpdateFilesystemQuotaByName failed. fs: %v quota: %v err: %v",
			testFsName, quota, err.Error()))
	}
}

func TestOpenAPI_UpdateFilesystemQuotaByName(t *testing.T) {
	efaasConf := testEfaasConf()

	quota := 5 * size.GiB
	err := UpdateFilesystemQuotaByName(efaasConf, testInstName, testFsName, quota)
	if err != nil {
		t.Fatal(fmt.Sprintf("UpdateFilesystemQuotaByName failed. fs: %v quota: %v err: %v",
			testFsName, quota, err.Error()))
	}
}

func TestOpenAPI_AddFilesystem(t *testing.T) {
	efaasConf := testEfaasConf()

	fsName := testFsName
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

	filesystem := efaasapi.DataContainerAdd{
		Name:        fsName,
		HardQuota:   int64(10 * size.GiB),
		QuotaType:   QuotaTypeFixed,
		Description: fmt.Sprintf("Filesystem %v", fsName),
		Accessors:   accessors,
		Snapshot:    snapshot,
	}

	err := AddFilesystem(efaasConf, testInstName, filesystem)
	if err != nil {
		t.Fatal(fmt.Sprintf("AddFilesystem failed: %v", err.Error()))
	}
}

func TestOpenAPI_DeleteFilesystem(t *testing.T) {
	efaasConf := testEfaasConf()

	err := DeleteFilesystem(efaasConf, testInstName, testFsName)
	if err != nil {
		t.Fatal(fmt.Sprintf("DeleteFilesystem %v failed: %v", testFsName, err.Error()))
	}
}

func TestOpenAPI_ListSnapshotsByFsName(t *testing.T) {
	efaasConf := testEfaasConf()

	snapshots, err := ListSnapshotsByFsName(efaasConf, testInstName, testFsName)
	if err != nil {
		t.Fatal(fmt.Sprintf("ListSnapshotsByFsName %v failed: %v", testFsName, err.Error()))
	}
	for _, snap := range snapshots {
		t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
	}
}

func TestOpenAPI_ListInstanceSnapshots(t *testing.T) {
	efaasConf := testEfaasConf()

	snapshots, err := ListInstanceSnapshots(efaasConf, testInstName)
	if err != nil {
		t.Fatal(fmt.Sprintf("ListInstanceSnapshots %v failed: %v", testInstName, err.Error()))
	}
	for _, snap := range snapshots {
		t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
	}
}

func TestOpenAPI_GetSnapshotByName(t *testing.T) {
	efaasConf := testEfaasConf()

	snap, err := GetSnapshotByFsAndName(efaasConf, testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetSnapshot failed: %v", err.Error()))
	}
	t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
}

func TestOpenAPI_GetSnapshotById(t *testing.T) {
	efaasConf := testEfaasConf()

	snap, err := GetSnapshotById(efaasConf, testSnapId)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetSnapshot failed: %v", err.Error()))
	}
	t.Logf("Snap %v (%v): %#v", snap.Id, snap.Name, snap)
}

func TestOpenAPI_CreateSnapshot(t *testing.T) {
	efaasConf := testEfaasConf()

	// Create snapshot
	snapshot := efaasapi.Snapshot{
		Name:      testSnapName,
		Retention: 3.0,
	}
	err := CreateSnapshot(efaasConf, testInstName, testFsName, snapshot)
	if err != nil {
		t.Fatal(fmt.Sprintf("CreateSnapshot failed - %v", err.Error()))
	}

	// Verify snapshot creation
	snapshots, err := ListSnapshotsByFsName(efaasConf, testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName failed: %v", err.Error())
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

func TestOpenAPI_DeleteSnapshot(t *testing.T) {
	efaasConf := testEfaasConf()

	// Delete snapshot
	err := DeleteSnapshot(efaasConf, testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("DeleteSnapshot %v failed: %v", testSnapName, err.Error()))
	}

	// Verify snapshot has been deleted
	snapshots, err := ListSnapshotsByFsName(efaasConf, testInstName, testFsName)
	if err != nil {
		t.Fatalf("ListSnapshotsByFsName failed: %v", err.Error())
	}

	for _, snap := range snapshots {
		if snap.Name == testSnapName {
			t.Fatalf("Snapshot %v found when it should've been deleted", testSnapName)
		}
	}
}

func TestGetShareWithFs(t *testing.T) {
	efaasConf := testEfaasConf()

	share, err := GetShareWithFs(efaasConf, testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetShareWithFs failed for snapshot %v", testSnapName))
	}

	t.Logf("Share on snapshot %v: %+v", testSnapName, *share)
}

func TestGetShare(t *testing.T) {
	efaasConf := testEfaasConf()

	share, err := GetShare(efaasConf, testInstName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("GetShareWithFs failed for snapshot %v", testSnapName))
	}

	t.Logf("Share on snapshot %v: %+v", testSnapName, *share)
}

func TestOpenAPI_CreateShareWithFs(t *testing.T) {
	efaasConf := testEfaasConf()

	err := CreateShareWithFs(efaasConf, testInstName, testFsName, testSnapName, testShareName)
	if err != nil {
		t.Fatal(fmt.Sprintf("CreateShareWithFs failed: %v", err.Error()))
	}

	share, err := GetShareWithFs(efaasConf, testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to get share on snapshot %v", err.Error()))
	}

	if share.Name != testShareName {
		t.Fatal(fmt.Sprintf("Share %v not found on snapshot %v", testShareName, testSnapName))
	}
}

func TestOpenAPI_CreateShare(t *testing.T) {
	efaasConf := testEfaasConf()

	err := CreateShare(efaasConf, testInstName, testSnapName, testShareName)
	if err != nil {
		t.Fatal(fmt.Sprintf("CreateShareWithFs failed: %v", err.Error()))
	}

	share, err := GetShare(efaasConf, testInstName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to get share on snapshot %v", err.Error()))
	}

	if share.Name != testShareName {
		t.Fatal(fmt.Sprintf("Share %v not found on snapshot %v", testShareName, testSnapName))
	}
}

func TestOpenAPI_DeleteShare(t *testing.T) {
	efaasConf := testEfaasConf()

	err := DeleteShare(efaasConf, testInstName, testFsName, testSnapName)
	if err != nil {
		t.Fatal(fmt.Sprintf("DeleteShare from snapshot %v failed: %v", testSnapName, err.Error()))
	}

	share, err := GetShareWithFs(efaasConf, testInstName, testFsName, testSnapName)
	if err == nil {
		t.Fatal(fmt.Sprintf("Received snapshot share when it was supposed to have been deleted: %#v", *share))
	}
}

func createDefaultInstance(efaasConf *efaasapi.Configuration, instanceName string) (err error) {
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
		Name:                     instanceName,
		Description:              "eFaaS instance description",
		ServiceClass:             "tiny-testing", // capacity-optimized
		Region:                   "us-east1",
		Zone:                     "us-east1-b",
		ProvisionedCapacityUnits: 1,
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
