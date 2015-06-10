### systemctl status postgresql-9.4.service提示`Failed to start SYSV: Starts and stop`

估计是未用systemctl启动导致.

先停止postgres(`pg_ctl stop -D /opt/PostgreSQL/9.4/data`),再用`sudo systemctl start postgresql-9.4.service`重新启动即可.

### 大小写

创建database/table时,postgres会自动将数据库名,表名和字段名转为小写;执行sql时也会将表名和字段名转为小写.

## 常用命令

psql:

- psql登入 : `psql -U user_name`
- 列出所有的数据库 : `\l`
- 切换数据库 : `\c db_name`
- 列出当前数据库下的数据表 : `\d`
- 列出指定表的所有字段 : `\d tb_name`
- 查看指定表的基本情况 : `\d+ tb_name`
- psql登出 ： `\q`

postgres:

- 启动数据库 : `pg_ctl start -D $PGDATA`,$PGDATA是postgres的数据目录
- 停止数据库 : `pg_ctl stop -D $PGDATA [-m SHUTDOWN-MODE]`,SHUTDOWN-MODE:smart(默认),等所有连接中止才关闭,如果client连接不中止,则无法关闭数据库;fast,主动断开client连接,回滚已有事务,然后正常关闭;immediate,立即停止数据库进程,直接退出,下次启动数据库需进行crash-recovery操作.
