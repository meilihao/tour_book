# traceroute
显示网络数据包传输到指定主机的路径信息, 追踪数据传输路由状态.

mtr: 结合了traceroute和ping的功能，用于实时监控网络路径、延迟和丢包情况

## 格式

    traceroute [选项] [远程目标] [数据包大小]

## 选项
- -i <网络接口> : 使用指定网络接口发送数据包
- -I : 使用ICMP. traceroute默认使用tcp或udp
- -n : 直接使用ip而不使用hostname
- -v : 命令执行的详细过程
- -w <timeout> : 设置等待远程目标回应的时间
- -x : 开启或关闭对数据包的正确性检验
- -s <source ip> : 设置本机主机发送数据包的ip地址
- -g <gateway> : 设置来源的路由网关, 最多8个

## example
```
# traceroute -i eth0 -s 192.168.0.60.251 -w 10 www.baidu.com 100 # 它会对每个路由节点做ICMP的回应时间测试.
```

## FAQ
### 响应中出现`*`
1. 该节点没有在指定时间内响应
1. 该节点不接收/不返回ICMP包

