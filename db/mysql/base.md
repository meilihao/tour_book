# base
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