#!/usr/bin/env bash

: ${NAMESPACE:="default"}

CONTAINER_NAME=csi-ecfsplugin
set -x
kubectl logs -l app=${CONTAINER_NAME} -c ${CONTAINER_NAME} --namespace ${NAMESPACE}
set +x

