# fs
## CIFS, SMB, NFS
SMB(Server Message Block，即服务(器)消息块) 是 IBM 公司在 80年代中期发明的一种文件共享协议. 它只是系统之间通信的一种方式（协议）. 目前最新版是`v3.1.1`.
CIFS是微软的Common Internet file system的缩写, 是 SMB 协议的一种特殊实现, 不常用.
Samba 也是 SMB 协议的实现, 常用于windows与类unix间的文件共享.
NFS是SUN为Unix开发的网络文件系统, 提供类unix间的文件共享. 目前最新版本是`v4.2`. NFSv4用户验证采用“用户名+域名”的模式，与Windows AD验证方式类似，NFSv4强制使用Kerberos验证方式.（Kerberos与Windows AD都遵循相同RFC1510标准），这样方便windows和`*nix`环境混合部署.

### NFS
安装:
```
$ sudo apt install nfs-kernel-server
$ sudo cat /proc/fs/nfsd/versions
$ sudo apt install nfs-common # Install NFS client 
$ sudo yum install nfs-utils # Install NFS client
$ nfsstat -s # server使用的nfs version
$ nfsstat -c # client使用的nfs version
$ nfsstat -m # 在client端已挂载的nfs信息
$ nfsstat -4 # 查看NFS版本4的状态
$ sudo systemctl status nfs-kernel-server
$ showmount -e 192.168.0.83 # 在 Client 端查看server端(192.168.0.83)共享出来的目录
$ sudo mount -t nfs -o vers=4.2 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享, 挂载成功后不能访问请检查nfs server端的权限
$ sudo mount -t nfs4 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享
$ sudo mount -o v4.2 192.168.0.83:/usr/local/mypool/p11 /mnt # 用指定版本的nfs挂载共享
$ df -h #查看挂载情况
$ sudo umount /mnt
```

NFS server 的配置选项在 /etc/default/nfs-kernel-server 和 /etc/default/nfs-common 里.
`/etc/exports`是用来管理NFS共享目录的使用权限与安全设置的地方. 特别注意的是，NFS本身设置的是网络共享权限，整个共享目录的权限还和目录自身的系统权限有关.
/var/lib/nfs/etab                      记录NFS共享出来的目录的完整权限设定值
/var/lib/nfs/xtab                      记录曾经登录过的客户端信息

>　nfs server指定使用的版本: `/etc/default/nfs-kernel-server`的`RPCMOUNTDOPTS="--manage-gids -V 4.2"`.

#### FS系统守护进程
- nfsd ：它是基本的NFS守护进程，主要功能是管理客户端是否能够登录服务器
- rpc.mountd ：它是RPC安装守护进程，主要功能是管理NFS的文件系统. 当客户端顺利通过nfsd登录NFS服务器后，在使用NFS服务所提供的文件前，还必须通过文件使用权限的验证. 它会读取NFS的配置文件/etc/exports来对比客户端权限.
- lockd : 用在管理档案的锁定 (lock) 用途. 当多个客户端同时尝试写入某个档案时， 需要lockd 来解决多客户端同时写入的问题. 但 lockd 必须要同时在客户端与服务器端都开启才行. 此外， lockd 也常与 rpc.statd 同时启用.
- statd : 检查文件的一致性，与lockd有关. 若发生因为客户端同时使用同一档案造成档案可能有所损毁时， statd 可以用来检测并尝试恢复该档案. 与 lockd 同样的，这个功能必须要在服务器端与客户端都启动才会生效.
- rpc.idmapd : 名字映射后台进程
- rpcbind : 主要功能是进行端口映射工作. 当客户端尝试连接并使用RPC服务器提供的服务（如NFS服务）时，rpcbind会将所管理的与服务对应的端口提供给客户端，从而使客户可以通过该端口向服务器请求服务, 因此rpcbind必须在nfs前启动.

#### 其他相关命令

1. exportfs

使/etc/exports的配置立刻生效，该命令格式如下：

　　# exportfs [-aruv]

　　-a 全部挂载或卸载 /etc/exports中的内容
　　-r 重新读取/etc/exports 中的信息 ，并同步更新/etc/exports、/var/lib/nfs/xtab
　　-u 卸载单一目录（和-a一起使用为卸载所有/etc/exports文件中的目录）
　　-v 在export的时候，将详细的信息输出到屏幕上。

具体例子：
　　# exportfs -au 卸载所有共享目录
　　# exportfs -rv 重新加载共享所有目录并输出详细信息

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

#### /etc/exports
格式：`<输出目录> [客户端1(选项: 访问权限,用户映射,其他)] [客户端2(选项: 访问权限,用户映射,其他)] ...`
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

	    all_squash：将远程访问的所有普通用户及所属组都映射为匿名用户或用户组（nfsnobody）, 可由anonuid/anongid指定
	    no_all_squash：与all_squash取反（默认设置）
	    root_squash：将root用户及所属组都映射为匿名用户或用户组（默认设置: 客户端 root 的身份会由 root_squash 的设定压缩成 nfsnobody， 如此对服务器的系统会较有保障）
	    no_root_squash：与rootsquash取反, 允许使用 root 身份来操作服务器的文件系统
	    anonuid=xxx：将远程访问的所有用户都映射为匿名用户，并指定该用户为本地用户（UID=xxx, 该 UID 必需要存在于你的 /etc/passwd 当中）
	    anongid=xxx：将远程访问的所有用户组都映射为匿名用户组账户，并指定该匿名用户组账户为本地用户组账户（GID=xxx）

	1. 其它选项

	    secure：限制客户端只能从小于1024的tcp/ip端口连接nfs服务器（默认设置）
	    insecure：允许客户端从大于1024的tcp/ip端口连接服务器
	    sync：将数据同步写入内存缓冲区与磁盘中，效率低，但可以保证数据的一致性
	    async：将数据先保存在内存缓冲区中，必要时才写入磁盘
	    wdelay：检查是否有相关的写操作，如果有则将这些写操作一起执行，这样可以提高效率（默认设置）
	    no_wdelay：若有写操作则立即执行，应与sync配合使用
	    subtree：若输出目录是一个子目录，则nfs服务器将检查其父目录的权限(默认设置)
	    no_subtree：即使输出目录是一个子目录，nfs服务器也不检查其父目录的权限，这样可以提高效率

### FAQ
#### wrong fs type, bad option, bad superblock on
`是/sbin/mount下面缺少挂载nfs格式的文件，应该是mount.nfs[xxx]，而该文件由nfs-common提供，所以需要nfs-common工具`,解决方案:
```
# apt install nfs-common
# yum install nfs-utils
```