#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

POD_MANIFEST=$1
: ${POD_CLEANUP_MANIFEST:="${MYPATH}/pod-cleanup-data.yaml"}
: ${SNAP_MANIFEST:="${MYPATH}/snapshot.yaml"}
: ${POD_MANIFEST:="${MYPATH}/pod-restore-snap.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc-restore-snap.yaml"}
: ${NAMESPACE:="default"}

echo "Deleting the snapshot"
assert_cmd kubectl delete -f ${SNAP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting the main pod"
assert_cmd kubectl delete -f ${POD_MANIFEST} --namespace ${NAMESPACE}

echo "Creating the cleanup pod"
assert_cmd kubectl create -f ${POD_CLEANUP_MANIFEST} --namespace ${NAMESPACE}
#kubectl wait --for=condition=Ready -f ${POD_CLEANUP_MANIFEST} --timeout=2m  --namespace ${NAMESPACE} # Container skips Ready

echo "Waiting for the cleanup pod to complete"
POD_CLEANUP_NAME=$(kubectl get -f ${POD_CLEANUP_MANIFEST} -o go-template='{{.metadata.name}}'  --namespace ${NAMESPACE})
i=0; while [[ $(kubectl get pod ${POD_CLEANUP_NAME} -o go-template='{{(index .status.containerStatuses 0).state.terminated.reason}}' --namespace ${NAMESPACE}) != "Completed" ]]; do sleep 1; let i+=1; echo -n .; if [[ $i -ge 5 ]]; then echo -e "\nDone"; break; fi; done

echo "Deleting the cleanup pod"
assert_cmd kubectl delete -f ${POD_CLEANUP_MANIFEST} --namespace ${NAMESPACE}

echo "Deleting the pvc"
assert_cmd kubectl delete -f ${PVC_MANIFEST} --namespace ${NAMESPACE}

echo "Waiting for the pv/pvc to be deleted"
wait=1
while [[ ${wait} != 0 ]]; do
    kubectl get -f ${PVC_MANIFEST}
    if [[ $? != 0 ]]; then
        wait=0
    else
        sleep 1
    fi
done

echo "Pod delete completed"
exec_cmd kubectl get pv,pvc,volumesnapshot,volumesnapshotcontent,pod --namespace ${NAMESPACE}
