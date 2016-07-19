#!/bin/sh

set -e
cd `dirname $0`


sudo rm -rf /tmp/docker-registry-viewer/bin
docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer:ro \
    -v /tmp/docker-registry-viewer/bin:/go/bin mkdym/golang:alpine-git \
    /bin/sh -c "go get -v github.com/gin-gonic/gin && go install github.com/mkdym/docker-registry-viewer"

sudo rm -rf ./bin
sudo mv /tmp/docker-registry-viewer/bin/ ./bin
docker build -t docker-registry-viewer:v1.0 .
