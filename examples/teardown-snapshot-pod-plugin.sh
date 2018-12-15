#!/usr/bin/env bash

set -x
./delete-snapshot.sh $1
./delete-pod.sh $1
../deploy/teardown-plugin.sh
set +x

