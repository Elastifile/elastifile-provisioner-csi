#!/usr/bin/env bash

MYPATH=$(dirname $0)

SNAPSHOT_MANIFEST=$1
: ${SNAPSHOT_MANIFEST:="${MYPATH}/snapshot.yaml"}
: ${NAMESPACE:="default"}

kubectl create -f ${SNAPSHOT_MANIFEST} --namespace ${NAMESPACE}
