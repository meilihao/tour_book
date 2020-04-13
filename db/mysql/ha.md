# ha
推荐**tidb**.

ha需要解决的两个问题:
1. 数据共享或数据同步

    - 数据共享: SAN(storage area network)
    - 数据同步: rsync或DRBD
1. 故障转移 : 当server死机或出现错误时自动切换到其他备用server而不影响业务.

## 解决方案
### 1. 主从复制
同步方式: binlog复制
故障转移: keepalived

> 从服务器只读

### 2. MMM(Master-Master Replication Manager for MySQL)
提供了MySQL主主复制配置的监控, 故障转移和管理的一套可伸缩的脚本套件.

典型: 双主多从架构

通过MySQL复制实现两个server互为主从, 且在任何时刻只有一个节点可写, 避免多点写入冲突; 可写入的server故障时触发切换.

故障转移: keepalived

### 3. heartbeat/SAN
故障转移: 高可用集群软件heartbeat.
数据共享: SAN存储共享数据. 正常状态由主节点挂载读写; 故障时, heartbeat通过仲裁设备将主节点挂载的存储释放, 然后让备用节点挂载.

> 基于心跳, 可能脑裂

### 4. heartbeat/DRDB
故障转移: 高可用集群软件heartbeat.
数据共享: 基于块的数据同步软件DRBD

> 基于心跳, 可能脑裂

### 5. MySQL Cluster
复杂, 实际应用不多.

## mysql复制
常用架构:
1. 一主一从 : 最常见
1. 一主多从 : 写操作不频繁; 查询频繁
1. 主主互备 : 写操作频繁 
1. 双主多从 : 读写都频繁

原则:
1. 同一时刻只有一个server在写
1. 一个主sever可有多个从server
1. 无论主从必须保证自身的id唯一, 否则双主互备会出问题

    当从库的io_thread发现binlog event的源与自己的server-id相同时，就会跳过该event, 导致数据遗失.
1. 从server可以级联

## 读写分离
中间件.