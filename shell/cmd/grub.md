# grub
参考:
- [grub2详解(翻译和整理官方手册)](https://www.lagou.com/lgeduarticle/9097.html)

grub2按以下方式为磁盘及分区命名：
1. 不论磁盘是 IDE（PATA）、SATA 或 SCSI，也不论它们的普遍称谓是 hd 或 sd，所有磁盘一律称为 hd
1. 设备分区名称从1开始
1. 一个数字代表磁盘；第二个数字（存在的话）代表分区, 用逗号分隔
1. grub2使用img文件，不再使用grub中的stage1、stage1.5和stage2


## 其他
要导入支持ext文件系统的模块时，只需`insmod ext2`即可，实际上也没有ext3和ext4对应的模块.

kernel的完整启动参数列表见[这里](http://redsymbol.net/linux-kernel-boot-parameters). 本文只列出几个常用的：
- init=   ：指定Linux启动的第一个进程init的替代程序。
- root=   ：指定根文件系统所在分区，在grub中，该选项必须给定。
- ro,rw   ：启动时，根分区以只读还是可读写方式挂载。不指定时默认为ro。
- initrd  ：指定init ramdisk的路径。在grub中因为使用了initrd或initrd16命令，所以不需要指定该启动参数。
- rhgb    ：以图形界面方式启动系统。
- quiet   ：以文本方式启动系统，且禁止输出大多数的log message。
- net.ifnames=0：用于CentOS 7，禁止网络设备使用一致性命名方式。
- biosdevname=0：用于CentOS 7，也是禁止网络设备采用一致性命名方式。
             ：只有net.ifnames和biosdevname同时设置为0时，才能完全禁止一致性命名，得到eth0-N的设备名。

