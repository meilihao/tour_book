# sar

## 描述

sar（System Activity Reporter系统活动情况报告）是目前 Linux 上最为全面的系统性能分析工具之一，可以从多方面对系统的活动进行报告，包括：文件的读写情况、系统调用的使用情况、磁盘I/O、CPU效率、内存使用状况、进程活动及IPC有关的活动等.

> 采样间隔配置: `/etc/cron.d/sysstat`
> 可用gnuplot绘成图.

## install/安装

```shell
# ubuntu
sudo apt-get install sysstat
```

## 选项

- -b : 统计i/o和传输速率
- -B : 统计分页
- -d : 统计每个设备的活动
- -I : 统计中断
- -m : 统计电源管理
- -n { key[,...] | ALL}: 统计网络

	关键词可以是：
		DEV	 网卡
		EDEV	 网卡 (错误)
		NFS	 NFS 客户端
		NFSD	 NFS 服务器
		SOCK	 Sockets (套接字)	(v4)
		IP	IP  流	      (v4)
		EIP	 IP 流	   (v4) (错误)
		ICMP	 ICMP 流	(v4)
		EICMP	 ICMP 流	(v4) (错误)
		TCP	 TCP 流  (v4)
		ETCP	 TCP 流  (v4) (错误)
		UDP	 UDP 流  (v4)
		SOCK6	 Sockets (套接字)	(v6)
		IP6	 IP 流	   (v6)
		EIP6	 IP 流	   (v6) (错误)
		ICMP6	 ICMP 流 (v6)
		EICMP6 ICMP 流 (v6) (错误)
		UDP6	UDP 流   	(v6)
- -o file : 将命令结果以二进制格式存放在名为file的文件中
- -q : 查看运行队列中的进程数、系统上的进程大小、平均负载等
- -r : 查看内存使用状况
- -R : 内存统计
- -S : 统计swap使用率
- -u [ALL] : 显示CPU利用率
- -v : 报告inode状态, 文件状态和其他内核表
- -w : 统计进程创建和系统切换活动
- -W : 显示交换分区状态
- -y : TTY 设备状况

## 说明

cpu输出项说明：
- all 表示统计信息为所有 CPU 的平均值
- 其他项与iostat一致

分页:
- pgpgin/s : 每秒系统从磁盘置入分页的总量(KB)
- pgpgout/s : 每秒系统移出分页到磁盘的总量(KB)
- fault/s : 每秒系统产生的分页错误(major+minor)的数量.
- majflt/s : 系统每秒产生主要错误数量，需要从磁盘加载一个内存分页
- pgfree/s : 系统每秒放置在空闲列表的分页数量
- pgscank/s : 每秒kswapd守护进程扫描的分页数量
- pgscand/s: 每秒直接扫描的分页数量
- pgsteal/s: 每秒系统从缓存回收的分页数量
- %vmeff: pgsteal/pgscan 度量分页回收效率. 太低说明虚拟内存有问题, 如果在interval时间内没有分页被扫描, 则为0

io/传输速率:
- tps：每秒发送到物理设备的 I/O 请求总量
- rtps：每秒发送到物理设备的读请求总量
- wtps：每秒发送到物理设备的写请求总量
- bread/s：每秒从物理设备读入的数据量，单位为：块/s. 块相当于扇区, 即512B
- bwrtn/s：每秒向物理设备写入的数据量，单位为：块/s

块设备状况: 参考iostat
- rd_sec/s : iostat's rsec/s.
- wr_sec/s : iostat's wsec/s.

中断:
- Int : 中断号
- SUM：每秒的中断总数
- ALL：前16个中断的统计数据
- XALL：所有中断的统计数据

电源管理:
CPU
    MHz：瞬时时钟频率
FAN : 风扇速度
    rpm：每分钟转速
    drpm：当前转速与下限的差异
    DEVICE：传感器名称
FREQ : cpu时钟频率
    wghMHz:CPU时钟频率加权值，MHz，需要编译cpufreq-stats驱动
