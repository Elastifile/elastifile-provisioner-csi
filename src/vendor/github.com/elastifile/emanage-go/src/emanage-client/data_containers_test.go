package emanage

// We're skipping all tests as they are requiring real machine
// Emanage test suite should have two types of unitests:
// 1. unitests vs mock-able management (mock rest-api).
// 2. real tests vs real machine
// func TestDcGetAll(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := NewClient(getEmanageAddress(t))

// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	dcs, err := mgmt.DataContainers.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(dcs)
// }

// func TestDcGetAllBySearch(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	dcs, err := mgmt.DataContainers.GetAll(&DcGetAllOpts{
// 		GetAllOpts: GetAllOpts{
// 			Search: optional.NewString("default"),
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("dcs info:\n%#v", dcs)
// }

// func TestDcGetFull(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	dcs, err := mgmt.DataContainers.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dcId := dcs[0].Id
// 	dcInfo, err := mgmt.DataContainers.GetFull(dcId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(dcs)
// 	spew.Dump(dcInfo)
// 	t.Logf("dc full info:\n%#v", dcInfo)
// }

// func TestDcCreate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := NewClient(getEmanageAddress(t))
// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	policies, err := mgmt.Policies.GetAll(&GetAllOpts{PerPage: optional.NewInt(1)})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	name := fmt.Sprintf("dc-%s", uuid.New())
// 	dc, err := mgmt.DataContainers.Create(name, policies[0].Id, &DcCreateOpts{
// 		SoftQuota: 1024 << 10,
// 		HardQuota: 1024 << 11,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dc, err = mgmt.DataContainers.Update(&dc, &DcUpdateOpts{
// 		SoftQuota: dc.SoftQuota.Bytes / 2,
// 		HardQuota: dc.HardQuota.Bytes * 2,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("dc full info:\n%#v", dc)
// }

// func TestDcDeleteAllByPrefix(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	mgmt := NewClient(getEmanageAddress(t))

// 	if err := mgmt.Sessions.Login("admin", "changeme"); err != nil {
// 		t.Fatal("login failed\n", err)
// 	}
// 	defer func() {
// 		if err := mgmt.Sessions.Logout(); err != nil {
// 			t.Fatal("logout failed\n", err)
// 		}
// 	}()

// 	dcs, err := mgmt.DataContainers.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(dcs)

// 	dcPrefix := "dc-"
// 	counter := 0
// 	for i, dc := range dcs {
// 		t.Logf("%v: DC Name=%s", i, dc.Name)
// 		if strings.HasPrefix(dc.Name, dcPrefix) {
// 			t.Logf("Going to delete DC Name=%s", dc.Name)
// 			_, err := mgmt.DataContainers.Delete(&dc)
// 			if err != nil {
// 				t.Fatal(err)
// 			} else {
// 				counter++
// 			}
// 		}
// 	}

// 	t.Logf("%v were deleted", counter)
// }

// TODO: need to test create with value passed to all args
// tODO: need to test create with pointer to nill passed to all optional args
