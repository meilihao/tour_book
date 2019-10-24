# zfs
参考:
- [在 Oracle® Solaris 11.2 中管理 zfs 文件系统](https://docs.oracle.com/cd/E56344_01/html/E53918/index.html)
- [zfs Administration](https://pthree.org/2012/12/04/zfs-administration-part-i-vdevs/)
- [FreeBSD Handbook's zfs](https://www.freebsd.org/doc/handbook/)

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
> zfs已实现零管理, 通过zfs daemon(zed,zfs Event Daemon)实现, 因此无需将pool信息写入`/etc/fstab`, pool配置在`/etc/zfs/zpool.cache`里.
> zfs pool使用zfs_vdev_scheduler来调度io.

## 概念
pool : 存储设备的逻辑分组, 它是ZFS的基本构建块，可将其存储空间分配给数据集.
dataset : zfs文件系统的组件即文件系统、克隆、快照和卷被统称为数据集
mirror : 一个虚拟设备存储相同的两个或两个以上的磁盘上的数据副本，在一个磁盘失败的情况下,相同的数据是可以用其他磁盘上的镜像.
resilvering ：在恢复设备时将数据从一个磁盘复制到另一个磁盘的过程
snapshot : 快照, 是文件系统或卷的只读副本. 在zfs中，快照几乎可以即时创建，而且最初不会额外占用池中的磁盘空间
scrub : 用于一致性检验. 其他文件系统会使用fsck.

### 建议
1. 使用整个磁盘来启用磁盘写入高速缓存并使维护更轻松
1. 使用zfs冗余(raidz,mirror)

   1. mirror, 使用镜像磁盘对
   1. raidz, 为每个vdev组合3-9个磁盘
   1. 不要在同一个pool中混用raidz和mirror, 因为会难以管理且性能受到影响
1. 避免在存储池中使用磁盘分片而是应使用整块磁盘, 已避免潜在复杂性.
1. 使用热备以减少硬件故障而导致的停机时间
1. 使用相同大小的磁盘
1. pool最小为8G.

### zfs虚拟设备(zfs vdevs)
一个 VDEV 是一个meta-device，代表着一个或多个物理设备. zfs 支持 7 中不同类型的 VDEV：
- disk, 默认, 比如HDD, SDD, PCIe NVME等等
- File : 预先分类的文件，为*.img的文件，可以作为一个虚拟设备载入zfs
- Mirror : 标准的 RAID1 mirror
- zfs 软件RAID : raidz=raidz1(raid5, >=3 disk)/2(raid6, >=4 disk)/3(>=5 disk), 非标准的基于分布式奇偶校验的软件RAID. 速度: raid0 > raid1 > raidz1 > raidz2 > raidz3
- Hot Spare : 用于热备 zfs 的软件 raid
- Cache : 用于2级自适应的读缓存的设备 (zfs L2ARC), 提供在 memory 和 disk的缓冲层, 用于改善静态数据的随机读写性能
- Log : zfs Intent Log(zfs ZIL/SLOG, zfs意图日志,一种对于 data 和 metadata 的日志机制，先写入然后再刷新为写事务), 用于崩溃恢复, 最好配置并使用快速的 SSD来存储ZIL, 以获得更佳性能. ZIL支持mirror. ZIL也可认为是zfs的写缓存.

VDEV始终是动态条带化的. 一个 device 可以被加到 VDEV, 但是不能移除.

## zpool
zfs支持分层组织filesystem, 每个filesystem仅有一个父级, 而且支持属性继承, 根为池名称.

> zpool描述了存储的物理特性, 因此必须在创建filesystem前创建池.
> 如果pool没有配置log device, pool会自行为ZIL预留空间.
> raid-z配置无法附加其他磁盘; 无法分离磁盘, 但将磁盘替换为备用磁盘或需要分离备用磁盘时除外; 无法移除除log device或cache device外的device
> 创建pool时, 不能使用其他pool的组件(vdev, 文件系统或卷), 否则会造成死锁.

```sh
$ sudo zpool create pool-test /dev/sdb /dev/sdc /dev/sdd # 创建了一个零冗余的RAID-0存储池, zfs 会在`/`中创建一个目录,目录名是pool name 
$ sudo zpool [option] list # 显示系统上pools的列表, `-o`只显示指定列,`-H`隐藏列头
$ sudo zpool status <pool> # 查看pool的状态,read/write列显示io错误次数, cksum列显示无法更正的校验和错误的次数. `-v`输出详细信息, `-x`仅显示有错误或因其他原因不可用的pool
$ sudo zpool destroy <pool> # 销毁pool
$ sudo zpool destroy <pool>/data-set # 销毁dataset
$ sudo zpool upgrade [<pool> | -a] # 更新 zfs 时，就需要更新指定/全部池
$ sudo zpool add <pool> /dev/sdx # 将驱动器添加到池中
$ sudo zpool exprot <pool> # 如果要在另一个系统上使用该存储池，则首先需要将其导出. zpool命令将拒绝导入任何尚未导出的存储池
$ sudo zpool export oldname # 重命名已创建的zpool的过程分为2个步骤, export + import
$ sudo zpool import oldname newname
$ sudo zpool create <pool> mirror <device-id-1> <device-id-m1> mirror <device-id-2> <device-id-m2> # 创建RAID10
$ sudo zpool add <pool> log mirror <device-id-1> <device-id-2> # 添加 SLOG
$ sudo zpool add <pool> spare devices # 添加热备, 大小应>=max(pool's vdev),且无法移除当前正在使用的热备. 移除用`zpool remove`
$ sudo zpool add <pool> cache <device-id> # 添加L2ARC
$ sudo zpool iostat -v <pool> N # 每隔N秒输出一次pool的io状态
$ sudo zpool remove <pool> mirror-1 # 移除mirror/不在使用的热备
$ sudo zpool attach <pool> <existing-device> <new-device> # 将新设备追加到已有vdev
$ sudo zpool detach  # 分离设备, 对象必须是mirror中的设备/raidz中已由其他物理设备或备用设备替换的设备
$ sudo zpool split <pool> <new-pool> [device] # 拆分pool, 仅适用mirror设备, 通过`-R`可指定新池的挂载点
$ sudo zpool offline [option] <pool> <device> # 离线zfs设备, `-t`表示临时离线, 重启后会重新恢复到online.
$ sudo zpool online [option] <pool> <device> # 上线zfs设备, 新设备上线后会同步. `-e`可扩展LUN(即使用更大容量设备时, 使用完整大小), 默认不扩展.
$ sudo zpool clear <pool> [devices] # 池设备故障时清理错误, 默认清理池内的所有设备错误.
$ sudo zpool replace <pool> replaced-device [new-device] # 替换存储池中的设备, 比如使用热备盘
$ sudo zpool get all <pool> # 获取pool的所有属性
$ sudo zpool get <property> <pool> # 获取pool的指定属性
$ sudo zpool set <property=value> <pool> # 设置pool的指定属性
$ sudo zpool history [pool] # 显示zfs和zpool命令的使用日志, `-l`使用长格式(追加: 用户名, 主机名执行操作的域), `-i`显示更详细的内部日志可用于诊断
$ sudo zpool iostat # 列出pools的io统计信息, `-v`显示pool包含的vdev的io统计信息, `-l`更多信息
$ sudo zpool export [option] <pool> # 导出pool, 该pool必须处于不活动状态, 导出后该pool在系统中不可见. `-f`强制取消已挂载的filesystem
$ sudo zpool import [option] [pool/id-number] # 导入pool, 导入时允许重命名,允许只读导入. `-m`表示导入mirror log, 默认不导入.`-d`导入非标准路径的设备/由file构成的pool.`-D`恢复已销毁的pool, 追加`-f`即`-Df`可表示已销毁的pool中某设备不可用也可恢复.
$ sudo zpool upgrade # 升级pool, 以使用新版的zfs功能. `-v`表示当前zfs支持的功能, `-a`表示升级到最新的zfs.
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
创建pool: `zpool create -f -m <mount> <pool> [raidz（2 | 3）| mirror] <ids>`

参数:
- f : 强制创建pool, 用于解决"EFI标签"错误
- m : 指定挂载点, 默认是`/`. 如果挂载点目录存在数据会保存
- o : 设置kv属性, 比如ashift
- n : 仅模拟运行, 并输出日志, 用于检查错误, 但不是100%正确, 比如在配置中使用了相同的设备做mirror.

> `blockdev --getpbsz /dev/sdXY`获取扇区大小, zfs默认是512, 与设备扇区不一致会导致性能下降, 比如具有4KB扇区大小的高级格式磁盘，建议将ashift更改为12(2^12=4096). 更改[ashift](https://github.com/zfsonlinux/zfs/wiki/faq#advanced-format-disks)选项的唯一方法是重新创建池.

### zpool scrub
只要读取数据或zfs遇到错误，`zpool scrub`就会在可能的情况下进行静默修复.

```sh
$ zpool scrub <pool> # 检修 pool
$ zpool scrub -s <pool> # 取消正在运行的检修
```

> 建议是每周/月检修一次.
> 替换磁盘而同步数据很耗时, 替换间执行`zpool scrub`有利于替换设备运行正常且数据写入正确.

## zfs
```sh
$ sudo zfs list # 显示系统上pools/filesystems的列表, `-r`递归显示fs及其子fs, `-o`指定要显示的属性; `-t`指定显示的类型, 比如filesystem, volume, share, snapshot.`-H`表示脚本模式: 不输出表头并用单个tab分隔各列; `-p`:精确显示数值; `-d`: 与`-r`连用,限制递归深度; `-s`按指定列升序排序; `-S`:与`-s`类似, 但以降序排序.
$ sudo zfs get [ all | property[,property]...] <pool> # 获取pool的参数. `-s`指定要显示的source类型; `-H`输出信息去掉标题, 并用tab代替空格来分隔
$ sudo zfs set atime = off <pool> # 设置pool参数
$ sudo zfs set compression=gzip-9 mypool # 设置压缩的级别
$ sudo zfs inherit -rS atime  <pool> # 重置参数到default值. `-r`以递归的方式应用inherit子命令
$ sudo zfs get keylocation <pool>/<filesystem> # 获取filesystem属性
$ sudo zfs set acltype = posixacl <pool> / <filesystem> # 使用ACL
$ sudo zfs set sharenfs=on <pool> # 通过nfs共享pool
$ sudo zfs create -o mountpoint=none mypool/test/storage # 创建未挂载的dataset, 常用于zfs recv的场景.
$ sudo zfs set sharenfs=on <pool>/<filesystem> # 通过nfs共享filesystem
$ sudo zfs set mountpoint=/<pool>/... <pool>/... # 设置挂载点, 设置后会立即挂载.
$ sudo zfs destroy <pool>/.../<filesystem> # 销毁文件系统, 此时fs必须是不活动的. `-r`表示递归销毁, `-R`表示递归销毁这些快照及其clone, `-d`销毁带保持标志的快照
$ sudo zfs rename <old-path> <new-path> # 重命名fs
$ sudo mount -o <pool>/.../<filesystem> # 挂载fs
$ sudo unmount <pool>/.../<filesystem> # 取消挂载fs, 此时fs必须是不活动的. `-f`强制取消挂载
$ sudo zfs create -V 5gb system1/vol # 创建5g大小的卷(创建卷时，会自动将预留空间设置为卷的初始大小，以确保数据完整性)
```

zfs list的属性可参考[freebsd zfs#Native Properties](https://www.freebsd.org/cgi/man.cgi?zfs(8)), `man zfs`与实际zfs包含的功能不一致, 有缺失, 比如`createtxg`, 部分属性说明:
- guid : 对象的全局唯一标识符（GUID)
- createtxg : 创建时的事务id. bookmark与snapshot有相同的createtxg. 该属性常用于`zfs send/recv`
- creation : 创建时间, unix时间戳

> 物理卷(Physical Volume, PV)：操作系统识别到的物理磁盘(或者RAID提交的逻辑磁盘LUN), 物理卷可以是一个磁盘，也可以是磁盘中的一个分区.
> volume : 通常是指逻辑卷, 是逻辑卷组(VG, 由若干PV组成)上的一块空间, 上面没有文件系统.

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

### zfs snapshot
zfs snapshot（快照）是 zfs 文件系统或卷的**只读**拷贝(即无法修改属性). 可以用于保存 zfs 文件系统的特定时刻的状态，在以后可以用于恢复该快照，并回滚到备份时的状态.

> snap直接占用pool, 不使用后备存储.
> 快照放在对应fs的`.zfs/snapshot`目录里.

```sh
$ sudo zfs snapshot -r mypool/projects@snap1 # 创建 mypool/projects 文件系统的快照. `-r`表示递归创建(即为所有后代文件系统创建快照)
$ sudo zfs list -t snapshot # 查看所有的snapshots列表
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

> 创建clone几乎是即时的, 且最初不占用其他磁盘空间. 此外, 还可以对clone进行快照.
> 可用`zfs get/set`操作clone的属性.

```sh
# 克隆 mypool/projects，首先创建一个 snapshot 然后 clone
$ sudo zfs snapshot mypool/projects@snap1
$ sudo zfs clone mypool/projects@snap1 mypool/projects-clone
$ sudo destroy mypool/projects-clone # 销毁clone
```

### zfs send 和 receive
zfs send 将文件系统的快照写入stdout，然后流式传送到文件或其他机器. zfs receive 从stdin接收该 stream 然后写到 snapshot 拷贝，作为 zfs 文件系统. 这有利于通过拷贝或网络传送备份.

> zfs send默认是完整流(从fs创建开始到指定的快照为止的所有内容).
> 增量流: 包含一个快照与另一个快照之间的差异.

```sh
# 创建 snapshot 然后 save 到文件
$ sudo zfs snapshot -r mypool/projects@snap2
$ sudo zfs send mypool/projects@snap2 > ~/projects-snap.zfs  # `-c`使用压缩(如果mypool/projects是活动的则必须使用), `-n`表示模拟send, 实际不产生数据流, `-P`表示生成流的信息, 比如全量/增量, 数据流大小.`-v`: 向(stderr)发送流的详细信息, 包括每秒传输多少.
$ sudo zfs receive -F mypool/projects-copy < ~/projects-snap.zfs # 恢复,此时目标fs必须存在. `-F`表示(此时目标必须没有快照)忽略目标fs的改动(mypool/projects-copy), 全量的话是直接覆盖原有fs, 增量的话是回滚到该增量快照的起点后再应用增量. `-d`: (此时目标fs必须存在)去掉原快照名称中的pool name,使用目标fs name+剩余名称作为新名称.
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
zfs dedup 将丢弃重复数据块，并以到现有数据块的引用来替代. 这将节约磁盘空间，这**需要大量的内存**. 内存中的去重记录表需要消耗大约 ~320 bytes/block. 表格尺寸越大，写入时就会越慢.

```
$ sudo zfs set dedup=on mypool/projects # 启用去重
```

## FAQ
### quota于refquota区别
如果对 tank/home 数据集设置了quota，则 tank/home 及其所有后代使用的总磁盘空间量不能超过该配额.
如果对 tank/home 数据集设置了refquota，则 tank/home 磁盘空间量不能超过该配额, 但不包括后代所占用的空间, 比如其快照/clone.

### reservation
如果 tank/home 指定了预留空间，则 tank/home 及其所有后代都会使用该预留空间. 预留空间计入父文件系统的已用磁盘空间, 并计入父文件系统的配额和预留空间.

### 文件读取很慢
原因:
1. pool的一块镜像盘的状态变成DEGRADED, 更换磁盘即可