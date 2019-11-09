# linux postgres 好用的gui client
[dbeaver](https://dbeaver.io/)

## 迁移db
### sqlite3 -> mysql
1. 使用正则`DEFAULT\s{1,}""` -> `DEFAULT ''`
1. 表名/字段名的`"` -> ```
1. 自增的`autoincrement` -> `auto_increment`
1. 使用正则`\.\d{1,8}\+08:00`替换datetime的时间精度和时区.

### 删除所有表但不删库的方法
`SELECT CONCAT('drop table ',table_name,';') FROM information_schema.`TABLES` WHERE table_schema='数据库名';`