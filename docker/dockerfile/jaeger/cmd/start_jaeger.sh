#!/bin/bash

# start cassandra : from https://github.com/docker-library/cassandra/blob/ca3c9df03cab318d34377bba0610c741253b0466/3.11/Dockerfile
/docker-entrypoint.sh cassandra -f 2>&1 > /cassandra-schema/cassandra.log &

# test cassandra is ok : from https://github.com/uber/jaeger/blob/master/plugin/storage/cassandra/schema/docker.sh
CQLSH=${CQLSH:-"/usr/bin/cqlsh"}
CQLSH_HOST=${CQLSH_HOST:-"cassandra"}
CASSANDRA_WAIT_TIMEOUT=${CASSANDRA_WAIT_TIMEOUT:-"60"}
DATACENTER=${DATACENTER:-"dc1"}
KEYSPACE=${KEYSPACE:-"jaeger_v1_${DATACENTER}"}
MODE=${MODE:-"test"}

total_wait=0
while true
do
  ${CQLSH} -e "describe keyspaces"
  if (( $? == 0 )); then
    break
  else
    if (( total_wait >= ${CASSANDRA_WAIT_TIMEOUT} )); then
      echo "Timed out waiting for Cassandra."
      exit 1
    fi
    echo "Cassandra is still not up at ${CQLSH_HOST}. Waiting 1 second."
    sleep 1s
    ((total_wait++))
  fi
done

echo "Generating the schema for the keyspace ${KEYSPACE} and datacenter ${DATACENTER}"

# inid cassandra
MODE="${MODE}" DATACENTER="${DATACENTER}" KEYSPACE="${KEYSPACE}" /cassandra-schema/create.sh | ${CQLSH}

# start jaeger
/jaeger/cmd/collector 2>&1 > collector.log --span-storage.type=cassandra --cassandra.keyspace="${KEYSPACE}" &
/jaeger/cmd/agent 2>&1 > agent.log --http-server.host-port=:6831 --collector.host-port=127.0.0.1:14267 &
/jaeger/cmd/query 2>&1 > query.log --query.static-files=/jaeger/jaeger-ui/build/ --span-storage.type=cassandra --cassandra.keyspace="${KEYSPACE}"