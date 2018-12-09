# Deploy/teardown Elastifile ECFS provisioner on Kubernetes 1.11 

There are several helper scripts, that can be used to deploy/teardown the configuration and make your life a bit easier in general

## Deploy/create
* `deploy-plugin.sh` - deploys the provisioner plugin
* `create-pod.sh` - creates a pod, takes pod manifest as an optional argument 
Demonstrates how to create a pod using a volume created by ECFS provisioner
* `deploy-plugin-create-pod.sh` - deploys the plugin, then calls `create-pod.sh`
IMPORTANT: Make sure the secrets and the configmap manifests are updated before running this script
Demonstrates full functionality using a single command
* `make-deploy-plugin-create-pod.sh` - builds the plugin, then calls `deploy-plugin-create-pod.sh`.
Only useful during plugin development

## Teardown

* `teardown-plugin.sh` - removes the provisioner plugin
* `teardown-pod-plugin.sh` - deletes the pod, then calls `teardown-plugin.sh`

## Troubleshooting

* `logs.sh` tails the output of the plugin container
* `logs-all.sh` shows the output of all the plugin containers, recommended to be used in combination with grep
* `exec-bash.sh` logs into the plugin's container and runs bash

## Miscellaneous
* `switch_gke_cluster.sh` - Changes the default K8s cluster/zone, gives the current user admin privileges on the cluster 
