# oracle
ref:
- [Oracle Database版本](https://en.wikipedia.org/wiki/Oracle_Database)
- [Oracle Database Documentation](https://docs.oracle.com/en/database/oracle/oracle-database/index.html)

> Oracle Database Express Edition 是Oracle Database的简装版, 只有一个数据库实例XE, 因此也叫Oracle Database XE. 它拥有正式版的所有功能, 只是在内存和数据大小上做了限制, 适合初学者用来学习Oracle.

## 组件
- Oracle Database : 部署oracle db单机
- Oracle Grid Infrastructure : 部署oracle db cluster

## 部署
- [CentOS 7.4下安装Oracle 11.2.0.4数据库的方法](https://cloud.tencent.com/developer/article/1721536)

## cmd
```bash
# --- 登录db
su - oracle
sqlplus / as sysdba
```

## 数据监控em
- https://localhost:1158/em