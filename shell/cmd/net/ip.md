# iproute2
- [新的网络管理工具 ip 替代 ifconfig 零压力](http://www.linuxstory.org/replacing-ifconfig-with-ip/)
- [ip命令以及与net-tools的映射](https://linux.cn/article-3144-1.html)
- [放弃 ifconfig，拥抱 ip 命令](https://linux.cn/article-13089-1.html)
- [iproute2](https://github.com/mozhuli/SDN-Learning-notes/blob/master/linux/iproute2.md)


## 组件
- ip : 用于管理路由表和网络接口
- tc : 用于流量控制管理
- ss : 用于转储套接字统计信息
- lnstat : 用于转储linux网络统计信息
- bridge : 用于管理网桥地址和设备
- nstat : 类似于netstat, 但比它提供更多的信息

    ```bash
    nstat -a
    nstat --json
    ```

    ```bash
    $ strace -e open nstat 2>&1 > /dev/null|grep /proc
    open("/proc/uptime", O_RDONLY)          = 4
    open("/proc/net/netstat", O_RDONLY)     = 4
    open("/proc/net/snmp6", O_RDONLY)       = 4
    open("/proc/net/snmp", O_RDONLY)        = 4

    $ strace -e open netstat -s 2>&1 > /dev/null|grep /proc
    open("/proc/net/snmp", O_RDONLY)        = 3
    open("/proc/net/netstat", O_RDONLY)     = 3
    ```

    参考:
    - [Linux network metrics: why you should use nstat instead of netstat](https://loicpefferkorn.net/2016/03/linux-network-metrics-why-you-should-use-nstat-instead-of-netstat/)
    - [Linux network statistics reference](https://loicpefferkorn.net/2018/09/linux-network-statistics-reference/)

# ip route

## 描述

管理路由表的工具

## 例
```bash
# ip route [list/show] # 显示系统路由, 或使用`route -n`
default via 192.168.0.1 dev bond0 
192.168.0.0/24 dev bond0  proto kernel  scope link  src 192.168.0.141 # 如果进程没有bind一个源地址，将会使用src域里面的源地址作为数据包的源地址进行发送; 但是如果进程提前bind了，命中了这个条目，但仍然会使用进程bind的源地址作为数据包的源地址. 因此这里的src只是一个建议的作用
# ip -4 addr show dev eth0 # 获取eth0上的ipv4
# ip route get 192.168.0.141 # 获取指定目的IP的路由信息
192.168.0.141 dev eth0  src 192.168.0.121 
    cache
# ip route add default via 192.168.1.254   # 设置系统默认路由
# ip route add default via  192.168.0.254  dev eth0        # 设置默认网关为192.168.0.254
# ip route add 192.168.4.0/24  via  192.168.0.254 dev eth0 # 设置192.168.4.0网段的网关为192.168.0.254,数据走eth0接口
# ip route del 192.168.4.0/24   # 删除192.168.4.0网段的网关
# ip route del default          # 删除默认路由
# ip route delete 192.168.1.0/24 dev eth0 # 删除路由
# ip route flush 192.168.1.0/24 # 清理所有192.168.1.0/24相关的所有路由
# ip route flush cache  #   清空路由表项缓存，下次通信时内核会查main表（即命令route输出的表）以确定路由
# ip neigh # 查看显示内核的ARP表(ip-mac映射, 本机不会缓存自己ip的arp映射), 与`nmap -sP 192.168.0.0/24 `即可查到某个ip的mac
# ip neigh add 192.168.1.100 lladdr 00:0c:29:c0:5a:ef dev eth0 # 添加arp映射
# ip neigh flush dev wlp3s0 # 清除arp缓存
# ip neigh del 192.168.1.100 dev eth0 # 删除arp映射
# ip neigh show 192.168.0.167 # 查看对应ip的mac, 前提是内核的ARP表有该记录, 没有则先ping一下
# arp 192.168.0.167 # 查看ip对应的mac, 但arp已淘汰
# arping -I ens160 192.168.16.38 # 反查mac, 需同局域网, by `apt install arping`, ens160必须是up
# ip link set dev ens33 multicast on # 启用多播
# ip maddr # 显示多播地址
# route -n # 靠前的优先
# arp -an # 查看数据走向
```

网关是路由出口的位置, 不一定和本机网段相同.

路由说明:
```bash
# ip route
default via 192.168.88.1 dev enp2s0 proto static metric 100 # 1
169.254.0.0/16 dev virbr0 scope link metric 1000 linkdown 
192.168.88.0/24 dev enp2s0 proto kernel scope link src 192.168.88.236 metric 100 # 1
192.168.122.0/24 dev virbr0 proto kernel scope link src 192.168.122.1 linkdown 
# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface # 1
0.0.0.0         192.168.88.1    0.0.0.0         UG    100    0        0 enp2s0
169.254.0.0     0.0.0.0         255.255.0.0     U     1000   0        0 virbr0
192.168.88.0    0.0.0.0         255.255.255.0   U     100    0        0 enp2s0 # 2
192.168.122.0   0.0.0.0         255.255.255.0   U     0      0        0 virbr0
```

`# 1`表示发往任何地方的数据包, 都是经过enp2s0, 并由网关192.168.88.1发出
`# 2`表示发往192.168.88.0/24的数据包, 都由enp2s0发出, 其src表示enp2s ip是192.168.88.236, metric 100是路由距离即到达指定网络需要的中转数

# ip addr
## example
```bash
ip addr show     # 显示网卡IP信息
ip addr add 192.168.0.1/24 dev eth0 # 设置eth0网卡IP地址192.168.0.1, 需要ip link set eth0 down/up重启网卡
ip addr del 192.168.0.1/24 dev eth0 # 删除eth0网卡IP地址
ip link delete cilium_vxlan # 删除网络接口
ip [-j -details] link show
```

# ip rule
ip rule list的这三张路由表，又称为路由规则. 只不过路由规则在路由表的基础上增加了优先级的概念. 优先级可以从具体路由表条目前的数字得出. 数字越低，优先级越高.

```bash
# ip rule [list] # 查看系统有哪些路由表
0:	from all lookup local # 0即优先级
32766:	from all lookup main
32767:	from all lookup default
# ip rule add from 192.168.1.0/24 table 1 prio 10  在table 1添加rule且优先级是10
```

# ip link
```bash
ip -s link list # 显示更加详细的设备信息, 信息from /sys/class/net/${interface}/statistics
ip link show                     # 显示网络接口信息
ip link set eth0 up             # 开启网卡
ip link set eth0 down            # 关闭网卡
ip link set eth0 promisc on      # 开启网卡的混合模式
ip link set eth0 promisc offi    # 关闭网卡的混个模式
ip link set eth0 txqueuelen 1200 # 设置网卡队列长度
ip link set eth0 mtu 1400        # 设置网卡最大传输单元
# ip link delete cilium_net@cilium_host # cilium_net@cilium_host： cilium_host可能是cilium_net的secondary ip
Cannot find device "cilium_net@cilium_host"
# ip link delete cilium_net
```

# ip netns
参考:
- [Linux ip netns 命令](https://www.cnblogs.com/sparkdev/p/9253409.html)

用来管理 network namespace

# vxlan
ref:
- [使用 VXLAN 为虚拟机创建虚拟第 2 层域](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/assembly_using-a-vxlan-to-create-a-virtual-layer-2-domain-for-vms_configuring-and-managing-networking)

## FAQ
### 能ping通, 但对端不能接受数据包(包括SYN)
本机的ARP表中的对端mac可能错误

### ifconfig ip显示不完整/缺ip
用`ip addr`, `ifconfig`已淘汰.

### 多网卡同IP和同网卡多IP技术
参考:
- [多网卡同IP和同网卡多IP技术](https://www.jianshu.com/p/c3278e44ee9d)

#### 多网卡同IP技术
参考:
- [Bonding vs Team](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/configuring_and_managing_networking/configuring-network-teaming_configuring-and-managing-networking)
- [linux bond设备删除,删除修改bond](https://blog.csdn.net/weixin_33976326/article/details/116748742)

    ```bash
    # echo -bond0 > /sys/class/net/bonding_masters # `-`表示删除
    ```

将多个网卡端口绑定为一个，可以提升网络的性能. 在linux系统上有两种技术可以实现:Linux 网络组和bond.

网络组(Teaming, RHEL7开始使用): 是将多个网卡聚合在一起方法，从而实现冗错和提高吞吐量,网络组不同于旧版中bonding 技术，能提供**更好的性能和扩展性**: 支持ipv6, 至多8张网卡, 负载均衡, hash加密. 网络组由内核驱动和teamd 守护进程实现.

> 负载均衡模式和轮询模式最大的区别就在于，轮询一定是进程两个交替分配，而负载均衡模式会根据目前的空闲网卡来进行分配进程，不一定是两个交替，在大数据进行网络传输等时，效果会更好.

确定内核是否支持 bonding:
```sh
# cat /boot/config-4.15.0-30deepin-generic |grep -i bonding
CONFIG_BONDING=m
```

team安装:
```bash
apt install libteam-utils
dnf install teamd NetworkManager-team
```

team与bond的mod对应关系:
1. roundrobin 【mode 0】轮转策略 (balance-rr)
1. activebackup【mode 1】活动-备份（主备）策略
1. loadbalance【mode 2】限定流量
1. broadcast【mode 3】广播策略
1. lacp (implements the 802.3ad Link Aggregation ControlProtocol)【mode 4】

team:
```bash
# --- 创建网络组team接口
# nmcli connection add type team con-name team0 ifname team0 config
 {"runner":{"name":"activebackup"}} ipv4.method manual ipv4.addresses 192.168.32.100/24 connection.autoconnect yes
# --- 创建port接口
# nmcli connection add con-name team0-ens33 type team-slave ifname ens33 master team0 # 把 ens33网卡创建为 team0网络组的子接口
# nmcli connection add con-name team0-ens37 type team-slave ifname ens37 master team0 # 把 ens37 网卡创建为 team0网络组的子接口
# --- 启用team0下创建的两个子接口
# nmcli connection up team0-ens33
# nmcli connection up team0-ens37
# --- 其他
# nmcli connection down team0 # 删除网络组team
# teamdctl team0 state # 查看当前网络组team主网卡和热备网卡哪个在工作
# nmcli connection delete team0-ens33 / nmcli connection delete team0-ens37
# nmcli connection show # 查看网络组, 删除后就看不到了
```

手动配置team:
```bash
# vim /etc/sysconfig/network-scripts/team
DEVICE="team"
DEVICETYPE="Team"
ONBOOT="yes"
BOOTPROTO=static
NETMASK=255.255.255.0
IPADDR=ip
GATEWAY=网关
TEAM_CONFIG='{"runner": {"name": "activebackup"}}'
# vim /etc/sysconfig/network-scripts/ifcfg-em1 # 编辑文件ifcfg- em1
DEVICE="em1"
DEVICETYPE="TeamPort"
ONBOOT="yes"
TEAM_MASTER="team"
# vim /etc/sysconfig/network-scripts/ifcfg-em2 # 编辑文件ifcfg- em2
DEVICE="em2"
DEVICETYPE="TeamPort"
ONBOOT="yes"
TEAM_MASTER="team"
```

#### 同网卡多IP技术
有两种实现：
- 早期的 ip alias : ip alias 是由 ifconfig 命令来创建和维护的，ifconfig显示的格式为`eth0:N`(即单独的网络接口). alias IP就是在网卡设备上绑定的第二个及以上的IP
- 现在的secondary ip(**推荐**) : ip addr add 创建的辅助IP，不能通过ifconfig查看，但是通过ifconfig创建的别名IP却可以在ip addr show 命令查看

    > ip addr show时secondary ip的记录里会出现关键字`secondary`


特性:
1. 在一个网络接口上可配置多个primary地址和多个secondary地址
1. 对一个特定的网络掩码（比如网络掩码为/24），只能有一个Primary地址
1. 当删除一个primary地址时, 相关的secondary地址也会被删除. 但通过配置`net.ipv4.conf.<eth0>.promote_secondaries=1`后, 当前primary地址被删除时, secondary地址会提升为primary地址
1. 当主机为本地生成的流量选择源IP地址时，只考虑Primary地址

![secondary ip的kernel描述](/misc/img/shell/20170125120034149.gif)

每个节点代表的 IP 地址标识一个网段，这个节点的 IP 就是这个网段的 Primary 地址，它下面所带的 IP 就是这个网段的 Secondary 地址，也就是说一个网卡可以带有各个节点所带链表长度之和个 IP 地址，而且这些 IP 不是线形的，而是上述的吊链结构.

## 配置网络的工具
nmtui 通过字符界面来配置网络.

> nmtui配置后需`vim /etc/sysconfig/network-scripts/ifcfg-eno16777736`配置`ONBOOT=yes`来支持重启仍激活网卡

> 配置网络需激活: `systemctl restart network`

## 查看网卡细节: ethtool
## 启用网卡
```bash
ifconfig <interface> up # 推荐
ifup <interface> # 部分网卡不支持, 因为该interface未在/etc/network/interfaces文件中明确定义.
```

## 双网卡同网段IP实现
双网卡同网段IP, 默认仅能使用一个.

> linux网关也是默认只能使用一个.

解决方法, 用工具 iproute2 把两个网卡分到两个不同路由表:
```conf
echo "210    local100" >> /etc/iproute2/rt_tables
echo "220    local200" >> /etc/iproute2/rt_tables

ip route add 192.168.1.0/24 dev wlo0 src 192.168.1.11 table local100
ip route add 192.168.1.0/24 dev eno1 src 192.168.1.22 table local200
ip route add default dev wlo0 table local100
ip route add default dev eno1 table local200

ip rule add from 192.168.1.11 table local100
ip rule add from 192.168.1.22 table local200

ip route flush cache
```

example2:
ref:
- [Linux 多网关应用场景配置](https://typefo.com/linux/linux-multiple-gateway.html)

Linux 多网关应用场景, 比如机房服务器有 3 块网卡, eth0 为内网IP, eth1 为电信公网IP, eth2 为联通公网IP, 一般情况下服务器只能配置一个默认网关, 外网客户端只能通过其中一个公网IP访问服务器, 通过配置 Linux 原路返回路由功能, 来实现客户端从哪个网卡进来就从哪个网卡出去. 也就是电信用户访问服务器的电信公网 IP 然后从电信网卡原路返回, 联通用户访问联通公网 IP 然后从联通网卡返回, 服务器本身就可以通过默认的内网网关访问外网.


配置 Linux 多网关需求:
- eth0 内网网卡，IP 地址 192.168.1.100， 网关 192.168.1.1
- eth1 电信网卡，IP 地址 114.216.29.65，网关 114.216.29.1
- eth2 联通网卡，IP 地址 60.30.128.15，网关 60.30.128.1
```bash
# --- 配置内网 IP 地址和默认网关
$ ip link set eth0 up
$ ip addr add 192.168.1.100/24 dev eth0
$ ip route add default via 192.168.1.1
# --- 配置电信和联通网卡 IP 地址

$ ip link set eth1 up
$ ip addr add 114.216.29.65/24 dev eth1
$ ip link set eth2 up
$ ip addr add 60.30.128.15/24 dev eth2
# --- 添加电信和联通两个路由表: 编辑 /etc/iproute2/rt_tables 配置文件, 添加两个编号 251 和 252 的路由表条目, tel 为电信 cnc 为联通.
$ vim /etc/iproute2/rt_tables
251 tel
252 cnc

# --- 配置电信和联通的原路路由

$ ip route flush table tel
$ ip route add default via 114.216.29.1 dev eth1 src 114.216.29.65 table tel
$ ip rule add from 114.216.29.65 table tel
$ ip route flush table cnc
$ ip route add default via 60.30.128.1 dev eth1 src 60.30.128.15 table cnc
$ ip rule add from 60.30.128.15 table cnc
# --- 以上就基本配置好了电信和联通的多线原路返回路由
```

### `ip route`配置gateway时报`Nexthop has invalid gateway`
解决方法: 先将网卡up并给其配上ip再配置gateway即可.

### 使用udev重命名网卡(**不推荐修改**)
**通过添加kernel参数`net.ifnames=0 biosdevname=0`的方式可能不成功**. RHEL7开始不支持这种修改方式.

这里通过修改udev:
```bash
# vim /etc/udev/rules.d/70-persistent-net.rules # 也有叫`70-custom-ifnames.rules`
SUBSYSTEM=="net", ACTION=="add", DRIVERS=="?*", ATTR{address}=="00:0c:29:30:be:cd", ATTR{type}=="1", NAME="eth0" # need root. ATTR{type} 参数值 1 定义了接口类型为 Ethernet.
# mv /etc/sysconfig/network-scripts/ifcfg-ens192 /etc/sysconfig/network-scripts/ifcfg-eth0
# vim /etc/sysconfig/network-scripts/ifcfg-eth0 # 将内容里的ens192替换为eth0
# reboot
```

udev 设备管理器支持很多不同的命名方案. 默认情况下, udev 根据固件、拓扑和位置信息分配固定名称. 它有以下优点：
1. 设备名称完全可预测
1. 即使添加或删除了硬件，设备名称也会保持不变，因为不会进行重新枚举
1. 替换问题/旧硬件不会引发重命名

RHEL9使用`/usr/lib/systemd/network/99-default.link`中的 NamePolicy 设置设备名称. 如果手动配置 udev 规则来更改内核设备名称，则这些规则优先. [网络设备重命名启用规则](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/how-the-network-device-naming-works_consistent-network-interface-device-naming).

### Oracle linux 7.9网络配置
gateway: /etc/sysconfig/network
ip: /etc/sysconfig/network-scripts/ifcfg-xxx # 这里的GATEWAY会被忽略, 因此最好时删除它

### net device online要求
1. 已插入网线且链路ok
2. ifup xxx

### 对调网卡
通过`/etc/udev/rules.d/70-persistent-net.rules`对调名称即可

### 同网段ping/ssh均正常, 不同网段不能通信
即192.168.0.191和192.168.16.78不能通信

因为目标端的route没有配置, 在16.78配上默认网关192.168.16.2后通信正常

### 同网段直连不能通信
env:
- A : 192.168.2.100
- B : eth0:192.168.2.101,用于连外网;eth1:192.168.2.102,用于直连. 即双网卡在同一网段.
- 连接: A与B的eth1直连

问题其实是出在路由表. 系统里面有一个路由表，当设置IP的时候就会同时设置路由表, 当需要访问外面的时候, 系统会去路由表里面查询，当查询到第一个匹配的项目时就应用这个项, 也就是从这条路径走了.

如果系统有两个独立网卡, 并且这两个网卡的IP属于同一个子网, 数据包默认发到该网段的默认路由(通常是后启动的网卡即`ip addr`的输出顺序)所在的eth0上, 又因为eth0上没有192.168.2.102, 从而不通.

解决方法:
1. 直连换其他网段
1. 设置策略路由, 指定IP路由的走向: 给每个网卡分配单独的路由表, 并且通过 ip rule 来指定

    ```bash
    ip route add to 192.168.0.0/16 dev eth1 table 20/ip route add via 192.168.2.254 dev eth1 table 20 //路由表20 走eth1
    ip rule add from 192.168.2.102/32 table 20   //源IP为192.168.2.102 走路由表20
    ```
1. down B的eth0

    eth0下线时会删除旧路由, 从而使得默认路由在eth1上

### `ping 192.168.16.159`报`connect: Network is unreachable`
现象:
1. ipmi登入, 16.159能ping通网关(192.168.16.2)
1. 0.241和16.159不能相互ping通, 以前可以

原因: 修改网络时, 16.159的默认网关被删除.
解决方案: `ip route add default via 192.168.16.2 dev eth0`