IN : 电压输入
    inV：电压输入伏特
    %in：相对输入，100%为输入到达上限，0为下限
    DEVICE：传感器名称
TEMP : 设备温度
    degC：摄氏度
    %temp：相对温度，100%意味着温度已达到上限
    DEVICE：传感器名称
USB : sar统计时当前插入系统所有usb的快照
    BUS：根集线器数量
    Idvendor：供应商ID
    Idprod：产品ID
    Maxpower：最大功耗mA
    Manufact：制造商名称
    Product：产品名称

网络: 参考*Linux性能优化大师p62*
key :
- DEV : 网络设备
	- IFACE : 网络接口名称
	- rxpck/s ：每秒接收包数
	- txpck/s ：每秒发送包数
	- rxkB/s  ：每秒接收数据量，KB
	- txkB/s  ：每秒发送数据量，KB
	- rxcmp/s ：每秒接收压缩数据包数量
	- txcmp/s ：每秒发送压缩数据包数量
	- rxmcst/s：每秒接收多播数据包数量
	- %ifutl : 网络接口的使用率，即半双工模式下为 (rxkB/s+txkB/s)/Bandwidth，而全双工模式下为 max(rxkB/s, txkB/s)/Bandwidth
- EDEV : 网络设备故障
	- IFACE : 网络接口名称
	- rxerr/s ：每秒接收坏数据包的数量
	- txerr/s ：当传输数据包时每秒发生错误的数量
	- coll/s  ：当传输数据包时每秒发生冲突的数量
	- rxdrop/s：因为linux缓冲区空间不足, 每秒丢弃接收数据包的数量
	- txdrop/s：因为linux缓冲区空间不足, 每秒丢弃传输数据包的数量
	- txcarr/s：每秒发送时发生载波错误的数量
	- rxfram/s：每秒接收时发生帧同步错误的数量
	- rxfifo/s：每秒接收时发生FIFO溢出错误的数量
	- txfifo/s: 每秒发送时发生FIFO溢出错误的数量
- NFS : NFS客户端活动
- NFSD : NFS服务端活动
- SOCK : IPv4 socket
	- totsck ：系统使用的socket总数
	- tcpsck ：当前使用中的tcp数量
	- udpsck ：当前使用中的udp数量
	- rawsck ：当前使用中的raw数量
	- ip-frag：当前队列中IP分片数量
	- tcp-tw ：处于TIME_WAIT的tcp socket数量
- IP : IPv4数据包
	- irec/s  ：每秒收到数据包数量，包括错误接收
	- fwddgm/s：收到的需要转发数据包数量，非最终目的IP
	- idel/s  ：收到数据包成功传递到IP协议的数量，包括ICMP
	- orq/s   ：本地IP协议栈(包括ICMP)到传输请求ip(ipOutRequests)的数据包总数
	- asmrq/s ：收到的需要重组的IP分片数量
	- asmok/s ：成功重组的的IP数据包数量
	- fragok/s：成功分片的IP数据包数量
	- fragcrt/s：生成的IP数据包分片数量
- EIP : IPv4 网络错误
	- ihdrerr/s：因为IP头错误，丢弃的输入数据包，包括checksum、版本、格式、TTL、options等
	- iadrerr/s: 因为目标地址非法，丢弃的输入数据包
	- iukwnpr/s: 成功接收但因为协议错误丢弃的数据包
	- idisc/s  : 成功接收IP数据报，但最终被丢弃的数量, 比如缓冲区空间不足
	- odisc/s  : 预计会成功发送，但最终被丢弃的数量, 比如冲区空间不足
	- onort/s  : 因为路由不正确丢弃的IP数据报数量
	- asmf/s   : 通过重组算法检测到的故障数量，不一定是丢弃的分片数量
	- fragf/s  : 因为无法分片丢弃的IP数据报数量
- TCP
	- active/s ：每秒连接从CLOSED转到SYN-SENT状态的次数
	- passive/s：每秒连接从LISTEN转到SYN-RCVD状态的次数
	- iseg/s   ：每秒接收到的分段数，包括错误的分段
	- oseg/s   ：每秒发送的分段数，不包含重发的字节数
