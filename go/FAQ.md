## map

### 不能赋值,提示"cannot assign to ..."
```go
type person struct {
       id   int
       name string
}
persons := make(map[int]person)
persons[0] = person{1, "a"}
persons[0].id=2 // 报错
```
map里面的元素是值拷贝,当前存入的person是不可寻址的,所以不能用persons[0].id这种方法赋值,解决方法:
1. 把element类型换成`*person`
2. 重新赋值
```go
p := person[0]
p.id = 2
person[0] = p
```

类似:
```go
// 接口口转型返回临时对象,只有使用用指针才能修改其状态
u := person{1, "Tom"}
var vi, pi interface{} = u, &u
// vi.(person).name = "Jack" // Error: cannot assign to vi.(person).name
pi.(*person).name = "Jack"
fmt.Printf("%v\n", vi.(person))
fmt.Printf("%v\n", pi.(*person))
```

## interface{}

### 接口断言
```go
package main

import "fmt"

func main() {
	var t T1 = 10
	var i1 I1 = t
	i1.Speak()

	var i2 I2
	i2 = i1
	i2.Speak()

	if i1, ok := i2.(I1); ok { // i2接口里的实际值是否实现了(I1)接口, ok==true
		i1.Step()
	}
}

type T1 int
type T2 int
type T3 int

func (t T1) Speak() {
	fmt.Println("sk", t)
}

func (t T2) Speak() {
	fmt.Println("sk", t)
}

func (t T3) Speak() {
	fmt.Println("sk", t)
}

func (t T1) Step() {
	fmt.Println("sp", t)
}

func (t T3) Step() {
	fmt.Println("sp", t)
}

type I1 interface {
	Speak()
	Step()
}

type I2 interface {
	Speak()
}
```

## http
### 获取http.DefaultClient.Do处理30x Redirect的路径
```go
var DefaultClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		fmt.Println(req.URL)

		via[0].Header.Set("XXX", req.URL.String())

		return nil
	},
}
```

## 其他
### go test -bench时跑了单元测试
go跑`go test -bench`前会跑`go test -run`保证代码的正确性,如果单元测试有错(比如调用了`t.Error()`),后面的性能测试就不会执行.

### 值接收者调用方法时的地址和变量地址不一致
如果使用值接收者声明方法，调用时会使用这个值的一个**副本**来执行.

### goto
如果必须使用 goto，应当**只使用正序的标签**（标签位于 goto 语句之后），且**Labal和 goto 语句之间不能定义新变量(该变量与Labal是相同的代码块级别)**，否则会导致编译失败.

### internal包（内部包）
有些时候 我们希望一些包并非能被所有外部包所导入，但却能被其“临近”的包所导入和访问. 因此Go 1.4引入了"internal"包的概念，导入这种internal包的规则约束如下:

如果导入代码本身不在以"internal"目录的父目录为root的目录树中，那么 不允许其导入路径(import path)中包含internal元素。

例如：
    – `a/b/c/internal/d/e/f`只可以被以`a/b/c`为根的目录树下的代码导入，不能被`a/b/g`下的代码导入
    – `$GOROOT/src/net/http/internal`只能被`net/http`和`net/http/*`的包所导入
    – `$GOPATH/src/mypkg/internal/foo`只能被`$GOPATH/src/mypkg*`包的代码所导入

对于Go 1.4该规则首先强制应用于$GOROOT下。Go 1.5将扩展应用到$GOPATH下.

注：Go 1.4 取消了$GOROOT/src/pkg，标准库都移到$GOROOT/src下了.

### go test 禁用cache
```sh
$ GOCACHE=off go test -v   util.go util_test.go
```

### exec.Command 报错: signal: interrupt
新创建的进程将与发起`exec.Command`的进程位于同一进程组中.这意味着默认情况下， signal将广播到`exec.Command`创建的进程中.

解决方法:
1. 使用类型SysProcAttr属性在命令之前强制新创建的进程位于其自己的进程组中.

```
cmd := exec.Command("sh","-c","xxx")
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid:true}
err := cmd.Run()
```

### exec.Command和nohup
node_server:
```go
func main() {
	for {
		log.Println(os.Getenv("uid"))

		time.Sleep(2*time.Second)
	}
}
```

原有代码
```go
for ... {
detailCmd = fmt.Sprintf("cd %s && env uid=%d nohup ./node_server > %s 2>&1 &",
				"/home/chen/tmpfs",
				0,
				"0.log",
			)

			cmd := exec.Command("/bin/bash", "-c", detailCmd)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			out, err := cmd.CombinedOutput() // 阻塞
}
```
执行的命令是`/bin/bash -c "cd /home/chen/tmpfs && env uid=0 nohup ./node_server > 0.log 2>&1 &"`, 原因未知.
```bash
$ ps -ef|grep node_server
chen     11005     1  0 15:42 pts/6    00:00:00 /bin/bash -c cd /home/chen/tmpfs && env uid=0 nohup ./node_server > 0.log 2>&1 &
chen     11006 11005  0 15:42 pts/6    00:00:00 ./node_server
```

