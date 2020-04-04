# dhcp
dhcp是一种基于 UDP 协议且仅限于在局域网内部使用的网络协议，主要用途是为局域网内部的设备或网络供应商自动分配 IP 地址等网络参数.

DHCP 服务器会自动把 IP 地址、子网掩码、网关、DNS 地址等网络信息分配给有需要的客户端，而且当客户端
的租约时间到期后还可以自动回收所分配的 IP 地址，以便交给新加入的客户端.

DHCP 术语:
- 作用域：一个完整的 IP 地址段，DHCP 协议根据作用域来管理网络的分布、分配 IP地址及其他配置参数
- 超级作用域：用于管理处于同一个物理网络中的多个逻辑子网段。超级作用域中包含了可以统一管理的作用域列表
- 排除范围：把作用域中的某些 IP 地址排除，确保这些 IP 地址不会分配给 DHCP客户端
- 地址池：在定义了DHCP的作用域并应用了排除范围后，剩余的用来动态分配给DHCP客户端的 IP 地址范围
- 租约：DHCP 客户端能够使用动态分配的 IP 地址的时间
- 预约：保证网络中的特定设备总是获取到相同的 IP 地址

dhcpd 是 Linux 系统中用于提供 DHCP 协议的服务程序.

dhcpd 服务程序配置文件中使用的常见参数以及作用:
ddns-update-style [类型] 
定义 DNS 服务动态更新的类型，类型包括 none（不支
持动态更新）、interim（互动更新模式）与 ad-hoc（特
殊更新模式）
[allow | ignore] client-updates 允许/忽略客户端更新 DNS 记录
default-lease-time [21600] 默认超时时间
max-lease-time [43200] 最大超时时间
option domain-name-servers [8.8.8.8] 定义 DNS 服务器地址
option domain-name ["domain.org"] 定义 DNS 域名
range 定义用于分配的 IP 地址池
option subnet-mask 定义客户端的子网掩码
option routers 定义客户端的网关地址
broadcase-address[广播地址] 定义客户端的广播地址
ntp-server[IP 地址] 定义客户端的网络时间服务器（NTP）
nis-servers[IP 地址] 定义客户端的 NIS 域服务器的地址
Hardware[网卡物理地址] 指定网卡接口的类型与 MAC 地址
server-name[主机名] 向 DHCP 客户端通知 DHCP 服务器的主机名
fixed-address[IP 地址] 将某个固定的 IP 地址分配给指定主机
time-offset[偏移误差] 指定客户端与格林尼治时间的偏移差