- ETCP
	- atmptf/s ：每秒连接从SYN-SENT或SYN-RCVD直接转为CLOSED，加上从SYN-RCVD直接转为LISTEN状态的次数
	- estres/s ：每秒连接从ESTABLISHED或CLOSE-WAIT直接转为CLOSED状态的次数
	- retrans/s：每秒重传的分段数，即TCP分段包含了一个或多个之前的字节
	- isegerr/s：每秒接收的错误分段数, 比如错误的tcp校验和
	- orsts/s  ：每秒发送的RST分段数
...

队列长度和平均负载:
- runq-sz ：运行队列的长度（等待运行的进程数）
- plist-sz ：进程列表中进程（processes）和线程（threads）的数量
- ldavg-1 ：最后1分钟的系统平均负载（System load average）
- ldavg-5 ：过去5分钟的系统平均负载
- ldavg-15 ：过去15分钟的系统平均负载

内存:
- frmpg/s ：系统每秒释放内存分页数量，-表示分配的数量
- bufpg/s ：系统每秒使用额外内存作为缓冲区的数量，-表示系统使用较少的分页作为缓冲区
- campg/s ：系统每秒使用额外内存作为缓存的数量，-表示缓存中有较少的分页

内存使用:
- kbmemfree：可用空闲内存数量，KB
- kbmemused：已使用内存数量，KB
- %memused：使用率
- kbbuffers：内核用作缓冲区的内存数量，KB
- kbcached：内核用作缓存的内存数量，KB
- kbcommit：当前工作负载需要的内存数量，KB，是对需要多少RAM/swap的估计，以保证永远不会内存不足
- %commit：相对于总内存(ram+swap)，当前工作负载需要的内存百分比，当内核过量使用内存时会大于100%
- kbactive：活跃内存数量，KB，最近使用的内存通常不会回收
- kbinact：非活跃内存数量，KB，最近很少使用，更可能被回收
- kbdirty：等待回写到磁盘的内存数量，KB

Swap使用率:
- kbswpfree：空闲swap空间数量，KB
- kbswpused：已使用swap空间数量，KB
- %swpused：已使用swap空间百分比
- kbswpcad：已缓存的swap(内存换出换入一次, 当前仍在swap中, 此时内存不用换出, 可节省IO)，KB
- %swpcad：已缓存的sawap/已使用的swap

CPU利用率: 参考mpstat

报告inode状态, 文件状态和其他内核表:
- dentunusd：在目录缓存中未使用的缓存条目数量
- file-nr：系统使用的文件数量
- inode-nr: 系统使用的inode数量
- pty-nr: 系统使用的伪终端数量

Swap状态:
- pswpin/s ：每秒换入的swap分页数量
- pswpout/s：每秒换出的swap分页数量

进程创建和系统切换活动:
- proc/s ：每秒创建任务总数
- cswch/s：每秒上下文切换总数

## 例

```shell
$ sar 16 4 # sar 采集间隔 采集次数
Linux 3.16.0-34-generic (localhost) 	2015年04月20日 	_x86_64_	(8 CPU)

13时34分08秒     CPU     %user     %nice   %system   %iowait    %steal     %idle
13时34分24秒     all      2.03      0.00      0.62      0.04      0.00     97.31
13时34分29秒     all      2.23      0.00      0.74      0.00      0.00     97.03
13时34分35秒     all      2.10      0.00      0.73      0.04      0.00     97.13
平均时间:     all      2.08      0.00      0.67      0.03      0.00     97.22
```

## FAQ

1 . `无法打开 /var/log/sysstat/saXX: 没有那个文件或目录 Please check if data collecting is enabled in /etc/default/sysstat`

>方法1(**ubuntu推荐**):修改`/etc/default/sysstat`文件， 将 ENABLED 设置为 true,再重启系统.
>方法2:执行` sudo sar -o XX`创建文件
>方法3(fedora22):执行`sudo systemctl enable sysstat && sudo systemctl start sysstat`创建文件
