# ftp
FTP 服务器是按照 FTP 协议在互联网上提供文件存储和访问服务的主机, FTP 客户端则是向服务器发送连接请求，以建立数据传输链路的主机.

FTP 协议有下面两种工作模式:
1. FTP 服务器主动向客户端发起连接请求
1. FTP 服务器等待客户端发起连接请求（默认工作模式）

# vsftpd
vsftpd（ very secure ftp daemon，非常安全的 FTP 守护进程）是一款运行在 Linux 操作系统上的 FTP 服务程序，不仅完全开源而且免费.

vsftpd 服务程序的主配置文件是`/etc/vsftpd/vsftpd.conf`.

vsftpd 服务程序常用的参数以及作用:
listen=[YES|NO] 是否以独立运行的方式监听服务
listen_address=IP 地址 设置要监听的 IP 地址
listen_port=21 设置 FTP 服务的监听端口
download_enable＝[YES|NO] 是否允许下载文件
userlist_enable=[YES|NO] 设置用户允许功能
userlist_deny=[YES|NO] 设置用户禁止功能
max_clients=0 最大客户端连接数， 0 为不限制
max_per_ip=0 同一 IP 地址的最大连接数， 0 为不限制
anonymous_enable=[YES|NO] 是否允许匿名用户访问
anon_upload_enable=[YES|NO] 是否允许匿名用户上传文件
anon_umask=022 匿名用户上传文件的 umask 值
anon_root=/var/ftp 匿名用户的 FTP 根目录
anon_mkdir_write_enable=[YES|NO] 是否允许匿名用户创建目录
anon_other_write_enable=[YES|NO] 是否开放匿名用户的其他写入权限（包括重命名、删除等操作权限）
anon_max_rate=0 匿名用户的最大传输速率（字节/秒）， 0 为不限制
local_enable=[YES|NO] 是否允许本地用户登录 FTP
local_umask=022 本地用户上传文件的 umask 值
local_root=/var/ftp 本地用户的 FTP 根目录
chroot_local_user=[YES|NO] 是否将用户权限禁锢在 FTP 目录，以确保安全
local_max_rate=0 本地用户最大传输速率（字节/秒）， 0 为不限制

vsftpd 作为更加安全的文件传输协议服务程序，允许用户以 3 种认证模式登录 FTP 服务器:
1. 是最不安全的一种认证模式，任何人都可以无须密码验证而直接登录到 FTP 服务器

	config:
	```bash
	anonymous_enable=YES 允许匿名访问模式
	anon_umask=022 匿名用户上传文件的 umask 值
	anon_upload_enable=YES 允许匿名用户上传文件
	anon_mkdir_write_enable=YES 允许匿名用户创建目录
	anon_other_write_enable=YES 允许匿名用户修改目录名称或删除目录
	```
1. 是通过 Linux 系统本地的账户密码信息进行认证的模式，相较于匿名开放模式更安全，而且配置起来也很简单。但是如果黑客破解了账户的信息，就可以畅通无阻地登录 FTP 服务器，从而完全控制整台服务器

	config:
	```bash
	anonymous_enable=NO 禁止匿名访问模式
	local_enable=YES 允许本地用户模式
	write_enable=YES 设置可写权限
	local_umask=022 本地用户模式创建文件的 umask 值
	userlist_deny=YES 启用“禁止用户名单”，名单文件为 ftpusers 和 user_list
	userlist_enable=YES 开启用户作用名单文件功能
	```

	vsftpd服务程序所在的目录中默认存放着两个名为“用户名单” 的文件（ ftpusers 和 user_list）, 只要里面写有某位用户的名字，就不再允许这位用户登录到 FTP 服务器上.

	如果把上面主配置文件中 userlist_deny 的参数值改成 NO，那么 user_list 列表就变成了强制白名单。 它的功能与之前完全相反，只允许列表内的用户访问，拒绝其他人的访问.

	在采用本地用户模式登录 FTP 服务器后，默认访问的是该用户的家目录，而且该目录的默认所有者、所属组都是该用户自己，因此不存在写入权限不足的情况.
