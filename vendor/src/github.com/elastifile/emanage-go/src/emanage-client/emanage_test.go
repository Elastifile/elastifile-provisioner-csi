package emanage

import (
	"fmt"
	"net/url"
	"testing"

	"optional"
)

const DC_NAME_TEST = "AAA-emanage-client-DC"

func TestCreateDC(t *testing.T) {
	CreateDCTest(t)
}

func CreateDCTest(t *testing.T) DataContainer {

	fmt.Println("Starting Create DC test")
	EMSClient := getLoggedinClientTest("10.11.209.226")

	// create DC
	dcCreateOpts := DcCreateOpts{Name: DC_NAME_TEST, Dedup: 0, Compression: 0}
	dc, err := EMSClient.DataContainers.Create(DC_NAME_TEST, 1, &dcCreateOpts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created dc: %v ID: %v", dc.Name, dc.Id)
	return dc
}

func TestDeleteDC(t *testing.T) {
	DeleteDC(t)
}

func DeleteDC(t *testing.T) bool {

	fmt.Println("Starting Create DC test")
	EMSClient := getLoggedinClient("10.11.209.226")

	// get DC
	opts := DcGetAllOpts{
		GetAllOpts: GetAllOpts{
			Search: optional.NewString(DC_NAME_TEST),
		},
	}
	dcs, err := EMSClient.DataContainers.GetAll(&opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(dcs) == 1 {
		//delete DC
		EMSClient.DataContainers.Delete(&dcs[0])
		return true
	} else {
		fmt.Println("More than 1 DC found")
	}
	return false
}

func TestEMSClient(t *testing.T) {

	fmt.Println("Starting EMS EMSClient test")
	EMSClient := getLoggedinClient("35.241.215.32")

	// get all hosts
	opts := &GetAllOpts{}
	hostArr, err := EMSClient.Hosts.GetAll(opts)
	if err != nil {
		fmt.Printf("Error getting all hosts: %s", err)
		return
	}

	// print
	for i, host := range hostArr {
		fmt.Printf("Host# %v: %v, Role: %v\n", i, host.Name, host.Role)
	}
}

func getLoggedinClientTest(IP string) *Client {

	// make EMSClient
	eurl := &url.URL{
		Scheme: "http",
		Host:   IP,
	}
	EMSClient := NewClient(eurl)

	// login
	err := EMSClient.Sessions.Login("admin", "changeme")
	if err != nil {
		fmt.Printf("Error logging in: %s", err)
		return nil
	} else {
		fmt.Println("Logged in")
	}
	return EMSClient
}
