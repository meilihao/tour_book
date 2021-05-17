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

### `sudo systemctl status postgresql.service` failed

使用`dnf install postgresql-server`安装了Fedora22 repo默认提供的postgres9.4,但无法启动.
原因：postgres未初始化,启动前运行`postgresql-setup initdb`即可.

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