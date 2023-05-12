# cnosdb
## example
```bash
# cnosdb singleton --config /etc/cnosdb/cnosdb.conf
# cnosdb-cli -P 31007 # connect to cnosdb.conf#http_listen_addr
> create database test;
> \c test
> DESCRIBE DATABASE oceanic_station;
> SHOW TABLES;
> CREATE TABLE t (
   k BIGINT,
   v STRING, -- 这个`,`不能省略
);
> DESCRIBE TABLE t;
> INSERT INTO t (time, k, v) VALUES (1666165200290401000, 56, 'test');
> select * from t;
```

> [cnosdb 2.2.0 rpm](https://www.cnosdb.com/download/)安装好后其bin没有可执行权限