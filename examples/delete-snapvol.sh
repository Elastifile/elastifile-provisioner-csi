#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPVOL_MANIFEST=$1
: ${SNAPVOL_MANIFEST:="${MYPATH}/volume-from-snapshot.yaml"}

kubectl delete -f ${SNAPVOL_MANIFEST}
