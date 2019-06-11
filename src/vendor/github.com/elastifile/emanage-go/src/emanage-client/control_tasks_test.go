package emanage

// We're skipping all tests as they are requiring real machine
// Emanage test suite should have two types of unitests:
// 1. unitests vs mock-able management (mock rest-api).
// 2. real tests vs real machine
// func TestControlTasksGetAll(t *testing.T) {
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

// 	tasks, err := mgmt.ControlTasks.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(tasks)
// }

// func TestControlTasksMonitor(t *testing.T) {
// 	// To generate control tasks for monitor test we'll force reset the system
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

// 	done, _ := mgmt.ControlTasks.Monitor()
// 	defer close(done)

// 	system, _, err := mgmt.Systems.GetById(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = system.ForceReset()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestControlTasksGetRecent(t *testing.T) {
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

// 	task, err := mgmt.ControlTasks.GetRecent()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Got recent task, ID: %v, Description: %v", task.ID, task.Name)
// }
