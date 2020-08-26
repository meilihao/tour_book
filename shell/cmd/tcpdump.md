# tcpdump

## 描述

根据使用者的定义对网络上的数据包进行截获的包分析工具,常与wireshark组合使用(tcpdump捕获数据,再用wireshark分析)

信息来自驱动.

## 选项
- -d ：把编译过的数据包编码转换成可阅读的格式，并在stdout输出
- -dd ：把编译过的数据包编码转换成C语言的格式，并在stdout输出
- -ddd ：把编译过的数据包编码转换成十进制数字的格式，并在stdout输出
- -D : 输出tcpdum可用的接口列表
- -e : 显示链路层头部
- -f ：用数字显示网际网络地址
- -F : 使用文件作为输入的过滤表达式, 忽略命令行上的表达式
- -G rotate_second : 指定`-w`的轮转时间
- -i : 指定监听的网络接口
- -n : 不将主机地址转换成名称, 避免dns查询
- -nn : 不要转换协议和端口号到名称
- -O : 不使用数据包匹配代码优化器
- -q ：快速输出，仅列出少数的传输协议信息
- -r file : 从file(由`-w`创建)中读取信息. 如果来源是stdin则file是`-`
- -S : 输出绝对的TCP序列号,而不是相对的
- -T <数据包类型> ：强制将表达方式所指定的数据包转译成指定的数据包类型
- -v ：详细显示指令执行过程
- -vv ：更详细显示指令执行过程
- -x ：**用十六进制字码列出数据包内容**
- -X : **用hex格式输出数据包的报头和内容**
- -w : 将捕获的数据写入文件

## 字段
tcpdump停止捕获数据包时:
- packets captured : tcp接收并处理的数据包的数量
- packets received by filter : 取决于运行tcpdump的os及该os的相关配置
- packets dropped by kernel : 由于缓冲区空间不足而丢包的数量

## 例子
    # tcpdump -i lo host 127.0.0.1 and port 8210  -w out.cap # 捕获发送到127.0.0.1:8210的数据
    # tcpdump host 210.27.48.1 and / (210.27.48.2 or 210.27.48.3 /) # 截获主机210.27.48.1 和主机210.27.48.2 或210.27.48.3的通信
    # tcpdump ip host 210.27.48.1 and ! 210.27.48.2 # 获取主机210.27.48.1除了和主机210.27.48.2之外所有主机通信的ip包
    # tcpdump -i eth0 host hostname and dst port 80 -x #  目的端口是80 # 列出送到80端口的数据包
    # tcpdump -i eth1 src host 211.167.237.199 and dst port 1467 # 211.167.237.199通过ssh源端口连接到221.216.165.189的1467端口
    # tcpdump -e -i ens160 arp # 截获arp
    # tcpdump -vv -eqtnni ens160 arp # 截获收到的arp

注意:
1. 主机有多个网卡时(比如192.168.0.166, 192.168.0.168), 发生数据包的ip是不确定的.
1. 即使主机的端口没有被监听但也能收到SYN包.

## FAQ
### flag
即tcp的flag:
- [S]： SYN（开始连接）
- [.]: ACK/没有 Flag
- [P]: PUSH（推送数据）
- [F]: FIN （结束连接）
- [R]: RST（重置连接）
- [U] : (URG)
- [W] : (ECN CWR)
- [E] : (ECN-Echo)

### arp解读
```
52:54:00:fc:28:03 > ff:ff:ff:ff:ff:ff, ARP, length 60: Ethernet (len 6), IPv4 (len 4), Request who-has 192.168.0.67 tell 192.168.0.167, length 46 # 192.168.0.167发出arp req因此192.168.0.167对应52:54:00:fc:28:03
...
00:50:56:84:83:bd > 52:54:00:fc:28:03, ARP, length 60: Ethernet (len 6), IPv4 (len 4), Request who-has 192.168.0.167 (52:54:00:fc:28:03) tell 192.168.0.197, length 46 # 同上192.168.0.197对应00:50:56:84:83:bd 
52:54:00:fc:28:03 > 00:50:56:84:83:bd, ARP, length 60: Ethernet (len 6), IPv4 (len 4), Reply 192.168.0.167 is-at 52:54:00:fc:28:03, length 46 # 192.168.0.167 响应我是 52:54:00:fc:28:03
```

### ARP 广播拿到错误 MAC 地址
参考:
-  [ARP 广播拿到错误 MAC 地址](http://xy.am/2015/04/19/arp/)
- [arp_filter](https://lwn.net/Articles/45386/#arp_filter)

当有一个 ARP 广播请求一个 MAC 地址时，本机会查看自己是否拥有这个 IP，只要本机任何端口拥有该IP，就会随机返回一个网口的 MAC 地址, 即机器不管你的请求通过哪个网卡进来，反正请求能进到这台机器，就行了.

处理方法: 将 /proc/sys/net/ipv4/conf/all/arp_filter 设为 1 即可

> 也有可能是ip冲突, 解决方法: 关掉相关的电脑, 看还能否收到错误的arp, 如果还能收到则再用linux/windows远程去连接(发送错误mac的ip)一下看看.

