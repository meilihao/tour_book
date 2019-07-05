# redis
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
3中方式:
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

> 如果使用主从复制，那么要确保master激活了持久化，或者确保它不会在当掉后自动重启. 原因：slave是master的完整备份，因此如果master通过一个空数据集重启，slave也会被清掉.

在Redis2.6以后，slave只读模式是默认开启的，我们可以通过配置文件中的slave-read-only选项配置.

### 主从复制中的key过期问题
Redis处理key过期有惰性删除和定期删除两种机制，而在配置主从复制后，slave服务器就没有权限处理过期的key，这样的话，对于在master上过期的key，在slave服务器就可能被读取，因为master会累积过期的key，积累一定的量之后，发送del命令到slave，删除slave上的key.

业务层采用expireat timestamp 方式，这样命令传送到从库就没有影响, 前提是主从的时间要同步.