# ip route

## 描述

管理路由表的工具

## 例

    ip route show # 显示系统路由
    ip route add default via 192.168.1.254   # 设置系统默认路由
    ip route list                   # 查看路由信息
    ip route add 192.168.4.0/24  via  192.168.0.254 dev eth0 # 设置192.168.4.0网段的网关为192.168.0.254,数据走eth0接口
    ip route add default via  192.168.0.254  dev eth0        # 设置默认网关为192.168.0.254
    ip route del 192.168.4.0/24   # 删除192.168.4.0网段的网关
    ip route del default          # 删除默认路由
    ip route delete 192.168.1.0/24 dev eth0 # 删除路由
    # ip neigh # 查看显示内核的ARP表(ip-mac映射, 本机不会缓存自己ip的arp映射), 与`nmap -sP 192.168.0.0/24 `即可查到某个ip的mac
    # ip neigh add 192.168.1.100 lladdr 00:0c:29:c0:5a:ef dev eth0 # 添加arp映射
    # ip neigh flush dev wlp3s0 # 清除arp缓存
    # ip neigh del 192.168.1.100 dev eth0 # 删除arp映射
    # ip neigh show 192.168.0.167 # 查看对应ip的mac, 前提是内核的ARP表有该记录, 没有则先ping一下
    # arp 192.168.0.167 # 查看ip对应的mac, 但arp已淘汰

# ip addr
## example
```bash
ip -s link list # 显示更加详细的设备信息
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

## FAQ
### 能ping通, 但对端不能接受数据包(包括SYN)
本机的ARP表中的对端mac可能错误

### ifconfig bond ip显示不完整
用`ip addr`, `ifconfig`已淘汰.

### 多网卡同IP和同网卡多IP技术
参考:
- [多网卡同IP和同网卡多IP技术](https://www.jianshu.com/p/c3278e44ee9d)

#### 多网卡同IP技术
将多个网卡端口绑定为一个，可以提升网络的性能. 在linux系统上有两种技术可以实现:Linux 网络组和bond.

网络组(Teaming, RHEL7开始使用): 是将多个网卡聚合在一起方法，从而实现冗错和提高吞吐量,网络组不同于旧版中bonding 技术，能提供**更好的性能和扩展性**，网络组由内核驱动和teamd 守护进程实现.

确定内核是否支持 bonding:
```sh
# cat /boot/config-4.15.0-30deepin-generic |grep -i bonding
CONFIG_BONDING=m
```

#### 同网卡多IP技术
有两种实现：
- 早期的 ip alias
- 现在的secondary ip

ifconfig显示的格式为`eth0:N`(即单独的网络接口),`ip addr`则是网络接口属性里的一条记录.

## 配置网络的工具
nmtui 通过字符界面来配置网络.

> nmtui配置后需`vim /etc/sysconfig/network-scripts/ifcfg-eno16777736`配置`ONBOOT=yes`来支持重启仍激活网卡

> 配置网络需激活: `systemctl restart network`