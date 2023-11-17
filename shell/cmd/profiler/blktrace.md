# blktrace
参考:
- [IO神器blktrace使用介绍](https://developer.aliyun.com/article/698568)
- [利用blktrace分析IO性能](http://linuxperf.com/?p=161)
- [Linux命令拾遗-使用blktrace分析io情况](https://www.cnblogs.com/codelogs/p/16060775.html)

blktrace包安装后有blktrace、blkparse、btt、blkiomon这4个命令，其中blktrace负责采集I/O事件数据，blkparse负责将每一个I/O事件数据解析为纯文本方便阅读，btt、blkiomon负责统计分析

## example
```bash
# blktrace -d /dev/sda -o - |blkparse -i -
```

# blkparse
## 输出
```bash
# blktrace -d /dev/sda -o - |blkparse -i -
  8,0    3        1     0.000000000   594  A FWFSM 43781957 + 17 <- (8,4) 18782021
  8,0    3        2     0.000001967   594  Q FWFSM 43781957 + 17 [(null)]
  8,0    3        3     0.000006490   594  G FWFSM 43781957 + 17 [(null)]
  8,0    3        4     0.000031145   195  D  FN [kworker/3:1H]
  8,0    2        7     0.001628764     0  C  FN 0 [0]
  8,0    2        8     0.001671798   323  D WSM 43781957 + 17 [kworker/2:1H]
  8,0    2        9     0.001843458     0  C WSM 43781957 + 17 [0]
  8,0    2       10     0.001863050   323  D  FN [kworker/2:1H]
  8,0    2       11     0.002735974     0  C  FN 0 [0]
  8,0    2       12     0.002742812     0  C WSM 43781957 [0]
  8,0    3        5     1.477349073   233  A  WM 26379008 + 32 <- (8,4) 1379072
  8,0    3        6     1.477351508   233  Q  WM 26379008 + 32 [xfsaild/sda4]
  8,0    3        7     1.477362690   233  G  WM 26379008 + 32 [xfsaild/sda4]
  8,0    3        8     1.477364537   233  P   N [xfsaild/sda4]
  8,0    3        9     1.477375143   233  A  WM 26713248 + 32 <- (8,4) 1713312
  8,0    3       10     1.477376067   233  Q  WM 26713248 + 32 [xfsaild/sda4]
  8,0    3       11     1.477379594   233  G  WM 26713248 + 32 [xfsaild/sda4]
  8,0    3       12     1.477389324   233  A  WM 26713792 + 8 <- (8,4) 1713856
  8,0    3       13     1.477390200   233  Q  WM 26713792 + 8 [xfsaild/sda4]
  8,0    3       14     1.477393235   233  G  WM 26713792 + 8 [xfsaild/sda4]
...
```

字段:
1. 第一个字段：`8,0`是设备号 major device ID和minor device ID
1. 第二个字段：`3`表示CPU
1. 第三个字段：`11`为序列号
1. 第四个字段：`0.009507758`是Time Stamp, 即时间偏移
1. 第五个字段：PID 本次IO对应的进程PID
1. 第六个字段：Event，**反映了IO进行到了哪一步**
1. 第七个字段：R表示 Read， W是Write，D表示block，B表示Barrier Operation
1. 第八个字段：223490+56，表示的是起始block number 和 number of blocks，即常说的Offset 和 Size
1. 第九个字段：进程名

其中第六个字段非常有用：每一个字母都代表了IO请求所经历的某个阶段:
- A ： remap 对于栈式设备，进来的 I/O 将被重新映射到 I/O 栈中的具体设备
- X ： split 对于做了 Raid 或进行了 device mapper(dm) 的设备，进来的 I/O 可能需要切割，然后发送给不同的设备
- Q ： queued I/O 进入 block layer，将要被 request 代码处理（即将生成 I/O 请求）
- G ： get request I/O 请求（request）生成，为 I/O 分配一个 request 结构体
- M ： back merge 之前已经存在的 I/O request 的终止 block 号，和该 I/O 的起始 block 号一致，就会合并，也就是向后合并
- F ： front merge 之前已经存在的 I/O request 的起始 block 号，和该 I/O 的终止 block 号一致，就会合并，也就是向前合并
- I ： inserted I/O 请求被插入到 I/O scheduler 队列
- S ： sleep 没有可用的 request 结构体，也就是 I/O 满了，只能等待有 request 结构体完成释放
- P ： plug 当一个 I/O 入队一个空队列时，Linux 会锁住这个队列，不处理该 I/O，这样做是为了等待一会，看有没有新的 I/O 进来，可以合并
- U ： unplug 当队列中已经有 I/O request 时，会放开这个队列，准备向磁盘驱动发送该 I/O。这个动作的触发条件是：超时（plug 的时候，会设置超时时间）；或者是有一些 I/O 在队列中（多于 1 个 I/O）
- D ： issued I/O 将会被传送给磁盘驱动程序处理
- C ： complete I/O 处理被磁盘处理完成

这些Event中常见的出现顺序如下:
Q – 即将生成IO请求
|
G – IO请求生成
|
I – IO请求进入IO Scheduler队列
|
D – IO请求进入driver
|
C – IO请求执行完毕

# btt
blktrace 拿到的数据是 per cpu 的, 需要合并数据再用btt 来做整个 IO 的延迟分析. btt 提供对 I/O 在 I/O 栈中不同区域所花费时间的分析.

由于每个Event都有出现的时间戳，根据这个时间戳就可以计算出 I/O 请求在每个阶段所消耗的时间，比如从Q事件到C事件的时间叫Q2C，那么常见阶段如下：
```
Q------------>G----------------->I----------------------------------->D----------------------------------->C
|-Q time(Q2G)-|-Insert time(G2I)-|------schduler wait time(I2D)-------|--------driver,disk time(D2C)-------|
|------------------------------- await time in iostat output(Q2C) -----------------------------------------|
```
阶段说明:
1. Q2G: I/O进入block layer到生成 I/O 请求所消耗的时间，包括 remap 和 split 的时间.

   Q2G 很高则意味着队列中同时有大量请求。这可能代表该存储无法承担 I/O 负载
1. G2I: I/O 请求进入 I/O Scheduler(request queue)所消耗的时间，包括 merge 的时间
1. Q2M: I/O进入block层到该I/O被和已存在的I/O请求合并的时间
1. I2D: I/O 请求在 I/O Scheduler 中等待的时间，可以作为 I/O Scheduler 性能的指标
1. M2D: I/O合并成I/O请求到分发到设备驱动的时间
1. D2C: I/O 请求在 Driver 和硬件上所消耗的时间，可以作为硬件性能的指标.

   D2C 很高，那么该设备服务请求的时间就很长。这可能表示该设备只是超载了（可能是由共享资源造成的），或者是因为发送给给设备的负载未经优化

   D2C 的 AVG 越大，代表设备越慢，延迟越高; D2C 的 MAX 越大，代表设备的 tail latency 越高.
1. Q2C: 整个 I/O 请求所消耗的时间(Q2I + G2I + I2D + D2C = Q2C)，相当于 iostat 的 await, 代表了内核软件栈引入的延迟
1. Q2Q: 相邻2个 I/O 请求进入通用块层的I/O间隔

   如果 Q2Q 比 Q2C 大很多，则意味着程序没有以快速连续法式发出请求. 因此性能问题可能与 I/O 子系统无关

在上述过程中，Q2M、M2D两个阶段不是必然发生的，只有可以merge的I/O才会进行合并