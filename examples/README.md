# Deploy/teardown Elastifile ECFS provisioner on Kubernetes 1.11 

There are several helper scripts that can be used to deploy/teardown the configuration and make your life a bit easier in general by reducing the number of commands you have to type and demonstrating the provisioner usage end-to-end

## Deploy/create

* `../deploy/deploy-plugin.sh` - deploys the provisioner plugin

* `create-pod.sh` - creates a pod, takes pod manifest as an optional argument 

Demonstrates how to create a pod and mount volume provisioned by ECFS plugin
* `deploy-plugin-create-pod.sh` - calls `deploy-pugin.sh`, followed by `create-pod.sh`

Demonstrates end to end functionality - plugin deployment, pvc and pod creation, mounting the volume with generating some I/O - using a single command\*

`PLUGIN_TAG=v0.1.0 MGMT_ADDR=35.195.163.246 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 ./deploy-plugin-create-pod.sh`
* `make-deploy-plugin-create-pod.sh` - builds the plugin, then calls `deploy-plugin-create-pod.sh`.

Useful during plugin development to build and push the plugin images, followed by end-to-end flow resulting in I/O\*

\* IMPORTANT: make sure to pass the relevant plugin settings to this script, e.g.
`PLUGIN_TAG=v0.1.0 MGMT_ADDR=35.195.163.246 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 ./deploy-plugin-create-pod.sh`

## Teardown

* `../deploy/teardown-plugin.sh` - removes the provisioner plugin

* `delete-pod.sh` - deletes the pod and the pvc

* `teardown-pod-plugin.sh` - calls `delete-pod.sh`, followed by `teardown-plugin.sh`

## Troubleshooting

* `logs.sh` tails the output of the plugin container

* `logs-all.sh` shows the output of all the plugin containers, recommended to be used in combination with grep

* `exec-bash.sh` logs into the plugin's container and runs bash

## Miscellaneous

* `switch_gke_cluster.sh` - Changes the default K8s cluster/zone, gives the current user admin privileges on the cluster 
