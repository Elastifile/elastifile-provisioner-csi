FROM ubuntu:16.04
LABEL maintainers="Kubernetes Authors"
LABEL description="ECFS CSI Plugin"

COPY ecfsplugin /ecfsplugin

RUN apt-get update && apt-get -y install nfs-common telnet

RUN chmod +x /ecfsplugin && \
    mkdir -p /var/log/elastifile

ENTRYPOINT ["/ecfsplugin"]

