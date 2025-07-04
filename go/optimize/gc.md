# gc
参考:
- [GC 的认识](https://github.com/qcrao/Go-Questions/blob/master/GC/GC.md)
- [Go 语言设计与实现 | 7.2 垃圾收集器](https://draven.co/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/)

GC(Garbage Collection，即垃圾回收)，是一种自动管理内存的机制.

当程序向操作系统申请的内存不再需要时，垃圾回收器主动将其回收并供其他代码进行内存申请时复用或者将其归还给操作系统，这种针对内存资源的自动回收过程，即为垃圾回收, 而负责垃圾回收的程序组件，即为垃圾回收器.

gc优缺点:
- 优点: 降低开发者心智负担, 保障代码不易出现内存泄露
- 缺点: 需要占用资源进行垃圾回收, 影响系统性能及吞吐

## 常见的 GC 实现方式有哪些？Go 语言的 GC 使用的是什么？
所有的 GC 算法其存在形式可以归结为追踪（Tracing）和引用计数（Reference Counting）这两种形式的混合运用.
- 追踪式 GC

    从根对象出发，根据对象之间的引用信息，一步步推进直到扫描完毕整个堆并确定需要保留的对象，从而回收所有可回收的对象. Go、 Java、V8 对 JavaScript 的实现等均为追踪式 GC.

- 引用计数式 GC

    每个对象自身包含一个被引用的计数器，当计数器归零时自动得到回收. 因为此方法缺陷较多，在追求高性能时通常不被应用. Python、Objective-C 等均为引用计数式 GC.

目前比较常见的 GC 实现方式包括：
- 追踪式，分为多种不同类型，例如：
        
    - 标记清扫：从根对象出发，将确定存活的对象进行标记，并清扫可以回收的对象
    - 标记整理/标记压缩：为了解决内存碎片问题而提出，在标记过程中，将对象尽可能整理到一块连续的内存上
    - 增量式：将标记与清扫的过程分批执行，每次执行很小的部分，从而增量的推进垃圾回收，达到近似实时、几乎无停顿的目的
    - 增量整理：在增量式的基础上，增加对对象的整理过程
    - 分代式：将对象根据存活时间的长短进行分类，存活时间小于某个值的为年轻代，存活时间大于某个值的为老年代，永远不会参与回收的对象为永久代。并根据分代假设（如果一个对象存活时间不长则倾向于被回收，如果一个对象已经存活很长时间则倾向于存活更长时间）对对象进行回收
- 引用计数：根据对象自身的引用计数来回收，当引用计数归零时立即回收. 其最大的缺点是无法解决循环引用问题.

关于各类方法的详细介绍及其实现不在本文中详细讨论。对于 Go 而言，Go 的 GC 目前使用的是**无分代（对象没有代际之分）、不整理（回收过程中不对对象进行移动与整理）、并发（与用户代码并发执行）的三色标记清扫算法**. 原因在于：
1. 对象整理的优势是解决内存碎片问题以及“允许”使用顺序内存分配器. 但 Go 运行时的分配算法基于 tcmalloc，基本上没有碎片问题, 并且顺序内存分配器在多线程的场景下并不适用. Go 使用的是基于 tcmalloc 的现代内存分配算法，对对象进行整理不会带来实质性的性能提升.
1. 分代 GC 依赖分代假设，即 GC 将主要的回收目标放在新创建的对象上（存活时间短，更倾向于被回收），而非频繁检查所有对象。但 **Go 的编译器会通过逃逸分析将大部分新生对象存储在栈上（栈直接被回收），只有那些需要长期存在的对象才会被分配到需要进行垃圾回收的堆中**。也就是说，分代 GC 回收的那些存活时间短的对象在 Go 中是直接被分配到栈上，当 goroutine 死亡后栈也会被直接回收，不需要 GC 的参与，进而分代假设并没有带来直接优势。并且 Go 的垃圾回收器与用户代码并发执行，使得 STW 的时间与对象的代际、对象的 size 没有关系。Go 团队更关注于如何更好地让 GC 与用户代码并发执行（使用适当的 CPU 来执行垃圾回收），而非减少停顿时间这一单一目标上。


## [gc log](https://godoc.org/runtime)
参考:
- [GODEBUG之gctrace解析](http://cbsheng.github.io/posts/godebug%E4%B9%8Bgctrace%E8%A7%A3%E6%9E%90/)

调试GC打开`GODEBUG=gctrace=1`, 相关log:
```
scvg: 0 MB released
scvg: inuse: 3, idle: 59, sys: 63, released: 58, consumed: 4 (MB)
gc 252 @4316.062s 0%: 0.013+2.9+0.050 ms clock, 0.10+0.23/5.4/12+0.40 ms cpu, 16->17->8 MB, 17 MB goal, 8 P
```
scvg: 将gctrace设置为大于0的任何值时使得gc在内存释放回系统时输出的摘要信息. 将内存释放回系统的过程称为清理(scavenging).
scvg解读:
- inuse: 在使用或在部分使用的span大小
- idle: 等待清理的span大小
- sys: 系统映射的内存
- released: 释放回系统的内存
- consumed: 从系统申请的内存

gc解读:
- gc 252 ： 这是第252次gc
- @4316.062s ： 这次gc的markTermination阶段完成后距离runtime启动到现在的时间, 即程序启动以来的秒数
- 0% ：到目前为止，gc运行所用的CPU时间占总CPU的百分比
- 0.013+2.9+0.050 ms clock ：按顺序分成三部分，0.013表示mark阶段(设置GC清除停止及GC标记开始(需停止程序))的STW时间（单P的）; 2.9表示并发标记用的时间（所有P的）; 0.050表示markTermination阶段(设置GC标记结束(需停止程序))的STW时间（单P的）
- 0.10+0.23⁄5.4⁄12+0.40 ms cpu ：gc占用cpu时间,按顺序分成三部分，0.10表示整个进程在mark阶段STW停顿时间(`0.013*8`)；0.23⁄5.4/12有三块信息，0.23是mutator assists占用的时间，5.4是dedicated mark workers+fractional mark worker占用的时间，12是idle mark workers占用的时间. 这三块时间加起来会接近`2.9*8`；0.40 ms表示整个进程在markTermination阶段STW停顿时间(`0.050 * 8`), 8是P的个数.
- 16->17->8 MB ：按顺序分成三部分，16表示开始mark阶段前的heap_live大小；17表示开始markTermination阶段前的heap_live大小；8表示live heap的大小
- 17 MB goal：表示下一次触发GC的内存占用阀值是17MB，等于8MB * 2，向上取整
- 8 P ：本次gc共有多少个P

## FAQ
### GODEBUG
参考
- [用 GODEBUG 看调度跟踪](https://github.com/EDDYCJY/blog/blob/master/tools/godebug-sched.md)

GODEBUG 变量可以控制运行时内的调试变量，参数以逗号分隔，格式为`name=val`, 比如：
- gctrace=1: 让位于Go程序中的运行时在每次GC执行时输出此次GC相关的跟踪信息
- schedtrace：设置 schedtrace=X 参数可以使运行时在每 X 毫秒发出一行goroutine调度器的摘要信息到标准 err 输出中
- scheddetail：设置 schedtrace=X 和 scheddetail=1 可以使运行时在每 X 毫秒发出一次详细的调度器的多行信息，信息内容主要包括调度程序、处理器、OS 线程 和 Goroutine 的状态

### GOGC
除了显式调用runtime.GC强制运行GC外, Go还提供了一个可以调节GC触发时机的环境变量:GOGC.

GOGC是一个整数值 ,默认为100。这个值代表一个比例值,100表示100%. 这个比例值的分子是上一次GC结束后到当前时刻**新分配**的堆内
存大小(设为M),分母则是上一次GC结束后堆内存上的**活动对象**内存数据(live data,可达内存)的大小(设为N).

> 在两次GC之间新分配的堆内存在第二次GC启动的时候不一定都是活动对象占用的内存,也可能是刚分配后不久就处于非活动状态了(没有指针指向这个内存对象)

Go运行时实时监控当前堆内存状态,如果当前堆内存的M/N的值等于GOGC/100,则会再次触发运行GC.

如果程序在性能和响应延迟方面遇到瓶颈,可以大胆地通过调整GOGC的值进行调优,直到找到并验证最适合该Go程序的GOGC值. 

### gc历史
ref:
- [go-GC演化史](https://haoxuebing.github.io/go%E8%BF%9B%E9%98%B6/go-GC%E6%BC%94%E5%8C%96%E5%8F%B2.html)
- [Golang垃圾回收(GC)介绍](https://study.disign.me/article/202502/35.Golang%E5%9E%83%E5%9C%BE%E5%9B%9E%E6%94%B6(GC)%E4%BB%8B%E7%BB%8D.md)

1. go1.5前: STW垃圾回收器

    在GC期间停止程序带来的延迟过大让很多对实时有严格要求的程序无法接受,并且在多核时代这会带来计算资源的严重浪费
1. go1.5: 基于三色标记清除的并发垃圾回收器+屏障机制

    将GC延迟降到10ms以下, 但因在并发标记过程中,只要没有停止程序,用户程序就可以继续分配内存, 这导致了Go运行时无法精确控制堆内存大小
1. go1.8: 三色标记法+混合写屏障机制

### 抢占
Go1.13及以前版本的抢占是协作式的,只在有函数调用的地方才能插入抢占代码(埋点). Go1.14版本增加了基于系统信号的异步抢占调度, 解决deadloop`for{}`无法被抢占.