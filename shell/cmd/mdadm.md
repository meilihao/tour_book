# mdadm
linux下用于构建、管理和监控RAID阵列(也被称为md设备)的工具.

在创建好RAID以后，可以将RAID信息保存到 `/etc/mdadm/mdadm.conf` 文件中，这样在下次操作系统重新启动时，系统会自动加载这个文件来启用RAID.
可以通过 `/proc/mdstat` 文件来查看RAID的简洁信息 

## 格式

    mdadm [模式] [设备] [选项]

## 选项

- -a yes : 自动在/dev/下创建对应的RAID阵列设备
- -l N : 指定RAID级别为N
- -n N : 指定硬盘数量
- --detail : 查看RAID的详细信息 
- -C : 创建
- -v : 显示过程
- -f dev : 模拟设备损坏
- -r : 移除设备
- -Q : 查看摘要信息
- -D : 查看详细信息, 包括重建状态(Rebuild Status)
- -S : 停止 RAID 磁盘阵列
- -x N : 表示有几块备用盘

## 例

    # sudo mdadm --create /dev/md0 -a yes -l 5 -n 3 -x 1 /dev/sdb /dev/sdc /dev/sdd /dev/sde # 创建一个RAID5, 且有一块备用盘
    # mdadm --detail /dev/md0 # 查看RAID的详细信息
    # sudo mdadm --misc --stop /dev/md0 # 关闭RAID
    # sudo mdadm --manage /dev/md0 --remove /dev/sdd # 移除RAID阵列中的磁盘
    # sudo mdadm --manage /dev/md0 --add /dev/sdd # 追加磁盘
    # echo "/dev/md0 /RAID ext4 defaults 0 0" >> /etc/fstab # 并把挂载信息写入到配置文件中，使其永久生效
    # mdadm /dev/md0 -f /dev/sdb # 模拟坏盘
    # mdadm /dev/md0 -r /dev/sdb # 移除坏盘
    # mdadm /dev/md0 -a /dev/sdb # 将新盘换入md0(需先卸载再重新挂载)

## FAQ
### raid 重建进度
方法:
1. `cat /proc/mdstat`
1. `mdadm -D /dev/md<N>`的`Resync Status`

# [dmsetup](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/logical_volume_manager_administration/device_mapper#dm-mappings)
在 Linux 系统中，使用 Device Mapper (DM) 可以创建虚拟块设备，并将其用于分区或其他目的.

## 选项
ref:
- [TABLE FORMAT](https://man7.org/linux/man-pages/man8/dmsetup.8.html)

- --table : `<start_sector> <num_sectors> <target_type> <target_args>/start length mapping [mapping_parameters...]`

    - start_sector：起始扇区号: 在设备映射器表的第一行中，start 参数必须等于 0. 一行中的 start + length 参数必须与下一行中的 start 相同

        ps: target_type=linear时, 如果start一直为0, 则是通过target_args.offeset + num_sectors自动推算的
    - num_sectors：扇区数量 : 设备映射器中的大小总是以扇区（512 字节）指定
    - target_type：目标类型（如 linear、striped、error 等）
    - target_args：目标参数（如物理设备路径和偏移量）, 即在映射表的行中指定哪个映射参数取决于行中指定 mapping 类型

## example
```bash
# dmsetup create my_device --size 10G --table "0 20971520 linear /dev/sda 0" # my_device, 将创建的设备的名字; --size 10G 指定创建的设备大小; --table 定义了设备的映射表. `0 20971520 linear /dev/sda 0`表示将 /dev/sda 的前 10GB（20971520 块，每块 512 字节）映射到 my_device
# dmsetup deps : 显示设备的依赖关系
# dmsetup table  : 显示当前设备的table信息
# dmsetup remove my_device # 删除 Device Mapper 设备
# dmsetup ls # 列出所有当前活动的 Device Mapper 设备
```