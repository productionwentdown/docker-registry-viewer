#!/bin/sh

set -x
set -e
cd `dirname $0`


COMMAND="export CGO_ENABLED=0 && go build -o /docker-bin/docker-registry-viewer github.com/mkdym/docker-registry-viewer"


sudo rm -rf /tmp/docker-registry-viewer/docker-bin
docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer:ro \
    -v /tmp/docker-registry-viewer/docker-bin:/docker-bin golang:1.7-alpine \
    /bin/sh -c "$COMMAND"

sudo rm -rf ./docker-bin
sudo mv /tmp/docker-registry-viewer/docker-bin/ ./docker-bin

docker build -t mkdym/docker-registry-viewer:v1.1 .
