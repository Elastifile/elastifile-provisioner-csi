#!/usr/bin/env bash

: ${POD_ID:=1}
: ${LABEL:=csi-ecfsplugin}
: ${CONTAINER_NAME:=csi-ecfsplugin}
: ${NAMESPACE:="default"}

POD_NAME=$(kubectl get pods -l app=${LABEL} -o=name --namespace ${NAMESPACE} | head -n ${POD_ID} | tail -n 1)

if [[ -z "${POD_NAME}" ]]; then
    echo "Failed to detect pod name in namespace ${NAMESPACE}"
    exit 1
fi

function get_pod_status() {
	echo -n $(kubectl get ${POD_NAME} -o jsonpath="{.status.phase}" --namespace ${NAMESPACE})
}

while [[ "$(get_pod_status)" != "Running" ]]; do
	sleep 1
	echo "Waiting for ${POD_NAME} (status $(get_pod_status))"
done

echo "Showing logs for pod ${POD_NAME}"
# kubectl logs -f ${POD_NAME} -c ${CONTAINER_NAME} --namespace ${NAMESPACE}
set -x
kubectl logs ${POD_NAME} -c ${CONTAINER_NAME} --namespace ${NAMESPACE} $@
set +x
