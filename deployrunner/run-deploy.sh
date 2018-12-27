#!/usr/bin/env bash

# Usage:
# To deploy ECFS CSI plugin:
#   1. Edit deploy-runner-pod.yaml to include parameters specific to your ECFS cluster, e.g. Elastifile Management Console address
#   2. ./run-deploy.sh
# To tear down the plugin, run:
#   1. RUN_TEARDOWN=true ./run-deploy.sh

# Default values
: ${CREATE_MANIFESTS:=true}
: ${RUN_DEPLOY:=true}
: ${RUN_TEARDOWN:=false}
: ${TEARDOWN_MANIFESTS:=true}

MYNAME=$(basename $0)
MYPATH=$(dirname $0)
DEPLOYMENT_BASE=${MYPATH}

MANIFESTS=(cluster-admin-rbac.yaml)

# Handle initial permissions
if which gcloud > /dev/null 2>&1; then
    echo "Found gcloud"
    CRB_NAME=default-admin
    if ! kubectl get clusterrolebinding ${CRB_NAME} > /dev/null 2>&1; then
        echo "Creating clusterrolebinding"
        kubectl create clusterrolebinding ${CRB_NAME} --clusterrole=cluster-admin --user=$(gcloud config get-value account)
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
        kubectl create -f "${DEPLOYMENT_BASE}/${MANIFEST}"
    done
fi

# Run the deploy job
DEPLOY_POD_MANIFEST="${DEPLOYMENT_BASE}/deploy-runner-pod.yaml"
if ${RUN_DEPLOY} || ${RUN_TEARDOWN}; then
    echo "Running containerized deploy script"
    # Note: the contents of this manifest include environment variables required to properly deploy the plugin
    kubectl create -f ${DEPLOY_POD_MANIFEST}
    kubectl wait --for=condition=Ready -f ${DEPLOY_POD_MANIFEST} --timeout=1m

    kubectl cp ${TMP_KUBE_CONFIG} default/deploy-runner:/root/.kube/config
    if ${RUN_TEARDOWN}; then
        kubectl exec -it deploy-runner deploy/teardown-plugin.sh
    else
        kubectl exec -it deploy-runner deploy/deploy-plugin.sh
    fi

    kubectl delete -f ${DEPLOY_POD_MANIFEST}
fi

# Teardown
if ${TEARDOWN_MANIFESTS}; then
    for MANIFEST in ${MANIFESTS[@]}; do
        echo "Deleting ${MANIFEST}"
        kubectl delete -f "${DEPLOYMENT_BASE}/${MANIFEST}"
    done
    kubectl delete -f ${DEPLOY_POD_MANIFEST} > /dev/null 2>&1
    kubectl delete clusterrolebinding ${CRB_NAME} > /dev/null 2>&1
fi
