# base
## 存储过程
ref:
- [MySQL中的存储过程（详细篇）](https://zhuanlan.zhihu.com/p/679169773)

存储过程：（PROCEDURE）是事先经过编译并存储在数据库中的一段SQL语句的集合. **强烈不推荐使用**

基本语法:
- 存储过程中的参数分别是 in，out，inout三种类型:

	1. in代表输入参数（默认情况下为in参数），表示该参数的值必须由调用程序指定
	1. ou代表输出参数，表示该参数的值经存储过程计算后，将out参数的计算结果返回给调用程序
	1. inout代表即时输入参数，又是输出参数，表示该参数的值即可有调用程序制定，又可以将inout参数的计算结果返回给调用程序
- 存储过程中的语句必须包含在BEGIN和END之间。
- DECLARE中用来声明变量，变量默认赋值使用的DEFAULT，语句块中改变变量值，使用SET 变量=值

```sql
SHOW   {  PROCEDURE  |  FUNCTION  }  status -- 查看列表
SHOW   CREATE   {  PROCEDURE  |  FUNCTION  }  sp_name -- 看不到检查权限
```

## 表空间
> [mysql通过配置项：innodb_file_per_table指定MySQL使用独立表空间, MySQL5.6.6及以后的版本默认值是ON, MySQL5.6.5以前的版本默认值是OFF](http://uusama.com/922.html).
> [The CREATE TABLESPACE statement is not supported by MariaDB](https://mariadb.com/kb/en/create-tablespace/)

	创建tablespace语句(必须包含`ENGINE=xxx`)会成功, 但不会有实际效果.

```sql
mysql> create tablespace big_data_in_mysql add datafile 'first.ibd' [ENGINE [=] engine_name]; -- 未指定存储目录, 所以使用默认存储路径. 这个表空间所对应的元数据存放在first.ibd这个文件中
mysql> create tablespace test_tablespace add datafile '/var/lib/test_mysql_tablespace/first.ibd';
mysql> select * from information_schema.INNODB_SYS_TABLESPACES; -- 查看表空间
mysql> create table t2(c1 int primary key) tablespace test_tablespace;
mysql> DROP TABLESPACE tablespace_name [ENGINE [=] engine_name];
```