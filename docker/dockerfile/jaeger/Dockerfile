# docker build . -t "jaeger:201709"
# docker run --name jaeger -d -p6831:6831/udp -p5778:5778 -p16686:16686 -p14268:14268 jaeger:201709
# issue :
# 1. sh脚本启动的进程无法收到signal
FROM cassandra:latest

# Agent zipkin.thrift compact
#EXPOSE 5775/udp

# Agent jaeger.thrift compact
EXPOSE 6831/udp

# Agent jaeger.thrift binary
#EXPOSE 6832/udp

# Agent config HTTP
EXPOSE 5778

# Collector HTTP
EXPOSE 14268

# Web HTTP
EXPOSE 16686

COPY ./schema /cassandra-schema/
COPY ./cmd /jaeger/cmd
COPY ./jaeger-ui/build/ /jaeger/jaeger-ui/build

VOLUME /var/lib/cassandra

# start jaeger
CMD ["/jaeger/cmd/start_jaeger.sh"]

# ⋊> ~/jaeger tree                                                                                                                                                                               13:35:19
# .
# ├── cmd
# │   ├── agent
# │   ├── collector
# │   ├── query
# │   └── start_jaeger.sh
# ├── Dockerfile
# ├── jaeger-ui
# │   └── build
# │       ├── asset-manifest.json
# │       ├── favicon.ico
# │       ├── index.html
# │       └── static
# │           ├── css
# │           │   ├── main.8789ae43.css
# │           │   └── main.8789ae43.css.map
# │           ├── js
# │           │   ├── 1.acaf8887.chunk.js
# │           │   ├── 1.acaf8887.chunk.js.map
# │           │   ├── main.685a3c52.js
# │           │   └── main.685a3c52.js.map
# │           └── media
# │               ├── flags.9c74e172.png
# │               ├── icons.674f50d2.eot
# │               ├── icons.912ec66d.svg
# │               ├── icons.af7ae505.woff2
# │               ├── icons.b06871f2.ttf
# │               ├── icons.fee66e71.woff
# │               └── jaeger-logo.c596f5b8.svg
# └── schema
#     ├── create.sh
#     ├── docker.sh
#     └── v001.cql.tmpl