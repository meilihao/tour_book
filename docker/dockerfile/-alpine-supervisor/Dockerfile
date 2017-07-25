# sudo docker build . -t="tiny-env:latest"
# 不推荐使用,golang程序无法运行
FROM alpine:latest

MAINTAINER chenz057 <chenz057@zhixubao.com>

# ENV SUPERVISOR_ROOT=/etc/supervisor
# ENV SUPERVISOR_CONF_DIR=/etc/supervisor/conf.d
# ENV SUPERVISOR_CONF_FILE=/etc/supervisor/supervisord.conf

RUN echo -e 'https://mirrors.ustc.edu.cn/alpine/v3.6/main\nhttps://mirrors.ustc.edu.cn/alpine/v3.6/community' > /etc/apk/repositories
RUN apk add --no-cache supervisor
RUN rm /etc/supervisord.conf
RUN mkdir -p /etc/supervisor/conf.d

COPY supervisord.conf /etc/supervisor/supervisord.conf

