#!/usr/bin/env bash

MYPATH=$(dirname $0)

POD_MANIFEST=$1
: ${POD_MANIFEST:="${MYPATH}/pod-io.yaml}"
: ${PVC_MANIFEST:="${MYPATH}/pvc.yaml}"

kubectl create -f ${PVC_MANIFEST}
kubectl create -f ${POD_MANIFEST}

