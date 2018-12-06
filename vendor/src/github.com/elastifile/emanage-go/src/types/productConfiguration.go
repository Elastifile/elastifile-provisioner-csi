package types

// import (
// 	"helputils"
// 	"sysapi"
// 	"strings"
// )

// type ProdactConfigurations []ProdactParameters

// func (conf ProdactConfigurations) ToString() string {
// 	v := conf
// 	var str []string
// 	for _, params := range v {
// 		d := params
// 		str = append(str, d.ToString())
// 	}
// 	return "[" + strings.Join(str, " ") + "]"
// }

// func (fullConf ProdactConfigurations) GetConfiguration(sys *sysapi.System) (ProdactParameters, error) {
// 	isBareMetal := sys.IsBareMetal()
// 	for _, conf := range fullConf {
// 		if conf.SysType == "BM" && isBareMetal {
// 			return conf, nil
// 		}
// 		if conf.SysType == "HCI" && !isBareMetal {
// 			return conf, nil
// 		}
// 	}
// 	panic("Didn't find matching product configuration")
// }

// type ProdactParameters struct {
// 	SysType    string            `yaml:"SysType"`
// 	Boundaries ProdactBoundaries `yaml:"Boundaries"`
// 	Timeouts   ProdactTimeouts   `yaml:"Timeouts"`
// }

// func (s ProdactParameters) ToString() string {
// 	return "{" + strings.Join(helputils.StructToStrings(s), " ") + "}"
// }

// type ProdactTimeouts struct {
// 	HaOperations int `yaml:"HaOperations"` //in Minutes
// 	VheadRevive  int `yaml:"VheadRevive"`  //in Minutes
// }

// type ProdactBoundaries struct {
// 	MaxFoldersPerCluster int `yaml:"MaxFoldersPerCluster"`
// 	MaxFilesPerFolder    int `yaml:"MaxFilesPerFolder"`
// }

// // type TestsConfigurations map[string]TestParameters

// // type TestParameters struct {
// // 	SysType    string            `yaml:"SysType"`
// // 	Performance Performance `yaml:"Performance"`
// // }

// // type PerformanceTools struct {
// // 	ErunNumberOfFiles   int `yaml:"ErunNumberOfFiles"`
// // 	ErunNumberOfClients int `yaml:"ErunNumberOfClients"`
// // }
// // type Performance struct {
// // 	InWorkingSet PerformanceTools `yaml:"InWorkingSet"`
// // 	AbovePstore  PerformanceTools `yaml:"AbovePstore"`
// // 	AboveDstore  PerformanceTools `yaml:"AboveDstore"`
// // 	AboveMD      PerformanceTools `yaml:"AboveMD"`
// // }
