#!/usr/bin/env bash

# Usage:
# GCP_PROJECT_ID=399 create_gke_cluster.sh

: ${K8S_VERSION:="1.13.6"}
# Alpha versions can be found at https://cloud.google.com/kubernetes-engine/docs/release-notes

: ${GKE_CLUSTER_NAME:="cluster-1"}
: ${GCP_REGION:="us-east1"}
: ${GCP_ZONE:="us-east1-b"}

: ${GCP_HOST_PROJECT_NAME:="elastifile-gce-lab-c934"} # Project that shares its VPC
: ${GCP_SERVICE_PROJECT_NAME:="elastifile-gce-lab-c946"} # Project that hosts the GKE cluster
: ${GKE_VPC:="shared"}
: ${GCP_SECONDARY_RANGE_CLUSTER:="cluster"} # Pods live here
: ${GCP_SECONDARY_RANGE_SERVICES:="service"} # Has to be present in the shared VPC's subnet

: ${GCP_NETWORK:="projects/${GCP_HOST_PROJECT_NAME}/global/networks/${GKE_VPC}"}
: ${GCP_SUBNETWORK:="projects/${GCP_HOST_PROJECT_NAME}/regions/${GCP_REGION}/subnetworks/${GKE_VPC}"}

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
ENABLE_ALPHA_CLUSTER="--no-enable-autoupgrade --no-enable-autorepair --enable-kubernetes-alpha"
gcloud beta container --project ${GCP_SERVICE_PROJECT_NAME} clusters create ${GKE_CLUSTER_NAME} --zone ${GCP_ZONE} --username ${GKE_CLUSTER_USER} --cluster-version ${K8S_VERSION} --machine-type ${GKE_NODE_VM} --image-type ${GKE_NODE_IMAGE} --disk-type ${GKE_NODE_DISK_TYPE} --disk-size ${GKE_NODE_DISK_SIZE} --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes ${GKE_NODE_COUNT} --enable-cloud-logging --enable-cloud-monitoring --network ${GCP_NETWORK} --subnetwork ${GCP_SUBNETWORK} --addons HorizontalPodAutoscaling,HttpLoadBalancing,KubernetesDashboard --no-issue-client-certificate --metadata disable-legacy-endpoints=true --enable-ip-alias  --cluster-secondary-range-name ${GCP_SECONDARY_RANGE_CLUSTER} --services-secondary-range-name ${GCP_SECONDARY_RANGE_SERVICES} --default-max-pods-per-node "110"
set +x

