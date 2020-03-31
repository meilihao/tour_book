# dns
DNS（Domain Name System，域名系统）是一项用于管理和解析域名与 IP 地址对应关系的技术.

DNS 提供了下面三种类型的服务器:
- 主服务器：在特定区域内具有唯一性，负责维护该区域内的域名与 IP 地址之间的对应关系
- 从服务器：从主服务器中获得域名与 IP 地址的对应关系并进行维护，以防主服务器宕机等情况
- 缓存服务器：通过向其他域名解析服务器查询获得域名与 IP 地址的对应关系，并将经常查询的域名信息保存到服务器本地，以此来提高重复查询时的效率

DNS 域名解析方式:
- 递归查询 :  DNS 服务器在收到用户发起的请求时，必须向用户返回一个准确的查询结果. 如果 DNS 服务器本地没有存储与之对应的信息，则该服务器需要询问其他服务器，并将返回的查询结果提交给用户.
- 迭代查询 : DNS 服务器在收到用户发起的请求时，并不直接回复查询结果，而是告诉另一台 DNS 服务器的地址，用户再向这台 DNS 服务器提交请求，这样依次反复，直到返回查询结果.

本机-> local dns是递归查询
local dns -> 权威dns是迭代

## bind
```bash
# yum install bind-chroot
# cp -a /var/named/named.localhost linux.com.zone # 获取正向解析的模板
# cp -a /var/named/named.loopback 192.168.10.arpa # 获取反向解析的模板
# systemctl restart named
```

bind 服务程序中有下面这三个比较关键的文件:
- 主配置文件（/etc/named.conf）：定义 bind 服务程序的运行
- 区域配置文件（/etc/named.rfc1912.zones）：用来保存域名和 IP 地址对应关系的所在位置. 类似于图书的目录，对应着每个域和相应 IP 地址所在的具体位置，当需要查看或修改时，可根据这个位置找到相关文件
- 数据配置文件目录（/var/named）：该目录用来保存域名和 IP 地址真实对应关系的数据配置文件

> bind 服务程序的名称为 named.

> bind支持通过ip地域来分别进行dns处理的分离解析技术.

正向解析是指根据域名（主机名）查找到对应的 IP 地址.
反向解析的作用是将用户提交的 IP 地址解析为对应的域名信息，它一般用于对某个 IP 地址上绑定的所有域名进行整体屏蔽，屏蔽由某些域名发送的垃圾邮件.

互联网中的绝大多数 DNS 服务器（超过 95%）都是基于 BIND 域名解析服务搭建的，而bind 服务程序通过 TSIG (RFC 2845) 保证了 DNS 服务器之间传输域名区域信息的安全性. TSIG主要是利用了密码编码的方式来保护区域信息的传输（Zone Transfer）.

## 方案
- DNS over TLS
- DNS over HTTPS

> 均是基于IETF RFC, [这也是不选用DNSCrypt的原因之一](https://tenta.com/blog/post/2017/12/dns-over-tls-vs-dnscrypt).
>
> Cloudflare公共DNS解析器使用[开源Knot解析器](https://blog.cloudflare.com/dns-resolver-1-1-1-1/)
>
> [相关RFC](https://blog.cloudflare.com/dns-resolver-1-1-1-1/)

类似软件:
- [隐私优先的 DNS 解决方案 Tenta DNS](https://github.com/tenta-browser/tenta-dns)