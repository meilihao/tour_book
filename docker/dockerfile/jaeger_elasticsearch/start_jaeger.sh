#!/bin/bash

# env
ESServerUrls=${ESServerUrls:-"http://localhost:9200"}

# start jaeger
cmd/collector 2>&1 > collector.log --span-storage.type=elasticsearch --es.server-urls="${ESServerUrls}" &
cmd/agent 2>&1 > agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
cmd/query 2>&1 > query.log --query.static-files=jaeger-ui/build/ --span-storage.type=elasticsearch --es.server-urls="${ESServerUrls}" &