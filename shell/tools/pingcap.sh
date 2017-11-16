#!/usr/bin/env bash

pingcap='tidb-latest-linux-amd64'

echo 'start pd'
nohup $pingcap/bin/pd-server --data-dir=$pingcap/pd 2>&1 >$pingcap/pb.log &
sleep 5

echo 'start tikv'
nohup $pingcap/bin/tikv-server --pd="127.0.0.1:2379" --data-dir=$pingcap/tikv 2>&1 >$pingcap/tikv.log &
sleep 5

echo 'start tidb'
nohup $pingcap/bin/tidb-server --store=tikv --path="127.0.0.1:2379" 2>&1 >$pingcap/tidb.log &

echo 'start ok!'