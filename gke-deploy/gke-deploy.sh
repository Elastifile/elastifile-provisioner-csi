#!/usr/bin/env bash

MYPATH=$(dirname $0)

: ${NAMESPACE:="ekfs-ns"}

source ${MYPATH}/functions.sh

function safe_create () {
    manifest="$1"
    exec_cmd kubectl delete -f "${manifest}" --namespace ${NAMESPACE}
    assert_cmd kubectl create -f "${manifest}" --namespace ${NAMESPACE}

}

function safe_create_template () {
    manifest="$1"
    NAMESPACE=${NAMESPACE} envsubst < "${manifest}" | kubectl delete -f -
    NAMESPACE=${NAMESPACE} envsubst < "${manifest}" | kubectl create -f -
    if [[ $? != 0 ]]; then
        echo "Failed to create StorageClass"
        exit 1
    fi
}

# Cluster-wide resources are not supported by GKE schema.yaml method - create them here

safe_create_template "${MYPATH}/templates/storageclass.yaml"

safe_create "${MYPATH}/snapshotclass.yaml"
safe_create "${MYPATH}/csidriver.yaml"
safe_create "${MYPATH}/csinodeinfo.yaml"

#NAMESPACE=${NAMESPACE} envsubst < "${MYPATH}/templates/storageclass.yaml" | kubectl delete -f -
#NAMESPACE=${NAMESPACE} envsubst < "${MYPATH}/templates/storageclass.yaml" | kubectl create -f -
#if [[ $? != 0 ]]; then
#    echo "Failed to create StorageClass"
#    exit 1
#fi

#assert_cmd kubectl create -f "${MYPATH}/snapshotclass.yaml" --namespace ${NAMESPACE}
#assert_cmd kubectl create -f "${MYPATH}/csidriver.yaml"
#assert_cmd kubectl create -f "${MYPATH}/csinodeinfo.yaml"

echo "Manifests created"
