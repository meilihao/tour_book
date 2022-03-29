# perf
参考:
- [Linux性能分析工具Perf简介](https://segmentfault.com/a/1190000021465563)
- [利用perf剖析Linux应用程序](https://blog.gmem.cc/perf)
- [Linux性能优化实战学习笔记：第四十九讲](https://www.cnblogs.com/luoahong/p/11577395.html)

## perf report
Children/Self: 如果在record时收集了调用链, 则Overhead可以在Children、Self两个列中显示. Children显示子代函数的样本计数、Self显示函数自己的样本计数.

## FAQ
### 软中断ksoftirqd/n 占用CPU 过高排查
ref:
- [性能分析（5）- 软中断导致 CPU 使用率过高的案例](https://cloud.tencent.com/developer/article/1678685)

```bash
# perf top # 全局查看
# perf top -p 45558 # 查看进程pid=45558的性能信息
# perf top -C 1 -e cpu-clock # 对cpu0进行性能分析的展示，它能显示出各模块具体函数占用的cpu比例
# perf record -a -g -p 9 -- sleep 30 # 9为高ksoftirqd的pid
# perf report
# perf record -C 1 -e cpu-clock -g sleep 10 # 记录10s中cpu0的性能数据
```

> 部分平台上perf会引起奇怪的crash。使用 -e cpu-clock 参数可以避免这个问题.

显示前4个cpu的软中断情况: `watch -d "/bin/cat /proc/softirqs | /usr/bin/awk 'NR == 1{printf \"%-15s %-15s %-15s %-15s %-15s\n\",\" \",\$1,\$2,\$3,\$4}; NR > 1{printf \"%-15s %-15s %-15s %-15s %-15s\n\",\$1,\$2,\$3,\$4,\$5}'"`

分析得到: TIMER（定时中断）、SCHED（内核调度）、RCU（RCU 锁）等这几个软中断都在不停变化

以网络中断导致ksoftirqd高cpu的`perf record -a -g -p 9 -- sleep 30`(9是ksoftirqd pid)展示的调用栈举例:
- net_rx_action 和 netif_receive_skb，表明这是接收网络包（rx 表示 receive）
- br_handle_frame ，表明网络包经过了网桥（br 表示 bridge）
- br_nf_pre_routing ，表明在网桥上执行了 netfilter 的 PREROUTING（nf 表示netfilter），而我们已经知道 PREROUTING 主要用来执行 DNAT，所以可以猜测这里有 DNAT 发生
- br_pass_frame_up，表明网桥处理后，再交给桥接的其他桥接网卡进一步处理。比如，在新的网卡上接收网络包、执行 netfilter 过滤规则等等

### 火焰图
- 轴表示采样数和采样比例。**一个函数占用的横轴越宽，就代表它的执行时间越长**. 同一层的多个函数，则是按照字母来排序.
- 纵轴表示调用栈，由下往上根据调用关系逐个展开. 换句话说，上下相邻的两个函数中，下面的函数，是上面函数的父函数。这样，**调用栈越深，纵轴就越高**.
- **要注意图中的颜色，并没有特殊含义，只是用来区分不同的函数**

> 火焰图是动态的矢量图格式，所以它还支持一些动态特性。比如，鼠标悬停到某个函数上时，就会自动显示这个函数的采样数和采样比例。而当你用鼠标点击函数时，火焰图就会把该层及其上的各层放大，方便你观察这些处于火焰图顶部的调用栈的细节.

根据性能分析的目标来划分，火焰图可以分为下面这几种:

- on-CPU 火焰图：表示 CPU 的繁忙情况，用在 CPU 使用率比较高的场景中。
- off-CPU 火焰图：表示 CPU 等待 I/O、锁等各种资源的阻塞情况。
- 内存火焰图：表示内存的分配和释放情况。
- 热/冷火焰图：表示将 on-CPU 和 off-CPU 结合在一起综合展示。
- 差分火焰图：表示两个火焰图的差分情况，红色表示增长，蓝色表示衰减。常用来比较不同场景和不同时期的火焰图，以便分析系统变化前后对性能的影响情况。

从 perf record 记录生成火焰图的工具常用 https://github.com/brendangregg/FlameGraph, 它要生成火焰图，其实主要需要三个步骤：
- 执行 perf script, 将 perf record 的记录转换成可读的采样记录
- 执行 stackcollapse-perf.pl 脚本，合并调用栈信息
- 执行 flamegraph.pl 脚本，生成火焰图

组合起来就是`perf script -i /root/perf.data | ./stackcollapse-perf.pl --all |  ./flamegraph.pl > ksoftirqd.svg`, 执行成功后，使用浏览器打开 ksoftirqd.svg ，就可以看到生成的火焰图

### 跳转亲和性
ref:
- [CPU利用率高的定位思路和方法](https://freesion.com/article/49901216629/)

```bash
echo 8 > /proc/irq/96/smp_affinity # //8=2^3, 表示绑定到cpu3
```