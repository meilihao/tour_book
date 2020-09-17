# multipath
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