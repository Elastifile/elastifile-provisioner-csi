#!/usr/bin/env bash

MYPATH=$(dirname $0)

: ${NAMESPACE:="ekfs-ns"}

source ${MYPATH}/functions.sh

# Create cluster-wide resources, which are not supported by GKE schema.yaml method

NAMESPACE=${NAMESPACE} envsubst < "${MYPATH}/templates/storageclass.yaml" | kubectl create -f -
if [[ $? != 0 ]]; then
    echo "Failed to create StorageClass"
    exit 1
fi

assert_cmd kubectl create -f "${MYPATH}/snapshotclass.yaml" --namespace ${NAMESPACE}
assert_cmd kubectl create -f "${MYPATH}/csidriver.yaml"
assert_cmd kubectl create -f "${MYPATH}/csinodeinfo.yaml"

echo "Manifests created"
