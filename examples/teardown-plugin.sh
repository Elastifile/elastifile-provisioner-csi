#!/bin/bash

deployment_base="${1}"

if [[ -z $deployment_base ]]; then
	deployment_base="../deploy"
fi

cd "$deployment_base" || exit 1

objects=(csi-ecfsplugin-attacher csi-ecfsplugin-provisioner csi-ecfsplugin csi-snapshotter-rbac csi-snapshotter snapshotclass storageclass csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac configmap secret)

for obj in ${objects[@]}; do
    echo "=== Deleting ${obj}"
	kubectl delete -f "./$obj.yaml"
done
