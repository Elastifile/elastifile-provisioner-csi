#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPSHOT_MANIFEST=$1
: ${SNAPSHOT_MANIFEST:="${MYPATH}/snapshot.yaml"}

kubectl delete -f ${SNAPSHOT_MANIFEST}
