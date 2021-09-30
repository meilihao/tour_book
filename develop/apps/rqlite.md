# rqlite
参考:
- [rqlite 使用及注意事项](https://blog.arstercz.com/rqlite-%E4%BD%BF%E7%94%A8%E5%8F%8A%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A1%B9/)

rqlite 是一款轻量级且分布式的关系型数据库, 实际上它是在 SQLite 的基础上实现的关系型数据库, 加上 raft 协议实现了一致性分布式. 另外它还提供了 http 相关的接口. 不同于 etcd 和 consul 这两款键值数据库, rqlite 是真正的关系型数据库, 也支持事务等. 从这方面看, 在对数据有分布式一致性且支持事务功能的需求中, rqlite 非常适合, 而且很轻量, 方便部署, 不像 MySQL, PostgrelSQL 等显得有些笨重. 当然和传统数据库比起来 rqlite 也有自身的缺点.

> 实际上 raft 协议本身就要求了必须通过 leader 节点才能写入.

## 部署
- [creating-a-cluster](https://github.com/rqlite/rqlite/blob/master/DOC/CLUSTER_MGMT.md#creating-a-cluster)

```bash
# mkdir rqlite_demo && cd rqlite_demo
# mkdir data{1,2,3}
# rqlited -node-id 1 -http-addr localhost:4001 -raft-addr localhost:4002 ./data1
# rqlited -node-id 2 -http-addr localhost:4011 -raft-addr localhost:4012 -join http://localhost:4001 ./data2
# rqlited -node-id 3 -http-addr localhost:4021 -raft-addr localhost:4022 -join http://localhost:4001 ./data3
# curl localhost:4001/status?pretty # 获取cluster状态
# rqlite -H localhost # rqlite的命令行工具, 类似psql
```

## 性能
rqlite 默认使用 sqlite 的特性in-memory SQLite database 来尽可能的增加性能. 如果觉得性能不是那么重要, 可以在 rqlited 启动的时候增加 on-disk 选项使 rqlite 使用 file-based SQLite database 特性.

同样 rqlite 使用 in-memory 特性不会有数据有丢失的风险, 因为每个节点上的 raft 日志已经保存了相关的数据.