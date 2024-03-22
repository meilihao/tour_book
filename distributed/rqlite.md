# rqlite
ref:
- [rqlite Design](https://rqlite.io/docs/design/)


## 部署

rqlite默认通过`http://localhost:4001`访问server.

### 单节点
server:
```bash
# cd rqlite-v8.23.0-linux-amd64
# mkdir -p data
# ./rqlited -node-id 1 data1 # node1
```

client:
```bash
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