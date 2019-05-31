### systemctl status postgresql-9.4.service提示`Failed to start SYSV: Starts and stop`

估计是未用systemctl启动导致.

先停止postgres(`pg_ctl stop -D /opt/PostgreSQL/9.4/data`),再用`sudo systemctl start postgresql-9.4.service`重新启动即可.

### 大小写

创建database/table时,postgres会自动将数据库名,表名和字段名转为小写;执行sql时也会将表名和字段名转为小写.

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

>[PostgreSQL 中的客户端认证](https://scarletsky.github.io/2017/04/26/client-authentication-in-postgresql/)
>[postgres认证相关文档](http://postgres.cn/docs/9.5/auth-methods.html#AUTH-PASSWORD)
>[pg_hba.conf文件](http://www.postgres.cn/docs/9.5/auth-pg-hba-conf.html)
>postgres会根据pg_hba.conf中规则出现的位置，从上到下依次匹配.规则中的METHOD是指定如何处理客户端的认证,常用的有ident，md5，password，trust，reject，pam:
>- ident是Linux下PostgreSQL默认的local认证方式，凡是能正确登录服务器的操作系统用户（注：不是数据库用户）就能使用本用户映射的数据库用户不需密码登录数据库。用户映射文件为pg_ident.conf，这个文件记录着与操作系统用户匹配的数据库用户，如果某操作系统用户在本文件中没有映射用户，则默认的映射数据库用户与操作系统用户同名。比如，服务器上有名为user1的操作系统用户，同时数据库上也有同名的数据库用户，user1登录操作系统后可以直接输入psql，以user1数据库用户身份登录数据库且不需密码。很多初学者都会遇到psql -U username登录数据库却出现“username ident 认证失败”的错误，明明数据库用户已经createuser。原因就在于此，使用了ident认证方式，却没有同名的操作系统用户或没有相应的映射用户。解决方案：1、在pg_ident.conf中添加映射用户；2、改变认证方式。
>- md5是常用的密码认证方式，如果你不使用ident，最好使用md5.密码是以md5形式传送给数据库，较安全，且不需建立同名的操作系统用户
>- password是以明文密码传送给数据库，建议不要在生产环境中使用
>- trust是只要知道数据库用户名就不需要密码或ident就能登录，建议不要在生产环境中使用
>- reject是拒绝认证
>- pam 	使用操作系统提供的可插入认证模块服务(PAM)来认证

```
// DATABASE指定多个数据库时以逗号分隔.`all`只有在没有其他的符合条目时才代表“所有”，因为`all`的优先级最低.
local   all             postgres                                peer
local   all             all                                     md5
```

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
> alter user postgres with password 'postgres' # 为postgres创建密码
> psql -h localhost -p 5432 -U postgres -W # 使用密码登录
```

### postgres日志

postgres日志是在`${PGDATA}/pg_log`文件夹中.

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