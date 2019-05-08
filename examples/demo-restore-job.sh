#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

assert_cmd ${MYPATH}/create-job.sh
assert_cmd ${MYPATH}/create-restore-job.sh
assert_cmd ${MYPATH}/delete-restore-job.sh
assert_cmd ${MYPATH}/delete-job.sh

