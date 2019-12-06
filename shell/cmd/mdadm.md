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

## 例

    # sudo mdadm --create /dev/md0 -a yes -l 5 -n 3 /dev/sdb /dev/sdc /dev/sdd # 创建一个RAID5
    # mdadm --detail /dev/md0 # 查看RAID的详细信息
    # sudo mdadm --misc --stop /dev/md0 # 关闭RAID
    # sudo mdadm --manage /dev/md0 --remove /dev/sdd # 移除RAID阵列中的磁盘
    # sudo mdadm --manage /dev/md0 --add /dev/sdd # 追加磁盘