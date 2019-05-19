package efaas

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	efaasapi "csi-provisioner-elastifile/ecfs/efaas-api"
)

var jsonData = []byte(``)

func TestREST(t *testing.T) {
	t.Parallel()

	res, err := demo1(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}

func TestAPI(t *testing.T) {
	client, err := GetEfaasClient(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	res, err := apiCallGet(client, InstancesURL)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}

func TestSwaggerLowLevelAPI(t *testing.T) {
	client, err := GetEfaasClient(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	apiConf := efaasapi.NewConfiguration()
	apiConf.BasePath = BaseURL
	apiConf.AccessToken = client.GoogleIdToken
	apiConf.Debug = true
	apiConf.DebugFile = "/tmp/api-debug.log"

	// Insecure transport
	defaultTransport := http.DefaultTransport.(*http.Transport)
	apiConf.Transport = &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // TODO: FIXME before deploying to production
	}
	apiConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %v", client.GoogleIdToken))

	res, err := apiConf.APIClient.CallAPI("https://bronze-eagle.gcp.elastifile.com/api/v1/regions", "GET",
		nil, apiConf.DefaultHeader, nil, nil, "", nil)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}
	t.Logf("RES: %+v", res)
}

func TestOpenAPI_CreateInstance(t *testing.T) {
	efaasConf, err := NewEfaasConf(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	snapshots := efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 0.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "10.142.0.0/20", // TODO: Detect the range via K8s OR get it at deploy time
		AccessRights: "readWrite",
	}

	accessors := efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainer{
		Name:        "dc1",                          // Filesystem name
		Description: fmt.Sprintf("Filesystem desc"), // Filesystem description
		QuotaType:   "fixed",                        // Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota.
		HardQuota:   10 * 1024 * 1024 * 1024,        // Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size.
		Snapshot:    snapshots,                      // Snapshot object
		Accessors:   accessors,                      // Defines the access rights to the File System. This is a listof access rights configured by the client for the file system.
	}

	payload := efaasapi.Instances{
		Name:                     "jean-instance1",
		Description:              "eFaaS instance description",
		ServiceClass:             "capacity-optimized-az",
		Region:                   "us-east1",
		Zone:                     "us-east1-b",
		ProvisionedCapacityUnits: 3,
		Network:                  "default",
		Filesystems:              []efaasapi.DataContainer{filesystem},
	}

	op, resp, err := instancesAPI.CreateInstance(ProjectId, payload, "")
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	if resp.StatusCode > http.StatusAccepted {
		t.Fatal("HTTP request failed", "status code", resp.StatusCode, "status", resp.Status)
	}

	t.Logf("Opration: %#v", op)
	t.Logf("Response: %#v", resp)
	t.Logf("Response payload: %v", fmt.Sprint(string(resp.Payload)))

	t.Logf("Waiting for operation id %v ...", op.Id)

	err = WaitForOperationStatusComplete(efaasConf, op.Id, time.Hour)
	if err != nil {
		t.Fatal("WaitForOperationStatusComplete failed", "err", err)
	}
}

func TestOpenAPI_GetInstance(t *testing.T) { // Works (with update of int32 to int64)
	efaasConf, err := NewEfaasConf(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	inst, resp, err := instancesAPI.GetInstance("test-instance--efb4feee-1", ProjectId)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	if resp.StatusCode >= 300 {
		t.Fatal("HTTP request failed", "status code", resp.StatusCode, "status", resp.Status)
	}

	t.Logf("Instance: %#v", inst)
	t.Logf("Response: %#v", resp)
	t.Logf("Response payload: %v", fmt.Sprint(string(resp.Payload)))
}
