#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

: ${SNAP_MANIFEST:="${MYPATH}/snapshot.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc-restore-snap.yaml"}
: ${POD_MANIFEST:="${MYPATH}/pod-restore-snap.yaml"}
: ${NAMESPACE:="default"}

assert_cmd kubectl create -f ${SNAP_MANIFEST} --namespace ${NAMESPACE}
assert_cmd kubectl create -f ${PVC_MANIFEST} --namespace ${NAMESPACE}
assert_cmd kubectl create -f ${POD_MANIFEST} --namespace ${NAMESPACE}

TIMEOUT=5m
echo "Waiting for the pod to become Ready for up to ${TIMEOUT}"
assert_cmd kubectl wait --for=condition=Ready -f ${POD_MANIFEST} --timeout=${TIMEOUT} --namespace ${NAMESPACE}
