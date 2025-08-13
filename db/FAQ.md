# linux postgres 好用的gui client
[dbeaver](https://dbeaver.io/)

## db管理工具
- [heidisql](https://www.heidisql.com/)
- [DBeaver](https://dbeaver.io/)
- [Navicat Premium Lite](https://www.navicat.com.cn/products/navicat-premium-lite)

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

可能情况:
1. 有时所需jar已存在, 但dbeaver还是报找不到driver: 清空`库`标签页, 重新添加jar即可
2. 有时在线下载驱动完成后, dbeaver还是报找不到driver, 重启dbeaver

### 多租户数据隔离
ref:
- [多租户系统设计方案](https://zhuanlan.zhihu.com/p/718325767)

1. 独立数据库 (Separate Database)

	这是最彻底的隔离方式，每个租户都有自己的独立数据库
1. 独立 Schema (Separate Schema)

	将所有租户的数据放在同一个数据库中，但为每个租户创建独立的 Schema（模式）
1. 共享表（Shared Table）

	所有租户的数据都存储在同一张表中，通过一个 tenant_id 字段来区分不同租户的数据

如何选择合适的隔离方案:
特性	独立数据库 (Separate Database)	独立 Schema (Separate Schema)	共享表 (Shared Table)
数据隔离性	最高（物理隔离）	高（逻辑隔离）	最低（应用层隔离）
运维复杂度	最高	中等	最低
硬件成本	最高	中等	最低
可扩展性	低（随租户数线性增长）	中等	最高（适合海量租户）
容灾能力	强（单个租户可独立恢复）	中等	弱（难以隔离恢复）
典型场景	金融、医疗、大型企业应用	中型 SaaS、平台	博客、论坛、轻量级应用