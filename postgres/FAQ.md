### systemctl status postgresql-9.4.service提示`Failed to start SYSV: Starts and stop`

估计是未用systemctl启动导致.

先停止postgres(`pg_ctl stop -D /opt/PostgreSQL/9.4/data`),再用`sudo systemctl start postgresql-9.4.service`重新启动即可.

### 大小写

创建database/table时,postgres会自动将数据库名,表名和字段名转为小写;执行sql时也会将表名和字段名转为小写.
