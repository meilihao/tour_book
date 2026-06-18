# list

## 高可用
- [Linux网卡bond的七种模式详解](https://blog.51cto.com/linuxnote/1680315)

## tools
- [dropwatch - 监听系统丢包信息工具](https://cloud.tencent.com/developer/article/1638140)
- [命令行 DNS 查询工具，支持 DNS-over-TLS 和 DNS-over-HTTPS](https://github.com/mr-karan/doggo)
- [NextTrace : 源轻量级可视化路由追踪工具](https://www.nxtrace.org/)
- [BGP 工具探索](https://linux.cn/article-13857-1.html)

## bgp
- [Syntropy的初创公司提出了一种去中心化自主路由协议（DARP），旨在取代BGP](https://linux.cn/article-13204-1.html)

## 网络I/O指标
ref:
- [网络I/O：丢包、重传、连接数，三个最容易搞混的指标](https://mp.weixin.qq.com/s/Zz07Ci9nWnEKQzyk2HAlRQ)

1. 丢包率（Packet Loss）: 网卡接收/发送队列溢出导致的丢包
    ifconfig输出的`RX/TX drops` / `cat /proc/net/dev`


    丢包位置:
    1. 网络栈: 内核网络栈最底层的丢包统计

        1. 网卡Ring Buffer → RX/TX drops
        
            网卡收到数据包，先放进Ring Buffer（接收队列）。如果Ring Buffer满了（比如中断处理不及时），新到的包就被丢掉。这个丢包会体现在ifconfig的RX drops里

            `ifconfig eth0 | grep -i drop`/`cat /proc/net/dev | grep eth0`
        1. 内核qdisc（排队规则）→ TX drops

            数据包从协议栈往下走，要经过qdisc（通常是pfifo_fast或mqprio）。如果qdisc队列满了，包被丢掉。这个丢包有些会体现在ifconfig TX drops里，但不一定——取决于网卡驱动的实现

            `tc -s qdisc show dev eth0`: 看qdisc的dropped计数; 如果backlog很高，说明qdisc在积压
        1. Socket Receive Buffer → 应用层丢包

            这是最容易被忽略的。数据包已经通过了网卡和内核协议栈，放进了Socket接收缓冲区（sk_rcvbuf）。但应用层读取得不够快，缓冲区满了，内核就会偷偷丢掉新到的包，也不会通知应用

            `cat /proc/net/snmp | grep -E "Tcp:|Ip:"` : 看TCP丢包统计（包括Socket Buffer溢出）
            1. TcpExtListenOverflows = 连接队列满
            1. TcpExtTCPBacklogDrop = 半连接队列丢包
            1. IpInReceives - IpInDelivers = 被丢掉的数据包
            
            `cat /proc/sys/net/core/rmem_default`: 默认接收缓冲区
            `cat /proc/sys/net/core/rmem_max`: 最大接收缓冲区
    1. 其他

        1. Socket Buffer溢出
        1. conntrack表满都不会体现在网卡drops里

    网卡`drops=0`≠ 没有丢包. Socket Buffer溢出导致的丢包，应用层感知到的是"连接慢"或"超时"，但网卡层面完全看不出来
1. 重传率（Retransmission Rate）: TCP报文丢失后触发重传的比例

    `cat /proc/net/snmp | grep -E "Tcp:|^Tcp"`:
    1. TcpOutSegs = 发出的TCP报文总数
    1. TcpRetransSegs = 重传的报文数
    1. TcpExtTCPTimeouts = 超时重传次数
    1. 重传率 = TcpRetransSegs / TcpOutSegs
    
    更精确的（包含各连接的重传） : `ss -i`, 看rto/rtt/cwnd/retrans 这几个字段

    重传率 = TcpRetransSegs / TcpOutSegs, 它是TCP层面的可靠性，不是"丢包"。TCP保证可靠传输的方式就是——丢了就重传

    |场景|网卡drops|TCP重传率|用户感知|
    |网络真正丢包|可能高|一定高|超时/慢|
    |网络正常，但对端慢|0|可能高（ACK延迟→触发重传）|慢|
    |网络正常，应用正常|0|< 1%|正常|

    重传的触发条件（不止丢包）：
    1. 超时重传：发送方在规定时间内没收到ACK，就重传。这个"规定时间"是RTO（Retransmission Timeout），动态计算，基于RTT
    2. 快速重传：收到3个重复ACK，不等超时直接重传（不需要等RTO）。这是最主要的重传类型
    3. SACK重传：Selective Acknowledgment，可以精确告诉发送方"哪段丢了"，避免重传已经收到的数据  

    重传率经验说明
    - < 1% : 正常，网络质量好
    - 1% ~ 5% : 需要关注，可能有网络抖动
    - > 5% : 异常，用户体验明显受损

    重传≠网卡看到丢包: 对端慢、ACK丢失也会触发重传
1. 连接数（Connections）: ESTABLISHED / TIME_WAIT / CLOSE_WAIT

    ss -s / netstat -an

    TIME_WAIT的危害：
    1. 占用本地端口（作为客户端连接后端时）—— 一台机器默认只有约28k个临时端口（/proc/sys/net/ipv4/ip_local_port_range）
    2. 占用内存（内核为每个TIME_WAIT保留一个小的控制块）

    CLOSE_WAIT：最容易出问题的状态, 被动关闭一方会进入CLOSE_WAIT——对端已经关闭连接（发了FIN），但本端应用层没有调用close()，连接就一直卡在CLOSE_WAIT.

    CLOSE_WAIT持续上涨 = 应用Bug——连接池没正确关闭、代码里忘了close、异常路径没处理。这个状态不会自动消失，只能重启应用或修复代码

    调优方式（不推荐tw_recycle，已废弃）
    ```
    # 推荐：允许复用TIME_WAIT的连接（仅对客户端有效） 
    # echo 1 > /proc/sys/net/ipv4/tcp_tw_reuse 
    # 不推荐：tw_recycle已在Linux 4.12移除，且在NAT环境下不安全 
    # echo 1 > /proc/sys/net/ipv4/tcp_tw_recycle
    # 增大本地端口范围 echo "10000 65535" > /proc/sys/net/ipv4/ip_local_port_range
    ```

    看各状态连接数:
    ```bash
    # -- 各状态连接数统计
    # ss -s # 或 netstat -an | awk '/^tcp/ {print $6}' | sort | uniq -c 
    # -- 看CLOSE_WAIT是否异常
    # netstat -an | grep CLOSE_WAIT | wc -l # 如果持续 > 100，需要排查应用
    ```

    K8s环境里，Pod的网络流量都要经过iptables + conntrack（连接跟踪表）。这个表有大小限制，满了之后新连接直接被丢弃，表现为"网络时好时坏"，但网卡层面完全看不出问题

    确认conntrack表满:
    ```
    # 看当前conntrack表使用量
    # cat /proc/sys/net/netfilter/nf_conntrack_count
    # 看conntrack表上限
    # cat /proc/sys/net/netfilter/nf_conntrack_max # 使用量 / 上限 > 80% 就要告警
    # 看是否有表满丢包
    # dmesg | grep "nf_conntrack: table full" # 或者 journalctl -k | grep conntrack
    ```

    调优:
    ```bash
    # -- 调大conntrack表上限（按需，比如调到65k）
    # echo 131072 > /proc/sys/net/netfilter/nf_conntrack_max
    # -- 减少conntrack表项存活时间（默认不活动432000秒=5天，太长了）
    # echo 86400 > /proc/sys/net/netfilter/nf_conntrack_tcp_timeout_established
    ```

    经验值：K8s节点上nf_conntrack_max建议设为内存MB数 × 16（比如16G内存 → 262144)

1. RTT（Round Trip Time）: 数据包往返时延，直接影响用户体验

    ss -i / cat /proc/net/tcp
1. 带宽使用率 : 网卡吞吐 / 网卡带宽上限

    sar -n DEV 1 / ifstat

网络问题排查
1. 确认用户感知

    - 是"完全不通"还是"很慢"？
    - 是所有用户还是部分用户？

1. 看重传率（比丢包率更重要）

    cat /proc/net/snmp | grep TcpRetransSegs
    重传率 > 5% → 网络质量有问题，用户会感知到慢

1. 看丢包（网卡层面）

    ifconfig看RX/TX drops
    drops在涨 → 网卡Ring Buffer或qdisc有问题，调大rx/tx ring size

1. 看连接数异常

    ss -s看各状态
    CLOSE_WAIT持续上涨 → 应用Bug，连接没关闭
    TIME_WAIT太高 → 客户端端口耗尽，开tcp_tw_reuse

1. K8s环境额外看conntrack

    cat /proc/sys/net/netfilter/nf_conntrack_count
    接近上限 → 调大nf_conntrack_max

1. 抓包确认（终极手段）

    tcpdump -i eth0 -w /tmp/capture.pcap
    用Wireshark分析重传、乱序、零窗口

关键经验：网络问题80%可以通过ss -i和/proc/net/snmp定位，不需要抓包

## FAQ
### 丢包排查
参考:
- [Linux 系统 UDP 丢包问题分析思路](https://cloud.tencent.com/developer/article/1638140)

    方法1: dropwatch
    方法2: perf

        ```bash
        # perf record -g -a -e skb:kfree_skb
        # perf script
        ```