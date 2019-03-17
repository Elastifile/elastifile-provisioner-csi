#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

POD_MANIFEST=$1
: ${POD_MANIFEST:="${MYPATH}/pod-clone.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc-clone.yaml"}
: ${POD_CLEANUP_MANIFEST:="${MYPATH}/pod-cleanup-data.yaml"}
: ${NAMESPACE:="default"}

echo "Deleting the main pod"
assert_cmd kubectl delete -f ${POD_MANIFEST} --namespace ${NAMESPACE}

echo "Creating the cleanup pod"
assert_cmd kubectl create -f ${POD_CLEANUP_MANIFEST} --namespace ${NAMESPACE}
#kubectl wait --for=condition=Ready -f ${POD_CLEANUP_MANIFEST} --timeout=2m  --namespace ${NAMESPACE} # Container skips Ready

echo "Waiting for the cleanup pod to complete"
MAX_ATTEMPTS=5
POD_CLEANUP_NAME=$(kubectl get -f ${POD_CLEANUP_MANIFEST} -o go-template='{{.metadata.name}}'  --namespace ${NAMESPACE})
for ((i = 0; i < $MAX_ATTEMPTS; i++)); do
    echo -n .
    POD_STATUS=$(kubectl get pod ${POD_CLEANUP_NAME} -o go-template='{{(index .status.containerStatuses 0).state.terminated.reason}}' --namespace ${NAMESPACE})
    if [[ "$POD_STATUS" == "Completed" ]]; then
        echo "Cleanup completed"
        break
    fi
    sleep 1
done

echo "Deleting the cleanup pod"
assert_cmd kubectl delete -f ${POD_CLEANUP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting pvc"
assert_cmd kubectl delete -f ${PVC_MANIFEST} --namespace ${NAMESPACE}
echo "Waiting for the pv/pvc to be deleted"
sleep 120

echo "Pod delete completed"
exec_cmd kubectl get pv,pvc,volumesnapshot,volumesnapshotcontent,pod --namespace ${NAMESPACE}
