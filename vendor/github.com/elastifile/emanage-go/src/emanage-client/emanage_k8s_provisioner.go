package emanage

import (
	"net/url"
	"fmt"
)

const DC_NAME = "DC_CREATED_BY_PROVISIONER"


func CreateDC(IP string)  {

	fmt.Printf("Create DC on %v", IP)
	EMSClient := getLoggedinClient(IP)

	// create DC
	dcCreateOpts := DcCreateOpts{Name: DC_NAME, Dedup: 0, Compression: 0}
	dc, err := EMSClient.DataContainers.Create(DC_NAME, 1, &dcCreateOpts)
	if err != nil {
		fmt.Printf("Error creating DC: %v", err)
	}
	fmt.Printf("Created dc: %v ID: %v", dc.Name, dc.Id)
}

func getLoggedinClient(IP string) *Client {

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

