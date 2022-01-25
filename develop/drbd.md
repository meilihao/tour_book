# DRBD
参考:
- [DRBD使用](https://documentation.suse.com/zh-cn/sle-ha/11-SP4/html/SLE-ha-guide/cha-ha-drbd.html)

DRBD (Distributed Replicated Block Device)分布式复制块设备是一种基于软件和网络(tcp/ip和RDMA)的**块复制**存储解决方案, 即多写.

DRBD是Linux内存存储层中的一个分布式存储系统，具体来说两部分构成:
1. 内核模板，主要用于虚拟一个块设备
1. 用户控件管理程序，主要用于和DRBD内核模块通信，以管理DRBD资源

    在DRBD中，资源主要包含DRBD设备，磁盘配置，网络配置等.

一个DRBD系统有两个及以上的节点构成.

DRBD设备在整个DRBD系统中位于逻辑块设备上，文件系统(及buffer/cache)之下，在文件系统和磁盘之间形成了一个中间层，当用户在主节点的文件系统中写入数据时，数据被正式写入磁盘前会被DRBD系统截获，同时，DRBD会通知用户控件管理程序把这些数据复制一份, 写入远程主机的DRBD镜像.

![DRBD的工作流程图](http://s3.51cto.com/wyfs02/M00/25/7A/wKiom1NgbRPzAB4LAABNdBRd5XE362.gif)

DRBD在远程传输上支持三种模式：
1. 异步(protocol A)：所谓异步就是指数据只需要发给本地的TCP/IP协议栈就可以了，本地存完就OK；而DRBD只需要把数据放到TCP/IP协议栈，放到发送队列中准备发送就返回了；这种方式更高效
1. 半同步(protocol B)：数据已经发送到对方的TCP/IP协议栈上，对方的TCP/IP协议栈已经把数据接收下来了就返回，数据存不存下来就不管了
1. 同步(protocol C)：数据必须确保对方把数据写入对方的磁盘再返回的就叫同步；**这种方式数据更可靠, 因此最常用**

> DRBD8.0以后的版本支持双主模式, 此时需要共享的集群文件系统(GFS，OCFS2等)来解决并发读写问题, 通过集群文件系统的分布式锁机制解决集群中两个主节点同时操作数据的问题.

> **应用写drbd device和drbd写底层disk的bio偏移量是相同的**.

## cmd
```bash
# drbdadm status my_res : 查看资源状态, 包含role
# --- 创建drbd
# drbdadm -c /etc/drbd.conf create-md r0 --force # 创建drbd metadata
# drbdadm -c /etc/drbd.conf up r0 # 创建drbd device
# drbdadm -c /etc/drbd.conf primary r0 --force # 设为primary
# drbdadm -c /etc/drbd.conf status r0
# --- 其他
# drbdadm -c /etc/drbd.conf down r0 # down drbd device
# drbdadm -c /etc/drbd.conf secondary r0 --force # set secondary role
## --- 扩容
# zfs set volsize=20G mypool/test
# drbdadm resize r0
# e2fsck -f /dev/drbd0
# resize2fs /dev/drbd0
```

## FAQ
### drbd 9 role自动变primary
drbd 9在所有都是secondary情况下, 某个drbd device一旦写入数据会自动变成primary
### 查看drbd的同步状态
- `cat /proc/drbd`, 其中的索引是drbd设备的次设备号.

### drbd的底层设备为什么还能挂载成功即不受drbd顶层设备的影响
1. `meta-disk internal`

    drbd将磁盘分为了两部分: 前半部分是数据区,后半部是元数据区. 因为fs在磁盘的前半段不影响挂载, 就相当于创建fs时指定了范围即只使用部分磁盘空间.
1. `meta-disk <dev>`

    就相当于fs使用了整个磁盘, 还是不影响fs挂载

1. drbd设备配置

    1. addr里不能使用localhost, 但可以使用"127.0.0.1"
    1. addr中的端口必须在[1~65535]中

### drbd secondary 设备挂载报错
drbd规定mount操作只能在primary节点进行.

### requested minor out of range
drbd设备超过限制, 目前了解最大是2^20, 已验证过的最大值是150000

### conflicting use of IP 'xxx:65534'
该端口虽然是空闲的, 但已配置在其他xxx.res中, 因此还是不能使用.

### 查看版本
`modinfo drbd`

### 多primariy
/etc/drbd.d/global-common.conf设置allow-two-primaries, 由上层gfs2实现锁防止多写.

### drbd未向底层disk(/dev/zd0)写入数据
env: drbd 9.1.2

r0.res: node2的hostname是不存在的, disk与node1相同， address是127.0.0.1， 端口与node1相同. protocol是C.

用dd向/dev/drbd0写入数据, zd0相应位置没有数据, 说明drbd0缓存了数据(能缓存多大多久, 未知)， `sync`也没有效果, 但当`down r0`时drbd会将缓存数据写入zd0.

### `unknown filesysem type 'drbd'`
直接挂载drbd的底层disk报该错. 解决方法: 创建drbd device来操作该disk.

### [drbd matedata大小](https://linbit.com/drbd-user-guide/drbd-guide-9_0-cn/#s-meta-data-size)

### 谁占用drbd
用`cat /sys/kernel/debug/drbd/resources/r46/volumes/0/openers`查看.

### drbdadm create-md报错`... device-minor 'device-minor:<node_peer_hostname>:40' first used here`
本端为对端配置的drbd index是40, 但实际上40已在对端存在.

### drbd device配置中是否所有node上的index都必须相同?
可以不同.