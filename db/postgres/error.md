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
