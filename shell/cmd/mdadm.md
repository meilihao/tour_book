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
- -D : 查看详细信息
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
    # mdadm /dev/md0 -a /dev/sdb # 将新盘换入md0(需先卸载再重新挂载)