# ha
参考:
- [pacemaker+corosync/heartbeat对比及资源代理RA脚本](https://www.cnblogs.com/clsblog/p/6202869.html)
- [<<DRBD权威指南——基于Corosync+Heartbeat技术构建网络RAID>>]

[驱动、开发者和Linux厂商，以及整个开源高可用集群社区，都已经转移到了基于Corosync 2.x+Pacemaker的HA堆栈上](http://www.linux-ha.org/wiki/Site_news), [Heartbeat](http://www.linux-ha.org)已名存实亡.

## example
- [ZFS High-Availability NAS](https://github.com/ewwhite/zfs-ha/wiki)

## corosync
coreosync在传递信息的时候可以通过一个简单的配置文件来定义信息传递的方式和协议等. 因此corosync是用于高可用环境中的提供通讯服务的，它位于高可用集群架构中的底层（Message Layer），扮演着为各节点（node）之间提供心跳信息传递这样的一个角色.

## pacemaker
pacemaker是一个开源的高可用资源管理器(CRM)，位于HA集群架构中资源管理、资源代理(RA)这个层次，它不能提供底层心跳信息传递的功能，要想与对方节点通信需要借助底层的心跳传递服务，将信息通告给对方. 通常它与corosync的结合方式有两种:
- pacemaker作为corosync的插件运行
- pacemaker作为独立的守护进程运行

pacemaker负责仲裁指定谁是活动节点、IP地址的转移、本地资源管理系统.

当故障节点修复后，资源返回来称为failback，当故障节点修复后，资源仍在备用节点，称为failover.

## corosync & pacemaker高可用解决方案