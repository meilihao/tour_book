## FAQ

1 . 显示sql历史

```shell
# 方法1:
# 修改/etc/mysql/my.cnf
# 设置[mysqld]节的general_log和general_log_file,再重启mysql即可
# ///
# 方法2:
# 用此方法mysql重启后会失效
# 进入mysql console
show variables like '%general_log%';
set global general_log=on;
# 注意global general_log_file在Linux中只能设置到 /tmp 或 /var 文件夹下，设置其他路径会报错
set global general_log_file='/var/log/mysql/mysql.log';
```