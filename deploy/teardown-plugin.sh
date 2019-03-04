#!/usr/bin/env bash

MYPATH=$(dirname $0)

: ${NAMESPACE:="default"} # K8s namespace to use for CSI plugin deployment

source ${MYPATH}/functions.sh

deployment_base="${1}"
if [[ -z ${deployment_base} ]]; then
	deployment_base="../deploy"
fi

cd "${deployment_base}"
assert $? "Path not found: ${deployment_base}"

objects=(csi-ecfsplugin-attacher csi-ecfsplugin-provisioner templates/csi-ecfsplugin templates/csi-snapshotter-rbac templates/csi-snapshotter snapshotclass templates/storageclass templates/csi-attacher-rbac templates/csi-provisioner-rbac templates/csi-nodeplugin-rbac templates/configmap templates/secret)

for obj in ${objects[@]}; do
    log_info "Deleting ${obj}"
	exec_cmd kubectl delete -f "./$obj.yaml --namespace ${NAMESPACE}"
done

pushd ${deployment_base}
exec_cmd ./delete_crd.sh
popd
