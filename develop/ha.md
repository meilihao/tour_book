# cluster
cluster是一组协同工作的服务集合, 用来提供比单一服务更稳定, 高效及扩展性的服务平台. 它通常由两个及以上的server组成.

集群特定:
- 高可用性

    出现服务故障时, 自动从故障节点切换到备用节点(常用做法是漂移vip), 从而提供不间断的服务, 保障业务的持续性.
- 可扩展性

    动态加入节点, 增强cluster的整体性能

常见ha(high availability cluster)的形式: 双机热备(active/standby), 双机互备, 多机互备等.
常见lb(load balance cluster): haproxy ,lvs集群, F5 Networks等.
常见分布式计算集群: hadoop, spark.

# ha
参考:
- [pacemaker+corosync/heartbeat对比及资源代理RA脚本](https://www.cnblogs.com/clsblog/p/6202869.html)
- [<<DRBD权威指南——基于Corosync+Heartbeat技术构建网络RAID>>]

[驱动、开发者和Linux厂商，以及整个开源高可用集群社区，都已经转移到了基于Corosync 2.x+Pacemaker的HA堆栈上](http://www.linux-ha.org/wiki/Site_news), [Heartbeat](http://www.linux-ha.org)已名存实亡.

## ha概念
- node

    运行心跳进程的一个独立主机(host), 即为节点(node).
- resource

    一个node可控制的实体, 当该节点发生故障时, 这些资源可被其他node接管, 常见的资源有:
    - 磁盘, 文件系统
    - ip
    - 应用服务
- event

    cluster中可能发生的事
- action

    event发生时ha的响应行为, 通常使用脚本来处理.

## example
- [ZFS High-Availability NAS](https://github.com/ewwhite/zfs-ha/wiki)

## corosync
coreosync在传递信息的时候可以通过一个简单的配置文件来定义信息传递的方式和协议等. 因此corosync是用于高可用环境中的提供通讯服务的，它位于高可用集群架构中的底层（Message Layer），扮演着为各节点（node）之间提供心跳信息传递这样的一个角色.

## pacemaker
参考:
- [Pacemaker Administration](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/index.html)
- [从头开始搭建集群 在Fedora上面创建主/主和主备集群](https://clusterlabs.org/pacemaker/doc/deprecated/zh-CN/Pacemaker/1.1-plugin/html-single/Clusters_from_Scratch/index.html)

pacemaker是一个开源的高可用集群资源管理器(CRM)，位于HA集群架构中资源管理、资源代理(RA)这个层次，它不能提供底层心跳信息传递的功能，要想与对方节点通信需要借助底层的心跳传递服务，将信息通告给对方. 通常它与corosync的结合方式有两种:
- pacemaker作为corosync的插件运行
- pacemaker作为独立的守护进程运行

pacemaker负责仲裁指定谁是活动节点、IP地址的转移、本地资源管理系统.

当故障节点修复后，资源返回来称为failback，当故障节点修复后，资源仍在备用节点，称为failover.

[pacemaker组建重命名(1.x->2.x)](https://wiki.clusterlabs.org/wiki/Pacemaker_2.0_Daemon_Changes):
- attrd 	pacemaker-attrd

    node attribute manager
- cib	pacemaker-based

    群集信息库管理者. 包含所有群集选项，节点，资源，他们彼此之间的关系和现状的定义. 会同步更新到所有群集节点

    使用xml描述集群的配置及集群中所有资源的当前状态
- crmd 	pacemaker-controld

    集群资源管理守护进程. 主要是消息代理的PEngine和LRM，还选举一个领导者（DC）统筹活动（包括启动/停止资源）的集群
- lrmd 	pacemaker-execd

    本地资源管理守护进程(local resource agent executor). 它提供了一个通用的接口支持的资源类型。直接调用资源代理（脚本）
- stonithd 	pacemaker-fenced

    STONITH(Shoot the Other Node in the Head), 强制使节点下线.

    因为如果一个节点没有相应，但并不代表它没有在提供服务，100%保证数据安全的做法就是在允许另外一个节点操作数据之前，使用STONITH来保证节点真的下线了.
- pacemaker_remoted 	pacemaker-remoted

    remote resource agent executor
- pengine 	pacemaker-schedulerd
    
    策略引擎, action scheduler

    主要负责将CRM发过来的一些信息按照配置文件中的各种设置（基于目前的状态和配置）计算集群的下一个状态.

![Pacemaker Architecture 2.x](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-internals.png)

pacemaker支持的集群类型:
- Active/Passive(主从)

    ![](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-active-passive.png)
- Active/Active

    ![](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-active-active.png)
- N+1
- N+M
- N-to-1
- N-to-N

资源代理(resource agent)是一种标准化的集群接口, packmaker通过该接口对集群资源进行操作.

pacemaker支持的资源代理:
- LSB(linux standard base resource agents), 即sysv/systemd脚本
- ocf(open cluster framework)资源代理, 是对LSB资源代理的扩展, 在`/usr/lib/ocf/resource.d/provider`
- STONITH

集群的全生命周期管理工具:
- pcs

    专用于pacemaker+corosync的设置工具,有CLI和web-based GUI界面
- crmsh, **推荐**

    是管理pacemaker的命令行界面

> [pcs与crmsh命令比较](https://www.cnblogs.com/gzxbkk/p/7305227.html)或[pcs-crmsh-quick-ref.md](https://github.com/ClusterLabs/pacemaker/blob/master/doc/pcs-crmsh-quick-ref.md)

## corosync & pacemaker高可用解决方案
参考:
- [Corosync+Pacemaker+crmsh构建Web高可用集群](https://www.cnblogs.com/cloudos/p/8336529.html)

环境准备for 每个node:
```bash
# ---每个node都执行
# chronyc makestep # 同步时间, corosync的要求, 最好也设置crontab.
# vim /etc/hosts # 设置hosts
...
192.168.0.41 node11
192.168.0.42 node12
# apt install nginx
# vim /var/www/html/index.nginx-debian.html # 修改nginx default page, 在<title>中加入每个node的hostname
# curl localhost # 验证上一步的修改是否已生效
# apt install pacemaker corosync pcs crmsh
# systemctl start pcsd && systemctl enable pcsd
# echo hacluster:123456 | chpasswd # for ubuntu. `echo "123456" | passwd --stdin hacluster` for centos, 为hacluster用户(created by pcs)设置密码
# ---在任一节点执行
# pcs cluster auth node11  node12 -u hacluster # 认证节点
# pcs cluster setup --force --name mycluster node11  node12 # 配置集群, 会自动创建/etc/corosync/corosync.onf, 并同步到所有node
# pcs cluster start --all # 启动cluster
# pcs cluster enable --all # 设置自动启动
# pcs status # 检测cluster status
...
Current DC: node11 (version 1.1.18-2b07d5c5a9) - partition with quorum # 当前的仲裁节点(Current DC, DC即为Designated Co-ordinator）为node11, 这个节点负责向集群中的节点发出一系列指令，使各个资源按照定义（存储在cib数据库中）启动或停止
...
# corosync-cfgtool -s # 验证集群状态, 仅显示当前node的状态
# corosync-cmapctl # 检查集群成员关系
# pcs status corosync # 检查corosync状态
# pcs property set stonith-enabled=false # pcs/crmsh可能会提示stonith resource未定义, 由于没有stonith设备可以通过修改集群全局属性将此error先关掉
# --- 基础的集群已经配置完毕
```

查看crm ra支持的资源:
```bash
crm(live)# cd ra
crm(live)ra# classes # 查看cluster支持的类型
lsb
ocf / .isolation heartbeat pacemaker redhat
service
stonith
systemd
crm(live)ra# list lsb # 查看该类别下可用的ra
acpid                     apparmor                  apport                    atd                       console-setup.sh          corosync                  cron                      cryptdisks
cryptdisks-early          dbus                      ebtables                  grub-common               heartbeat                 hwclock.sh                irqbalance                iscsid
keyboard-setup.dpkg-bak   keyboard-setup.sh         kmod                      logd                      lvm2                      lvm2-lvmetad              lvm2-lvmpolld             lxcfs
lxd                       mdadm                     mdadm-waitidle            networking                nfs-common                nginx                     nmbd                      open-iscsi
open-vm-tools             openhpid                  pacemaker                 pcsd                      plymouth                  plymouth-log              procps                    resolvconf
rpcbind                   rsync                     rsyslog                   samba-ad-dc               screen-cleanup            smbd                      ssh                       sysstat
udev                      ufw                       unattended-upgrades       uuidd                     zfs-share                 
crm(live)ra# list ocf heartbeat
AoEtarget            AudibleAlarm         CTDB                 ClusterMon           Delay                Dummy                EvmsSCC              Evmsd                Filesystem           ICP
IPaddr               IPaddr2              IPsrcaddr            IPv6addr             LVM                  LVM-activate         LinuxSCSI            MailTo               ManageRAID           ManageVE
NodeUtilization      Pure-FTPd            Raid1                Route                SAPDatabase          SAPInstance          SendArp              ServeRAID            SphinxSearchDaemon   Squid
Stateful             SysInfo              VIPArip              VirtualDomain        WAS                  WAS6                 WinPopup             Xen                  Xinetd               ZFS
anything             apache               asterisk             aws-vpc-move-ip      aws-vpc-route53      awseip               awsvip               clvm                 conntrackd           db2
dhcpd                dnsupdate            docker               eDir88               ethmonitor           exportfs             fio                  galera               garbd                iSCSILogicalUnit
iSCSITarget          ids                  iface-bridge         iface-vlan           iscsi                jboss                kamailio             ldirectord           lvmlockd             lxc
minio                mysql                mysql-proxy          nagios               named                nfsnotify            nfsserver            nginx                oraasm               oracle
oralsnr              ovsmonitor           pgagent              pgsql                pingd                portblock            postfix              pound                proftpd              rabbitmq-cluster
redis                rkt                  rsyncd               rsyslog              scsi2reservation     sfex                 sg_persist           slapd                symlink              syslog-ng
tomcat               varnish              vmware               vsftpd               zabbixserver         
crm(live)ra# info ocf:heartbeat:IPaddr # 查看该ra的help
```

创建/修改资源:
```bash
# crm
crm(live)# configure
crm(live)configure# primitive webip ocf:heartbeat:IPaddr params ip=192.168.10.160 cidr_netmask=24 # 定义webip资源
crm(live)configure# edit webip # 修改webip资源
crm(live)configure# primitive webserver systemd:nginx # 定义webserver资源
crm(live)configure# group webservice webip webserver # 定义组资源，组名为webservice，使webip和webserver服务在同一节点上. **组中在前面的先启动**
crm(live)configure# verify # 验证配置
crm(live)configure# commit # 提交配置
...
crm(live)node# standby node11 # 模拟node下线, 查看webip所在
crm(live)node# online node11 # node重新上线
crm(live)node# cd ..
crm(live)# exit
# crm_verify -L -V # 检查crm配置是否正确
# systemctl stop pacemaker corosync # on node12, 其上的pacemaker resource会全部撤销
# ip addr # on node11, 看到webip出现
```

或
```bash
# pcs resource create vip ocf:heartbeat:IPaddr2 ip=192.168.56.24 cidr_netmask=24 op monitor interval=30s
# crm node standby # 将node设为standby
```

corosync+pacemaker集群默认对节点高可用，但是对于节点上资源的运行状态无法监控，因此，需要配置集群对于资源的监控，在资源因意外情况下，无法提供服务时，对资源提供高可用.

```bash
# crm configure 
crm(live)configure# monitor webip 30s:20s
crm(live)configure# monitor webserver 30s:100s
crm(live)configure# verify
crm(live)configure# commit
```

# keepalived
linux下轻量级ha解决方案, 主要通过虚拟路由冗余来实现ha. 与heartbeat相比, 没有其功能强大, 但在某些场景下更简单合适.
keepalived起初是为LVS设计, 专门监控(根据3~5层的交换机制)集群中各节点的状态, 通过剔除故障节点和加入正常节点来实现ha.
keepalived支持VRRP(Virtual Router Redundancy Protocol, 虚拟路由器冗余协议, 为了解决静态路由器出现的单点故障).

keepalived运行机制:
- 网络层

    使用ICMP探测节点状态
- 传输层

    利用tcp的端口连接和扫描来判断节点状态
- 应用层

    通常是自定义keepalived来处理.

keepalived配置: `/etc/Keepalived/Keepalived.conf`, 可分为:
- 全局配置(global configuration)
- vrrpd配置
- lvs配置

## vrrp
将两台及以上的物理路由器虚拟成一个虚拟路由器, 它有一个唯一标识vrid(与一组ip关联), 同一时刻仅有一台物理路由器对外提供服务.

# lb
## lvs(base kernel)
参考:
- [LVS原理篇：LVS简介、结构、四种模式、十种算法](https://blog.csdn.net/lcl_xiaowugui/article/details/81701949)

LVS的模型中有两个角色：
- 调度器:Director Server，又称为Dispatcher，Balancer

    调度器主要用于接受用户请求
- 真实主机:Real Server，简称为RS

    用于真正处理用户的请求

将所在角色的IP地址分为以下4种：
- Director Virtual IP:调度器用于与客户端通信的IP地址，简称为VIP
- Director IP:调度器用于与RealServer通信的IP地址，简称为DIP
- Real Server : 后端主机的用于与调度器通信的IP地址，简称为RIP
- CIP：Client IP，访问客户端的IP地址所

lvs的ip负载均衡是通过ipvs实现的, 具体机制分3种:
- NAT

    用户请求LVS到达director，director将请求的报文的目的IP改为RIP和将报文的目标端口也改为realserver的相应端口，最后将报文发送到realserver上，realserver将数据返回给director，director再将响应报文的源地址和源端口换成vip和相应端口, 最后把数据发送给用户

    优点:
    - 支持端口映射
    - RS可以使用任意操作系统
    - 节省公有IP地址
    
        RIP和DIP都应该使用同一网段私有地址，而且RS的网关要指向DIP. 使用nat另外一个好处就是后端的主机相对比较安全.

    缺点：
    - 效率低: 请求和响应报文都要经过Director转发;极高负载时，Director可能成为系统瓶颈
- FULLNAT


    FULLNAT模式也不需要DIP和RIP在同一网段
    FULLNAT和NAT相比的话：会保证RS的回包一定可到达LVS
    FULLNAT需要更新源IP，所以性能正常比NAT模式下降10%
- TUN

    基于ip隧道技术.
    用户请求LVS到达director，director通过IP-TUN技术将请求报文的包封装到一个新的IP包里面，目的IP为VIP(不变)，然后director将报文发送到realserver，realserver基于IP-TUN解析出来包的目的为VIP，检测网卡是否绑定了VIP，绑定了就处理这个包，如果在同一个网段，将请求直接返回给用户，否则通过网关返回给用户；如果没有绑定VIP就直接丢掉这个包

    优点：

    - RIP,VIP,DIP都应该使用公网地址，且RS网关不指向DIP
    - 只接受进站请求，请求报文经由Director调度，但是响应报文不需经由Director, 解决了LVS-NAT时的问题，减少负载

    缺点：
    - 不指向Director所以不支持端口映射
    - RS的OS必须支持隧道功能
    - 隧道技术会额外花费性能，增大开销, 同时运维麻烦, 一般不用.

- DR

    当Director接收到请求之后，通过调度方法选举出RealServer, 将目标地址的MAC地址改为RealServer的MAC地址. RealServer接受到转发而来的请求，发现目标地址是VIP, RealServer配置在lo接口上, 处理请求之后则使用lo接口上的VIP响应CIP.

    优点:
    - RIP可以使用私有地址，也可以使用公网地址, 只要求DIP和RIP的地址在同一个网段内.
    - 请求报文经由Director调度，但是响应报文不经由Director. 性能最高最好, 因此**最常用**.
    - RS可以使用大多数OS

    缺点：
    - 不支持端口映射
    - 不能跨局域网

四种模式的比较:
- 是否需要VIP和realserver在同一网段
    
    DR模式因为只修改包的MAC地址，需要通过ARP广播找到realserver，所以VIP和realserver必须在同一个网段，也就是说DR模式需要先确认这个IP是否只能挂在这个LVS下面；其他模式因为都会修改目的地址为realserver的IP地址，所以不需要在同一个网段内
- 是否需要在realserver上绑定VIP
    
    realserver在收到包之后会判断目的地址是否是自己的IP
    DR模式的目的地址没有修改，还是VIP，所以需要在realserver上绑定VIP
    IP TUN模式值是对包重新包装了一层，realserver解析后的包的IP仍然是VIP，所以也需要在realserver上绑定VIP
- 四种模式的性能比较
    
    DR模式、IP TUN模式都是在包进入的时候经过LVS，在包返回的时候直接返回给client；所以二者的性能比NAT高
    但TUN模式更加复杂，所以性能不如DR
    FULLNAT模式不仅更换目的IP还更换了源IP，所以性能比NAT下降10%
    性能比较：DR>TUN>NAT>FULLNAT

LVS的八种调度方法
- 静态方法:仅依据算法本身进行轮询调度

    - RR:Round Robin,轮调 : 均等地对待每一台服务器，不管服务器上的实际连接数和系统负载
    - WRR:Weighted RR，加权论调 : 加权，手动让能者多劳
    - SH:SourceIP Hash : 与目标地址散列调度算法类似，但它是根据源地址散列算法进行静态分配固定的服务器资源. 来自同一个IP地址的请求都将调度到同一个RealServer
    - DH:Destination Hash : 根据目标 IP 地址通过散列函数将目标 IP 与服务器建立映射关系，出现服务器不可用或负载过高的情况下，发往该目标 IP 的请求会固定发给该服务器. 不管IP，请求特定的东西，都定义到同一个RS上

- 动态方法:根据算法及RS的当前负载状态进行调度

    - LC:least connections(最小链接数) : 链接最少，也就是Overhead最小就调度给谁; 假如都一样，就根据配置的RS自上而下调度
    - WLC:Weighted Least Connection (加权最小连接数) : 这个是LVS的默认算法. 调度器可以自动问询真实服务器的负载情况，并动态调整权值
    - SED:Shortest Expection Delay(最小期望延迟) : WLC算法的改进. 不考虑非活动链接，谁的权重大，优先选择权重大的服务器来接收请求，但权重大的机器会比较忙
    - NQ:Never Queue : SED算法的改进
    - LBLC:Locality-Based Least-Connection,基于局部的的LC算法

        请求数据包的目标 IP 地址的一种调度算法，该算法先根据请求的目标 IP 地址寻找最近的该目标 IP 地址所有使用的服务器，如果这台服务器依然可用，并且有能力处理该请求，调度器会尽量选择相同的服务器，否则会继续选择其它可行的服务器
    - lblcr, 复杂的基于局部性最少的连接算法 : 记录的不是要给目标 IP 与一台服务器之间的连接记录，它会维护一个目标 IP 到一组服务器之间的映射关系，防止单点服务器负载过高

## haproxy
基于应用实现的软负载均衡, 基于tcp和http.

## FAQ
### Cannot use default route w/o netmask
该错误是`systemctl restart pacemaker`时`journalctl -f`截获的.

原因: 定义`ocf:heartbeat:IPaddr`时未指定`cidr_netmask`
