# ss

## 描述

ss是iproute2包的一部分（控制TCP/IP网络和流量的工具）.iproute2的目标是替代先前用于配置网络接口、路由表和管理ARP表的标准Unix网络工具套装（通常称之为“net-tools”,其自2001年起便不再更新和维护）.ss工具用于导出套接字统计，它可以显示与netstat类似的信息，且可以显示更多的TCP和状态信息.因为它直接从内核空间获取信息，所以会更快.ss的操作和netstat很像，所以这让它可以很容易就可以取代netstat.

## 语法格式

```
ss [options] [ FILTER ]
```

## 选项

- -n : 不要尝试解析服务名
- -r : 尝试解析数字的地址/端口
- -a : 显示所有套接字
- -l : 显示监听套接字
- -p : 显示使用该套接字的进程
- -s : 打印统计数据
- -t : 只显示TCP套接字
- -u : 只显示UDP套接字
- -d : 只显示DCCP套接字
- -w : 只显示RAW套接字
- -x : 只显示Unix域套接字
- -f FAMILY : 显示FAMILY套接字的类型。目前支持下面这些族：unix、inet、inet6、link、netlink
- -A QUERY : 指定要列出的套接字列表，通过逗号分隔。可以识别下面的标识符：all、inet、tcp、udp、raw、unix、packet、netlink、unixdgram、unixstream、packetraw、packetdgram
- -o STATUS : 列出指定状态的套接字

其他:

- 可使用-4 标志来显示IPv4链接，-6标志来显示IPv6链接
- ss使用IP地址筛选
      ss src ADDRESS_PATTERN
src：表示来源
ADDRESS_PATTERN：表示地址规则
- ss使用端口筛选：
      ss dport OP PORT
OP:是运算符
PORT：表示端口
dport：表示过滤目标端口、相反的有sport

	OP运算符如下：
<= or le : 小于等于 >= or ge : 大于等于
== or eq : 等于
!= or ne : 不等于端口
< or lt : 小于这个端口 > or gt : 大于端口


## 例

    # ss -a # 查看这台服务器上所有的socket连接
    # ss -pl |grep 8000 #查询本机打开8000端口的进程信息
    # ss -o state established '( dport = :http or sport = :http )' #ss列出所有http连接中的连接,包含对外提供的80，以及访问外部的80,用于获取http并发连接数，监控中常用到