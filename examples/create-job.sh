#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../deploy/functions.sh

: ${JOB_MANIFEST:="${MYPATH}/job-io.yaml"}
: ${PVC_MANIFEST:="${MYPATH}/pvc.yaml"}
: ${NAMESPACE:="default"}

assert_cmd kubectl create -f ${PVC_MANIFEST} --namespace ${NAMESPACE}
assert_cmd kubectl create -f ${JOB_MANIFEST} --namespace ${NAMESPACE}

echo "Waiting for the job to complete"
assert_cmd kubectl wait --for=condition=complete -f ${JOB_MANIFEST} --timeout=2m --namespace ${NAMESPACE}

JOB_NAME=$(kubectl get -f ${JOB_MANIFEST} -o go-template='{{.metadata.name}}')
echo "Job logs:"
kubectl logs job.batch/${JOB_NAME}
