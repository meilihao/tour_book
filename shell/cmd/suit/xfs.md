# xfs
XFS是一个开源的（GPL）日志文件系统，它已被CentOS/RHEL 7采用，成为其默认的文件系统. 它支持在XFS文件系统已挂载的情况下可以进行扩展, 但不支持缩减.

## xfs_info
查看xfs fs的详细信息

## xfs_growfs
xfs fs扩容

```bash
# xfs_growfs /dev/centos/root -D 1986208 # `-D`指定大小(blocks)，否则xfs_growfs将会自动扩展XFS文件系统到最大的可用大小
```

## 其他命令
参考:
- [修改云盘的UUID](https://help.aliyun.com/document_detail/199739.html)

```bash
- xfs_admin: 调整 xfs 文件系统的各种参数
- xfs_copy: 拷贝 xfs 文件系统的内容到一个或多个目标系统（并行方式）
- xfs_db: 调试或检测 xfs 文件系统（查看文件系统碎片等）
- xfs_check: 检测 xfs 文件系统的完整性
- xfs_bmap: 查看一个文件的块映射
- xfs_repair: 尝试修复受损的 xfs 文件系统
- xfs_fsr: 碎片整理
- xfs_quota: 管理 xfs 文件系统的磁盘配额
- xfs_metadump: 将 xfs 文件系统的元数据 (metadata) 拷贝到一个文件中
- xfs_mdrestore: 从一个文件中将元数据 (metadata) 恢复到 xfs 文件系统
- xfs_freeze : 暂停（-f）和恢复（-u）xfs 文件系统
- xfs_logprint: 打印xfs文件系统的日志
- xfs_mkfile: 创建xfs文件系统
- xfs_ncheck: generate pathnames from i-numbers for XFS
- xfs_rtcp: XFS实时拷贝命令 
- xfs_io: 调试xfs I/O路径
```

```bash
# xfs_admin -U generate /dev/sdb7 # 使用xfs_admin为分区生成新的uuid. 其他方法: 1. `tune2fs -U c1b9d5a2-f162-11cf-9ece-0020afc76f16 /dev/sda5` 2. `uuidgen | xargs tune2fs /dev/sda5 -U`
# xfs_admin -lu /dev/mapper/vg_test-lv_test # 查看lable和uuid
# xfs_repair -n block-device # `-n`仅探测不修复
# xfs_repair [-L] /dev/sda # `-L`通常是xfs已损坏时强制修复使用, 是修复xfs文件系统的最后手段, 因为它会清空日志, 会丢失用户数据和文件
```

## FAQ
### mount报`Superblock has unknown read-only compatible features (0x8) enabled`
ref:
- [Mounting XFS created by RHEL9 fails in RHEL8 even in read-only mode](https://bugzilla.redhat.com/show_bug.cgi?id=2046431)

完整日志:
```log
[  170.593306] XFS (dm-4): Superblock has unknown incompatible features (0x8) enabled.
[  170.593309] XFS (dm-4): Filesystem cannot be safely mounted by this kernel.
[  170.593315] XFS (dm-4): SB validate failed with error -22.
```

### 定位features
在xfsprogs源码:
1. `grep -ri "features"`

    发现`mp->m_sb.sb_features_ro_compat |= XFS_SB_FEAT_RO_COMPAT_INOBTCNT`
1. `grep -ri "XFS_SB_FEAT_RO_COMPAT_INOBTCNT"`

    发现`libxfs/xfs_format.h:#define XFS_SB_FEAT_RO_COMPAT_INOBTCNT (1 << 3)`