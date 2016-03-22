# go tool pprof

在Go语言中，我们可以通过标准库的代码包runtime和runtime/pprof中的程序来生成三种包含实时性数据的概要文件，分别是
1. CPU概要文件
1. 内存概要文件
1. 程序阻塞概要文件

参考:
- [go tool pprof](https://github.com/hyper-carrot/go_command_tutorial/blob/master/0.12.md)
- [google-perftools](http://google-perftools.googlecode.com/svn/trunk/doc/)

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
