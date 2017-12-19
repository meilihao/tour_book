# sudo docker build . -t="tiny-env:latest"
# golang程序请使用`CGO_ENABLED=0`
FROM alpine:latest

MAINTAINER chenz057 <chenz057@zhixubao.com>

ENV SUPERVISOR_CONF_FILE=/etc/supervisord.conf
ENV SUPERVISOR_CONF_DIR=/etc/supervisor.d

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
RUN echo -e 'https://mirrors.ustc.edu.cn/alpine/v3.7/main\nhttps://mirrors.ustc.edu.cn/alpine/v3.7/community' > /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata supervisor \
    && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && mkdir -p $SUPERVISOR_CONF_DIR

CMD ["supervisord", "--nodaemon", "--configuration", "/etc/supervisord.conf"]