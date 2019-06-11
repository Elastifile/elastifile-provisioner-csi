# Elastifile's ECFS CSI provisioner

[Container Storage Interface (CSI)](https://github.com/container-storage-interface/) provisioner for Elastifile ECFS.

## Overview

Elastifile ECFS provisioner plugin implements an interface between CSI enabled Container Orchestrator (CO) and ECFS cluster.

The plugin allows dynamically provisioning ECFS volumes, creating volume snapshots, creating volumes from volume snapshots, and attaching them to workloads.

Current implementation of ECFS CSI plugin was tested in Kubernetes environment (requires Kubernetes 1.11+).

* For details about configuring and deploying the plugin, see [docs/deploy.md](docs/deploy.md).
* For example use of the plugin, e.g. creating a volume or a snapshot, see [docs/examples.md](docs/examples.md)
* For development information, see [docs/develop.md](docs/develop.md)
