package emanage

// func TestTennantsGetAll(t *testing.T) {
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

// 	tenants, err := mgmt.Tenants.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("tenants info:\n%#v", tenants)
// }
