# ip route

## 描述

管理路由表的工具

## 例
```bash
# ip route [list/show] # 显示系统路由
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
# arping -I ens160 192.168.16.38 # 反查mac, by `apt install arping`
```

网关是路由出口的位置, 不一定和本机网段相同.

# ip addr
## example
```bash
ip -s link list # 显示更加详细的设备信息, 信息from /sys/class/net/${interface}/statistics
ip link show                     # 显示网络接口信息
ip link set eth0 up             # 开启网卡
ip link set eth0 down            # 关闭网卡
ip link set eth0 promisc on      # 开启网卡的混合模式
ip link set eth0 promisc offi    # 关闭网卡的混个模式
ip link set eth0 txqueuelen 1200 # 设置网卡队列长度
ip link set eth0 mtu 1400        # 设置网卡最大传输单元
ip addr show     # 显示网卡IP信息
ip addr add 192.168.0.1/24 dev eth0 # 设置eth0网卡IP地址192.168.0.1, 需要ip link set eth0 down/up重启网卡
ip addr del 192.168.0.1/24 dev eth0 # 删除eth0网卡IP地址
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

# ip netns
参考:
- [Linux ip netns 命令](https://www.cnblogs.com/sparkdev/p/9253409.html)

用来管理 network namespace

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

将多个网卡端口绑定为一个，可以提升网络的性能. 在linux系统上有两种技术可以实现:Linux 网络组和bond.

网络组(Teaming, RHEL7开始使用): 是将多个网卡聚合在一起方法，从而实现冗错和提高吞吐量,网络组不同于旧版中bonding 技术，能提供**更好的性能和扩展性**，网络组由内核驱动和teamd 守护进程实现.

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