# go tool pprof
ref:
- [go tool pprof](https://github.com/hyper-carrot/go_command_tutorial/blob/master/0.12.md)
- [google-perftools](http://google-perftools.googlecode.com/svn/trunk/doc/)

Go内置了对代码进行性能剖析的工具:pprof, 其源自Google Perf Tools工具套件.

注意: **采样不是免费的,因此一次采样尽量仅采集一种类型的数据,不要同时采样多种类型的数据,避免相互干扰采样结果**

在Go语言中，我们可以通过标准库的代码包runtime和runtime/pprof中的程序来生成三种包含实时性数据的概要文件，分别是
1. CPU概要文件

	Go运行时会每隔一段短暂的时间(10ms)就中断一次(由SIGPROF信号引发)并记录当前所有goroutine的函数栈信息(存入cpu.prof)

1. 内存概要文件

	堆内存分配的采样频率可配置,默认每1000次堆内存分配会做一次采样(存入mem.prof)
1. 锁竞争(mutex.prof)

	锁竞争采样数据记录了当前Go程序中互斥锁争用导致延迟的操作
1. 程序阻塞概要文件(block.prof)

	该类型采样数据记录的是goroutine在某共享资源(一般是由同步原语保护)上的阻塞时间,包括从无缓冲channel收发数据、阻塞在一个已经被其他goroutine锁住的互斥锁、向一个满了的channel发送数据或从一个空的channel接收数据等

ps: 锁竞争和程序阻塞在默认情况下是不启用的

使用pprof对程序进行性能剖析的工作一般分为两个阶段:数据采集和数据剖析.

Go目前主要支持两种性能数据采集方式:
1. 通过性能基准测试进行数据采集

	比如
	```go
	go test -bench . xxx_test.go -cpuprofile=cpu.prof
	go test -bench . xxx_test.go -memprofile=mem.prof
	go test -bench . xxx_test.go -blockprofile=block.prof
	go test -bench . xxx_test.go -mutexprofile=mutex.prof
	```

	> 一旦开启性能数据采集(比如传入-cpuprofile),go test的`-c`选项便会自动开启
1. 独立程序的性能数据采集

	通过标准库runtime/pprof和runtime包提供的低级API对独立程序进行性能数据采集

## 采集

### 手工配置采集
```go
// --- cpu ---
import (  
    "runtime/pprof"  // 引用pprof package  
    "os"  
)  
func main() {  
    f, _ := os.Create("profile_file")  
    // 在默认情况下，在启动CPU使用情况记录操作之后，
    // 运行时系统就会以每秒100次的频率将取样数据写入到CPU概要文件
    pprof.StartCPUProfile(f)  // 开始cpu profile，结果写到文件f中  
    // pprof.StopCPUProfile函数通过把CPU使用情况取样的频率设置为0来停止取样操作,
    // 并且，只有当所有CPU使用情况记录都被写入到CPU概要文件之后，
    // pprof.StopCPUProfile函数才会退出,从而保证了CPU概要文件的完整性.
    defer pprof.StopCPUProfile()  // 结束profile  
    ...  
}

// --- mem ---
// 内存概要文件用于保存在用户程序执行期间的内存使用情况.
// 这里所说的内存使用情况，其实就是程序运行过程中堆内存的分配情况。
// 对内存使用情况进行取样的程序会假定取样间隔在用户程序的运行期间内都是一成不变的，
// 并且等于runtime.MemProfileRate变量的当前值,无法中途变更.
import (
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	f, _ := os.Create("profile_file")

	// 分析器将会在每分配指定的字节数量(memProfileRate)后对内存使用情况进行取样.
	// runtime.MemProfileRate变量的默认值是512 * 1024，即512K个字节.
	// 只有当我们显式的将0赋给runtime.MemProfileRate变量之后，才会取消取样操作.
  memProfileRate := 1
	runtime.MemProfileRate = memProfileRate
	// 在默认情况下，内存使用情况的取样数据只会被保存在运行时内存中,
	// 保存到文件的操作只能由我们自己来完成
	defer func() {
		pprof.WriteHeapProfile(f)
	}()

	...
}

// --- mutex ---
func main() {
	f, _ := os.Create("profile_file")

	defer func() {
		pprof.Lookup("mutex").WriteTo(f, 0)
	}()
}

// --- block ---
// 程序阻塞概要文件用于保存用户程序中的Goroutine阻塞事件的记录
import (
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	f, _ := os.Create("profile_file")

	// 分析器的取样间隔，单位是次,
	// 分析器会在每发生几次Goroutine阻塞事件时对这些事件进行取样,默认是1.
	// 将这个取样间隔设置为0或者负数，那么取样操作就会被取消.
	blockProfileRate := 1
	runtime.SetBlockProfileRate(rate)
	// 通过函数pprof.Lookup将保存在运行时内存中的内存使用情况记录取出，
	// 并在记录的实例上调用WriteTo方法将记录写入到文件中
	defer func() {
		pprof.Lookup("block").WriteTo(f, 0)
	}()
  ...
}
```

### go test采集
```
go test -bench=. -cpuprofile=cpu.prof
```
`-cpuprofile`为指定CPU概要文件的文件名,此时go test同时生成xxx.test和cpu.prof文件.

## 分析
web访问http://ip:6060/debug/pprof/即可

### cpu分析
```
go tool pprof water-http-routing-benchmark.test  cpu.prof
```
flat  flat%   sum%        cum   cum%
- col1 : 当前函数运行时间
- col2 : 当前函数运行时间占总时间的百分百(col1/总时间)
- col3 : 累积百分比=本项col2+上一项col2
- col4 : 累计运行时间,即当前函数以及当前函数直接或间接调用的函数运行所占的时间.
- col5 : 累计运行时间占总时间的百分比(col4/总时间)
- col6 : 函数名


### mem
ref:
- [golang 内存泄漏分析案例](https://www.cnblogs.com/zhanchenjin/p/17101573.html)

`go tool pprof -alloc_space/-inuse_space http://ip:8899/debug/pprof/heap`
优先使用-inuse_space来分析，因为直接分析导致问题的现场比分析历史数据肯定要直观的多，一个函数alloc_space多不一定就代表它会导致进程的RSS高.

1. 直接heap分析当前的内存:
```bash
# go tool pprof http://ip:8899/debug/pprof/heap
(pprof) top
(pprof) list <func> # 泄漏函数, 可以看到哪行泄漏了
```

2. 使用base能够对比两个profile文件的差别
```bash
# go tool pprof http://ip:8899/debug/pprof/heap # get base
# go tool pprof http://ip:8899/debug/pprof/heap # get latest
# go tool pprof -base pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.001.pb.gz pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
```

3. 查看heap图片: `go tool pprof -png http://localhost:6060/debug/pprof/heap > heap.png`

### goroutine泄露
ref:
- [golang 内存泄漏分析案例](https://www.cnblogs.com/zhanchenjin/p/17101573.html)

方法与上节mem类似

## GC
- [如何监控 golang 程序的垃圾回收](https://holys.im/2016/07/01/monitor-golang-gc/)
