#!/bin/bash

CONTAINER_NAME=csi-ecfsplugin
kubectl logs -l app=${CONTAINER_NAME} -c ${CONTAINER_NAME}

