#!/usr/bin/env bash

deployment_base="${1}"

if [[ -z ${deployment_base} ]]; then
	deployment_base="../deploy"
fi

cd "$deployment_base" || exit 1

objects=(csi-ecfsplugin-attacher csi-ecfsplugin-provisioner templates/csi-ecfsplugin csi-snapshotter-rbac csi-snapshotter snapshotclass storageclass csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac templates/configmap templates/secret)

for obj in ${objects[@]}; do
    echo "=== Deleting ${obj}"
	kubectl delete -f "./$obj.yaml"
done

pushd ${deployment_base}
./delete_crd.sh
popd

