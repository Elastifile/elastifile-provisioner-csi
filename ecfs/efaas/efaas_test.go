package efaas

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

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

func TestSwaggerAPI(t *testing.T) {
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
