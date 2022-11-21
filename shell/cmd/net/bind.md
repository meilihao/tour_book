# bind
dns服务器类型:
1. 主服务器: 在特定区域内具有唯一性，负责维护该区域内的域名与 IP 地址之间的对应关系
1. 从服务器: 从主服务器中获得域名与 IP 地址的对应关系并进行维护，以防主服务器宕机等情况
1. 缓存服务器: 通过向其他域名解析服务器查询获得域名与 IP 地址的对应关系，并将经常查询的域名信息保存到服务器本地，以此来提高重复查询时的效率

	DNS 缓存服务器是一种不负责域名数据维护的 DNS 服务器

DNS 域名解析服务采用分布式的数据结构来存放海量的“区域数据” 信息，在执行用户发起的域名查询请求时，具有递归查询和迭代查询两种方式:
1. 递归查询, 是指 DNS 服务器在收到用户发起的请求时，必须向用户返回一个准确的查询结果. 如果 DNS 服务器本地没有存储与之对应的信息，则该服务器需要询问其他服务器，并将返回的查询结果提交给用户.
1. 迭代查询, DNS 服务器在收到用户发起的请求时，并不直接回复查询结果，而是告诉另一台DNS 服务器的地址，用户再向这台 DNS 服务器提交请求，这样依次反复，直到返回查询结果

域名解析记录类型:
- A 将域名指向一个 IPv4 地址
- CNAME 将域名指向另外一个域名
- AAAA 将域名指向一个 IPv6 地址
- NS 将子域名指定由其他 DNS 服务器解析
- MX 将域名指向邮件服务器地址
- SRV 记录提供特定的服务的服务器
- TXT 文本内容一般为 512 字节，常作为反垃圾邮件的 SPF（ Sender Policy Framework，发送方策略框架） 记录
- CAA CA 证书颁发机构授权校验
- 显性 URL 将域名重定向到另外一个地址
- 隐性 URL 与显性 URL 类似，但是会隐藏真实目标地址

bind配置:
- 主配置文件(/etc/named.conf): 这些参数用来定义 bind 服务程序的运行

	```bash
	# dnf install bind-chroot # bind 服务程序的名称为 named
	# vim /etc/named.conf
	..

	options {
		listen-on port 53 { any; }; // any: 监听所有端口
		...
		allow-query { any; }; // any: 允许所有人对本服务器发送 DNS 查询请求
	...
	};
	...
	include "/etc/named.rfc1912.zones"; # 区域配置文件
	...
	# named-checkconf # 检查主配置
	# named-checkzone # 检查区域配置
	```
- 区域配置文件(/etc/named.rfc1912.zones): 用来保存域名和 IP 地址对应关系的所在位置. 类似于图书的目录，对应着每个域和相应 IP 地址所在的具体位置，当需要查看或修改时，可根据这个位置找到相关文件。

	可配置:
	- 正向解析: 域名解析为 IP 地址

		```bash
		# vim /etc/named.rfc1912.zones
		zone "linuxprobe.com" IN {
			type master;
			file "linuxprobe.com.zone"; # 具体解析规则
			allow-update {none;};
		};
		# cp -a /var/named/named.localhost linuxprobe.com.zone
		# vim linuxprobe.com.zone
		# systemctl restart named
		# nslookup
		```
	- 反向解析: 将 IP 地址解析为域名

		```bash
		# vim /etc/named.rfc1912.zones
		zone "10.168.192.in-addr.arpa" IN {
			type master;
			file "192.168.10.arpa";
			allow-update {none;};
		};
		# cp -a /var/named/named.loopback 192.168.10.arpa
		# vim 192.168.10.arpa
		```
- 数据配置目录(/var/named): 该目录用来保存域名和 IP 地址真实对应关系的数据配置文件

