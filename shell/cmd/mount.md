# mount
挂载文件系统, 即把硬件设备与目录进行关联. 

## 格式
`mount 文件系统 挂载目录`

## 选项
- -a 挂载所有在/etc/fstab 中定义的文件系统. 它会在执行后自动检查/etc/fstab文件中有无疏漏被挂载的设备文件，如果有，则进行自动挂载操作
- -t 指定文件系统的类型

## example
```bash
# mount -t vfat /dev/sda1 /mnt/usb #挂载usb
# mount -t iso9660 /dev/hda /mnt/cdrom #挂载光盘
#  mount  /dev/cdrom /mnt/cdrom #挂载光盘
```