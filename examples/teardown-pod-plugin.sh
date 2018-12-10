#!/usr/bin/env bash

set -x
./delete-pod.sh $1
../deploy/teardown-plugin.sh
set +x

