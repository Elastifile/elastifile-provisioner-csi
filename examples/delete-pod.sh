#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-io.yaml"}
: ${PVC_MANIFEST:="pvc.yaml"}
: ${POD_CLEANUP_MANIFEST:="pod-cleanup-data.yaml"}

echo "Deleting pod"
kubectl delete -f ${POD_MANIFEST}

echo "Waiting for the pod to be deleted - may take some time due to a large amount of I/O"
#kubectl wait --for=delete -f ${POD_MANIFEST} --timeout=2m
sleep 180

echo "Creating cleanup pod"
#kubectl wait --for=condition=Released -f ${PVC_MANIFEST} --timeout=0; kubectl create -f ${POD_CLEANUP_MANIFEST}
kubectl create -f ${POD_CLEANUP_MANIFEST}

echo "Waiting for cleanup pod to be gone"
#kubectl wait --for=delete -f ${POD_CLEANUP_MANIFEST} --timeout=2m
sleep 90
kubectl delete -f ${POD_CLEANUP_MANIFEST}

echo "Deleting pvc"
kubectl delete -f ${PVC_MANIFEST}
echo "Waiting for the pv/pvc to be deleted"
sleep 120

echo "Pod delete completed"
kubectl get pv,pvc,volumesnapshot,volumesnapshotcontent,pod

