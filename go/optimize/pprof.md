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

### `net/http/pprof`
Go在net/http/pprof包中还提供了一种更为高级的针对独立程序的性能数据采集方式, 这种方式尤其适合那些内置了HTTP服务的独立程序.

如果独立程序的代码中没有使用http包的默认请求路由器DefaultServeMux,那么就需要重新在新的路由器上为pprof包提供的性能数据采集方法注册服务端点:
```go
mux := http.NewServeMux()
mux.HandleFunc("/debug/pprof/", pprof.Index)
mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
```

使用:
1. 访问`http://localhost:8080/debug/pprof/`
1. 以cpu举例: web上点击cpu采集后, 该服务默认会发起一次持续30秒的性能采集,得到
的数据文件会由浏览器自动下载到本地. 自定义采集时长60s用`http://localhost:8080/debug/pprof/profile?seconds=60`

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
go test -bench=. [-benchtime=2s] -cpuprofile=cpu.prof [-memprofile=mem.prof] [-blockprofile=block.prof]
```
`-cpuprofile`为指定CPU概要文件的文件名,此时go test同时生成xxx.test和cpu.prof文件.

## 分析
web访问http://ip:6060/debug/pprof/即可

### cpu分析
```
go tool pprof water-http-routing-benchmark.test  cpu.prof
```
flat  flat%   sum%        cum   cum%
- flat : 当前函数在采集过程中的运行时间
- flat% : 当前函数运行时间占总时间的百分百(flat/总时间)
- sum% : 累积百分比=本行flat% + 排在该行前面**所有行**的flat%值的累加和
- cum : 累计运行时间,即当前函数以及当前函数直接或间接调用的函数运行所占时间的总和. 越是接近函数调用栈底层的代码, 其cum列的值越大
- cum% : 累计运行时间占总时间的百分比(cum/总时间)
- col6 : 函数名

pprof字命令:
1. top [<N>]: 默认N=10, topN命令的输出结果默认按flat(flat%)从大到小的顺序输出
1. `top -cum`: 按cum值从大到小的顺序输出
1. `list <包.函数名>`: 列出函数对应的源码, 比如`list main.main`
1. png : 生成png, 需要依赖graphviz

	cum%较大的叶子节点(用黑色粗体标出)
1. svg

其他命令:
1. `go tool pprof standalone_app cpu.prof`: 通过net/http/pprof注册的性能采集数据服务端点获取数据并剖析
1. `go tool pprof http://localhost:8080/debug/pprof/profile`
1. `go tool pprof -http=:9090 pprof_standalone1_cpu.prof` : Web图形化方式

	VIEW:
	- Top视图等价于命令行交互模式下的topN命令输出
	- Source视图等价于命令行交互模式下的list命令输出
	- Flame Graph视图即火焰图, Go 1.10版本在go工具链中添加了对火焰图的支持. 通过火焰图可以快速、准确地识别出执行最频繁的代码路径, 因此它多用于对CPU类型采集数据的辅助剖析(其他类型的性能采样数据也有对应的火焰图,比如内存分配)

		go tool pprof在浏览器中呈现出的火焰图与标准火焰图有些差异: 它是**倒置**的,即调用栈最顶端的函数在最下方。在这样一幅倒置火焰图中, **y轴表示函数调用栈,每一层都是一个函数。调用栈越深,火焰越高**。倒置火焰图每个函数调用栈的最下方就是正在执行的函数,上方都是它的父函数.

		**火焰图的x轴表示抽样数量, 如果一个函数在x轴上占据的宽度越宽,就表示它被抽样到的次数越多,即执行的时间越长**。倒置火焰图就是看最下面的哪个函数占据的**宽度最大,这样的函数可能存在性能问题**.

### mem
ref:
- [golang 内存泄漏分析案例](https://www.cnblogs.com/zhanchenjin/p/17101573.html)

参数:
- `-alloc_space`: 当前pprof将呈现程序运行期间所有内存分配的采样数据(即使该分配的内存在最后一次采样时已经被释放)
- `-inuse_space`: 内存数据采样结束时依然在用的内存

	优先使用-inuse_space来分析，因为直接分析导致问题的现场比分析历史数据肯定要直观的多，一个函数alloc_space多不一定就代表它会导致进程的RSS高.

1. 直接heap分析当前的内存:
```bash
# go tool pprof [-alloc_space/-inuse_space] http://ip:8899/debug/pprof/heap
(pprof) top
(pprof) list <func> # 泄漏函数, 可以看到哪行泄漏了
(pprof) sample_index = inuse_space # 实现切换内存数据呈现类型
```

> 最新方式: `-alloc_space`->`-sample_index=alloc_space`

3. 使用base能够对比两个profile文件的差别
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

## expvar辅助定位性能问题
使用expvar获取应用的内省(introspection)数据(json), 辅助定位性能问题.

> 应用的内省(introspection)数据: 反映应用运行状态的数据

Go标准库中的expvar包提供了一种输出应用内部状态信息的标准化方案, 这个方案标准化了以下三方面内容:
1. 数据输出接口形式
1. 输出数据的编码格式
1. 用户自定义性能指标的方法

### go内置
如果应用程序本身并没有使用默认DefaultServeMux,那么需要手动将expvar包的服务端点注册到应用程序所使用的“路由
器”上:
```go
import "expvar"

mux := http.NewServeMux()
mux.Handle("/debug/vars", expvar.Handler())
```

访问: `http://localhost:8080/debug/vars`

字段:
1. cmdline: 应用名
1. memstats: 输出的数据对应的是runtime.Memstats结构体(会随go版本而变化),反映的是应用在运行期间堆内存分配、栈内存分配及GC的状态

### 自定义输出
所有实现了`Var`接口类型的变量都可以被发布(`expvar.Publish`)并作为输出的应用内部状态的一部分.

在设计能反映Go应用内部状态的自定义指标时,经常会设计下面两类指标:
- 测量型: 这类指标是数字,支持上下增减

	。定期获取该指标的快照。常见的CPU、内存使用率等指标都可归为此类型。在业务层面,当前网站上的在线访客数量、当前业务系统平均响应延迟等都属于这类指标。
- 计数型:这类指标也是数字,它的特点是随着时间的推移,其数值不会减少. 虽然它们永远不会减少,但有时可以将其重置为零,它们会再次开始递增

	系统正常运行时间、某端口收发包的字节数、24小时内入队列的消息数量等都是此类指标。计数型指标的一个优势在于可以用来计算变化率:将T+1时刻获取的指标值与T时刻的指标值做比较,
	
expvar包提供了对常用指标类型的原生支持

开源工具expvarmon支持将从expvar输出的数据以基于终端的图形化方式展示出来, 这种方式可以让开发者以最快的速度看到自定义的指标数据.