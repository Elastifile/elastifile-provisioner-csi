FROM alpine

RUN apk add --update bash && apk add --update curl
RUN apk add --update libintl && apk add --virtual build_deps gettext && cp /usr/bin/envsubst /usr/local/bin/envsubst && apk del build_deps

COPY kubectl /usr/local/bin/kubectl
RUN mkdir /root/.kube && chmod +x /usr/local/bin/kubectl
