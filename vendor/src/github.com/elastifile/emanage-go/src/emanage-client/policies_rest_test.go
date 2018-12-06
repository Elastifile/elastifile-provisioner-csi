package emanage

// func TestPoliciesGetAll(t *testing.T) {
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

// 	policies, err := mgmt.Policies.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("policies info:\n%#v", policies)
// }

// func TestPoliciesGetFull(t *testing.T) {
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

// 	policyFull, err := mgmt.Policies.GetFull(policies[0].Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("policy info:\n%#v", policyFull)
// }

// func TestPoliciesCreate(t *testing.T) {
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

// 	name := fmt.Sprintf("policy-%s", uuid.New())
// 	policy, err := mgmt.Policies.Create(name, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf(" created policy:\n%#v", policy)
// }
