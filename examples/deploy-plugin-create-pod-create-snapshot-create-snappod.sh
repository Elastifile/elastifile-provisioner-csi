#!/usr/bin/env bash

# Usage examples:
# MGMT_ADDR=10.11.209.228 NFS_ADDR=172.16.0.1 PLUGIN_TAG=v0.1.0 ./deploy-plugin-create-pod-create-snappod.sh
# MGMT_ADDR=10.11.209.228 NFS_ADDR=172.16.0.1 PLUGIN_TAG=v0.1.0 SNAPSHOT_MANIFEST=snapshot.yaml SNAPSHOT_DELAY=30 ./deploy-plugin-create-pod-create-snappod.sh

: ${POD_MANIFEST:="pod-io.yaml"}
: ${SNAPSHOT_MANIFEST:="snapshot.yaml"}
: ${SNAPSHOT_DELAY:=60}
: ${SNAPPOD_MANIFEST:="snapvol-pod-mount.yaml"}

../deploy/deploy-plugin.sh
./create-pod.sh ${POD_MANIFEST}
echo Waiting for ${SNAPSHOT_DELAY} before taking a snapshot
sleep ${SNAPSHOT_DELAY}
./create-snapshot.sh ${SNAPSHOT_MANIFEST}
echo Waiting for ${SNAPSHOT_DELAY} before using the snapshot
sleep ${SNAPSHOT_DELAY}
./create-snappod.sh ${SNAPPOD_MANIFEST}
