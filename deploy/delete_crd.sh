#!/bin/bash

MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

assert_cmd kubectl delete -f csidriver.yaml
assert_cmd kubectl delete -f csinodeinfo.yaml
