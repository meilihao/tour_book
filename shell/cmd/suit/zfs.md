# zfs
参考:
- [在 Oracle® Solaris 11.2 中管理 zfs 文件系统](https://docs.oracle.com/cd/E56344_01/html/E53918/index.html)
- [zfs Administration](https://pthree.org/2012/12/04/zfs-administration-part-i-vdevs/)
- [FreeBSD Handbook's zfs](https://www.freebsd.org/doc/handbook/)
- [ZFS磁盘空间管理(分配、释放原理)](https://blog.csdn.net/beiya123/article/details/80393720)
- [man pages](https://zfsonlinux.org/manpages/0.8.4/index.html)
- [Btrfs vs ZFS 实现 snapshot 的差异](https://farseerfc.me/zhs/btrfs-vs-zfs-difference-in-implementing-snapshots.html)
- [zfs使用jemalloc(未验证)](https://gitlab.openebs100.io/openebs/cstor/commit/043de97885a79533d1ee262a9e6d23cb9f604c1c)
- [OpenZFS开源文件系统2.0+：持久化L2ARC读缓存、ZIL写缓存提速](https://zhuanlan.zhihu.com/p/338227098)
- [ZFS──瑞士军刀般的文件系统](https://www.eaimty.com/2020/02/zfs-file-system.html)
- [ZFS 分层架构设计](https://farseerfc.me/zhs/zfs-layered-architecture-design.html)
- [zfs share(nas)](https://wiki.debian.org/ZFS#File_Sharing)
- [在Linux上安装和使用ZFS](https://www.escapelife.site/posts/caf259ea.html)
- [Zfs_ondiskformat.pdf](http://www.giis.co.in/Zfs_ondiskformat.pdf)
- [ZFS IOPS limit](https://dokuwiki.fl8.jp/01_linux/13_storage/31_zfs_iops_limit)

	cgcreate + cgset
- [zfs空间估算 : ZFS / RAIDZ Capacity Calculator (beta)](https://wintelguy.com/zfs-calc.pl)

```sh
# ubuntu 18.04
$ sudo apt install zfsutils-linux # 安装zfs
# centos看这里https://github.com/zfsonlinux/zfs/wiki/RHEL-and-CentOS, 可能需先安装DKMS. 可能需要重启.
# 以下三种方式获取zfs version
$ modinfo zfs
$ cat /sys/module/zfs/version
$ dmesg | grep -i zfs
```

zfs有两个工具: zpool和zfs. zpool用于维护zfs pools, zfs负责维护zfs filesystems.

> zfs可充当卷管理器
> zfs已实现零管理, 通过zfs daemon(zed,zfs Event Daemon, zfs事件守护进程, 将监听任何zfs生成的kernel event, 配置在`/etc/zfs/zed.d/zed.rc`)实现, 因此无需将pool信息写入`/etc/fstab`, pool配置在`/etc/zfs/zpool.cache`里(可通过`zpool get "all" <pool>`找到`zpool.cache`的位置).
> zfs pool使用zfs_vdev_scheduler来调度io.
> zdb 是zpool的调试工具.
> zgenhostid : generate and store a hostid in /etc/hostid
> zvol_wait :  导入pool时, 等待/dev/zvol下的所有符号链接都创建后返回

> [zfs-2.1.10 : Removed Python 2 and Python 3.5- support](https://github.com/openzfs/zfs/releases/tag/zfs-2.1.10)

## 概念
参考:
- [ZFS 术语](https://docs.oracle.com/cd/E26926_01/html/E25826/ftyue.html)
- [zpoolconcepts.7](https://openzfs.github.io/openzfs-docs/man/master/7/zpoolconcepts.7.html)

pool : 存储设备的逻辑分组, 它是ZFS的基本构建块，可将其存储空间分配给数据集.
dataset : zfs文件系统的组件即文件系统、克隆、快照和卷被统称为数据集
mirror : 一个虚拟设备存储相同的两个或两个以上的磁盘上的数据副本，在一个磁盘失败的情况下,相同的数据是可以用其他磁盘上的镜像.
resilvering ：在恢复设备时将数据从一个磁盘复制到另一个磁盘的过程
snapshot : 快照, 是文件系统或卷的只读副本. 在zfs中，快照几乎可以即时创建，而且最初不会额外占用池中的磁盘空间
scrub : 用于一致性检验. 其他文件系统会使用fsck.
thin: zfs支持thin provisioning

### version
ref:
- [ZFS RAIDz expansion 抢先体验](https://ntzyz.space/zh-cn/post/testing-zfs-raidz-expansion/)
- [ZFS新特性：引入RAID-Z Expansion功能](https://www.163.com/dy/article/GVSNHHR00511CUMI.html)

	目前适用于 RAIDZ-1/2/3

- [2.3](https://github.com/openzfs/zfs/releases/tag/zfs-2.3.0)

	- [RAIDZ Expansion](https://github.com/openzfs/zfs/pull/15022): 向raidz添加新盘以扩容空间, 但不改变raid level

		`zpool attach p raidz1-0 /dev/vdf`, 一次仅能attach一个新盘, 且不能移除该盘, 除非destroy整个pool. 不同raidz间可以并发attach

		扩展过程是动态的，即使在扩展时，数据访问也不会受到干扰，但整个扩展过程需要一些时间，无法即刻获得新空间. 同时[它比正常扩展raid组来说会损失部分空间](https://louwrentius.com/zfs-raidz-expansion-is-awesome-but-has-a-small-caveat.html).
	- Fast Dedup: 更快的dedup

## 优化
- [Running PostgreSQL using ZFS and AWS EBS](https://bun.uptrace.dev/postgres/tuning-zfs-aws-ebs.html)

- zfs事务间隔
- dataset sync属性

## reg
- [dataset name](https://openzfs.github.io/openzfs-docs/man/8/zfs.8.html)

### 建议
1. 使用整个磁盘来作为磁盘写入高速缓存(ZIL)并使维护更轻松
1. 使用zfs冗余(raidz,mirror)

   1. mirror, 使用镜像磁盘对
   1. raidz, 为每个vdev组合3-9个磁盘(不包括校验盘)
   1. 不要在同一个pool中混用raidz和mirror, 因为会难以管理且性能受到影响
1. 避免在存储池中使用磁盘分片而是应使用整块磁盘, 已避免潜在复杂性.
1. 使用热备以减少硬件故障而导致的停机时间
1. 使用相同大小的磁盘以便在各个设备之间平衡 I/O
1. pool最小为8G.
1. [重删和压缩比较消耗资源, 特别是重删, 要求20G RAM/T, 我们(共40G RAM)碰到过原先200M/s的速度开了重删和压缩后变成20M/s](https://constantin.glez.de/2011/07/27/zfs-to-dedupe-or-not-dedupe/).
1. 使用超过80%后建议扩容否则可能遭遇性能问题
1. 任一fs均可使用全部的可用空间, 限制方法是设置reservation(预留) and quota(配额).
1. 压缩推荐使用lz4(速度)和zstd(压缩比).

	```bash
	# zpool set feature@lz4_compress=enabled data
	# zfs set compression=lz4 data/datafs
	# zfs get compressratio data/datafs
	```
1. [zfs内存需求](https://www.truenas.com/docs/scale/introduction/scalehardwareguide/#memory-sizing)

### zfs虚拟设备(zfs vdevs)
参考:
- [ZFS高速缓存介绍：ZIL和L2ARC](https://www.xiangzhiren.com/archives/288)
- [OpenZFS – Understanding ZFS vdev Types](https://klarasystems.com/articles/openzfs-understanding-zfs-vdev-types/)

一个 VDEV 是一个meta-device，代表着一个或多个物理设备. zfs 支持 7 中不同类型的 VDEV：
- disk, 默认, 比如HDD, SDD, PCIe NVME等等
- File : 预先分类的文件，为*.img的文件，可以作为一个虚拟设备载入zfs
- Mirror : 标准的 RAID1 mirror
- zfs 软件RAID : raidz=raidz1(raid5, >=3 disk)/2(raid6, >=4 disk)/3(raid7, >=5 disk), 非标准的基于分布式奇偶校验的软件RAID.

	性能对比

	Stripe > Mirror
	Stripe > RAIDZ1 > RAIDZ2 > RAIDZ3
	
	数据可靠性:

	Mirror > Stripe
	RAIDZ3 > RAIDZ2 > RAIDZ1 > Stripe

	推荐:
	- raidz: 3<=2n+1<=11
	- raidz2: 4<=2n+2<=12
- Hot Spare : 用于**热备** zfs 的软件 raid, 当正在使用的磁盘发生故障后，Spare磁盘将马上代替此故障盘.
- Cache : 用于2级自适应的**读缓存**的设备 (zfs L2ARC), 提供在 memory 和 disk的缓冲层, 用于改善静态数据的随机读写性能

	zfs实际有两层缓存: ARC(在内存) + L2ARC(高速存储设备上, 比如Flash，SSD，Nvme等)

	为 ZFS 池配置 L2ARC 缓存后，ZFS 将从 ARC 缓存中删除的数据存储在 L2ARC 缓存中。因此，可以将更多数据保存在缓存中以加快访问速度.

	ARC内存占用可用[`arc_summary.py`](https://github.com/openzfs/zfs/tree/master/cmd/arc_summary)获取, 或使用arcstat命令
- Log : zfs Intent Log(zfs意图日志), 是记录两次完整事务语义提交之间的日志，用来加速实现 fsync 之类的文件事务语义, 是一种对于 data 和 metadata 的日志机制，先写入然后再刷新为写事务), 持久性写缓存, 用于崩溃恢复, 最好配置并使用快速的 SSD来存储ZIL, 以获得更佳性能. ZIL支持mirror. ZIL也可认为是zfs的**写缓存**.

	在zfs中, ZIL是一种机制，而SLOG则是一种设备，也就是说SLOG是可选的并不是必须的，但是ZIL则是默认的。在存储池中如果没有单独设置SLOG设备时，ZIL机制也是存在的, 会在存储池中划分一部分空间出来进行处理.

	SLOG 之于 ZIL 有点像 L2ARC 之余 ARC ， L2ARC 是把内存中的 ARC 放入额外的高速存储设备，而 **SLOG 是把原本和别的数据块存储在一起的 ZIL 放到额外的高速存储设备**.

	**默认情况下，ZFS 会分配一小部分池用于存储写入缓存**。它被称为ZIL或ZFS 意图日志。在将数据写入物理硬盘之前，它会存储在 ZIL 中。为了最大限度地减少写入操作的数量并减少数据碎片，数据在 ZIL 中进行分组，并在达到某个阈值时刷新到物理硬盘驱动器。它更像是一个写缓冲区而不是缓存

	由于 ZFS 使用池的一小部分来存储 ZIL，因此它共享 ZFS 池的带宽。这可能会对 ZFS 池的性能产生负面影响。

	要解决此问题，您可以使用快速 SSD 作为 SLOG 设备。如果 ZFS 池上存在 SLOG 设备，则将 ZIL 移动到 SLOG 设备。ZFS 不会再将 ZIL 数据存储在池中。因此，不会在 ZIL 上浪费任何池带宽。

	还有其他好处。如果应用程序通过网络（即 VMware ESXi、NFS）写入 ZFS 池，ZFS 可以快速将数据写入 SLOG 并向应用程序发送数据已写入磁盘的确认。然后，它可以像往常一样将数据写入速度较慢的硬盘。这将使这些应用程序更具响应性.

	请注意，通常情况下，ZFS 不会从 SLOG 中读取。ZFS 仅在断电或写入失败的情况下从 SLOG读取数据. 确认的写入仅暂时存储在那里，直到它们被刷新到较慢的硬盘驱动器。它只是为了确保在断电或写入失败的情况下，确认的写入不会丢失，并且它们会尽快刷新到永久存储设备.

	另请注意，在没有 SLOG 设备的情况下，ZIL 将用于相同目的
- Special: from zfs 0.8

	用于存储文件系统的元数据和小块数据。通过将元数据和小块数据存储在**高速设备（如 SSD, nvme）**上，可以显著提高 ZFS 文件系统的性能

	查看specila使用情况: `zdb -bb mypool`
	移除special是异步过程(zpool remove返回1且stderr有内容), zfs会迁移数据, 见zpool status的remove; 移除过程不能对相关盘执行zpool labelclear否则会导致pool I/O suspend
	zpool remove可能需要多次(即使使用了`-w`), 一次remove多个special mirror时, zpool status显示completed, 实际可能只移除了一个mirror, 因此需要重复执行以删除剩余mirror. 从zpool remove执行效果来看, 一次只能移除一个vdev
	zpool存在special情况下移除cache vdev报"cannot determine indirect size of xxx: no such device in pool", 此时试试将vdev改为绝对路径

VDEV始终是动态条带化的. 一个 device 可以被加到 VDEV, 但是不能移除.

> zpool create 允许raidz1/2/3最低使用2/3/4块盘, 与标准的raid不一致, 推测: raidz是raid的软件解决方案内部,它自行解决了相关问题, 但还是推荐参照标准raid来配置硬盘.

## zpool
zfs支持分层组织filesystem, 每个filesystem仅有一个父级, 而且支持属性继承, 根为池名称.

> zpool描述了存储的物理特性, 因此必须在创建filesystem前创建池.
> 如果pool没有配置log device, pool会自行为ZIL预留空间.
> raid-z配置无法附加其他磁盘; 无法分离磁盘, 但将磁盘替换为备用磁盘或需要分离备用磁盘时除外; 无法移除除log device或cache device外的device
> 创建pool时, 不能使用其他pool的组件(vdev, 文件系统或卷), 否则会造成死锁.
> dataset 可以是 ZFS 池、文件系统、快照、卷和克隆. 它是可以存储和检索数据的 ZFS层.
> 在创建池时, ZFS 会根据**最小磁盘的容量**进行调整，以确保冗余和数据一致性

```sh
$ sudo zpool create [-m /mnt/p] pool-test /dev/sdb /dev/sdc /dev/sdd # 创建了一个零冗余的RAID-0存储池, zfs 会在`/`中创建一个目录,目录名是pool name. 如果mountpoint不可用(比如挂载点的父目录只读)也能创建成功, 仅是不挂载了而已
$ sudo zpool [option] list # 显示系统上pools的列表, `-o`只显示指定列,`-H`隐藏列头. size是所有磁盘的大小, free是剩余未被使用的磁盘大小. 看pool实际可用大小用`zfs get all <pool>`的availabled, 已用used.
$ sudo zpool status [-D] [-L] <pool> # 查看pool的状态,read/write列显示读写io时的错误次数, cksum列显示设备对读取请求返回损坏数据(校验和错误)的次数. `-v`输出详细信息, `-D`, dedup信息;`-x`仅显示有错误或因其他原因不可用的pool; `-L`显示vdev的真实设备名
$ sudo zpool destroy <pool> # 销毁pool
$ sudo zpool destroy <pool>/data-set # 销毁dataset
$ sudo zpool upgrade [<pool> | -a] # 更新 zfs 时，就需要更新指定/全部池
$ sudo zpool add <pool> /dev/sdx # 将驱动器添加到池中
$ sudo zpool exprot <pool> # 如果要在另一个系统上使用该存储池，则首先需要将其导出. zpool命令将拒绝导入任何尚未导出的存储池
$ sudo zpool export oldname
$ sudo zpool import <name> # 导入pool
$ sudo zpool import oldname newname # 重命名已创建的zpool的过程分为2个步骤, export + import
$ sudo zpool import -d /dev/disk/by-uuid # scan pool
$ sudo zpool import -F <poolname> -d /dev/disk/by-uuid # `-F`强制导入, 可解决突然断电导致的pool数据损坏(在truenas scale上碰到多次该损坏而导致系统启动进入Initramfs的问题)
$ sudo zpool create <pool> mirror <device-id-1> <device-id-m1> mirror <device-id-2> <device-id-m2> # 创建RAID10
$ sudo zpool add <pool> log mirror <device-id-1> <device-id-2> # 添加 SLOG
$ sudo zpool add <pool> -f spare /dev/sdf # 添加热备, 大小应>=max(pool's vdev),且无法移除当前正在使用的热备. 移除用`zpool remove`
$ sudo zpool add <pool> cache <device-id> # 添加L2ARC
$ sudo zpool iostat -v <pool> N # 每隔N秒输出一次pool的io状态. 第一次输出的数值可能较大, 因为那是`(pool_imported_total_count-0)/(time.now - pool_imported_at)`, 第二次开始才是`(pool_imported_total_count-lasttime_total_count)/(time.now-lasttime)`
$ sudo zpool iostat <pool> <vdev> # pool中vdev的iostate
# sudo zpool iostat -vly 1 1
$ sudo zpool remove <pool> mirror-1 # 移除mirror/不在使用的热备
$ sudo zpool remove <pool> /dev/sdf
$ sudo zpool attach <pool> <existing-device> <new-device> # 将新设备追加到已有vdev
$ sudo zpool detach  # 分离设备, 对象必须是mirror中的设备/raidz中已由其他物理设备或备用设备替换的设备
$ sudo zpool split <pool> <new-pool> [device] # 拆分pool, 仅适用mirror设备, 通过`-R`可指定新池的挂载点
$ sudo zpool offline [option] <pool> <device> # 离线zfs设备, `-t`表示临时离线, 重启后会重新恢复到online.
$ sudo zpool online [option] <pool> <device> # 上线zfs设备, 新设备上线后会同步. `-e`可扩展LUN(即使用更大容量设备时, 使用完整大小), 默认不扩展.
$ sudo zpool clear <pool> [devices] # 池设备故障时清理错误, 默认清理池内的所有设备错误.
$ sudo zpool replace <pool> replaced-device [new-device] # 替换存储池中的设备, 比如使用热备盘
$ sudo zpool get all <pool> # 获取pool的所有属性, 其他free属性并不是pool的剩余空间, 剩余空间应该用`zfs get available <pool>`获取
$ sudo zpool get <property> <pool> # 获取pool的指定属性
$ sudo zpool set <property=value> <pool> # 设置pool的指定属性
$ sudo zpool history [pool] # 显示zfs和zpool命令的使用日志, `-l`使用长格式(追加: 用户名, 主机名执行操作的域), `-i`显示更详细的内部日志可用于诊断
$ sudo zpool iostat # 列出pools的io统计信息, `-v`显示pool包含的vdev的io统计信息, `-l`更多信息
$ sudo zpool export [option] <pool> # 导出pool, 该pool必须处于不活动状态, 导出后该pool在系统中不可见. `-f`强制取消已挂载的filesystem
$ sudo zpool import [option] [pool/id-number] # 导入pool, 导入时允许重命名,允许只读导入. `-m`表示导入mirror log, 默认不导入.`-d`导入非标准路径的设备/由file构成的pool.`-D`恢复已销毁的pool, 追加`-f`即`-Df`可表示已销毁的pool中某设备不可用也可恢复.
$ sudo zpool upgrade # 升级pool, 以使用新版的zfs功能. `-v`表示当前zfs支持的功能, `-a`表示升级到最新的zfs.
$ blkid /dev/sdg # 检查盘是否被zfs使用过
/dev/sdg: LABEL="t" ... TYPE="zfs_member" # t是pool name, zpool destroy后该信息还保留.
$ zpool set cachefile=/etc/zfs/zpool.cache tank # 强制更新pool.cache
$ zpool set cachefile = / etc / zfs / zpool.cache # 适用于故障转移配置(高可用)的设置, 此时必须由高可用软件显示import pool
$ zpool events [-v] # 获取zpool事件. `-v`, 更详细.
$ zpool import -d /dev # 查找未导入的pool
```

pool status:
- ONLINE : 正常
- UNAVAIL : 池的元数据遭到破坏, 或者若干个设备不可用且没有足够的副本支持其继续运行
- DEGRADED : 池中若干设备发生了故障, 因为冗余配置, 其数据仍然可用.
- SUSPENDED : 池正在等待恢复设备连接, 问题解决前该pool一直处于wait状态.

vdev status:
- DEGRADED : 虚拟设备出现过故障, 但仍可用. 常见于mirror或raidz设备缺少一个或多个组成设备. 此时pool的容错能力可能已遭损害.
- OFFLINE : 已脱机
- ONLINE : 正常工作
- REMOVED : 已物理移除了设备, 移除检测依赖硬件
- UNAVAIL : 无法打开vdev. 如果顶层的vdev为UNAVAIL, 则无法访问池中的任何设备.

pool的status是由其所有顶层vdev的status决定的. 如果pool处于UNAVAIL/SUSPENDED则完全无法访问该pool. pool处于DEGRADED时, 于正常情况比可能无法实现相同的数据冗余或吞吐量.

mirror/raidz设备不能从pool中删除, 但可增删不活动的hot spares(热备), cache, log device.

### zpool create
创建pool: `zpool create -f -m <mount> <pool> [raidz | raidz2 | raidz3 | mirror] <ids>`

参数:
- f : 强制创建pool, 用于解决"EFI标签"错误
- m : 指定挂载点, 默认是`/`. 如果挂载点目录存在数据会保存
- o : 设置kv属性, 比如ashift
- n : 仅模拟运行, 并输出日志, 用于检查错误, 但不是100%正确, 比如在配置中使用了相同的设备做mirror.

> `blockdev --getpbsz /dev/sdXY`获取扇区大小, zfs默认是512, 与设备扇区不一致会导致性能下降, 比如具有4KB扇区大小的高级格式磁盘, 此时zfs会将ashift更改为12(2^12=4096). zfs默认ashift=0, 由其根据扇区大小决定ashift的具体值. 更改[ashift](https://github.com/zfsonlinux/zfs/wiki/faq#advanced-format-disks)选项的唯一方法是重新创建池. `zdb -U /data/zfs/zpool.cache | grep ashift`可查看ashift default的具体值.

> ashift仅在创建pool或添加新vdev时使用(此时仅针对新vdev). `zpool get all <pool>`可查询ashift, 可用`smartctl -a /dev/sdb`查看`sector size`

### zpool scrub
只要读取数据或zfs遇到错误，`zpool scrub`就会在可能的情况下进行静默修复.

```sh
$ zpool scrub <pool> # 检修 pool, 异步执行, 命令启动zpool后台检查后会立即返回. 此时zpool status输出的scan会显示scrub的执行信息.
$ zpool scrub -s <pool> # 取消正在运行的检修
```

> 建议是每周/月检修一次.
> 替换磁盘而同步数据很耗时, 替换间执行`zpool scrub`有利于替换设备运行正常且数据写入正确.

resilvering 是替换磁盘后重建冗余组的过程. 它是一个低优先级的操作系统进程, 在一个存储系统非常繁忙时需要更多时间.

## zfs
```sh
$ zfs list -t all -r <dataset>
$ sudo zfs list # 显示系统上pools/filesystems的列表, `-r`递归显示fs及其子fs, `-o`指定要显示的属性,比如space(空间使用); `-t`指定显示的类型, 比如filesystem, volume, share, snapshot.`-H`表示脚本模式: 不输出表头并用单个tab分隔各列; `-p`:精确显示数值; `-d`: 与`-r`连用,限制递归深度; `-s`按指定列升序排序; `-S`:与`-s`类似, 但以降序排序.
$ sudo zfs get [ all | property[,property]...] <pool> # 获取pool的参数. `-s`指定要显示的source类型; `-H`输出信息去掉标题, 并用tab代替空格来分隔
$ sudo zfs set atime = off sync=disable <pool> # 设置pool参数
$ sudo zfs set compression=gzip-9 mypool # 设置压缩的级别
$ sudo zfs inherit -rS atime  <pool> # 重置参数到default值. `-r`以递归的方式应用inherit子命令
$ sudo zfs get keylocation <pool>/<filesystem> # 获取filesystem属性
$ sudo zfs get volsize -o value -H -p <pool>/<vol> # 获取zvol大小
$ sudo zfs set acltype = posixacl <pool> / <filesystem> # 使用ACL
$ sudo zfs set sharenfs=on <pool> # 通过nfs共享pool
$ sudo zfs create -o mountpoint=none mypool/test/storage # 创建未挂载的dataset, 常用于zfs recv的场景.
$ sudo zfs set sharenfs=on <pool>/<filesystem> # 通过nfs共享filesystem
$ sudo zfs set mountpoint=/<pool>/... <pool>/... # 设置挂载点, 设置后会立即挂载.
$ sudo zfs destroy filesystem|volume # 销毁文件系统/volume, 此时dataset必须是不活动的. `-r`表示递归销毁(子代+snap+自身), `-R`表示`-r`+clone(此时该clone不能是busy即正在使用), `-f`删除前先unmout fs,仅用于fs, `-n`模拟删除, `-v`输出销毁过程的细节
$ sudo zfs rename <old-path> <new-path> # 重命名fs
$ sudo mount -o <pool>/.../<filesystem> # 挂载fs
$ sudo unmount <pool>/.../<filesystem> # 取消挂载fs, 此时fs必须是不活动的. `-f`强制取消挂载
$ sudo zfs create [-s] -V 5gb system1/vol # 创建5g大小的卷(创建卷时，会自动将预留空间设置为卷的初始大小，以确保数据完整性), `-s`创建精简卷(refreservation<volsize, 非精简是两者相等), 有点类似EMC存储的thin provisioning卷, 使用时(延迟)分配空间, 因此分配的大小可超过实际存储的大小. [blocksize支持: The default blocksize for volumes is 8 Kbytes. Any power of 2 from 512 bytes to 128 Kbytes is valid](https://linux.die.net/man/8/zfs)
$ sudo zfs get mountpoint mypool
$ zfs list -o space # 显示used属性
$ zfs get checksum tank/ws
$ zfs list -t all -o space 显示所有space属性(AVAIL USED USEDSNAP USEDDS REFRESERV USEDCHILD)
```

zfs list的属性可参考[freebsd zfs#Native Properties](https://www.freebsd.org/cgi/man.cgi?zfs(8)), `man zfs`与实际zfs包含的功能不一致, 有缺失, 比如`createtxg`, 部分属性说明:
- guid : 对象的全局唯一标识符（GUID)
- createtxg : 创建时的事务id. bookmark与snapshot有相同的createtxg. 该属性常用于`zfs send/recv`
- creation : 创建时间, unix时间戳
- REFER : 即referenced
- recordsize: fs块大小. 更改记录大小只会对新文件产生影响, 它不会对现有文件产生任何影响. 之前将recordsize从128k改到16k, 发现其fs性能没变化(212M/s); 但创建fs时recordsize初始为16k, 性能为144M/s.

> 物理卷(Physical Volume, PV)：操作系统识别到的物理磁盘(或者RAID提交的逻辑磁盘LUN), 物理卷可以是一个磁盘，也可以是磁盘中的一个分区.
> volume : 通常是指逻辑卷, 是逻辑卷组(VG, 由若干PV组成)上的一块空间, 上面没有文件系统.
> zfs fs/volume名称的parent必须是pool/fs.

与大多数其他文件系统不同，zfs具有可变的记录大小，或者通常称为块大小, 默认情况下，zfs上的记录大小为128KiB，这意味着它将根据要写入的文件大小动态分配从512B到128KiB任意大小的块.

RDBMS倾向于实现自己的缓存算法，通常类似于zfs自己的ARC. 因此最好禁用zfs对数据库文件数据的缓存，并让数据库自己完成.

### zfs create
```sh
# zfs create <pool>/<filesystem>/... # 在pool下创建filesystem(必须使用完整路径), filesystem除了快照外，还可以提高控制级别, 比如配额(可通过`zfs list`或`df`查看).
```

参数:
- -o : 设置fs的属性
- -f : 销毁时取消挂载, 取消共享
- -r : 递归销毁
- -R : 递归销毁, 包括依赖其的clones
- -m : 设置挂载点

> 文件系统默认挂载在pool下, 除非指定了mountpoint属性
> 为了能够创建和mount filesystem，zpool中不得预先存在相同名称的目录.
> zfs 创建逻辑卷并关闭压缩再格式化成ext4并挂载, 使用`dd if=/dev/zero`写入数据, `zfs get all xxx`上统计的存储相关属性不正确, 需换其它的输入源,比如使用磁盘复制`dd if=/dev/sda1`. zfs vol挂载后df显示的大小是`volsize`而非`refreservation`.
> 创建的vol会对应到`/dev/{dataset_path}和/dev/zvol/{dataset_path}`文件, 推荐使用`/dev/zvol/{dataset_path}`. `/dev/zvol/{dataset_path}`仅是`/dev/zdXXX`的软连接, 由`udev rule(/lib/udev/rules.d/60-zvol.rules)`创建
> [zvol创建volume时,kernel会先创建`/dev/zdXXX`,再通过udev rule异步创建`/dev/zvol/{dataset_path}`, zfs create时可用`udevadm monitor`监控到这些时序信息](https://github.com/openzfs/zfs/issues/7875), 因为linux udev是异步创建device, 因此`/dev/zvol/{dataset_path}`可能没立即出现, 应该轮询一下.

> `zvol_id /dev/zd<N>`可获取该zvol的dataset

> ZFS Management API (libzfs)使用/dev/zfs.

### zfs snapshot
ref:
- [ZFS Daily Tips (3) - snapshot & checkpoint](https://zhuanlan.zhihu.com/p/75755509)

zfs snapshot（快照）是 zfs 文件系统或卷的**只读**拷贝(即无法修改属性). 可以用于保存 zfs 文件系统的特定时刻的状态，在以后可以用于恢复该快照，并回滚到备份时的状态.

创建一个快照意味着记录文件系统 vnode 并跟踪它们.

> snap直接占用pool, 不使用后备存储.
> 快照放在对应fs的`.zfs/snapshot`目录里.

```sh
$ sudo zfs snapshot -r mypool/projects@snap1 # 创建 mypool/projects 文件系统的快照. `-r`表示递归创建(即为所有后代文件系统创建快照)
$ sudo zfs list -t snapshot # 查看所有的snapshots列表. USED, 快照占用空间, REFER, 快照引用的数据所站空间.
$ sudo zfs list -t snapshot -d 1  mypool/p13 # 只显示指定dataset的snaps
$ sudo zfs list -t snapshot -r  mypool/p13 # 只显示指定dataset的snaps
$ sudo zfs rollback mypool/projects@snap1 # 回滚快照
$ sudo zfs destroy mypool/projects@snap1 # 移除snapshot
$ sudo zfs destroy mypool/projects@% # %表示限定范围, 其两边为空默认表示最早~最晚
$ sudo zfs hold keep mypool/home@today # 保持快照, `-r`表示递归
$ sudo zfs holds mypool/home # 显示保持的快照的列表
$ sudo zfs release keep mypool/home # 释放快照的保持标志, 之后可用`zfs destroy销毁`.`-r`表示递归
$ sudo zfs rename mypool/home@today yesterday # 重命名快照. `-r`表示递归
$ sudo zfs diff  [-FHt] <snapshot> [snapshot|filesystem] # 显示差异
```

缺省情况下，`zfs rollback`无法回滚到除最新快照以外的快照. 要回滚到之前的快照，必须通过指定 -r 选项来销毁所有中间快照.
如果存在任何中间快照的**克隆**，还要使用 -R 选项销毁克隆.

### zfs clone
一个zfs clone是文件系统或卷的可写的拷贝. 一个 zfs clone 只能从 zfs snapshot中创建，该 snapshot 不能被销毁，直到该 clones 也被销毁为止.

> **只要克隆存在,就无法销毁原始快照**
> 创建clone几乎是即时的, 且最初不占用其他磁盘空间. 此外, 还可以对clone进行快照.
> 可用`zfs get/set`操作clone的属性.
> 可通过查询dataset@sanp的clones属性获知它有哪些clone.

```sh
# 克隆 mypool/projects，首先创建一个 snapshot 然后 clone
$ sudo zfs snapshot mypool/projects@snap1
$ sudo zfs clone mypool/projects@snap1 mypool/projects-clone
$ zfs get origin mypool/projects-clone # 获取origin snap
$ zfs get all mypool/projects@snap1 # clones字段显示基于它创建的clone, 该snap没有clone时可能没有该属性
$ zfs get all -r mypool/projects
$ sudo destroy mypool/projects-clone # 销毁clone
```

### zfs send 和 receive
zfs send 将文件系统的快照写入stdout，然后流式传送到文件或其他机器. zfs receive 从stdin接收该 stream 然后写到 snapshot 拷贝，作为 zfs 文件系统. 这有利于通过拷贝或网络传送备份.

> zfs send默认是完整流(从fs创建开始到指定的快照为止的所有内容).
> 增量流: 包含一个快照与另一个快照之间的差异.

```sh
# 创建 snapshot 然后 save 到文件
$ sudo zfs snapshot -r mypool/projects@snap2
$ sudo zfs send mypool/projects@snap2 > ~/projects-snap.zfs  # `-c`使用压缩(如果mypool/projects是活动的则必须使用), `-n`表示模拟send, 实际不产生数据流, `-P`表示生成流的信息, 比如全量/增量, 数据流大小.`-v`: 向(stderr)发送流的统计信息, 包括每秒传输多少.
$ sudo zfs receive -F mypool/projects-copy < ~/projects-snap.zfs # 恢复,此时目标dataset必须存在. `-F`表示(此时目标必须没有快照)忽略目标dataset的改动(mypool/projects-copy), 全量的话是直接覆盖原有dataset, 增量的话是回滚到该增量快照的起点后再应用增量. `-d`: (此时目标dataset必须存在)去掉原快照名称中的pool name,使用目标dataset name+剩余名称作为新名称.
$ sudo zfs send -i @old_snap1  ool/dana@new_snap2 # `-i`增量发送,`-I`将一组增量快照合并为一个快照,`-R`表示复制 zfs 文件系统和其后代.
$ sudo zfs send pool/dana@snap1 | ssh system2 zfs recv pool/dana
$ zfs send ... | gzip | <network> |   gunzip | zfs recv otherpool/new-f # 中间可使用压缩, 或其他更快的压缩, 比如lz4.
$ zfs send --compressed tank/my-fs@today | ... # 最佳方式, 前提条件是发送端dataset已启用压缩属性
$ sudo zfs send -i system1/dana@snap2 system1/dana@snap3 | ssh sys2 zfs recv -F newsys/dana
```

接收文件系统快照时的要求：
- 接收增量快照时，目标文件系统必须已存在
- 将取消挂载文件系统和所有后代文件系统
- 文件系统在接收期间不可访问
- 目标系统上不得存在名称与要接收的源文件系统相同的文件系统. 如果文件系统名称在目标系统上已存在，请先重命名文件系统.

> [zfs send/receive过程中中止可能产生隐藏的clone(`zfs list`不可见, 但`zfs list <clone_dataset>`可见)](https://www.reddit.com/r/zfs/comments/h7suck/how_do_you_list_hidden_clones_eg_filesystemrecv/), 命名是`<dataset>/%recv`, 可`zfs destroy`

#### zfs resume
参考:
- [Resuming ZFS send](https://www.oshogbo.vexillium.org/blog/66/)
- [zfs send receive resume compressed](https://forum.proxmox.com/threads/zfs-send-receive-resume-compressed.33612/)

演示:
```sh
192.168.0.41 # zfs send -Pv storage/b@zrepl_20191016_071546_000 | ssh root@192.168.0.42 "zfs receive -s mypool/test/resume" # `mypool/test/resume`
192.168.0.41 # ... 未完成就强制中断
192.168.0.42 # zfs get -H -o value receive_resume_token mypool/test/resume # receive_resume_token 仅在zfs recv中断的情况下出现在属性中
1-c2d6054a6-d0-789c636064000310a500c4ec50360710e72765a5269740c418b0c9a7a515a7968064ccc460f26c48f2499525a9c540fa80a39b1836fd25f9e9a599290c0cfadf0ddc972975d93920c97382e5f312735319188a4bf28b12d353f5931caa8a520b72e28d0c0c2d0d0d0ccde20dcc0d4d4d80948101030c0000c6571c90
192.168.0.41 # zfs send -vt 1-c2d6054a6-d0-789c636064000310a500c4ec50360710e72765a5269740c418b0c9a7a515a7968064ccc460f26c48f2499525a9c540fa80a39b1836fd25f9e9a599290c0cfadf0ddc972975d93920c97382e5f312735319188a4bf28b12d353f5931caa8a520b72e28d0c0c2d0d0d0ccde20dcc0d4d4d80948101030c0000c6571c90 | ssh root@192.168.0.42 zfs receive -s mypool/test/resume
```

如果想放弃resume 可执行`zfs receive -A otherpool/new-fs`即可清除.

### zfs share
用于通过nfs或smb协议共享和发布zfs文件系统, 通过创建fs时设置share.nfs/share.smb属性来实现.

> openzfs share时使用的是sharenfs/sharesmb.

```sh
$ sudo zfs create -o share.nfs=on tank/sales # 创建共享的fs, 默认为可读可写. share.nfs属性会被继承到任何后代fs中.
$ sudo zfs create -o share.nfs.ro=\* tank/sales/logs # 设置fs只读
$ sudo zfs get share.nfs.all tank/sales # 获取所有share.nfs属性
```

### zfs Ditto Blocks，重复块
Ditto blocks 创建更多的冗余拷贝.

对于只有一个设备的storage pool，ditto blocks are spread across the device, trying to place the blocks at least 1/8 of the disk apart. 对于多设备的 pool，zfs 试图分布 ditto blocks 到多个独立的 VDEVs, 可设置1~3 份拷贝.

```sh
$ sudo zfs set copies=3 mypool/projects
```

### zfs Deduplication，文件去重
zfs dedup 将丢弃重复数据块，并以到现有数据块的引用来替代. 这将节约磁盘空间，这**需要大量的内存和cpu**. 内存中的去重记录表DDT (deduplication tables)需要消耗大约 ~320 bytes/block. 表格尺寸越大，写入时就会越慢.

```
$ sudo zfs set dedup=on mypool/projects # 启用去重
```

每当数据写入时, zfs都会将其与已记录的块比较, 如果存在相同, 则它不会写入物理数据, 反而会添加一些元信息代替写入, 从而节省很多和大量的磁盘空间.

一旦开启重删, 即使关闭它, import pool时DDT还是会被载入内存. 解决方法: 
1. 备份数据, 重新创建pool(此时不开重删)再恢复数据
2. 将数据挪到没有重删的pool上

`zpool list`的`DEDUP`列显示`1.00x`表示没有重删.

**dedup**只对开启后的新数据块有效. 根据经验: **不推荐使用, 不如用压缩**

## [属性](https://openzfs.github.io/openzfs-docs/man/7/zfsprops.7.html)
ref:
- [ZFS 属性介绍](https://docs.oracle.com/cd/E24847_01/html/819-7065/gazss.html)

- avail: 文件系统中的可用空间

	- fs : avail = refquota + usedbysnapshots - used
- used : 只读属性，表明此数据集及其所有后代占用的磁盘空间量.
	
	可根据此数据集的配额和预留空间来检查该值. 使用的磁盘空间不包括数据集的预留空间，但会考虑任何后代数据集的预留空间. 数据集占用其父级的磁盘空间量以及以递归方式销毁该数据集时所释放的磁盘空间量应为其使用空间和预留空间的较大者.

	在父数据集的已用磁盘空间计算中会考虑预留空间，并会针对其配额、预留空间或同时针对两者进行计数.

	used = usedbychildren + usedbydataset + usedbyrefreservation + usedbysnapshots

	- usedbychildren(usedchild) : 只读属性，用于标识此数据集的后代占用的磁盘空间量；如果所有数据集后代都被销毁，将释放该空间
	- usedbydataset(usedds) : 只读属性(已用空间)，用于标识数据集本身占用的磁盘空间量；如果在先销毁所有快照并删除所有 refreservation 预留空间后销毁数据集，将释放该空间.
	- usedbyrefreservation(refreserv) : 只读属性，用于标识针对数据集设置的 refreservation 占用的磁盘空间量；如果删除 refreservation，将释放该空间
	- usedbysnapshots(usedsnap) : 只读属性，用于标识数据集的快照占用的磁盘空间量. 特别是，如果此数据集的所有快照都被销毁，将释放该磁盘空间. 请注意，此值不是简单的快照已用空间的总和，因为多个快照可能引用了相同的块
- quota : 限制**数据集及其后代**可占用的磁盘空间量. 一旦达到该容量后，不管存储池的可用空间有多大，都无法再将数据写入该数据集.

	该属性可对已使用的磁盘空间量强制实施硬限制，**包括后代（含文件系统和快照）占用的所有空间**. 对已有配额的数据集的后代设置配额不会覆盖祖先的配额，但会施加额外的限制. **不能对卷设置配额**，因为 volsize 属性可用作隐式配额.
- refquota : 限制该数据集可以使用的磁盘空间量, 该限制不包括后代所占用的磁盘空间
- reservation : 预留空间, 从池中分配的保证可供数据集**及其后代**使用的磁盘空间. 预留空间受存储池中的可用空间容量限制.
- refreservation :  预留空间来保证用于**该数据集**的磁盘空间，该空间不包括快照和克隆使用的磁盘空间. 此预留空间计算在父数据集的使用空间内，并会针对父数据集的配额和预留空间进行计数.

	如果 usedbydataset 空间低于此值，则认为数据集正在使用 refreservation 指定的空间量(也是数据集的初始大小是refreservation). usedbyrefreservation 数字表示该额外空间，并添加到数据集占用的总 used 空间，继而占用父数据集的用量、配额和预留空间. 它会随实际使用而减少(但也可能不变), 这将通过确保提前预留用于未来写操作的空间，防止数据集过度承载池资源
- referenced : 此数据集可访问的数据量，这些数据可以与池中其他数据集共享，也可以不与其他数据集共享. 创建快照或克隆时，首先会引用与创建时所基于的文件系统或快照相同的空间量，因为其内容相同.
- volsize(for volumn)/refquota (for filesystem): 允许dataset的最大大小, 它们应该>=refreservation.
- logicalreferenced : dataset在逻辑上使用的空间.
- logicalused : dataset及其后代在逻辑上使用的空间.
- written

	- snapshot: 卷的当前快照与上一张快照的referenced变化量, 即快照间数据变化量

> Reservation是最小值, quota是最大值.

## zdb
ref:
- [Turbocharging ZFS Data Recovery](https://www.delphix.com/blog/openzfs-pool-import-recovery)

选项:
- -G, --dump-debug-msg: zdb结束前dump出zfs_dbgmsg内容

- `zdb -l /dev/sdj` : 查看磁盘上的zpool信息
- `zdb -dddddddd testpool` : 查看写入的range, 去重表（DDT）等
- `zdb -c -eFX -G <pool>`:  检查无法import的pool, 很慢
- `zdb -dep <disk_dir/disk> -G <pool>`: 检查无法import的pool, 慢
- `zdb -V -ep <disk_dir/disk> -G <pool>`: 检查无法import的pool, 快
- `/proc/spl/kstat/zfs/dbgmsg` : import log, 按时间戳变化分割日志

## lib
ref:
- [zpool-status](https://pypi.org/project/zpool-status/)

	json output

- [JSON output support for zfs and zpool commands #16217](https://github.com/openzfs/zfs/pull/16217)

	zfs 2.3.0

## 周边
- [zfsbackup-go](https://github.com/someone1/zfsbackup-go)

## FAQ
### 加密
ref:
- [Encrypting ZFS File Systems](https://docs.oracle.com/cd/E26502_01/html/E29007/gkkih.html)

必须创建dataset时就启用, 否则就无法启用了, 方法:`-o encryption=on -o keyformat=passphrase`.

解锁源dataset后, 其快照视图和子数据集会自动解锁.

> 是否loaded key, 看keystatus属性(available/unavailable)

```bash
zfs load-key xxx
zfs mount xxx
zfs umount xxx
zfs unload-key xxx
```

### quota于refquota区别
如果对 tank/home 数据集设置了quota，则 tank/home 及其所有后代使用的总磁盘空间量不能超过该配额.
如果对 tank/home 数据集设置了refquota，则 tank/home 磁盘空间量不能超过该配额, 但不包括后代所占用的空间, 比如其快照/clone.

```bash
# zfs set quota=1.5G data/engineering
# zfs set refquota=1G data/engineering
# zfs set reservation=800M data/engineering
```

### reservation
如果 tank/home 指定了预留空间，则 tank/home 及其所有后代都会使用该预留空间. 预留空间计入父文件系统的已用磁盘空间, 并计入父文件系统的配额和预留空间.

### 文件读取很慢
原因:
1. pool的一块镜像盘的状态变成DEGRADED, 更换磁盘即可

### share
1. nfs
```
# zfs set sharenfs=on rpool/fs1
# zfs set acltype=posixacl rpool/fs1
# zfs set aclinherit=passthrough rpool/fs1 # 当前openzfs 的 aclinherit属性不支持posixacl
```

**acl针对的是目录**

nfs配置见[develop/nas.md](/develop/nas.md)

1. samba
```
# zfs set sharesmb=on rpool/fs1
```

### zfs 2.3.0编译
ref:
- [ZFS RAIDz expansion 抢先体验](https://ntzyz.space/zh-cn/post/testing-zfs-raidz-expansion/)
- [ZFS RAID-Z Expansion，RAID-Z 单盘扩展前瞻测试使用](https://www.truenasscale.com/2022/02/12/562.html)

ps: v2.3.1已在2025.3.11发布

```bash
# apt install linux-image-6.9.7+bpo-amd64 linux-headers-6.9.7+bpo-amd64
# apt install alien autoconf automake build-essential debhelper-compat dh-autoreconf dh-dkms dh-python dkms fakeroot gawk git libaio-dev libattr1-dev libblkid-dev libcurl4-openssl-dev libelf-dev libffi-dev libpam0g-dev libssl-dev libtirpc-dev libtool libudev-dev parallel po-debconf python3 python3-all-dev python3-cffi python3-dev python3-packaging python3-setuptools python3-sphinx uuid-dev zliblg-dev
# git clone https://github.com/openzfs/zfs
# cd ./zfs
# git checkout zfs-2.3.0-rc3
# sh autogen.sh
# ./configure
# make -s -j$(nproc) native-deb-utils / make -s -j$(nproc) && make deb
# cd ..
# rm openzfs-zfs-dracut_*.deb openzfs-zfs-initramfs_2.3.0-1_all.deb
# dpkg -i ./openzfs-zfs-zed_2.3.0-1_amd64.deb ./openzfs-zfs-dkms_2.3.0-1_all.deb ./openzfs-libuutil3_2.3.0-1_amd64.deb ./openzfs-libzfs6_2.3.0-1_amd64.deb ./openzfs-libnvpair3_2.3.0-1_amd64.deb ./openzfs-zfsutils_2.3.0-1_amd64.deb
# zpool import -a
# zpool upgrade <pool>
```

### zfs 2.0.0编译
参考:
- [官方Custom Packages制作指导](https://openzfs.github.io/openzfs-docs/Developer%20Resources/Custom%20Packages.html#)
- [arch zfs-linux](https://aur.archlinux.org/cgit/aur.git/tree/PKGBUILD?h=zfs-linux)
- [arch zfs-utils](https://aur.archlinux.org/cgit/aur.git/tree/PKGBUILD?h=zfs-utils)
- [在 Ubuntu 20.04 上引入最新版本的 OpenZFS](https://qiita.com/yamakenjp/items/380ea5bb338940b5dc55)

```bash
# # env deepin v20 amd64
# # [Building ZFS](https://openzfs.github.io/openzfs-docs/Developer%20Resources/Building%20ZFS.html)
# apt install build-essential autoconf automake libtool gawk alien fakeroot dkms libblkid-dev uuid-dev libudev-dev libssl-dev zlib1g-dev libaio-dev libattr1-dev libelf-dev linux-headers-$(uname -r) python3 python3-dev python3-setuptools python3-cffi libffi-dev
# cd <zfs-2.0.0 src>
# ./autogen.sh
# ./configure --prefix=/usr --sysconfdir=/etc --sbindir=/usr/bin --libdir=/usr/lib \
                --datadir=/usr/share --includedir=/usr/include --with-udevdir=/usr/lib/udev \
                --libexecdir=/usr/lib --with-config=all --enable-systemd # configure前必须安装alien否则`make deb`会报错; 不能将`--with-python`设为"no", 否则`make deb`根据rpm spec构建rpm时会报错"configure: error: Unknown --with-python value ':'"
# # 下面自行编译再install / 直接打包选一种即可
# make -s -j$(nproc) # 自行编译 / make deb # [zfs会先通过构建rpm再通过alien将rpm转成deb](https://github.com/openzfs/zfs/issues/10168)
# make deb
# dpkg -i \
        zfs-initramfs_2.0.0-1_amd64.deb \
        zfs_2.0.0-1_amd64.deb \
        python3-pyzfs_2.0.0-1_amd64.deb \
        libzpool4_2.0.0-1_amd64.deb \
        libzfs4_2.0.0-1_amd64.deb \
        libuutil3_2.0.0-1_amd64.deb \
        libnvpair3_2.0.0-1_amd64.deb \
        kmod-zfs-$(uname -r)_2.0.0-1_amd64.deb # from 0.8.1
```

> `make deb`因为会先构建rpm的原因, 导致会根据`rpm/xxx/yyy.spec.in`重新编译zfs, 且编译参数由rpm spec指定.

> 需要test时`dpkg -i`追加`zfs-test_2.0.0-1_amd64.deb`

> ZFS模块可以通过两种方式加载到内核，DKMS和kmod, 区别在于: 如果安装基于DKMS的ZFS模块，然后由于某种原因更新了操作系统的内核，则可用再次重新编译ZFS内核模块(需要相关的源码), 否则它将无法工作; 但是基于kmod的ZFS模块仅针对特定版本的kernel; 基于RHEL/Centos的Kabi升级内核无需处理. 具体可参考[Custom Packages里的说明](https://openzfs.github.io/openzfs-docs/Developer%20Resources/Custom%20Packages.html)

## FAQ
### [zfs test](https://openzfs.github.io/openzfs-docs/Developer%20Resources/Building%20ZFS.html)
相关脚本在`zfs-test_*.deb`中, 安装后执行即可`/usr/share/zfs && ./zfs-tests.sh -vx`即可.

### Error
#### libzfs.h: No such file or directory
```sh
$ sudo apt install libzfslinux-dev
```

### zfs rename后mkfs 报`No such file or directory`
env: zfs 0.7.7
zfs rename volume后`/dev/zvol/{datapath}`不变, mkfs时`No such file or directory`, 怀疑是bug.

zfs 0.8.1 rename后`/dev/zvol/{datapath}`会跟着变化, 且mkfs正常.

### zfs周边
- 复制

	- z3 : `pip install z3`,  Z3备份与恢复的基本原理是围绕zfs send和zfs receive的管道来实现的.
	- zrepl
### pool is busy
`fuser -vm /dev/zd640`
`lsblk` # using by drbd, lvm, kvm. lvm可用`ls /dev/mapper/*`查看, 用`dmsetup remove /dev/dm-2`删除; kvm可用`virsh ls+ virtsh shutdown`命令处理.
`lsof -t /dev/zd640`
`lsof -t <pool_mount_path>`
`blkid /dev/zd640`: 看看该盘是什么类型的磁盘. 如果zd640是lvm盘, 而lvm.conf未配置过过滤该设备, 那么就可能被lvm占用

### zfs 类clone 无法挂载
`XFS (zd32): Filesystem has duplicate UUID adf19c69-ebc4-4622-97e2-1ab899f8f5c3 - can't mount` from `syslog`.

操作步骤:
1. 创建volume, 格式化为xfs
1. 写入数据, 再执行sync(保证数据落盘)

	不用fsync, 是因为fsync不是递归的, 不能确保指定路径的子文件全部落盘.
1. 创建volume2, snap, zfs send/receive volume到volume2
1. mount volume2报错

> 此时使用`mount -o rw,nouuid /dev/zd64  /mnt`可以成功

解决方法:
```bash
# xfs_repair -L /dev/zd64
# xfs_admin -U generate /dev/zd64
```

### [blk_update_request: I/O error, dev sdc, sector 824769880 且 zpool status提示有write和cksum错误](https://github.com/openzfs/zfs/issues/3785)
可用`hdparm --read-sector 824769880 /dev/sdc`尝试多次读取.


### pool I/O is currently suspend
pool只有raidz0(1块), 没有其他盘, 物理移除该盘后, `zpool destroy`时报该错误.

查阅相关资料后发现只能通过reboot来解决, reboot后该pool消失.

### zfs编译时configure报`error: cannot guess build type; you must specify one`
原因: 自带`config.guess`提示当前编译类型无法找到

解决方法一：指定平台，手动编译
```bash
./configure --build=arm-linux
make -s -j$(nproc)
```

解决方法二：替换ZFS自带`config.guess`
```bash
mv config/config.guess config/config.guess.bak
cp /usr/share/automake-1.13/config.guess config/
make -s -j$(nproc)
```

### `zpool create x /dev/sdd`报`... invalid feature 'redaction_bookmarks'`
直接`apt remove zfsutils-linux`后直接通过`dpkg -i *.deb`安装了自打包的zfs 2.0.0导致创建pool时报了这个错误. 运行`zpool --version`输出了包含`0.8.3`即旧版本的zfs-kmod信息.

解决方法: 重启即可.

其实上面的升级zfs版本的步骤有问题, 应该: 删除旧版zfs后重启, 再安装新版本.

### 写满测试
场景: 4g pool创建10g 精简zvol, 用如下2种方式尝试写满:

1. `fio -filename=/dev/zd0 -rw=write -ioengine=psync -iodepth=16 -numjobs=1 -ramp_time=30 -direct=1 -runtime=300 -time_based -group_reporting -bs=1MB -size=10GB -name=test`

	- zfs 0.7.7(ubuntu 16.04)

		一次写完不报错; 3次写完报"No Space", 之后fio(参数同上), dd(direct)写都报该错, 但dd(no direct, 1g)不报错, 此时推测是dd异步写没处理报错的原因.


	- zfs 2.0.0(ubuntu 20.04, 自deb打包)

		每次写完10G不报错

1. `dd if=/dev/urandom of=/dev/zd0 bs=1M count=10240 oflag=direct`

	- zfs 0.7.7

		用同上参数的fio写立马报错; `dd`也直接报错

	- zfs 2.0.0

		用同上参数的fio写立马报错; `dd`再写三次都会写入少量数据, 之后就无法写入而是直接报错


还是上面的精简vol格式化成ext4后批量cp文件:
	- zfs 0.7.7

		写入3950M后变成read-only fs. umount后重新挂载报错: "can't read superblock", 因此禁止做nas.

		将该不能重新挂载的vol做成iscsi, 发现iscsi client能挂载但不能格式成ext4, 会报错.

		删除该vol, 重新创建一个相同大小的精简vol, 做成iscsi后, client能挂载也能格式化成ext4, 批量cp发现写入了**9072M**, 根据时间戳cp后写入的文件报错"no such file or directory", 但进入挂载目录发现该文件有大小87M. umount后再mount也报"can't read superblock".

	- zfs 2.0.0

		写入4577M后变成read-only fs, umount后重新挂载变正常, 可删除, 有空间后仍可创建文件, 但批量cp文件后又变成只读fs.

还是上面的pool, 创建非精简vol 3.5g并格式化成ext4后cp文件:
	- zfs 0.7.7

		写满后不会变成只读而是报"No space left on device", 部分文件大小为0(写入部分后再遇到空间不足, 文件会创建但大小为0). 删除部分文件后, 仍可再cp文件

	- zfs 2.0.0

		写满后不会变成只读而是报"No space left on device", 部分文件大小为0(写入部分后再遇到空间不足, 文件会创建但大小为0). 删除部分文件后, 仍可再cp文件

### zfs 在磁盘开始的位置
参考:
- [C语言读取GPT分区信息 # PMBR](https://blog.csdn.net/qq_37734256/article/details/88384750)

env: zfs 2.0.0

zfs使用一块盘前会将该盘划为两个分区(gpt), 一个8M([Solaris reserved 1](https://github.com/openzfs/zfs/issues/6110), 假设是sdd9), 一个剩余空间(假设是sdd1), 分区表中8M在后.

```bash
zdb -l /dev/sdd1 # 获取zfs label
hexdump -C -n 102400 /dev/sdd1 # 结合上面的label信息, 发现zfs从0x3fd8开始使用磁盘
```

通过`wipefs -a /dev/sdd`抹除分区表, 并将mbr内容全部置为`1`, 再创建pool发现, 磁盘的mbr内容已变化即pmbr被重置了, 因此制作磁盘签名时要回写.

根据PMBR的定义可回写位置是0~446B内.


**推荐: truenas scale方案**, 使用`/dev/disk/by-partuuid`中的partuuid作为磁盘签名.

### zpool加密
参考:
- [Encrypting ZFS File Systems](https://docs.oracle.com/cd/E53394_01/html/E54801/gkkih.html)
- [How-To: Using ZFS Encryption at Rest in OpenZFS (ZFS on Linux, ZFS on FreeBSD, …)](https://blog.heckel.io/2017/01/08/zfs-encryption-openzfs-zfs-on-linux/)

```bash
zpool create -O encryption=on -O keylocation=file:///root/keys/hdd256.key \
             -O keyformat=raw \
             mypool /dev/disk/by-id/mydisk
zfs get encryption,keystatus,keysource,pbkdf2iters mypool
```

`-O pbkdf2iters=350000`用于`-O keyformat=passphrase`选项, 迭代passphrase来保证安全.

`head -c 32 /dev/urandom > /dev/shm/enc3key` for `-O encryption=aes-256-gcm`

import pool时, keystatus可显示加密pool的状态.

`zfs load-key -L file:///xxx.key -n <pool>`加载秘钥.


### zfs create超慢
检查pool avail的大小.

### zfs删除离线pool
> env: zfs 0.8.6

- raid0 ： x64/arm64 zfs module均panic
- raiz1 ： x64能强制删除; arm64, 第一次普通删除阻塞, 第二次强制删除报错, 之后无论普通/强制删除均报错.

### zfs snap消耗的空间
创建snap时, 它使用根pool文件系统和snap的空间, 以及可能的引用它父级快照的空间.

```bash
zpool set listsnapshots=on rpool # zfs list时也输出snap信息, 默认是off
```

### arc_summary
arc_summary:
```
ARC size (current):
	Max size: 按配置最多使用的内存量
```

### arcstat
arcstat 1
arcstat -o /tmp/a.log 2 10
arcstat -s "," -o /tmp/a.log 2 10
arcstat -v
arcstat -f time,hit%,dh%,ph%,mh% 1

最大 ARC 缓存内存（c）、当前 ARC 缓存大小（arcsz）、从 ARC 缓存中读取的数据（read）等信息

字段:
- c : the target size of the ARC in bytes, 当前arc大小
- c_max : the maximum size of the ARC in bytes
- size : the current size of the ARC in bytes:  arc总大小

> arcstat from `/proc/spl/kstat/zfs/arcstats`

### 调整arc大小
`/etc/modprobe.d/zfs.conf`:
```bash
options zfs zfs_arc_max=<memory_size_in_bytes>
```

> `echo $((5*2**30))`=`python3 -c "print(5*2**30)"`=`5368709120`

zfs_arc_max=0在arm64上可能不生效, 导致arc内存使用超过50%, 此时可追加zfs_arc_sys_free(保留空余内存的下限)以限制.

遇到:
1. arm64 sdb(in zfs)的`Time spent Doing I/Os`(by node_exporter)一直是100%, 导致mem使用突破了zfs_arc_max和zfs_arc_sys_free的限制, 平时sdb io没问题的情况下, 没有出现该问题

	该问题也可能是其他原因引发, 但目前观察到zfs arc异常时间刚好和sdb的io异常时间重合

### 清理cache
```
# sync # 设置/proc/sys/vm/drop_caches的前提, 清理dirty object即`/proc/meminfo`中的Dirty.
# echo 3 > /proc/sys/vm/drop_caches
# echo 1 > /proc/sys/vm/drop_caches # 会忽略arc
```

/proc/sys/vm/drop_caches, 见`man proc`:
- 1: 清理page cache. 影响free中的buffer和cached, 进而影响了used和free
- 2: 清理dentry和inode. 影响free中的buffer
- 3: = 1 + 2

[zfs_arc_min是0时, Linux buffer/cache可能会将 ARC 从内存中逐出, 因此推荐设置zfs_arc_min.](https://serverfault.com/questions/857350/zfs-arc-cache-and-linux-buffer-cache-contention-ubuntu-16-04)

### arcstat数据来源
`/proc/spl/kstat/zfs/arcstats`, **arc使用的memory算入`free`的`used`**

### zfs使用内存
`cat /proc/spl/kmem/slab |awk '{a+=$3}END{print a}'`, **这部分memory不包括arc, 仅是zfs objects**

### zfs命令卡住
- 用iostat发现磁盘写入超过20MB/s时, `%util`已超过90%, await也挺高, 是磁盘性能不够导致. 磁盘空闲后, 该命令可很快成功.
- 服务器存在fc多路径盘, fc离线后, 存储软件将pool使用的盘的软连接指向了其多路径中的某条路径的磁盘, 从而导致zfs异常, 同时导致使用本地盘创建的其他pool创建卷时也卡住了

	syslog里存在分区信息错误的record, 这里是`kernel: GPT:disk_guids don't match.\nkernel: GPT:partition_entry_array_crc32 values don't match\nkernel: GPT: Use GNU Parted to correct GPT errors.`

	解决方法: 重新扫描fc, 待磁盘上线后重建软连接, 可能需要重启系统.

### zfs fs写放大35左右%
ref:
- [ZFS Recordsize 和 TXG 探索](https://tim-tang.github.io/blog/2016/06/01/zfs-recordsize-txg)

原始文件124M, 拷入fs

- recordsize=4k: 168M左右
- recordsize=8k: 148M左右 
- recordsize=16k及以上: 123M左右

ssd pool且recordsize=4k不存在该问题.

> truenas 22.02.4(22.9.27)上也发现该问题.

### truenas+vmware `Disable physical block size reporting`
如果initiator不支持超过 4K 的物理块大小值, 则设置. 该设置还可以防止在 ESXi 中使用此共享时出现恒定块大小警告.

vmkwarning.log:
```
WARNING: ScsiPath: 4394: The Physical block size "131072" reported by the path vmhba64:C0:T4:L0 is not supported. The only supported physical blocksizes are 512 and 4096
WARNING: ScsiDeviceIO: 6462: The Physical block size "131072" reported by the device naa.6589cfc0000000572b71f35019e9c31f is not supported. The only supported physical blocksizes are 512 and 4096
```

### zfs fs可用空间可能出现负数
ref:
- [quota: extend quota for dataset](https://github.com/openzfs/zfs/pull/13839)

直接置为0即可

### arc_prune导致cpu stuck
env:
- kernel: 4.19.90
- cpu: FT-2000+/64

在arm环境容易遇到, 问题环境内存为63G, 遇到时的状态:
1. free 3G, buffer/cache 25G, available 28G, swap(1G) used 1G
1. free 11G, buffer/cache 19G, available 30G, swap(1G) used 1G
1. free 21G, buffer/cache 11G, available 31G, swap(1G) used 1G

观察发现, 在内存不足的情况下更容易遇到.

### zpool create报"/dev/vdc is part of active pool 'p'"
blkid查看vdc提示`LABEL="p" UUID="801774493520823192" UUID_SUB="16736042990918707269" BLOCK_SIZE="512" TYPE="zfs_member"`

解决(未验证):
1. `zpool labelclear /dev/vdc`
1. `wipefs -a /dev/vdc`

### zpool status解析
1. v2.3: 支持json输出
2. v2.3之前: 每行以`\t`开头未换行内容, 否则为field(比如pool, state, scan, remove, config, errors等).

### 获取dataset iostat
`/proc/spl/kstat/zfs/<pool>/objset-xxx`, 里面包含dataset_name

### special规划
ZFS 并没有硬性要求 special vdev 的容量，它的容量需求主要取决于以下数据：
1. 元数据（文件系统元数据、ZIL 指针等）
1. 小文件（small_blocks，通常 ≤ recordsize，默认 ≤ 128KB）
1. Deduplication 表（如果启用去重）
1. 其他特殊数据（如 resilver 状态）但以下是建议和实战经验总结

ps: 如果 special 空间不足，ZFS 会回退到普通存储，导致性能下降

建议和实战经验总结:
目标	建议容量比例
基本使用	至少为主池容量的 0.5%~2%
小文件为主的工作负载	2%~5% 或更多（甚至到 10%）
所有元数据和小数据块都放入 special vdev（完全 offload）	足够容纳所有元数据 + 小块数据的总和

实际容量大小依据：
1. 元数据量

	- 多数 ZFS 文件系统元数据体积大概在主池使用量的 1~2%，但具体看数据类型和文件数。

- Small block threshold

	默认情况下，小于 128 KiB 的块可以被写入 special vdev。

	这个阈值可以通过 special_small_blocks 属性调整，比如设为 64K

- 使用估算

	可以使用如下命令查看已有的元数据和小块使用量估算`zdb -Lbbbs <poolname>`, 输出中的 metadata 和 small blocks 使用量可以作为 special vdev 规划依据

### zpool移除mirror special报`all top-level vdevs must have the same sector size and not be raidz`
ref:
- [zpool cannot remove vdev #14312](https://github.com/openzfs/zfs/issues/14312#issuecomment-2296222799)

已检查`zdb -C <pool>`的ashift和`lsblk -o NAME,PHY-SEC,LOG-SEC`, 没有问题. 应该是top-level vdevs存在raidz导致的