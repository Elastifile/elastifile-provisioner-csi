#!/usr/bin/env bash

# Template expansion variables with default values
: ${PLUGIN_TAG:="dev"} # Docker image tag
: ${MGMT_ADDR:="UNDEFINED"} # Management address
: ${MGMT_USER:="admin"} # Management user
: ${MGMT_PASS:="Y2hhbmdlbWU="} # Management user's password (base64 encoded)
: ${NFS_ADDR:="10.255.255.1"} # NFS load balancer's address
: ${NAMESPACE:="default"} # K8s namespace to use for CSI plugin deployment
: ${EKFS:="false"} # Optional. If true, there's no need to specify MGMT_ADDR and NFS_ADDR
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

exec_cmd which kubectl
assert $? "kubectl not found"

exec_cmd which envsubst
assert $? "envsubst not found"

log_info "Checking permissions"
exec_cmd kubectl auth can-i create clusterrolebinding
assert $? "ERROR: Current user/sa doesn't have enough permissions to create clusterrolebinding"

if [[ -n "${K8S_USER}" ]]; then
    log_info "Assigning cluster role cluster-admin to ${K8S_USER}"
    # On repeat runs clusterrolebinding already exists and it's ok for it to fail with AlreadyExists
    exec_cmd kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user ${K8S_USER} ${DRY_RUN_FLAG} > /dev/null 2>&1
else
    log_info \$K8S_USER not specified - assuming the script is running under service account with cluster-admin role
fi

OBJECTS=(templates/configmap templates/secret templates/csi-attacher-rbac templates/csi-provisioner-rbac templates/csi-nodeplugin-rbac templates/csi-snapshotter-rbac templates/csi-snapshotter csi-ecfsplugin-attacher csi-ecfsplugin-provisioner templates/storageclass templates/csi-ecfsplugin snapshotclass)

pushd ${DEPLOYMENT_BASE}
exec_cmd ./create_crd.sh
popd

for OBJ in ${OBJECTS[@]}; do
    if [[ "${OBJ}" == *"templates"* ]]; then
        log_info "Creating ${OBJ} from template"
        PLUGIN_TAG=${PLUGIN_TAG} NAMESPACE=${NAMESPACE} MGMT_ADDR=${MGMT_ADDR} MGMT_USER=${MGMT_USER} MGMT_PASS=${MGMT_PASS} NFS_ADDR=${NFS_ADDR} EKFS=${EKFS} envsubst < "${DEPLOYMENT_BASE}/${OBJ}.yaml" | kubectl create -f - --namespace ${NAMESPACE} ${DRY_RUN_FLAG}
        assert $? "Failed to create ${OBJ} from template"
    else
        log_info "Creating ${OBJ}"
	    exec_cmd kubectl create -f "${DEPLOYMENT_BASE}/${OBJ}.yaml" --namespace ${NAMESPACE} ${DRY_RUN_FLAG}
	    EXIT_CODE=$?
        if [[ ${EXIT_CODE} != 0 && ${OBJ} == "snapshotclass" ]]; then
            # Workaround for the race between VolumeSnapshotClass CRD creation in external-snapshotter and its use in snapshotclass.yaml
            CRD="volumesnapshotclasses.snapshot.storage.k8s.io"
            MAX_RETRIES=15
            for ((attempt = 2; attempt < MAX_RETRIES+2; attempt++)); do
                echo -n .
                kubectl get crd ${CRD} > /dev/null 2>&1
                if [[ $? == 0 ]]; then
                    echo
                    log_info "Resolved the above failure - found CRD ${CRD} on attempt #${attempt}"
                    exec_cmd kubectl create -f "${DEPLOYMENT_BASE}/${OBJ}.yaml" --namespace ${NAMESPACE} ${DRY_RUN_FLAG}
                    break
                fi
                sleep 1
            done
        else
            assert ${EXIT_CODE} "Command execution failed"
        fi
    fi
done
