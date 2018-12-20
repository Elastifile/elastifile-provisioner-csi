#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="snapvol-pod-mount.yaml"}

kubectl delete -f ${POD_MANIFEST}

