# dmesg
显示开机信息

kernel会将开机信息保存在系统缓冲区中(ring buffer)以及`/var/log/dmesg`.

## 选项
- -c : 显示开机信息后,清空ring buffer.
- -n : 设置记录信息的level
- -s : 设置缓冲区的大小, 默认是8192