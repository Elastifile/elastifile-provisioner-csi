#!/usr/bin/env bash

# Usage:
# GCP_PROJECT_ID=399 create_gke_cluster.sh

DEFAULT_GCP_REGION=$(gcloud config get-value compute/region) # e.g. europe-west1
DEFAULT_GCP_ZONE=$(gcloud config get-value compute/zone) # europe-west1-b

: ${GCP_PROJECT_ID:="340"}
: ${K8S_VERSION:="1.11.5"}

: ${GKE_CLUSTER_NAME:="cluster-1"}
: ${GCP_REGION:="${DEFAULT_GCP_REGION}"}
: ${GCP_ZONE:="${DEFAULT_GCP_ZONE}"}

: ${GCP_PROJECT_NAME:="elastifile-gce-lab-c"${GCP_PROJECT_ID}}
: ${GCP_NETWORK:="projects/${GCP_PROJECT_NAME}/global/networks/vpc-c${GCP_PROJECT_ID}"}
: ${GCP_SUBNETWORK:="projects/${GCP_PROJECT_NAME}/regions/${GCP_REGION}/subnetworks/vpc-c${GCP_PROJECT_ID}-subnet"}

: ${GKE_CLUSTER_USER:="admin"}
: ${GKE_NODE_VM:="custom-1-1536"}
: ${GKE_NODE_IMAGE:="COS"} # Can be COS|UBUNTU
: ${GKE_NODE_DISK_TYPE:="pd-standard"}
: ${GKE_NODE_DISK_SIZE:="50"} # GB
: ${GKE_NODE_COUNT:="3"}

if [[ "${GCP_ZONE}" != *"${GCP_REGION}"* ]]; then
  echo "Bad configuration - GCP zone and region mismatch (region: ${GCP_REGION} zone: ${GCP_ZONE})"
  exit 3
fi

set -x
EXPERIMENTAL_FLAGS="--no-issue-client-certificate --no-enable-ip-alias --metadata disable-legacy-endpoints=true"
gcloud beta container --project ${GCP_PROJECT_NAME} clusters create ${GKE_CLUSTER_NAME} --zone ${GCP_ZONE} --username ${GKE_CLUSTER_USER} --cluster-version ${K8S_VERSION} --machine-type ${GKE_NODE_VM} --image-type ${GKE_NODE_IMAGE} --disk-type ${GKE_NODE_DISK_TYPE} --disk-size ${GKE_NODE_DISK_SIZE} --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes ${GKE_NODE_COUNT} --enable-cloud-logging --enable-cloud-monitoring --network ${GCP_NETWORK} --subnetwork ${GCP_SUBNETWORK} --addons HorizontalPodAutoscaling,HttpLoadBalancing,KubernetesDashboard ${EXPERIMENTAL_FLAGS}
set +x
