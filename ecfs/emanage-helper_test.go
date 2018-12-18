package main

import (
	"fmt"
	"testing"

	"github.com/go-errors/errors"

	"csi-provisioner-elastifile/ecfs/co"
)

func fakeEmsConfig() (err error) {
	var dataConfigMap = map[string]string{ // Should match configmap manifest keys
		//"managementAddress":  "https://10.11.209.228",
		"managementAddress":  "https://35.205.94.59",
		"managementUserName": "admin",
		"nfsAddress":         "172.16.0.1", // On-prem default
		//"nfsAddress": "10.255.255.1", // GCP default
	}

	// TODO: If Update fails - create the config map
	// Current workaround - kubectl create -f deploy/configmap.yaml
	err = co.CreateConfigMap("default", configMapName, dataConfigMap)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed UpdateConfigMap with %+v", dataConfigMap), 0)
	}

	var dataSecrets = map[string][]byte{ // Should match secrets manifest keys
		"password": []byte("Y2hhbmdlbWU="),
	}

	// TODO: If Update fails - create the secrets
	// Current workaround - kubectl create -f deploy/secret.yaml
	err = co.CreateSecrets("default", secretsName, dataSecrets)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed UpdateSecrets with %+v", dataSecrets), 0)
	}

	return
}

func TestGetSnapshotByName(t *testing.T) {
	var snapshotName = "vs-111-222"
	var ems emanageClient

	err := fakeEmsConfig()
	if err != nil {
		t.Fatal("GetSnapshotByName failed: ", err)
	}

	snapshot, err := ems.GetSnapshotByName(snapshotName)
	if err != nil {
		t.Fatal("GetSnapshotByName failed: ", err)
	}

	t.Log("TestGetSnapshotByName", "snapshot.Name", snapshot.Name, "snapshot.ID", snapshot.ID, "snapshot", *snapshot)
}
