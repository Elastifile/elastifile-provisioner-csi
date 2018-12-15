#!/usr/bin/env bash

# Usage example:
# MGMT_ADDR=10.11.209.228 NFS_ADDR=172.16.0.1 PLUGIN_TAG=v0.1.0 ./make-deploy-plugin-create-pod.sh

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-with-io.yaml"}

pushd ..
time make all
popd
./deploy-plugin-create-pod.sh ${POD_MANIFEST}
