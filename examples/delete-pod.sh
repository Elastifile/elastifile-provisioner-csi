#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-io.yaml"}
: ${PVC_MANIFEST:="pvc.yaml"}
: ${POD_CLEANUP_MANIFEST:="pod-cleanup-data.yaml"}

echo "Deleting the main pod"
kubectl delete -f ${POD_MANIFEST}

echo "Creating the cleanup pod"
kubectl create -f ${POD_CLEANUP_MANIFEST}
#kubectl wait --for=condition=Ready -f ${POD_CLEANUP_MANIFEST} --timeout=2m # Container skips Ready

echo "Waiting for the cleanup pod to complete"
POD_CLEANUP_NAME=$(kubectl get -f ${POD_CLEANUP_MANIFEST} -o go-template='{{.metadata.name}}')
i=0; while [[ $(kubectl get pod ${POD_CLEANUP_NAME} -o go-template='{{(index .status.containerStatuses 0).state.terminated.reason}}') != "Completed" ]]; do sleep 1; let i+=1; echo -n .; if [[ $i -ge 5 ]]; then echo -e "\nDone"; break; fi; done

echo "Deleting the cleanup pod"
kubectl delete -f ${POD_CLEANUP_MANIFEST}

echo "Deleting pvc"
kubectl delete -f ${PVC_MANIFEST}
echo "Waiting for the pv/pvc to be deleted"
sleep 120

echo "Pod delete completed"
kubectl get pv,pvc,volumesnapshot,volumesnapshotcontent,pod
