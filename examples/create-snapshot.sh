#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPSHOT_MANIFEST=$1
: ${SNAPSHOT_MANIFEST:="${MYPATH}/snapshot.yaml"}

kubectl create -f ${SNAPSHOT_MANIFEST}