改进
```go
for ... {
detailCmd := fmt.Sprintf("env uid=%d nohup ./node_server > %s 2>&1 &", // 这里必须要重定向,否则cmd无法退出
			0,
			"0.log",
		)
cmd := exec.Command("/bin/bash","-c",detailCmd) // 命令推荐使用绝对路径. 因为nohup启动进程时再执行exec.Command会因为$PATH为空, 而报`executable file not found in $PATH`
		cmd.Dir="/home/chen/tmpfs" // 设置 working directory
		// cmd.Env=[]string{"uid=0"} // 设置env, 不推荐. 因为自行设置导致不会继承父进程的env, 从而导致某些命令出错, 比如`lame`.
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		out, err := cmd.CombinedOutput()
}
```

### hash of unhashable type model.CustomerProfile
```go
{
	...
	se := e.NewSession() //xorm
	defer se.Close()

	err = se.Begin()
	if err != nil {
		sugar.Fatal(err)
	}

	cp := Profile{
		GroupId: 10,
		Phone:   "111",
		Name:    "test1",
		Belong:  "test1",
		Expand: KVExpand{
			Keys:   make([]string, 0),
			Values: make([]string, 0),
		},

		Status:  1,
	}
	_, err = se.Insert(cp)
	if err != nil {
		sugar.Fatal(err) // 报错: panic: runtime error: hash of unhashable type model.CustomerProfile
	}

	sugar.Debug(se.Commit())
}


type Profile struct {
	GroupId int    `xorm:"pk"` // 主账号id
	Phone   string `xorm:"pk"`    // 客户phone
	Name    string                // 客户名称
	Belong  string              // 客户所属公司
	Expand    KVExpand  `xorm:"json" json:"expand"` // 客户资料, kv对形式
	Status    int
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type KVExpand struct{
	Keys []string `json:"keys"`
	Values []string `json:"values"`
}
```

错误信息:
```log
panic: runtime error: hash of unhashable type model.CustomerProfile

goroutine 15 [running]:
github.com/go-xorm/xorm.(*Session).innerInsert.func1(0xd32b00, 0xc0002ca140)
	/home/chen/git/go/src/github.com/go-xorm/xorm/session_insert.go:411 +0x1b9
github.com/go-xorm/xorm.(*Session).innerInsert(0xc0000cf8f8, 0xd32b00, 0xc0002ca140, 0x1, 0x0, 0x0)
	/home/chen/git/go/src/github.com/go-xorm/xorm/session_insert.go:532 +0x195d
github.com/go-xorm/xorm.(*Session).Insert(0xc0000cf8f8, 0xc0000cf6a8, 0x1, 0x1, 0x0, 0x0, 0x0)
	/home/chen/git/go/src/github.com/go-xorm/xorm/session_insert.go:83 +0x5bc
```

当把一个 interface{} 类型值作为键添加到一个字典值的时候，Go语言会先获取这个 interface{} 类型值的实际类型（即动态类型），
然后再使用与之对应的 hash 函数对该值进行 hash 运算，也就是说，interface{} 类型值总是能够被正确地计算出 hash 值.
但是字典类型的键不能是函数类型、字典类型或切片类型，否则会引发一个运行时恐慌，并提示如下：
`panic: runtime error: hash of unhashable type <某个函数类型、字典类型或切片类型的名称>`

解决方法:
将Profile.Expand定义为`*KVExpand`即可正确计算hash值.

