# redis

现在基本没有看过还有项目使用 Memcached 来做缓存，都是直接用 Redis.

## 竞品
- [Dragonfly是一种针对现代应用程序负荷需求而构建的内存数据库](https://github.com/dragonflydb/dragonfly/blob/main/README.zh-CN.md)

### 重写redis
- [rudis](https://github.com/sleeprite/rudis)

## 通讯协议
Redis 服务器与客户端通过 RESP（Redis Serialization Protocol）协议通信, 该协议是文本协议. 它是会浪费流量，不过它的优点在于直观，非常的简单，解析性能极其的好.

> redis v6应该会采用RESP v3.

## pipeline
目的: 将多次io往返时间压缩为一次

pipeline 不适用于执行顺序有依赖关系的一批命令, 可用lua脚本代替.

## 事务
> 事务提供了一种“将多个命令打包， 然后一次性、按顺序地执行”的机制， 并且事务在执行的期间不会主动中断.

Redis 通过 MULTI 、 DISCARD 、 EXEC 和 WATCH,UNWATCH 五个命令来实现事务功能.

通过WATCH 命令监听指定的 Key，当调用 EXEC 命令执行事务时，如果一个被 WATCH 命令监视的 Key 被 其他`客户端/Session` 修改的话，整个事务都不会被执行, 同一个session的话该事务还是会被执行.

> MULTI 后可输入若干命令，并将它们放到队列中，直到调用了 EXEC 命令后，再按FIFO执行所有的命令; 也可通过 DISCARD 取消一个事务，它会清空事务队列中保存的所有命令.

> 如果在执行 WATCH 命令之后， EXEC 命令或 DISCARD 命令先被执行了的话，那么就不需要再执行 UNWATCH 了.

结论:
- Redis 具备了一定的原子性，但不支持回滚
- Redis 不具备 ACID 中一致性的概念
- Redis 具备隔离性
- Redis 通过一定策略可以保证持久性

替代方法:
1. lua : from v2.6

    一段 Lua 脚本可以视作一条命令执行，一段 Lua 脚本执行过程中不会有其他脚本或 Redis 命令同时执行，保证了操作不会被其他指令插入或打扰.

    不过，如果 Lua 脚本运行时出错并中途结束，出错之后的命令是不会被执行的。并且，出错之前执行的命令是无法被撤销的，无法实现类似关系型数据库执行失败可以回滚的那种原子性效果。因此，严格来说的话，通过 Lua 脚本来批量执行 Redis 命令实际也是不完全满足原子性的.

1. Redis functions : from v7.0

    比 Lua 更强大的脚本

pipeline 和 Redis 事务的对比：
1. 事务是原子操作，pipeline 是非原子操作. 两个不同的事务不会同时运行，而 pipeline 可以同时以交错方式执行
1. Redis 事务中每个命令都需要发送到服务端，而 Pipeline 只需要发送一次，请求次数更少

### 原子性
Redis 开始事务 multi 命令后，Redis 会为这个事务生成一个**队列**，每次操作的命令都会按照顺序插入到这个队列中. 这个队列里面的命令不会被马上执行，直到 exec 命令提交事务，所有队列里面的命令会被一次性，并且排他的进行执行.

**从严格的意义上来说 Redis 并不具备原子性: 事务中间出现了失败，并没有进行回滚**, 本质是Redis 是完成操作之后才会进行 AOF 日志记录，AOF 日志的定位只是记录操作的指令记录, 而没有类似mysq的undo日志.

### 一致性
因为 Redis 并不具备回滚，也就不具备传统意义上的原子性，所以 Redis 也应该不具备传统的一致性.

### 隔离性
Redis 因为是单线程操作，所以在隔离性上有天生的隔离机制.

### 持久性
Redis 是否具备持久化，这个取决于 Redis 的持久化模式：
- 纯内存运行，不具备持久化，服务一旦停机，所有数据将丢失；
- RDB(snapshotting) 模式，取决于 RDB 策略，只有在满足策略才会执行 Bgsave，异步执行并不能保证 Redis 具备持久化

    创建rdb的方法:
    1. save : 同步保存操作，会阻塞 Redis 主线程；
    1. bgsave : fork 出一个子进程，子进程执行，不会阻塞 Redis 主线程，默认选项
- AOF(append-only file) 模式，只有将 appendfsync 设置为 always，程序才会在执行命令同步保存到磁盘，这个模式下，Redis 具备持久化.（将 appendfsync 设置为 always，只是在理论上持久化可行，但一般不会这么操作）

    > Redis 6.0 之后已经默认是开启了

    与快照持久化相比，AOF 持久化的实时性更好.

    Redis AOF 持久化机制是在执行完命令之后再记录日志，这和关系型数据库（如 MySQL）通常都是执行命令之前记录日志（方便故障恢复）不同

    BGREWRITEAOF的AOF 重写阻塞: 缓冲区中新数据写到新文件的过程中会产生阻塞
- RDB 和 AOF 的混合持久化（Redis 4.0 新增）

与 RDB 持久化相比，AOF 持久化的实时性更好, 但aof文件比rdb大且恢复速度慢. AOF 持久化的 appendfsync 策略为 no、everysec 时都会存在数据丢失的情况。always 下可以基本是可以满足持久性要求的，但性能太差，实际开发过程中不会使用.

总结:
1. Redis 保存的数据丢失一些也没什么影响的话，可以选择使用 RDB
1. 不建议单独使用 AOF，因为时不时地创建一个 RDB 快照可以进行数据库备份、更快的重启以及解决 AOF 引擎错误
1. 如果保存的数据要求安全性比较高的话，建议同时开启 RDB 和 AOF 持久化或者开启 RDB 和 AOF 混合持久化

## 持久化
在对 Redis 进行恢复时，RDB 快照直接读取磁盘即可恢复，而 AOF 需要对所有的操作指令进行重放进行恢复，这个过程有可能非常漫长.

Redis 4.0 之后，引入了新的持久化模式，混合持久化，将 RDB 的文件和局部增量的 AOF 文件相结合来使用.

### RDB
RDB : redis database, 在指定的时间间隔内将内存中的数据集快照写入磁盘, 即快照就是一次全量的备份.

Redis 在进行 RDB 的快照生成有两种方法: Save/Bgsave. 由于 Redis 是单进程单线程，直接使用 Save，Redis 会进行一个庞大的文件 IO 操作而阻塞线上的业务, 因此通常采用Bgsave.

在使用 Bgsave 的时候，Redis 会 Fork 一个子进程，快照的持久化就交给子进程去处理，而父进程继续处理线上业务的请求.

### AOF
AOF : Append Only File, 数据操作修改的指令记录日志，类比 MySQL 的 Binlog

不过 AOF 日志也有两个比较大的问题：
1. AOF 的日志会随着时间递增，如果一个数据量大运行的时间久，AOF 日志量将变得异常庞大
1. AOF 在做数据恢复时，由于重放的量非常庞大，恢复的时间将会非常的长

针对上述的问题，Redis 在 2.4 之后也使用了 bgrewriteaof 对 AOF 日志进行瘦身: bgrewriteaof 命令用于异步执行一个 AOF 文件重写操作. 重写会创建一个当前 AOF 文件的体积优化版本.

## 删除策略
Redis 支持三种不同的删除策略:
1. 定时删除：在设置键过去的时间同时，创建一个定时器，让定时器在键过期时间来临，立即执行对键的删除操作

    对内存是友好的，但是对 CPU 的时间是最不友好的，特别是在业务繁忙的时候.
1. 惰性删除：放任键过期不管，但是每次从键空间获取键时，都会检查该键是否过期，如果过期的话，就删除该键

    惰性删除对内存来说又是最不友好的.
1. 定期删除：每隔一段时间，程序都要对数据库进行一次检查，删除里面的过期键，至于要删除多少过期键，由算法而定

    定期删除策略是前两种策略的一个整合和折中, 通过合理的删除执行的时长和频率，来达到合理的删除过期键.

## 主从复制机制
参考:
- [一篇文章带你深入解析Redis主从复制机制！](https://www.jianshu.com/p/fb579e89d0c2)

Redis的主从复制机制是指可以让从服务器(slave)能精确复制主服务器(master)的数据.

一台master服务器可以对应多台slave服务器, 另外，slave服务器也可以有自己的slave服务器，这样的服务器称为sub-slave,而这些sub-slave通过主从复制最终数据也能与master保持一致.

redis可以通过Sentinel系统管理多个Redis服务器，当master服务器发生故障时，Sentineal系统会根据一定的规则将某台slave服务器升级为master服务器,继续提供服务，实现故障转移，保证Redis服务不间断.

查看role: `redis-cli info|grep role`

### 原理
master服务器会记录一个replicationId的伪随机字符串，用于标识当前的数据集版本，还会记录一个当数据集的偏移量offset，不管master是否有配置slave服务器，replication Id和offset会一直记录并成对存在，可以通过`info replication`命令查看.

当master与slave正常连接时，slave使用PSYNC命令向master发送自己记录的旧master的replication id和offset，而master会计算与slave之间的数据偏移量，并将缓冲区中的偏移数量同步到slave，此时master和slave的数据一致.

而如果slave引用的replication太旧了，master与slave之间的数据差异太大，则master与slave之间会使用全量复制的进行数据同步.

### 复制方式
3种方式:
1. 当master服务器与slave服务器正常连接时，master服务器会发送数据命令流给slave服务器,将自身数据的改变复制到slave服务器
1. 当因为各种原因master服务器与slave服务器断开后，slave服务器在重新连上master服务器时会尝试重新获取断开后未同步的数据即部分同步，或者称为部分复制
1. 如果无法部分同步(比如初次同步)，则会请求进行全量同步，这时master服务器会将自己的rdb文件发送给slave服务器进行数据同步，并记录同步期间的其他写入，再发送给slave服务器，以达到完全同步的目的，这种方式称为全量复制

### 配置
> 主从复制的开启，完全是在从节点发起的；不需要我们在主节点做任何事情.

配置过程:
1. 在slave配置复制
1. 等待复制完成, 见`复制进度`
1. 必要时, 在slave设置SLAVEOF NO ONE, 将slave转为master

假设Redis的master服务器地址为192.168.0.101

两种方式：
1. 通过已连接到slave的redis-cli执行: `slaveof 192.168.1.101 6379` 
    
    如果master设置了密码： 在redis-cli先执行`config set masterauth xxx`
2. 配置 slave 的redis.conf: `slaveof 192.168.1.101 6379`

    如果master设置了密码： 同时设置redis.conf的masterauth

> **如果使用主从复制，那么要确保master激活了持久化，或者确保它不会在当掉后自动重启**. 原因：slave是master的完整备份，因此如果master通过一个空数据集重启，slave也会被清掉.

在Redis2.6以后，slave只读模式是默认开启的，我们可以通过配置文件中的slave-read-only选项配置.

> 如果当前服务器已经是某个主服务器(master server)的从属服务器，那么执行 SLAVEOF host port 将使当前服务器停止对旧主服务器的同步，丢弃旧数据集，转而开始对新主服务器进行同步.

> 利用`SLAVEOF NO ONE`不会丢弃同步所得数据集这个特性，可以在主服务器失败的时候，将从属服务器用作新的主服务器, 从而实现无间断运行.

> redis主从复制是[异步的](http://www.redis.cn/topics/cluster-tutorial.html).

### 复制进度
参考:
- [Redis主从同步与故障切换，有哪些坑？](https://new.qq.com/omn/20201125/20201125A0GFNT00.html)

通过redis 的`INFO replication`命令查看主库接收写命令的进度信息（master_repl_offset）和从库复制写命令的进度信息（slave_repl_offset）, 即`diff=master_repl_offset-slave_repl_offset`, diff=0为复制完成.

### 主从复制中的key过期问题
Redis处理key过期有惰性删除和定期删除两种机制，而在配置主从复制后，slave服务器就没有权限处理过期的key，这样的话，对于在master上过期的key，在slave服务器就可能被读取，因为master会累积过期的key，积累一定的量之后，发送del命令到slave，删除slave上的key.

业务层采用expireat timestamp 方式，这样命令传送到从库就没有影响, 前提是主从的时间要同步.

### 解决Redis主从架构数据丢失
通过降低min-slaves-max-lag参数的值，可以避免在发生故障时大量的数据丢失，一旦发现延迟超过了该值就不会往master中写入数据.

这种解决数据丢失的方法是降低系统的可用性来实现的.

## 哨兵
为了解决主从模式的Redis集群不具备自动容错和恢复能力的问题，Redis从2.6开始提供哨兵模式.

哨兵集群的主要作用如下：
- 监控所有服务器是否正常运行：通过发送命令返回监控服务器的运行状态（心跳），处理监控主服务器、从服务器外，哨兵之间也相互监控
- 故障切换：当哨兵监测到master宕机，会自动将某个slave切换到master，然后通过发布订阅模式, 通知其他的从服务器，修改配置文件，让它们切换master. 同时那台有问题的旧主也会变成新主的从节点.

优点：
- 哨兵模式是基于主从模式的，解决主从模式中master故障不可以自动切换master的问题

缺点：
- 浪费资源，集群里所有节点保存的都是**全量数据**，数据量过大时，主从同步会严重影响性能
- Redis主机宕机后，投票选举结束之前，Redis会开启保护机制，禁止写操作，直到选举出了新的Redis主机
- 只有一个master库执行写请求，写操作会单机性能瓶颈影响

redis sentinal致力于高可用

## cmd
```bash
redis-cli keys "*"
```

## cluster
参考:
- [Redis Cluster数据分片实现原理、及请求路由实现](https://www.huaweicloud.com/articles/38e2316d01880fdbdd63d62aa26b31b4.html)
- [Redis Cluster 会丢数据吗？](https://segmentfault.com/a/1190000021147037)

在Redis3.0中，Redis也提供Redis Cluster.

Redis集群目前无法做到数据库选择，默认在0数据库.

Redis Cluster 并没有使用一致性哈希，采用的是 哈希槽分区，每一个键值对都属于一个 hash slot（哈希槽）。当客户端发送命令请求的时候，需要先根据 key 通过上面的计算公式找到的对应的哈希槽，然后再查询哈希槽和节点的映射关系，即可找到目标 Redis 节点.

Redis Cluster是一种服务器Sharding技术(分片和路由都是在服务端实现)，采用多主多从，每一个分区都是由一个Redis主机和多个从机组成，片区和片区之间是相互平行的. Redis Cluster集群采用了P2P的模式，完全去中心化.

Redis Cluster 不保证强一致性，在一些特殊场景，客户端即使收到了写入确认，还是可能丢数据的:
- 异步复制

    在 master 写成功，但 slave 同步完成之前，master 宕机了，slave 变为 master，数据丢失.

    wait 命令可以增强这种场景的数据安全性, 但并不保证强一致性且影响性能. 它会阻塞当前 client 直到之前的写操作被指定数量的 slave 同步成功.
- 网络分区

    如果网络分区很短暂，那数据就不会丢失；但分区后一个 master 继续接收写请求， 如果在网络分区期间集群重新选举了某个slave节点为主节点，那么数据就会丢失了.

    可以设置节点过期时间，减少 master 在分区期间接收的写入数量，降低数据丢失的损失.

    > 过去一定时间(maximum window, 最大时间窗口即节点过期时间）后，分区的多数边就会进行选举，slave 成为 master， 这时分区少数边的 master 在达到过期时间后，就被认为是故障的，进入 error 状态，停止接收写请求，可以被 slave 取代.

优点:
- 更加方便地添加和移除节点，增加节点时，只需要把其他节点的某些哈希槽挪到新节点就可以了，当移除节点时，只需要把移除节点上的哈希槽挪到其他节点就行了，不需要停掉Redis任何一个节点的服务，采用一致性哈希算法时增加和移除节点需要rehash

redis sentinal致力于扩展性

## FAQ
### redis做mq
方案:
- Pub/Sub : 非常适合应用在在即时通信、游戏、消息通知等业务上, 但对于无法容忍数据丢失，消息可能积压的场景不太适合

    缺点：

    1. 消息没有持久化的机制. 在Pub/Sub模型中，消费者是和连接（Connection）绑定的，当消费者的连接断掉（网络原因或者消费者进程crash）后，再次重连，那么Channel中的消息将永久消失（对于该消费者而言），也就是说Pub/Sub模型缺少消息回溯的机制
    1. 消费消息的速度和消费者的数量成反比. 在Redis的实现中，Redis会把Channel中的消息逐个（Linear）推送给每个消费者，因此当消费者的数量达到一定规模时，服务器的性能将线性下降，因此每个消费者获取到消息的延迟也线性增长
    1. 当生产者产生消息的速度远大于消费者的消费能力的时候（此时可以简单地理解为消息积压），消费者会被强制断开连接，因此会造成消息的丢失，这个特性可以详见redis的配置(`client-output-buffer-limit`)
- List

    优点:

    1. 消息可以持久化
    1. 消息可以积压

    缺点:

    1. 在消费者消费之前，对消息进行处理，把该消息写入到若干个队列中，这样能支持多个消费者同时消费，但是数据却被拷贝了多次
    1. 重复消费时只能再次加入List
- ZSet
    优点:

    1. 消息可以持久化
    1. 支持消息重复消费(ZRANGEBYSCORE不会删除消息)

    缺点:

    1. 消息的顺序. score至关重要，这关系到消息的先后顺序. 一种可行的方案是使用timestamp+seq作为score，这样能够保证消息的顺序
    1. 重复消息的添加. member即message body, 由于有序集合的限制, 无法添加重复member
- stream

### `READONLY You can't write against a read only slave`
Redis 主服务器（master server）具有读写的权限，而 从服务器（slave master）默认 只具有读的权限, 如果在从服务器中写数据则会报该错误.

解决方法:
1. 梳理逻辑不在slave上写
1. 运行在slave上写

    - 修改 redis.conf 配置文件中的参数`slave-read-only  yes`，将 yes 修改为 no, 然后保存并重启 redis 服务， 此时从服务器就具备了读写权限
    - 在redis-cli中执行`config set slave-read-only no`