# route

route命令可修改路由表, ICMP重定向报文也可修改路由表.

## example
```bash
# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
192.168.0.0      0.0.0.0         255.255.192.0   U     0      0        0 eth0
0.0.0.0         192.168.0.1      0.0.0.0         UG    100    0        0 eth0
```

Flags各项的含义：
- U 该路由可用
- G 该路由是一个网关，如果没有该标志，则是直接路由
- H 该路由是一个主机，如果没有该标志，则是一个网络
- D 该路由是由重定向报文创建的
- M 该路由被ICMP重定向报文修改过