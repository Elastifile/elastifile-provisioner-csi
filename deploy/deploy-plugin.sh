#!/usr/bin/env bash

# Template expansion variables with default values
: ${PLUGIN_TAG:="dev"} # Docker image tag
: ${MGMT_ADDR:="UNDEFINED"} # Management address
: ${MGMT_USER:="admin"} # Management user
: ${MGMT_PASS:="Y2hhbmdlbWU="} # Management user's password (base64 encoded)
: ${NFS_ADDR:="10.255.255.1"} # NFS load balancer's address
# In order to set one of the above values, run this script prefixed by the variable assignment. For example:
# PLUGIN_TAG=v0.1.0 MGMT_USER=manager ./deploy-plugin.sh

# Other variables
MYNAME=$(basename $0)
MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

DEPLOYMENT_BASE="${1}"
: ${DRY_RUN:=false}
: ${DEPLOYMENT_BASE:="${MYPATH}"}

DEFAULT_K8S_USER=${USER}
if which gcloud > /dev/null 2>&1; then
    log_info Found gcloud
    DEFAULT_K8S_USER=$(gcloud config get-value account)
fi

: ${K8S_USER:=${DEFAULT_K8S_USER}}

DRY_RUN_FLAG=""
if [[ "$DRY_RUN" = true ]]; then
    log_info "WARNING: DRY RUN"
    DRY_RUN_FLAG="--dry-run"
fi

if [[ ! -d "${DEPLOYMENT_BASE}" ]]; then
    log_error "Deployment base ${DEPLOYMENT_BASE} not found. If not running the the default location, please override \$DEPLOYMENT_BASE"
    exit 1
fi

if ! which kubectl; then
    log_error "kubectl not found"
    exit 1
fi

if ! which envsubst; then
    log_error "envsubst not found"
    exit 1
fi

if [[ -z "${K8S_USER}" ]]; then
    log_info \$K8S_USER not specified - assuming the script is running under service account with cluster-admin role
    echo "Checking permissions"
    if ! kubectl auth can-i create clusterrolebinding; then
        echo "ERROR: Current user/sa doesn't have enough permissions"
        exit 1
    fi

else
    kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user ${K8S_USER} ${DRY_RUN_FLAG}
fi

OBJECTS=(templates/configmap templates/secret csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac csi-snapshotter-rbac csi-ecfsplugin-attacher csi-ecfsplugin-provisioner csi-snapshotter storageclass snapshotclass templates/csi-ecfsplugin)

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
