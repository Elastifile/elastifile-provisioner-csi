# `examples` directory includes manifests and scripts that demonstrate the use of Elastifile CSI plugin 

There are several helper scripts that can be used to deploy/teardown the configuration and make your life a bit easier in general by reducing the number of commands you have to type and demonstrating all parts of the provisioner use as well as end-to-end workflow.

## Deploy/create

* `../deploy/deploy-plugin.sh` - deploys the provisioner plugin

* `create-job.sh` - creates a pvc and a job that consumes the pvc and writes some data 

Demonstrates how to create a job and mount volume provisioned by Elastifile CSI plugin

* `create-restore-job.sh` - creates a snapshot on an existing volume, restores it to a new volume, then writes to a new file while reading fom an existing one 

Demonstrates how to create a volume based on an existing snapshot, and 

* `deploy-plugin-create-pod.sh` - calls `deploy-pugin.sh`, followed by `create-pod.sh`

Demonstrates end to end functionality - plugin deployment, pvc and pod creation, mounting the volume with generating some I/O

* `deploy-plugin-create-pod-create-snapshot.sh`

Demonstrates end to end functionality - plugin deployment, pvc and pod creation, mounting the volume with generating some I/O, then taking a snapshot

* `deploy-plugin-create-pod-create-snapshot-create-snappod.sh`

Demonstrates end to end functionality - plugin deployment, pvc and pod creation, mounting the volume with generating some I/O, taking a snapshot, creating a volume from the snapshot, and creating a pod mounting the new volume

`PLUGIN_TAG=v0.1.0 MGMT_ADDR=35.195.163.246 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 ./deploy-plugin-create-pod.sh`
* `make-deploy-plugin-create-pod.sh` - builds the plugin, then calls `deploy-plugin-create-pod.sh`.

Useful during plugin development to build and push the plugin images, followed by end-to-end flow resulting in I/O
* IMPORTANT: make sure to pass the relevant plugin settings to this script, e.g.
`PLUGIN_TAG=v0.1.0 MGMT_ADDR=35.195.163.246 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 ./deploy-plugin-create-pod.sh`

## Teardown

* `../deploy/teardown-plugin.sh` - removes the provisioner plugin

* `delete-job.sh` - deletes the job and the pvc (removes any data on the volume)

* `delete-restore-job.sh` - deletes the snapshot, the job and the pvc (removes any data on the volume)

* `delete-snapshot.sh` - deletes the snapshot (can't have an export)

* `delete-snappod.sh` - deletes the pvc created from a volume snapshot, and the pod that mounts it

* `teardown-pod-plugin.sh` - calls `delete-pod.sh`, followed by `teardown-plugin.sh`

* `teardown-snapshot-pod-plugin.sh` - calls `delete-snapshot.sh`, then `delete-pod.sh`, followed by `teardown-plugin.sh`

* `teardown-snappod-snapshot-pod-plugin.sh` - calls ./delete-snappod.sh, `delete-snapshot.sh`, `delete-pod.sh`, `teardown-plugin.sh`

## Troubleshooting

* `logs.sh` tails the output of the plugin container

* `logs-all.sh` shows the output of all the plugin containers, recommended to be used in combination with grep

* `exec-bash.sh` logs into the plugin's container and runs bash

## Miscellaneous

* `switch_gke_cluster.sh` - Changes the default K8s cluster/zone, gives the current user admin privileges on the cluster 

* `create_gke_cluster.sh` - Creates a GKE cluster for testing purposes. See inside the script for the values that can be configured
