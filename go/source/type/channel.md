# channle
- [图解Go的channel底层原理](https://blog.csdn.net/i6448038/article/details/89303398)
- [多图详解Go中的Channel源码](https://www.luozhiyun.com/archives/427)

```go
// /runtime/chan.go
type hchan struct {
	qcount   uint           // total data in the queue // 队列中数据总数
	dataqsiz uint           // size of the circular queue // 环形队列的size即构造 channel 时指定的 buf 大小
	buf      unsafe.Pointer // points to an array of dataqsiz elements // 指向环形队列
	elemsize uint16 // 元素大小
	synctest bool // true if created in a synctest bubble
	closed   uint32 // 是否关闭, 0表示未关闭
	timer    *timer // timer feeding this chan
	elemtype *_type // element type // 元素类型
	sendx    uint   // send index // 已发送元素的索引位置
	recvx    uint   // receive index // 已接收元素的索引位置
	recvq    waitq  // list of recv waiters // 等待接收的 goroutines
	sendq    waitq  // list of send waiters // 等待发送的 goroutines

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct { // 是双向链表
	first *sudog
	last  *sudog
}
```

channel是异步的, 有3种状态:
1. nil
2. active, 正常的channel
3. closed