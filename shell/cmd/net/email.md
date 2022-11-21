# email
电子邮件系统基于邮件协议来完成电子邮件的传输，常见的邮件协议:
- 简单邮件传输协议(simple mail transfer protocol, SMTP) : 用于发送和中转发出的电子邮件，占用服务器的 TCP/25 端口。
- 邮局协议版本3(post office protocol3) : 用于将电子邮件存储到本地主机，占用服务器的 TCP/110 端口。
- Internet消息访问协议版本4(internet message access protocol4) : 用于在本地主机上访问邮件，占用服务器的 TCP/143 端口

在电子邮件系统中，为用户收发邮件的服务器名为邮件用户代理（ Mail User Agent，
MUA）. 

邮件投递代理（ MailDelivery Agent， MDA），其工作职责是把来自于邮件传输代理（ Mail Transfer Agent， MTA）
的邮件保存到本地的收件箱中。其中，这个 MTA 的工作职责是转发处理不同电子邮件服务供应商之间的邮件，把来自于 MUA 的邮件转发到合适的 MTA 服务器.

电子邮件的传输过程: MUA->MTA->MTA->MUA

企业级邮件系统的附加模块:
- 添加反垃圾与反病毒模块 : 它能够很有效地阻止垃圾邮件或病毒邮件对企业信箱的干扰。
- 对邮件加密 : 可有效保护邮件内容不被黑客盗取和篡改。
- 添加邮件监控审核模块 ：可有效地监控企业全体员工的邮件中是否有敏感词， 是否有透露企业资料等违规行为
- 保障稳定性: 电子邮件系统的稳定性至关重要，运维人员应做到保证电子邮件系统的稳定运行，并及时做好防范分布式拒绝服务（ Distributed Denial of Service， DDoS）攻击的准备

# Postfix
Postfix 服务程序的邮件收发能力强于 Sendmail 服务，而且能自动增加、减少进程的数量来保证电子邮件系统的高性能与稳定性。
另外， Postfix 服务程序由许多小模块组成，每个小模块都可以完成特定的功能，因此可在生产工作环境中根据需求灵活搭配.

配置项:
- myhostname 邮局系统的主机名
- mydomain 邮局系统的域名
- myorigin 从本机发出邮件的域名名称
- inet_interfaces 监听的网卡接口
- mydestination 可接收邮件的主机名或域名
- mynetworks 设置可转发哪些主机的邮件
- relay_domains 设置可转发哪些网域的邮件

Dovecot 是一款能够为 Linux 系统提供 IMAP 和 POP3 电子邮件服务的开源服务程序，安
全性极高，配置简单，执行速度快，而且占用的服务器硬件资源也较少，因此是一款值得推
荐的收件服务程序。

## example
```bash
# --- 配置服务器主机名称，需要保证服务器主机名称与发信域名保持一致
# hostnamectl set-hostname mail.linuxprobe.com
# --- 为电子邮件系统提供域名解析
# firewall-cmd --permanent --zone=public --add-service=dns
# firewall-cmd --reload
# dnf install bind-chroot
# cat /etc/named.conf
...
options {
	listen-on port 53 { any; };
	allow-query { any; };
...
# cat /etc/named.rfc1912.zones
zone "linuxprobe.com" IN {
	type master;
	file "linuxprobe.com.zone";
	allow-update {none;};
};
# cp -a /var/named/named.localhost /var/named/linuxprobe.com.zone
# vim /var/named/linuxprobe.com.zone
# systemctl restart named
# systemctl enable named
# --- 部署postfix
# dnf install postfix
# vim /etc/postfix/main.cf
...
myhostname = mail.linuxprobe.com # 保存服务器的主机名称
...
mydomain = linuxprobe.com # 用来保存邮件域的名称
...
myorigin = $mydomain # 定义发出邮件的域
...
inet_interfaces = all # 定义网卡监听地址。可以指定要使用服务器的哪些 IP 地址对外提供电子邮件服务；也可以干脆写成 all，表示所有 IP 地址都能提供电子邮件服务
...
mydestination = $myhostname, $mydomain # 定义可接收邮件的主机名或域名列表。这里可以直接调用前面定义好的 myhostname 和 mydomain 变量（如果不想调用变量，也可以直接调用变量中的值）
...
# --- 创建电子邮件系统的登录账户. Postfix 与 vsftpd 服务程序一样，都可以调用本地系统的账户和密码，因此在本地系统创建常规账户即可
# useradd liuchuan
# echo "linuxprobe" | passwd --stdin liuchuan
Changing password for user liuchuan.
passwd: all authentication tokens updated successfully.
# systemctl restart postfix
# systemctl enable postfix
# --- 部署dovecot
# dnf install -y dovecot
# vim /etc/dovecot/dovecot.conf
...
protocols = imap pop3 lmtp # 支持的协议
disable_plaintext_auth = no # 允许用户使用明文进行密码验证。之所以这样操作，是因为 Dovecot 服务程序为了保证电子邮件系统的安全而默认强制用户使用加密方式进行登录，
而由于当前还没有加密系统，因此需要添加该参数来允许用户的明文登录
...
login_trusted_networks = 192.168.10.0/24 # 设置允许登录的网段地址
...
# vim /etc/dovecot/conf.d/10-mail.conf # 配置邮件格式与存储路径。在 Dovecot 服务程序单独的子配置文件中，定义一个路径，用于指定要将收到的邮件存放到服务器本地的哪个位置
...
mail_location = mbox:~/mail:INBOX=/var/mail/%u
...
# firewall-cmd --permanent --zone=public --add-service=imap
# firewall-cmd --permanent --zone=public --add-service=pop3
# firewall-cmd --permanent --zone=public --add-service=smtp
# firewall-cmd --reload
# --- 用Foxmail/Thunderbird发送邮件
# dnf install -y thunderbird
# --- 在邮件服务器查看邮件
# dnf install mailx
# mailx
# quit
```

别名:
```bash
# cat /etc/aliases # 定义用户别名与邮件接收人的映射。除了使用本地系统中系统账户的名称外, 还可以自行定义一些别名来接收邮件
dream: root
...
# newaliases # 让新的用户别名配置文件立即生效
```