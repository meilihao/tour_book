# iproute2
- [新的网络管理工具 ip 替代 ifconfig 零压力](http://www.linuxstory.org/replacing-ifconfig-with-ip/)
- [ip命令以及与net-tools的映射](https://linux.cn/article-3144-1.html)
- [放弃 ifconfig，拥抱 ip 命令](https://linux.cn/article-13089-1.html)
- [iproute2](https://github.com/mozhuli/SDN-Learning-notes/blob/master/linux/iproute2.md)

## 网卡工作模式 
网卡有以下几种工作模式，通常网卡会配置广播和多播模式：
- 广播模式（Broad Cast Model）:它的物理地址地址是 0Xffffff 的帧为广播帧，工作在广播模式的网卡接收广播帧。它将会接收所有目的地址为广播地址的数据包，一般所有的网卡都会设置为这个模式
- 多播传送（MultiCast Model）：多播传送地址作为目的物理地址的帧可以被组内的其它主机同时接收，而组外主机却接收不到。但是，如果将网卡设置为多播传送模式，它可以接收所有的多播传送帧，而不论它是不是组内成员。当数据包的目的地址为多播地址，而且网卡地址是属于那个多播地址所代表的多播组时，网卡将接纳此数据包，即使一个网卡并不是一个多播组的成员，程序也可以将网卡设置为多播模式而接收那些多播的数据包。
- 直接模式（Direct Model）:工作在直接模式下的网卡只接收目地址是自己 Mac 地址的帧。只有当数据包的目的地址为网卡自己的地址时，网卡才接收它。
- 混杂模式（Promiscuous Model）:工作在混杂模式下的网卡**接收所有的流过网卡的帧，抓包程序就是在这种模式下运行的**。网卡的缺省工作模式包含广播模式和直接模式，即它只接收广播帧和发给自己的帧。如果采用混杂模式，网卡将接受同一网络内所有所发送的数据包，这样就可以到达对于网络信息监视捕获的目的。它将接收所有经过的数据包，这个特性是编写网络监听程序的关键


## 组件
- ip : 用于管理路由表和网络接口

    - link: 配置网络设备
    - addr: 管理某个网络设备与协议(ipv4/ipv6)相关的地址
    - addrlabel: ipv6的地址标签, 主要用于RFC3484中描述的ipv6地址的选择
    - route: 管理路由
    - rule: 管理策略路由
    - neigh: 管理neighbor/arp表
    - ip tunel: 隧道配置
    - ip maddr: 多播地址管理
    - mroute: 多播路由管理
    - monitor : 状态监控. 比如监控ip地址和路由的状态
    - xfrm: 设置xfrm. xfrm是一个ip框架, 可以转换数据包的格式
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
ref:
- [iproute2路由配置（ip rule、ip route、traceroute）](https://www.cnblogs.com/liugp/p/16395089.html)

    net-tools和iproute2的大致对比
- [Linux的策略路由](https://linuxgeeks.github.io/2017/03/17/170119-Linux%E7%9A%84%E7%AD%96%E7%95%A5%E8%B7%AF%E7%94%B1/)

## 描述

管理路由表的工具

Linux 最多可以支持 255 张路由表，其中有 4 张表是内置的:
- 表 255 本地路由表（Local table） 本地接口地址，广播地址，已及 NAT 地址都放在这个表. 该路由表**由系统自动维护，管理员不能直接修改**
- 表 254 主路由表（Main table） 如果没有指明路由所属的表，所有的路由都默认都放在这个表里，一般来说，旧的路由工具（如 route）所添加的路由都会加到这个表. 一般是普通的路由
- 表 253 默认路由表 （Default table） 一般来说默认的路由都放在这张表，但是如果特别指明放的也可以是所有的网关路由
- 表 0 保留

## 例
基本语法: ip route add [prefix] via [gateway] dev [device] [additional options]

常用参数及解释
- add     增加路由
- del     删除路由
- -net    设置到某个网段的路由
- -host   设置到某台主机的路由
- gw      出口网关 IP地址, 推荐使用**via**
- prefix:             指定要添加的目标网络或主机。格式为 IP地址/子网掩码长度，例如 192.168.1.0/24 或 0.0.0.0/0（默认路由）。
- via [gateway]:      指定路由数据包的下一跳网关 IP 地址。例如 via 192.168.1.1。这是数据包应通过的中间路由器。
- dev [device]:       指定要通过的网络设备接口。例如 dev eth0。这用于指示数据包应发送到的物理或虚拟网络接口。
- src [source-address]:   指定用于该路由的数据包的源地址。例如 src 192.168.1.100。这在**多源地址**环境中非常有用。
- metric [value]:     指定路由的优先级（度量值）。度量值越小，优先级越高。用于在多条路由可用时确定使用哪条路由。
- table [table_id]:   指定要添加路由的路由表。默认情况下，路由添加到主路由表（table 254）。
- scope [scope]:      指定路由的范围。常用范围包括：
    global：         全球范围，适用于整个互联网。
    link：               链路范围，适用于直连网络。
    host：               主机范围，仅适用于特定主机。
- proto [protocol]:   指定路由协议。常用协议包括：
    static： 静态路由。
    kernel： 由内核生成的路由。
    boot：       在系统启动时添加的路由。
    dhcp：       由 DHCP 获取的路由。
    onlink:     强制将目标视为直接连接，即使没有相关 ARP 条目。这在某些需要绕过常规 ARP 检查的情况下有用。
    nexthop:    指定多个下一跳以实现负载均衡或冗余
- weight：指定每个下一跳的权重，用于流量分配

    比如`ip route add 192.168.1.0/24 nexthop via 192.168.1.1 dev eth0 weight 1 nexthop via 192.168.1.2 dev eth1 weight 2`

```bash
# cat /etc/iproute2/rt_tables
#
# reserved values
#
255     local
254     main
253     default
0       unspec
#
# local
#
#1      inr.ruhep
# ip route [list/show] [dev xxx] [table yyy] # ip route = ip route show table main, 显示系统路由, 或使用`route -n`
default via 192.168.0.1 dev bond0 
192.168.0.0/24 dev bond0  proto kernel  scope link  src 192.168.0.141 metric 100 # 如果进程没有bind一个源地址，将会使用src域里面的源地址作为数据包的源地址进行发送; 但是如果进程提前bind了，命中了这个条目，但仍然会使用进程bind的源地址作为数据包的源地址. 因此这里的src只是一个建议的作用. metric 为路由指定所需跃点数的整数值（范围是 1 ~ 9999），它用来在路由表里的多个路由中选择与转发包中的目标地址最为匹配的路由。所选的路由具有最少的跃点数。跃点数能够反映跃点的数量、路径的速度、路径可靠性、路径吞吐量以及管理属性.
# ip route ls/show table local/255 # id/name from /etc/iproute2/rt_tables
# ip route list 192.168.182.0/24 # 查看指定网段的路由
# ip route show table all # 显示所有路由表的路由
# ip -4 addr show dev eth0 # 获取eth0上的ipv4
# ip route get from <本机ip> <目标ip> # 检查数据包的路由决策信息, 即确认访问路径.
# ip route get 192.168.0.141 # 获取指定目的IP的路由信息
192.168.0.141 dev eth0  src 192.168.0.121 
    cache
# ip route get 8.8.8.8 from 172.31.13.130 iif veth-1
# ip route add default via 192.168.1.254   # 设置系统默认路由
# ip route add default via  192.168.0.254  dev eth0        # 设置默认网关为192.168.0.254
# ip route add 192.168.4.0/24  via  192.168.0.254 dev eth0 # 设置192.168.4.0网段的网关为192.168.0.254,数据走eth0接口
# ip route del 192.168.4.0/24   # 删除192.168.4.0网段的网关
# ip route del default          # 删除默认路由
# ip route delete 192.168.1.0/24 dev eth0 # 删除路由
# ip route add 192.168.1.10/32 dev eth0 table 100 # 主机路由
# ip route show match 192.168.88.2 [table xxx] # 匹配包含或等于指定网段的路由
# ip route show match 1.2.3.1 table 110
# ip route show root 1.2.0.0/16 [table 110] # 查看在**指定的网段**内的路由条目 
# ip route flush # 清空默认表路由
# ip route flush table tab1 # 清空指定表路由
# ip route flush 192.168.1.0/24 # 清理所有192.168.1.0/24相关的所有路由
# ip route flush cache  #   清空路由表项缓存，下次通信时内核会查main表（即命令route输出的表）以确定路由
# ip neigh # 查看显示内核的ARP表(即同局域网的ip-mac映射, 本机不会缓存自己ip的arp映射), 与`nmap -sP 192.168.0.0/24 `即可查到某个ip的mac
192.168.88.137 dev eno1 lladdr 8c:d0:b2:c3:94:62 STALE # STALE, 下一次发送前需要重新确认mac
192.168.88.2 dev eno1 lladdr 00:1a:8c:58:45:18 REACHABLE 
# ip neigh add 192.168.1.100 lladdr 00:0c:29:c0:5a:ef dev eth0 # 添加arp映射
# ip neigh flush dev wlp3s0 # 清除arp缓存
# ip neigh del 192.168.1.100 dev eth0 # 删除arp映射
# ip neigh show 192.168.0.167 # 查看对应ip的mac, 前提是内核的ARP表有该记录, 没有则先ping一下
# arp 192.168.0.167 # 查看ip对应的mac, 但arp已淘汰
# arping -I ens160 192.168.16.38 # 反查mac, 需同**局域网**(比如同网段ip), by `apt install arping`, ens160必须是up
# ip link set dev ens33 multicast on # 启用多播
# ip maddr # 显示多播地址
# route -n # 靠前的优先. route已淘汰, 但理解`Kernel IP routing table`直观
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
default         gateway         0.0.0.0         UG    0      0        0 eth0
10.0.0.10       10.139.128.1    255.255.255.255 UGH   0      0        0 eth0
10.139.128.0    0.0.0.0         255.255.224.0   U     0      0        0 eth0
# arp -an # 查看数据走向
```

route说明:
- Destination	目标网络或目标主机. Destination 为 default（0.0.0.0）时，表示这个是默认网关，所有数据都发到这个网关
- Gateway	网关地址，**0.0.0.0 表示当前记录对应的 Destination 跟本机在同一个网段，通信时不需要经过网关**

    **Destination的0.0.0.0和Gateway的0.0.0.0不是同一个意思, 不要将Gateway的0.0.0.0理解为下一跳是gateway的意思**
- Genmask	Destination 字段的网络掩码，Destination 是主机时需要设置为 255.255.255.255，是默认路由时会设置为 0.0.0.0
- Flags	标记

    含义:
    - U 路由是活动的
    - H 目标是个主机
    - G 需要经过网关
    - R 恢复动态路由产生的表项
    - D 由路由的后台程序动态地安装
    - M 由路由的后台程序修改
    - ! 拒绝路由
- Metric	路由距离，到达指定网络所需的中转数，是大型局域网和广域网设置所必需的 （不在Linux内核中使用）
- Ref	路由项引用次数 （不在Linux内核中使用）
- Use	此路由项被路由软件查找的次数
- Iface	网卡名字，例如 eth0

网关是路由出口的位置, 不一定和本机网段相同.

路由说明:
```bash
# ip route
default via 192.168.88.1 dev enp2s0 proto static metric 100 # 1. via, 网关.
169.254.0.0/16 dev virbr0 scope link metric 1000 linkdown 
192.168.88.0/24 dev enp2s0 proto kernel scope link src 192.168.88.236 metric 100 # 2
192.168.122.0/24 dev virbr0 proto kernel scope link src 192.168.122.1 linkdown 
# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         192.168.88.1    0.0.0.0         UG    100    0        0 enp2s0 # 1
169.254.0.0     0.0.0.0         255.255.0.0     U     1000   0        0 virbr0
192.168.88.0    0.0.0.0         255.255.255.0   U     100    0        0 enp2s0 # 2
192.168.122.0   0.0.0.0         255.255.255.0   U     0      0        0 virbr0
```

`# 1`表示发往任何地方的数据包, 都是经过enp2s0, 并由网关192.168.88.1发出
`# 2`表示发往192.168.88.0/24的数据包, 都由enp2s0发出, 其src表示enp2s ip是192.168.88.236, metric 100是路由距离即到达指定网络需要的中转数

ip route get注意事项:
1. ip route get 会检查是否开启 IP 转发功能 ，如果返回 No route to host ，需要先开启 IP 转发功能(`echo 1 > /proc/sys/net/ipv4/ip_forward`)
1. 从 iptables/netfilter 架构可知, 路由决策是后于 DNAT(nat/PREROUTING) 和先于 SNAT(nat/POSTROUTING) 的， 因此 from 后面的 IP 须是原始 IP (SNAT 时) ，或是转换后的 IP (DNAT 时)
1. 如果出现错误 RTNETLINK answers: Invalid cross-device link ，则需要禁用 rp_filter

# ip addr
## example
```bash
ip addr show [dev ens33]   # 显示网卡IP信息
ip addr add 192.168.0.1/24 dev eth0 # 设置eth0网卡IP地址192.168.0.1, 需要ip link set eth0 down/up重启网卡
ip addr del 192.168.0.1/24 dev eth0 # 删除eth0网卡IP地址
ip link delete cilium_vxlan # 删除网络接口
ip [-j -details] link show
```

# ip rule
策略rule主要包含三个信息，即`rule的优先级，条件，路由表`,  其中rule的优先级数字越小表示优先级越高，然后是满足什么条件下由指定的路由表来进行路由. 在linux系统启动时，内核会为路由策略数据库配置三条缺省的规则，即rule 0，rule 32766， rule 32767（数字是rule的优先级）. 当有一个数据包需要路由时，系统会从优先级最高的规则开始检查，如果没有匹配的规则，就会继续检查下一条规则，直到找到一个匹配的规则. 如果没有任何自定义规则匹配，最后会使用优先级为 32767 的默认规则.

ip rule list的记录又称为路由规则. 只不过路由规则在路由表的基础上增加了优先级的概念. 优先级可以从具体路由表条目前的数字得出.

rule 和 route 的关系: 规则指向路由表，多个规则可以引用一个路由表，而且某些路由表可以没有策略指向它.

参数:
- from [ADDRESS[/MASK]]：指定源地址或源地址范围。只匹配来自这些地址的数据包。可以是单个 IP 地址或 CIDR 范围
- to [ADDRESS[/MASK]]：  指定目的地址或目的地址范围。只匹配发送到这些地址的数据包。可以是单个 IP 地址或 CIDR 范围
- iif [NAME]：        指定入接口。只匹配从该接口接收到的数据包。接口名例如 eth0
- oif [NAME]：        指定出接口。只匹配将从该接口发送出去的数据包。接口名例如 eth0
- priority [PREFERENCE]：指定规则的优先级。数值越小优先级越高。规则按优先级顺序进行匹配。规则按从小到大的顺序进行匹配      
- fwmark [MARK]：         指定防火墙标记。只匹配带有该标记的数据包。由 iptables 或 nftables 设置的数据包标记
- uidrange [UID-START-UID-END]：匹配数据包的用户 ID 范围。例如， 0-65535
- sport [PORT[/MASK]]：  指定源端口或源端口范围。只匹配来自这些端口的数据包。可以是单个端口或端口范围
- dport [PORT[/MASK]]：  指定目标端口或目标端口范围。只匹配发送到这些端口的数据包。可以是单个端口或端口范围
- tos [TOS]：             指定服务类型 (Type of Service)。只匹配带有此 TOS 值的数据包。TOS 值用来指定数据包的优先级
- ipproto [PROTOCOL]：   指定 IP 协议。只匹配使用该协议的数据包。例如， tcp， udp， icmp 等

action:
- table TABLE_ID: 指定路由表. 数据包匹配此规则后，将在指定的路由表中查找路由。默认路由表 ID 为 main (254)， local (255)， default (253) 等. 此处的table也可以替换为lookup，效果是一样的
- nat: 透明网关
- blackhole/reject: 丢弃匹配的数据包
- unreachable: 丢弃匹配的数据包，并发送 NET UNREACHABLE 的 ICMP 信息
- prohibit: 丢弃匹配的数据包，并发送 COMM.ADM.PROHIITED 的 ICMP 信息

```bash
# ip rule [list] # 查看系统有哪些策略路由
0:	from all lookup local # 0即优先级, 数字越小，代表优先级别越高
32766:	from all lookup main
32767:	from all lookup default
# echo '1024    tab1' >> /etc/iproute2/rt_tables # 添加路由表
# ip route list table main # 查看main路由表
# ip rule add [from 0/0] table 1 pref 32800 # 规则匹配的对象是所有的数据包，动作是选用路由表1的路由，这条规则的优先级是32800
# ip rule add from 192.168.182.0/24 table tab1 prio 220 # 根据来源端IP来决定数据包参考哪个路由表发送出去
# ip rule add from 192.168.182.247 [nat 192.168.182.130] table tab1 prio 320 # 在table 1添加rule且优先级是10
# ip rule add to 192.168.183.1 table tab1  # 根据目的端IP来决定数据包参考哪个路由表发送出去
# ip rule add to 192.168.183.0/24 table tab1
# ip rule add dev eth0 table 140 # 根据网卡设备决定路由表
# ip rule del from 192.168.182.0/24 table tab1 prio 220 # 删除策略路由
# ip rule del from 192.168.182.247 [nat 192.168.182.130] table tab1 prio 320
# ip rule del from 192.168.10.10  # 根据明细条目删除
# ip rule del prio 32765          # 根据优先级删除
# ip rule del table wt            # 根据表名称来删除
# ip rule add fwmark 1 prio 10 tab 10 # 自定义基于标签的规则和对应的路由表. `iptables -t mangle -A PREROUTING -p tcp -m multiport --dports 80 -s 192.168.24.0/24 -j MARK --set-mark 1`可以添加标签
# ip rule flush # 清空所有规则
```

example:
```bash
echo "192 net_192 " >> /etc/iproute2/rt_tables

#清空net_192路由表
ip route flush table net_192
# 添加一个路由规则到 net_192 表，这条规则是 net_192 这个路由表中数据包默认使用源 IP 172.31.192.201 通过 ens4f0 走网关 172.31.192.254
ip route add default via 172.31.192.254 dev ens4f0 src 172.31.192.201 table net_192
#来自 172.31.192.201 的数据包，使用 net_192 路由表的路由规则
ip rule add from 172.31.192.201 table net_192
```

# ip link
```bash
ip -s link list # 显示更加详细的设备信息, 信息from /sys/class/net/${interface}/statistics
ip link show                     # 显示网络接口信息
ip link set eth0 up             # 开启网卡, 前提是物理链路必须正常, 否则没效果
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
### 丢包
- [Linux内核常见的网络丢包场景分析](https://mp.weixin.qq.com/s/vdW0L7nEdfrxSJ_9VGviaA)
- [云网络丢包故障定位，看这一篇就够了](https://mp.weixin.qq.com/s?__biz=MzUyODY4NTA2Mw==&mid=2247551473&idx=1&sn=5d4bec2468585d719c82f53243c58344&source=41#wechat_redirect)

### [指定broadcast address](https://library.netapp.com/ecmdocs/ECMP1155586/html/GUID-0B34BAE2-D6F2-4E2D-A84F-D417EAA737E2.html)
通常无需手动设置广播地址, 它是根据ip地址和netmask自动确定. 实际ipv4可以指定广播地址(`ifconfig e3a broadcast 192.0.2.25`), ipv6不支持广播地址.

### 能ping通, 但对端不能接受数据包(包括SYN)
本机的ARP表中的对端mac可能错误

### ifconfig ip显示不完整/缺ip
用`ip addr`, `ifconfig`已淘汰.

### 多网卡同IP和同网卡多IP技术
参考:
- [多网卡同IP和同网卡多IP技术](https://www.jianshu.com/p/c3278e44ee9d)

#### 多网卡同IP技术
参考:
- [team与bonding的比较](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/networking_guide/sec-comparison_of_network_teaming_to_bonding)
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

bond mode 1/4区别:
1. mode 1:

    主备模式, 只有一个网卡在使用中
    
    优点就是很安全, 两块网卡同时坏的概率很低
    缺点则是利用率低下，只有50%的利用率

    应用场景一般是服务器的管理口, 管理口一般没有太高的网络需求, 稳定第一
1. mode 4:

    链路聚合模式, 相当于两块小网卡合并一起当作一个大网卡用, 类似1+1=2
    前置条件1: 交换机需要支持IEEE802.3ad(链路聚合标准), 并且在交换机上进行相应配置
    前置条件2: ethtool支持获取每个slave的速率和双工设定

    应用场景一般是业务网了, 需要的大的带宽的情况比较适合

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
- [Routing Tables](http://linux-ip.net/html/routing-tables.html)
- [Linux 策略路由](https://www.infvie.com/ops-notes/linux-policy-routing.html)

    多运营商ip

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
解决方法: 先将网卡**up**并给其配上ip再配置gateway即可.

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

### bond0使用的eth0和eth1 mac启动时小概率对调
ref:
- [小斗CentOS7.x网卡名称错乱、及网卡启动失败](https://www.houzhibo.com/archives/684)

1. 通过ip link(与`/sys/class/net/eth0/address`一致)获取eth0的mac
2. 通过`/sys/class/net`的eth0软连接路径或`udevadm info /sys/class/net/eth0 |grep ID_PATH`(需去除前缀`pci-`)获取id_path
3. 通过`dmesg |grep <id_path>`获取mac, 对比步骤1获取的mac, 发现步骤1获取的mac是dmesg里eth1的

同时eth0和eth1的udev规则在/etc/udev/rules.d/70-persistent-net.rules中, 是配置正确的. 且`dmesg|grep udev`没有发现报错信息.

正常环境dmesg时序:
1. kernel枚举device
1. 网卡rename
1. 创建bond

检查了问题环境的dmesg, 没有发现异常, 创建bond时也使用了正确的eth0

对调原因: 未知

> ID_PATH 提供了设备的总线拓扑路径，用于唯一标识设备在系统中的位置, 这个路径反映了设备是如何连接到系统的，涉及总线、端口、设备号等信息.

观察到在vmware上vm(非bond网卡)也出现了对调.

### 一键接管后, 原机ip ping不通
env:
- 原机: 16.142
- 接管机: 16.124, mac=52:54:00:7F:EC:01
- 一键接管server: 16.159
- ping机: 0.234

原机上系统日志报: 16.142和接管机上的mac冲突, 导致原机网关route消失. 具体日志为`NetworkManager[799]: ... conflict for address 192.168.16.142 detected with host 52:54:00:7F:EC:01 on interface 'ens3'`

推测: 一键接管后, 接管机启动时的ip还是16.142, 它启动过程中才会由接管agent去修改ip为16.124
改进方案:
- 原机网卡配置信息中添加mac地址进行绑定, 已测试
- 接管前接管机的系统盘上的ip信息要已被替换, 未测试

### 获取active_interface
`ip route get 8.8.8.8 | awk 'NR==1 {print $5}'`

### mac克隆
通过nmtui即可, 需要重启

### centos7 `systemctl restart netowrk`报`TNETLINK answers: File exists`
network 和 NetworkManager 服务有冲突, 使用其中一个即可, 推荐NetworkManager.

### 不添加路由前能ping通, 添加路由后ping不通
env:
- 192.168.16.165, gw: 192.168.16.1

现象: ping 192.168.8.2成功, 添加route `192.168.8.0/24 dev eth0 proto static metric 100`后ping不成功

原因: 参考route输出的Gateway说明. 没有该route时, 匹配到默认路由, 走了gateway, 因此能ping通; 有该规则时, 匹配到该规则, 因为`Gateway=0.0.0.0`, 表示当前记录对应的 Destination 跟本机在同一个网段，通信时不需要经过网关, 因此ping不通.