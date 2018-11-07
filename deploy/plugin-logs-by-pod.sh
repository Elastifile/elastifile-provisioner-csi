PODNAME=$1
kubectl logs ${PODNAME} -c csi-ecfsplugin
