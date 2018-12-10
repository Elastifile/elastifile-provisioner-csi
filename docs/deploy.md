# Deploying Elastifile's ECFS CSI provisioner

## Deployment with Kubernetes

Requires Kubernetes 1.11+

Your Kubernetes cluster must allow privileged pods (i.e. `--allow-privileged` flag must be set to true for both the API server and the kubelet). Moreover, as stated in the [mount propagation docs](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation), the Docker daemon of the cluster nodes must allow shared mounts.

YAML manifests are located under `deploy/`.

## Configuration

### Deploy plugin
```bash
PLUGIN_TAG=v0.1.0 MGMT_ADDR=10.10.10.10 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 ./deploy-plugin.sh
```

### Volume defaults
In `deploy/storageclass.yaml`, `parameters` such as `userMapping` can be customized to suit your persistent volume needs

### Teardown plugin
```bash
./teardown-plugin.sh
```

### Verifying the deployment in Kubernetes

After successfully completing the steps above, you should see output similar to this:
```bash
$ kubectl get pod,storageclass
NAME                               READY   STATUS    RESTARTS   AGE
pod/csi-ecfsplugin-attacher-0      1/1     Running   0          37s
pod/csi-ecfsplugin-provisioner-0   1/1     Running   0          35s
pod/csi-ecfsplugin-rvzz2           2/2     Running   0          31s
pod/csi-ecfsplugin-wkbhz           2/2     Running   0          31s
pod/csi-ecfsplugin-wkpxx           2/2     Running   0          31s

NAME                                             PROVISIONER            AGE
storageclass.storage.k8s.io/elastifile           csi-ecfsplugin         32s
storageclass.storage.k8s.io/standard (default)   kubernetes.io/gce-pd   3h
```

You can try deploying a demo pod from `examples/` to test the deployment further.
The recommended manifests are
* `pod-with-volume.yaml` - creates a pvc and a pod, which mounts the resulting volume on /mnt
* `pod-with-io.yaml` - similar to `pod-with-volume.yaml`, but also starts `dd` to generate some traffic, which can be observed via stats of ECFS management console

### Notes on volume deletion

Upon PVC deletion, ECFS Data Container is going to be deleted.
In case there's data in the Data Container, volume deletion will success, but the Data Container will be kept intact. 
