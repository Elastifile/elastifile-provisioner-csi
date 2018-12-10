#!/usr/bin/env bash

set -x
kubectl delete -f pod-with-volume.yaml
../deploy/teardown-plugin.sh
set +x

