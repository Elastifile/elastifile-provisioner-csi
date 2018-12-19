package types

import (
	"fmt"
	"strings"
)

type TestConfigs map[string]map[string]TestParameters

func (testConfig TestConfigs) Parameters(sysType string, configName string) (TestParameters, error) {
	if len(testConfig) == 0 {
		return TestParameters{}, fmt.Errorf("Test configs not supplied cant run test")
	}
	if test, ok := testConfig[configName]; ok {
		if params, ok := test[sysType]; ok {
			return params, nil
		} else {
			if params, ok = test["Any"]; ok {
				return params, nil
			} else {
				return TestParameters{}, fmt.Errorf("system type %s is not specified for config name %s", sysType, configName)
			}
		}

	}
	return TestParameters{}, fmt.Errorf("missing parameter configuration for test %v", configName)
}

func (conf TestConfigs) ToString() string {
	v := conf
	var str []string
	for test, sysType := range v {
		str = append(str, test+":")
		for key, value := range sysType {
			str = append(str, key+":")
			str = append(str, value.ToString())
		}
	}
	return "[" + strings.Join(str, " ") + "]"
}

func (conf TestConfigs) Validate() error {
	c := conf
	for test, config := range c {
		for sysType, _ := range config {
			if sysType == "DSM" || sysType == "Any" || sysType == "HCI" {
				continue
			}
			return fmt.Errorf(fmt.Sprintf("Test configs - Test %v has an unknown system type %v", test, sysType))
		}
	}
	return nil
}
