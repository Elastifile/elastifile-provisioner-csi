#!/usr/bin/env bash

# Default values
: ${DRY_RUN:=false}
: ${CREATE_MANIFESTS:=true}
: ${RUN_DEPLOY:=true}
: ${TEARDOWN_MANIFESTS:=true}


MYNAME=$(basename $0)
MYPATH=$(dirname $0)
DEPLOYMENT_BASE=${MYPATH}

DRY_RUN_FLAG=""
if ${DRY_RUN}; then
    echo "WARNING: DRY RUN"
    DRY_RUN_FLAG="--dry-run"
fi

MANIFESTS=(cluster-admin-rbac.yaml)

# Handle initial permissions
if which gcloud > /dev/null 2>&1; then
    echo "Found gcloud"
    CRB_NAME=default-admin
    if ! kubectl get clusterrolebinding ${CRB_NAME} > /dev/null 2>&1; then
        echo "Creating clusterrolebinding"
        kubectl create clusterrolebinding ${CRB_NAME} --clusterrole=cluster-admin --user=$(gcloud config get-value account) ${DRY_RUN_FLAG}
    fi
fi

echo "Checking permissions"
if ! kubectl auth can-i create clusterrolebinding; then
    echo "ERROR: This user doesn't have enough permissions"
    exit 1
fi

# Prepare kubeconfig
TMP_KUBE_CONFIG=/tmp/config
echo "Dumping current kube config into ${TMP_KUBE_CONFIG}"
kubectl config view --minify --flatten > ${TMP_KUBE_CONFIG}

# Setup
if ${CREATE_MANIFESTS}; then
    for MANIFEST in ${MANIFESTS[@]}; do
        echo "Creating ${MANIFEST}"
        kubectl create -f "${DEPLOYMENT_BASE}/${MANIFEST}" ${DRY_RUN_FLAG}
    done
fi

# Run the deploy job
DEPLOY_MANIFEST="${DEPLOYMENT_BASE}/deploy-runner-pod.yaml"
if ${RUN_DEPLOY}; then
    echo "Running containerized deploy script"
    # Note: the contents of this manifest include environment variables required to properly deploy the plugin
    kubectl create -f ${DEPLOY_MANIFEST} ${DRY_RUN_FLAG}
    kubectl wait --for=condition=Ready -f ${DEPLOY_MANIFEST} ${DRY_RUN_FLAG} --timeout=1m

    kubectl cp ${TMP_KUBE_CONFIG} default/deploy-runner:/root/.kube/config
    kubectl exec -it deploy-runner deploy/deploy-plugin.sh

    kubectl delete -f ${DEPLOY_MANIFEST} ${DRY_RUN_FLAG}
fi

# Teardown
if ${TEARDOWN_MANIFESTS}; then
    for MANIFEST in ${MANIFESTS[@]}; do
        echo "Deleting ${MANIFEST}"
        kubectl delete -f "${DEPLOYMENT_BASE}/${MANIFEST}" ${DRY_RUN_FLAG}
    done
    kubectl delete -f ${DEPLOY_MANIFEST} ${DRY_RUN_FLAG} > /dev/null 2>&1
    kubectl delete clusterrolebinding ${CRB_NAME} > /dev/null 2>&1
fi
