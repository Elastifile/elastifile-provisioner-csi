
${MANIFEST}=$1
: ${MANIFEST:="pod-with-io.yaml"}

pushd ..
time make all
popd
./plugin-deploy.sh
kubectl create -f ${MANIFEST}
echo Examine the cluster state with the following command:
echo kubectl get pod,pvc

