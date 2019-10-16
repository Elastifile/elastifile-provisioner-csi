# Elastifile CSI provisioner

[Container Storage Interface (CSI)](https://github.com/container-storage-interface/) plugin for Elastifile ECFS/EFAAS.

## Overview

Elastifile CSI plugin implements an interface between a CSI enabled Container Orchestrator (CO) and Elastifile cluster (SP).

The plugin allows dynamically provisioning Elastifile volumes, creating volume snapshots, creating volumes from volume snapshots, and attaching them to workloads.

This implementation of Elastifile CSI plugin was tested with Kubernetes environment 1.14.

Status: Beta

* For details about configuring and deploying the plugin with EMS, see [docs/deploy-ems.md](docs/deploy-ems.md).
* For details about configuring and deploying the plugin with eFaaS, see [docs/deploy-efaas.md](docs/deploy-efaas.md).
* For example use of the plugin, e.g. creating a volume or a snapshot, see [docs/examples.md](docs/examples.md)
* For development information, see [docs/develop.md](docs/develop.md)
