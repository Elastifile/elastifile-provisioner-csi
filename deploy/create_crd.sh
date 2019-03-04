#!/usr/bin/env bash

MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

assert_cmd kubectl create -f csidriver.yaml --validate=false
assert_cmd kubectl create -f csinodeinfo.yaml --validate=false
