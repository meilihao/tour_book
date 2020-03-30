# lvm
参考:
- [LVM的基本概念和部署](http://xintq.net/2014/07/30/LVM%E7%9A%84%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5%E5%92%8C%E9%83%A8%E7%BD%B2/)

## 描述
一种硬盘设备资源管理技术, 允许用户对硬盘资源进行动态调整.

物理卷处于 LVM 中的最底层，可以是物理硬盘、硬盘分区或者 RAID 磁盘阵列.

物理卷(PV，Physical Volume)由基本单元(PE，Physical Extent, 大小默认为 4MB)组成-> pv构成卷组(VG，Volume Group)-> vg可切分成逻辑(LV，Logical Volume）即逻辑磁盘, lv是PE的倍数.

## 常用的 LVM 部署命令
功能/命令 物理卷管理 卷组管理 逻辑卷管理
扫描 pvscan vgscan lvscan
建立 pvcreate vgcreate lvcreate
显示 pvdisplay vgdisplay lvdisplay
删除 pvremove vgremove lvremove
扩展 - vgextend lvextend
缩小 - vgreduce lvreduce

LVM 还具备有`快照卷`功能，该功能类似于虚拟机软件的还原时间点功能. LVM 的快照卷功能有两个特点：
- 快照卷的容量必须等同于逻辑卷的容量；
- 快照卷仅一次有效，一旦执行还原操作后则会被立即自动删除

### example
```bash
# pvcreate /dev/sdb /dev/sdc # 创建vg
# vgcreate storage /dev/sdb /dev/sdc # # vg添加磁盘
# lvcreate -n vo -l 37 storage # 用37个pe创建lv
# mkfs.ext4 /dev/storage/vo
# lvextend -L 290M /dev/storage/vo # lv扩容, 但需先umount
# e2fsck -f /dev/storage/vo # 检查文件系统的完整性，并重置硬盘容量
# resize2fs /dev/storage/vo 120M # 把逻辑卷 vo 的容量减小到 120MB, 有数据丢失风险, 且缩容前必须检查文件系统的完整性
# lvcreate -L 120M -s -n SNAP /dev/storage/vo # 使用-s 参数生成一个快照卷，使用-L 参数指定切割的大小, SNAP为快照卷名称
# lvconvert --merge /dev/storage/SNAP # 要对逻辑卷进行快照还原(SNAP包含原卷信息, 因此不用指定原卷), 需先umount原卷.
# lvremove /dev/storage/vo # 删除lv, 但需先umount
# vgremove storage # 删除vg
# pvremove /dev/sdb /dev/sdc # 删除物理卷设备 
```