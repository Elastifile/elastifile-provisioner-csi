#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPVOL_MANIFEST=$1
: ${SNAPVOL_MANIFEST:="${MYPATH}/snapvol-pod-mount.yaml"}

kubectl create -f ${SNAPVOL_MANIFEST}
kubectl exec -it demo-snap-pod ls /mnt
