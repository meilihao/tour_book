# dhcpd
dhcpd 是 Linux 系统中用于提供 DHCP 的服务程序.

```bash
# dnf install -y dhcp-server
# cat /etc/dhcp/dhcpd.conf
# 
# DHCP Server Configuration file.
# see /usr/share/doc/dhcp-server/dhcpd.conf.example
# see dhcpd.conf(5) man page
#
```

## 配置
一个标准的配置文件应该包括全局配置参数、子网网段声明、地址配置选项以及地址配置参数。其中，全局配置参数用于定义 dhcpd 服务程序的整体运行参数；子网网段声明用于
配置整个子网段的地址属性.

参数:
- ddns-update-style 类型 定义 DNS 服务动态更新的类型，类型包括 none（不支持动态更新）、 interim（互动更新模式）与ad-hoc（特殊更新模式）
- allow/ignore client-updates 允许/忽略客户端更新 DNS 记录
- default-lease-time 21600 默认超时时间
- max-lease-time 43200 最大超时时间
- option domain-name-servers 8.8.8.8 定义 DNS 服务器地址
- option domain-name "domain.org" 定义 DNS 域名
- range 定义用于分配的 IP 地址池
- option subnet-mask 定义客户端的子网掩码
- option routers 定义客户端的网关地址
- broadcast-address 广播地址 定义客户端的广播地址
- ntp-server IP 地址 定义客户端的网络时间服务器（ NTP）
- nis-servers IP 地址 定义客户端的 NIS 域服务器的地址
- hardware 硬件类型 MAC 地址 指定网卡接口的类型与 MAC 地址
- server-name 主机名 向 DHCP 客户端通知 DHCP 服务器的主机名
- fixed-address IP 地址 将某个固定的 IP 地址分配给指定主机
- time-offset 偏移差 指定客户端与格林尼治时间的偏移差

### 自动管理 IP 地址
```bash
# vim /etc/dhcp/dhcpd.conf
ddns-update-style none; # 设置 DNS 服务不自动进行动态更新
ignore client-updates; # 忽略客户端更新 DNS 记录
subnet 192.168.10.0 netmask 255.255.255.0 { # 作用域为 192.168.10.0/24 网段
	range 192.168.10.50 192.168.10.150; # IP 地址池为 192.168.10.50-150
	option subnet-mask 255.255.255.0; # 定义客户端默认的子网掩码
	option routers 192.168.10.1; # 定义客户端默认的gateway
	option domain-name "linuxprobe.com"; # 定义默认的搜索域
	option domain-name-servers 192.168.10.1; # 定义客户端的 DNS 地址
	default-lease-time 21600;
	max-lease-time 43200;
}
# systemctl start dhcpd
# systemctl enable dhcpd
# firewall-cmd --zone=public --permanent --add-service=dhcp
# firewall-cmd --reload
```

### 分配固定 IP 地址
根据网卡mac来固定ip

```bash
# vim /etc/dhcp/dhcpd.conf
ddns-update-style none;
ignore client-updates;
subnet 192.168.10.0 netmask 255.255.255.0 {
	range 192.168.10.50 192.168.10.150;
	option subnet-mask 255.255.255.0;
	option routers 192.168.10.1;
	option domain-name "linuxprobe.com";
	option domain-name-servers 192.168.10.1;
	default-lease-time 21600;
	max-lease-time 43200;
	host linuxprobe {
		hardware ethernet 00:0c:29:dd:f2:22; # mac
		fixed-address 192.168.10.88;
	}
}
# systemctl restart dhcpd
```