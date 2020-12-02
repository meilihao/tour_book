# mkfs
创建文件系统

## 格式

    mkfs [-V] [-t fstype] Device

## 选项
- -t <fstype> : 系统自动寻找`mkfs.fstype`来执行, 比如`mkfs.ext4`

## FAQ
### mke2fs
mke2fs是专门用于管理ext系列文件系统的一个专门的工具.

### e2fsprogs
e2fsprogs是一个ext系列文件系统工具（Ext2Filesystems Utilities），它包含了诸如创建、修复、配置、调试ext系列文件系统的标准工具.

### 查询当前系统支持的fs类型/mke2fs.conf文件中没有定义类型 xfs 的文件系统
```bash
# cat /proc/filesystems |grep -i xfs # 检查kernel是否支持xfs
```