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
- %iowait：CPU **空闲（idle）**的时间里，有多少比例是因为在等磁盘I/O而空闲的

    iowait是CPU指标，只说明**CPU在等I/O, 不是磁盘在"忙"的时间占比**, 可能只是单个慢查询，磁盘整体很闲. 要看 iostat await + avgqu-sz
- %steal：管理程序(hypervisor)为另一个虚拟进程提供服务而等待虚拟 CPU 的百分比
- %idle：CPU空闲时间百分比

> 如果%iowait的值过高，表示硬盘存在I/O瓶颈，%idle值高，表示CPU较空闲. 如果%idle值高但系统响应慢时，有可能是CPU等待分配内存，此时应加大内存容量. %idle值如果持续低于10，那么系统的CPU处理能力相对较低，表明系统中最需要解决的资源是CPU. %system持续较高可能是网络和驱动程序堆栈上存在瓶颈.

disk属性值说明：
- rrqm/s: 设备请求队列中, 每秒进行 merge 的读操作数目. 即 rmerge/s
- wrqm/s: 设备请求队列中, 每秒进行 merge 的写操作数目. 即 wmerge/s
- `r/s`/`w/s`: 每秒完成的读/写请求的次数即IOPS(合并后的). 随机I/O看这个，顺序I/O看kB/s
- rsec/s: 每秒读扇区数, 每个扇区512B. 即 rsect/s. ssd没有该项
- wsec/s: 每秒写扇区数. 即 wsect/s. ssd没有该项
- `rkB/s`/`wkB/s`: 每秒读/写吞吐量. 是 `rsect/s`/`wsect/s` 的一半，因为每扇区大小为512字节. 顺序读写的瓶颈看这个
- avgrq-sz: 平均每次设备I/O操作的数据大小 (以扇区为单位). 
- avgqu-sz: 平均I/O队列(正在处理+等待的I/O数)深度. 这个比%util更能反映磁盘繁忙程度 
- await: 平均I/O响应时间（队列等待+实际处理）, 包括请求在队列中的耗时和svctm. 高不一定等于磁盘慢，可能是队列太长

    await = (所有I/O的总等待时间 + 所有I/O的总服务时间) / I/O总数
    
    await 高，可能是队列排太长（并发I/O多），不一定是磁盘处理慢, 看 avgqu-sz 来判断
- r_await: 平均每次设备read操作的等待时间 (毫秒). 
- w_await: 平均每次设备write操作的等待时间 (毫秒). 
- svctm: 平均I/O服务时间（实际处理时间）

    假设I/O是串行处理的, 在SSD/NVMe上，这个假设不成立，svctm 的值严重失真, 因此现代iostat已废弃

    想知道磁盘的真实服务时间，用 fio 做基准测试，或者在应用层用 blktrace + blkparse 追踪每个I/O的实际处理时间
- %util: 采样周期内磁盘"有I/O在处理"的时间占比, 就是io使用率，即被io消耗的cpu百分比, 越小表示磁盘越空闲, 持续大于90%, 需重视, 说明产生的I/O请求太多，I/O系统已经满负荷，该磁盘可能存在瓶颈. 

    %util = (磁盘有I/O在处理的时间) / (采样周期), %util=100% 只说明"磁盘每一微秒都在处理I/O"，但不代表"没有更多I/O可以被并行处理":
    - 机械盘（HDD）是串行处理I/O的，同一时刻只能处理一个I/O。所以 %util=100% 确实意味着磁盘饱和了，新I/O必须排队等
    - SSD/NVMe 是并行处理I/O的。NVMe 可以同时处理 数万到数十万个 I/O（看队列深度）

        SSD/NVMe上%util=100%是正常现象
        SATA SSD需看NCQ深度, 可能还有并行能力
        NVMe应完全忽略它

        avgqu-sz 持续接近设备的最大队列深度，且 await 在增长 → 才是真正的SSD/NVMe瓶颈

await和svctm是一对相对的数据, await是i/o的处理时间, 包括队列时间和操作时间, 一般系统i/o处理时间应小于5ms, 一旦超过20ms, server会感觉卡顿. syctm表示设备i/o操作的服务时间, 一般await大于svctm, 它们差值越小, 说明 I/O 几乎没有等待时间即队列时间越短, 性能越好. 如果 await 远大于 svctm，说明I/O 队列太长，io响应太慢，则需要进行必要优化. 如果avgqu-sz比较大，也表示有当量io在等待. 

tps和吞吐量:
- tps : 每秒发送到物理设备的I/O请求次数
- kB_read/s : 每秒从设备读取的块数(512B/s)
- kB_wrtn/s : 每秒向设备写入的块数(512B/s)
- kB_read : 读取的总块数
- kB_wrtn : 写入的总块数

> iostat 工具是 sysstat 包的一部分

实战结论：
- iowait 高 + CPU 还有富余 → 可能是少数的I/O密集型进程在拖，整体还好
- iowait 高 + CPU 也跑满 → 系统整体I/O压力大，需要查
- iowait 低 + 应用很卡 → 可能是网络I/O（iowait不统计网络等待）

磁盘I/O排查决策树
1. Step 1：确认是不是I/O问题（别被iowait骗了）
    iostat -x 1 看 await 和 avgqu-sz，不是看 %util 或 iowait
    await < 10ms（SSD）或 < 100ms（HDD）→ I/O不慢，问题在别处
1. Step 2：区分是吞吐量瓶颈还是IOPS瓶颈

    rkB/s 很高但 r/s 不高 → 顺序读写，吞吐量瓶颈，看 kB/s 是否接近磁盘理论带宽
    r/s 很高但 rkB/s 不高 → 随机I/O，IOPS瓶颈，看 r/s 是否接近磁盘最大IOPS
1. Step 3：找是哪些进程在搞I/O

    iotop 或 pidstat -d 1 看每个进程的 I/O 吞吐量
    iotop -o 只看有I/O活动的进程，更清晰
1. Step 4：查I/O调度器是否合适

    cat /sys/block/sda/queue/scheduler 看当前调度器。
    HDD 建议 mq-deadline；SSD/NVMe 建议 none（或 kyber 高负载场景）
1. Step 5：深入追踪（必要时）
    blktrace -d /dev/sda -o - | blkparse -i - 追踪每个I/O从发起到完成的完整路径，精确到微秒级，可以看到I/O在哪个环节耗时最多

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
    # iostat -x -d -m 1 zd224 zd224 # 1, 刷新间隔(s)
    # pidstat -d 1 # 展示I/O统计，每秒更新一次