#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

POD_MANIFEST=$1
: ${POD_MANIFEST:="${MYPATH}/pod-io.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc.yaml"}
: ${NAMESPACE:="default"}

assert_cmd kubectl create -f ${PVC_MANIFEST} --namespace ${NAMESPACE}
assert_cmd kubectl create -f ${POD_MANIFEST} --namespace ${NAMESPACE}

echo "Waiting for the pod to become Ready"
assert_cmd kubectl wait --for=condition=Ready -f ${POD_MANIFEST} --timeout=2m --namespace ${NAMESPACE}
