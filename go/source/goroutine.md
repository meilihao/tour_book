# goroutine
goroutine是由Go运行时管理的用户层轻量级线程. 相较于操作系统线程,goroutine的资源占用和使用代价都要小得多.

> 析goroutine调度器的源码可参考雨痕所著的《Go语言学习笔记》.

每个goroutine的初始栈大小仅为2KB, 见[$GOROOT/src/runtime/stack.go]()的_StackMin.

## GPM
关于G、P、M的定义,可以参见`$GOROOT/src/runtime/runtime2.go`:
- G: 代表goroutine,存储了goroutine的执行栈信息、goroutine状态及goroutine的任务函数等。另外G对象是可以重用的
- P: 代表逻辑processor,P的数量决定了系统内最大可并行的G的数量(前提:系统的物理CPU核数>=P的数量). P中最有用的是其拥有的各种G对象队列、链表、一些缓存和状态
- M: 代表着真正的执行计算资源. 在绑定有效的P后,进入一个调度循环;而调度循环的机制大致是从各种队列、P的本地运行队列中获取G,切换到G的执行栈上并执行G的函数,调用goexit做清理工作并回到M, 如此反复. M并不保留G状态,这是G可以跨M调度的基础

## G的抢占调度
M是如何让G停下来并调度下一个可运行的G的呢? 答案是:**G是被抢占调度的**. 除非极端的无限循环或死循环,否则只要G调用函数,Go运行时就有了抢占G的机会. 因为是Go编译器在函数的入口处插入了一个运行时的函
数调用:runtime.morestack_noctxt. 它会检查是否需要扩容连续栈, 并进入抢占调度的逻辑中. 一旦所在goroutine被置为可被抢占的,那么抢占调度代码就会剥夺该goroutine的执行权,将其让给其他goroutine.

调用函数可能不会触发抢占的原因(by `go tool objdump -S go-binary > xx.s/go tool compile -S xx.go >xx.s/go build -gcflags '-S' xx.go > xx.s`):
1. 内联(inline)优化的结果
1. 为被调用函数位于调用树的leaf(叶子)位置,编译器可以确保其不再有新栈帧生成,不会导致栈分裂或超出现有栈边界,于是就不用插入morestack

    `go build -gcflags '-S'`可看到被调用函数后出现`NOSPLIT`的标志, 说明它使用的栈是固定大小(24字节),不会再分裂(split)或超出现有栈边界

    因此解决的最简单的方法就是再封装一层

在Go程序启动时,运行时会启动一个名为sysmon的M(一般称为监控线程),该M的特殊之处在于它无须绑定P即可运行(以g0这个G的形式). 该M在整个Go程序的运行过程中至关重要.

sysmon每20us~10ms启动一次,主要完成如下工作:
1. 释放闲置超过5分钟的span物理内存
1. 如果超过2分钟没有垃圾回收,强制执行
1. 将长时间未处理的netpoll结果添加到任务队列
1. 向长时间运行的G任务发出抢占调度
1. 收回因syscall长时间阻塞的P

sysmon将向长时间运行的G任务发出抢占调度,这由函数retake(`$GOROOT/src/runtime/proc.go`)实现: 如果一个G任务运行超过10ms,sysmon就会认为其运行时间太久而发出抢占式调度的请求。一旦G的抢占标志位被设为true,那么在这个G下一次调用函数或方法时,运行时便可以将G抢占并移出运行状态,放入P的本地运行队列中(如果P的本地运行队列已满, 那么将放在全局运行队列中),等待下一次被调度.

channel阻塞或网络I/O情况下的调度:
如果G被阻塞在某个channel操作或网络I/O操作上,那么G会被放置到某个等待队列中,而M会尝试运行P的下一个可运行的G。如果此时P
没有可运行的G供M运行,那么M将解绑P,并进入挂起状态。当I/O操作完成或channel操作完成,在等待队列中的G会被唤醒,标记为
runnable(可运行),并被放入某个P的队列中,绑定一个M后继续执行。

系统调用阻塞情况下的调度:
如果G被阻塞在某个系统调用上,那么不仅G会阻塞,执行该G的M也会解绑P(实质是被sysmon抢走了),与G一起进入阻塞状态。如果
此时有空闲的M,则P会与其绑定并继续执行其他G;如果没有空闲的M,但仍然有其他G要执行,那么就会创建一个新M(线程)。当系统
调用返回后,阻塞在该系统调用上的G会尝试获取一个可用的P,如果有可用P,之前运行该G的M将绑定P继续运行G;如果没有可用的P,那么
G与M之间的关联将解除,同时G会被标记为runnable,放入全局的运行队列中,等待调度器的再次调度。

## 调试信息输出标志字符串
ref:
    - [关于Go调度器调试信息输出的详细信息,可以参考Dmitry Vyukov的 "Debugging Performance Issues in Go Programs"](https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs)

Go提供了调度器当前状态的查看方法:使用Go运行时环境变量GODEBUG, 比如`GODEBUG=schedtrace=1000`就是每1000ms打印输出一次goroutine调度器的状态.

`SCHED 6016ms: gomaxprocs=4 idleprocs=0 threads=26 spinningthreads=0 idlethreads=20 runqueue=1 [3 4 0 10]`解读
- SCHED:调试信息输出标志字符串,代表本行是goroutine调度器相关信息的输出
- 6016ms:从程序启动到输出这行日志经过的时间
- gomaxprocs:P的数量
- idleprocs:处于空闲状态的P的数量。通过gomaxprocs和idleprocs的差值, 就可以知道当前正在执行Go代码的P的数量
- threads:操作系统线程的数量,包含调度器使用的M数量,加上运行时自用的类似sysmon这样的线程的数量
- spinningthreads:处于自旋(spin)状态的操作系统线程数量
- idlethread:处于空闲状态的操作系统线程的数量
- runqueue=1:Go调度器全局运行队列中G的数量
- [3 4 0 10]:分别为4个P的本地运行队列中的G的数量

输出每个goroutine、M和P的详细调度信息可用`GODEBUG=schedtrace=1000,scheddetail=1`.