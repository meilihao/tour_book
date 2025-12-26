# iproute2
- [新的网络管理工具 ip 替代 ifconfig 零压力](http://www.linuxstory.org/replacing-ifconfig-with-ip/)
- [ip命令以及与net-tools的映射](https://linux.cn/article-3144-1.html)
- [放弃 ifconfig，拥抱 ip 命令](https://linux.cn/article-13089-1.html)
- [iproute2](https://github.com/mozhuli/SDN-Learning-notes/blob/master/linux/iproute2.md)
- [Red Hat Enterprise Linux 7 网络指南](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/networking_guide/)
- [iproute and Routing Tables](https://datahacker.blog/industry/technology-menu/networking/routes-and-rules/iproute-and-routing-tables)
- [**Linux 路由实现原理**](https://blog.csdn.net/CoderTnT/article/details/124856179)

## 网卡工作模式 
网卡有以下几种工作模式，通常网卡会配置广播和多播模式：
- 广播模式（Broad Cast Model）:它的物理地址地址是 0Xffffff 的帧为广播帧，工作在广播模式的网卡接收广播帧。它将会接收所有目的地址为广播地址的数据包，一般所有的网卡都会设置为这个模式
- 多播传送（MultiCast Model）：多播传送地址作为目的物理地址的帧可以被组内的其它主机同时接收，而组外主机却接收不到。但是，如果将网卡设置为多播传送模式，它可以接收所有的多播传送帧，而不论它是不是组内成员。当数据包的目的地址为多播地址，而且网卡地址是属于那个多播地址所代表的多播组时，网卡将接纳此数据包，即使一个网卡并不是一个多播组的成员，程序也可以将网卡设置为多播模式而接收那些多播的数据包。
- 直接模式（Direct Model）:工作在直接模式下的网卡只接收目地址是自己 Mac 地址的帧。只有当数据包的目的地址为网卡自己的地址时，网卡才接收它。
- 混杂模式（Promiscuous Model）:工作在混杂模式下的网卡**接收所有的流过网卡的帧，抓包程序就是在这种模式下运行的**。网卡的缺省工作模式包含广播模式和直接模式，即它只接收广播帧和发给自己的帧。如果采用混杂模式，网卡将接受同一网络内所有所发送的数据包，这样就可以到达对于网络信息监视捕获的目的。它将接收所有经过的数据包，这个特性是编写网络监听程序的关键

## 路由流程
路由表用于决定数据包从哪个网口发出，其判断依据则主要是IP地址.

Linux可以配置很多很多策略，数据包将依次通过各个策略，一旦匹配某个策略则进一步应用策略对应的路由表，如果当前路由表无法匹配到路由则继续执行后续策略匹配.

路由的基本流程为: 收到数据包之后，解析出目的IP，判断是否是本机IP。如果是本机IP，则交由上层传输层处理。如果不是本机IP，则通过查找路由表找到合适的网络接口将IP数据包转发出去.

Linux上通过路由规则和路由表配合来实现路由流程, 处理逻辑如下:
1. 按路由规则优先级, 根据规则匹配条件找到需要匹配的路由表
2. 根据路由表中条目进行匹配的结果进行转发
3. 若路由表中没有匹配到满足的路由条目，则处理下一路由规则

ps: linux路由表的id仅标识唯一, 不表示优先级; 策略路由有优先级, 数值小优先.

在默认情况下，Linux 上的转发功能是关闭的，如果发现收到的网络包不属于自己就会将其丢弃.

但在某些场景下，例如对于容器网络来说，Linux 需要转发本机上其它网络命名空间中过来的数据包，需要手工开启转发:
```bash
# sysctl -w net.ipv4.ip_forward=1
# sysctl net.ipv4.conf.all.forwarding=1
```

开启后，Linux 就能像路由器一样对不属于本机（严格地说是本网络命名空间）的 IP 数据包进行路由转发了.

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

man文档:
- ip
- ip-route

# ip route
ref:
- [iproute2路由配置（ip rule、ip route、traceroute）](https://www.cnblogs.com/liugp/p/16395089.html)

    net-tools和iproute2的大致对比
- [Linux的策略路由](https://linuxgeeks.github.io/2017/03/17/170119-Linux%E7%9A%84%E7%AD%96%E7%95%A5%E8%B7%AF%E7%94%B1/)

## 描述

管理路由表的工具

Linux 最多可以支持 255 张路由表，其中有 4 张表是内置的:
- 表 255 本地路由表(local): 本地接口地址，广播地址，已及 NAT 地址都放在这个表. 该路由表**由kernel自动维护，管理员不能直接修改**
- 表 254 主路由表(main): 如果没有指明所属的路由表的所有路由都默认都放在这个表里，一般来说，旧的路由工具（如 route）所添加的路由都会加到这个表. 一般是普通的路由即所有非策略路由
- 表 253 默认路由表(default): 所有其他路由表都没有匹配到的情况下，根据该表中的条目进行处理
- 表 0 保留

> Netfilter 处理网络包的先后顺序：接收网络包，先 DNAT，然后查路由策略，查路由策略指定的路由表做路由，然后 SNAT，再发出网络包

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
192.168.0.0/24 dev bond0  proto kernel  scope link  src 192.168.0.141 metric 100 # 如果进程没有bind一个源地址，将会使用src域里面的源地址作为数据包的源地址进行发送; 但是如果进程提前bind了，命中了这个条目，但仍然会使用进程bind的源地址作为数据包的源地址. 因此这里的src只是一个建议的作用. metric 为路由指定所需跃点数的整数值（范围是 1 ~ 9999），它用来在路由表里的多个路由中选择与转发包中的目标地址最为匹配的路由。所选的路由具有最少的跃点数。跃点数能够反映跃点的数量、路径的速度、路径可靠性、路径吞吐量以及管理属性. metric最小，优先级最高.
# ip route add table 3 via 10.0.0.1 dev ethX # 添加路由表, ethx 是 10.0.0.1 所在的网卡, 3是路由表的编号
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
# ip route del default table 30000  # 删除路由表中的路由项
# ip route del 192.168.8.0/24 table 30000  # 删除路由表中的路由项
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
策略rule主要包含三个信息，即`rule的优先级(Priority)，条件(Selector)，路由表(Action)`,  其中rule的优先级数字越小表示优先级越高，然后是满足什么条件下**由指定的路由表来进行路由**. 在linux系统启动时，内核会为路由策略数据库配置三条缺省的规则，即rule 0，rule 32766， rule 32767（数字是rule的优先级）. 当有一个数据包需要路由时，系统会从优先级最高的规则开始检查，如果没有匹配的规则，就会继续检查下一条规则，直到找到一个匹配的规则. 如果没有任何自定义规则匹配，最后会使用优先级为 32767 的默认规则.

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

>  lookup [xxx] : 表示搜索xxx路由表，1-252之间的数字或名称

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
# ip rule del table <wt|30000>            # 根据表名称/id来删除
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
### 静态路由与`ip和网卡`的关系
静态路由 是手动配置的路由规则，用于指定数据包从源地址到目标地址的路径.

每条静态路由通常包括以下信息：
1. 目标网络：数据包要到达的目标网络地址（如 192.168.2.0/24）
1. 网关：数据包需要经过的下一个路由器的 IP 地址（如 192.168.1.1）
1. 网络接口：数据包从哪个网络接口发送出去（如 eth0）

静态路由是针对 IP 地址 的，而不是直接针对网卡（网络接口）。静态路由的作用是告诉操作系统如何将数据包发送到特定的目标网络或 IP 地址.

静态路由的工作原理:
1. 当操作系统收到一个数据包时，会检查数据包的目标 IP 地址。
1. 根据路由表（包括静态路由）决定数据包的下一跳（网关）和出口接口。
1. 如果目标 IP 地址匹配静态路由规则，数据包将通过指定的接口发送到指定的网关

三者关系:
1. 静态路由本质上是为数据包指明如何到达目标地址，它依赖于 **IP 地址**，而不是直接与网卡绑定
1. 网卡是数据包的出口: 静态路由需要指定一个网络接口来发送数据包

    在大多数情况下，静态路由需要指定网卡，因为网卡是数据包的出口。

    在以下情况下，可以不显式指定网卡：
    1. 点对点链路
    1. 默认路由
    
        如果主机只有一个网卡，操作系统会自动通过该网卡发送数据包，而不需要显式指定网卡
    1. 多网卡环境中的隐式选择

        如果主机有多个网卡，但**网关的 IP 地址只与其中一个网卡在同一子网中**，操作系统会自动选择该网卡作为出口
    1. VPN 或隧道接口

    操作系统会根据网关的 IP 地址和路由表自动选择合适的网卡作为出口
1. 网卡的 IP 地址和网关：**静态路由的网关必须与网卡的 IP 地址在同一子网中**

    原因:
    1. ARP 协议只能在同一子网内工作，用于解析网关的 MAC 地址
    2. 操作系统需要知道如何将数据包发送到网关，而网关必须在本地网络中

    ps: 在点对点链路等特殊情况下，网关可以与网卡的 IP 地址可以不在同一子网中. 因为点对点链路不需要 ARP 协议，数据包直接通过链路发送到对端设备

### 直接路由
直接路由是指目标网络与主机或路由器**直接相连**的情况. 操作系统会自动生成直接路由，无需手动配置.

静态路由与直接路由的区别:
特性	静态路由	直接路由
配置方式	手动配置	自动生成
适用网络	跨网络通信（目标网络不在本地子网）	直连网络（目标网络在本地子网）
是否需要网关	需要指定网关	无需网关
路由表生成	手动添加到路由表	自动添加到路由表
适用场景	小型网络或网络拓扑稳定的环境	目标网络与主机或路由器直接相连
动态调整	不支持动态调整	随网络接口状态自动调整

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
- [**NIC team 在 Red Hat Enterprise Linux 9 中已弃用**](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/configuring-network-teaming_configuring-and-managing-networking#proc_migrating-a-network-team-configuration-to-network-bond_configuring-network-teaming)

    - [If You Like Bonding, You Will Love Teaming](https://www.redhat.com/en/blog/if-you-bonding-you-will-love-teaming)

        原因: **行业普遍接受bond, 同时redhat希望专注于一个解决方案**

        NIC team将在rhel 10移除
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
ref:
- [多网卡配置同一网段IP情况解析](https://blog.csdn.net/whenloce/article/details/88095474)

多网卡分别配置同网段的不同ip, 为什么配置的部分ip会不通? 这通常是由于 路由冲突 或 ARP 问题 引起的:
1. 路由冲突

    当多个网卡配置了相同网段的 IP 地址时，操作系统会为每个网卡生成一条直接路由（direct route）.
    这些直接路由的目标网络相同，但出口接口不同，导致路由表出现冲突
    操作系统无法确定数据包应该从哪个网卡发送出去，从而导致部分 IP 地址无法通信
2. ARP 问题

    ARP（地址解析协议）用于将 IP 地址解析为 MAC 地址
    当多个网卡配置了相同网段的 IP 地址时，ARP 请求和响应可能会混淆. 因为其他主机ping 本机ip返回的是相同mac, 通过systemtap跟踪到问题在fib_validate_source(), 导致处于相同网段的几个IP之间，在进行arp reply的时候，因为直连路由指向同网段, 此时会使用优先级高(metric最小)的那条直连路由对应的mac来响应.

    例如，如果 eth0 和 eth1 都连接到同一交换机，交换机会将 ARP 请求广播到所有端口，导致 ARP 表项冲突.

解决方法, 用工具 iproute2 把两个网卡分到两个不同路由表(其实也可是是新建table+main, 本质还是分到两个不同路由表):
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

### ip route add default dev wlo0和ip route add default  gw 10.0.0.1 dev wlo0区别?
|命令	|含义	|使用场景|
|ip route add default dev wlo0	|设置默认路由，不指定网关，数据包将通过 wlo0 接口转发到本地网络或直接相连的设备|	适用于可以直接通过 wlo0 访问目标网络的情况（如直接连接的局域网）|
|ip route add default gw 10.0.0.1 dev wlo0	|设置默认路由，通过网关 10.0.0.1 转发数据包，数据包通过 wlo0 接口|	适用于需要通过网关（如路由器）转发数据包到其他网络或互联网的情况|

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

其他情况:
1. 关机下, 插拔了bond的网线

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

### nmtui配置ip后, 过段时间后会消失
手动配置ipv4后, `IPv4 CONFIGURATION`仍是Automatic导致与dhcp冲突, 改为Manual即可

### 配置静态路由
ref:
- [在 ifcfg 文件中配置静态路由](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/networking_guide/sec-configuring_static_routes_in_ifcfg_files#sec-Configuring_Static_Routes_in_ifcfg_files)

静态路由用于定义目标网络的路由路径. 如果尝试到达的任何网络没有被一个静态路由定义覆盖，则只能通过默认网关访问.

`route-<interface>` 文件的格式有两种:
- 传统格式([Network/Netmask Directives Format, 使用网络/网络掩码指令格式的静态路由](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/networking_guide/sec-configuring_static_routes_in_ifcfg_files#sec-Configuring_Static_Routes_in_ifcfg_files))：每行定义一个路由规则
- IP命令格式：使用 ip route 命令的语法

    ```bash
    <目标网络> via <网关> dev <接口>
    default via 192.168.1.1 dev eth0 # 默认路由通过网关 192.168.1.1，使用 eth0 接口. default=0.0.0.0/0
    192.168.2.0/24 via 192.168.1.1 dev eth0 # 发送到 192.168.2.0/24 网络的流量通过网关 192.168.1.1，使用 eth0 接口
    ```

路由的优先级：多个路由配置可以通过目标网络的匹配度来优先选择. 例如，192.168.1.0/24 会优先匹配比 0.0.0.0/0 更具体的路由

配置完成后，需要重启网络服务或接口以使配置生效

### /etc/sysconfig/network和/etc/sysconfig/network-scripts配置的关系
ref:
- [4.6. 配置默认网关](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/networking_guide/sec-configuring_the_default_gateway#sec-Configuring_the_Default_Gateway)
- [`man 5 nm-settings-ifcfg-rh`]()

对于新部署的系统（尤其是 RHEL 8 和 CentOS 8 及以后版本），建议使用 NetworkManager 或 systemd-networkd 来进行网络配置. 且在 CentOS Stream 9 和 RHEL 9 中，/etc/sysconfig/network-scripts/ 目录已经不再默认存在，网络配置被完全转移到 NetworkManager 或 systemd-networkd. 同时ifup和ifdown命令在现代的Linux发行版中已经逐渐被废弃或替换, 可使用`ip link`代替

ps:
1. ifup/ifdown 是 Debian/Ubuntu 等用来管理网络接口的启动和停止的命令. 它不仅仅是操作接口，还会根据 /etc/network/interfaces 或 /etc/netplan（在 Ubuntu 20.04 及以后的版本中）中的网络配置来管理接口, 比如触发一系列与网络配置相关的操作（如关闭 DHCP 客户端，清除路由等）.
1. ifconfig 是旧版的网络配置工具，已经被 ip 命令所取代. ifconfig down 只是 简单地禁用网络接口，它不会触发任何配置文件的读取，也不会执行复杂的网络配置操作, 主要是将网络接口从操作系统层面上禁用

在Linux系统中，/etc/sysconfig/network 和 /etc/sysconfig/network-scripts 是用于网络配置的两个重要目录，通常用于基于Red Hat的发行版（如CentOS、RHEL、Fedora等）。它们之间的关系如下：

1. /etc/sysconfig/network
作用：这个文件用于设置全局网络配置，适用于整个系统

常见配置项：
- NETWORKING=yes|no：是否启用网络功能
- HOSTNAME=<hostname>：设置系统的主机名
- GATEWAY=<gateway-ip>：设置默认网关（通常也可以在接口配置文件中设置）
- NISDOMAIN=<nis-domain>：设置NIS域名（如果使用NIS服务）

2. /etc/sysconfig/network-scripts
作用：这个目录包含每个网络接口的独立配置文件，用于配置具体的网络接口（如eth0、eth1等）

文件命名规则：ifcfg-<interface-name>，例如 ifcfg-eth0 是eth0接口的配置文件

常见配置项：
- DEVICE=<interface-name>：指定网络接口的名称（如eth0）。
- BOOTPROTO=dhcp|static|none：设置IP地址的获取方式（DHCP或静态IP）。
- IPADDR=<ip-address>：设置静态IP地址。
- NETMASK=<subnet-mask>：设置子网掩码。
- GATEWAY=<gateway-ip>：设置默认网关（如果未在全局配置中设置）。
- ONBOOT=yes|no：是否在系统启动时启用该接口。
- DNS1=<dns-server-ip> 和 DNS2=<dns-server-ip>：设置DNS服务器。

3. 两者的关系
全局与局部配置：
- /etc/sysconfig/network 是全局配置文件，适用于整个系统的网络设置
- /etc/sysconfig/network-scripts 是局部配置文件，针对每个网络接口进行独立配置

优先级：
- **如果某个配置项（如网关）同时在全局配置文件和接口配置文件中设置，接口配置文件中的值通常会覆盖全局配置**

依赖关系：
- 如果 /etc/sysconfig/network 中的 NETWORKING=no，则无论接口配置文件如何设置，网络功能都不会启用
- 接口配置文件中的 ONBOOT=yes 必须设置为 yes，才能在系统启动时启用该接口

### NetworkManager / systemd-networkd / netplan
NetworkManager 是一个功能强大的网络管理工具, 支持动态管理和复杂的网络配置, 旨在为桌面和移动环境提供易用的图形界面和功能强大的网络管理能力.

systemd-networkd 是 systemd 组件的一部分，主要用于管理系统的网络配置. 它适合于服务器、嵌入式设备和其他不需要图形界面的环境, 主要提供基础的网络配置功能，较为轻量级和稳定.

Netplan是Canonical(Ubuntu)开发的做为ubuntu Linux发行版上默认的网络配置命令行工具. Netplan 使用 YAML 描述文件来配置网络，然后通过这些描述为任何给定的底层呈现工具(主要就是systemd-networkd和networkmanager二种工具)生成必要的配置选项.

## FAQ
### 诡异/莫明其妙网络排查顺序
1. ip冲突, **大概率**

    表现:
    1. tcp连接总是莫名断开
    1. 均能ping通目标ip, 但有些电脑能连上该ip, 有些则不能连接(抓包发现是syn包没有响应)
    1. 抓包中出现很多`TCP Spurious Retransmission`, `TCP Dup ACK`, `TCP Retransmission`包

2. mac重复