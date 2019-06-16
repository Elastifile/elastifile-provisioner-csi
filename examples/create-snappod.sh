#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPVOL_MANIFEST=$1
: ${SNAPVOL_MANIFEST:="${MYPATH}/snapvol-pod-mount.yaml"}
: ${NAMESPACE:="default"}

kubectl create -f ${SNAPVOL_MANIFEST} --namespace ${NAMESPACE}
kubectl exec -it demo-snap-pod --namespace ${NAMESPACE} ls /mnt
