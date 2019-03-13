package emanage_test

// func TestEMSGetList(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	emsList, err := mgmt.EmanageVMs.Get()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(emsList)
// }

// func TestEMSGetVIPFromList(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	emsList, err := mgmt.EmanageVMs.Get()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(emsList)

// 	if emsList != nil && len(emsList) > 1 {
// 		vip, err := emsList.GetVIP()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		fmt.Printf("Found VIP:'%v'\n", vip)
// 	}
// }

// func TestEMSGetPassive(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	emsList, err := mgmt.EmanageVMs.Get()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(emsList)

// 	if emsList != nil {
// 		passiveEMS, err := emsList.GetPassiveEmanageIP()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		fmt.Printf("Found passive emanage:'%v'\n", passiveEMS)
// 	}
// }

// func TestEMSGetActive(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	emsList, err := mgmt.EmanageVMs.Get()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(emsList)

// 	if emsList != nil {
// 		emanageIPs := getEmanageIps(t)
// 		fmt.Printf("emanage IP addresses list:%v\n", emanageIPs)

// 		activeEMS, err := emsList.ActiveEmanageIP()

// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		fmt.Printf("Found active emanage:'%v'\n", activeEMS)
// 	}
// }

// func TestEMSMonitor(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	monitor, err := mgmt.GetEMSMonitor()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(monitor)
// }
