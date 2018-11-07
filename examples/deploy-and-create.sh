./plugin-deploy.sh
kubectl create -f pod-with-volume.yaml
sleep 5
kubectl get pod,pvc
