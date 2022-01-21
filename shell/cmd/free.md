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

## example
```
$ watch -n 3 -d free
```

## `/proc/meminfo`
- MemFree: 空闲内存数

    MemUsed(已被用掉的内存)=MemTotal-MemFree
- MemAvailable

    MemAvailable≈MemFree+Buffers+Cached，它是内核使用特定的算法计算出来的，是一个估计值. 它与MemFree的关键区别点在于，MemFree是说的系统层面，MemAvailable是说的应用程序层面

    之前遇到过`kylin v10 + 飞腾主板 + cpu>=16`的机器`MemAvailable<MemFree`, 估计与保留内存有关, 可参考[Linux内存管理 (25)内存sysfs节点解读](cnblogs.com/arnoldlu/p/8568330.html)
- Buffer：缓冲区内存数
- Cache：缓存区内存数
- Shared：多个进程共享的内存空间

公式:
- OS Mem total = OS Mem used + OS Mem free
- APP buffers/cache used = OS Mem used - OS Mem buffers - OS Mem cached
- APP buffers/cache free = OS Mem free + OS Mem buffers + OS Mem cached
- APP buffers/cache total = APP buffers/cache used + APP buffers/cache free

## FAQ
### buffer/cache区别
buffers是指用来给**块设备**做的缓冲大小，它只记录文件系统的**metadata以及 追踪瞬时页面(tracking in-flight pages)**.

cached是用来给**文件内容**做缓冲