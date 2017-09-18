# docker build . -t "gobuilder:1.8.3"
FROM golang:1.8.3-stretch

MAINTAINER meilihao <563278383@qq.com>

RUN apt-get update && apt-get install -y --no-install-recommends  ca-certificates && rm -rf /var/lib/apt/lists/*

COPY ./build.sh /build.sh

ENTRYPOINT ["/build.sh"]

