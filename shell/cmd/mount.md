# mount
挂载文件系统, 即把硬件设备与目录进行关联. 

## 格式
`mount 文件系统 挂载目录`

## 选项
- -a 挂载所有在/etc/fstab 中定义的文件系统. 它会在执行后自动检查/etc/fstab文件中有无疏漏被挂载的设备文件，如果有，则进行自动挂载操作
- -L<lable> : 磁盘分区标识的别名.
- -n : 不将挂载信息记录在/etc/mtab中
- -r : 以只读方式加载设备
- -t 指定文件系统的类型

    - ext4
    - iso9660
- -w : 已可读写的方式挂载设备, 默认
- -v : 输出详细信息
- -vvv : 比`-v`更详细
- -o : 指定挂载文件系统时的选项

    - async 以非同步的方式执行文件系统的输入输出动作
    - atime 每次存取都更新inode的存取时间，默认设置，取消选项为noatime
    - auto 必须在/etc/fstab文件中指定此选项. 指定-a参数时，会加载设置为auto的设备，取消选取为noauto
    - defaults 使用默认的选项默认选项为rw、suid、dev、exec、anto nouser与async
    - dev 可读文件系统上的字符或块设备，取消选项为nodev
    - exec 可执行二进制文件，取消选项为noexec
    - noatime 每次存取时不更新inode的存取时间
    - noauto 设置此选项, 就无法使用-a来加载
    - nodev 不读文件系统上的字符或块设备
    - noexec 无法执行二进制文件
    - nosuid 关闭set-user-identifier(设置用户ID)与set-group-identifer(设置组ID)设置位
    - nouser 使用户无法执行加载操作，默认设置
    - iocharset=xxx : 指定mount分区时使用的字符集
    - codepage=xxx : 指定mount分区时使用的内码表
    - remount 重新加载设备通常用于改变设备的设置状态
    - ro 以只读模式加载
    - rw 以可读写模式加载
    - suid 启动set-user-identifier(设置用户ID)与set-group-identifer(设置组ID)设置位，取消选项为nosuid
    - sync 以同步方式执行文件系统的输入输出动作
    - user 可以让一般用户加载设备

## example
```bash
# mount -t vfat /dev/sda1 /mnt/usb #挂载usb
# mount -t iso9660 /dev/hda /mnt/cdrom #挂载光盘
# mount  /dev/cdrom /mnt/cdrom #挂载光盘
# mount -v /export/sdc_share
mount: /srv/dev-disk-by-path-pci-0000-00-10.0-scsi-0-0-2-0-part1/test/ bound on /export/sdc_share.
```

# umount
卸载文件系统

## 格式
`nmount 挂载目录`