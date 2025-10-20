# 常用命令
psql 中的元命令是指以反斜线开头的命令, psql 提供丰富的元命令, 能够便捷地管理数据库 

- `su postgres -c psql` = `sudo -u postgres psql`
- `sudo -u postgres psql`
- `psql -h 127.0.0.1 -p 5432 -U user -d dbname` : 连接数据库
- `psql -U user -d dbname -W` : 没有`-h`即连接数据库 by unix socket
- `psql -E ...` : 获取元命令对应的 SQL
- `psql -v xxx=yyy ...` : `-v`传递变量
- `\encoding [编码名称]` : 显示或设定用户端编码
- `\?` : 列出所有的元命令
- `\h [NAME]` : help
- `\q` : 退出 psql
- `\c dbname` : 切换数据库
- `\l` : 列举数据库
- `\dt` : 列举表
- `\d tblname` : 查看表结构
- `\di[+]` : 查看索引大小
- `\db+` : 查看表空间
- `\x` : 以列显示的开关
- `\timing on/off` : 显示执行时长
- `\conninfo` : 显示连接信息
- `\! date` : 执行shell命令, 主要命令前有空格
- `\i file` : 执行file中的sql, = `psql -s file`
- `\pset border 0/1/2：设置执行结果的边框样式` : 显示执行结果的边框
- `\c` : 查看当前数据库和用户. = `select current_user;`
- `\c db` : 进入指定db
- `\c db username` :切到某个db的某个角色
- `\dn` : 列出当前库下所有schema
- `\d` : 查看当前数据库下的所有表、视图和序列
- `\dt` : 只查看数据库中的所有表
- `\d tb_name` : 查看表结构定义
- `\dt+ tb_name` : 查看表大小等属性
- `\db[+]` : 查看表空间
- `\du` : 列出所有用户及其用户权限
- `\ds` : 查看用户自定义序列
- `\df` : 查看用户自定义函数
- `\sf xxx` : 查看函数定义
- `\timing` : 显示 SQL 执行时间
- `\watch [seconds]` : 反复执行当前查询缓冲区的SQL命令,直到 SQL 被中止或执行失败
- `\echo :PROMPT1` : 查看提示符格式
- `SHOW config_file;` : 显示postgresql.conf路径
- `SHOW data_directory;` : 显示pg data目录
- `SHOW unix_socket_directories` : 显示unix socket路径
- `SHOW hba_file;`

其他:
1. psql 支持箭头键上下翻历史 SQL 命令, 需要readline

script:
```bash
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$QUERY" # 在script中使用pg

QUERIES=("SELECT * FROM product;" "SELECT * FROM orders;")
for QUERY in "${QUERIES[@]}"
do
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$QUERY"
done

if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$QUERY"; then
    echo "Failed to execute query."
    exit 1
fi
```

## schema
```sql
CREATE SCHEMA sales [AUTHORIZATION sales_user]; --- 指定所有者
ALTER TABLE sales.orders SET SCHEMA archive; --- 挪动schema
```

