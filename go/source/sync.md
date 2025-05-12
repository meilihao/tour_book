# sync
ref:
- [go锁mutex与RWMutex](https://cloud.tencent.com/developer/article/2481775)
- [多图详解Go的互斥锁Mutex](https://www.luozhiyun.com/archives/413)

```go
// $GOROOT/src/sync/mutex.go
type Mutex struct {
    state int32
    sema uint32
}

// /internal/sync/mutex.go
const (
	mutexLocked = 1 << iota // mutex is locked // 已锁定
	mutexWoken // 是否有协程抢占锁
	mutexStarving // 进入饥饿状态
	mutexWaiterShift = iota
    ...
}
```

Mutex由两个字段state和sema组成:
- state:表示当前互斥锁的状态
- sema:用于控制锁状态的信号量

Mutex互斥锁有两种模式:
- 正常模式（Normal Mode）

    正常模式是Mutex互斥锁的默认模式。在正常模式下，等待者按先进先出 (FIFO) 的顺序排队，但唤醒的等待者并不拥有互斥锁，而是与新到达的 Goroutine 竞争所有权.

    性能最好
- 饥饿模式（Starvation Mode）

    饥饿模式是一种公平的模式。当某个goroutine连续多次尝试获取锁但一直失败时，Mutex可能会切换到饥饿模式。在饥饿模式下，互斥锁的所有权直接从解锁的 Goroutine 移交给队列最前面的等待者. 新到达的 Goroutine 不会尝试获取互斥锁，而是会将自己排在等待队列的尾部

    取锁公平, 但性能不一定最佳

Mutex和RWMutex适用场景:
1. 在并发量较小的情况下,互斥锁性能更好;随着并发量增大,互斥锁的竞争激烈,导致加锁和解锁性能下降
1. 读写锁的读锁性能并未随并发量的增大而发生较大变化, 性能始终恒定在40ns左右
1. 在并发量较大的情况下,读写锁的写锁性能比互斥锁、读写锁的读锁都差,并且随着并发量增大,其写锁性能有继续下降的趋势

总结: 读写锁适合应用在具有一定并发量且读多写少的场合

sync.Pool是一个数据对象缓存池,它具有如下特点:
1. 它是goroutine并发安全的,可以被多个goroutine同时使用
1. 放入该缓存池中的数据对象的生命是暂时的,随时都可能被垃圾回收掉
1. 缓存池中的数据对象是可以重复利用的,这样可以在一定程度上降低数据对象重新分配的频度,减轻GC的压力
1. sync.Pool为每个P(goroutine调度模型中的P)单独建立一个local缓存池,进一步降低高并发下对锁的争抢