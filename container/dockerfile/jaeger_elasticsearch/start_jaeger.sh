#!/bin/bash
# 构建jaeger-ui需使用[releases](https://github.com/jaegertracing/jaeger-ui/tags)版本, 否则yarn build时可能报错.

export JaegerRoot=${JaegerRoot:-"/home/chen/opt/jaeger"}
export ESServerUrls=${ESServerUrls:-"http://localhost:9200"}

export SPAN_STORAGE_TYPE=elasticsearch # 只有设置SPAN_STORAGE_TYPE=elasticsearch后,collector和query才显示es的配置参数

# --- start jaeger, before v1.21.0
# nohup ${JaegerRoot}/cmd/collector 2>&1 > ${JaegerRoot}/cmd/collector.log --es.server-urls="${ESServerUrls}" &
# nohup ${JaegerRoot}/cmd/agent 2>&1 > ${JaegerRoot}/cmd/agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
# nohup ${JaegerRoot}/cmd/query 2>&1 > ${JaegerRoot}/cmd/query.log --query.static-files=${JaegerRoot}/jaeger-ui/build/ --es.server-urls="${ESServerUrls}" &

# --- start jaeger, v1.21.0, work with opentelemetry-collector
nohup ${JaegerRoot}/jaeger-collector 2>&1 > ${JaegerRoot}/collector.log --es.server-urls="${ESServerUrls}" --collector.grpc-server.host-port=":14250" &
nohup ${JaegerRoot}/jaeger-query 2>&1 > ${JaegerRoot}/query.log --query.static-files=${JaegerRoot}/jaeger-ui/build/ --es.server-urls="${ESServerUrls}" &

# jaeger-collector --collector.grpc.tls.enabled只能是true. 默认不设置即为grpc insecure=true

