#!/usr/bin/env bash

# Usage:
# GCP_PROJECT_ID=399 create_gke_cluster.sh

: ${K8S_VERSION:="1.14.6"}
# Alpha versions can be found at https://cloud.google.com/kubernetes-engine/docs/release-notes

: ${GKE_CLUSTER_NAME:="cluster-1"}
: ${GCP_REGION:="us-east1"}
: ${GCP_ZONE:="us-east1-b"}

: ${GCP_PROJECT_NAME:="golden-eagle-dev-consumer10"}
: ${GKE_VPC:="default"}
: ${GCP_NETWORK:="projects/${GCP_PROJECT_NAME}/global/networks/${GKE_VPC}"}
: ${GCP_SUBNETWORK:="projects/${GCP_PROJECT_NAME}/regions/${GCP_REGION}/subnetworks/${GKE_VPC}"}

: ${GKE_CLUSTER_USER:="admin"}
: ${GKE_NODE_VM:="n1-standard-1"}
: ${GKE_NODE_IMAGE:="COS"} # Can be COS|UBUNTU
: ${GKE_NODE_DISK_TYPE:="pd-standard"}
: ${GKE_NODE_DISK_SIZE:="50"} # GB
: ${GKE_NODE_COUNT:="3"}

if [[ "${GCP_ZONE}" != *"${GCP_REGION}"* ]]; then
  echo "Bad configuration - GCP zone and region mismatch (region: ${GCP_REGION} zone: ${GCP_ZONE})"
  exit 3
fi

set -x
EXPERIMENTAL_FLAGS="--no-issue-client-certificate --no-enable-ip-alias --metadata disable-legacy-endpoints=true --no-enable-autoupgrade --no-enable-autorepair --enable-kubernetes-alpha"
gcloud beta container --project ${GCP_PROJECT_NAME} clusters create ${GKE_CLUSTER_NAME} --zone ${GCP_ZONE} --username ${GKE_CLUSTER_USER} --cluster-version ${K8S_VERSION} --machine-type ${GKE_NODE_VM} --image-type ${GKE_NODE_IMAGE} --disk-type ${GKE_NODE_DISK_TYPE} --disk-size ${GKE_NODE_DISK_SIZE} --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes ${GKE_NODE_COUNT} --enable-stackdriver-kubernetes --network ${GCP_NETWORK} --subnetwork ${GCP_SUBNETWORK} --addons HorizontalPodAutoscaling,HttpLoadBalancing ${EXPERIMENTAL_FLAGS}
set +x
