# rqlite
参考:
- [rqlite 使用及注意事项](https://blog.arstercz.com/rqlite-%E4%BD%BF%E7%94%A8%E5%8F%8A%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A1%B9/)
- [rqlite Design](https://rqlite.io/docs/design/)

rqlite 是一款轻量级且分布式的关系型数据库, 实际上它是在 SQLite 的基础上实现的关系型数据库, 加上 raft 协议实现了一致性分布式. 另外它还提供了 http 相关的接口. 不同于 etcd 和 consul 这两款键值数据库, rqlite 是真正的关系型数据库, 也支持事务等. 从这方面看, 在对数据有分布式一致性且支持事务功能的需求中, rqlite 非常适合, 而且很轻量, 方便部署, 不像 MySQL, PostgrelSQL 等显得有些笨重. 当然和传统数据库比起来 rqlite 也有自身的缺点.

> 实际上 raft 协议本身就要求了必须通过 leader 节点才能写入.

## 部署
ref:
- [creating-a-cluster](https://github.com/rqlite/rqlite/blob/master/DOC/CLUSTER_MGMT.md#creating-a-cluster)

rqlite默认通过`http://localhost:4001`访问server.

### 单节点
server:
```bash
# cd rqlite-v8.23.0-linux-amd64
# mkdir -p data1
# ./rqlited -node-id1 data1 # node1
```

client:
```bash
# curl localhost:4001/status?pretty # 获取cluster状态
# ./rqlite
Welcome to the rqlite CLI.
Enter ".help" for usage hints.
Connected to http://127.0.0.1:4001 running version v8.23.0
127.0.0.1:4001> CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)
0 row affected
127.0.0.1:4001> INSERT INTO foo(name) VALUES("fiona")
1 row affected
127.0.0.1:4001> SELECT * FROM foo
+----+-------+
| id | name  |
+----+-------+
| 1  | fiona |
+----+-------+
127.0.0.1:4001> .status
```

### cluster
在单节点基础上, 追加节点node2

```bash
# ./rqlited -node-id 2 -http-addr localhost:4003 -raft-addr localhost:4004 -join localhost:4002 data2 # node2
```

此时rqlite执行`INSERT INTO foo(name) VALUES("fiona3")`可以成功, 因为leader是node1. 此时ctrl+c node2, 再次插入报`leader not found`

再次启动node2后, 再次插入报成功. 因为`>=2/2+1`.


再追加node3:
```bash
# ./rqlited -node-id 3 -http-addr localhost:4005 -raft-addr localhost:4006 -join localhost:4002 data3
```

此时停止node2, node3. rqlite插入/查询均报`leader not found`.

再启动node2, rqlite插入/查询均正常了.

## [构建](https://rqlite.io/docs/install-rqlite/building-from-source/)

## 性能
rqlite 默认使用 sqlite 的特性in-memory SQLite database 来尽可能的增加性能. 如果觉得性能不是那么重要, 可以在 rqlited 启动的时候增加 on-disk 选项使 rqlite 使用 file-based SQLite database 特性.

同样 rqlite 使用 in-memory 特性不会有数据有丢失的风险, 因为每个节点上的 raft 日志已经保存了相关的数据.