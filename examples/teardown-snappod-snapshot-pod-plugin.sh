#!/usr/bin/env bash

set -x
./delete-snappod.sh $1
./delete-snapshot.sh $1
echo Waiting for snapshot to be deleted in Ealstifile and for his fact to be registered by K8s
sleep 180
./delete-pod.sh $1
../deploy/teardown-plugin.sh
set +x

