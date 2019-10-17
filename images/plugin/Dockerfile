FROM ubuntu:16.04
LABEL maintainers="jeanspector@google.com"
LABEL description="Elastifile CSI Plugin"

COPY ecfsplugin /ecfsplugin
COPY entrypoint.sh /entrypoint.sh

# curl is needed to install root CA
RUN apt-get update && \
    apt-get -y install nfs-common telnet curl

RUN chmod +x /ecfsplugin && \
    mkdir -p /var/log/elastifile

ENTRYPOINT ["/entrypoint.sh", "/ecfsplugin"]
