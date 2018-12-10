#!/usr/bin/env bash

kubectl delete -f pod-with-volume.yaml
./teardown-plugin.sh

