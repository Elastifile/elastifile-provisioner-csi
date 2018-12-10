#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-with-io.yaml"}

kubectl delete -f ${POD_MANIFEST}

