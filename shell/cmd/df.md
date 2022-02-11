# df

## 描述

查看linux系统的磁盘空间占用情况

## 选项

- -a : 包含全部的文件系统
- -h : 以易读的方式显示结果
- -i : 显示文件系统分区的inode信息
- -k : 以kb显示结果, 默认
- -m : 以mb显示结果
- -T : 显示磁盘分区的文件系统类型

## FAQ
### inode耗尽
`df -h`显示有空间剩余, 但创建文件报"No space", 通过`df -i`显示inode耗尽.

通过`for i in /*; do echo $i; find $i | wc -l; done`查看文件最多的目录, 如果确定目录范围，把`/*`写的具体点, 最终可定位导致该问题目录.

> inodes的大小在磁盘格式化分区时确定，跟分区的大小相关，分区越大，inodes越大，反之亦然.

### df已100%但`du -s -h`和`df -i`都有空余
之前删除过一个7g多的文件怀疑是删除文件未释放导致, 执行`lsof|grep deleted|grep xsession-errors`查看到有记录.

解决方法:
1. kill占用该文件的进程
2. reboot

> df基于statfs的系统调用, 直接分析superblock获取分区使用情况, 因此其运行速度不受文件多少影响.

> du基于文件, 它会对每个文件调用fstat的系统调用来获取大小.