# fuser

## 描述

显示正在用一个文件的进程有哪些.

## 选项
- -m : 显示所有使用指定文件系统或块设备的进程
- -v : 输出更多信息

## example
```
$ fuser -m /home/ubuntu/test/smb <=> lsof /home/ubuntu/test/smb
/home/ubuntu/test/smb:  3307c # 进程3307在使用
$ grep /home/ubuntu/test/smb /proc/*/mounts
```
