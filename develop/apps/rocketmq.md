# rocketmq
ref:
- [RocketMQ DLedger多副本即主从切换的实现原理](https://www.yisu.com/zixun/598918.html)
- [解读 RocketMQ 5.0 全新的高可用设计](https://my.oschina.net/yunqi/blog/10086167)

> 从4.5开始RocketMQ集群从原先的支持主从同步升级到可支持主从切换

RocketMQ DLedger集群支持主从切换, 它至少需要3台机器.

## FAQ
### ledger是如何实现主从自动切换的
Dledger自己就有一套CommitLog机制，如果使用了它，它接到数据第一步就是写入自己的CommitLog.

所以，引入Dledger技术，其实就是使用Dledger的CommitLog来替换掉Broker自己的CommitLog, 然后Broker仍然可以基于Dledger的CommitLog，把消息的位置信息保存到ConsumeQueue中.

### Dledger的数据同步机制
Dledger是通过Raft协议进行多副本同步的，简单来讲，数据同步分为两个阶段，uncommitted阶段和committed阶段。

首先，当Leader接到消息数据后，会先标记消息为uncommitted状态，然后通过Dledger的组件把uncommitted状态的消息发送给Follower上的DledgerServer。

接着Follower接到消息后，会发送一个ack给Leader上的DledgerServer，然后如果Leader发现超过半数的Follower已经给自己返回了ack，那么就认为同步成功了，这时候把状态改为committed。

然后再发消息给Follower，将Follower上的状态也改为committed。

这就是基于Dledger的数据同步机制.