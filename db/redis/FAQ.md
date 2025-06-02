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
- 设置了过期时间的key即从已设置过期时间的数据集（server.db[i].expires）中挑选

    - voltile-lru(v3.0前的默认策略)：淘汰最久未使用的键值
    - voltile-lfu(from v4.0)：淘汰最少使用的键值
    - volatile-ttl：优先淘汰更早过期的键值
    - volatile-random：随机淘汰任意键值
- 所有key即从数据集（server.db[i].dict）中挑选

    - allkeys-lru：淘汰最久未使用的键值
    - allkeys-lfu(from v4.0)：淘汰最少使用的键值
    - allkeys-random：随机淘汰任意键值
- 不淘汰

    - no-enviction（驱逐, v3.0及以后默认策略）：禁止驱逐数据

建议:
1. keys访问频率相近: allkeys-random
1. keys访问类似正态分布: allkeys-lru

LRU和LFU都是内存管理的页面置换算法:
- LRU：最近最少使用(最长时间)淘汰算法（Least Recently Used）. LRU是淘汰最长时间没有被使用的页面
- LFU：最不经常使用(最少次)淘汰算法（Least Frequently Used）. LFU是淘汰一段时间内，使用次数最少的页面

### 大量 key 集中过期问题
1. 给 key 设置随机过期时间
2. 开启 lazy-free（惰性删除/延迟释放）

    lazy-free 特性是 Redis 4.0 开始引入的，指的是让 Redis 采用异步方式延迟释放 key 使用的内存，将该操作交给单独的子线程处理，避免阻塞主线程

建议不管是否开启 lazy-free，我们都尽量给 key 设置随机过期时间

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

Redlock 算法的思想是让客户端向 Redis 集群中的多个独立的 Redis 实例依次请求申请加锁，如果客户端能够和半数以上的实例成功地完成加锁操作，那么我们就认为，客户端成功地获得分布式锁，否则加锁失败.

实际项目中不建议使用 Redlock 算法，成本和收益不成正比，可以考虑基于 Redis 主从复制+哨兵模式实现分布式锁.

### MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error
Redis被配置为保存数据库快照，但它目前不能持久化到硬盘, 通常是权限问题.

> 我这里是存在所属用户是root的/var/lib/redis/dump.rdb, 导致以普通用户redis运行的redis无法替换同名文件.

解决方案：
- 关闭Redis快照持久化

    `127.0.0.1:6379> config set stop-writes-on-bgsave-error no`
- 处理权限

    > CONFIG SET dir /tmp/some/directory/other/than/var # 更换redis dir
    > CONFIG SET dbfilename temp.rdb

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

### 主从同步停止
在slave端执行`info replication`, 看`master_link_status`状态, `down`为已停止.

看slave日志发现slave连接master成功, 但之后会报`MASTER aborted replication with an error: NOAUTH Authentication required`, 在slave配置`masterauth 主库的密码`即可.

