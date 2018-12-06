package emanage

// func TestStatisticsGetAllNil(t *testing.T) {
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

// 	stats, err := mgmt.Statistics.GetAll(nil)
// 	if err != nil {
// 		t.Fatal("get statistics failed\n", err)
// 	}
// 	t.Logf("last stats entry is:\n%#v", stats[len(stats)-1])
// }

// func TestStatisticsGetAllAfterT(t *testing.T) {
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

// 	stats, err := mgmt.Statistics.GetAll(nil)
// 	if err != nil {
// 		t.Fatal("get statistics failed\n", err)
// 	}

// 	startTime := stats[1].TimeStamp
// 	stats, err = mgmt.Statistics.GetAll(etime.NewNilableTime(startTime))
// 	if err != nil {
// 		t.Fatal("get statistics after time failed\n", err)
// 	}
// 	t.Logf("last stats entry is:\n%#v", stats[len(stats)-1])
// }

// func TestStatisticsGetFull(t *testing.T) {
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

// 	startTime := etime.NewNilableTime(time.Now().Add(-5 * time.Hour))
// 	stats, err := mgmt.Statistics.GetAll(startTime)
// 	if err != nil {
// 		t.Fatal("get statistics failed\n", err)
// 	}
// 	statId := stats[0].Id
// 	t.Log("statId is: ", statId)

// 	fullStat, err := mgmt.Statistics.GetFull(statId)
// 	if err != nil {
// 		t.Fatal("mgmt.Statistics.GetFull failed\n", err)
// 	}
// 	t.Logf("full stat is:\n%#v", fullStat)
// }
