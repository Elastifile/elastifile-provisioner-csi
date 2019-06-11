package types_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"types"
)

//gb test -v -run TestDefultTest

// func TestDefultProduct(t *testing.T) {
// 	var productParams ProdactConfigurations
// 	filename := "/home/yonir/Documents/workspace/elfs-system/tesla/src/elastifile/tesla/tests/libs/configurations/defultProductConfiguration.yaml"
// 	err := loadProductFromYAMLFile(filename, &productParams)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("TestParams: '%s'", productParams.ToString())
// }

// func TestDefultTest(t *testing.T) {
// 	var testParams types.TestConfigs
// 	filename := "/home/yonir/Documents/workspace/elfs-system/tesla/tests/config/defultTestsConfiguration.yaml"
// 	err := loadTestsFromYAMLFile(filename, &testParams)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("TestParams: '%s'", testParams.ToString())

// 	testP, err := testParams.Parameters("BM", "TestA")
// 	if err != nil {
// 		t.Fatal("Failed getting test")
// 	} else {
// 		t.Logf("TestA/BM: '%s'", testP.ToString())
// 	}

// 	testP, err = testParams.Parameters("BM", "XXX")
// 	if err == nil {
// 		t.Fatal("Expected Fail on test name")
// 	} else {
// 		t.Log("Failed on test name - expected")
// 	}

// }

func readFromFile(filename string) (data []byte, err error) {
	if filename == "" {
		return nil, errors.New("No configuration file was specified")
	}

	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Configuration file does not exist, file='%v'", filename)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return
}

// func loadProductFromYAMLFile(filename string, tc *ProdactConfigurations) (err error) {
// 	data, err := readFromFile(filename)
// 	if err != nil {
// 		return err
// 	}
// 	err = yaml.Unmarshal(data, &tc)
// 	if err != nil {
// 		return fmt.Errorf("Unmarshal error=%v", err)
// 	}

// 	return nil
// }

func loadTestsFromYAMLFile(filename string, tc *types.TestConfigs) (err error) {
	data, err := readFromFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &tc)
	if err != nil {
		return fmt.Errorf("Unmarshal error=%v", err)
	}

	return nil
}
