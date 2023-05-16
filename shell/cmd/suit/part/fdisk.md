# fdisk
管理硬盘设备.  推荐[**使用`parted`的gpt分区**](https://wiki.archlinux.org/index.php/Partitioning_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)).

磁盘分区的id:
- 83 : 主分区和逻辑分区
- 5 : 扩展分区
- 82 : swap
- 7 : ntfs分区

heads: 磁盘面数
sectors: 扇区数
cylinders: 柱面数
容量=heads*(sectors*512)*cylinders

boot: 是否为引导分区

## 选项
- m 查看全部可用的参数
- n 添加新的分区
- d 删除某个分区信息
- l 列出所有可用的分区类型
- t 改变某个分区的类型
- p 查看分区信息
- w 保存并退出
- q 不保存直接退出
- u : 与`-l`联用, 用扇区数取代柱面数, 用来表示每个分区的起始地址

## other
`partprobe`命令: 分区后手动将分区信息同步到内核，而且一般推荐连续两次执行该命令，效果会更好.

## FAQ
### fdisk修改ntfs磁盘签名报`WARNING: Re-reading the partition table failed with error 16: Device or resource busy...the next reboot or after you run partprobe(8) or kpartx(8)`
ref:
- [How do I use new fdisk table without reboot (kpartx)?](https://unix.stackexchange.com/questions/117949/how-do-i-use-new-fdisk-table-without-reboot-kpartx)

	文章推荐使用`partx -u`, 未测试

出错命令:
```
# ntfsfix /dev/zd0p1
# ntfsfix /dev/zd0p2
# lsblk # 看到/dev/zd0<p1,p2>没有被mount
# fdisk /dev/zd0 <<EOF # 出错
x
i
0xe24cfedd
w
EOF
```

尝试过fdisk前调用`kpartx -l -u -s -v /dev/zd0`刷新内核分区信息, 但发现其p1,p2重复注册了分区(lsblk时看到zd0下多出两个重复分区p1,p2), 导致destroy zd0报busy, 解决方法: 先`kpartx -d /dev/zd0`再destroy.

解决方法: `partprobe /dev/zd0`, 报错明细减少但之后还有小概率出现过一次. 可试试`partprobe /dev/zd0`+`partx -u /dev/zd0 (= partx -d /dev/zd0 + partx -a /dev/zd0, 用"-u"还是有更小概率复现)`/`partx -d /dev/zd0`(直接先删除zd0在kernel分区cache中的信息, 此时lsbk看不到zd0的分区, 再由之后的fdisk触发刷新内核分区信息, 之后lsblk能看到zd0的分区, 未再出现)+`time.sleep(5)`+fdisk操作.
