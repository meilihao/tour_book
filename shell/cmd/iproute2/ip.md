# ip route

## 描述

管理路由表的工具

## 例

    # ip route add 192.168.16.0/24 via 192.168.88.2 dev enp2s0 # 192.168.88.2为gw
    # ip neigh # 查看显示内核的ARP表(ip-mac映射, 本机不会缓存自己ip的arp映射), 与`nmap -sP 192.168.0.0/24 `即可查到某个ip的mac
    # ip neigh add 192.168.1.100 lladdr 00:0c:29:c0:5a:ef dev eth0 # 添加arp映射
    # ip neigh flush dev wlp3s0 # 清除arp缓存
    # ip neigh del 192.168.1.100 dev eth0 # 删除arp映射
    # ip neigh show 192.168.0.167 # 查看对应ip的mac, 前提是内核的ARP表有该记录, 没有则先ping一下
    # arp 192.168.0.167 # 查看ip对应的mac, 但arp已淘汰


## FAQ
### 能ping通, 但对端不能接受数据包(包括SYN)
本机的ARP表中的对端mac可能错误