# numa

# numactl

## 参数
- --hardware : 输出硬件信息

# numastat
> `yum install numactl`

为进程和os显示每个NUMA(Non-Uniform Memory Architecture)节点的内存统计信息.

> 信息来源于: `/sys/devices/system/node/node${n}/numastat`. numa_miss和other_node高预示着numa问题. 如果进程的内存没有分配在本地节点上, 可通过renice调整或使用cpu亲和力.
> numa node 相关信息会在`/sys/devices/system/node/node${n}`里.

linux默认开启numa自动平衡策略. 开关用`echo 0(关闭)/1(开启) > /proc/sys/kernel/numa_balancing`.

内存消耗型应用可关闭numa自动平衡策略, 避免出现一个cpu耗光, 另一个cpu还有大量可用, 而系统却去使用swap, 导致性能严重降低.

## 选项

- -c xxx ：查看相关进程的numa内存使用情况. 基于数据内容动态收缩列的宽度，内存被四舍五入为MB，对更密集显示多个NUMA节点的系统有帮助
- -m ：生成每个节点的内存使用信息,格式与`/proc/meminfo`类似
- -n ：显示原始数据统计信息,与`numastat`类似, 以MB为单位, 格式也略有不同.
- -p <PID>/<pattern> ：为指定进程显示每个节点的内存分配信息
- -s [<node>]：以降序排序，列出最大内存消费者,没有指定node时是以总数列排序
- -v ：显示详细信息
- -s ：忽略值为0的项, 但四舍五入为0的列还是会显示

## 字段
- numa_hit : 使用本节点的内存次数
- numa_miss : 预期使用被节点内存, 但进程被调度到其他不同节点的次数. 每个numa_miss在另一个节点上有一个numa_foreign
- numa_foreign : 预期使用其他节点内存, 但使用本地内存次数 .每个numa_foreign在另一个节点上有一个numa_miss
- interleave_hit : 交叉分配使用的内存中使用本节点内存的次数
- local_node : 在当前节点运行的程序使用本节点的内存次数
- other_node : 在其他节点运行的程序使用本节点的内存次数