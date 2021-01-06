# parted
gpt分区工具. 与fdisk类似, 也分为两种模式: 命令模式(直接进行分区, 适合编程使用)和交互模式.

> 同类的有gdisk, 即fdisk的gpt版, 命令与fdisk类似.

>  重启或使用`partprobe -s `让kernel刷新分区表, 即将新的分区表变更同步至kernel.

> sfdisk是fdisk的非交互式变体.