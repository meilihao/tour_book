### FAQ

#### 显示sql历史
参考：
- [mysql日志介绍](https://andblog.cn/?p=928)

```shell
# 方法1:
# 修改/etc/mysql/my.cnf
# 设置[mysqld]节的general_log和general_log_file,再重启mysql即可
# ///
# 方法2:
# 用此方法mysql重启后会失效
# 进入mysql console
show variables like '%general_log%';
set global general_log=on;
set global log_output={FILE|TABLE|NONE}	 #设置日志的输出方式，可以是输出到文件中：FILE，或者输出到数据库的表中，或者设置为：FILE,TABLE，两个都记录
set global general_log_file='/var/log/mysql/mysql.log'; # 注意global general_log_file在Linux中只能设置到 /tmp 或 /var 文件夹下，设置其他路径会报错. general_log_file是相对路径时在data dir目录下(`/var/lib/mysql`)
```

#### 查询mysql使用到的目录
```sql
SHOW VARIABLES WHERE Variable_Name LIKE "%dir";
```

#### 修改用户密码

```shell
mysqladmin -u 用户名 -p password "新密码"
```

mariadb 10.4 初始化密码:
```
$ sudo mysql -u root
> set password for 'root'@'localhost' = PASSWORD('xxx');
```

#### 查询sql_mdoe
```sql
select @@sql_mode,@@GLOBAL.SQL_MODE;
```
### 查看warning信息

    show warnings;

#### 查看表结构

	desc table_name;

#### 查看建表语句

	show create table table_name;

#### 查看schema(数据库)下的所有table

	show tables;

#### 查看实例下的所有数据库

	show databases;

#### 切换scheme

	use scheme_name;

#### 查看mysql变量

	show variables;

##### mysql备份

其他参考: http://www.cnblogs.com/Cherie/p/3309456.html

```
1) 导出整个数据库
    mysqldump -u 用户名 -p 数据库名 > 导出的文件名
    mysqldump -u user_name -p123456 database_name > outfile_name.sql

2) 导出一个表
    mysqldump -u 用户名 -p 数据库名 表名> 导出的文件名
    mysqldump -u user_name -p database_name table_name > outfile_name.sql

3) 导出一个数据库结构
    mysqldump -u user_name -p -d --add-drop-table database_name > outfile_name.sql
    -d 没有数据 –add-drop-table 在每个create语句之前增加一个drop table

4) 带语言参数导出
    mysqldump -uroot -p --default-character-set=latin1 --set-charset=gbk --skip-opt database_name > outfile_name.sql
1. 备份多个db

    mysqldump -uroot -p123456 --single-transaction --databases db1 db2> dbs.sql

ps:

备份时排除系统表：

    mysql -e "show databases;" -uroot -proot | grep -Ev "Database|information_schema|mysql|performance_schema" | xargs mysqldump -uroot -proot --databases > /home/chen/mysql_dump.sql

数据库导入:

mysql -h 10.6.208.183 -u test2 -p  -P 3310 目的数据库名称 < test.sql;也可以直接在mysql命令行下面用source导入(先用use进入到某个数据库，mysql>source /home/xxx/test.sql，后面的参数为sql文件).注意,**导入前应先确保目的数据库存在**.
```

### 基于模板表创建新表
`CREATE TABLE if not EXISTS bill_<N> like bill"`

### 大小写

数据库中表名用小写,程序中表名用大写开头,mysql报错`Table '数据库名.表名' doesn't exist`.

原因:

MySQL在Linux下数据库名、表名、列名、别名大小写规则是这样的：

- 数据库名与表名是严格区分大小写的；
- 表的别名是严格区分大小写的；
- 列名与列的别名在所有的情况下均是忽略大小写的；
- 变量名也是严格区分大小写的；

注:MySQL在Windows下都不区分大小写.

在 MySQL 中，数据库和表对就于那些目录下的目录和文件。因而，操作系统的敏感性决定数据库和表命名的大小写敏感。这就意味着数据库和表名在 Windows 中是大小写不敏感的，而在大多数类型的 Unix 系统中是大小写敏感的。

在MySQL的配置文件中my.ini [mysqld] 中增加`lower_case_table_names = 1`(0：区分大小写;1：不区分大小写)即可,这样MySQL 将在创建与查找时将所有的表名自动转换为小写字符,不过不推荐这种方法.

推荐的命名规则是：在定义数据库、表、列的时候全部采用小写字母加下划线的方式，不使用任何大写字母.

#### 字符集

查询数据库支持的编码

    show character set;

查看mysql当前使用的编码(注意要在未use db时查看，否则看到的是当前db使用的编码)

    status;

查看数据库编码：

    SHOW CREATE DATABASE db_name;

查看表编码：

    SHOW CREATE TABLE tbl_name;

查看字段编码：

    SHOW FULL COLUMNS FROM tbl_name;

修改数据库字符集：

    ALTER DATABASE db_name DEFAULT CHARACTER SET character_name [COLLATE ...];

把表默认的字符集和所有字符列（CHAR,VARCHAR,TEXT）改为新的字符集：

    ALTER TABLE tbl_name CONVERT TO CHARACTER SET character_name [COLLATE ...]
    //如：ALTER TABLE logtest CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;

只是修改表的默认字符集：

    ALTER TABLE tbl_name DEFAULT CHARACTER SET character_name [COLLATE...];
    //如：ALTER TABLE logtest DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

修改字段的字符集：

    ALTER TABLE tbl_name CHANGE c_name c_name c_type CHARACTER SET character_name [COLLATE ...];
    //如：ALTER TABLE logtest CHANGE title title VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_general_ci;

注:修改字符集时无法将原先存入latin1字符集字段中的中文转为utf-8,即latin1不兼容utf8.

[配置默认编码为utf8](https://mariadb.com/kb/en/mariadb/setting-character-sets-and-collations/)

参考:

- [理解和解决 MySQL 乱码问题](https://linux.cn/article-5028-1.html)
- [十分钟搞清字符集和字符编码](https://linux.cn/article-5027-1.html)

### 自增

插入时,如果在自增字段上指定了数值,那么将由指定数值取代默认的自增值.

### 启动脚本

[mariadb](https://mariadb.com/kb/zh-cn/iniciando-e-parando-mariadb-automaticamente/)脚本在`/etc/init.d/mysql`.

### 主键

    Alter table tb_name add primary key(id);
    Alter table tb_name drop primary key;

### warning

#### update limit warning

[官方文档](https://dev.mysql.com/doc/refman/5.6/en/replication-features-limit.html),可忽略该警告.但推荐update时不使用limit.

### update

update重复执行相同语句(即同一语句多次执行),返回mysql_affected_rows为0,[文档在这](https://mariadb.com/kb/en/mariadb/mysql_affected_rows/).

### mysql workbench

### 技巧

- 在`SQL Editor`窗口只执行选中的sql语句,使用快捷键F9.

### 关键字

    alter table Avatar drop column `Using`//报错

因为`Using`是mysql的关键字,这里是建表时错误使用造成的,需用反引号包裹关键字,再进行操作.

    alter table Avatar drop column `Using`

## CRUD出错

### 插入空值出错

mysql sql_mode包含"STRICT_TRANS_TABLES"时(严格模式),db不允许插入空值.在建表时该字段添加default即可解决.

### 插入时间"0000-00-00"报错

mysq的sql_mode使用了"NO_ZERO_DATE".文档:[mariadb](https://mariadb.com/kb/en/mariadb/datetime/),[mysql5.7](http://dev.mysql.com/doc/refman/5.7/en/sql-mode.html#sqlmode_no_zero_date)

### int(M) 含义

int(M) 表示使用integer数据类型，而M表示**最大显示宽度**,与存储空间无关.

### 从 innodb 的索引结构分析，为什么索引的 key 长度不能太长
key 太长会导致一个页当中能够存放的 key 的数目变少，间接导致索引树的页数目变多，索引层次增加，从而影响整体查询变更的效率

### MySQL 的数据如何恢复到任意时间点？
恢复到任意时间点以定时的做全量备份，以及备份增量的 binlog 日志为前提. 恢复到任意时间点首先将全量备份恢复之后，再此基础上回放增加的 binlog 直至指定的时间点.

### mariadb 10.4 系统root + root@localhost登录时无需密码
[mariadb 10.4就是这么设计的.](https://mariadb.org/authentication-in-mariadb-10-4/)

### Found 1 prepared transactions! It means that mysqld was not shut down properly last time and critical recovery information (last binlog or tc.log file) was manually deleted after a crash. You have to start mysqld with --tc-heuristic-recover switch to commit or rollback pending transactions.
```sh
# step1, 交易回滚
# mysqld --tc-heuristic-recover=ROLLBACK
# step2, 重新启动mysql
# systemctl start mysql
```

### Specified key was too long; max key length is 767 bytes
数据库表采用utf8编码，其中varchar(255)的column进行了唯一键索引,而mysql默认情况下单个列的索引不能超过767位(不同版本可能存在差异),于是utf8字符编码下，255*3 byte 超过限制.

### json操作
```sql
-- conds = conds.And(Expr(`JSON_CONTAINS(mcode,'["` + code + `"]','$')`))
select * FROM `raw_license` WHERE JSON_CONTAINS(mcode,'["04932265479995"]') ;  -- mcode是array, 检索成员的格式必须是`'[xxx ,...]'`
```

### grant
GRANT 命令的常见格式以及解释
命令 作用
GRANT 权限 ON 数据库.表单名称 TO 账户名@主机名 | 对某个特定数据库中的特定表单给予授权
GRANT 权限 ON 数据库.*TO 账户名@主机名 | 对某个特定数据库中的所有表单给予授权
GRANT 权限 ON*.*TO 账户名@主机名 | 对所有数据库及所有表单给予授权
GRANT 权限 1,权限 2 ON 数据库.*TO 账户名@主机名 | 对某个数据库中的所有表单给予多个授权
GRANT ALL PRIVILEGES ON *.*TO 账户名@主机名 | 对所有数据库及所有表单给予全部授权(需谨慎操作)

查看权限: `SHOW GRANTS FOR luke@localhost;`

### 全同步复制
全同步复制（full sync replication）是指当主库执行完一个事务后，需要确保所有的从库都执行了该事务才返回给客户端. 因为需要等待所有的从库都执行完该事务才能返回，所以全同步复制的性能较差.

MySQL自身不支持同步复制，需要用到第三方工具如DRBD（sync模式）等实现同步复制，严格来说，把半同步复制技术默认（或人为）全部应用到所有从库上也算是全同步复制.

[MySQL](https://dev.mysql.com/doc/mysql-replication-excerpt/8.0/en/replication.html)/[mariadb](https://mariadb.com/kb/en/semisynchronous-replication/)仅支持两种复制:
- Asynchronous replication
- Semi-synchronous replication

### mariadb 10.4 忘记密码
```bash
# vim /lib/systemd/system/mariadb.service
ExecStart=/usr/sbin/mysqld --skip-grant-tables ...
# systemctl daemon-reload
# systemctl restart mariadb.service
# mysql -h localhost -u root
> flush privileges; -- 先刷新一下权限表, 否则会报: xxx is running with the --skip-grant-tables option so it cannot execute this statement
> ALTER USER 'root'@'localhost' IDENTIFIED BY 'passowrd';
> exit
```

记得删除skip-grant-tables, 并重启mariadb.

### sql
```sql
> CREATE USER luke@localhost IDENTIFIED BY 'linuxprobe';
```

### 允许MySQL的root用户远程登录
```sql
> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'password';
> FLUSH PRIVILEGES;
```

### master_use_gtid介绍
master_use_gtid是当Slave连接到Master时，Master将从哪个GTID开始给Slave复制 Binlog.

`master_use_gtid = { slave_pos | current_pos | no }`有3种选项：

- no: 传统复制(即非gtid复制, **不推荐**)
- slave_pos：slave将Master最后一个GTID的position复制到本地，Slave主机可通过gtid_slave_pos变量查看最后一个GTID的position
- current_pos：使用当前机器的gtid_current_pos(`select @@gtid_current_pos;`)作为master_use_gtid. 假设有AB两台主机，A是Master，当A故障后，B成为Master，A修复后以Slave的身份重新添加，A之前从没担任过slave角色，所以没有之前复制的GTID号，此时gtid_slave_pos为空，为了能让A能自动添加为Slave，此时就用到该选项(因为A原先有数据了). 该选项是大多数情况下使用的选项，因为他简单易用同，不必在意服务器之前是Master还是Slave角色。但要注意不要让从服务器在binlog日志中写入事务

建议在服务器上启用gtid_strict_mode(不让从服务器在binlog日志中写入事务)，这样非Master产生的事务将被拒绝. 如果从服务器没有开启binlog上面两种方式等价.

### ERROR 2002 (HY000):Can't connect to MySQL server on 'xxx' (115)
mariadb server仅监听了127.0.0.1.

```bash
# vim /etc/mysql/my.cnf
[mysqld]
bind-address=0.0.0.0
```

### show slave status报: Error: connecting slave requested to start from GTID X-X-X, which is not in the master's binlog
slave gtid_slave_pos比master早了.

### MySQL RESET MASTER与RESET SLAVE
删除所有index file 中记录的所有binlog 文件，将日志索引文件清空，创建一个新的日志文件，这个命令通常仅仅用于第一次用于搭建主从关系的时的主库

注意reset master 不同于purge binary log的两处地方:
1. reset master 将删除日志索引文件中记录的所有binlog文件，创建一个新的日志文件 起始值从000001 开始，然而purge binary log 命令并不会修改记录binlog的顺序的数值
2. reset master 不能用于有任何slave 正在运行的主从关系的主库。因为在slave 运行时刻 reset master 命令不被支持，reset master 将master 的binlog从000001 开始记录,slave 记录的master log 则是reset master 时主库的最新的binlog,从库会报错无法找的指定的binlog文件

reset slave 将使slave 清除主从复制关系的位置信息. 该语句会删除master.info文件和relay-log.info 文件以及所有的relay log 文件(即使relay log中还有SQL没有被SQL线程apply完), 并重新启用一个新的relay log文件.

使用reset slave之前必须使用stop slave 命令将复制进程停止.

> RESET SLAVE有个问题，它虽然删除了上述文件，但内存中的change master信息并没有删除. RESET SLAVE ALL还会删除内存中的连接信息，此时执行start slave会报错.

### mariadb无法启动报`InnoDB: Missing FILE_CREATE, FILE_DELETE or FILE_MODIFY before FILE_CHECKPOINT for tablespace 4044`
env: ubuntu 22.04 + mariadb 10.6.7

[解决方法](https://github.com/MariaDB/mariadb-docker/issues/267#issuecomment-801152453)

[`[ERROR] InnoDB: Missing FILE_CREATE, FILE_DELETE or FILE_MODIFY before FILE_CHECKPOINT during Mariabackup prepare`](https://jira.mariadb.org/browse/MDEV-27711)

还是不行就卸载重新安装.

### Datetime时间误差1s

描述: datetime类型的时间存入db有时会出现误差,误差为`+1s`.

mysql 5.7 有该问题
mariadb 10.1 没有

原因:
mysql 5.7 支持[fractional seconds part(fsp,即小数秒)](http://dev.mysql.com/doc/refman/5.7/en/fractional-seconds.html)进行了四舍五入且没有warning或error的提示.
mariadb 10.1 则截断小数部分,但有warning提示.

### FULLTEXT INDEX (`xxx`) WITH PARSER ngram,查询出错
表中一列为name,有两行name的值均为`查查`,可是每次查询只能查到一行.原因:name中文字不同,但是字形很相近.
```sql
mysql>  select id,name,hex(name) from dispatcher where id in (10057,10026);
+-------+--------+--------------+
| id    | name   | hex(name)    |
+-------+--------+--------------+
| 10026 | 查查   | E69FA5E69FA5 |
| 10057 | 査査   | E69FBBE69FBB |
+-------+--------+--------------+
2 rows in set (0.01 sec)
```

> 査 zhā 姓氏;査，拼音chá.
> 查看十六进制用`hex(xxx),推荐`,二进制用`bin(xxx),仅用于整数`

### @@GLOBAL.GTID_PURGED can only be set when @@GLOBAL.GTID_EXECUTED is empty
mysqldump备份db再恢复时报错,解决:
```sql
reset master;
```

### 一个表中可以有多个自增列?
一个表中只能有一个自增列

### mariadb(>=10.1.48)备份还原
ref:
 - [Full Backup and Restore with Mariabackup](https://mariadb.com/kb/en/full-backup-and-restore-with-mariabackup/)
 - [Incremental Backup and Restore with Mariabackup](https://mariadb.com/kb/en/incremental-backup-and-restore-with-mariabackup/)

    mariadb 10.1和>=10.2的增量还原有区别, 见`Incremental Backup and Restore with Mariabackup`

关键[选项](https://www.hanzz.red/archives/mysql%E5%A4%87%E4%BB%BD%E4%B8%8E%E6%81%A2%E5%A4%8D):
- --copy-back : 做数据恢复时将备份数据文件拷贝到MySQL服务器的datadir. 使用该选项, 则下次还可用其还原数据且不用再`--prepare`, 因此再次prepare全备时不报错, 但prepare增量时会报`This target seems to be already prepared.`
- --move-back : 这个选项与–copy-back相似, 唯一的区别是它不拷贝文件, 而是移动文件到目的地. 这个选项会移除backup文件，用时候必须小心.

    `mariabackup --move-back`后再次使用时`mariabackup --prepare`会报错(我这里是直接core dump)

其他:
1. 备份文件中的`xtrabackup_checkpoints`的`backup_type`可表明当次备份是全备还是增量
1. **还原时mariabackup需要root权限**. `--prepare`是检查用于还原的备份的数据文件一致性

```bash
# --- 全量备份
mariabackup --backup --target-dir=/var/mariadb/backup/ --user=root --password=123456

systemctl stop mariadb
rm -rf /var/lib/mysql/* # 确保还原前为空目录
mariabackup --prepare --target-dir=/var/mariadb/backup
mariabackup --copy-back --target-dir=/var/mariadb/backup/
chown -R mysql:mysql /var/lib/mysql/
systemctl start mariadb

# --- 增量备份
mariabackup --backup --target-dir=/var/mariadb/backup/ --user=root --password=123456
mariabackup --backup --target-dir=/var/mariadb/inc1/ --incremental-basedir=/var/mariadb/backup/ --user=root --password=123456
# 基于上次增量备份做增量备份
mariabackup --backup --target-dir=/var/mariadb/inc2/ --incremental-basedir=/var/mariadb/inc1/ --user=root --password=123456

systemctl stop mariadb
rm -rf /var/lib/mysql/* # 确保还原前为空目录
# 准备全量备份文件
mariabackup --prepare --target-dir=/var/mariadb/backup
# 准备增量备份文件
mariabackup --prepare --target-dir=/var/mariadb/backup --incremental-dir=/var/mariadb/inc1 # 检查增备前必须检查全备, 否则会报错`applying incremental backup need a prepared target`
# 恢复数据
mariabackup --copy-back --target-dir=/var/mariadb/backup/ --incremental-dir=/var/mariadb/inc1
# 修改数据文件权限
chown -R mysql:mysql /var/lib/mysql/
systemctl start mariadb
```

### varchar 和 char 的区别
- char 表示变长，char 表示长度固定，未满填充空格，超出国定长度则拒绝插入并提示错误信息
- 存储容量不同: 对char来说，最多能存放的字符个数为255. 对于varchar,最多能存放的字符个数是65532
- 存储速度不同: char长度固定，存储速度会比varchar快一些，但在空间上会占用额外的空间，属于一种空间换时间的策略. varchar空间利用率会更高些

## in和exists一般用于子查询
使用exists时会先进行外表查询，将查询到的每一行数据都带入内表查询中看是否满足条件；使用in一般会先进行内表查询获取结果集，然后对外表查询匹配结果集，返回数据.

in在内表查询或者外表查询过程中都会用到索引
exsits仅在内表查询时会用到索引

一般来说，当子查询的结果集比较大，外表较小时用exist效率更高；当子查询的结果较小，外表较大时，使用in效率更高.

对于not in 和 not exists, not exists效率比not in 效率高，与子查询的结果集无关，因为not in 对于内外表都进行了全表扫描，没有使用到索引. not exists的子查询中可以用到表上的索引.

### mariadb client 11.8连接aliyun rdb报`ERROR 2026 (HY000): TLS/SSL error: SSL is required, but the server does not support it`
追加`--skip-ssl`