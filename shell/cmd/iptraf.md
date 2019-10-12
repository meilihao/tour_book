# iptraf
iptraf是一个基于ncurses的IP局域网监控器，用来生成包括TCP信息、UDP计数、ICMP和OSPF信息、以太网负载信息、节点状态信息、IP校验和错误等等统计数据.

## 字段
- IP traffic monitor : 网络流量按tcp连接统计
- General interface statistics : ip流量按网络接口统计
- Detailed interface statistics : 网络流量按协议统计
- Statistical breakdowns : 网络流量按tcp/udp端口或数据包大小统计
- LAN station monitor : 网络流量按数据链路层地址统计

## 选项
- -i iface : 指定接口或所有接口(`-i all`)
- -g : 立即开始生成网络接口的概要状态信息
- -d iface : 在指定接口上立即启动详细的统计
- -s iface : 在指定网络接口上立即开始监视TCP和UDP网络流量信息
- -z iface : 在指定网络接口上显示包大小的计数
- -l iface : 在指定网络接口或所有接口上启动监视局域网工作信息
- -t timeout : 指定iptraf监视的时间, 以分钟为单位
- -B : 将stdout重新定向到`/dev/null`，关闭stdin，将程序作为后台进程运行
- -L logfile : 指定一个文件用于记录命令行的所有log，默认文件是地址：/var/log/iptraf
- -I interval : 指定记录log的时间间隔(单位是minute), 不包括IP traffic monitor
- -f : 清空所有计数器和所有锁. 仅用于从异常终止或系统崩溃中恢复回来.
