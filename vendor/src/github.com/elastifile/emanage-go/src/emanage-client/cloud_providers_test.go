package emanage_test

import (
//	emanage "emanage-client"
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"
)

const (
	envName   string = "TESLA_SYSTEM_ELAB_DATA_EMANAGE_VIP"
	emanageIp string = "18.196.196.178"
)

func TestCloudProvidersGetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName, emanageIp)
	mgmt := startEManage(t)

	providers, err := mgmt.CloudProviders.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(providers)
}

func TestCloudProvidersGetById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName, emanageIp)
	mgmt := startEManage(t)

	providers, err := mgmt.CloudProviders.GetById(1)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(providers)
}

//func TestCloudProvidersUpdate(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//	}
//	os.Setenv(envName, emanageIp)
//	mgmt := startEManage(t)
//
//	clp, err := mgmt.CloudProviders.GetById(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	size := int(111)
//	opt := emanage.CloudProvidersUpdOpts{
//		CloudProviderCommonOpts: emanage.CloudProviderCommonOpts{
//			Project:         clp.Project,
//			LocalNumOfDisks: 2,
//			StorageType:     "local",
//		},
//		LocalDiskSize: size,
//	}
//
//	providers, err := mgmt.CloudProviders.Update(1, &opt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	spew.Dump(providers)
//}
