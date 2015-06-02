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

mysql -h 10.6.208.183 -u test2 -p  -P 3310 < test.sql;也可以直接在mysql命令行下面用source导入(先用use进入到某个数据库，mysql>source /home/xxx/test.sql，后面的参数为sql文件)
```