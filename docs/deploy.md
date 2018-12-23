# Deploying Elastifile's ECFS CSI provisioner

## Deployment requirements

Requires Kubernetes 1.11+

Your Kubernetes cluster must allow privileged pods (i.e. `--allow-privileged` flag must be set to true for both the API server and the kubelet). Moreover, as stated in the [mount propagation docs](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation), the Docker daemon of the cluster nodes must allow shared mounts.

`kubectl` should be available in $PATH and configured to point to the K8s cluster to which you're interested in deploying the provisioner

`envsubst` should be available in $PATH

Deployment scripts and YAML manifests are located under [deploy](../deploy) directory, and the rest of the document assumes that this is where you're located

## Configuration

### Deploy plugin
```bash
PLUGIN_TAG=v0.2.0 MGMT_ADDR=10.10.10.10 MGMT_USER=admin MGMT_PASS=Y2hhbmdlbWU= NFS_ADDR=10.255.255.1 K8S_USER=user@example.com ./deploy-plugin.sh
```

These values may be set by the user:
* PLUGIN_TAG - The version of the Elastifile ECFS CSI Provisioner you're interested in
* MGMT_ADDR - The IP address or DNS name you use to connect to Elastifile ECFS Management Console with
* MGMT_USER - The username you log into Elastifile ECFS Management Console with
* MGMT_PASS - The password for $MGMT_USER (base64 encoded)
* NFS_ADDR - The address you use to mount your Elastifile ECFS instance
* K8S_USER - The username with permissions to administer your Kubernetes cluster 

Some the above variables are optional:
* MGMT_USER defaults to admin
* K8S_USER default behavior is to use your gcloud credentials if 'gcloud' binary is available, $USER otherwise

### Volume defaults
If you're interested in tweaking volume creation defaults, please edit [storagelass.yaml](../deploy/storageclass.yaml) to suit you needs.

Each value under `parameters` that is expected to be modified by the user, e.g. `userMapping`, has a comment explaining its meaning and - where applicable - the values it takes.

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

You can deploy a demo pod from `examples/` to test the deployment further.
The recommended manifests are
* `pod-with-volume.yaml` - creates a pvc and a pod, which mounts the resulting volume on /mnt
* `pod-with-io.yaml` - similar to `pod-with-volume.yaml`, but also starts `dd` to generate some traffic, which can be observed via stats of ECFS management console

### Notes on volume deletion

Upon PVC deletion, ECFS Data Container is going to be deleted.

In case there's data in the Data Container, it will be kept intact to prevent accidental data loss. 

## Further reading

For examples on the use of the plugin, see [docs/examples.md](../docs/examples.md)
