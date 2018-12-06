package emanage_test

import (
	"testing"
)

// func TestLicenseGet(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	license, err := system.GetLicense()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(license)
// }

// func TestLicenseUpload1(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	licenseOpts := emanage.LicenseOpts{
// 		License: "order_number=8567951 " +
// 			"expiration_date=04.04.2018 " +
// 			"raw_capacity=50TB " +
// 			"hosts=3 " +
// 			"signature=f1zN2FF1VEZpDOTtcJGMD8tvq7OmGGdCCSIV4nww/ttWYfVQ0d/x76lol6bfZn0pwY4cgXd9SGaznOBNw4+xiatMYu+SzHXTQwXXzgXnUGDQbJwbkeBfWxZfNK6fxcTx3ui8JhPYdmRrjme1mOymI5ZI/BWlsfE1RXeAEDqaBEw=",
// 	}

// 	license, err := system.UploadLicense(&licenseOpts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(license)
// }

// func TestLicenseUpload2(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	licenseOpts := emanage.LicenseOpts{
// 		License: "order_number=8567952 " +
// 			"expiration_date=01.01.2020 " +
// 			"raw_capacity=100TB " +
// 			"hosts=4 " +
// 			"signature=eG0II1HuNmUF48lqQZ+hx2W5HCeQ5Rulhlb1IlZJB/JfaMIpEuY/qeAf2GxF 3OH8RbcVfOPnjNo0DvT1rjDFq0QAxI18jBvg85REk56jVfqwrzcDQmh16PMC PALjZD1WMOPmERAKwQFhAoBn8JADnFouf3cWwo2vHwVc4bmyLVc=",
// 	}

// 	license, err := system.UploadLicense(&licenseOpts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(license)
// }

// func TestLicenseLoadFromFile(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	licenseFileContent := "order_number=8567952\n" +
// 		"expiration_date=01.01.2020\n" +
// 		"raw_capacity=100TB\n" +
// 		"hosts=4\n" +
// 		"signature=eG0II1HuNmUF48lqQZ+hx2W5HCeQ5Rulhlb1IlZJB/JfaMIpEuY/qeAf2GxF 3OH8RbcVfOPnjNo0DvT1rjDFq0QAxI18jBvg85REk56jVfqwrzcDQmh16PMC PALjZD1WMOPmERAKwQFhAoBn8JADnFouf3cWwo2vHwVc4bmyLVc=\n"

// 	filename := "TestLicenseLoadFromFile"
// 	f, err := ioutil.TempFile("", filename)
// 	check(t, err)
// 	defer func() { _ = os.Remove(f.Name()) }()

// 	err = ioutil.WriteFile(f.Name(), []byte(licenseFileContent), os.ModePerm)
// 	check(t, err)

// 	var licenseOpts emanage.LicenseOpts
// 	err = licenseOpts.LoadFromFile(f.Name())
// 	check(t, err)

// 	spew.Dump(licenseOpts)

// 	mgmt := startEManage(t)

// 	system, _, err := mgmt.Systems.GetById(sysId)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	license, err := system.UploadLicense(&licenseOpts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	spew.Dump(license)
// }

///////////////////////////////////////////////////////////////////////////////
func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
