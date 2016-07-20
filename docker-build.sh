#!/bin/sh

set -x
set -e
cd `dirname $0`


COMMAND_GET="go get -v github.com/gin-gonic/gin"
COMMAND_GET_GOPM="gopm get -v -g github.com/gin-gonic/gin"
COMMAND_BUILD="go build -o /docker-bin/docker-registry-viewer github.com/mkdym/docker-registry-viewer"

COMMAND="${COMMAND_GET} && ${COMMAND_BUILD}"


sudo rm -rf /tmp/docker-registry-viewer/docker-bin
docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer:ro \
    -v /tmp/docker-registry-viewer/docker-bin:/docker-bin mkdym/golang:alpine-git \
    /bin/sh -c "$COMMAND"

sudo rm -rf ./docker-bin
sudo mv /tmp/docker-registry-viewer/docker-bin/ ./docker-bin

docker build -t mkdym/docker-registry-viewer:v1.0 .
