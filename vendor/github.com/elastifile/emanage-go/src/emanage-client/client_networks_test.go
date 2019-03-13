package emanage_test

// func TestClientNetworksGet(t *testing.T) {
// 	mgmt := startEManage(t)

// 	cns, err := mgmt.ClientNetworks.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	spew.Dump(cns)

// 	_, sysDetails, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	spew.Dump(sysDetails.NfsAddress)
// }

// func TestClientNetworksCRUD(t *testing.T) {
// 	mgmt := startEManage(t)

// 	cnets, err := mgmt.ClientNetworks.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	cns := cnets[0]

// 	cnName := fmt.Sprintf("CN-%x", md5.Sum([]byte(uuid.New())))
// 	ips := []string{
// 		"172.16.0.1",
// 		"172.16.0.2",
// 	}

// 	cn, err := mgmt.ClientNetworks.Create(&emanage.ClientNetwork{
// 		Name:        cnName,
// 		Vlan:        123,
// 		Subnet:      "172.16.0.0",
// 		Range:       16,
// 		IpAddresses: ips,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if cn.Name != cnName || cn.Vlan != 123 || !strSliceEqual(cn.IpAddresses, ips) {
// 		t.Fatal(fmt.Errorf("CN returned from create is not current: %+v", cn))
// 	}
// 	t.Logf("created CN: %+v", cn)

// 	cnets, err = mgmt.ClientNetworks.GetAll()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(cnets) == 0 {
// 		t.Fatal("Expected at least one network client")
// 	}

// 	cns2 := cnets[0]
// 	if len(cns2.IpAddresses) != len(cns.IpAddresses)+1 {
// 		t.Fatal(fmt.Errorf("CN list length did not increase by 1 - before: %d, after: %d", len(cns.IpAddresses), len(cns2.IpAddresses)))
// 	}
// 	t.Logf("verified CN list length increase by 1 to: %d", len(cns2.IpAddresses))

// 	cn, err = mgmt.ClientNetworks.GetById(cn.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if cn.Name != cnName || cn.Vlan != 123 || !strSliceEqual(cn.IpAddresses, ips) {
// 		t.Fatal(fmt.Errorf("CN returned from read (get) is not current: %+v", cn))
// 	}
// 	t.Logf("got CN: %+v", cn)

// 	cn.Vlan = 234
// 	ips = append(ips, "172.16.0.3")
// 	cn.IpAddresses = ips

// 	cn, err = mgmt.ClientNetworks.Update(&cn)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if cn.Name != cnName || cn.Vlan != 234 || !strSliceEqual(cn.IpAddresses, ips) {
// 		t.Fatal(fmt.Errorf("CN returned from read (get) is not current: %+v", cn))
// 	}
// 	t.Logf("updated CN: %+v", cn)

// 	cn, err = mgmt.ClientNetworks.GetById(cn.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if cn.Name != cnName || cn.Vlan != 234 || !strSliceEqual(cn.IpAddresses, ips) {
// 		t.Fatal(fmt.Errorf("CN returned from read (get) is not current: %+v", cn))
// 	}
// 	t.Logf("got updated CN: %+v", cn)

// 	cn, err = mgmt.ClientNetworks.Delete(cn.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("deleted CN: %+v", cn)

// 	_, err = mgmt.ClientNetworks.GetById(cn.Id)
// 	if err == nil {
// 		t.Fatal(fmt.Errorf("CN %d was not deleted", cn.Id))
// 	}
// 	t.Logf("verified missing CN: %+v", cn)
// }

func strSliceEqual(s1 []string, s2 []string) bool {
	for i, str := range s1 {
		if i == len(s2) {
			return len(s1) == len(s2)
		} else if str != s2[i] {
			return false
		}
	}
	return true
}
