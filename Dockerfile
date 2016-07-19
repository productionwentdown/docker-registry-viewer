FROM alpine:3.3

MAINTAINER juxiaoheng "juxiaoheng@wps.cn"


COPY ./bin /docker-bin/
COPY ./resources /docker-bin/

EXPOSE 49110

WORKDIR /docker-bin/
ENTRYPOINT ["/docker-bin/docker-registry-viewer"]

