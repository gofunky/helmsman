#  This is a docker image for the helmsman test container
#  It can be pulled from gofunky/helmsman:test
FROM golang:1.10.3-alpine3.7 as builder

ENV KUBE_LATEST_VERSION v1.11.0
ENV HELM_VERSION v2.10.0-rc.3

RUN apk add --update ca-certificates \
    && apk add --update -t deps curl \
    && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl 

RUN apk add --update --no-cache ca-certificates git \
    && apk add --update -t deps curl tar gzip make bash \
    && curl -L https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar zxv -C /tmp \
    && mv /tmp/linux-amd64/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && rm -rf /tmp/linux-amd64
    
RUN go get github.com/BurntSushi/toml \ 
    && go get github.com/goreleaser/goreleaser 

WORKDIR src/helmsman
