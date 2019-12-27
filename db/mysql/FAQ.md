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

ps:

备份时排除系统表：

    mysql -e "show databases;" -uroot -proot | grep -Ev "Database|information_schema|mysql|performance_schema" | xargs mysqldump -uroot -proot --databases > /home/chen/mysql_dump.sql

数据库导入:

mysql -h 10.6.208.183 -u test2 -p  -P 3310 目的数据库名称 < test.sql;也可以直接在mysql命令行下面用source导入(先用use进入到某个数据库，mysql>source /home/xxx/test.sql，后面的参数为sql文件).注意,**导入前应先确保目的数据库存在**.
```

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

### [mariadb 离线下载地址](https://mariadb.com/downloads/)

### Found 1 prepared transactions! It means that mysqld was not shut down properly last time and critical recovery information (last binlog or tc.log file) was manually deleted after a crash. You have to start mysqld with --tc-heuristic-recover switch to commit or rollback pending transactions.
```sh
# step1, 交易回滚
# mysqld --tc-heuristic-recover=ROLLBACK
# step2, 重新启动mysql
# systemctl start mysql
```

### Specified key was too long; max key length is 767 bytes
数据库表采用utf8编码，其中varchar(255)的column进行了唯一键索引,而mysql默认情况下单个列的索引不能超过767位(不同版本可能存在差异),于是utf8字符编码下，255*3 byte 超过限制.