## 配置主从
主从dns:
```bash
# -- 主配置
# vim /etc/named.rfc1912.zones
zone "linuxprobe.com" IN {
	type master;
	file "linuxprobe.com.zone";
	allow-update { 192.168.10.20; }; # 配置allow-update可指定dns从服务器
};
# --- 从配置
# vim /etc/named.rfc1912.zones
zone "linuxprobe.com" IN {
	type slave;
	masters { 192.168.10.10; }; # masters表示可以有多个主服务器
	file "slaves/linuxprobe.com.zone";
};
```

从dns服务器获取的数据保存在`/var/named/slaves`

## 配置缓存服务器
```bash
# vim /etc/named.conf
...
options {
	...
	forwarders { 8.8.8.8; }; # 配置上级 DNS 服务器(即可获取数据配置文件的服务器)
...
# systemctl restart named
# firewall-cmd --permanent --zone=public --add-service=dns
# firewall-cmd --reload
# nslookup
```

## 分离解析
分离解析即让位于不同地理范围内的用户就近访问服务(服务已在多地域部署).

```bash
# vim /etc/named.conf
zone "." IN { # 需要删除: 配置的 DNS 分离解析功能与 DNS 根服务器配置参数有冲突
 type hint;
 file "named.ca";
};
# vim /etc/named.rfc1912.zones
acl "china" { 122.71.115.0/24; };
acl "america" { 106.185.25.0/24; };
view "china"{
	match-clients { "china"; };
	zone "linuxprobe.com" {
		type master;
		file "linuxprobe.com.china";
	};
};
view "america" {
	match-clients { "america"; };
	zone "linuxprobe.com" {
		type master;
		file "linuxprobe.com.america";
	};
};
# --- 建立数据配置文件 from 模板
# cd /var/named
# cp -a named.localhost linuxprobe.com.china
# cp -a named.localhost linuxprobe.com.america
# vim linuxprobe.com.china
# vim linuxprobe.com.America
# systemctl restart named
```

## 安全
bind 服务程序为了提供安全的解析服务，已经对 TSIG（见 RFC 2845） 加密机制提供了支持。

TSIG 主要是利用了密码编码的方式来保护区域信息的传输（ Zone Transfer），即 TSIG 加密机制保证了 DNS 服务器之间传输域名区域信息的安全性.

```bash
# dnssec-keygen -a HMAC-MD5 -b 128 -n HOST master-slave # 生成一个主机名称为 master-slave 的 128 位 HMAC-MD5 算法的密钥文件。在执行该命令后默认会在当前目录中生成公钥和私钥文件
Kmaster-slave.+157+62533
# cd /var/named/chroot/etc/
# vim transfer.key
key "master-slave" {
algorithm hmac-md5;
secret "NI6icnb74FxHx2gK+0MVOg==";
};
# chown root:named transfer.key # 为了安全起见，需要将文件的所属组修改成 named，并将文件权限设置得要小一点， 然后设置该文件的一个硬链接，并指向/etc 目录
# chmod 640 transfer.key
# ln transfer.key /etc/transfer.key
# --- 主配置
# vim /etc/named.conf
...
include "/etc/transfer.key";
options {
...
	allow-transfer { key master-slave; };
...
[root@linuxprobe~]# systemctl restart named
# --- 从配置, 先配置transfer.key, 同主, transfer.key也来源于主
# rm -rf /var/named/slaves/* && systemctl restart named && ls /var/named/slaves/ # 可验证主设置TSIG后, 从是否还能同步
# vim /etc/named.conf
...
include "/etc/transfer.key";
...
server 192.168.10.10
{
keys { master-slave; };
};
...
# systemctl restart named
# ls /var/named/slaves/
```

dnssec-keygen参数:
- -a 指定加密算法，包括 RSA MD5（ RSA）、 RSA SHA1、 DSA、 NSEC3RSASHA1、 NSEC3DSA 等
- -b 密钥长度（ HMAC-MD5 的密钥长度在 1～512 位之间）
- -n 密钥的类型（ HOST 表示与主机相关）