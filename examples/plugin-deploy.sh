#!/bin/bash

deployment_base="${1}"

if [[ -z $deployment_base ]]; then
	deployment_base="../deploy"
fi

cd "$deployment_base" || exit 1

kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)

objects=(configmap secret csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac csi-ecfsplugin-attacher csi-ecfsplugin-provisioner storageclass csi-ecfsplugin)

for obj in ${objects[@]}; do
	kubectl create -f "./$obj.yaml"
done
