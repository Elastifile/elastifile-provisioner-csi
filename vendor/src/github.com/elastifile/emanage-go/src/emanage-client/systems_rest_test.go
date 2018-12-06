package emanage_test

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"config"
	"emanage-client"
	"types"
)

const sysId = 1

func startEManage(t *testing.T) (mgmt *emanage.Client) {
	conf, err := config.FromTagsAndEnvironment()
	if err != nil {
		t.Fatal("Loading configuration failed\n", err)
	}

	baseURL := conf.EmanageURL()

	mgmt = emanage.NewClient(baseURL)
	t.Log("logging in -", "url", baseURL, "conf", conf.Tesla.Emanage)
	if err = mgmt.Sessions.Login(conf.Tesla.Emanage.Username, conf.Tesla.Emanage.Password); err != nil {
		t.Fatal("login failed\n", err)
	}
	t.Log("started emanage client at", baseURL)

	return mgmt
}

func getEmanageIps(t *testing.T) []string {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	conf := &types.Config{}
	err := config.FromAllSources(conf, nil, nil)
	if err != nil {
		t.Fatal("Loading configuration failed\n", err)
	}

	return conf.System.Elab.EmanageIps()
}

func TestSystemUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	os.Setenv(envName, emanageIp)
	nameServer := "stam2018.elastifile.com"
	mgmt := startEManage(t)
	_, details, err := mgmt.Systems.GetById(sysId)
	if err != nil {
		t.Fatal(err)
	}
	opts := emanage.SystemDetails{
		Name:            details.Name,
		NameServer:      nameServer,
		DeploymentModel: "cloud_amazon",
	}
	details, err = mgmt.Systems.Update(sysId, &opts)

	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(details)

}

func TestSystem_Setup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	os.Setenv(envName, emanageIp)
	mgmt := startEManage(t)
	system, _, err := mgmt.Systems.GetById(sysId)
	if err != nil {
		t.Fatal(err)
	}

	system.Setup(nil, true)

}
func TestSystem_AcceptEULA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	os.Setenv(envName, emanageIp)
	mgmt := startEManage(t)
	system, _, err := mgmt.Systems.GetById(sysId)
	if err != nil {
		t.Fatal(err)
	}
	system.AcceptEULA()
}

// func TestSystemGetAll(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	systems, err := mgmt.Systems.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump("systems info:\n%#v", systems)
// }

// func TestSystemGetById(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	_, details, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(details)
// }

// func TestSystemStatus(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	details, err := system.GetDetails()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("*** Status: %v\n", details.Status)
// }

// func TestSystemListReports(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	listReports, err := system.ListReports()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmt.Printf("%+v", listReports)
// }

// func TestSystemCreateReport(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	uuid, createReports, err := system.CreateReportForAllNodes(emanage.ReportTypeMinimal)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmt.Printf("%s\n", uuid[0])
// 	fmt.Printf("%+v\n", createReports)

// 	downDetails, e := system.PrepareReportFromAllNodes(uuid[0])
// 	if e != nil {
// 		t.Fatal(e)
// 	}

// 	fmt.Printf("%+v\n", downDetails)

// 	sysDetails, err := system.DeleteReportOnAllNodes(uuid[0])

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmt.Printf("%+v\n", sysDetails)

// }

// func TestSystemForceReset(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = system.ForceReset()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestSystemStart(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	opts := emanage.SystemStartOpts{
// 		SkipTests: true,
// 	}
// 	details, err := system.Start(opts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(details)
// }

// func TestSystemShutdown(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err := system.Shutdown()
// Expect(err).ToNot(HaveOccurred())

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(details)
// }

// func TestSystemAnswerFile(t *testing.T) {
// 	// integration test
// 	if true {
// 		t.Skip("skipping integration test")
// 	}

// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	answer, err := system.AnswerFile()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// answerFile.ByNode("esx5a").Services.DataStore.Devices
// 	dstores, err := answer.ByHost("esx5d.lab.il.elastifile.com")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, r := range dstores.DStoreServices() {
// 		fmt.Printf("service id: %v\n", r.ID.Type)
// 		fmt.Printf("service devices: %v\n", r.DeviceID)
// 	}
// }

// func TestSystemCapacity(t *testing.T) {
// 	// integration test
// 	if true {
// 		t.Skip("skipping integration test")
// 	}

// 	mgmt := startEManage(t)
// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	cap, err := system.Capacity()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(cap)
// }

// func TestUpgradeStart(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := startEManage(t)
// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	body := emanage.UpgradeOpts{Type: "rolling",
// 		DegradedReplication: true,
// 		AdminPasswd:         "changeme",
// 		SkipTest:            true,
// 	}
// 	upgr, err := system.UpgradeStart(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(upgr)
// }

// func TestUpgradePause(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = system.UpgradePause()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
// func TestUpgradeResume(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := startEManage(t)
// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = system.UpgradeResume()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
