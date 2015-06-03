### FAQ

#### 显示sql历史

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
# 注意global general_log_file在Linux中只能设置到 /tmp 或 /var 文件夹下，设置其他路径会报错
set global general_log_file='/var/log/mysql/mysql.log';
```

#### 修改用户密码

```shell
mysqladmin -u 用户名 -p password "新密码"
```

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
    mysqldump -u user_name -p -d –add-drop-table database_name > outfile_name.sql
    -d 没有数据 –add-drop-table 在每个create语句之前增加一个drop table

4) 带语言参数导出
    mysqldump -uroot -p –default-character-set=latin1 –set-charset=gbk –skip-opt database_name > outfile_name.sql

数据库导入:

mysql -h 10.6.208.183 -u test2 -p  -P 3310 目的数据库名称 < test.sql;也可以直接在mysql命令行下面用source导入(先用use进入到某个数据库，mysql>source /home/xxx/test.sql，后面的参数为sql文件).注意,**导入前应先确保目的数据库存在**.
```

#### 大小写

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