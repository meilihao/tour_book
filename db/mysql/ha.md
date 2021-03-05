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

### 2. MMM(Master-Master Replication Manager for MySQL), 已被MHA替代
提供了MySQL主主复制配置的监控, 故障转移和管理的一套可伸缩的脚本套件.

典型: 双主多从架构

通过MySQL复制实现两个server互为主从, 且在任何时刻只有一个节点可写, 避免多点写入冲突; 可写入的server故障时触发切换.

故障转移: keepalived

### 3. heartbeat/SAN
故障转移: 高可用集群软件heartbeat.
数据共享: SAN存储共享数据. 正常状态由主节点挂载读写; 故障时, heartbeat通过仲裁设备将主节点挂载的存储释放, 然后让备用节点挂载.

> 基于心跳, 可能脑裂

### 4. heartbeat/DRDB
参考:
- [MySQL高可用之DRBD](https://wxy0327.blog.csdn.net/article/details/103070764)

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

### 复制方式
1. 异步复制（Asynchronous replication）

    MySQL默认的复制即是异步的，主库在执行完客户端提交的事务后会立即将结果返给给客户端，并不关心从库是否已经接收并处理，这样就会有一个问题，主如果crash掉了，此时主上已经提交的事务可能并没有传到从上，如果此时，强行将从提升为主，可能导致新主上的数据不完整.
1. 全同步复制（Fully synchronous replication）, mysql不支持

    指当主库执行完一个事务，所有的从库都执行了该事务才返回给客户端. 因为需要等待所有从库执行完该事务才能返回，所以全同步复制的性能必然会收到严重的影响.
1. 半同步复制（Semisynchronous replication）, **推荐**

    介于异步复制和全同步复制之间，主库在执行完客户端提交的事务后不是立刻返回给客户端，而是等待至少一个从库接收到并写到relay log中才返回给客户端. 相对于异步复制，半同步复制提高了数据的安全性，同时它也造成了一定程度的延迟，这个延迟最少是一个TCP/IP往返的时间. 所以，半同步复制最好在低延时的网络中使用.

## 读写分离
中间件.

## 半复制+GDIT
```conf
# vim /etc/mysql/my.cnf
[mysqld]
bind-address=0.0.0.0

binlog-ignore-db=mysql
binlog-ignore-db=information_schema
binlog-ignore-db=performance_schema
# 或
binlog_do_db=my_database

[mariadb]
...
rpl_semi_sync_master_enabled=ON
rpl_semi_sync_slave_enabled=ON
rpl_semi_sync_master_wait_point=AFTER_SYNC

gtid_strict_mode=ON

log-bin
server_id=1001/1002
binlog-format=mixed
```

create user 'repl'@'%' identified by '123456';
GRANT REPLICATION SLAVE, REPLICATION SLAVE ADMIN, REPLICATION MASTER ADMIN, REPLICATION SLAVE ADMIN, BINLOG MONITOR, SUPER ON *.* to 'repl'@'%';

CHANGE MASTER TO MASTER_HOST='192.168.16.44',MASTER_PORT=3306,MASTER_USER='repl',MASTER_PASSWORD='123456',master_use_gtid=current_pos;start slave;show slave status;


相关命令:
```sql
> show binlog status \G; -- 查看master status, 查看当前使用的binlog file. 等价于以前的`show master status`
> SHOW BINLOG EVENTS in 'mysqld-bin.000002'; -- 查看binlog内容, 不加指定log_name(mysqld-bin.000002)时只能显示第一个binlog
> SHOW relay EVENTS; -- 查看relaylog内容

> show slave status \G; -- 查看slave status

> show variables like '%gtid%' -- 查看gtid

-- 开始slave
> SET GLOBAL gtid_slave_pos = "0-1-2"; -- BINLOG_GTID_POS 函数得到的位置, gtid_slave_pos为空时, 从master的第一个gtid开始复制
> CHANGE MASTER TO master_host="127.0.0.1", master_port=3310, master_user="root", master_use_gtid=slave_pos;
> START SLAVE;
```

> Slave的执行状态（最后一个执行的 GTID）被记录在 mysql.gtid_slave_pos 系统表中