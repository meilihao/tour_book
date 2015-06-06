### systemctl status postgresql-9.4.service提示`Failed to start SYSV: Starts and stop`

估计是未用systemctl启动导致.

先停止postgres(`pg_ctl stop -D /opt/PostgreSQL/9.4/data`),再用`sudo systemctl start postgresql-9.4.service`重新启动即可.

### 大小写

创建database/table时,postgres会自动将数据库名,表名和字段名转为小写;执行sql时也会将表名和字段名转为小写.

## 常用命令

- psql登入 : `psql -U user_name`
- 列出所有的数据库 : `\l`
- 切换数据库 : `\c db_name`
- 列出当前数据库下的数据表 : `\d`
- 列出指定表的所有字段 : `\d tb_name`
- 查看指定表的基本情况 : `\d+ tb_name`
- psql登出 ： `\q`