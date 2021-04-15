# ntp (Network Time Protocol)
参考:
- [ntpd时钟同步服务](http://xstarcd.github.io/wiki/sysadmin/ntpd.html)

> **推荐使用新版ntp协议实现: chrony**

ntp是用来使计算机时间同步化的一种协议.

## example
```bash
# ntpdate -q 192.168.0.223 # 仅查询, 不set时间
```

## FAQ
### ntpdate `-t`参数
通过实践: `-t`参数设置越大, ntpdate执行时间越长, 推荐使用`5`.

### ntp no server suitable for synchronization found
在ntp客户端运行ntpdate 10.0.0.106，出现`no server suitable for synchronization found`的错误.
在ntp客户端用ntpdate -d 10.0.0.106查看，发现有`Server dropped: Strata too high`的错误，并且显示`stratum 16`, 而正常情况下stratum这个值得范围是`0~15`.

这是因为NTP server还没有和其自身或者它的server同步上.

`/etc/ntp.conf`
```conf
# ss -anlp |grep 123 # 查看ntp使用的端口
# 计算本ntp server 与上层ntpserver的频率误差
driftfile /var/lib/ntp/drift

# 使用上层的internet ntp服务器
# server cn.pool.ntp.org prefer

restrict default kod nomodify notrap nopeer noquery # 默认拒绝任何操作即任何ip4地址、ip6地址  不能修改、不能trap远程登录、不能尝试对等、不能校对时间
restrict -6 default kod nomodify notrap nopeer noquery # restrict -6 表示IPV6地址的权限设置
restrict 127.0.0.1 # 允许本地所有操作
restrict -6 ::1    # # 允许本地所有操作(ipv6)
restrict 192.168.0.0 mask 255.255.0.0 nomodify notrap  # 允许的局域网络段或单独ip
restrict 0.0.0.0 mask 0.0.0.0 notrust           # 拒绝没有认证的用户端
server 192.168.1.12 

## 以下的定义是让NTP Server和其自身保持同步，如果在/ntp.conf中定义的server都不可用时，将使用local时间作为ntp服务提供给ntp客户端
server  127.127.1.0     # 127.127.1.0 is the local ntpd server address
fudge   127.127.1.0 stratum 10

#日志文件
logfile /var/log/ntp.log
```

restrict 用于控制相关权限, 用法为： `restrict [ 客户端IP ]  mask  [ IP掩码 ]  [参数]`, 其中IP地址也可以是default ，default 就是指所有的IP, 其他参数有以下几个：
- nomodify：户端不能使用NTPC与ntpq这两个程序来修改server的时间参数，但client仍可通过该server来进行网路校时
- notrust ：拒绝未认证的用戶端
- noquery ：客户端不能使用ntpq，ntpc来查询ntp服务器，等于不提供校对时间服务
- notrap ：不提供陷阱这个远端事件邮箱（远程事件日志）的功能
- nopeer ：用于阻止主机尝试与服务器对等
- kod ： 访问违规时发送 KoD 包，向不安全的访问者发送Kiss-Of-Death报文

在ntp server上重新启动ntp服务后，**ntp server自身或者与其server的同步的需要一个时间段，这个过程可能是5分钟**，在这个时间之内在客户端运行ntpdate命令时会产生no server suitable for synchronization found的错误.

> `ntpdate -u 192.168.1.12` : 仅更新系统时间不包括硬件

监控ntp server完成了和自身同步的过程可使用`watch ntpq -p`, `reach`达到一定数值时再同步就不会报该错误:
```log
Every 2.0s: ntpq -p

     remote           refid      st t when poll reach   delay   offset  jitter
==============================================================================
 192.168.0.73    .INIT.          16 u    -  256    0    0.000    0.000   0.000
*LOCAL(0)        .LOCL.          10 l   48   64  377    0.000    0.000   0.000
```

字段:
- remote 响应这个请求的NTP服务器的名称
- refid NTP服务器使用的更高一级服务器的名称

	- LOCAL : 本机(当没有远程节点或服务器可用时)
	- GPS : GPS
- st 正在响应请求的NTP服务器的级别
- when 上一次成功请求之后到现在的秒数
- poll : 下次更新在多少秒后

	本地和远程服务器下一次同步在多少时间之后，单位秒， 在一开始运行NTP的时候这个poll值会比较小，服务器同步的频率大，可以尽快调整到正确的时间范围，之后poll值会逐渐增大，同步的频率也就会相应减小
- reach 已经向上层ntp服务器要求更新的次数. 用来测试能否和服务器连接，是一个八进制值，每成功连接一次它的值就会增加
- delay 网络延迟. 从本地机发送同步要求到ntp服务器的往返时间
- offset : 时间补偿. 主机通过NTP时钟同步与所同步时间源的时间偏移量，单位为毫秒，offset越接近于0，主机和ntp服务器的时间越接近
- jitter : 系统时间与bios时间差. 统计了在特定个连续的连接数里offset的分布情况. 简单地说这个数值的绝对值越小，主机的时间就越精确

> ntpd是步进式的逐渐调整时间，而ntpdate是断点更新