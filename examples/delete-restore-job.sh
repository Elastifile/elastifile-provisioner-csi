#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

: ${JOB_CLEANUP_MANIFEST:="${MYPATH}/job-cleanup-restore-snap-data.yaml"}
: ${SNAP_MANIFEST:="${MYPATH}/snapshot.yaml"}
: ${JOB_MANIFEST:="${MYPATH}/job-restore-snap.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc-restore-snap.yaml"}
: ${NAMESPACE:="default"}

echo "Deleting the snapshot"
assert_cmd kubectl delete --wait -f ${SNAP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting the main job"
assert_cmd kubectl delete --wait -f ${JOB_MANIFEST} --namespace ${NAMESPACE}

echo "Creating the cleanup job"
assert_cmd kubectl create -f ${JOB_CLEANUP_MANIFEST} --namespace ${NAMESPACE}
kubectl wait --for=condition=complete -f ${JOB_CLEANUP_MANIFEST} --timeout=2m  --namespace ${NAMESPACE}

echo "Deleting the cleanup job"
assert_cmd kubectl delete --wait -f ${JOB_CLEANUP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting the pvc"
assert_cmd kubectl delete --wait -f ${PVC_MANIFEST} --namespace ${NAMESPACE}
#kubectl wait --for=delete -f ${PVC_MANIFEST} --timeout=2m  --namespace ${NAMESPACE}

echo "Job delete completed"
exec_cmd kubectl get pv,pvc,volumesnapshot,volumesnapshotcontent,job --namespace ${NAMESPACE}
