# numastat
为进程和os显示每个NUMA(Non-Uniform Memory Architecture)节点的内存统计信息.

> 信息来源于: `/sys/devices/system/node/node${n}/numastat`. numa_miss和other_node高预示着numa问题. 如果进程的内存没有分配在本地节点上, 可通过renice调整或使用cpu亲和力.
> numa node 相关信息会在`/sys/devices/system/node/node${n}`里.

## 选项

- -c ：基于数据内容动态收缩列的宽度，内存被四舍五入为MB，对更密集显示多个NUMA节点的系统有帮助
- -m ：生成每个节点的内存使用信息,格式与`/proc/meminfo`类似
- -n ：显示原始数据统计信息,与`numastat`类似, 以MB为单位, 格式也略有不同.
- -p <PID>/<pattern> ：为指定进程显示每个节点的内存分配信息
- -s [<node>]：以降序排序，列出最大内存消费者,没有指定node时是以总数列排序
- -v ：显示详细信息
- -s ：忽略值为0的项, 但四舍五入为0的列还是会显示

## 字段
- numa_hit : 预期的内存在这个节点上成功分配
- numa_miss : 在这个节点上分配内存，但进程使用了不同的节点. 每个numa_miss在另一个节点上有一个numa_foreign
- numa_foreign : 为该节点准备的内存实际被分配到了别的节点上，每个numa_foreign在另一个节点上有一个numa_miss
- interleave_hit : 预期的在这个节点上成功交错分配的内存
- local_node : 当一个进程在节点上运行的时候，在这个节点上分配的内存
- other_node : 当一个进程运行在其它节点上，在这个节点分配的内存