> 参考: [Go语言之自定义集合Set](https://www.jb51.net/article/89736.htm)

### struct 比较
进行结构体比较时候，只有相同类型的结构体才可以比较，结构体是否相同不但与属性类型个数有关，还与**属性顺序**相关

### cannot take the address of 常量
```go
package main

const cl = 100
var bl = 123

func main() {
	println(&bl, bl)
	println(&cl, cl)
}
```
### Go 1.9 新特性 Type Alias
```go
func main() {
	type MyInt2 = int
	var i int = 9
	var i2 MyInt2 = i
	fmt.Println(i2)
}
```
```go
package main

import "fmt"

type User struct {
}

type MyUser2 = User // 因为MyUser2完全等价于User，所以具有其所有的方法，并且其中一个新增了方法，另外一个也会有

func (i User) m2() {
	fmt.Println("User.m2")
}

func main() {
	var i2 MyUser2
	i2.m2()
}
```

```go
package main

import "fmt"

type T1 struct {
}
func (t T1) m1(){
    fmt.Println("T1.m1")
}
type T2 = T1
type MyStruct struct {
    T1
    T2
}
func main() {
    my:=MyStruct{}
    my.m1() // ambiguous selector my.m1
}
```
MyStruct没有m1方法, 且T1,T2有重复的方法, 编译器不知道该选择哪个. 改为:
```go
my.T1.m1()
my.T2.m1()
```

### Happens Before
参考:
- [Go 内存模型和Happens Before关系](https://zhuanlan.zhihu.com/p/29108170)

为什么需要定义Happens Before关系来保证内存操作的可见性呢？原因是没有限制的情况下，编译器和CPU使用的各种优化，会对此造成影响.

具体的来说就是操作重排序和CPU CacheLine缓存同步：
- 操作重排序 :　现代CPU通常是流水线架构，且具有多个核心，这样多条指令就可以同时执行。然而有时候出现一条指令需要等待之前指令的结果，或是其他造成指令执行需要延迟的情况。这个时候可以先执行下一条已经准备好的指令，以尽可能高效的利用CPU。操作重排序可以在两个阶段出现：
	- 编译器指令重排序
	- CPU乱序执行
- CPU 多核心间独立Cache Line的同步问题 : 多核CPU通常有自己的一级缓存和二级缓存，访问缓存的数据很快. 但是如果缓存没有同步到主存和其他核心的缓存，其他核心读取缓存就会读到过期的数据.

因此Happens Before就是对编译器和CPU的限制，禁止违反Happens Before关系的指令重排序及乱序执行行为，以及必要的情况下保证CacheLine的数据更新等.

Go 中定义的Happens Before保证:
1. 单线程
在单线程环境下，所有的表达式，按照代码中的先后顺序，具有Happens Before关系. 这并不是说编译器或者CPU不能做重排序，只要优化没有影响到Happens Before关系就是可以的.
1. init函数
	- 如果包P1中导入了包P2，则P2中的init函数Happens Before 所有P1中的操作
	- main函数Happens After 所有的init函数
1. Goroutine
Goroutine的创建Happens Before所有此Goroutine中的操作
Goroutine的销毁Happens After所有此Goroutine中的操作
1. Channel
	- 对一个元素的send操作Happens Before对应的receive 完成操作
	- 对channel的close操作Happens Before receive 端的收到关闭通知操作
	- 对于Unbuffered Channel，对一个元素的receive 操作Happens Before对应的send完成操作
	- 对于Buffered Channel，假设Channel 的buffer 大小为C，那么对第k个元素的receive操作，Happens Before第k+C个send完成操作, 可以看出上一条Unbuffered Channel规则就是这条规则C=0时的特例

首先注意这里面，send和send完成，这是两个事件，receive和receive完成也是两个事件.

然后，Buffered Channel这里有个坑，它的Happens Before保证比UnBuffered 弱，这个弱只在【在receive之前写，在send之后读】这种情况下有问题, 而【在send之前写，在receive之后读】，这样用是没问题的. 这也是我们通常写程序常用的模式，千万注意这里不要弄错！

1. Lock
Go里面有Mutex和RWMutex两种锁，RWMutex除了支持互斥的Lock/Unlock，还支持共享的RLock/RUnlock.

	- 对于一个Mutex/RWMutex，设n < m，则第n个Unlock操作Happens Before第m个Lock操作
	- 对于一个RWMutex，存在数值n，RLock操作Happens After 第n个UnLock，其对应的RUnLockHappens Before 第n+1个Lock操作

简单理解就是这一次的Lock总是Happens After上一次的Unlock，读写锁的RLock Happens After上一次的UnLock，其对应的RUnlock Happens Before 下一次的Lock

1. Once
once.Do中执行的操作，Happens Before 任何一个once.Do调用的返回

> 如果是前端或者使用node的程序员，那么你压根就不需要清楚这些，毕竟目前js始终只有一个线程在跑

### Golang中实现典型的fork调用
https://github.com/moby/moby/tree/master/pkg/reexec

Golang中没有提供 fork 调用，只有

syscall.ForkExec
os.StartProcess
exec.Cmd

这三个都类似于 fork + exec，但是没有类似C中的fork调用，在fork之后根据返回的pid 然后进入不同的函数。原因在：

https://stackoverflow.com/questions/28370646/how-do-i-fork-a-go-process/28371586#28371586

简要翻译一下：
- fork 早出现在只有进程，没有线程的年代
- C中是自行控制线程，这样fork之后才不会发生紊乱。一般都是单线程fork之后，才会开始多线程执行。
- Go中多线程是runtime自行决定的，所以Go中没有提供单纯的fork，而是fork之后立即就exec执行新的二进制文件

> C fork结果: 父进程fork后返回子进程的pid, 子进程在fork成功后返回0.

### channel
####  close 安全
单写多读 : 确保通道写安全的最好方式是由负责写通道的协程(生产者)自己来关闭通道，读通道的协程不要去关闭通道
多写单读 : 使用到内置 sync 包提供的 WaitGroup 对象，它使用计数来等待指定事件完成
#### 多路通道
select : 同时管理多个通道读写
select default : 非阻塞读写, 它是决定通道读写操作阻塞与否的关键

### 锁
#### RWMutex
写锁是拍他锁，加写锁时会阻塞其它协程再加读锁和写锁，读锁是共享锁，加读锁还可以允许其它协程再加读锁，但是会阻塞加写锁

读写锁在写并发高的情况下性能退化为普通的互斥锁

### database/sql
当数据库字段内容为`NULL`时, golang driver不会调用相应字段的`Scan()`