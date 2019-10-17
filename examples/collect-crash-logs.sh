#!/usr/bin/env bash

: ${NAMESPACE:="elastifile-csi"}
: ${LOGDIR:="/tmp/crash-logs-"$(date +%s)}
: ${TARBALL:="/tmp/crash-logs.tgz"}
: ${LOGSCRIPT:="${LOGDIR}/fetch-logs.sh"}

mkdir -p ${LOGDIR}
pushd ${LOGDIR}

export NAMESPACE
kubectl get pod -n ${NAMESPACE} -o go-template='{{range .items}}{{$podName := .metadata.name}}{{range .status.containerStatuses}}{{if ge .restartCount 1}}{{print "kubectl logs -p " $podName " -c " .name " -n $NAMESPACE > " $podName "--" .name ".log\n"}}{{end}}{{end}}{{end}}' > ${LOGSCRIPT}

bash -x ${LOGSCRIPT}

cd ..
tar czvf ${TARBALL} $(basename ${LOGDIR})

popd

rm -rf ${LOGDIR}

echo Done - logs were saved as ${TARBALL}
