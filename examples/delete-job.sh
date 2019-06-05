#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

: ${JOB_MANIFEST:="${MYPATH}/job-io.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc.yaml"}
: ${JOB_CLEANUP_MANIFEST:="${MYPATH}/job-cleanup-data.yaml"}
: ${NAMESPACE:="default"}

echo "Deleting the main job"
assert_cmd kubectl delete --wait -f ${JOB_MANIFEST} --namespace ${NAMESPACE}

echo "Creating the cleanup job"
assert_cmd kubectl create -f ${JOB_CLEANUP_MANIFEST} --namespace ${NAMESPACE}
kubectl wait --for=condition=complete -f ${JOB_CLEANUP_MANIFEST} --timeout=2m  --namespace ${NAMESPACE}

echo "Deleting the cleanup job"
assert_cmd kubectl delete --wait -f ${JOB_CLEANUP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting pvc"
assert_cmd kubectl delete --wait -f ${PVC_MANIFEST} --namespace ${NAMESPACE}

echo "Delete completed"
