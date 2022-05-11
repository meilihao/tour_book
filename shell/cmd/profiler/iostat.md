# iostat
参考:
- [深入理解iostat](https://bean-li.github.io/dive-into-iostat/)

用于块设备.

iostat数据的来源是Linux操作系统的`/proc/diskstats`, 从第四个字段开始，介绍的是该设备的相关统计:
- 主设备号
- 从设备号
- 设备名
- (rd_ios) : 读操作的次数
- (rd_merges):合并读操作的次数。如果两个读操作读取相邻的数据块，那么可以被合并成1个。
- (rd_sectors): 读取的扇区数量
- (rd_ticks):读操作消耗的时间（以毫秒为单位）。每个读操作从__make_request()开始计时，到end_that_request_last()为止，包括了在队列中等待的时间。
- (wr_ios):写操作的次数
- (wr_merges):合并写操作的次数
- (wr_sectors): 写入的扇区数量
- (wr_ticks): 写操作消耗的时间（以毫秒为单位）
- (in_flight): 当前未完成的I/O数量。在I/O请求进入队列时该值加1，在I/O结束时该值减1。 注意：是I/O请求进入队列时，而不是提交给硬盘设备时
- (io_ticks)该设备用于处理I/O的自然时间(wall-clock time)
- (time_in_queue): 对字段#10(io_ticks)的加权值

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

> 如果%iowait的值过高，表示硬盘存在I/O瓶颈，%idle值高，表示CPU较空闲. 如果%idle值高但系统响应慢时，有可能是CPU等待分配内存，此时应加大内存容量. %idle值如果持续低于10，那么系统的CPU处理能力相对较低，表明系统中最需要解决的资源是CPU. %system持续较高可能是网络和驱动程序堆栈上存在瓶颈.

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
- %util: 一秒中有百分之多少的时间用于 I/O 操作,就是io使用率，即被io消耗的cpu百分比, 越小表示磁盘越空闲, 持续大于90%, 需重视, 说明产生的I/O请求太多，I/O系统已经满负荷，该磁盘可能存在瓶颈.

r/s和w/s就是磁盘的iops, 分别是读iops和写iops.
rkB/s和wkB/s是磁盘的吞吐

await和svctm是一对相对的数据, await是i/o的处理时间, 包括队列时间和操作时间, 一般系统i/o处理时间应小于5ms, 一旦超过20ms, server会感觉卡顿. syctm表示设备i/o操作的服务时间, 一般await大于svctm, 它们差值越小, 说明 I/O 几乎没有等待时间即队列时间越短, 性能越好. 如果 await 远大于 svctm，说明I/O 队列太长，io响应太慢，则需要进行必要优化. 如果avgqu-sz比较大，也表示有当量io在等待. 

tps和吞吐量:
- tps : 每秒发送到物理设备的I/O请求次数
- kB_read/s : 每秒从设备读取的块数(512B/s)
- kB_wrtn/s : 每秒向设备写入的块数(512B/s)
- kB_read : 读取的总块数
- kB_wrtn : 写入的总块数

> iostat 工具是 sysstat 包的一部分

## 选项

- -c : 显示cpu利用率
- -d : 查看设备的利用率
- -k : 以kb为单位进行输出
- -p [设备名] : 查看所有/具体设备和它的分区的 I/O 统计
- -x : 显示扩张统计信息
- -m : 以 MB 为单位而不是 KB 查看所有设备的统计. 默认以 KB 显示输出
- -n : 显示nfs使用情况
- -N : 查看 LVM 磁盘 I/O 统计报告
- -t : 显示报告的时间戳

## 例

    # iostat
    # iostat 5 2 # 打算以 5 秒捕获的间隔捕获两个报告, iostat [Interval] [Number Of Reports], 使用特定的间隔输出
    # iostat -x -d -m 1 zd224 # 1, 刷新间隔(s)
    # pidstat -d 1 # 展示I/O统计，每秒更新一次