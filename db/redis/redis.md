# redis

## 通讯协议
Redis 服务器与客户端通过 RESP（Redis Serialization Protocol）协议通信, 该协议是文本协议. 它是会浪费流量，不过它的优点在于直观，非常的简单，解析性能极其的好.

> redis v6应该会采用RESP v3.

## 事务
> 事务提供了一种“将多个命令打包， 然后一次性、按顺序地执行”的机制， 并且事务在执行的期间不会主动中断.

Redis 通过 MULTI 、 DISCARD 、 EXEC 和 WATCH,UNWATCH 五个命令来实现事务功能.

> 如果在执行 WATCH 命令之后， EXEC 命令或 DISCARD 命令先被执行了的话，那么就不需要再执行 UNWATCH 了.

结论:
- Redis 具备了一定的原子性，但不支持回滚
- Redis 不具备 ACID 中一致性的概念
- Redis 具备隔离性
- Redis 通过一定策略可以保证持久性

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
- RDB 模式，取决于 RDB 策略，只有在满足策略才会执行 Bgsave，异步执行并不能保证 Redis 具备持久化
- AOF 模式，只有将 appendfsync 设置为 always，程序才会在执行命令同步保存到磁盘，这个模式下，Redis 具备持久化.（将 appendfsync 设置为 always，只是在理论上持久化可行，但一般不会这么操作）

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
假设Redis的master服务器地址为192.168.0.101

两种方式：
1. 通过已连接到slave的redis-cli执行: `slaveof 192.168.1.101 6379` 
    
    如果master设置了密码： 在redis-cli先执行`config set masterauth xxx`
2. 配置 slave 的redis.conf: `slaveof 192.168.1.101 6379`

    如果master设置了密码： 同时设置redis.conf的masterauth

> **如果使用主从复制，那么要确保master激活了持久化，或者确保它不会在当掉后自动重启**. 原因：slave是master的完整备份，因此如果master通过一个空数据集重启，slave也会被清掉.

在Redis2.6以后，slave只读模式是默认开启的，我们可以通过配置文件中的slave-read-only选项配置.

### 主从复制中的key过期问题
Redis处理key过期有惰性删除和定期删除两种机制，而在配置主从复制后，slave服务器就没有权限处理过期的key，这样的话，对于在master上过期的key，在slave服务器就可能被读取，因为master会累积过期的key，积累一定的量之后，发送del命令到slave，删除slave上的key.

业务层采用expireat timestamp 方式，这样命令传送到从库就没有影响, 前提是主从的时间要同步.

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