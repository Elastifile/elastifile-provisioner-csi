#!/usr/bin/env bash

# Usage example:
# MGMT_ADDR=10.11.209.228 NFS_ADDR=172.16.0.1 PLUGIN_TAG=v0.1.0 ./deploy-plugin-create-pod.sh
: ${NAMESPACE:="default"}

../deploy/deploy-plugin.sh
./create-pod.sh $1

