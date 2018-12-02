#!/bin/bash

POD_MANIFEST=$1
: ${POD_MANIFEST:="pod-with-io.yaml"}

./deploy-plugin.sh
./create-pod.sh $1