## ddl
ref:
- [创建和删除索引](https://cloud.tencent.com/document/product/1129/39829)

```sql
--- table
alter table xxx drop column yyy;
alter table xxx add column yyy smallint default 0 not null;
--- index
create unique index  tbl_bb_index  on  tbl_bb(id,name);
drop index xxx;
--- constraint
alter table xxx drop constraint yyy;
```

> index和constraint是不同的, 删除命令有差异.

## 其他命令
- `select version()`
- `show autovacuum` : 查看auto vacuum状态

    `SELECT name, setting FROM pg_settings WHERE name LIKE 'autovacuum%';`查看auto vacuum参数:
    - `autovacuum_analyze_scale_factor = 0.05` : 当表的 5% 被修改时触发autovacuum分析

    可通过使用系统视图来监控清理活动并评估表的健康状况(比如autovacuum来不及, 需要手动vacuum): `SELECT relname, last_vacuum, last_autovacuum, n_dead_tup FROM pg_stat_user_tables WHERE n_dead_tup > 1000;`

    查看autovacuum状态: `SELECT pid, datname, relname, query, state FROM pg_stat_activity WHERE query LIKE '%autovacuum%';`(输出结果会包含本身)

    对于非常大的表，考虑进行分区以提高清理效率. 这将表分成更小、可管理的部分. 比如`CREATE TABLE sales.orders_part ( order_id INT, order_date DATE, total_due NUMERIC ) PARTITION BY RANGE (order_date);`
- `show config_file` : 查看配置文件
- `show hba_file`
- `show ident_file`
- `show all`: 查看所有pg配置参数或使用`select * from pg_settings;`
- `show archive_command` : 查看指定参数
- `show transaction_isolation;` : 查看隔离基本
- `vacuum test` : vacuum test表
- `vacuum full test` : vacuum test表, 是通过独占锁表, 并重写整个表来回收额外的空间
- `VACUUM (PARALLEL 2) test;` : vacuum test表, 并发2
- `VACUUM ANALYZE test;` : vacuum test表并分析表(ANALYZE会更新查询规划器的表统计信息, 应该能体现vacuum前后的查询性能)
- `select * from pg_stat_activity;`: 查看PostgreSQL 进程信息,每一个进程在视图中存在一条记录
- `select * from pg_locks where granted is not true;` : 查看锁等待信息
- `select name,setting from pg_settings where name in('synchronous_commit','synchronous_standby_names');` : 查看配置
- `select pg_size_pretty(pg_database_size(db_name)); ` : 查看db大小
- `select pg_database.datname, pg_size_pretty (pg_database_size(pg_database.datname)) AS size from pg_database;` : 查看所有数据库的大小
- `select pg_size_pretty(pg_relation_size(table_name))`: 查看表大小
- `alter table xxx owner to new_owner;` # 修改表owner

## 内置函数
```psql
# SELECT current_setting('server_version_num'); -- 查看server version
# select current_timestamp;
# \df pg_start_backup # 查看pg_start_backup函数
# select pg_xlogfile_name(pg_switch_xlog()); -- 切换wal并输出pg_xlogfile_name
# select pg_current_wal_lsn(); -- 获得当前wal日志写入位置 # pg10前的版本需要将函数中的wal替换为xlog
# select pg_walfile_name(pg_current_wal_lsn()); -- 转换wal日志位置为文件名
# select pg_walfile_name_offset(pg_current_wal_lsn()); -- 转换wal日志位置为文件名和偏移量
# select pg_wal_switch() : 切换当前wal日志并归档到归档目录
# SELECT pg_create_checkpoint() : 创建checkpoint, 触发wal归档
# SELECT * FROM pg_stat_archiver; # 查看wal归档情况
```

## 添加用户
```sql
postgres=# CREATE DATABASE mytestdb [owner postgres];
CREATE DATABASE
postgres=# CREATE USER mytestuser WITH ENCRYPTED PASSWORD '123456'; # 或`create user root with password 'password';`
CREATE ROLE
postgres=# GRANT ALL PRIVILEGES ON DATABASE mytestdb to mytestuser;
GRANT
ALTER DATABASE name OWNER TO new_owner; -- 修改db owner
ALTER TABLE <tablename> OWNER TO <username>; -- 修改table owner
postgres=# drop user xxx;
```

## sql结果保存到file
```sql
postgres=# \o /var/lib/pgsql/data/output.txt -- 前提: db user有对应目录保存文件的权限
postgres=# select * from dummy_table;
postgres=# \o
```

## 表空间
database在表空间里时, 其表,索引,序列等对象默认在该表空间里.

```bash
cat $PGDATA/tablespace_map # 查看表空间映射位置, from pg 9.5. 在10.5上是用$PGDATA/pg_tblspc, 但pg_basebackup备份的base.tar.gz里有tablespace_map
```

```sql
# create tablespace tbs_test owner postgres location '/usr/local/pgdata'; # 会在$PGDATA/pg_tblspc下有一个连接文件xxx, 指向/usr/local/pgdata. pgdata的所有者必须是postgres
# CREATE DATABASE logistics TABLESPACE ts_primary; -- 在表空间内建库
# select d.datname,p.spcname from pg_database d, pg_tablespace p where d.datname='lottu01' and p.oid = d.dattablespace; --查看dbname的默认表空间
# create table test(a int) tablespace tbs_test; --在表空间内建表
# \db[+] [<tablespace_name>] --罗列表空间, `+`表示更多细节, 比如空间大小
# select * from pg_tablespace; -- 查看表空间
# select pg_tablespace_location(16385); -- 查看表空间的location
# select pg_tablespace_size('pg_default'); -- 查看表空间大小
# select pg_relation_filepath('<表名/索引名/etc...>'); # 先进入database
# select spcname, pg_size_pretty(pg_tablespace_size(spcname)) from pg_tablespace; -- 查看各个表空间的大小
# alter table test_tsp03 set tablespace tsp01; -- 将表从一个表空间移到另一个表空间, 期间会锁表(在这个期间涉及到的对象将被锁定, 不可访问)
# drop tablespace if exists tbs_test; -- 删除表空间. 删除表空间前必须要删除该表空间下的所有数据库对象，否则无法删除
-- 为表和索引指定新的表空间
postgres=> ALTER TABLE foo SET TABLESPACE pg_default;
postgres=> ALTER INDEX foo_idx SET TABLESPACE pg_default;
-- 使用如下语句将一个表空间中的所有表或索引移至另一个表空间, 相应对象会被锁定, 直至完成
postgres=> ALTER TABLE ALL IN TABLESPACE myspace SET TABLESPACE pg_default;
postgres=> ALTER INDEX ALL IN TABLESPACE myspace SET TABLESPACE pg_default;
```

## 元数据
```sql
SELECT u.datname  FROM pg_catalog.pg_database u where u.datname='xxx'; # 检查是否存在数据库xxx
```

## dump
```bash
pg_dump -h localhost -U postgres testdb > db.sql
psql -h localhost -U postgres -d testdb < db.sql
```

## install
env: ubuntu 24.04

```bash
$ apt install postgresql postgresql-contrib
$ sudo sed -i '/^host/s/ident/md5/' /etc/postgresql/16/main/pg_hba.conf
$ sudo sed -i '/^local/s/peer/trust/' /etc/postgresql/16/main/pg_hba.conf
$ echo "host all all 0.0.0.0/0 md5" | sudo tee -a /etc/postgresql/16/main/pg_hba.conf
$ sudo systemctl restart postgresql
$ sudo -u postgres psql
postgres=# ALTER USER postgres PASSWORD 'postgres';
```