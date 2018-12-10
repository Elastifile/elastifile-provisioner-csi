#!/usr/bin/env bash

MYNAME=$(basename $0)
MYPATH=$(dirname $0)
DEPLOYMENT_BASE="${1}"
: ${DRY_RUN:=false}
: ${DEPLOYMENT_BASE:="../deploy"}

# Template expansion values
: ${PLUGIN_TAG:="dev"}
: ${MGMT_ADDR:="1.1.1.1"}
: ${MGMT_USER:="admin"}
: ${MGMT_PASS:="Y2hhbmdlbWU="} # base64
: ${NFS_ADDR:="10.255.255.1"}

source ${MYPATH}/functions.sh

DRY_RUN_FLAG=""
if [ "$DRY_RUN" = true ]; then
    log_info "WARNING: DRY RUN"
    DRY_RUN_FLAG="--dry-run"
fi

test -d "${DEPLOYMENT_BASE}" || exit 1

kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account) ${DRY_RUN_FLAG}

OBJECTS=(templates/configmap templates/secret csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac csi-snapshotter-rbac csi-ecfsplugin-attacher csi-ecfsplugin-provisioner csi-snapshotter snapshotclass storageclass templates/csi-ecfsplugin)

pushd ${DEPLOYMENT_BASE}
./create_crd.sh
popd

for OBJ in ${OBJECTS[@]}; do
    if [[ "${OBJ}" == *"templates"* ]]; then
        log_info "Creating ${OBJ} from template"
        PLUGIN_TAG=${PLUGIN_TAG} MGMT_ADDR=${MGMT_ADDR} MGMT_USER=${MGMT_USER} MGMT_PASS=${MGMT_PASS} NFS_ADDR=${NFS_ADDR} envsubst < "${DEPLOYMENT_BASE}/${OBJ}.yaml" | kubectl create -f - ${DRY_RUN_FLAG}
    else
        log_info "Creating ${OBJ}"
	    kubectl create -f "${DEPLOYMENT_BASE}/${OBJ}.yaml" ${DRY_RUN_FLAG}
    fi
done

