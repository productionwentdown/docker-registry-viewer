#!/bin/sh

set -e
cd `dirname $0`


docker run -ti --rm \
    -v `pwd`:/go/src/github.com/mkdym/docker-registry-viewer \
    -v `pwd`/bin:/go/bin mkdym/golang:alpine-git \
    /bin/sh -c "go get -v github.com/gin-gonic/gin && go install github.com/mkdym/docker-registry-viewer"
docker build -t docker-registry-viewer:v1.0 .
