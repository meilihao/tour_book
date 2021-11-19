# 常用命令
- `su postgres -c psql` = `sudo -u postgres psql`
- `sudo -u postgres psql`
- `psql -h 127.0.0.1 -p 5432 -U user -d dbname` : 连接数据库
- `\encoding [编码名称]` : 显示或设定用户端编码
- `\q` : 退出 psql
- `\c dbname` : 切换数据库
- `\l` : 列举数据库
- `\dt` : 列举表
- `\d tblname` : 查看表结构
- `\di` : 查看索引

## 添加用户
```sql
postgres=# CREATE DATABASE mytestdb;
CREATE DATABASE
postgres=# CREATE USER mytestuser WITH ENCRYPTED PASSWORD '123456'; # 或`create user root with password 'password';`
CREATE ROLE
postgres=# GRANT ALL PRIVILEGES ON DATABASE mytestdb to mytestuser;
GRANT
```

## sql结果保存到file
```sql
postgres=# \o /var/lib/pgsql/data/output.txt -- 前提: db user有对应目录保存文件的权限
postgres=# select * from dummy_table;
postgres=# \o
```