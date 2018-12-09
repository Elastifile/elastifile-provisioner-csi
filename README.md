# Elastifile's ECFS CSI provisioner

[Container Storage Interface (CSI)](https://github.com/container-storage-interface/) provisioner for Elastifile ECFS.

## Overview

Elastifile ECFS provisioner plugin implements an interface between CSI enabled Container Orchestrator (CO) and ECFS cluster.
It allows dynamically provisioning ECFS volumes and attaching them to workloads.
Current implementation of ECFS CSI plugin was tested in Kubernetes environment (requires Kubernetes 1.11+).

* For details about configuration and deployment of the plugin, see documentation under [docs/deploy.md](docs/deploy.md).

* For example usage of the plugin, see examples under [examples/README.md](examples/README.md)

* For development information, see [docs/develop.md](docs/develop.md)

* Deployment manifests are located under [deploy/](deploy)
