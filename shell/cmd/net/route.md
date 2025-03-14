# route

route命令可修改路由表, ICMP重定向报文也可修改路由表.

## example
```bash
# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
default         gateway         0.0.0.0         UG    0      0        0 eth0
10.0.0.10       10.139.128.1    255.255.255.255 UGH   0      0        0 eth0
10.139.128.0    0.0.0.0         255.255.224.0   U     0      0        0 eth0
```

route说明:
- Destination	目标网络或目标主机. Destination 为 default（0.0.0.0）时，表示这个是默认网关，所有数据都发到这个网关
- Gateway	网关地址，**0.0.0.0 表示当前记录对应的 Destination 跟本机在同一个网段，通信时不需要经过网关**

    **Destination的0.0.0.0和Gateway的0.0.0.0不是同一个意思, 不要将Gateway的0.0.0.0理解为下一跳是gateway的意思**
- Genmask	Destination 字段的网络掩码，Destination 是主机时需要设置为 255.255.255.255，是默认路由时会设置为 0.0.0.0
- Flags	标记

    含义:
    - U 路由是活动的
    - H 该路由是一个主机，如果没有该标志，则是一个网络
    - G 该路由是一个网关，如果没有该标志，则是直接路由
    - R 恢复动态路由产生的表项
    - D 由该路由是由重定向报文创建的
    - M 该路由被ICMP重定向报文修改过
    - ! 拒绝路由
- Metric	路由距离，到达指定网络所需的中转数，是大型局域网和广域网设置所必需的 （不在Linux内核中使用）
- Ref	路由项引用次数 （不在Linux内核中使用）
- Use	此路由项被路由软件查找的次数
- Iface	网卡名字，例如 eth0

ps: **route仅显示主路由表**