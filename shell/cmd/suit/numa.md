# numa
# numastat
numastat用来查看某个（些） 进程或者整个系统的内存消耗在各个NUMA节点的分布情况

## 选项

- -n : 以MB显示, 默认是内存页数

输出字段说明:
- numa_hit表示成功地从该节点分配到的内存页数。
- numa_miss表示成功地从该节点分配到的内存页数， 但其本意是希望从别的节点分配， 失败以后退而求其次从该节点分配。
- numa_foreign与numa_miss互为“影子”， 每个numa_miss都来自另一个节点的numa_foreign。
- interleave_hit， 有时候内存请求是没有NUMA节点偏好的， 此时会均匀分配自各个节点（interleave） ， 这个数值就是这种情况下从该节点分配出去的内存页面数。
- local_node表示分配给运行在同一节点的进程的内存页数。
- other_node与上面相反。 local_node值加上other_node值就是numa_hit值

# numad
numad是一个可以自动管理NUMA亲和性（affinity） 的工具（同时也是一个后台进程）. 它实时监控NUMA拓扑结构（topology） 和资源使用， 并动态调整. 同时它还可以在启动一个程序前， 提供NUMA优化建议.

与numad功能类似， Kernel的auto NUMA balancing（/proc/sys/kernel/numa_balancing）也是进行动态NUMA资源的调节。 numad启动后会覆盖Kernel的auto NUMA balancing功能.

numad自己内部维护一个inclusive list和一个exclusive list； -p<pid>、 -x<pid>就是分别往这两个list里面添加进程id； -r<pid>就是从这两个list里面移除。 在numad首次启动时候， 可以重复多个-p或者-x； 启动后， 每次调用numad只能跟一个-p、 -x或者-r参数. 默认没有这些指定的话， numad会对系统所有进程进行NUMA资源优化.

> 与KSM的关联， 在于/sys/kernel/mm/ksm/m erge_nodes最好设置为0， 禁止KSM跨NUMA节点地同页合并.

## 选项
- -S 0/1 : 0表示只对inclusive list里面的进程进行NUMA优化； 1表示对除exclusive list以外的所有进程进行优化. 通常， -S与-p、 -x搭配使用.
- -R<cpu_list>， reserve， 指定一些CPU是numad不能染指的， numad不会在自动优化NUMA资源的时候把进程放到这些CPU上去运行。
- -t<百分比>， 它指示逻辑（logical） CPU（比如Intel HyperThread打开时） 运算能力对于它的HW core的比例。 这个值关系到numad内部调配资源时的计算， 默认是20%。
- -u<百分比>， numad最多能消耗每个NUMA节点多少资源（CPU和内存） ， 默认是85%。 numad毕竟不能取代内核调度器， 并不能接管系统里所有route的调度， 所以， 留有余地是必须的。 但当你确定将一个node专属（dedicate） 给一个进程时， 也可以设置-u 100%， 甚至超过100%， 但要小心。
- -C 0/1， 是否将NUMA节点的inactive file cache作为free memory对待。 默认为1， 表示进程的inactive file cache不纳入NUMA优化的考量， 即如果一个进程还有一些inactive filecache留在另一个节点上， numad也不会把它搬过来。
- -K 0/1， 控制是否将interleaved memory合并到一个NUMA节点。 默认是合并的， 但要注意， 合并到一个节点并不一定有最好的performance， 而应该根据实际的work-load来决定。 比如， 如果一个系统里主要就是一个大型的数据库应用程序（大量内存访问且地址随机） ， -K 1禁止numad合并interleaved memory反而有更好的性能.
- -m<百分比>， 它是一个阈值， 表示当内存中在本地节点的数量达到它所属进程的内存总量的多少时， numad停止对该进程的NUMA优化。
- -i<最小间隔： 最大间隔>， 最小值可以省略。 它设置numad 2次扫描系统情况的时间间隔。 通常用`-i 0`来终止（退出） numad.
- -H<时间间隔>， 它设置（override） 透明大页（见7.2节） 的扫描间隔时间。 默认地， numad会将/sys/kernel/mm/tranparent_hugepage/khugepaged/scan_sleep_millisecs值从默认的10000毫秒缩短为1000毫秒， 因为更激进的透明大页合并更有利于numad将页面在
NUMA节点之间迁移。
- -w<NCPUS[:MB]>， 它就是numad一次性地运行一下（而不是作为系统后台daemon） 以供咨询： “我有一个应用程序将要运行， 它会需要NCPUS个CPU， M兆内存，
numad， 你告诉我该把它放到哪个NUMA节点上运行好啊？ ”numad此时会返回一个合适的NUMA node list

# numactl
numactl则是主动地在程序起来时候就指定好它的NUMA节点. 它还可以设置共享内存/大页文件系统的内存策略， 以及进程的CPU和内存的亲和性.

## 选项
- --hardware， 列出来目前系统中可用的NUMA节点， 以及它们之间的距离（distance）
- --membind， 确保command执行时候内存都是从指定的节点上分配； 如果该节点没有足够内存， 返回失败。
- --cpunodebind， 确保command只在指定node的CPU上面执行。
- --phycpubind， 确保command只在指定的CPU上执行。
- --localalloc， 指定内存只从本地节点上分配。
- --preferred， 指定一个偏好的节点， command执行时内存优先从这个节点分配， 不够的话才从别的节点分配

## example
```bash
# --- 用numactl来控制启动一个客户机， 让它只运行在节点1上， 然后通过numastat来确认
numactl --membind=1 --cpunodebind=1 -- qemu-system-x86_64 -enable-kvm ...
numastat qemu-system
# --- 让客户机均匀地占用两个节点的资源
numactl --interleave=0,1 -- qemu-system-x86_64 -enable-kvm ...
numastat qemu-system
```