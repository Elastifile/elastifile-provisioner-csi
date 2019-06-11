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
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/golang/glog"
	"golang.org/x/oauth2/google"

	"ecfs/log"
)

const (
	envProjectNumber = "CSI_GCP_PROJECT_NUMBER"
	envEfaasUrl      = "EFAAS_URL"
	GoogleAuthURL    = "https://www.googleapis.com/oauth2/v4/token"
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
	glog.V(log.VERBOSE_DEBUG).Infof("Using encoded values %v", encodedValues)
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

	glog.V(log.VERBOSE_DEBUG).Infof("Google body: %v", string(body))
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

		// Token expiration time
		"exp": timestamp.Add(tokenExpiration).Unix(),

		// iss (Issuer) is the service account email
		"iss": jwtConf.Email, // "service-account-name@project-name.iam.gserviceaccount.com"

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

func ProjectNumber() string {
	// TODO: Check if the value can be obtained programmatically, e.g. https://github.com/googleapis/google-cloud-ruby/issues/1416
	projectNumber := os.Getenv(envProjectNumber)
	if projectNumber == "" {
		panic(fmt.Sprintf("GCP project number not specified - expected to be present in '%v' environment variable",
			envProjectNumber))
	}
	return projectNumber
}

// efaasBaseUrl returns the base eFaaS URL, e.g. https://bronze-eagle.gcp.elastifile.com
func efaasBaseUrl() string {
	projectNumber := os.Getenv(envEfaasUrl)
	if projectNumber == "" {
		panic(fmt.Sprintf("eFaaS URL not specified - expected to be present in '%v' environment variable",
			envEfaasUrl))
	}
	return projectNumber
}

// EfaasApiUrl returns the base URL for accessing eFaaS APIs, e.g. https://bronze-eagle.gcp.elastifile.com/api/v2
func EfaasApiUrl() string {
	return efaasBaseUrl() + "/api/v2"
}
