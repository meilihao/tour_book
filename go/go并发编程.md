# go并发编程
参考:
- [<<Go并发编程实践 - 晁岳攀.pdf>>]

## 基本同步原语
锁是一种并发编程中的同步原语（Synchronization Primitives），它能保证多个 Goroutine 在访问同一片内存时不会出现竞争条件（Race condition）等问题.

1. mutex

	- 非重入锁

	互斥锁有两种状态: 正常状态和饥饿状态.

	在正常状态下, 所有等待锁的goroutine按照FIFO顺序等待. 唤醒的goroutine不会直接拥有锁, 而是会和新请求锁的goroutine竞争锁的拥有.
	如果一个等待的goroutine超过1ms没有获取锁, 那么它将会把锁转变成饥饿模式.

	在饥饿模式下, 锁的拥有权将从unlock的gorutine直接交给等待队列中的第一个. 新来的goroutine将不会尝试去获得锁, 即使锁看起来是unlock状态, 也不会去尝试自旋操作, 而是放在等待队列的尾部.

	如果一个等待的goroutine获取了锁, 并且满足以下任一条件. 它会将锁的状态转换成正常状态:
	1. 它是队列中的最后一个
	1. 它等待的时间小于1ms

1. RWMutex

	- 可以被很多的reader持有, 或者被一个writer持有
	- 适合大并发read的场景
	- writer的Lock相对后续的reader的RLock优先级高
	- 禁止递归读锁

1. Cond

	Condition Variables是一组等待同一个条件的goroutine的容器.

	Monitor=Mutex + Condition Variables

	每个Cond和一个Locker相关联.
	改变条件或者调用Wait需要获取锁.

1. Waitgroup

	- 等待一组goroutine完成
	- Add可以是负数, 但Waitgroup计数器小于0会panic
	- 当计数器为0时, 阻塞在Wait方法的goroutine都会被释放
	- 可重用

1. Once

	- 只执行一次
	- 避免死锁
	- 即使f panic, Once()也认为它完成了.

1. Pool

	- 临时对象pool
	- 可能在任何时候任意的对象都可能被移除
	- 可以安全地并发访问
	- 装箱/拆箱

1. Sync.Map

	- 两个cases:
	
		- 设置一次, 多次读
		- 多个goroutine并发的读,写, 更新不同的key
	- 装箱/拆箱
	- Range进行遍历, 可能会加锁
	- 没有Len方法, 且也不会添加

## 扩展同步原语

1. ReenstrantLock
1. Semaphore

	golang.org/x/sync/semaphore
1. SingleFlight

	golang.org/x/sync/singleflignt
1. ErrGroup

	golang.org/x/sync/semaphore

	- Wait会等待所有的goroutine执行完后才释放
	- 如果想遇到第一个err就返回, 使用Context
1. SpinLock

	- 自旋锁
	- 公平性
	- 处理器忙等待
1. fslock

	github.com/juju/fslock, 跨进程的Mutex

1. concurrent-map

	github.com/orcaman/concurrent-map

## 原子操作

	Add的减法:
	- AddUint32(&x, ^uint32(0)), AddUint64(&x, ^uint64(0))
	- AddUint32(&x, ^uint32(0)), AddUint64(&x, ^uint64(0))

## channel
特殊情况

| 操作        | nil    |  not empty  |empty| full| not full|
| --------   | -----:   | :----: |
| Receive        | block      |   value    |block| value|value|
| Send        | block     |   write value   |write value|block|write value|
| close        | panic    |   closed, drained read, return zero value    |closed,return zero value|closed,drained read, return zero value|closed, drained read, return zero value|

