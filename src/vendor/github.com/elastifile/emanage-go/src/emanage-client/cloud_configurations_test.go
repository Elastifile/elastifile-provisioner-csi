package emanage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"types"
)

const (
	envName1   string = "TESLA_SYSTEM_ELAB_DATA_EMANAGE_VIP"
	emanageIp1 string = "35.158.56.46"
)

func TestPrint(t *testing.T) {
	fmt.Println(t)
}

func TestCloudConfigurationsGetAll(t *testing.T) {
	fmt.Println("test")
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName1, emanageIp1)
	mgmt := startEManage(t)

	providers, err := mgmt.CloudConfigurations.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range *providers {
		fmt.Println(p)
	}
}

func TestCloudConfigurationsGetById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName1, emanageIp1)
	mgmt := startEManage(t)

	providers, err := mgmt.CloudConfigurations.GetById(6)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(providers)
}

//func TestCloudConfigurationsCreate(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//	}
//	os.Setenv(envName1, emanageIp1)
//	mgmt := startEManage(t)
//
//	clc, err := mgmt.CloudConfigurations.GetById(6)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(clc)
//	opt := types.CloudConfigurationCreateOpts{
//		Name:              "TestUT123312",
//		NumOfDisks:        1,
//		StorageType:       string(types.LocalStorageType),
//		DiskSize:          10,
//		MinNumOfInstances: 3,
//		InstanceType:      string(types.LocalInstanceType),
//	}
//
//	providers, err := mgmt.CloudConfigurations.Create(&opt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	spew.Dump(providers)
//	fmt.Println(providers.ID)
//}

func TestCloudConfigurationsUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName1, emanageIp1)
	mgmt := startEManage(t)

	clc, err := mgmt.CloudConfigurations.GetById(6)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(clc)
	test := types.CloudConfigurationUpdateOpts{}
	test.ID = 6
	test.Name = "LiranTheKing"
	providers, err := mgmt.CloudConfigurations.Update(&test)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(providers)
	fmt.Println(providers.ID)
}
