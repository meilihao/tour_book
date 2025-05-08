# time

## Time
```go
// go1.24
type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
```

Time由三个字段组成, 同时表示两种时间——挂钟时间(wall time)和单调时间(monotonic time), 并且精度级别为纳秒.
- wall: 挂钟时间. 最高比特位是一个名为hasMonotonic的标志比特位. 

    - hasMonotonic=1: time.Time表示的即时时间中既包含挂钟时间, 也包含单调时间. 62~30bit表示秒(距1885年1月1日的秒数), 29~0(30bit)表示纳秒
    - hasMonotonic=0: 仅表示挂钟时间. 秒数(33bit)两部分均被置为0, 纳秒数(30bit)依旧用于表示挂钟时间的非整数秒部分
- ext: hasMonotonic=1时表示进程启动后的单调流逝时间(ns); hasMonotonic=0时整个用于表示挂钟时间的整秒部分, 其含义为距公元元年1月1日的秒数
- loc: 时区信息, 未设置时默认使用系统时区. 在Linux/macOS上,默认使用的是/etc/localtime指向的时区数据

通过time.Now函数获取的当前时间中则既包含挂钟时间,也包含单调时间. 它是runtime包的runtimeNow函数实现的. time.Now函数在获取当前时间时会考虑时区信息,如果TZ环境变量不为空,那么它将尝试读取该环境变量指定的时区信息并输出对应时区的即时时间表示. 如果TZ有误或为"", 它默认使用UTC时间.

Go 1.9版本中加入了对单调时间的支持,单调时间常被用于两个即时时间之间的比较和间隔计算.

## Timer
```go
// go1.24
type Timer struct {
	C         <-chan Time
	initTimer bool
}

// /src/runtime/runtime2.go
type p struct{
    ...
    // Timer heap.
	timers timers
    ...
}
```

Timer字段:
1. C : 用户层用户接收定时器触发事件的channel

它由runtime newTimer实现

创建Timer的方式:
1. NewTimer
1. AfterFunc
1. After

Go语言标准库提供的timer实质上是由Go runtime自行维护的,而不是操作系统级的定时器资源.

go1.9之前, runtime维护一个由互斥锁保护的全局最小堆(minheap),定时器最小堆的维护操作都要对其互斥锁进行加解锁操作,导致其性能和伸缩性很差。最新的定时器管理调度
方案(Go 1.14)抛弃了全局唯一最小堆方案,而是为每个P(goroutine调度器中的那个P)创建一个定时器最小堆,并通过网络轮询器(netpoller)在运行时调度的协助下对各个定时器最小堆进行统一管理和调度

Go官方文档建议只对如下两种定时器调用Reset方法:
1. 已经停止了的定时器(Stopped)
1. 已经触发过且Timer.C中的数据已经被读空

推荐写法:
```go
if !timer.Stop() {
	select {
		case <-timer.C: // 避免数据已被取走引发的阻塞
		default:
	}
}
timer.Reset(5 * time.Second)
```

## Format
go采用参考时间(reference time)方案