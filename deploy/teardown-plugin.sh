#!/usr/bin/env bash

MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

deployment_base="${1}"
if [[ -z ${deployment_base} ]]; then
	deployment_base="../deploy"
fi

cd "${deployment_base}"
assert $? "Path not found: ${deployment_base}"

objects=(csi-ecfsplugin-attacher csi-ecfsplugin-provisioner templates/csi-ecfsplugin csi-snapshotter-rbac csi-snapshotter snapshotclass storageclass csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac templates/configmap templates/secret)

for obj in ${objects[@]}; do
    log_info "Deleting ${obj}"
	exec_cmd kubectl delete -f "./$obj.yaml"
done

pushd ${deployment_base}
exec_cmd ./delete_crd.sh
popd
