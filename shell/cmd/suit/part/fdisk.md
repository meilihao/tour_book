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
