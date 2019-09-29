# goroutine
参考:
 - [Go调度器系列](http://lessisbetter.site/2019/03/10/golang-scheduler-1-history/)
 - [Go语言——goroutine并发模型](https://www.jianshu.com/p/f9024e250ac6)
 - [Goroutine并发调度模型深度解析](https://juejin.im/entry/5b2878c7f265da5977596ae2)
 - [Go 调度模型](https://wudaijun.com/2018/01/go-scheduler/)

## scheduler
调度器由三方面实体构成：
1. M：真正的内核OS线程，类似于 POSIX 的标准线程
1. G：goroutine，它拥有自己的栈、指令指针和维护其他调度相关的信息
1. P：代表调度上下文，可将其视为一个局部调度器，使Golang代码跑在一个线程上. P 的数量可由 runtime.GOMAXPROCS() 进行设置，它代表了真正的并发能力，即可有多少个 goroutine 同时运行

一个M对应一个P，一个P下面挂多个G，但一个时刻只有一个G在跑，其余都是放入等待队列，等待下一次切换时使用.
Goroutine调度器和OS调度器是通过M结合起来的，每个M都代表了1个内核线程，OS调度器负责把内核线程分配到CPU的核上执行.
每个M都有一个特殊的G即g0.用于执行调度，gc，栈管理等任务，所以g0的栈称为调度栈. g0使用的是os线程的栈, 其不会自动增长，不会被gc

GMP的可视化方法:
1. [`go tool trace`](https://mp.weixin.qq.com/s/nf_-AH_LeBN3913Pt6CzQQ)
1. [GODEBUG](http://lessisbetter.site/2019/03/26/golang-scheduler-2-macro-view/, https://colobu.com/2016/04/19/Scheduler-Tracing-In-Go/)

> G-P-M模型的定义放在src/runtime/runtime2.go里面，而调度过程则放在了src/runtime/proc.go里

### 队列
Go调度器有两个不同的运行队列：
- GRQ : 全局运行队列，尚未分配给P的G
- LRQ : 本地运行队列，每个P都有一个LRQ，用于管理分配给P执行的G

### 设计思想
调度器的有两大思想：
1. 复用线程：协程本身就是运行在一组线程之上，不需要频繁的创建、销毁线程，而是对线程的复用。在调度器中复用线程还有2个体现：
    1. work stealing，当本线程无可运行的G时，尝试从其他线程绑定的P偷取(一半)G，而不是销毁线程
    1. hand off，当本线程因为G进行系统调用阻塞时，线程释放绑定的P，把P转移给其他空闲的线程执行
1. 利用并行：GOMAXPROCS设置P的数量，当GOMAXPROCS大于1时，就最多有GOMAXPROCS个线程处于运行状态，这些线程可能分布在多个CPU核上同时运行，使得并发利用并行. 另外，GOMAXPROCS也限制了并发的程度，比如GOMAXPROCS = 核数/2，则最多利用了一半的CPU核进行并行.

调度器的两小策略：
1. 抢占：在coroutine中要等待一个协程主动让出CPU才执行下一个协程，在Go中，一个goroutine最多占用CPU 10ms，防止其他goroutine被饿死，这就是goroutine不同于coroutine的一个地方
1. 全局G队列：在新的调度器中依然有全局G队列，但功能已经被弱化了，当M执行work stealing从其他P偷不到G时，它可以从全局G队列(在`/src/runtime/runtime2.go`的`schedt`中)获取G

### goroutine切换
Go runtime会在下面的goroutine被阻塞的情况下运行另外一个goroutine：

- blocking syscall (for example opening a file)
- network input
- channel operations
- primitives in the sync package

## FAQ
### 线程模型
传统的协程库属于用户级线程模型，而goroutine和它的Go Scheduler在底层实现上其实是属于[混合线程模型](https://github.com/meilihao/programming-interface/master/process/process.md#进程和线程).

### 栈大小
[线程是有固定的栈的，默认是2MB(但进程即主线程是8MB)](http://man7.org/linux/man-pages/man3/pthread_create.3.html)，当然，不同系统可能大小不太一样，但是的确都是固定分配的.这个栈用于保存局部变量，用于在函数切换时使用.

但是对于goroutine这种轻量级的协程来说，一个大小固定的栈可能会导致资源浪费：比如一个协程里面只print了一个语句，那么栈基本没怎么用；当然，也有可能嵌套调用很深，那么可能也不够用.
所以go采用了动态扩张收缩的策略(g0除外)：初始化为2KB(`src/runtime/stack.go._StackMin`)，最大可扩张到1GB.

> [user stack的大小是固定的，Linux中默认为8192KB，运行时内存占用超过上限，程序会崩溃掉并报告segment错误](https://studygolang.com/articles/10597). 为了修复这个问题，我们可以调大内核参数中的stack size, 或者在创建线程时显式地传入所需要大小的内存块.

> [聊一聊goroutine stack](https://kirk91.github.io/posts/2d571d09/)
> Linux中栈默认为8192KB(`ulimit -a |grep  "Maximum stack size"`)

### goroutine没有id
go语言设计之初考虑的，防止被滥用，所以你不能在一个协程中杀死另外一个协程