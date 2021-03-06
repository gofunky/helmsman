FROM gofunky/kubernetes:v1.11.2 as kube

FROM golang:1.10.3 as builder
WORKDIR /go/src/
RUN go get github.com/BurntSushi/toml
RUN go get gopkg.in/yaml.v2
RUN git clone https://github.com/gofunky/helmsman.git
RUN go get github.com/gofunky/helmsman/gcs
RUN go get github.com/gofunky/helmsman/aws
#  build a statically linked binary so that it works on stripped linux images such as alpine/busybox.
RUN cd helmsman \
    && LastTag=$(git describe --abbrev=0 --tags) \
    && TAG=$LastTag-$(date +"%d%m%y") \
    && LT_SHA=$(git rev-parse ${LastTag}^{}) \
    && LC_SHA=$(git rev-parse HEAD) \
    && if [ ${LT_SHA} != ${LC_SHA} ]; then TAG=latest-$(date +"%d%m%y"); fi \
    && CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-X main.version='$TAG' -extldflags "-static"' .

# The image to keep
FROM alpine:3.8 
RUN apk add --update --no-cache ca-certificates git

ARG HELM_VERSION=v2.10.0-rc.3

RUN apk --no-cache update \
    && rm -rf /var/cache/apk/* \
    && apk add --update -t deps curl tar gzip make bash \
    && curl -L http://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar zxv -C /tmp \
    && mv /tmp/linux-amd64/helm /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && mkdir -p ~/.helm/plugins \
    && helm plugin install https://github.com/hypnoglow/helm-s3.git \
    && helm plugin install https://github.com/nouney/helm-gcs \
    && rm -rf /tmp/linux-amd64

COPY --from=kube /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=builder /go/bin/helmsman   /usr/local/bin/helmsman

WORKDIR /tmp
# ENTRYPOINT ["/bin/helmsman"]
