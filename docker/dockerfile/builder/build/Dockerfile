FROM ubuntu:20.04

MAINTAINER meilihao <563278383@qq.com>

WORKDIR /app

RUN set -x \
	&& sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && apt update \
    && apt install -y automake \
	&& apt autoremove && apt autoclean && rm -rf /var/lib/apt/lists/*

# COPY ./build.sh /app/build.sh
# ENTRYPOINT ["/app/build.sh"]