package emanage_test

import (
	"encoding/json"
	"testing"

	"emanage-client"
	"optional"
)

// See comment before first test at data_containers_tests
// func TestExportsGetAll(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	exports, err := mgmt.Exports.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("exports info:\n%#v", exports)
// }

// func TestExportsGetFull(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	exports, err := mgmt.Exports.GetAll(&emanage.GetAllOpts{PerPage: optional.NewInt(1)})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	exportFull, err := mgmt.Exports.GetFull(exports[0].Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("export info:\n%#v", exportFull)
// }

// func TestExportsCreateUpdate(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	policies, err := mgmt.Policies.GetAll(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var uid int
// 	var gid int

// 	for _, access := range emanage.ExportAccessModeValues {
// 		dcName := fmt.Sprintf("dc-%x", md5.Sum([]byte(uuid.New())))
// 		if dc, err := mgmt.DataContainers.Create(dcName, policies[0].Id, nil); err != nil {
// 			t.Fatal(err)
// 		} else {
// 			exName := fmt.Sprintf("ex-%x", md5.Sum([]byte(uuid.New())))
// 			uid++
// 			gid++

// 			if export, err := mgmt.Exports.Create(
// 				exName,
// 				&emanage.ExportCreateOpts{
// 					DcId:        dc.Id,
// 					Path:        "/",
// 					Access:      access,
// 					UserMapping: emanage.UserMappingAll,
// 					Uid:         optional.NewInt(uid),
// 					Gid:         optional.NewInt(gid),
// 				},
// 			); err != nil {
// 				t.Fatal(err)
// 			} else {
// 				t.Logf("created export %s access: %s, mapping: %s, uid: %d, gid: %d",
// 					export.Name, access, emanage.UserMappingAll, uid, gid)

// 				for _, access := range emanage.ExportAccessModeValues {
// 					uid++
// 					gid++

// 					_, err := mgmt.Exports.Update(
// 						&export,
// 						&emanage.ExportUpdateOpts{
// 							Path:        "/",
// 							Access:      access,
// 							UserMapping: emanage.UserMappingRoot,
// 							Uid:         optional.NewInt(uid),
// 							Gid:         optional.NewInt(gid),
// 						},
// 					)
// 					if err != nil {
// 						t.Fatal(err)
// 					}

// 					export, err := mgmt.Exports.GetFull(export.Id)
// 					if err != nil {
// 						t.Fatal(err)
// 					}
// 					t.Logf("updated export %s access: %s, mapping: %s, uid: %d, gid: %d",
// 						export.Name, access, emanage.UserMappingRoot, uid, gid)
// 				}
// 			}
// 		}
// 	}
// }

// func TestExportsDelete(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}
// 	mgmt := startEManage(t)

// 	name := fmt.Sprintf("ex-%x", md5.Sum([]byte(uuid.New())))
// 	export, err := mgmt.Exports.Create(
// 		name,
// 		&emanage.ExportCreateOpts{
// 			DcId: 1,
// 			Path: "/",
// 		},
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf(" created export:\n%#v", export)

// 	result, err := mgmt.Exports.Delete(&export)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf(" deleted export:\n%#v", result)
// }

func TestExportsPrintNoUidGid(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	em := &emanage.ExportCreateOpts{
		DcId: 1,
		Path: "/",
	}
	t.Logf("export:\n%#v", em)

	jsonBody, err := json.Marshal(em)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("body:\n%s", jsonBody)

}

func TestExportsPrintRootUidGid(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	em := &emanage.ExportCreateOpts{
		DcId: 1,
		Path: "/",
		Uid:  optional.NewInt(0),
		Gid:  optional.NewInt(0),
	}
	t.Logf("export:\n%#v", em)

	jsonBody, err := json.Marshal(em)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("body:\n%s", jsonBody)

}
