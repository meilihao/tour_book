# 常用命令
- `su postgres -c psql` = `sudo -u postgres psql`
- `sudo -u postgres psql`
- `psql -h 127.0.0.1 -p 5432 -U user -d dbname` : 连接数据库
- `\encoding [编码名称]` : 显示或设定用户端编码
- `\?` : help
- `\q` : 退出 psql
- `\c dbname` : 切换数据库
- `\l` : 列举数据库
- `\dt` : 列举表
- `\d tblname` : 查看表结构
- `\di` : 查看索引
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
- `\db` : 查看表空间
- `\du` : 列出所有用户及其用户权限
- `\ds` : 查看用户自定义序列
- `\df` : 查看用户自定义函数

## 其他命令
- `select version()`
- `show config_file` : 查看配置文件
- `show hba_file`
- `show ident_file`
- `show all`: 查看所有pg配置参数或使用`select * from pg_settings;`
- `show archive_command` : 查看指定参数
- `show transaction_isolation;` : 查看隔离基本
- `vacuum test` : vacuum test表
- `select * from pg_stat_activity;`: 查看会话
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
```

## 添加用户
```sql
postgres=# CREATE DATABASE mytestdb;
CREATE DATABASE
postgres=# CREATE USER mytestuser WITH ENCRYPTED PASSWORD '123456'; # 或`create user root with password 'password';`
CREATE ROLE
postgres=# GRANT ALL PRIVILEGES ON DATABASE mytestdb to mytestuser;
GRANT
ALTER DATABASE name OWNER TO new_owner; -- 修改db owner
ALTER TABLE <tablename> OWNER TO <username>; -- 修改table owner
```

## sql结果保存到file
```sql
postgres=# \o /var/lib/pgsql/data/output.txt -- 前提: db user有对应目录保存文件的权限
postgres=# select * from dummy_table;
postgres=# \o
```

## 表空间
```bash
cat $PGDATA/tablespace_map # 查看表空间映射位置, from pg 9.5. 在10.5上是用$PGDATA/pg_tblspc, 但pg_basebackup备份的base.tar.gz里有tablespace_map
```

```sql
# create tablespace tbs_test owner postgres location '/usr/local/pgdata'; # 会在$PGDATA/pg_tblspc下有一个连接文件xxx, 指向/usr/local/pgdata
# CREATE DATABASE logistics TABLESPACE ts_primary; -- 在表空间内建库
# create table test(a int) tablespace tbs_test; --在表空间内建表
# \db[+] [<tablespace_name>] --罗列表空间, `+`表示更多细节, 比如空间大小
# select * from pg_tablespace; -- 查看表空间
# select pg_tablespace_size('pg_default'); -- 查看表空间大小
# select spcname, pg_size_pretty(pg_tablespace_size(spcname)) from pg_tablespace; -- 查看各个表空间的大小
# alter table test_tsp03 set tablespace tsp01; -- 将表从一个表空间移到另一个表空间, 期间会锁表(在这个期间涉及到的对象将被锁定, 不可访问)
# drop tablespace if exists tbs_test; -- 删除表空间. 删除表空间前必须要删除该表空间下的所有数据库对象，否则无法删除
```

## 元数据
```sql
SELECT u.datname  FROM pg_catalog.pg_database u where u.datname='xxx'; # 检查是否存在数据库xxx
```