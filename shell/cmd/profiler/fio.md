# fio
参考:
- [磁盘测试工具FIO](https://www.cnblogs.com/klb561/p/11939355.html)
- [存储设备性能benchmark](https://pibenchmarks.com/)

    - [TheRemote/PiBenchmarks : Linux Storage Benchmarking Script](https://github.com/TheRemote/PiBenchmarks)

模拟各种情况的i/o基准测试工具, 支持 14 种不同的 I/O 引擎，包括: sync,mmap, libaio, posixaio, SG v3, splice, null, network, syslet, guasi, solarisaio, iouring 等等.

fio 分顺序读，随机读，顺序写，随机写，混合随机读写模式.

一般使用以下几个指标对存储设备的性能进行描述：
- IOPS：每秒读/写次数，单位为次（计数）. 存储设备的底层驱动类型决定了不同的 IOPS.
- 吞吐量：每秒的读写数据量，单位为MB/s
- 时延：IO操作的发送时间到接收确认所经过的时间，单位为秒

注意：
1. 请不要在系统盘上进行 fio 测试，避免损坏系统重要文件
2. fio测试建议在空闲的、未保存重要数据的硬盘上进行，并在测试完后重新制作文件系统。请不要在业务数据硬盘上测试，避免底层文件系统元数据损坏导致数据损坏
2. 测试硬盘性能时，推荐直接测试裸盘（如 vdb）；测试文件系统性能时，推荐指定具体文件测试（如 /data/file）

不同场景的测试公式基本一致，只有3个参数（读写模式，iodepth，blocksize）的区别. 常见用例如下：
- block=4k iodepth=1 随机读测试，能反映磁盘的时延性能
- block=128K iodepth=32 能反映峰值吞吐性能
- block=4k iodepth=32 能反映峰值IOPS性能

## 选项
- filename=/dev/sdb1   # 测试设备的设备文件
- direct=1             # 测试过程绕过机器自带的 buffer 使测试结果更真实
- rw=randwrite/randrw/randread/read(顺序)/write(顺序)  # 读写模式
- bs=16k               # 单次 io 的块文件大小为 16k
- bsrange=512-2048     # 同上，提定数据块的大小范围
- size=5G              # 本次的测试文件大小为 5g，以每次 4k 的 io 进行测试
- numjobs=30           # 本次的测试线程为 30 个
- runtime=1000         # 测试时间 1000 秒，如果不写则一直将 5g 文件分 4k 每次写完为止
- iodepth=1            # 请求IO队列的深度
- ioengine=psync       # io 引擎使用 psync 方式

    - libaio : Linux 原生的异步 I/O，这也是通常我们这边用的最多的测试盘吞吐和延迟的方法, 需要yum install libaio-devel/apt install libaio-dev
    - sync : 也就是最通常的 read / write 操作
    - vsync : 使用 readv / writev，主要是会将相邻的 I/O 进行合并
    - psync : 对应的 pread / pwrite
    - pvsync / pvsync2 : 对应的 preadv / pwritev，以及 preadv2 / p writev2
- rwmixwrite=30        # 在混合读写的模式下，写占 30%
- group_reporting      # 关于显示结果的，汇总每个进程的信息
- lockmem=1G           # 只使用 1g 内存进行测试
- zero_buffers         # 用 0 初始化系统 buffer
- name=xxx             # job的名称
- nrfiles=8            # 每个进程生成文件的数量
- norandommap          # 在随机io时有用, 默认随机io也会写入所有的size里描述的快, 加了之后就打破了这个限制, 有些快可能无法被read/write到, 有些则可能io多次, 能够更好地模拟用户场景
- refill_buffers       # 每次提交IO任务会重新填充
- randrepeat=0         # 随机序列是否重复
- size=100G            # io测试的寻址空间
- rwmixwrite=30        # 在混合读写的模式下，写占30%
- -fdatasync=1         # 落盘方式
- `--debug=all`        : debug log
- -thread : fio 默认使用 进程（Process） 模式来运行任务，而不添加 --thread 会导致内存占用大幅增加，甚至在嵌入式设备（如 MMC 卡应用场景）上导致 OOM（内存溢出）

## fio输出解读
slat / clat / lat, 这几个是 latency 指标:
- slat 就是 Submission latency，也就是提交到实际执行 I/O 的时间，在 sync 测试里面这个是没有的，因为 slat 就是 clat
- clat 就是 Completion latency，也就是从提交到完成的时间
- lat 就是 Total latency，包括 fio 从创建这个 I/O 单元到完成的总的时间

disk stats:
- ios : 总的 I/O 操作次数
- merge : 被 I/O 调度合并的次数
- ticks : 让磁盘保持忙碌的次数
- in_queue : 总的在磁盘队列里面的耗时
- util : 磁盘的利用率

fio 还支持将中间的操作输出到文件，然后使用工具绘制图表展示，通常就是设置 write_bw_log，write_bw_log 和 write_iops_log，然后使用 fio_generate_plots 来绘图，另外也可以用 fio2gnuplot 来绘制.

## example
```bash
# 顺序读
fio -filename=/dev/sda -direct=1 -iodepth 1 -thread -rw=read -ioengine=psync -bs=16k -size=200G -numjobs=30 -runtime=1000 -group_reporting -name=mytest

# 顺序写
fio -filename=/dev/sda -direct=1 -iodepth 1 -thread -rw=write -ioengine=psync -bs=16k -size=200G -numjobs=30 -runtime=1000 -group_reporting -name=mytest

# 随机读
fio -filename=/dev/sda -direct=1 -iodepth 1 -thread -rw=randread -ioengine=psync -bs=16k -size=200G -numjobs=30 -runtime=1000 -group_reporting -name=mytest

# 随机写
fio -filename=/dev/sda -direct=1 -iodepth 1 -thread -rw=randwrite -ioengine=psync -bs=16k -size=200G -numjobs=30 -runtime=1000 -group_reporting -name=mytest

# 混合随机读写
fio -filename=/dev/sda -direct=1 -iodepth 1 -thread -rw=randrw -rwmixread=70 -ioengine=psync -bs=16k -size=200G -numjobs=30 -runtime=100 -group_reporting -name=mytest -ioscheduler=noop

# 利用配置文件
# cat << EOF > fio.conf
[global]
ioengine=libaio
direct=1
thread=1
norandommap=1
randrepeat=0
runtime=60
ramp_time=6
size=1g
directory=/path/to/test

[read4k-rand]
stonewall
group_reporting
bs=4k
rw=randread
numjobs=8
iodepth=32

[read64k-seq]
stonewall
group_reporting
bs=64k
rw=read
numjobs=4
iodepth=8

[write4k-rand]
stonewall
group_reporting
bs=4k
rw=randwrite
numjobs=2
iodepth=4

[write64k-seq]
stonewall
group_reporting
bs=64k
rw=write
numjobs=2
iodepth=4
EOF

# 测试
fio fio.conf
```

[fs测试](https://juicefs.com/docs/zh/cloud/single_node_benchmark):
```bash
$ fio --name=big-file-sequential-read \
    --directory=/jfs \
    --rw=read --refill_buffers \
    --bs=256k --size=4G
$ fio --name=big-file-sequential-write \
    --directory=/jfs \
    --rw=write --refill_buffers \
    --bs=256k --size=4G
```

## FAQ
### io scheduler noop not found
`cat /sys/block/sda/queue/scheduler`, 查看当前系统是否没选中noop. 修改即可:`echo 'noop' > /sys/block/sda/queue/scheduler`.

### fio测试时报`fio: job startup hung? exiting.`
可能的问题:
1. 内存不足