### SETNX和SET NX区别
ref:
- [Redis分布式锁的进化之路](https://mbd.baidu.com/newspage/data/dtlandingsuper?nid=dt_4488365490908995111)

SETNX‌：
- ‌优点‌：实现简单，适用于需要确保键不存在时才设置值的场景，如分布式锁
‌- 缺点‌：没有设置过期时间，可能导致死锁问题。需要结合EXPIRE命令来设置过期时间，但这会破坏原子性
‌SET NX‌：
‌- 优点‌：作为SET命令的一部分，具有更好的灵活性和原子性，可以同时设置值和过期时间
‌- 缺点‌：需要结合其他命令（如EXPIRE）来实现过期时间管理，增加了实现的复杂性

如果仅需要简单的分布式锁功能，SETNX 可能已经足够；如果需要更多的功能，如设置键值对的同时设置有效时间等，SET 命令结合 NX 选项会更加合适.
### bigkey（大Key）
危害:
1. 客户端超时阻塞：由于 Redis 执行命令是单线程处理，然后在操作大 key 时会比较耗时，那么就会阻塞 Redis，从客户端这一视角看，就是很久很久都没有响应
1. 网络阻塞：每次获取大 key 产生的网络流量较大，如果一个 key 的大小是 1 MB，每秒访问量为 1000，那么每秒会产生 1000MB 的流量，这对于普通千兆网卡的服务器来说是灾难性的
1. 工作线程阻塞：如果使用 del 删除大 key 时，会阻塞工作线程，这样就没办法处理后续的命令

发现bigkey:
1. `--bigkeys` : `redis-cli -p 6379 --bigkeys -i 3`

    这种方式只能找出每种数据结构 top 1 bigkey
1. 使用 Redis 自带的 SCAN 命令

    |数据结构|命令|复杂度|结果（对应 key）|
    |String|STRLEN|O(1)|字符串值的长度|
    |Hash|HLEN|O(1)|哈希表中字段的数量|
    |List|LLEN|O(1)|列表元素数量|
    |Set|SCARD|O(1)|集合元素数量|
    |Sorted Set|ZCARD|O(1)|有序集合的元素数量|
    
    对于集合类型还可以使用 MEMORY USAGE 命令（Redis 4.0+），这个命令会返回键值对占用的内存空间
1. 通过分析 RDB 文件来找出 big key。这种方案的前提是你的 Redis 采用的是 RDB 持久化

    redis-rdb-tools：Python 语言写的用来分析 Redis 的 RDB 快照文件用的工具
    [rdb_bigkeys](https://github.com/weiyanwei412/rdb_bigkeys)：Go 语言写的用来分析 Redis 的 RDB 快照文件用的工具，性能更好
1. 借助公有云的 Redis 分析服务

处理:
1. 分割 bigkey：将一个 bigkey 分割为多个小 key
1. 手动清理：Redis 4.0+ 可以使用 UNLINK 命令来异步删除一个或多个指定的 key。Redis 4.0 以下可以考虑使用 SCAN 命令结合 DEL 命令来分批次删除
1. 采用合适的数据结构：例如，文件二进制数据不使用 String 保存、使用 HyperLogLog 统计页面 UV、Bitmap 保存状态信息（0/1）
1. 开启 lazy-free（惰性删除/延迟释放）：lazy-free 特性是 Redis 4.0 开始引入的，指的是让 Redis 采用异步方式延迟释放 key 使用的内存，将该操作交给单独的子线程处理，避免阻塞主线程

### hotkey（热Key）
如果一个 key 的访问次数比较多且明显多于其他 key 的话，那这个 key 就可以看作是 hotkey（热 Key）.

hotkey 出现的原因主要是某个热点数据访问量暴增.

处理 hotkey 会占用大量的 CPU 和带宽，可能会影响 Redis 实例对其他请求的正常处理.

发现hotkey:
1. 使用 Redis 自带的 --hotkeys 参数来查找

    前提:  Redis Server 的 maxmemory-policy 参数设置为 LFU 算法, 否则会报错
1. 使用 MONITOR 命令

    MONITOR 命令是 Redis 提供的一种实时查看 Redis 的所有操作的方式，可以用于临时监控 Redis 实例的操作情况，包括读写、删除等操作.
    
    由于该命令对 Redis 性能的影响比较大，因此禁止长时间开启 MONITOR（生产环境中建议谨慎使用该命令）
1. 工具

    [hotkey](https://gitee.com/jd-platform-opensource/hotkey) 这个项目不光支持 hotkey 的发现，还支持 hotkey 的处理
1. 根据业务情况提前预估

    根据业务情况来预估一些 hotkey，比如参与秒杀活动的商品数据等. 但存在无法预估的情况, 比如突发新闻
1. 业务代码中记录分析。

    在业务代码中添加相应的逻辑对 key 的访问情况进行记录分析. 不过，这种方式会让业务代码的复杂性增加，一般也不会采用.

1. 借助公有云的 Redis 分析服务

处理:
1. 读写分离：主节点处理写请求，从节点处理读请求
1. 使用 Redis Cluster：将热点数据分散存储在多个 Redis 节点上
1. 二级缓存：hotkey 采用二级缓存的方式进行处理，将 hotkey 存放一份到 JVM 本地内存中（可以用 Caffeine）
1. 公有云的 Redis 服务

    阿里云 Redis支持通过代理查询缓存功能（Proxy Query Cache）优化热点 Key 问题

### 慢查询命令
慢查询命令也就是那些命令执行时间较长的命令

Redis 中的大部分命令都是 O(1) 时间复杂度，但也有少部分 O(n) 时间复杂度的命令，例如
- KEYS *：会返回所有符合规则的 key
- HGETALL：会返回一个 Hash 中所有的键值对
- LRANGE：会返回 List 中指定范围内的元素
- SMEMBERS：返回 Set 中的所有元素
- SINTER/SUNION/SDIFF：计算多个 Set 的交集/并集/差集
- ...

不过， 这些命令并不是一定不能使用，但是需要明确 N 的值。另外，有遍历的需求可以使用 HSCAN、SSCAN、ZSCAN 代替.

还有一些时间复杂度可能在 O(N) 以上的命令，例如：
- ZRANGE/ZREVRANGE：返回指定 Sorted Set 中指定排名范围内的所有元素。时间复杂度为 O(log(n)+m)，n 为所有元素的数量，m 为返回的元素数量，当 m 和 n 相当大时，O(n) 的时间复杂度更小
- ZREMRANGEBYRANK/ZREMRANGEBYSCORE：移除 Sorted Set 中指定排名范围/指定 score 范围内的所有元素。时间复杂度为 O(log(n)+m)，n 为所有元素的数量，m 被删除元素的数量，当 m 和 n 相当大时，O(n) 的时间复杂度更小。
- ...

发现:
1. 在 redis.conf 文件中, 以使用 slowlog-log-slower-than 参数设置耗时命令的阈值，并使用 slowlog-max-len 参数设置耗时命令的最大记录条数

    `SLOWLOG GET <N>`查看, 默认最近的10条

    慢查询日志中的每个条目都由以下六个值组成：
    1. 唯一渐进的日志标识符
    1. 处理记录命令的 Unix 时间戳
    1. 执行所需的时间量，以微秒为单位
    1. 组成命令参数的数组
    1. 客户端 IP 地址和端口
    1. 客户端名称

### 缓存穿透
缓存穿透说简单点就是大量请求的 key 是不合理的，**根本不存在于缓存中，也不存在于数据库中**. 这就导致这些请求直接到了数据库上, 引发故障.

常见场景: 骇客攻击

处理:
1. 参数校验
1. 缓存无效 key

    需设置过期时间否则肯能导致跳过正常key
1. 布隆过滤器

    理论情况下添加到集合中的元素越多，误报的可能性就越大, 且存放在布隆过滤器的数据不容易删除
1. 限流

### 缓存击穿
请求的 key 对应的是 热点数据, 其**存在于数据库中，但不存在于缓存中（通常是因为缓存中的那份数据已经过期）**. 这就可能会导致瞬时大量的请求直接打到了数据库上, 引发故障.

常见场景: 秒杀活动

处理:
1. 永不过期（不推荐）：设置热点数据永不过期或者过期时间比较长
1. 提前预热（推荐）：针对热点数据提前预热，将其存入缓存中并设置合理的过期时间比如秒杀场景下的数据在秒杀结束之前不过期
1. 加锁（看情况）：在缓存失效后，通过设置互斥锁确保只有一个请求去查询数据库并更新缓存

### 缓存雪崩
缓存在同一时间**大面积**的失效，导致大量的请求都直接落到了数据库上，对数据库造成了巨大的压力, 引发故障.

常见场景: cache server宕机

处理:
1. 针对 Redis 服务不可用的情况：

    1. Redis 集群：采用 Redis 集群，避免单机出现问题整个缓存服务都没办法使用。Redis Cluster 和 Redis Sentinel 是两种最常用的 Redis 集群实现方案
    1. 多级缓存：设置多级缓存，例如本地缓存+Redis 缓存的二级缓存组合，当 Redis 缓存出现问题时，还可以从本地缓存中获取到部分数据
1. 针对大量缓存同时失效的情况：

    1. 设置随机失效时间（可选）：为缓存设置随机的失效时间，例如在固定过期时间的基础上加上一个随机值，这样可以避免大量缓存同时到期，从而减少缓存雪崩的风险
    1. 提前预热（推荐）：针对热点数据提前预热，将其存入缓存中并设置合理的过期时间，比如秒杀场景下的数据在秒杀结束之前不过期
    1. 持久缓存策略（看情况）：虽然一般不推荐设置缓存永不过期，但对于某些关键性和变化不频繁的数据，可以考虑这种策略。

### 缓存预热
常见的缓存预热方式有两种:
1. 使用定时任务来定时触发缓存预热的逻辑，将数据库中的热点数据查询出来并存入缓存中
1. 使用消息队列，比如 Kafka，来异步地进行缓存预热，将数据库中的热点数据的主键或者 ID 发送到消息队列中，然后由缓存服务消费消息队列中的数据，根据主键或者 ID 查询数据库并更新缓存。

### 保证缓存和数据库数据的一致性
Cache Aside Pattern（旁路缓存模式）:
1. 适合读请求比较多的场景, 需要同时维系 db 和 cache，并且是以 db 的结果为准,常用
1. 写: 更新数据库，然后直接删除缓存

    如果更新数据库成功，而删除缓存这一步失败的情况的话，简单说有两个解决方案：
    1. 缓存失效时间变短（不推荐，治标不治本）：让缓存数据的过期时间变短，这样的话缓存就会从数据库中加载数据。另外，这种解决办法对于先操作缓存后操作数据库的场景不适用
    1. 增加缓存更新重试机制（常用）：如果缓存服务当前不可用导致缓存删除失败的话，就隔一段时间进行重试，重试次数可以自己定。不过，这里更适合引入消息队列实现异步重试，将删除缓存重试的消息投递到消息队列，然后由专门的消费者来重试，直到成功。虽然说多引入了一个消息队列，但其整体带来的收益还是要更高一些.
1. 读: 从 cache 中读取数据，读取到就直接返回. cache 中读取不到的话，就从 db 中读取数据返回, 再把数据放到 cache 中

Read/Write Through Pattern:
1. cache 服务负责将此数据读取和写入 db，从而减轻了应用程序的职责

    少见: 常用的redis并没有提供 cache 将数据写入 db 的功能
1. 写: 先查 cache，cache 中不存在，直接更新 db. cache 中存在，则先更新 cache，然后 cache 服务自己更新 db（同步更新 cache 和 db）
1. 读: 从 cache 中读取数据，读取到就直接返回. 读取不到的话，先从 db 加载，写入到 cache 后返回响应

Write Behind Pattern（异步缓存写入）:
1. 由 cache 服务来负责 cache 和 db 的读写

    Read/Write Through 是同步更新 cache 和 db，而 Write Behind 则是只更新缓存，不直接更新 db，而是改为异步批量的方式来更新 db.

    这种方式对数据一致性带来了更大的挑战.

    消息队列中消息的异步写入磁盘、MySQL 的 Innodb Buffer Pool 机制都用到了这种策略。Write Behind Pattern 下 db 的写性能非常高，非常适合一些数据经常变化又对数据一致性要求没那么高的场景，比如浏览量、点赞量等.

### 基于Redis实现延时任务
基于 Redis 实现延时任务的功能无非就下面两种方案：
1. Redis 过期事件监听

    Redis 2.0 引入了发布订阅 (pub/sub) 功能。在 pub/sub 中，引入了一个叫做 channel（频道） 的概念，有点类似于消息队列中的 topic（主题）

    Redis 中有很多默认的 channel，这些 channel 是由 Redis 本身向它们发送消息的，而不是我们自己编写的代码。其中，`__keyevent@0__:expired` 就是一个默认的 channel，负责监听 key 的过期事件。也就是说，当一个 key 过期之后，Redis 会发布一个 key 过期的事件到`__keyevent@<db>__:expired`这个 channel 中

    这个功能被 Redis 官方称为 keyspace notifications ，作用是实时监控 Redis 键和值的变化

    缺点:
    1. 时效性差

        比如惰性删除场景
    1. 丢消息

        Redis 的 pub/sub 模式中的消息并不支持持久化. 当没有订阅者时，消息会被直接丢弃
    1. 多服务实例下消息重复消费
2. Redisson 内置的延时队列, 推荐

    Redisson 的延迟队列 RDelayedQueue 是基于 Redis 的 SortedSet 来实现的

    相比于 Redis 过期事件监听实现延时任务功能，这种方式具备下面这些优势：
    1. 减少了丢消息的可能：DelayedQueue 中的消息会被持久化，即使 Redis 宕机了，根据持久化机制，也只可能丢失一点消息，影响不大。当然了，也可以使用扫描数据库的方法作为补偿机制
    1. 消息不存在重复消费问题：每个客户端都是从同一个目标队列中获取任务的，不存在重复消费的问题

### 如何查看 Redis 内存碎片的信息
使用 info memory 命令即可查看 Redis 内存相关的信息, 详情见[这里](https://redis.io/commands/INFO)

Redis 内存碎片率的计算公式：mem_fragmentation_ratio （内存碎片率）= used_memory_rss (操作系统实际分配给 Redis 的物理内存空间大小)/ used_memory(Redis 内存分配器为了存储数据实际申请使用的内存空间大小)

mem_fragmentation_ratio （内存碎片率）的值越大代表内存碎片率越严重.

通常情况下，mem_fragmentation_ratio > 1.5 的话才需要清理内存碎片

Redis4.0及以后自带了内存整理，可以避免内存碎片率过大的问题
