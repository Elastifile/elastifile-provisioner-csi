package efaas

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
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

func GetOperation(efaasConf *efaasapi.Configuration, id string) (operation efaasapi.Operation, err error) {
	api := efaasapi.ProjectsprojectoperationApi{Configuration: efaasConf}
	ops, resp, err := api.GetOperation(id, ProjectId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get operation by id %v project %v", id, ProjectId), 0)
		return
	}

	if resp.StatusCode >= 300 {
		err = errors.Errorf("HTTP request failed - %v", resp.Status)
		return
	}

	if len(ops) > 1 {
		glog.Warningf("Bad ops count - %v for operation %v (%v)", len(ops), ops[0].Name, ops[0].Id)
	}

	op := ops[0]
	if len(op.Error_.Errors) > 0 {
		err = errors.Errorf("Operation %v (%v) failed - %#v", op.Name, op.Id, op.Error_.Errors)
		return
	}

	return op, nil
}

func WaitForOperationStatus(efaasConf *efaasapi.Configuration, id string, status string, timeout time.Duration) (err error) {
	for startTime := time.Now(); time.Since(startTime) <= timeout; time.Sleep(time.Second) {
		op, e := GetOperation(efaasConf, id)
		if e != nil {
			glog.V(log.DEBUG)
		}
		if op.Status == status {
			return nil
		}
	}

	return
}

func WaitForOperationStatusComplete(efaasConf *efaasapi.Configuration, id string, timeout time.Duration) (err error) {
	err = WaitForOperationStatus(efaasConf, id, EfaasOperationStatusDone, timeout)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

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
