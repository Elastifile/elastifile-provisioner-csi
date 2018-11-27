# Elastifile ECFS CSI 0.3.0

[Container Storage Interface (CSI)](https://github.com/container-storage-interface/) provisioner for Elastifile ECFS.

## Overview

Elastifile ECFS plugin implements an interface between CSI enabled Container Orchestrator (CO) and ECFS cluster.
It allows dynamically provisioning ECFS volumes and attaching them to workloads.
Current implementation of ECFS CSI plugin was tested in Kubernetes environment (requires Kubernetes 1.11+).

For details about configuration and deployment of the plugin, see documentation under `docs/`.

For example usage of the plugin, see examples under `examples/`.
