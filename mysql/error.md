# Datetime时间误差1s

描述: datetime类型的时间存入db有时会出现误差,误差为`+1s`.

mysql 5.7 有该问题
mariadb 10.1 没有

原因:
mysql 5.7 支持[fractional seconds part(fsp,即小数秒)](http://dev.mysql.com/doc/refman/5.7/en/fractional-seconds.html)进行了四舍五入且没有warning或error的提示.
mariadb 10.1 则截断小数部分,但有warning提示.

# FULLTEXT INDEX (`xxx`) WITH PARSER ngram,查询出错
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

# @@GLOBAL.GTID_PURGED can only be set when @@GLOBAL.GTID_EXECUTED is empty
mysqldump备份db再恢复时报错,解决:
```sql
reset master;
```
