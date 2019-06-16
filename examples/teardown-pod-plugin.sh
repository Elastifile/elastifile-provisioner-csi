#!/usr/bin/env bash

: ${NAMESPACE:="default"}

set -x
./delete-pod.sh $1
../deploy/teardown-plugin.sh
set +x

