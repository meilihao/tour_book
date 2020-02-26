# sudo docker build . -t="ubuntu-env:latest"
FROM ubuntu:17.04

MAINTAINER chenz057 <chenz057@zhixubao.com>

# ENV SUPERVISOR_ROOT=/etc/supervisor
# ENV SUPERVISOR_CONF_DIR=/etc/supervisor/conf.d
# ENV SUPERVISOR_CONF_FILE=/etc/supervisor/supervisord.conf

# COPY 01proxy /etc/apt/apt.conf.d/ # 设置apt proxy, 也可用`RUN export http_proxy=http://xxx && export https_proxy=https://xxx`处理

RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y supervisor && \
#   apt-get install -y iproute2 && \
    rm -rf /var/lib/apt/lists/*
# RUN rm /etc/supervisord.conf
# RUN mkdir -p /etc/supervisor/conf.d

# COPY supervisord.conf /etc/supervisor/supervisord.conf