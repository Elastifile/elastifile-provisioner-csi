package emanage_test

// func TestSnapshotsGetList(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	snapshotsList, err := mgmt.Snapshots.Get()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapshotsList)
// }

// func TestSnapshotsGetById(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	snapshot, err := mgmt.Snapshots.GetById(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapshot)
// }

// func TestSnapshotsUpdate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	snapshotData := emanage.Snapshot{
// 		ID:   1,
// 		Name: "UpdateName",
// 	}
// 	snapshot, err := mgmt.Snapshots.Update(&snapshotData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapshot)
// }

// func TestSnapshotsCreate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	snapshotData := emanage.Snapshot{
// 		Name:            "TestName",
// 		DataContainerID: 1,
// 	}
// 	snapshot, err := mgmt.Snapshots.Create(&snapshotData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapshot)
// }
// func TestSnapshotsConsistencyGroup(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	snapshotData := emanage.Snapshot{
// 		Name:                 "TestNameCG",
// 		ConsistencyGroupDcID: []int{1, 2, 3},
// 	}
// 	snapshot, err := mgmt.Snapshots.CreateConsistencyGroup(&snapshotData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapshot)
// }

// func TestSnapshotsDelete(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	snapshot, err := mgmt.Snapshots.GetById(11)

// 	tasks := snapshot.Delete()
// 	if tasks.Error() != nil {
// 		t.Fatal(tasks.Error())
// 	}
// 	err = tasks.Wait()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestSnapshotsStatistics(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	snapshot, err := mgmt.Snapshots.GetById(6)
// 	snapStat, err := snapshot.Statistics()
// 	// snapshot, err := mgmt.Snapshots.Delete(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapStat)
// }

// func TestSnapshotsCreateLock(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	snapshot, err := mgmt.Snapshots.GetById(6)
// 	snapshotLockData := emanage.SnapshotLock{
// 		Name: "snap_lock2",
// 	}
// 	snapLock, err := snapshot.CreateLock(&snapshotLockData)
// 	// snapshot, err := mgmt.Snapshots.Delete(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapLock)
// }

// func TestSnapshotsListLocks(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	snapshot, err := mgmt.Snapshots.GetById(6)
// 	snapLockList, err := snapshot.ListLocks()
// 	// snapshot, err := mgmt.Snapshots.Delete(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapLockList)
// }

// func TestSnapshotsDeleteLock(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)
// 	snapshot, err := mgmt.Snapshots.GetById(6)
// 	snapLockList, err := snapshot.DeleteLock(3)
// 	// snapshot, err := mgmt.Snapshots.Delete(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(snapLockList)
// }
