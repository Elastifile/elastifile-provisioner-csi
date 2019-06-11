FROM ubuntu:16.04
LABEL maintainers="jean.spector@elastifile.com"
LABEL description="ECFS CSI Plugin"

COPY ecfsplugin /ecfsplugin

# curl is needed to install root CA
RUN apt-get update && apt-get -y install nfs-common telnet curl

RUN chmod +x /ecfsplugin && \
    mkdir -p /var/log/elastifile

ENTRYPOINT ["/ecfsplugin"]

