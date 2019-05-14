package efaas

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/golang/glog"
)

type EfaasClient struct {
	*http.Client
	GoogleIdToken string
}

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
