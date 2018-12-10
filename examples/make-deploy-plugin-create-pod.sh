#!/usr/bin/env bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-with-io.yaml"}

pushd ..
time make all
popd
./deploy-plugin-create-pod.sh ${POD_MANIFEST}
