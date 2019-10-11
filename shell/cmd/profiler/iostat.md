# iostat

用于块设备.

## 描述

iostat 被用来报告**CPU**的统计和**设备与分区的输出/输出**的统计

推荐监控间隔: 1m.

cpu属性值说明：
- %user：CPU处在用户模式(application)下的时间百分比
- %nice：CPU花费在re-nicing进程(更改进程的执行顺序和优先级)上的时间百分比
- %system：CPU处在系统模式(kernel)下的时间百分比
- %iowait：CPU等待I/O操作的时间的百分比
- %steal：管理程序(hypervisor)为另一个虚拟进程提供服务而等待虚拟 CPU 的百分比
- %idle：CPU空闲时间百分比

> 如果%iowait的值过高，表示硬盘存在I/O瓶颈，%idle值高，表示CPU较空闲. 如果%idle值高但系统响应慢时，有可能是CPU等待分配内存，此时应加大内存容量. %idle值如果持续低于10，那么系统的CPU处理能力相对较低，表明系统中最需要解决的资源是CPU.%system持续较高可能是网络和驱动程序堆栈上存在瓶颈.

disk属性值说明：
- rrqm/s: 设备请求队列中, 每秒进行 merge 的读操作数目. 即 rmerge/s
- wrqm/s: 设备请求队列中, 每秒进行 merge 的写操作数目. 即 wmerge/s
- r/s: 每秒完成的读请求的次数(合并后的). 即 rio/s
- w/s: 每秒完成的写请求的次数(合并后的). 即 wio/s
- rsec/s: 每秒读扇区数, 每个扇区512B. 即 rsect/s. ssd没有该项
- wsec/s: 每秒写扇区数. 即 wsect/s. ssd没有该项
- rkB/s: 每秒读K字节数. 是 rsect/s 的一半，因为每扇区大小为512字节. 
- wkB/s: 每秒写K字节数. 是 wsect/s 的一半. 
- avgrq-sz: 平均每次设备I/O操作的数据大小 (以扇区为单位). 
- avgqu-sz: 平均I/O队列(即io等待中)长度. 
- await: 平均每次发出到设备的I/O请求直至到被服务的耗时 (毫秒), 包括请求在队列中的耗时和svctm.
- r_await: 平均每次设备read操作的等待时间 (毫秒). 
- w_await: 平均每次设备write操作的等待时间 (毫秒). 
- svctm: 平均每次设备I/O操作的服务(即响应)时间 (毫秒). 
- %util: 一秒中有百分之多少的时间用于 I/O 操作,就是io使用率，即被io消耗的cpu百分比

> 如果 %util 接近 100%，说明产生的I/O请求太多，I/O系统已经满负荷，该磁盘可能存在瓶颈. 如果 svctm 比较接近 await，说明 I/O 几乎没有等待时间；如果 await 远大于 svctm，说明I/O 队列太长，io响应太慢，则需要进行必要优化. 如果avgqu-sz比较大，也表示有当量io在等待. 

tps和吞吐量:
- tps : 每秒发送到物理设备的I/O请求次数
- kB_read/s : 每秒从设备读取的块数(512B/s)
- kB_wrtn/s : 每秒向设备写入的块数(512B/s)
- kB_read : 读取的总块数
- kB_wrtn : 写入的总块数

> iostat 工具是 sysstat 包的一部分

## 选项

- -d : 查看所有设备的 I/O 统计
- -p [设备名] : 查看所有/具体设备和它的分区的 I/O 统计
- -x : 显示所有设备的详细的 I/O 统计信息
- -m : 以 MB 为单位而不是 KB 查看所有设备的统计. 默认以 KB 显示输出
- -N : 查看 LVM 磁盘 I/O 统计报告
- -t : 显示时间

## 例

    # iostat
    # iostat 5 2 # 打算以 5 秒捕获的间隔捕获两个报告, iostat [Interval] [Number Of Reports], 使用特定的间隔输出
