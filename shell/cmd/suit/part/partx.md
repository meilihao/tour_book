# partx
修改磁盘分区表后, 无需重启, 用partx命令告诉内核, 分区已改动, 内核可以读入新的分区表信息.

## examples
```bash
partx -d /dev/sdb # 因为内核中存在部分未调整磁盘的信息，故先将所有信息清零
partx -a /dev/sdb # 添加调整后的磁盘分区信息
partx -s /dev/sdb # 显示磁盘分区信息
```