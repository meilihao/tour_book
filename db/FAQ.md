# linux postgres 好用的gui client
[dbeaver](https://dbeaver.io/)

## db管理工具
- [heidisql](https://www.heidisql.com/)
- [DBeaver](https://dbeaver.io/)

## 迁移db
### sqlite3 -> mysql
1. 使用正则`DEFAULT\s{1,}""` -> `DEFAULT ''`
1. 表名/字段名的`"` -> ```
1. 自增的`autoincrement` -> `auto_increment`
1. 使用正则`\.\d{1,8}\+08:00`替换datetime的时间精度和时区.

### 删除所有表但不删库的方法
`SELECT CONCAT('drop table ',table_name,';') FROM information_schema.`TABLES` WHERE table_schema='数据库名';`

### dbeaver缺失驱动
ref:
- [如何在无网络的情况下给Dbeaver安装数据库驱动](https://blog.csdn.net/Georgetwo/article/details/112390120)

	数据库->驱动管理器->选中一种db, 点右侧"编辑"-> 选中"库"标签页, 删除原有内容(可以看到需要的driver及其版本), 点右侧"添加文件", 加入驱动jar文件即可

- [miaridb : mariadb-java-client](https://mvnrepository.com/artifact/org.mariadb.jdbc/mariadb-java-client)
- [postgres : postgresql-42.2.20.jar](https://mvnrepository.com/artifact/org.postgresql/postgresql)
- [oracle]()

有时所需jar命令存在, 但dbeaver还是报找不到driver: 清空`库`标签页, 重新添加jar即可
