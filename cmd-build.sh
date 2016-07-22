#!/bin/sh

set -x
set -e
cd `dirname $0`


COMMAND="export CGO_ENABLED=0 && go build -o /cmd-bin/regtool github.com/mkdym/docker-registry-viewer/cmd"


sudo rm -rf /tmp/docker-registry-viewer/cmd-bin
docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer:ro \
    -v /tmp/docker-registry-viewer/cmd-bin:/cmd-bin mkdym/golang:alpine-git \
    /bin/sh -c "$COMMAND"

sudo rm -rf ./cmd-bin
sudo mv /tmp/docker-registry-viewer/cmd-bin/ ./cmd-bin


