#!/usr/bin/env bash

#: ${PROJECT:=elastifile-gce-lab-c906}
#: ${CLUSTER:=ekfs-cluster}
#: ${REGION:=us-central1}
#: ${ZONE:=us-central1-a}

: ${PROJECT:=elastifile-gce-lab-c945}
: ${CLUSTER:=ekfs-cluster}
: ${REGION:=us-central1}
: ${ZONE:=us-central1-a}

set -x

# Set project and zone
gcloud config set project ${PROJECT}
gcloud config set compute/region ${REGION}
gcloud config set compute/zone ${ZONE}

# Set cluster
gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE"

# Configure docker client
gcloud auth configure-docker

# Give current user admin privileges
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)

set +x
