#!/bin/bash

# var
JaegerRoot=${JaegerRoot:-"/home/chen/opt/jaeger"}
ESServerUrls=${ESServerUrls:-"http://localhost:9200"}

export SPAN_STORAGE_TYPE=elasticsearch # 只有设置SPAN_STORAGE_TYPE=elasticsearch后,collector和query才显示es的配置参数

# start jaeger
nohup ${JaegerRoot}/cmd/collector 2>&1 > ${JaegerRoot}/cmd/collector.log --es.server-urls="${ESServerUrls}" &
nohup ${JaegerRoot}/cmd/agent 2>&1 > ${JaegerRoot}/cmd/agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
nohup ${JaegerRoot}/cmd/query 2>&1 > ${JaegerRoot}/cmd/query.log --query.static-files=${JaegerRoot}/jaeger-ui/build/ --es.server-urls="${ESServerUrls}" &