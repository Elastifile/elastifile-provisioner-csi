package rest

import (
	"net/url"
	"os"
	"testing"
)

func getEmanageAddress(t *testing.T) *url.URL {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const envName = "TESLA_EMANAGE_SERVER"
	host := os.Getenv(envName)
	t.Logf("Env settings %v", host)
	//if host == "" {
	//	t.Fatalf("Environment variable %v not set", envName)
	//}

	baseURL := &url.URL{
		Scheme: "http",
		Host:   "10.11.209.226",
	}

	return baseURL
}

func TestLogin(t *testing.T) {
	// logging.Setup(logging_config.ConfigForUnitTest())

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	mgmt := NewDefaultSession(getEmanageAddress(t))

	t.Log("mgmt.Login")
	t.Logf("Session:, verify: %v", mgmt.Url())
	err := mgmt.Login("admin", "changeme")
	if err != nil {
		t.Fatal("login:", err)
	}

	t.Log("mgmt.Logout")
	err = mgmt.Logout()
	if err != nil {
		t.Fatal("logout:", err)
	}
}
