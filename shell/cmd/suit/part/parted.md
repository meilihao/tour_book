# parted
gpt分区工具. 与fdisk类似, 也分为两种模式: 命令模式(直接进行分区, 适合编程使用)和交互模式.

> 同类的有gdisk, 即fdisk的gpt版, 命令与fdisk类似.

>  重启或使用`partprobe -s `让kernel刷新分区表, 即将新的分区表变更同步至kernel; 或仅刷新指定设备`partprobe /dev/zd123`

> sfdisk是fdisk的非交互式变体.

## 选项
- m : 可解析格式

## 判断是否系统盘
- `parted /dev/sda print`输出的"Number"后是否存在`boot`/`swap`/`esp`
- 通过`/dev/disk/by-lable`

## example
```bash
# parted -l /dev/nvme0n1 | grep "Partition Table" # 查看磁盘分区方案(是否是gpt)
# parted /dev/sda print # 打印分区
# parted /dev/sda mklabel gpt # 设为gpt磁盘
# parted /dev/sdb mkpart primary 0 50% # 划取一半做个分区
```