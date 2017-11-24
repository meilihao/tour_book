#!/bin/bash

# env
ESServerUrls=${ESServerUrls:-"/usr/bin/cqlsh"}

# start jaeger
/jaeger/cmd/collector 2>&1 > collector.log --span-storage.type=cassandra --es.server-urls="${ESServerUrls}" &
/jaeger/cmd/agent 2>&1 > agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
/jaeger/cmd/query 2>&1 > query.log --query.static-files=/jaeger/jaeger-ui/build/ --span-storage.type=cassandra --cassandra.keyspace="${KEYSPACE}"