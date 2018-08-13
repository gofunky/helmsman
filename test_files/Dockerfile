#  This is a docker image for the helmsman test container
#  It can be pulled from gofunky/helmsman:test
FROM gofunky/kubernetes:v1.11.2 as kube

FROM golang:1.10.3-alpine3.7 as builder

ARG HELM_VERSION=v2.10.0-rc.3

COPY --from=kube /usr/local/bin/kubectl /usr/local/bin/kubectl

RUN apk add --update --no-cache ca-certificates git \
    && apk add --update -t deps curl tar gzip make bash \
    && curl -L https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar zxv -C /tmp \
    && mv /tmp/linux-amd64/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && rm -rf /tmp/linux-amd64
    
RUN go get github.com/BurntSushi/toml \ 
    && go get github.com/goreleaser/goreleaser 

WORKDIR src/helmsman
