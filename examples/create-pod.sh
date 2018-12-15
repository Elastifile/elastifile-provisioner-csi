#!/usr/bin/env bash

MYPATH=$(dirname $0)

POD_MANIFEST=$1
: ${POD_MANIFEST:="${MYPATH}/pod-with-io.yaml"}

kubectl create -f ${POD_MANIFEST}

