#!/bin/bash

deployment_base="${1}"

if [[ -z $deployment_base ]]; then
	deployment_base="../deploy"
fi

cd "$deployment_base" || exit 1

objects=(csi-ecfsplugin-attacher csi-ecfsplugin-provisioner csi-ecfsplugin storageclass csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac configmap secret)

for obj in ${objects[@]}; do
	kubectl delete -f "./$obj.yaml"
done
