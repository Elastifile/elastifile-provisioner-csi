#!/usr/bin/env bash

CONTAINER_NAME=csi-ecfsplugin
set -x
kubectl logs -l app=${CONTAINER_NAME} -c ${CONTAINER_NAME}
set +x

