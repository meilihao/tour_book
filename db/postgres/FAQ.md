# FQQ

### systemctl status postgresql-9.4.service提示`Failed to start SYSV: Starts and stop`

估计是未用systemctl启动导致.

先停止postgres(`pg_ctl stop -D /opt/PostgreSQL/9.4/data`),再用`sudo systemctl start postgresql-9.4.service`重新启动即可.

### 大小写

创建database/table时,postgres会自动将数据库名,表名和字段名转为小写;执行sql时也会将表名和字段名转为小写.

### 安装部署postgres 13
参考:
- [安装文档](https://www.postgresql.org/download/)

启动pg: `sudo systemctl start postgresql@13-main.service`或`pg_ctlcluster 13 main start`
修改postgres.conf: `vim /etc/postgresql/13/main/postgresql.conf`

## 常用命令

psql:

- psql登入 : `psql -U user_name`
- 列出所有的数据库 : `\l`
- 切换数据库 : `\c db_name`
- 列出当前数据库下的数据表 : `\d`
- 列出指定表的所有字段 : `\d tb_name`
- 查看指定表的基本情况 : `\d+ tb_name`
- psql登出 ： `\q`

postgres:

- 启动数据库 : `pg_ctl start -D $PGDATA`,$PGDATA是postgres的数据目录
- 停止数据库 : `pg_ctl stop -D $PGDATA [-m SHUTDOWN-MODE]`,SHUTDOWN-MODE:smart(默认),等所有连接中止才关闭,如果client连接不中止,则无法关闭数据库;fast,主动断开client连接,回滚已有事务,然后正常关闭;immediate,立即停止数据库进程,直接退出,下次启动数据库需进行crash-recovery操作.

###`psql: 致命错误:  用户 "postgres" peer 认证失败`

原因：postgres中缺省仅存在 postgres 用戶,且受制于pg_hba.conf中对Local采用peer方式验证用戶身份，具体原因可查看pg的日志.

- 使用postgres账户登入即可(pg9.4 默认规则是`local all all peer`)，即`sudo -u postgres psql`.
- 修改数据库实例的认证文件，比如`/var/lib/pgsql/data/pg_hba.conf`,再重启pg，根据新规则进行登入即可.

### pg_hba.conf
参考:
- [PostgreSQL 中的客户端认证](https://scarletsky.github.io/2017/04/26/client-authentication-in-postgresql/)
- [postgres认证相关文档](http://postgres.cn/docs/9.5/auth-methods.html#AUTH-PASSWORD)
- [pg_hba.conf文件](http://www.postgres.cn/docs/9.5/auth-pg-hba-conf.html)

pg根据pg_hba.conf匹配client的授权方式. postgres会根据pg_hba.conf中规则出现的位置，从上到下依次匹配.规则中的METHOD是指定如何处理客户端的认证方式, 匹配格式:`# TYPE DATABASE USER ADDRESS METHOD [option]`:

TYPE 连接类型，表示允许用哪些方式连接数据库，它允许以下几个值：
- local 通过 Unix socket 的方式连接, 因此没有ADDRESS信息
- host 通过 TCP/IP 的方式连接，它能匹配 SSL 和 non-SSL 连接
- hostssl 只允许 SSL 连接
- hostnossl 只允许 non-SSL 连接

> local/host实际就是连接数据库时的client ip类型, 不然`psql`的`-h`参数

DATABASE 可连接的数据库，它有以下几个特殊值：
- all 匹配所有数据库
- sameuser 可连接和用户名相同的数据库
- samerole 可连接和角色名相同的数据库
- replication 允许复制连接，用于集群环境下的数据库同步
- 除了上面这些特殊值之外，还可以写特定的数据库，可以用逗号 (,) 来分割多个数据库

USER 可连接数据库的用户，值有三种写法：
- all 匹配所有用户
- 特定数据库用户名
- 特定数据库用户组，需要在前面加上 + (如：+admin)

ADDRESS 可连接数据库的地址，有以下几种形式：
- all 匹配所有 IP 地址
- samehost 匹配该服务器的 IP 地址
- samenet 匹配该服务器子网下的 IP 地址
- ipaddress/netmask (如：172.20.143.89⁄32)，支持 IPv4 与 IPv6
- 如果上面几种形式都匹配不上，就会被当成是 hostname. 注意: 只有 host, hostssl, hostnossl 会应用该字段

METHOD 连接数据库时的认证方式，常见的有几个特殊值：
- trust 无条件通过认证, **建议不要在生产环境中使用**
- reject 无条件拒绝认证
- md5 用 md5 加密密码进行认证, 如果你不使用ident，最好使用md5.密码是以md5形式传送给数据库，较安全，且不需建立同名的操作系统用户
- password 用明文密码进行认证，不建议在不信任的网络中使用
- ident 从一个 ident 服务器 (RFC1413) 获得客户端的操作系统用户名并且用它作为被允许的数据库用户名来认证，只能用在 TCP/IP 的类型中 (即 host, hostssl, hostnossl)

    用户映射文件为pg_ident.conf，这个文件记录着与操作系统用户匹配的数据库用户，如果某操作系统用户在本文件中没有映射用户，则默认的映射数据库用户与操作系统用户同名.
- peer 从内核获得客户端的操作系统用户名并把它用作被允许的数据库用户名来认证，只能用于本地连接 (即 local)
- pam   使用操作系统提供的可插入认证模块服务(PAM)来认证

> ident 和 peer 都要求客户端操作系统中存在对应的用户
> 只有 md5 和 password 是需要密码的，其他方式都不需要输入密码认证
> DATABASE指定多个数据库时以逗号分隔.**`all`只有在没有其他的符合条目时才代表“所有”**，因为`all`的优先级最低.
> 可在psql中用`SHOW hba_file;`查找pg_hba.conf, 该文件有可能不存在需自行创建.
> pg11的pg_hba.conf在`/etc/postgresql/11/main/pg_hba.conf`
> pg11的数据在`/var/lib/postgresql/11/main`

密码登录example:
```conf
local bareos bareos md5
host all all all md5
```

### `sudo systemctl status postgresql.service` failed

使用`dnf install postgresql-server`安装了Fedora22 repo默认提供的postgres9.4,但无法启动.
原因：postgres未初始化,启动前运行`postgresql-setup initdb`即可.

> kylinv10上的pg10使用`postgresql-setup --initdb`

### psql: fe_sendauth: no password supplied
新安装的pg, 其postgres用户没有密码

解决:
```
> sudo -u postgres psql # 进入psql
> alter user postgres with password 'postgres'; # 为postgres创建密码
> psql -h localhost -p 5432 -U postgres -W # 使用密码登录
```

### postgres日志
`postgresql.conf`日志相关选项:
- `logging_collector = on` : 开启日志, 需重启pg
- `log_directory='pg_log' : 日志保存位置`${PGDATA}/pg_log`
- `log_statement = 'none' # none, ddl, mod, all`

    - none不记录
    - ddl记录所有数据定义命令，比如CREATE,ALTER,和DROP 语句
    - mod记录所有ddl语句,加上数据修改语句INSERT,UPDATE等
    - all记录所有执行的语句，将此配置设置为all可跟踪整个数据库执行的SQL语句
- `log_connections = off` : 是否记录连接日志
- `log_disconnections = off` : 是否记录连接断开日志

### 查找PGDATA

- `ps -aux|grep postgres`,查看pg的启动参数`-D`
- PGDATA环境变量

### 执行计划

使用explain,如`explain select * from table_name`.

### 显示NULL值
pgadmin3: `File->Options->Query tool->Results grid->Show NULL values as <NULL>`
psql: `\pset null NULL`

### 获取主键

使用["INSERT...RETURNING"](http://postgres.cn/docs/9.4/sql-insert.html),RETURNING后加想要返回的列名,加"*"表示返回所有列.

### Application Stack Builder 目标pg的配置

/etc/postgres-reg.ini

### sql日志

$PGDATA/postgresql.conf
```
log_statement = "all"
```

### 导出/导入

```
// 导出到文件xxx.sql
pg_dump -h localhost -U postgres -f xxx.sql
// 导入到newdatabase
psql -d newdatabase -U postgres -f mydatabase.sql
```

### pgadmin 安装
方式:
1. 官方docker镜像
2. Python Wheel

Python Wheel方式:
```shell
$ sudo pip3 install ./pgadmin4-3.2-py2.py3-none-any.whl
$ cd /usr/local/lib/python3.6/dist-packages/pgadmin4
$ sudo python3 ./pgAdmin4.py
```

### jsonb数组追加
```sql
-- log : '[]'::jsonb
UPDATE "logs" SET "log" = log || '{"kind":1,"text":"xxx"}'
```

### 修改serial起始值
```sql
alter sequence channel_id_seq restart with 5;
```

### duplicate key value violates unique constraint
明明`\d xxx_id_seq;`有内容, `select currval('xxx_id_seq');`却报错.

```sql
SELECT setval('tablename_id_seq', (SELECT MAX(id) FROM tablename)+1) -- 修正后`select currval`可用
```
serial key其实是由sequence实现的，当手动给serial列赋值的时候，sequence是不会自增, 因此不要给serial手工赋值

### 位操作
```sql
```

### `pq: relation "${table_name}" does not exist`即该表不存在
- connection string 错误
- sql的字段需要引号/需要指定schema

### [`~/.pgpass`密码文件](https://www.runoob.com/manual/PostgreSQL/libpq-pgpass.html)的使用
格式: `hostname:port:database:username:password`, 比如`192.168.0.102:5432:postgres:postgres:rootroot`

`.pgpass`的权限必须是`600`

### the database was initialized with lc_collate "zh_CN.UTF-8:, which is not recognized by setlocale().
`localedef -f UTF-8 -i zh_CN zh_CN.UTF-8`并再重启postgresql后

### postgresql.auto.conf
从PostgreSQL 9.4版本开始，新引入了postgresql.auto.conf配置文件，它作为postgresql.conf文件的补充，在参数配置格式上，它和postgresql.conf保持一致, 并优先于postgresql.conf, 能被`ALTER SYSTEM`语句修改.

## 模拟操作
### 插入可产生约2G wal日志的数据
```psql
create table t1(a int);
insert into t1 values (generate_series(1,10000000));
insert into t1 values (generate_series(1,10000000));
insert into t1 values (generate_series(1,10000000));
insert into t1 values (generate_series(1,10000000));
```

## 参数
### [postgresql.conf](https://postgresqlco.nf/doc/zh/param/max_replication_slots/)
- default_transaction_isolation='repeatable read' : 默认事务隔离级别

    在线修改方法:
    ```sql
    # -- 修改transaction_isolation必须在事务中进行
    # begin;
    # set default_transaction_isolation='repeatable read';
    # show default_transaction_isolation;
    # begin;
    # set transaction isolation level serializable;
    # show transaction_isolation;
    ```
- hot_standby = on  : 在备份时同时允许查询
- max_wal_senders = 10 : 需要设置为一个大于0的数，它表示主库最多可以有多少个并发的standby数据库

    一个备库通常只消耗一个wal_sender进程
    pg_basebackup也会消耗一个wal_sender进程

- wal_sender_timeout = 60s: 流复制超时时间

    中断那些停止活动超过指定毫秒数的复制连接. 这对发送服务器检测一个后备机崩溃或网络中断有用, 设置为0将禁用该超时机制.
- wal_keep_segments = 64 : 将为standby库保留64个WAL日志文件. 应当设置为一个尽量大的值, 以便备库落后主库时可通过主库保留的wal进行追回, 但是需要考虑磁盘空间(即归档空间)允许，一个WAL日志文件的大小是16M.
- max_standby_streaming_delay = 30s : 可选, 流复制最大延迟
- wal_receiver_status_interval = 10s : 可选，从向主报告状态的最大间隔时间
- hot_standby_feedback = on : 可选，出现错误复制, 向主机反馈

    设置hot_standby_feedback参数之后备库会定期向主库通知最小活跃事务id(xmin)，这样使得主库vacuum进程不会清理大于xmin值的事务。但是假如主备之间的网络突然中断，备库就无法向主库正常发送xmin值，如果时间够长，主库在这段时间内还是会清理无用元组，这样网络恢复后就可能发生冲突ERROR：canceling statement due to confilct with recovery.
- synchronous_commit : on代表同步复制；off则是异步，也是默认模式；除此中间还有几个模式也跟刷盘策略有关
- synchronous_standby_names: 配置**同步复制**的备库列表. 主库的这个值必须和同步备库recovery.conf的primary_conninfo参数的application_name设置一致

    假设synchronous_standby_names='*', 备库有两台, 那么查看pg_stat_replication时, 一台是sync, 另一台是potential.
- max_replication_slots = 10 : 为使用replication slot，必须大于0；replication slot作用是保证wal没有同步到standby之前不能从pg_xlog移走

    ref:
    - [postgresql之replication slot](https://blog.csdn.net/qq_35462323/article/details/106810187)

    replication slots 是从 postgresql 9.4 引入的，主要是提供了一种自动化的方法来确保主控机在所有的后备机收到 WAL 段之前不会移除它们，并且主控机也不会移除可能导致恢复冲突的行，即使后备机断开也是如此.

    为了防止 WAL 被删，SLOT restart_lsn 之后的WAL文件都不会删除. (wal_keep_segments 则是一个人为设定的边界, slot 是自动设定的边界(无上限), 所以使用 slot 并且下游断开后, 可能导致数据库的 WAL 堆积爆满).

    ```psql
    select pg_create_physical_replication_slot(‘gp1_a_slot’);  #创建replication slot
    select * from pg_replication_slots;                        #查询创建的replication slot
    ```

# recovery.conf
- standby_mode : 标记PG为STANDBY SERVER
- primary_conninfo : 标识主库信息
- trigger_file : 标识触发器文件

## 工具
### pg_rewind
ref:
- [PostgreSQL使用pg_rewind快速增量重建备库](https://bbs.huaweicloud.com/blogs/215956)
- [PgSQL · 特性分析 · 神奇的pg_rewind](http://mysql.taobao.org/monthly/2018/05/05/)


在常见的PostgreSQL双节点高可用构架中, 如果主库挂了且主备无延迟, 高可用系统会提升老备库为新主库对外服务. 而对于老主库, 则可以有很多处理策略, 例如：
1. 删掉，重搭新备库

    很显然，相比来说这种不是个很好的方案。 当数据量比较大时，重搭备库的时间成本太高，系统的可用性降低.
1. 降级为备库，继续服务

    因为老的主库挂掉的原因多种多样，甚至有可能是高可用系统的误判，而老主库也有可能是在挂掉之后又重新作为主库启动起来，这个时候降级并重搭流复制关系的操作就有可能失败（新的备库比新主库数据更超前）.
    为了解决这种情况，PostgreSQL 引入了pg_rewind工具.


    关闭pg会有一个checkpoint的点. 备库同步完成后，两边的数据是同样的一致性状态. 这样原来的主很容易就能做新主的备库. 但是如果在执行的过程中忘记了关闭主库，主库一直处于运行状态，那么这个旧主和新主在数据时间线上就不是一致的了, 就会导致后续搭建备库失败.

    这个命令就是通过新主去同步旧主，使这两个库处于一致性的状态。类似于PG旧主的一次向前回滚.

    它是为了防止库级别太大，重搭建库过于麻烦而引入的.

pg_rewind应用场景:
- 当2个pg实例时间线（timeline）出现分叉时，pg_rewind可以做实例间的同步, 比如主备切换后，老主库仍然运行，导致主备时间线不一致，老主库无法当做新主库的备库启动.

pg_rewind注意事项:
- pg_rewind运行过程中，会对比主（源）备（目标）的差异点，并把主库的差异点后的WAL日志传递给备库。所以，如果主库在差异点之后的WAL也丢失了，那么rewind是不会拷贝不存在的WAL日志的，所以此时备库仍然不会被成功做成standby。解决该问题需要用restore
- 需要使用超级用户

在[PostgreSQL 官方文档的介绍](https://www.postgresql.org/docs/10/static/app-pgrewind.html)中，pg_rewind 不光可以针对上面说的failover 的场景，理论上, 它可以使基于同一集群复制的任意两个副本（集群）进行同步.

pg_rewind 工具主要实现了从源集群到目的集群的文件级别数据同步, 即仅复制产生变化的数据块和一些文件：新数据文件、配置文件、WAL segments。和rsync的区别是，pg_rewind 不需要去读那些未变化的文件块，当数据量比较大而变化较小的时候，pg_rewind会更快.

值的注意的是，pg_rewind为了能够支持文件级别的数据同步，两个集群都打开如下参数：

- initdb初始化时启用`--data-checksums` 或者 wal_log_hints=on，参数说明详见[文档](https://www.postgresql.org/docs/10/static/runtime-config-wal.html#GUC-FULL-PAGE-WRITES)
- full_page_writes=on，参数说明详见[文档](https://www.postgresql.org/docs/10/static/runtime-config-wal.html#GUC-FULL-PAGE-WRITES)

以上几个参数打开，才能够保证通过WAL 日志恢复出来的数据是完整的，一致的，从而才能够实现文件级别的数据同步. 其实这2个参数设置的目的就是方便pg_rewind快速完成数据不一致的修复.

pg_rewind参数:
- -D directory --target-pgdata=directory

    此选项指定与源同步的目标数据目录。在运行pg_rewind之前，**必须干净关闭目标服务器**

- --source-pgdata=directory

    指定要与之同步的源服务器的数据目录的文件系统路径。**此选项要求干净关闭源服务器**

- --source-server=connstr

    指定要连接到源PostgreSQL服务器的libpq连接字符串。连接必须是具有超级用户访问权限的正常(非复制)连接。此选项要求**源服务器正在运行，而不是处于恢复模式**

- -n --dry-run

    除了实际修改目标目录之外，执行所有操作。

- -P --progress

    输出进展报告

- --debug

    输出很多Debug的信息。如果失败时，可以用此选项定位错误原因

pg_rewind的排除目录见[这里](https://github.com/postgres/postgres/blob/REL_12_STABLE/src/bin/pg_rewind/filemap.c).

