#### `(error) MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.`

在redis-cli运行`ping`命令时碰到(之前直接调用github.com/garyburd/redigo/redis.Pool.Close(),client退出时数据无权限保存).
原因:强制停止redis快照导致不能持久化引起的,具体原因可能是磁盘空间已满或持久化快照时没有权限(这种常见).

运行`config set stop-writes-on-bgsave-error no`命令或将`redis.conf # stop-writes-on-bgsave-error`设为no,可让redis忽略这个异常,不推荐.推荐是用适当权限重启redis-server.

#### 查看当前所有的key

运行`keys *`

#### `ERR Errors trying to SHUTDOWN. Check logs`

停止redis时碰到,可能是权限问题，因为在shutdown命令的时候，会进行save操作，而save需要操作dump.rdb文件，如果没有权限则会报这个错.

解决:先kill redis-server,再用适当权限(通常是sudo)重启.

#### github.com/garyburd/redigo/redis

redis.Pool参数：

- MaxIdle：最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
- MaxActive：最大的激活连接数，表示同时最多有N个连接
- IdleTimeout：最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
- Dial：建立连接的function
- Wait：当Wait=true时,在Pool达到MaxIdle限制(即没有空闲连接)时,Get()等待直到有连接返回Pool中.

 例如：
```
[Redis]
	Addr        = '127.0.0.1:6379'
	Password    = ''
	MaxIdle     = 5
	MaxActive   = 50
	IdleTimeout = 3600 # 1小时
	Wait        = true
```

#### redis-server.service: PID file /var/run/redis/redis-server.pid not readable (yet?) after start: No such file or directory
开始以为是`/var/run/redis/redis-server.pid`需要手动创建, 重启`redis-server.service`后发现还是提示该问题且`/var/run/redis/redis-server.pid`已被删除.

通过查查`redis-server.log`发现错误`Creating Server TCP listening socket ::1:6379: bind: Cannot assign requested address`, 考虑到是使用阿里云 ECS默认ipv6是关闭的, 修改redis.conf的bind参数,取消`::1`即可.

## FAQ
### redis 的数据淘汰策略
redis 内存数据集大小上升到一定大小的时候，就会施行数据淘汰策略。redis 提供 6种数据淘汰策略：
- voltile-lru：从已设置过期时间的数据集（server.db[i].expires）中挑选最近最少使用的数据淘汰
- volatile-ttl：从已设置过期时间的数据集（server.db[i].expires）中挑选将要过期的数据淘汰
- volatile-random：从已设置过期时间的数据集（server.db[i].expires）中任意选择数据淘汰
- allkeys-lru：从数据集（server.db[i].dict）中挑选最近最少使用的数据淘汰
- allkeys-random：从数据集（server.db[i].dict）中任意选择数据淘汰
- no-enviction（驱逐）：禁止驱逐数据

### [Redis主从同步策略](https://blog.csdn.net/sk199048/article/details/50725369)
主从刚刚连接的时候，进行全量同步；全同步结束后，进行增量同步. 当然，如果有需要，slave 在任何时候都可以发起全量同步. redis 策略是，无论如何，首先会尝试进行增量同步，如不成功，要求从机进行全量同步.

### Redis 常见的性能问题
1. Master写rdb，save命令调度rdbSave函数，会阻塞主线程的工作，当快照比较大时对性能影响是非常大的，会间断性暂停服务，所以Master最好不要写rdb
1. Master AOF持久化，如果不重写AOF文件，这个持久化方式对性能的影响是最小的，但是AOF文件会不断增大，AOF文件过大会影响Master重启的恢复速度. Master最好不要做任何持久化工作，包括rdb和AOF日志文件，特别是不要启用rdb做持久化,如果数据比较关键，某个Slave开启AOF备份数据，策略为每秒同步一次.
1. Master调用BGREWRITEAOF重写AOF文件，AOF在重写的时候会占大量的CPU和内存资源，导致服务load过高，出现短暂服务暂停现象
1. Redis主从复制的性能问题，为了主从复制的速度和连接的稳定性，Slave和Master最好在同一个局域网内