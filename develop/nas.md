# fs
参考:
- [NAS 最佳实践](https://help.aliyun.com/document_detail/132279.html)
- [NAS产品规格限制](https://www.alibabacloud.com/help/zh/doc-detail/122195.htm)

阿里云NAS支持情况: NFSv3.0/4.0+, SMB2.1+. nfs仅支持linux, smb仅支持windows.

总结:
- 跨平台挂载会因为字符集导致乱码
- smb2.0+ Protocol不支持unix通用权限, 导致mount.cifs挂载时权限显示不正确.

## CIFS, SMB, NFS
SMB(Server Message Block，即服务(器)消息块) 是 IBM 公司在 80年代中期发明的一种文件共享协议. 它只是系统之间通信的一种方式（协议）. 目前最新版是`v3.1.1`.
CIFS是微软的Common Internet file system的缩写, 是 SMB 协议的一种特殊实现, CIFS 实现的协议至今仍很少被使用. 大多数现代存储系统不再使用 CIFS，而是使用 SMB2 或 SMB3. 在 Windows 系统环境中，SMB2 和 SMB3 是事实使用的标准.
Samba 也是 SMB 协议的实现, 常用于windows与类unix间的文件共享.
NFS是SUN为Unix开发的网络文件系统, 提供类unix间的文件共享. 目前最新版本是`v4.2`. NFSv4用户验证采用“用户名+域名”的模式，与Windows AD验证方式类似，NFSv4可使用Kerberos验证方式.（Kerberos与Windows AD都遵循相同RFC1510标准），这样方便windows和`*nix`环境混合部署.

> nfs server端权限变化后client端无需重新mount即可生效.

autofs 自动挂载服务: 无论是 Samba 服务还是 NFS 服务，都要把挂载信息写入到/etc/fstab 中，这样远程共享资源就会自动随服务器开机而进行挂载. autofs 服务程序是一种 Linux 系统守护进程，当检测到用户视图访问一个尚未挂载的文件系统时，将自动挂载该
文件系统.

>  RHEL 7 开始不支持NFSv2

## NFS
参考:
- [管理权限组](https://help.aliyun.com/document_detail/27534.html)
- [acl](/shell/cmd/acl.md)
- [rhel 8 Chapter 3. Exporting NFS shares](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/deploying_different_types_of_servers/index#exporting-nfs-shares_Deploying-different-types-of-servers)
- [Common NFS Mount Options](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html-single/storage_administration_guide/index#ch-nfs)
- [Securing NFS](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/deploying_different_types_of_servers/index#securing-nfs_Deploying-different-types-of-servers)
- [如何在CentOS 8安装并配置NFS服务](https://www.myfreax.com/how-to-install-and-configure-an-nfs-server-on-centos-8/)
- [aAmazon Elastic File System(nas) : 文件系统中文件和目录的用户和组 ID 权限 即 rwx模型](https://docs.aws.amazon.com/zh_cn/efs/latest/ug/efs-ug.pdf)
- [pNFS](https://wenku.baidu.com/view/7cd3eee26294dd88d0d26b0c.html)
- [windows 支持nfs的版本](https://docs.microsoft.com/en-us/windows-server/storage/nfs/nfs-overview)

> NFS 客户端为内核的一部分，由于部分内核存在一些缺陷，会影响 NFS 的正常使用, 见[NFS 客户端已知问题](https://www.alibabacloud.com/help/zh/doc-detail/114129.htm)

> NFS v4.1开始支持[Parallel NFS (pNFS)](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html-single/storage_administration_guide/index#ch-nfs).

> **推荐使用以上命令通过 NFSv3 协议挂载，获得最佳性能. 如果应用依赖文件锁，也即需要使用多个client 同时编辑一个文件时使用 NFSv4 协议挂载**

> [nfsv4不再需要rpcbind, rpc.statd, lockd, rpc.mountd服务](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/deploying_different_types_of_servers/index#services-required-by-nfs_exporting-nfs-shares), 但其他rpc服务还是需要: `systemctl mask --now rpc-statd.service rpcbind.service rpcbind.socket`

> mount.nfs不支持bind client ip, 见FAQ的"unmatch host"

安装:
```
$ sudo apt install nfs-kernel-server
$ sudo apt install nfs-common # Install NFS client
$ sudo yum install nfs-utils # Install NFS client
$ sudo systemctl status nfs-kernel-server
$ systemctl start nfs-server # from centos7, 启动nfs
$ sudo cat /proc/fs/nfsd/versions # 查看nfs server支持的nfs protocol version, nfs服务需先启动
$ nfsstat --server # nfs server status
$ nfsstat -s # server使用的nfs version
$ nfsstat -c # client使用的nfs version
$ nfsstat -m # 在client端已挂载的nfs信息
$ nfsstat -4 # 查看NFS版本4的状态
$ showmount -e 192.168.0.83 # 在 Client 端查看server端(192.168.0.83)共享出来的目录

	- -e : 显示 NFS 服务器的共享列表
	- -a : 显示本机挂载的文件资源的情况
	- -v : 显示版本号
$ sudo mount -t nfs -o vers=4.2 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享, 挂载成功后不能访问请检查nfs server端的权限
$ sudo mount -t nfs4 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享
$ sudo mount -o v4.2 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享
$ df -h #查看挂载情况
$ sudo umount /mnt
$ cat /etc/exports
/usr/local/files/mypool/share  *(rw,sync,all_squash,anonuid=1037)
```

```bash
# from [手动挂载NFS文件系统](https://help.aliyun.com/document_detail/90529.html)
# 有利于提高同时发起的NFS请求数量
sudo echo "options sunrpc tcp_slot_table_entries=128" >> /etc/modprobe.d/sunrpc.conf
sudo echo "options sunrpc tcp_max_slot_table_entries=128" >> /etc/modprobe.d/sunrpc.conf
# 挂载时使用了noresvport参数，以规避NFS文件系统卡住的风险, 但[部分kernel不支持需检查否则mount时会报错](https://help.aliyun.com/document_detail/121264.html)
sudo mount -t nfs -o vers=3,nolock,proto=tcp,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport 3f0954ac37-kaf99.cn-shanghai.nas.aliyuncs.com:/ /mnt
sudo mount -t nfs -o vers=4,minorversion=0,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport 3f0954ac37-kaf99.cn-shanghai.nas.aliyuncs.com:/ /mnt
vim /etc/fstab
# from [自动挂载NFS文件系统](https://help.aliyun.com/document_detail/91476.html)
# 防止客户端在网络就绪之前开始挂载文件系统
file-system-id.region.nas.aliyuncs.com:/ /mnt nfs vers=4,minorversion=0,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,_netdev,noresvport 0 0
file-system-id.region.nas.aliyuncs.com:/ /mnt nfs vers=3,nolock,proto=tcp,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,_netdev,noresvport 0 0
```

NFS server 的配置选项在 /etc/default/nfs-kernel-server 和 /etc/default/nfs-common 里.
`/etc/exports`是用来管理NFS共享目录的使用权限与安全设置的地方. 特别注意的是，NFS本身设置的是网络共享权限，整个共享目录的权限还和目录自身的系统权限有关.
/var/lib/nfs/etab                      记录NFS共享出来的目录的完整权限设定值, from "man exportfs"
/var/lib/nfs/xtab                      记录曾经登录过的客户端信息

>　nfs server指定使用的版本: `/etc/default/nfs-kernel-server`的`RPCMOUNTDOPTS="--manage-gids -V 4.2"`.

`/ect/fatab`:
```
192.168.0.10:/nfs_share /mnt/nfs nfs defaults 0 0
```

### FS系统守护进程
- nfsd ：它是基本的NFS守护进程，主要功能是通过登入者ip, 用户id等管理客户端是否能够登录服务器

	支持`/etc/exports.d/*.exports`
- rpc.mountd ：主要功能是管理NFS的文件系统. 当客户端顺利通过nfsd登录NFS服务器后，在使用NFS服务所提供的文件前，还必须通过文件使用权限的验证. 它会读取NFS的配置文件/etc/exports来对比客户端权限.
- lockd : 用在管理档案的锁定 (lock) 用途. 当多个客户端同时尝试写入某个档案时， 需要lockd 来解决多客户端同时写入的问题. 但 lockd 必须要同时在客户端与服务器端都开启才行. 此外， lockd 也常与 rpc.statd 同时启用.
- statd : 检查文件的一致性，与lockd有关. 若发生因为客户端同时使用同一档案造成档案可能有所损毁时， statd 可以用来检测并尝试恢复该档案. 与 lockd 同样的，这个功能必须要在服务器端与客户端都启动才会生效.
- rpc.idmapd : 名字映射后台进程
- rpcbind : 主要功能是进行端口映射工作. 当客户端尝试连接并使用RPC服务器提供的服务（如NFS服务）时，rpcbind会将所管理的与服务对应的端口提供给客户端，从而使客户可以通过该端口向服务器请求服务, 因此rpcbind必须在nfs前启动.

### 其他相关命令

1. exportfs

允许root用户有选择地导出或取消导出目录，而无需重新启动NFS服务.

使/etc/exports的配置立刻生效，该命令格式如下：

　　# exportfs [-aruv]

　　-a 全部挂载或卸载 /etc/exports中的内容
　　-r 重新读取/etc/exports 中的信息 ，并同步更新/var/lib/nfs/xtab
　　-u 卸载单一目录（和-a一起使用为卸载所有/etc/exports文件中的目录）
　　-v 输出详细信息

具体例子:
		 # exportfs # 默认输出当前已导出文件系统的列表
　　# exportfs -au #  卸载所有共享目录
		 # exportfs -ra # 刷新nfs export, **推荐**. 已挂载的fs被取消export时,mounted端操作会导致报`ls: 无法访问'xxx': 过旧的文件控柄`
		 # exportfs -u 127.0.0.1:/scratch/test # 卸载单一目录
　　# exportfs -rv 重新加载共享所有目录并输出详细信息
		 # exportfs -o rw,no_root_squash 127.0.0.1:/scratch/test # 将/scratch/test共享给127.0.0.1, 信息不会写入`/etc/exports`, 但可用`showmount -e  ${nfs server ip}`查到

1. nfsstat

查看NFS的运行状态，对于调整NFS的运行有很大帮助。

1. rpcinfo

查看rpc执行信息，可以用于检测rpc运行情况的工具，利用rpcinfo -p 可以查看出RPC开启的端口所提供的程序有哪些

1. showmount

　　-a 显示已经于客户端连接上的目录信息
　　-e IP或者hostname 显示此IP地址共享出来的目录

1. netstat

可以查看出nfs服务开启的端口，其中nfs 开启的是2049，rpcbind 开启的是111，其余则是rpc开启的

最后注意两点，虽然通过权限设置可以让普通用户访问，但是挂载的时候默认情况下只有root可以去挂载，普通用户可以执行sudo。

NFS server 关机的时候一点要确保NFS服务关闭，没有客户端处于连接状态！通过showmount -a 可以查看，如果有的话用kill killall pkill 来结束.

### /etc/exports
格式：`export host1(options1) host2(options2) host3(options3) ...` from `man exports`
说明:
- 输出目录 : NFS系统中需要共享给客户机使用的目录
- 客户端 : 网络中可以访问这个NFS输出目录的计算机

	客户端常用的指定方式:

	- 指定ip地址的主机：192.168.0.200
	- 指定子网中的所有主机：192.168.0.0/24 192.168.0.0/255.255.255.0
	- 指定域名的主机：david.bsmart.cn
	- 指定域中的所有主机：*.bsmart.cn
	- 所有主机：*
选项：用来设置输出目录的访问权限、用户映射等

	NFS主要有3类选项：

	1. 访问权限选项 : 最终能不能读写，还是与文件系统的 rwx 及身份有关

	    设置输出目录只读：ro
	    设置输出目录读写：rw

	1. 用户映射选项

	    all_squash：将远程访问的所有普通用户及所属组都映射为匿名用户/用户组（默认是`nobody:nogroup`）, 可由anonuid/anongid指定
	    no_all_squash：与all_squash取反（默认设置）
	    root_squash：将root用户及所属组都映射为匿名用户/用户组（默认是`nobody:nogroup`， 如此对服务器的系统会较有保障）, 可由anonuid/anongid指定
	    no_root_squash：与rootsquash取反, 允许使用 root 身份来操作服务器的文件系统. 这个选项会留下严重的安全隐患，一般不建议采用.
	    anonuid=xxx：将远程访问的所有用户都映射为匿名用户，并指定该用户为本地用户（UID=xxx, 该 UID 必需要存在于你的 /etc/passwd 当中）
	    anongid=xxx：将远程访问的所有用户组都映射为匿名用户组账户，并指定该匿名用户组账户为本地用户组账户（GID=xxx）

	1. 其它选项

	    secure：限制客户端只能从小于1024的tcp/ip端口连接nfs服务器（默认设置）
	    insecure：允许客户端从大于1024的tcp/ip端口连接服务器
	    sync：将数据同步写入内存缓冲区与磁盘中，效率低，但可以保证数据的一致性
	    async：将数据先保存在内存缓冲区中，必要时才写入磁盘
	    wdelay：检查是否有相关的写操作，如果有则将这些写操作一起执行，这样可以提高效率（默认设置）
	    no_wdelay：若有写操作则立即执行，应与sync配合使用
	    subtree_check：若输出目录是一个子目录，则nfs服务器将检查其父目录的权限. 在客户端打开文件时重命名该文件会导致许多问题. 在几乎所有情况下，最好禁用子树检查.
	    no_subtree_check：即使输出目录是一个子目录，nfs服务器也不检查其父目录的权限，这样可以提高效率, (默认设置)

> nfs 支持使用no_acl来禁用acl.

### 身份映射(`/etc/idmapd.conf`)
NFS服务虽然不具备用户身份验证的功能，但是NFS提供了一种身份映射的机制来对用户身份进行管理. 当客户端访问NFS服务时，服务器会根据情况将客户端用户的身份映射成NFS匿名用户`nobody:nogroup`. `nobody:nogroup`是由linux中自动创建的一个用户账号，该账号不能用于登录系统，**专门用作服务的匿名用户账号**.

用户身份重叠, 是指在NFS服务采用默认设置（用户身份映射选项为root_squash）时，如果在服务器端赋予某个用户对共享目录具有相应权限，而且在客户端恰好也有一个具有相同uid的用户，那么当在客户端以该用户身份访问共享时，将自动具有服务器端对应用户的权限.

根源: 对于Linux系统而言，区分不同用户的唯一标识就是uid.

避免用户身份重叠：
1. 在设置NFS共享时，建议采用`all_squash`选项，将所有客户端用户均映射为`nobody:nogroup`.
1. 严格控制NFS共享目录的系统权限，尽量不用为普通用户赋予权限

## SAMBA
参考:
- [如何使用CIFS在Linux上挂载Windows共享](https://www.myfreax.com/how-to-mount-cifs-windows-share-on-linux/)
- [mount options with SMB share](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html-single/storage_administration_guide/index#ch-nfs)
- [SMB Mount Options](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html-single/storage_administration_guide/index#frequently_used_mount_options)
- [SMB on rhel 8](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html-single/deploying_different_types_of_servers/index#assembly_using-samba-as-a-server_Deploying-different-types-of-servers)
- [使用POSIX ACL控制Samba文件系统的访问](https://help.aliyun.com/document_detail/143007.html)

> 在rhel上，内核的cifs.ko文件系统模块提供了对SMB协议的支持. samba支持windows, mac, linux, 但linux推荐使用nfs.
> linux作为samba server实现多人分组共享, 只能使用acl. 步骤是: 1. 创建共享; 2. 组织用户 3. 清除acl, 再设置acl

```sh
$ sudo apt install samba samba-common smbclient cifs-utils # 安装samba
```

SMB 协议版本:
- SMB1：SMB1（也称为 CIFS）自 Windows NT 发布以来得到支持, 比如windows xp.
- SMB2：SMB2 自从 Windows Vista 发布以来得到支持，且为 SMB 的增强版本. SMB2 增加了将多重 SMB 操作功能组合到单个请求的功能，以减少网络数据包的数量并提高性能.
    SMB2 和 Large MTU：最大传输单元 (MTU) 是指可通过通讯协议的最大数据单元. 为利用最快的更快的接口，如 1- 或 10-gigabit 以太网，Large MTU 将最大传输单元提高至 1 megabyte (MB). 启用 Large MTU 可提高大文件传输的速度和效率，同时降低需处理的数据包数量.
- SMB3：SMB3 自 Windows 8 和 Windows Server 2012 发布以来得到支持, 它是 SMB 2 的增强版. SMB3 支持基于 AES 的文件加密传输，从而提高了对等文件传输的安全性.

> Windows Vista、Windows Server 2008 R2、Windows 7 和以上的版本支持 SMB2.
> Windows Server 2008 R2、Windows 7 和以上的版本支持 Large MTU.
> 确认 kernel 是否支持 CIFS 挂载：`grep -i cifs /boot/config-4.4.58-20180615.kylin.server.YUN+-generic`，y 或 m 表示支持即`CONFIG_CIFS=m`

### 组件
- smbd : 提供了文件和打印服务, 基于tcp.
- nmbd : 提供了NetBIOS名称服务和浏览支持，帮助SMB客户定位服务器，基于UDP. 它可以把linux系统共享的工作组名称和其ip对应起来, 否知就只能通过ip来访问共享文件.
- smbstatus ：列出目前 Samba 的联机状况， 包括每一条 Samba 联机的 PID, 分享的资源，使用的用户来源等等
- pdbedit : 管理用户数据

	- a : 用户名 建立 Samba 账户
    - -x : 用户名 删除 Samba 账户
    - -L : 列出账户列表
    - -Lv  : 列出账户详细信息的列表
- testparm : 检验配置文件 smb.conf 的语法正确与否
- smbclient : 查看其他计算机所分享出来的目录或打印机
- smbtree : 列出网络内其他计算机正在分享的内容, 类似于windows 网络邻居的显示效果.

> 在samba服务器端,权限由共享的目录的普通权限和smb.conf配置文件共同决定.
> SAMBA 使用的 NetBIOS 通讯协议
> SAMBA 仅只是 Linux 底下的一套软件，使用 SAMBA 来进行 Linux 文件系统时，还是需要以 Linux 系统下的 UID 与 GID 为准则. 也就是说，在 SAMBA 上面的使用者账号，必须要是 Linux 账号中的一个.

### 配置文件
- [/etc/samba/smb.conf](https://www.samba.org/samba/docs/current/man-html/smb.conf.5.html)

	samba的主要配置文件，基本上仅有这个文件，而且这个配置文件本身的说明非常详细. 主要的设置包括服务器全局设置，如工作组、NetBIOS名称和密码等级，以及共享目录的相关设置，如实际目录、共享资源名称和权限等两大部分

	```conf
	[global]
	server min protocol = SMB2 # 同`min protocol`, 也可指定具体版本[`server min protocol = SMB2_02`](https://www.samba.org/samba/docs/current/man-html/smb.conf.5.html#SERVERMAXPROTOCOL). [How to configure Samba to use SMBv2 and disable SMBv1 on Linux or Unix](https://www.cyberciti.biz/faq/how-to-configure-samba-to-use-smbv2-and-disable-smbv1-on-linux-or-unix/)
	client min protocol = SMB2
	client max protocol = SMB3
	load printers = yes # 是否加载打印机
	workgroup = WORKGROUP # 工作组，用来设定服务器所要加入的工作组或者域. 通常是配合windows使用的`WORKGROUP`.
	server string = Samba Server Version %v # 服务器简单介绍，%v显示samba版本号
	interfaces = lo eth0 192.168.12.2/24 192.168.13.2/24 # 服务器所监听的网卡名、IP地址
	hosts allow = 127. 192.168.12. 192.168.13. # 访问控制白名单，可以用一个IP表示，也可以用一个网段表示，多个参数以空格隔开
	log file = /var/log/samba/log.%m # 设置服务器日志文件的存储位置以及存储日志文件名称，%m 表示来访的主机名，即对每台访问服务器的机器都单独记录一个日志文件.
	log level = 3 # 0~10, 值越大越详细
	max log size = 5 # 定义日志文件的最大容量为 50KB
	security = user # 定义安全级别, 一共由四种级别：
	# - share：匿名共享，用户访问服务器不需要提供用户名和口令, 安全性差
	# - user：使用samba服务自我管理的帐号和密码进行用户认证，用户必须是本系统用户，但密码非/etc/shadow中的密码，而由samba自行管理的文件，其密码文件的格式由passdb bachend进行定义.
	# - server：由第三方服务进行统一认证
	# - domain：使用主域控制器进行认证，基于kerberos协议进行
	# - ADS: Active Directory Service, 是samba 3.0新增的身份验证方式
	passdb backend = tdbsam # 定义用户后台的类型，共有 3 种:
	# - smbpasswd：使用 smbpasswd 命令为系统用户设置 Samba 服务程序的密码. 使用smbpasswd命令来管理用户，要添加/管理的用户必须先是系统用户
	# - tdbsam使用一个数据库文件(`/var/lib/samba/private/passdb`)来建立用户数据库. 新版Samba的密码验证方式已使用tdbsam取代smbpasswd. 使用pdbedit命令来管理用户，要添加/管理的用户必须先是系统用户(**推荐**)
	# - ldapsam：基于 LDAP 服务进行账户验证
	load printers = yes #设置在 Samba 服务启动时是否共享打印机设备
	map to guest = bad user # 将samba sever所不能正确识别的用户都映射成guest用户
	guest account = user_name # samba默认将guest账户映射为nobody
	[josh] # 挂载时将使用的共享名称, 其相关的读写共权限与acl独立起作用
	comment = 共享的描述信息
    path = /samba/josh # 分享路径
    browseable = yes # 是否在“网上邻居”中可见
	writeable = true #该共享路径是否可写, read only的反义词
	write list = u1,u2 # 拥有写权限的用户列表（和writable不能同时使用）,会覆盖read only
    read only = no # 有效用户列表中指定的用户是否能够写入此共享
	read list = mary, @students  # 被授予对服务的只读访问权限的用户列表. 如果正在连接的用户在此列表中，则无论将`read only`选项设置什么，都将不授予他们写访问权限.
    force create mode = 0660 # 为此共享中新创建的文件设置权限, 会覆盖 create mode 设定的权限
    force directory mode = 2770 # 设置此共享中新创建的目录的权限
	force group = g1
	force user = u1 #  force group和force user强制规定创建的文件或文件夹的拥有者和组拥有者是谁. 一般这两个值来空，则表示拥有者和组拥有者为创建文件者.
    valid users = josh @sadmin # 允许访问共享的用户和组列表. 组以`@`为前缀, 其他所有用户都不能访问
    hosts allow = 192.168.115.0/24 127.0.0.1
    hosts deny = 0.0.0.0/0 # 当host deny 和hosts allow 字段同时出现并定义的内容相互冲突时，hosts allow 优先
	guest ok = no # 是否允许来宾帐号访问, 默认值为NO ，即设定在没有提交帐号和口令的情况下，是否允许访问此区段中定义的共享资源. 如同意guest帐号访问时，设为YES即是否允许匿名访问
	guest only = yes # 只允许用guest帐号访问
	public = yes # 是否允许匿名访问, 即是否"所有人可见", 这个开关有时候也叫guest ok，所以有的配置文件中出现guest ok = yes其实和public = yes是一样的
	invalid users = root # 设定不允许访问此共享资源的用户或组
    sync always = no # 写操作后是否立即进行sync 
	```

	在smb.conf中<section header>中有三个特殊的NAME，分别是global、homes和printers:
	- [global]：其属性选项是全局可见的，但是在需要的时候，我们可以在其他<section>中定义某些属性来覆盖[global]的对应选项定义.
	- [homes]：当客户端发起访问共享服务请求时，samba服务器就查询smb.conf文件是否定义了该共享服务，如果没有指定的共享服务<section>，但smb.conf文件定义了[homes]时，samba服务器会将请求的共享服务名看做是某个用户的用户名，并在本地的password文件中查询该用户，若用户名存在并且密码正确，则samba服务器会将[homes]这个<section>中的选项定义克隆出一个共享服务给客户端，该共享的名称是用户的用户名.
	- [printers]：用于提供打印服务. 当客户端发起访问共享服务请求时，没有特定的服务与之对应，并且[homes]也没有找到存在的用户，则samba服务器将把请求的共享服务名当做一个打印机的名称来进行处理.

	example:
	```ini
	# 参考: [Samba共享目录的多用户权限设置案例](https://www.cnblogs.com/kevingrace/p/5569993.html)
	[exchage] # 所有人都能读写，包括guest用户，但每个人不能删除别人的文件
	comment = Exchange File Directory
	path = /home/samba/exchange # 再加`chmod -R 1777 /home/samba/exchange`
	public = yes
	writable = yes
	[public] # 所有人只读这个文件夹的内容
	comment = Read Only Public
	path = /home/samba/public
	public = yes
	read only = yes
	[caiwu] # caiwu组和lingdao组的人能看到，network02也可以访问，但只有caiwu01有写的权限
	comment = caiwu
	path = /home/samba/caiwu
	public = no
	valid users = @caiwu,@lingdao,network02
	write list = caiwu01
	printable = no

	[lingdao] # 只有领导组的人可以访问并读写，还有network02也可以访问，但外人看不到那个目录
	comment = lingdao
	path = /home/samba/lingdao
	public = no
	browseable = no
	valid users = @lingdao,network02
	printable = no
	```

- /var/lib/samba/private/{passdb.tdb,secrets.tdb} 

	管理 Samba 的用户账号/密码时，会用到的数据库档案

### 使用
```sh
$  testparm -s # 检查smb.conf是否正确
$ smbclient -L //127.0.0.1 [-U josh]# 列出正在分享的内容
$ smbclient //192.168.0.141/{samba_share_name} # 默认以当前用户和字符界面模式访问samba_share_name
$ smbclient --user=share //192.168.66.198/share # 访问共享
$ sudo useradd -M -s /usr/sbin/nologin -G sambashare josh
# $ sudo smbpasswd -a josh # 设置用户密码将sadmin用户帐户添加到Samba数据库, 默认已启用账号. 可用`pdbedit -a -u ${user}`代替
# $ yes password |sudo smbpasswd -a ubuntu # 不用交互输入密码
# $ sudo smbpasswd -e josh # 启用账号josh
$ pdbedit -a username    #新建Samba账户, **username必须已存在**
$ pdbedit -x username    #删除Samba账户
$ pdbedit -v username    #显示账户详细信息
$ sudo pdbedit -L -v # 查看smbpasswd创建的samba用户
$ sudo systemctl restart smbd # 使**配置生效**
# smbcontrol all reload-config # 重新加载Samba配置, 使授权生效, **推荐**
$ sudo mount -t cifs //127.0.0.1/{samba_share_name} /mnt [-o username=josh -o password=xxx -o vers=2.0  -o uid=$(id -u),gid=$(id -g) ] # 挂载samba分享的内容, client端支持的smb protocol 版本可通过`man mount.cifs#vers查看`. samba使用samba_share_name, 而不像nfs那样的export路径. 未登录用户(密码登录)映射为nobody:nogroup, 否则用指定的username:username.
$ sudo mount | grep cifs # 挂载的详细参数, 可参考[通过云服务器ECS（Linux）访问SMB文件系统#挂载文件系统](https://www.alibabacloud.com/help/zh/doc-detail/128737.htm)
$ sudo smbstatus # 查看连接到samba server的client及使用的protocol version + samba server version, 映射的用户及用户组. version显示`Unknown`: 客户端支持的smb协议比smbd新.
```

注意(smbd v4.3.11):
- **刷新前已挂载目录(可写)在刷新后(已去掉可写)权限不变(仍可写)**, 使用`service samba restart/smbcontrol all reload-config`刷新也是同样的结果. 需client重新挂载生效.
- client 挂载时原有的ownner是root, server端`chown -R nobody:nogroup ${samba_share_name}`, client端的挂载仍是root,  `service samba restart/smbcontrol all reload-config`刷新后仍不变. 需client重新挂载生效.

on windows:
1. `win + R`, 输入`\\{samba_server_ip}`
1. 输入设置的samba账号, 进入共享目录
或`net use z: \\xxx-shanghai.nas.aliyuncs.com\myshare [用户名密码 /user:管理员权限的用户名]` #linux/windows未登录挂载时用户会被映射为`nobody:nogroup`; 登录挂载时因为samba登录没有组的概念, 因此用户会被映射为`username:username`(可通过smbstatus查看); 如果samba server是linux, 那么它还会带上支持组的权限; 新建文件归属于映射到的用户.

执行`net use`命令，检查挂载结果

> 清除windows网络邻居的连接(默认只能连接一个共享): `net use * /del /y`

on linux:
`/etc/fstab`:
```
//192.168.0.10/gacanepa /mnt/samba  cifs credentials=/root/smbcredentials,defaults 0 0
# smbcredentials:
# username=gacanepa
# password=XXXXXX
```

## FAQ
### 注意点
不能以 NFS 和 SMB 同时挂载同一个文件系统.

建议不要使用 Linux 作为客户端访问 SMB，因为存在一些操作上的问题. 例如支持的字符集、文件名的长度（Windows 支持255宽字符，Linux 支持255 UTF8 字节）等等.

但用户如果确实需要的话，可以在支持 SMB2 及以上的 kernel 上挂载.

### samba挂载乱码
根源: 支持的字符集不同.

解决方案:
- windows is server, linux mount

	sudo mount -t smbfs -o username=guest,codepage=cp936,iocharset=utf8 //192.168.0.38/movie /mnt/smb/ # 未测试

	参考:
	- [关于mount/samba/字符集的两篇好文](https://blog.zengrong.net/post/1019/)
- linux is server, windows mount

	# 未测试, 也不推荐修改linux的字符集, 这样可能会在linux上出现其他问题, 比如linux开始出现乱码.
	# locale -a | grep zh # 查找支持gbk字符集
	# export LANG="zh_CN.gb18030"
	#export LC_ALL="zh_CN.gb18030"


 **推荐windows和linux不挂载同一个samba/nfs共享**

### [windows挂载nfs的中文乱码问题的解决](https://support.huawei.com/enterprise/zh/knowledge/EKB1100039367)
因是windows内置的nfs挂载工具所支持的文字编码太有限了，不支持utf-8. 

使用第三方nfs 客户端，比如 ms-nfs41-client 软件.

### samba client挂载显示的ownner与server上的权限不一致
参考[文件和目录的属主及权限 from `man mount.cifs`](http://www.jinbuguo.com/man/mount.cifs.html).

**核心 CIFS 协议并不提供文件和目录的 unix 属主或权限信息, 而是采用了windows权限模型**。
正因为如此，文件和目录才会看上去像被 uid= 和 gid= 选项指定的用户和组所拥有，
并且其权限才会看上去和 file_mode 以及 dir_mode 指定的权限一致。
可以通过 chmod/chown 来修改这些值，但是并不会在服务器端产生真正的实际效果。

> windows权限不受影响.

如果服务器端支持Unix扩展，并且客户端也允许使用Unix扩展，文件和目录的 uid, gid, mode 将由服务器端提供。
因为 CIFS 通常由同一个用户挂载，所以不管是哪个用户访问此文件系统，所使用的 credentials 文件都是同一个。
这样，新创建的文件和目录其属主/属组就都根据同一个 credentials 文件中的连接用户来设置。

如果客户端和服务器端使用的 uid 和 gid 并不匹配，那么 forceuid 和 forcegid 选项就很有用处了。
注意，并没有强制改写 mode 的选项。
当指定了 forceuid 和/或 forcegid 后，文件和目录的权限就可能不能反映真正的权限了。

如果Unix扩展被nounix禁用(或者服务器端本身就不支持)，仍然有可能使用"dynperm"选项在服务器上模拟出来。
使用该选项后，新创建的文件和目录将看上去拥有了正确的权限。
不过这些权限并不真正存储在服务器端的文件系统上(仅在内存中)，因此可能会随时丢失(比如内核刷新了inode缓存)。
因此，我们不鼓励使用此选项。

还可以使用 noperm 选项在客户端完全越过权限检查。
但是**服务器端的权限检查是无法越过的**，服务器端将始终根据 credentials 文件中提供的用户信息进行权限检查，
而与客户端实际访问文件系统的用户无关。

[unix extensions仅支持smb 1](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_file_systems/assembly_mounting-an-smb-share-on-red-hat-enterprise-linux_managing-file-systems#con_unix-extensions-support_assembly_mounting-an-smb-share-on-red-hat-enterprise-linux), 因为安全问题已在之后的smb protocol中被禁用了.

因此**smb的权限是由client mounted显示的权限, 登录账户, server端权限**共同作用,推荐挂载时使用`-o uid=$(id -u),gid=$(id -g)`选项(默认是挂载者的uid/gid), 或samba**仅支持windows共享**.

> [SMB 1.0 由于协议设计的巨大差异导致在性能和功能上有严重的不足，并且只支持 SMB1.0 或更早协议版本的 Windows 产品已经完全退出微软支持的生命周期](https://www.alibabacloud.com/help/zh/doc-detail/122195.htm)

### wrong fs type, bad option, bad superblock on
`是/sbin/mount下面缺少挂载nfs格式的文件，应该是mount.nfs[xxx]，而该文件由nfs-common提供，所以需要nfs-common工具`,解决方案:
```
# apt install nfs-common
# yum install nfs-utils
```

### `service nfs-kernel-server start`报 Not starting NFS kernel daemon: no exports
`/etc/exports`为空导致.

### mount.nfs: timeout
通常是网络问题, ping一下网络.

### mount nfs on windows
```sh
> showmount -e 192.168.0.83
> mount 192.168.0.83:/usr/local/xxx/mypool/f14 Z:
```

### zfs nfs mount no acl
1. 设置acls后nfs挂载,创建目录成功, umount, 再使用samba挂载,取消挂载并重启smbd, 最后回到nfs挂载, 无法创建目录(报无权限), `getfacl`无法获取acls, 但`nfs4_getfacl`可以, 刷新nfs server(`exportfs -r`)后恢复, 但`getfacl`无法获取acls.

### [是否支持 NFS 和SMB 同时挂载一个文件系统](https://www.alibabacloud.com/help/zh/doc-detail/110839.htm)
不能以 NFS 和 SMB 同时挂载同一个文件系统.

建议不要使用 Linux 作为客户端访问 SMB，因为存在一些操作上的问题, 例如支持的字符集、文件名的长度（Windows 支持255宽字符，Linux 支持255 UTF8 字节）等等. 但如果确实需要的话，可以在支持 SMB2 和kernel 3.10.0-514 及以上的系统上挂载.

### `smbstatus`显示client的protocol version 是 Unknown (0x0311)
"Unknown (0x0311)" protocol is fixed in [`Samba 4.4.0`](https://bugzilla.samba.org/show_bug.cgi?id=11472).

### samba 无法创建文件???
env:
```
$ Linux 5.3.0-24-generic
$ Samba version 4.3.11-Ubuntu
$ mount.cifs version: 6.9
```
明明有写权限, 还是无法创建文件, windows server 2012和Linux 4.4.131-20190505.kylin.server-generic + mount.cifs version: 6.4则正常.

将mount.cifs version: 6.9降到6.4还是报同样的错.

### 禁用nfsv2, nfsv3, udp
```bash
# vim /etc/default/nfs-kernel-server
RPCNFSDOPTS="-N 2 -N 3 -U"
# systemctl restart nfs-kernel-server
$ sudo  mount -t nfs -o vers=2 192.168.0.141:/mnt/xfs nfs_xfs 
mount.nfs: Protocol not supported
```

> rhel7是`/etc/sysconfig/nfs`的`RPCNFSDARGS`

> RPCNFSDCOUNT 是 nfsd使用线程数

### `mount: can't find nfs in /etc/fstab`
```
$ mount -t nfs4 -o 192.168.0.141:/mnt/xfs nfs 
mount: can't find nfs in /etc/fstab
```

删除多余的`-o`

### refused mount request from 192.168.0.121 for /mnt/xfs (/mnt/xfs): unmatched host
参考:
- [NFS Mount over a Specific Interface](https://www.redhat.com/archives/fedora-list/2005-September/msg03442.html)

```bash
# mount -t nfs -o vers=3,clientaddr=192.160.0.31  192.168.0.141:/mnt/xfs nfs # 报错:`unmatched host`. 192.168.0.121与192.160.0.31是同一台电脑.
# ### nsf server: 0.141
# tcpdump -i bond0 src host 192.168.0.121
# ### nfs client
# tcpdump -i eth0 src host 192.168.0.121  and dst port 2049 # 2049是nfs server使用的端口
# tcpdump -i eth0 src host 192.168.0.121 and dst host 192.168.0.141 # 和上面的作用一样: 判断链路情况
```

通过tcpdump发现, 即使指定了clientaddr, 但mount.nfs还是使用了192.168.0.121.

> `unmatched host`仅在第一次请求时输出, 重复请求不输出, 此时重启nfs server后又可看到该错误, 预计nfsd的其他错误也会有这种情况.

> clientaddr在`man 5 nfs`

原因: 在mount.nfs代码中未发现bind指定ip的操作:
```c
// git clone -b ubuntu/trusty  https://git.launchpad.net/ubuntu/+source/nfs-utils # nfs-common 1:1.2.8-6ubuntu1.2
// nfs-utils/utils/mount/nfsmount.c : 787~817
// 最新版1.3.4-3也是这样
if (nfs_mount_data_version == 1) {
		/* create nfs socket for kernel */
		if (nfs_pmap->pm_prot == IPPROTO_TCP)
			fsock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
		else
			fsock = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);
		if (fsock < 0) {
			perror(_("nfs socket"));
			goto fail;
		}
		if (bindresvport(fsock, 0) < 0) {
			perror(_("nfs bindresvport"));
			goto fail;
		}
	}

#ifdef NFS_MOUNT_DEBUG
	printf(_("using port %lu for nfs deamon\n"), nfs_pmap->pm_port);
#endif
	nfs_saddr->sin_port = htons(nfs_pmap->pm_port);
	/*
	 * connect() the socket for kernels 1.3.10 and below only,
	 * to avoid problems with multihomed hosts.
	 * --Swen
	 */
	if (linux_version_code() <= MAKE_VERSION(1, 3, 10) && fsock != -1
	    && connect(fsock, (struct sockaddr *) nfs_saddr,
		       sizeof (*nfs_saddr)) < 0) {
		perror(_("nfs connect"));
		goto fail;
	}
```


### nfs debug
```bash
# rpcdebug -vh
# rpcdebug -m nfs -s all # Enable all NFS (client-side) debugging
# rpcdebug -m rpc -s call # only Enable RPC Call (client-side) debugging
# rpcdebug -m nfsd -s all # Enable NFSD (server-side) debugging
# ### Disable debugging
# rpcdebug -m nfs -c all
# rpcdebug -m nfsd -c all
```

rpcdebug module:
- nfs 	NFS client
- nfsd 	NFS server
- nlm 	Network Lock Manager Protocol(NLM)
- rpc 	Remote Procedure Call

rpcdebug选项:
- -m : module name to set or clear kernel debug flags
- -s : To set available kernel debug flag for a module
- -c : Clear Kernel debug flags

> 将nfsd日志输入syslog: `RPCNFSDOPTS="-d -s"`

### zfs xfs nas
env: 5.3.0-26-generic/4.4

> 在zfs fs (on 0.7.x)上直接使用acl容易出现莫名奇妙的问题, 且[zfs 还未支持NFSv4 ACL](https://github.com/openzfs/zfs/pull/9709). 当前思路是使用zfs vol+格式化作为磁盘, 在其上再设置nas, 整个共享使用一个账户, 再将客户端的用户加入对应的组即可.

> 读写权限 : 允许授权对象对文件系统进行只读或读写.

> nfs和smb不允许重合使用, 避免未知问题.

nfs:
```bash
# grep -i CONFIG_XFS_FS /boot/config-5.3.0-26-generic #  check kernel support xfs
# modinfo xfs # check kernel support xfs when CONFIG_XFS_FS=m
# modprobe xfs # kernel load xfs module
# lsmod |grep -i xfs # check xfs mod is loaded
# cat /proc/filesystems |grep -i xfs # check kernel support xfs

# dpkg -l |grep -i xfs # check packages for xfs 
# apt-get install xfsprogs

# grep -i acl /boot/config* check kernel support for POSIX_ACL, like: CONFIG_EXT4_FS_POSIX_ACL, CONFIG_XFS_POSIX_ACL
# grep -i nfs /boot/config* check kernel support for NFSv4. like: CONFIG_NFS_V4_1, CONFIG_NFS_V4_2

# sudo zfs create -V 5gb x/vol_xfs # vol@/dev/zvol/x/vol_xfs
# mkfs -t xfs /dev/zvol/x/vol_xfs
#  mkdir /mnt/xfs
# mount /dev/zvol/x/vol_xfs /mnt/xfs
# chown -R nobody: nogroup /mnt/xfs
# chmod 777 /mnt/xfs
# vim /etc/exports
/mnt/xfs 192.168.0.245(rw,all_squash,no_subtree_check,async)
/mnt/xfs 192.168.0.131(ro,all_squash,no_subtree_check,async)
# exportfs -ra

## on client @ 192.168.0.245
# gpasswd -a  ${USER} nogroup # 将当前用户加入nogroup
# id # 查看是否已加入nogroup
# mount -t nfs -o vers=4,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 192.168.0.141:/mnt/xfs nfs_xfs
# cd nfs_xfs
# touch a # is ok, 但有时第一次操作会卡几秒~几十秒钟
## on client @ 192.168.0.131
# gpasswd -a  ${USER} nogroup
# mount -t nfs -o vers=4,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 192.168.0.141:/mnt/xfs nfs_xfs
# cd nfs_xfs
# touch b
touch: cannot touch 'c': Read-only file system # is ok, because exported with ro
```

smb:
```bash
# sudo zfs create -V 5gb x/vol_smb # vol@/dev/zvol/x/vol_smb
# mkfs -t xfs /dev/zvol/x/vol_smb
#  mkdir /mnt/smb
# mount /dev/zvol/x/vol_smb /mnt/smb
# chown -R nobody: nogroup /mnt/smb
# chmod 777 /mnt/smb
# vim /etc/samba/smb.conf
/mnt/smb 192.168.0.245(rw,all_squash,no_subtree_check,async)
/mnt/smb 192.168.0.131(ro,all_squash,no_subtree_check,async)
# smbcontrol all reload-config

## on client @ 192.168.0.245
# gpasswd -a  ${USER} nogroup # 将当前用户加入nogroup
# id # 查看是否已加入nogroup
# mount -t nfs -o vers=4,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 192.168.0.141:/mnt/smb nfs_xfs
# cd nfs_xfs
# touch a # is ok, 但有时第一次操作会卡几秒~几十秒钟
## on client @ 192.168.0.131
# gpasswd -a  ${USER} nogroup
# mount -t nfs -o vers=4,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 192.168.0.141:/mnt/smb nfs_xfs
# cd nfs_xfs
# touch b
touch: cannot touch 'c': Read-only file system # is ok, because exported with ro
```