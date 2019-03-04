#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="snapvol-pod-mount.yaml"}
: ${NAMESPACE:="default"}

kubectl delete -f ${POD_MANIFEST} --namespace ${NAMESPACE}

