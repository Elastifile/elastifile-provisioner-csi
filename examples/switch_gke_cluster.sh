#!/bin/bash

: ${PROJECT:=elastifile-gce-lab-c340}
: ${CLUSTER:=cluster-1}
: ${ZONE:=europe-west1-b}


if [ "$1" == "poc" ]; then
    PROJECT=launcher-poc-207208
    CLUSTER=cluster-tmp
    ZONE=us-central1-a
fi

set -x

# Set project and zone
gcloud config set project ${PROJECT}
gcloud config set compute/zone ${ZONE}


# Set cluster
gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE"

# Configure docker client
gcloud auth configure-docker

# Give current user admin privileges
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)

set +x
