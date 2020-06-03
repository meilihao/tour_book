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

## FAQ
### Cannot use default route w/o netmask
该错误是`systemctl restart pacemaker`时`journalctl -f`截获的.

原因: 定义`ocf:heartbeat:IPaddr`时未指定`cidr_netmask`
