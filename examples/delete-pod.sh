#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-io.yaml"}
: ${PVC_MANIFEST:="pvc.yaml"}
: ${POD_CLEANUP_MANIFEST:="pod-cleanup-data.yaml"}

kubectl delete -f ${POD_MANIFEST}
kubectl delete -f ${PVC_MANIFEST}
kubectl create -f ${POD_CLEANUP_MANIFEST}

