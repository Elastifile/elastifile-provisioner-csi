#!/usr/bin/env bash

set -x
./delete-snapvol.sh $1
./delete-snapshot.sh $1
./delete-pod.sh $1
../deploy/teardown-plugin.sh
set +x

