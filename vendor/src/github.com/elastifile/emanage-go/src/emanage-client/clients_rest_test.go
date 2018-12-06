package emanage_test

// func TestClientsGetAll(t *testing.T) {
// 	mgmt := startEManage(t)

// 	clients, err := mgmt.Clients.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, cl := range clients {
// 		t.Logf("Name: %s, Id: %d, Status: %s", cl.Name, cl.Id, cl.Status)
// 	}
// }

// func TestClientsGetById(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	clients, err := mgmt.Clients.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	any := rand.Intn(len(clients))
// 	client, err := mgmt.Clients.GetById(clients[any].Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("client info:\n%#v", client)
// }

// func TestClientsGetStats(t *testing.T) {
// 	mgmt := startEManage(t)

// 	stats, err := mgmt.Clients.GetStats()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Got %d stats", len(stats))

// 	for i, st := range stats {
// 		t.Logf("%d. client ip: %s, stats: %s", i+1, st.ClientIp, st.Stats)
// 	}
// }