## 内存模型
参考:
- [Go 内存模型和Happens Before关系](https://zhuanlan.zhihu.com/p/29108170)

内存模型描述了线程/goroutine通过内存的交互, 以及对数据的共享使用.

> 如果是前端或者使用node的程序员，那么压根就不需要清楚这些，毕竟目前js始终只有一个线程在跑

为什么需要定义Happens Before关系来保证内存操作的可见性呢？原因是没有限制的情况下，编译器和CPU使用的各种优化，会对此造成影响.

具体的来说就是操作重排序和CPU CacheLine缓存同步：
- 操作重排序 :　现代CPU通常是流水线架构，且具有多个核心，这样多条指令就可以同时执行。然而有时候出现一条指令需要等待之前指令的结果，或是其他造成指令执行需要延迟的情况。这个时候可以先执行下一条已经准备好的指令，以尽可能高效的利用CPU。操作重排序可以在两个阶段出现：
	- 编译器指令重排序
	- CPU乱序执行
- CPU 多核心间独立Cache Line的同步问题 : 多核CPU通常有自己的一级缓存和二级缓存，访问缓存的数据很快. 但是如果缓存没有同步到主存和其他核心的缓存，其他核心读取缓存就会读到过期的数据.

因此Happens Before就是对编译器和CPU的限制，禁止违反Happens Before关系的指令重排序及乱序执行行为，以及必要的情况下保证CacheLine的数据更新等.

go内存模型:
- 定义了对同一个变量, 如何保证在一个goroutine对此变量读的时候, 能观察到其他goroutine对该变量的写.
- 修改一个同时被多个goroutine并发访问的变量的时候, 需要串行化访问. 通过channel或其他同步原语实现串行化访问.

memory order guarntee

	单个goroutine内:
	- 读写执行的顺序与程序定义的顺序一致
	- 乱序执行不影响程序的行为

happen before

	内存操作的偏序(partial order)

	定义了两个事件结果的先后顺序:
	- a -> b, a happens before b 或 b hanpens after a
	- a does not happen before b, and does not happen after b, a和b同时发生

### happen-before
1. 单线程
在单线程环境下，所有的表达式，按照代码中的先后顺序，具有Happens Before关系. 这并不是说编译器或者CPU不能做重排序，只要优化没有影响到Happens Before关系就是可以的.
1. init函数
	
	- 如果包P1中导入了包P2，则P2中的init函数Happens Before 所有P1中的操作
	- main函数Happens After 所有的init函数
1. goroutine
	- goroutine的创建happens before所有此goroutine中的操作
	- goroutine的销毁happens after所有此goroutine中的操作

1. Channel
	- 第n个send一定happen before第n个receive完成, 不管是buffered channel还是unbuffered channel
	- 对于capacity为m的channel, 第n个receive一定happen before第(n+m)send完成.
	- m=0 unbuffered, 第n个receiver一定happen before第n个send完成
	- 对channel的close操作happens Before receive 端的收到关闭通知, 得到通知意味着receive收到一个因为channel close而收到的零值.

	首先注意这里面，send和send完成，这是两个事件，receive和receive完成也是两个事件.

	然后，Buffered Channel这里有个坑，它的Happens Before保证比UnBuffered 弱，这个弱只在【在receive之前写，在send之后读】这种情况下有问题, 而【在send之前写，在receive之后读】，这样用是没问题的. 这也是我们通常写程序常用的模式，千万注意这里不要弄错！	


1. Lock=Mutex + RWMutex
	- 对于Mutex/RWMutex m, 第n个成功的m.Unlock一定happen before 第n+1 m.Lock方法调用的返回.
	- 对于RWMutex rw, 如果它的第n个rw.Lock已返回, 那么它的第n个成功的rw.Unlock的方法调用一定happen before任何一个rw.RLock方法调用的返回(它们happen after 第n个rw.Lock方法调用返回)
	- 对于RWMutex rw, 如果它的第n个rw.RLock已返回, 接着第m`(m<n)`个rw.RUnlock的方法调用一定happen before任意的rw.Lock(它们happen after 第n个rw.RLock方法调用返回之后)

简单理解就是这一次的Lock总是Happens After上一次的Unlock，读写锁的RLock Happens After上一次的UnLock，其对应的RUnlock Happens Before 下一次的Lock

1. Waitgroup:
	- 对于Waitgroup b, 对于其计数器不是0的时候, 假如此时刻后有一组wg.Add(n), 并且确信只有最后一组方法调用使其计数器最后复原为0, 那么这组wg.Add方法调用一定happen before这一时候之后发生的wg.Wait
	- wg.Done()=wg.Add(-1)

1. Once:
	- once.Do方法的执行一定happen before任何一个once.Do方法的返回

1. atomic:
	- 没有官方的保证
	- 建议不要依赖atomic保证内存的顺序

## FAQ
### channel vs Mutex
Channel:
- 传递数据的owner
- 分发任务单元
- 交流异步结果
- 任务编排

Mutex:
- cache
- 状态
- 临界区