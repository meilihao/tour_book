# lvm
参考:
- [LVM的基本概念和部署](http://xintq.net/2014/07/30/LVM%E7%9A%84%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5%E5%92%8C%E9%83%A8%E7%BD%B2/)

## 描述
一种硬盘设备资源管理技术, 允许用户对硬盘资源进行动态调整.

## 概念
- PV：physical volume 物理卷在逻辑卷管理系统最底层，可为整个物理硬盘或实际物理硬盘上的分区
- VG：volume group 卷组建立在物理卷上，一卷组中至少要包括一物理卷，卷组建立后可动态的添加卷到卷组中，类似磁盘的扩展分区，本身不能创建文件系统，需再划分成LV来使用，一个逻辑卷管理系统工程中可有多个卷组
- LV：logical volume 逻辑卷建立在卷组基础上，卷组中未分配空间可用于建立新的逻辑卷，逻辑卷建立后可以动态扩展和缩小空间
- PE：physical extent 物理区域是物理卷中可用于分配的最小存储单元，物理区域大小在建立卷组时指定，一旦确定不能更改，同一卷组所有物理卷的物理区域大小需一致，新PV加入到VG后，PE大小自动更改为VG中定义的PE大小

物理卷处于 LVM 中的最底层，可以是物理硬盘、硬盘分区或者 RAID 磁盘阵列.

物理卷(PV，Physical Volume)由基本单元(PE，Physical Extent, 大小默认为 4MB)组成-> pv构成卷组(VG，Volume Group)-> vg可切分成逻辑(LV，Logical Volume）即逻辑磁盘, lv是PE的倍数.

## 常用的 LVM 部署命令
PV管理工具

- pvs: 简要显示物理卷信息

    ```bash
    # pvs --reportformat json
    ```
- pvdisplay: 显示物理卷详细信息
- pvcreate: 创建物理卷
- pvremove: 移除物理卷

VG管理工具

- vgs: 简要显示卷组信息


    ```bash
    # vgs --reportformat json
    ```
- vgdisplay: 显示卷组详细信息
- vgcreate: 创建卷组
- vgextend: 扩展卷组
- vgreduce: 缩小卷组
- vgremove：删除卷组

LV管理工具

- lvs: 简要显示逻辑卷信息

    ```bash
    # lvs --reportformat json
    ```
- lvdisplay: 显示逻辑卷详细信息
- lvcreate: 创建逻辑卷
        
    - -l ： 指定逻辑卷的大小（LE数）
    - -L: 大小[mMgGtT]
    - -n: 指定创建卷名
    - -s: 指定创建为快照
    - -p: 权限[r|rw],默认rw
- lvextend: 扩展逻辑卷

    - -L: 指定大小
- resize2fs: 调整文件系统
- lvremove: 删除逻辑卷

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
# pvcreate /dev/sdb /dev/sdc # 创建pv
# vgcreate storage /dev/sdb /dev/sdc # 创建vg, 将pv加入vg
# lvcreate -L 20G --thinpool Data_Pool storage # 创建thin pool, 它是建立在vg上
# lvcreate -V 10G --thin -n thin_LV_data01 storage/Data_Pool # 创建thin volume
# lvcreate -n vo -l 37 storage # 用37个pe创建lv
# lvdisplay /dev/mapper/vo # 或 lvdisplay /dev/storage/vo
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

lvm可用`ls /dev/mapper/*`查看, 用`dmsetup remove /dev/dm-2`删除, [可用`vgchange -ay <vg>`重新激活](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/configuring_and_managing_logical_volumes/assembly_lvm-activation-configuring-and-managing-logical-volumes).

lvm会以`/dev/<vg>/<lv>`形式生成块设备的软连接, 这与zvol类似.

## vg配置
在`/etc/lvm/backup/<vg_name>`里.

## FAQ
### 从lv名称中获取vg名称
lv名称的第一个`-`前的字符为vg名称