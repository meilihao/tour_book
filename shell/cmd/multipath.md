# multipath
设备映射器多路径（ DM-Multipath） 可让将服务器节点和存储阵列间的多个 I/O 路径（这些路径在服务器中的主机总线适配器和设备储存控制器之间）配置为一个单一设备. 通常是在光纤通道 (FC) 或 iSCSI SAN 环境中.

> 本质: 同一时刻能通过多条路径访问到远端的同一个块设备, 因此multipath要求块设备的wwpn必须相同. 对于fc是存在多条光纤链路; 对于iscsi是target server存在多个ip, 且iscis target有多条portal或使用了通配符portal.

```sh
# yum install device-mapper-multipath -y
# apt install multipath-tools -y
# 配成开机启动
# 真正配置前需停止
# multipath -F # 删除之前的配置, 不建议
# lsblk
# /lib/udev/scsi_id --whitelisted --device=/dev/sda        #获取设备的wwid
# for i in `cat /proc/partitions | awk {'print $4'} |grep sd`; do
val=`/sbin/blockdev --getsize64 /dev/$i`; val2=`expr $val / 1073741824`; echo "/dev/$i:$val2 `/lib/udev/scsi_id -gud /dev/$i`"; done # 获取所有lun的wwid, VMware   Virtual disk没有wwid
# mpathconf --enable --find_multipaths y --with_module y --with_chkconfig y # 生成multipath.conf文件. 没删除过配置时不用重新生成
# vim /etc/multipath.conf
defaults {
    find_multipaths yes
    path_grouping_policy multibus
    user_friendly_names yes
}

blacklist {
    wwid xxx # 本地磁盘, 不需要多路径
}

multipaths {
    multipath { # 将不重复的wwid作为multipath内容追加到multipaths
        wwid xxx1
        alias yyy1
    }
    multipath {
        wwid xxx2
        alias yyy2
    }
    ...
}
# systemctl restart multipath-tools.service # 启动multipath, 启动后会在/dev/mapper下生成多路径逻辑盘(即/dev/dm*), 此时盘权限全部都是root用户, 有文章介绍要通过udev来修改权限, 不知能否通过chown修改.
# multipath -ll # 查看多路径状态
# ls -1 /dev/mapper/control  # 查看设备
```

## example
```bash
$ multipath -l # 查看多路径设备, 设备在`/dev/mapper`
$ multipath -F # 清除所有映射关系
$ multipath -f <multipath_device_name> # 清除指定映射关系
```

## FAQ
### find_multipaths=yes
如果将 find_multipaths 配置参数设定为 yes，那么multipath 将只在满足以下三个条件之一时创建设备：
- 至少有两个没有列入黑名单的路径使用同一WWID
- 用户使用 multipath 命令指定设备强制手动生成该设备
- 路径拥有与之前创建的多路径设备相同的 WWID（ 即使该多路径设备目前不存在）

旧版 multipath 总是尝试为每个没有明确放入黑名单的路径创建multipath 设备. 而find_multipaths这个功能可让大多数用户自动选择正确的路径创建多路径设备， 而无需编辑黑名单.