1. 更安全的一种认证模式，它需要为 FTP 服务单独建立用户数据库文件，虚拟出用来进行密码验证的账户信息，而这些账户信息在服务器系统中实际上是不存在的，仅供 FTP 服务程序进行认证使用。这样，即使黑客破解了账户信息也无法登录服务器， 从而有效降低了破坏范围和影响

	```bash
	# vim vuser.list # 分别创建 zhangsan 和 lisi 两个用户，密码均为 redhat
	zhangsan
	redhat
	lisi
	redhat
	# db_load -T -t hash -f vuser.list vuser.db
	# chmod 600 vuser.db && rm -f vuser.list
	# useradd -d /var/ftproot -s /sbin/nologin virtual # 创建一个可以映射到虚拟用户的系统本地用户
	# vim /etc/pam.d/vsftpd.vu # 新建一个用于虚拟用户认证的 PAM 文件, 其中 PAM 文件内的“db=” 参数为使用db_load 命令生成的账户密码数据库文件的路径，但不用写数据库文件的后缀
	auth required pam_userdb.so db=/etc/vsftpd/vuser
	account required pam_userdb.so db=/etc/vsftpd/vuser
	# mkdir /etc/vsftpd/vusers_dir/ && cd /etc/vsftpd/vusers_dir/ # 为虚拟用户设置不同的权限
	# touch lisi # 只能登录，没有其他权限
	# vusers_dir]# vim zhangsan # 写入允许的相关权限（使用匿名用户的参数）
	anon_upload_enable=YES
	anon_mkdir_write_enable=YES
	anon_other_write_enable=YES
	```

	config:
	```bash
	anonymous_enable=NO 禁止匿名开放模式
	local_enable=YES 允许本地用户模式
	guest_enable=YES 开启虚拟用户模式
	guest_username=virtual 指定虚拟用户账户
	pam_service_name=vsftpd.vu 指定 PAM 文件
	allow_writeable_chroot=YES 允许对禁锢的 FTP 根目录执行写入操作，而且不拒绝用户的登录请求
	pam_service_name=vsftpd.vu
	user_config_dir=/etc/vsftpd/vusers_dir
	```

	PAM 作为应用程序层与鉴别模块层的连接纽带，可以让应用程序根据需求灵活地在自身插入所需的鉴别功能模块。当应用程序需要 PAM 认证时，则需要在应用程序中定义负责认证的 PAM 配置文件，实现所需的认证功能.

## FAQ
### selinux
```bash
setsebool -P ftpd_full_access=on # 避免被selinux拦截
```

# TFTP
单文件传输协议（ Trivial File Transfer Protocol， TFTP）是一种基于 UDP 协议在客户端和服务器之间进行简单文件传输的协议。顾名思义，它提供不复杂、开销不大的文件传输服
务，可将其当作 FTP 协议的简化版本.

TFTP 的命令功能不如 FTP 服务强大，甚至不能遍历目录，在安全性方面也弱于 FTP 服务。而且，由于 TFTP 在传输文件时采用的是 UDP 协议，占用的端口号为 69，因此文件的传
输过程也不像 FTP 协议那样可靠。但是，因为 TFTP 不需要客户端的权限认证，也就减少了无谓的系统和网络带宽消耗，因此在传输琐碎（ trivial）不大的文件时，效率更高.

tftp-server 是服务程序， tftp 是用于连接测试的客户端工具， xinetd 是管理服务.

> xinetd 服务可以用来管理多种轻量级的网络服务，而且具有强大的日志功能。它专门用于控制那些比较小的应用程序的开启与关闭.

在RHEL 8 系统中, tftp 所对应的配置文件默认不存在, 需要用户根据示例文件（ /usr/share/doc/xinetd/sample.conf）自行创建.

```bash
# vim /etc/xinetd.d/tftp
service tftp
{
	socket_type = dgram
	protocol = udp
	wait = yes
	user = root
	server = /usr/sbin/in.tftpd
	server_args = -s /var/lib/tftpboot
	disable = no
	per_source = 11
	cps = 100 2
	flags = IPv4
}
```

## tftp
### 选项
- ? 帮助信息
- put 上传文件
- get 下载文件
- verbose 显示详细的处理信息
- status 显示当前的状态信息
- binary 使用二进制进行传输
- ascii 使用 ASCII 码进行传输
- timeout 设置重传的超时时间
- quit 退出