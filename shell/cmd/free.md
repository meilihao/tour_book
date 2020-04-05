# free

查看内存使用状况, 信息来自`/proc/meminfo`.

第一行(Mem):
- total : 物理内存总量
- used : 已使用的物理内存
- free : 空闲的物理内存
- shared : 多个进程共享的内存
- buff/cache : 块设备的读写缓冲+文件系统的cache所使用的内存
- available : 系统实际可用的内存

第二行(Swap):
- total : swap总量
- used : 已使用的swap
- free : 未使用的swap

## 选项
- -b : 以字节为单位显示
- -m : 以MB为单位显示
- -K : 以KB为单位显示
- -t : 显示内存总和列
- -s <秒数> : 以指定的间隔持续显示内存使用情况
- -o : 不显示系统缓冲区列