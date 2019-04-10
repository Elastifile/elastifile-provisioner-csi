#!/usr/bin/env bash

MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

# Obtained from https://github.com/kubernetes/kubernetes/blob/release-1.13/cluster/addons/storage-crds
# Not needed as of K8s 1.14
assert_cmd kubectl create -f csidriver.yaml --validate=false
assert_cmd kubectl create -f csinodeinfo.yaml --validate=false
