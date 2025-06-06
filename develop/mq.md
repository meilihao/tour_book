# mq
参考:
- [90%的Java程序员，都扛不住这波消息中间件的面试四连炮！](http://www.liuhaihua.cn/archives/587877.html)
- [云原生消息系统设计：NATS + Cloudevents](https://wbsnail.com/p/using-cloudevents-with-nats)

	[NATS Streaming Server supports clustering and data replication, implemented with the Raft consensus algorithm, for the purposes of high availability.](https://docs.nats.io/legacy/stan/intro/clustering)

使用MQ的场景有挺多的，但是比较核心的有3个：异步、解耦、削峰填谷.

## rabbitmq
rabbitmq解决消息丢失:
1. 持久化机制
1. 关闭RabbitMQ消费者的自动提交ack,在消费者处理完这条消息之后再手动提交ack

## rocketmq

### FAQ
#### broker卡死
内存和swap被耗光导致系统卡死.

解决方法:
修改broker的启动参数
```sh
$ vim vim rocketmq-all-4.3.0-bin-release/bin/runbroker.sh
...
# JAVA_OPT="${JAVA_OPT} -server -Xms8g -Xmx8g -Xmn4g" # 原始参数
JAVA_OPT="${JAVA_OPT} -server -Xms4g -Xmx4g -Xmn1g"
...
```

#### rocketmq-console: org.apache.rocketmq.remoting.exception.RemotingConnectException: connect to <null> failed
解决:
1. `export NAMESRV_ADDR=localhost:9876` 之后再运行`mvn spring-boot:run`
1. 编辑`rocketmq-console/src/main/resources/application.properties`的`rocketmq.config.namesrvAddr=`

### 消息堆积
原因:
1. 消费者因bug, 宕机, 网络等无法消费
2. 生产者推送大量消息给broker, 消费者消费不及时

解决:
1. 事前: 预估+压侧
1. 事中: 

	1. 临时扩容

		1. 添加消费者
		2. 添加消息队列数量
	2. 临时存储

		添加临时消费者, 暂存消息到db, 延后处理
	3. 业务降级

		1. 减小消息的产生
		2. 减少消费者业务逻辑
1. 事后

	1. 优化消费者, 提高消费速度, 比如并行
	2. 扩容

### 保证消息不被重复消费
原因:
1. 生产者重复产生消息, 比如接口的幂等操作
1. 消费者消费后且回复ack前mq挂了
1. 消费者消费后且回复ack前消费者挂了

解决:
1. 记录消费记录

	1. 插入消费记录, 状态是消费中, 且设置过期时间
	2. 如果执行成功, 状态改为完成
	3. 如果失败, 删除该记录并重试

### 保证消息不丢失
原因:
1. mq收到消息后, 保存失败
1. 队列中消息未持久化
1. 消费者开启自动应答, 且消费失败

解决:
1. 开启生产者确认
2. 持久化队列中的消息
3. 手动ack

### 保证消息队列的顺序性
1. 顺序消息按顺序进入同一个队列
2. 消息携带序号, 消费者消费前检查序号

### rocksmq持久化
RocketMQ 在持久化的设计上，采取的是消息顺序写、随机读的策略，利用磁盘顺序写的速度，让磁盘的写速度不会成为系统的瓶颈。并且采用 MMPP 这种“零拷贝”技术，提高消息存盘和网络发送的速度。极力满足 RocketMQ 的高性能、高可靠要求.

在 RocketMQ 持久化机制中，涉及到了三个角色：
- CommitLog：消息真正的存储文件，所有消息都顺序写入 CommitLog 文件中

	在 RocketMQ 中提供了同步刷盘和异步刷盘两种刷盘方式，可以通过 Broker 配置文中中的 flushDiskType 参数来设置（SYNC_FLUSH、ASYNC_FLUSH）:
	- 异步刷盘方式（默认）：消息写入到内存的 PAGECACHE中，就立刻给客户端返回写操作成功，当 PAGECACHE 中的消息积累到一定的量时，触发一次写操作，将 PAGECACHE 中的消息写入到磁盘中。这种方式吞吐量大，性能高，但是 PAGECACHE 中的数据可能丢失，不能保证数据绝对的安全
	- 同步刷盘方式：消息写入内存的 PAGECACHE 后，立刻通知刷盘线程刷盘，然后等待刷盘完成，刷盘线程执行完成后唤醒等待的线程，返回消息写成功的状态。这种方式可以保证数据绝对安全，但是吞吐量不大
- ConsumeQueue：消息消费逻辑队列，类似于 MySQL 中的二级索引
- IndexFile：消息索引文件，主要存储消息 Key 与 offset 对应关系，提升消息检索速度
