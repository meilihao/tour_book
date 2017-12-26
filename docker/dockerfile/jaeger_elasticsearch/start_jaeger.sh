#!/bin/bash

# env
JaegerRoot=${JaegerRoot:-"/app/jaeger"}
ESServerUrls=${ESServerUrls:-"http://localhost:9200"}

# start jaeger
${JaegerRoot}/cmd/collector 2>&1 > ${JaegerRoot}/cmd/collector.log --span-storage.type=elasticsearch --es.server-urls="${ESServerUrls}" &
${JaegerRoot}/cmd/agent 2>&1 > ${JaegerRoot}/cmd/agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
${JaegerRoot}/cmd/query 2>&1 > ${JaegerRoot}/cmd/query.log --query.static-files=${JaegerRoot}/jaeger-ui/build/ --span-storage.type=elasticsearch --es.server-urls="${ESServerUrls}"
