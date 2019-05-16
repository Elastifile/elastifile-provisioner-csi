package efaas

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/golang/glog"
	"golang.org/x/oauth2/google"
)

var (
	BaseURL      = "https://bronze-eagle.gcp.elastifile.com/api/v1"
	RegionsURL   = "https://bronze-eagle.gcp.elastifile.com/api/v1/regions"
	InstancesURL = "https://bronze-eagle.gcp.elastifile.com/api/v1/projects/276859139519/instances" // 934 = 276859139519
)

const (
	GoogleAuthURL = "https://www.googleapis.com/oauth2/v4/token"
)

type googleIdTokenResp struct {
	IdToken          string `json:"id_token,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func GetGoogleIdToken(token string) (id string, err error) {
	values := url.Values{}
	values.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	values.Add("assertion", token)
	client := &http.Client{}

	encodedValues := values.Encode()
	glog.V(1).Infof("Using encoded values %v", encodedValues)
	req, err := http.NewRequest("POST", GoogleAuthURL, strings.NewReader(encodedValues))
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create request to %+v", GoogleAuthURL), 0)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to read response to %+v", req), 0)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to read response to %+v", req), 0)
		return
	}

	var tokenResp = &googleIdTokenResp{}
	err = json.Unmarshal(body, tokenResp)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to unmarshal Google Id Token response %v", string(body)), 0)
		return
	}
	if tokenResp.Error != "" {
		err = errors.Errorf("Failed to get Google Id Token: %+v", tokenResp)
		return
	}

	glog.V(1).Infof("Google body: %v", string(body))
	return tokenResp.IdToken, nil
}

func publicFromPrivateKey(privateKey *rsa.PrivateKey) (pubKeyBase64 string, err error) {
	publicKey := privateKey.Public()
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to marshal public key"), 0)
		return
	}
	pubKeyBase64 = base64.StdEncoding.EncodeToString(pubKeyBytes)
	return
}

// GetEfassClaims returns JWT claims based on service account's json, which can be obtained via cloud console
func GetEfaasClaims(data []byte, scope string) (claims jwt.MapClaims, err error) {
	tokenExpiration := time.Hour

	jwtConf, err := google.JWTConfigFromJSON(data, scope)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get JWTConfigFromJSON"), 0)
		return
	}

	timestamp := time.Now()

	claims = jwt.MapClaims{
		// Issued At
		"iat": timestamp.Unix(),

		// Token expiration expires after one hour
		"exp": timestamp.Add(tokenExpiration).Unix(),

		// iss (Issuer) is the service account email

		"iss": jwtConf.Email, // "efaas-csi@elastifile-gce-lab-c934.iam.gserviceaccount.com"

		// target_audience is the URL of the target service
		"target_audience": "563209362155-dmktm1rt2snprao3te1a5gf0tk9l39i8.apps.googleusercontent.com", // eFaaS

		// aud must be Google token endpoints URL
		"aud": GoogleAuthURL,
	}

	return
}

func GetPrivateKeyFromJSON(data []byte) (privateKey *rsa.PrivateKey, err error) {
	conf, err := google.JWTConfigFromJSON(data, "")
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get JWTConfigFromJSON"), 0)
		return
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(conf.PrivateKey)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to parse private key from PEM: %v", string(conf.PrivateKey)), 0)
		return
	}

	err = privateKey.Validate()
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to validate private key: %v", privateKey), 0)
		return
	}

	return
}

func GetEfaasToken(data []byte) (googleIdToken string, err error) {
	scope := "563209362155-dmktm1rt2snprao3te1a5gf0tk9l39i8.apps.googleusercontent.com"
	claims, err := GetEfaasClaims(data, scope)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to eFaaS claims"), 0)
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := GetPrivateKeyFromJSON(data)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get private key from JSON"), 0)
		return
	}

	jwtTokenSigned, err := jwtToken.SignedString(privateKey)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get token signed with service account's key"), 0)
		return
	}

	googleIdToken, err = GetGoogleIdToken(jwtTokenSigned)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get Google Id jwtToken"), 0)
		return
	}

	return
}

func demo1(data []byte) (res []byte, err error) {
	client, err := GetEfaasClient(data)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS client"), 0)
		return
	}

	res, err = apiCallGet(client, InstancesURL)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS client"), 0)
		return
	}

	glog.Infof("AAAAA RESULT: %v", string(res))

	return
}

func EFaaS1() {
	// Credentials can be obtained from the Google Developer Console (https://console.developers.google.com).

	data := []byte(`
	{
	 "type": "service_account",
	 "project_id": "elastifile-gce-lab-c934",
	 "private_key_id": "5e0d188967e7f23ad77129ff4c9ab59889ccd25d",
	 "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCMBJyta1PEkd7q\nCLEYNdUBqk4Hlnw7mGXnByjao+4SOZi7mJ1NIAtYjptJ/rcPxjft+hxEba1a1DON\nUU7RuJ3eQk+kLVHdbD2D4noMw6VxJtuWnuyQ2V8v5ojv8kVvVSsbkDAQHVGKTe/8\nCEHxlekGoY0NC+KwWlUKmb7cv/B/2aD1eFsyV7ALE/YJmyFbbvtLrab+U5js04ER\nIWcE+gKlvAF7Xq9Iq6MucyjRvgPagz5RSP146HjbCPdJIz3ilcEL7idVGaZnnx/P\ncZAqYnYZAJTGBhi4fUEpAYR7KVUWIVfc9oXEKJDNwwBHnyyZMBPdYn9prs7xgrEL\ngA+WHPPZAgMBAAECggEACVNhUBee66+/hhzwFqm3NzYtnknCmoGK//k1GmLiv2oA\npzYB/BoPR2WwKByD+tP786i96zzW1/7cNCRfOI6wTRZjkY7HLhVAf6E8+c6qHUA2\nTfDl1rvzoBAdvMWJJGIqzdorqVcakDiirEmsgre2Xo+yAlVxUsehdGRLFw7dqNYv\nrINMqjE2W/SCd8jw2WmplmH+c0MvBKkving9CCNgFnvSMUGinv7y3Zvf2GpplvlC\nFdSFGGXxn1o6HbgrkovKn6EVZ8nP3JadG5evwjotEv1fcEu4vOKMq/jgvfxzscRf\ng9bfdhb3/oc+x43dsH3fR0axaImB7LKKgfu7w7vnJQKBgQDCmgAE7noPd0bt7Xg+\nrl44OgCHv3x0QY4lx0y07Yo1Bg1C72H8BCghr/5rxGUOSCGjoFYTVeLhCVIsYX+8\nxbtplxCJFAgN7lu48EyCgIpP7ppjf1a3Uh762O04BCMw0tXw22ich7d4KN5+r8L7\nOknRStrZYD89QjoUsSEYOK0wnwKBgQC4MePUNoBJEG+yhlMOpDz7mnf/F1U4gFQQ\nxD4stAEA1P/QuSgMb0snJJA3yT3dCL4W2DUxDCWOH/Wx3XnJy216+QR//8fHImCR\nYS4fjmaWlbMOKko1yeCtCLsNfA5uB5Yplrujn2o6v5BE52h3JCjW4qUqzZ6T9cBq\n0rQFacWwhwKBgBKLJDdUFjOFFTA08cFfUkEfXc+RsqVNXeNBs5CGFiZpVjgroXWn\nW7+iCqdwRoTu4K276JfdFkqFXdw2yjpNyUcNixjU3NOfBASCeXfyEbv+K54Rk0zS\nuXsD0s8ErenIHXTfI3/O+u+rTVBbJURVUJVuAZ63Ki+HMQupuVKai/5XAoGBALcp\nHSV8IKsHBhtfSR5JIT8MhoCKIjsyGOYnTrBDOrAqHkveor1iujetOx/OJI80T1oG\nGzavnnSqwTXiR2XrvO1IzDnADletjptiKGxGvSrGp6vRT8QXACzwfpjVIMA3GRI4\nClSVhBvxO7PY7N90fIvaCmX629LD0FgpN8weNu/nAoGAP4rXRr37757Q+c/qeKyU\nsmUCYeHj6w+GIkqJIhsDsj5tE8fLTyU87LF6hvscxYJCX9ZVycvhuzRBiFLkc9yo\nZUKC4SllFDw4Zl63RU7me3PnZHpomiNs0hk3fgqAME1Cx3Pn8NT6iptybSqk2kb7\nHOuPCeblZecVZU0UOPyQrWM=\n-----END PRIVATE KEY-----\n",
	 "client_email": "efaas-csi@elastifile-gce-lab-c934.iam.gserviceaccount.com",
	 "client_id": "102179953128561786237",
	 "auth_uri": "https://accounts.google.com/o/oauth2/auth",
	 "token_uri": "https://oauth2.googleapis.com/token",
	 "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	 "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/efaas-csi%40elastifile-gce-lab-c934.iam.gserviceaccount.com"
	}
	`)

	//data, err := ioutil.ReadFile("/path/to/key-file.json")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Initiate an http.Client, request will be authorized and authenticated on the behalf of jwt.Config.Email/Subject
	//ctx := context.Background()
	_, err := demo1(data)
	if err != nil {
		panic(err.Error())
	}
}
