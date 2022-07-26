## libpq.so.5: cannot open shared object file: No such file or directory
- 方法1:
```
// 推荐
sudo ln -s /opt/PostgreSQL/9.4/lib/libpq.so.5 /usr/lib
sudo ldconfig
```
- 方法2:
 1. 进入`/{pg_install_path}/lib`确认so是否存在
 2. 启用libpq.so.5
```
cd /etc/ld.so.conf.d
echo "/opt/PostgreSQL/9.4/lib" >>pgsql.conf
sudo ldconfig
```
但该方法可能会和mysql所需的so冲突,而导致mysql命令报错:
```
mysql: /opt/PostgreSQL/9.4/lib/libcrypto.so.1.0.0: no version information available (required by mysql)
mysql: /opt/PostgreSQL/9.4/lib/libssl.so.1.0.0: no version information available (required by mysql)
```

### `pg_basebackup -h localhost -U postgres ...`备份报`Ident authentication failed for user "postgres"`
env: pg 10.5

使用相同的参数, psql可登入. 将`-h localhost`换成`-h 192.168.16.246`后又报`FATAL:  no pg_hba.conf entry for replication connection from host "192.168.16.246", user "postgres", SSL off`, 按照["PostgreSQL 9.4.10执行pg_basebackup的bug及解决方案"](http://www.knockatdatabase.com/2021/10/29/postgresql-9-4-10-pg_basebackup-bug/)的说法是个bug.

解决方法:
1. 创建带复制权限的repuser用户而不是使用postgres用户

    ```sql
    CREATE ROLE repuser WITH REPLICATION LOGIN ENCRYPTED PASSWORD '123456';
    ```
1. 修改pg_hba.conf, 在配置开头加入`host replication repuser all md5`, 再重启pg即可
### `pg_basebackup -Ft`备份报`directory "/data/10/pg_tbs/tbs_mydb" exists but is not empty`
tbs_mydb是创建的表空间

解决方法:
1. 使用`-T`重新指定自定义表空间目录路径即可. 可用`\db`查看表空间的location

    `pg_basebackup -Fp ... -T OLDDIR=NEWDIR`, OLDDIR与NEWDIR不能相同, NEWDIR相当于备份时OLDDIR的备份了.
1. 使用`-Ft`由pg_basebackup自动处理表空间