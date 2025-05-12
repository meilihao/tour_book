# btrfs
参考:
- [boot loader对btrfs的支持](https://wiki.archlinux.org/index.php/Arch_boot_process)
- [^btrfs 使用指南 - 1 概念，创建，块设备管理，性能优化](https://developer.aliyun.com/article/50)

## 安装
`apt install btrfs-progs`

新创建的 Btrfs 文件系统会创建一个`subvol=/`的默认 Subvolume，即 Root Subvolume. 其 ID 为 5 (别名为 0), 这是一个 ID 和目录都预设好的 Subvolume, 可用`mount|grep btrfs`查看.

btrfs可以在 Subvolume 的基础上制作快照，几点需要注意：
- 默认情况下 Subvolume 的快照是可写的
- 快照是特殊的 Subvolume，具有 Subvolume 的属性。所以快照也可以通过 mount 挂载，也可以通过 btrfs property 命令设置只读属性
- 由于快照的本质就是一个 Subvolume ，所以可以在快照上面再做快照

## example
> /mnt/bfs是btrfs的mountpoint

```bash
# btrfs subvolume create sub1 # 必须在btrfs mountpoint下执行. 属于不同的 Subvolume 间的文件不能建立硬链接, 且rm命令无法删除subvolume
# btrfs subvolume del [-c] sub1 # 或指定id(`--subvolid 263 <mountpoint>`). 删除subvolume # 返回信息中的`Delete subvolume (no-commit)`，表示 Subvolume 被删除了，但没有提交, 即删除操作在内存里面生效了，但磁盘上的内容还没删，意味着如果这个时候系统 Crash 掉，这个 Subvolume 有可能还会回来. Btrfs 这样做的好处是删除速度很快，不会影响使用，缺点是有可能在后台 Commit 的过程中系统挂掉，导致 Commit 失败. 在删除 Subvolume 的时候指定 -c 参数，这样 btrfs命令会等提交完成之后再返回, 否则就需要`btrfs subvolume sync <mountpoint>`
# btrfs subvolume list /mnt/bfs # 获取subvolume的列表. `-s`表示只显示snapshot
# btrfs property get -ts /mnt/bfs/sub1 # 查看 Subvolume 的只读状态
# btrfs property set -ts /mnt/bfs/sub1/ ro true # 设置 Subvolume 的只读状态
# btrfs subvolume snapshot /mnt/bfs /mnt/bfs/snap-rootCreate # 创建快照. 默认可写, 除非使用`-r`指定只读
# btrfs subvolume del snap-rootDelete # 删除快照和删除 Subvolume 是一样的，即没有commit.
# btrfs subvolume set-default 256 /mnt/bfs # 设置默认subvolume, 256是其id. Btrfs 分区的默认 Subvolume，即在挂载磁盘的时候，可以只让分区中的指定 Subvolume 对用户可见.
# ---
# btrfs inspect-internal dump-tree /dev/sdm # 获取btree情况
# btrfs filesystem show /dev/sdb # 获取有关文件系统的详细信息
# btrfs filesystem df /mnt/bfs # 获取有关数据和元数据的更多详细信息
# btrfs filesystem resize -2g /mnt/bfs # 减少2GB
# btrfs filesystem resize +1g /mnt/bfs # 增加1GB
# btrfs filesystem resize max /mnt/bfs # 增加到最大可用空间
# btrfs device add /dev/sde /mnt/bfs # 添加/dev/sde到btrfs文件系统
# btrfs filesystem balance /mnt/bfs # 执行文件系统平衡，以便数据和元数据分布在所有设备上
# btrfs device delete /dev/sdc /mnt # 从btrfs文件系统在线删除完整的硬盘驱动器
# btrfs filesystem usage /mnt/bfs # 像 df 这样的用户空间工具可能不会准确的计算剩余空间 (由于并无分别计算文件和元数据的使用状况) . 推荐使用 btrfs filesystem usage 来查看使用状况
# btrfs filesystem defragment -r /mnt/bfs # 不带`-r`将致使仅整理该目录的子卷所拥有的元数据
# btrfs scrub start /mnt/bfs # 启动一个（后台运行的）对 /mnt/bfs 目录的文件系统的在线检查任务
# btrfs scrub status /mnt/bfs # 查看上诉任务的运行状态
```

btrfs 其他常用命令如下：
- btrfs 文件系统属性查看：btrfs filesystem show
- 调整文件系统大小：btrfs filesystem resize +10g MOUNT_POINT
- 添加硬件设备：btrfs filesystem add DEVICE MOUNT_POINT
- 均衡文件负载：btrfs blance status|start|pause|resume|cancel MOUNT_POINT
- 移除物理卷(联机、自动移动)：btrfs device delete DEVICE MOUNT_POINT
- 动态调整数据存放机制：btrfs balance start -dconvert=RAID MOUNT_POINT
- 动态调整元数据存放机制：btrfs balance start -mconvert=RAID MOUNT_POINT
- 动态调整文件系统数据数据存放机制：btrfs balance start -sconvert=RAID MOUNT_POINT
- 创建子卷：btrfs subvolume create MOUNT_POINT/DIR
- 列出所有子卷：btrfs subvolume list MOUNT_POINT
- 显示子卷详细信息：btrfs subvolume show MOUNT_POINT
- 删除子卷：btrfs subvolume delete MOUNT_POIN/DIR
- 创建子卷快照(子卷快照必须存放与当前子卷的同一父卷中)：btrfs subvolume snapshot SUBVOL PARVOL
- 删除快照同删除子卷一样：btrfs subvolume delete MOUNT_POIN/DIR
- 修改fs uuid: btrfstune -U $(uuidgen) /dev/sdb1

## snap
ref:
- [Btrfs 详解：快照](https://linux.cn/article-16287-1.html)
- [使用 Btrfs 快照进行增量备份](https://linux.cn/article-12653-1.html)

## 进阶
ref:
- [文件系统系列专题之 Btrfs](https://blog.csdn.net/zhuzongpeng/article/details/127115533)

## FAQ
### fstab
```bash
# cat /etc/fstab
# <file system> <mount point>   <type>  <options>       <dump>  <pass>
# / was on /dev/nvme0n1p3 during installation
UUID=ca75d517-3282-4a44-8f1b-5f19937034a5 /               btrfs   defaults,subvol=@rootfs 0       0
UUID=ca75d517-3282-4a44-8f1b-5f19937034a5 /home           btrfs   defaults,subvol=home    0       0
UUID=ca75d517-3282-4a44-8f1b-5f19937034a5 /opt            btrfs   defaults,subvol=opt     0       0
# /boot was on /dev/nvme0n1p2 during installation
UUID=2af5d19c-6376-47c6-b81d-9a4826f0069b /boot           ext4    defaults        0       2
# /boot/efi was on /dev/nvme0n1p1 during installation
UUID=3CA2-D74B  /boot/efi       vfat    umask=0077      0       1
```

### 优化
autodefrag+定期平衡=性能不滑坡