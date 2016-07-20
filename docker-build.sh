#!/bin/sh

set -e
cd `dirname $0`


sudo rm -rf /tmp/docker-registry-viewer/docker-bin
docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer:ro \
    -v /tmp/docker-registry-viewer/docker-bin:/docker-bin mkdym/golang:alpine-git \
    /bin/sh -c "go get -v github.com/gin-gonic/gin && go build -o /docker-bin/docker-registry-viewer github.com/mkdym/docker-registry-viewer"

sudo rm -rf ./docker-bin
sudo mv /tmp/docker-registry-viewer/docker-bin/ ./docker-bin

docker build -t mkdym/docker-registry-viewer:v1.0 .
