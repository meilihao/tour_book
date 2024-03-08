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

### 分布式锁
```
tryLock(){  
    SET Key UniqId Seconds // 官方文档上提醒后面的版本有可能去掉SETNX, SETEX, PSETEX,并用SET命令代替
}
release(){  
    EVAL(
      //LuaScript
      if redis.call("get",KEYS[1]) == ARGV[1] then
          return redis.call("del",KEYS[1])
      else
          return 0
      end
    )
}
```

这个方案是目前最优的分布式锁方案, 在单实例redis的场景下是安全的，但是如果在Redis集群环境下依然存在问题: 由于Redis集群**数据同步为异步**，假设在Master节点获取到锁后未完成数据同步情况下Master节点crash，此时在新的Master节点依然可以获取锁，所以多个Client同时获取到了锁, 除非业务场景可以接受, 那么这种小概率可以忽略. 

因此不建议使用Redis集群,而是使用业务分片的单机版redis, 因为目前而言, redis已经足够稳定.

### MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error
Redis被配置为保存数据库快照，但它目前不能持久化到硬盘, 通常是权限问题.

> 我这里是存在所属用户是root的/var/lib/redis/dump.rdb, 导致以普通用户redis运行的redis无法替换同名文件.

解决方案：
- 关闭Redis快照持久化

    `127.0.0.1:6379> config set stop-writes-on-bgsave-error no`
- 处理权限

    > CONFIG SET dir /tmp/some/directory/other/than/var # 更换redis dir
    > CONFIG SET dbfilename temp.rdb

### 复制进度
参考:
- [Redis主从同步与故障切换，有哪些坑？](https://new.qq.com/omn/20201125/20201125A0GFNT00.html)

通过redis 的`INFO replication`命令查看主库接收写命令的进度信息（master_repl_offset）和从库复制写命令的进度信息（slave_repl_offset）, 即`diff=master_repl_offset-slave_repl_offset`, diff=0为复制完成.

### redis cmd监控
`redis-cli monitor`

### 模拟redis aof文件损坏
ref:
- [Redis之AOF重写及其实现原理](https://blog.csdn.net/hezhiqiang1314/article/details/69396887)

```bash
cp /var/lib/redis/appendonly.aof . # 获取正常aof文件
redis-check-aof appendonly.aof # 获得size=59962
truncate appendonly.aof -s 59960
echo "y"|redis-check-aof --fix appendonly.aof
```

解决方法:
```
cat redis.service
[Service]
ExecStartPre=/usr/bin/bash -c "echo 'y'|redis-check-aof --fix /var/lib/redis/appendonly.aof||true" # 可能会丢少量数据
```

> 追加`||true`原因: appendonly.aof不存在或大小为0时, redis-check-aof会报错

> aof超过6G时遇到启动redis时, 因为redis-check-aof执行长, 导致systemd启动redis超时

    解决方法:
    1. edis-cli 手动重写

        ```bash
        > bgrewriteaof
        ```
    1. redis.conf

        ```conf
        auto-aof-rewrite-percentage 参数表示当当前 AOF 文件大小超过上次重写后 AOF 文件大小的百分比时，触发 AOF 重写机制，默认值为 100
        auto-aof-rewrite-min-size 参数表示当当前 AOF 文件大小超过指定值时，才可能触发 AOF 重写机制，默认值为 64 MB
        ```

        环境已应用了该配置, 但aof还是涨到了6G, 按照业务数据量应该是不可能, 不知为什么没生效, 又因为redis都是缓存数据, 直接删除aof即可 by `find /var/log/redis -mindepth 1 -name appendonly.aof -size +1G/+512M -delete`.

### 应用连接redis报`LOADING Redis is loading the dataset in memory`
如果在系统将数据集完全加载到内存并使 Redis 准备好之前, 连接请求到达就会触发该报错

解决方法: 应用等待, 或使用 info 命令检查redis是否正在加载.

### Redis 4.0 混合持久化
Redis 4.0前, aof是全量的日志; 4.0开始支持混合持久化, 此时aof是自持久化开始到持久化结束的这段时间发生的增量 AOF 日志，通常这部分 AOF 日志很小. 开启后aof rewrite的时候就直接把 rdb 的内容写到 aof 文件开头.

开启混合持久化: `aof-use-rdb-preamble yes`