### redis后台运行

默认情况下，redis不是在后台运行，需要在后台运行时，可把`redis.conf # daemonize`的值更改为yes.

### 日志文件位置

配置`redis.conf # logfile`即可.

### 停止redis

运行`redis-cli shutdown`

### dump.rdb存放位置

    dbfilename dump.rdb
    dir ./

默认dump.rdb存放位置是不固定的，而是存放在启动redis时的当前目录下.修改`dir`路径即可.