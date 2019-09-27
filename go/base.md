## Go语言回顾：从Go 1.0到Go 1.13
- 1.0,  2012.3  : Go语言的第一个版本，并发布了一份兼容性文档, 该文档保证未来的Go版本将保持向后兼容
- 1.1,  2013.5  : Goroutine调度器被重写了，重写后的调度器性能大幅提升 & 嵌入了一个竞态探测器(race detector)
- 1.3,  2014.6  : 栈内存分配采用连续段(contiguous segment)的分配模式以提升内存分配效率. 之前的分割栈分配方式(segment stack)存在频繁分配/释放栈段导致栈内存分配性能不稳定(较低)的问题.
- 1.5,  2015.8  : 完成自举 & 垃圾回收器被全面重构, 由于引入并发回收器，回收阶段带来的延迟大幅减少
- 1.8,  2017.2  : 延迟时间已经全面降到毫秒级别以下
- 1.11, 2018.8  : 引入Go modules
- 1.13, 2019.8  : 重新实现逃逸分析(escape analysis)使得更少地在堆上分配内存。

## 类型
```go
s := "众人a"
for i, v := range s {
	fmt.Printf("%d,%x\n", i, v)
}
for i, v := range []byte(s) {
	fmt.Printf("%d,%x\n", i, v)
}

// output:
// 0,4f17
// 3,4eba
// 6,61
// 0,e4
// 1,bc
// 2,97
// 3,e4
// 4,ba
// 5,ba
// 6,61
```

"\u4f17"="\U00004f17"="\xe4\xbc\x97":
- "\u4f17"和"\U00004f17"是unicode point,"\u"表示的rune值的高位2字节为0, `\U`是完整格式.
- "\xe4\xbc\x97": utf8编码

### 传值和传指针

**go中所有传递均是值拷贝**

```go
package main

import (
	"fmt"
	"unsafe"
)

type SocialClient struct {
	Name string
}

type SocialConf []SocialClient

func (s SocialConf) New() {
	s = make([]SocialClient, 0)
}

func (s SocialConf) Add(sc ...SocialClient) { //等价于 func Add(s SocialConf,sc ...SocialClient){}
	fmt.Printf("2: %p\n", &s)
	s = append(s, sc...)  // 此处会分配底层数组
	p := (*struct {
		array uintptr
		len   int
		cap   int
	})(unsafe.Pointer(&s)) //获取slice底层结构的内容

	fmt.Printf("2: output: %+v\n", p)
	fmt.Printf("2: %#v\n", s)
}

func main() {
	var s SocialConf
	fmt.Printf("0: %p\n", &s)
	s.New() //与下面的s.Add()的传值情况相同
	var sc1 = SocialClient{Name: "weibo"}
	var sc2 = SocialClient{Name: "qq"}
	s.Add(sc1, sc2) //相当于将s的底层结构拷贝一份传入"func (s SocialConf) Add()",拷贝内容变化,但原值s的内容不变
	p := (*struct {
		array uintptr
		len   int
		cap   int
	})(unsafe.Pointer(&s))
	fmt.Printf("0: output: %+v\n", p)
	fmt.Printf("0: %#v\n", s)
}
/*
0: 0xc20801e020
2: 0xc20801e040
2: output: &{array:833357996128 len:2 cap:2}
2: main.SocialConf{main.SocialClient{Name:"weibo"}, main.SocialClient{Name:"qq"}}
0: output: &{array:0 len:0 cap:0}
0: main.SocialConf(nil)
*/
```

```go
package main

import (
	"fmt"
)

type SocialClient struct {
	Name string
}

type SocialConf []SocialClient

func (s *SocialConf) New() {
	fmt.Printf("%#v\n", s)
	*s = make([]SocialClient, 0)
}

func (s *SocialConf) Add(sc ...SocialClient) {
	*s = append(*s, sc...)
}

func main() {
	var s *SocialConf
	s = new(SocialConf) //关键,因为原先s只是一个值为nil的指针
	s.New()
	var sc1 = SocialClient{Name: "weibo"}
	var sc2 = SocialClient{Name: "qq"}
	s.Add(sc1, sc2)
	fmt.Printf("%#v\n", *s)
}
/*
&main.SocialConf(nil)
main.SocialConf{main.SocialClient{Name:"weibo"}, main.SocialClient{Name:"qq"}}
*/
```

### 初始化

```go
package main

import (
	"fmt"
)

type T struct {
	Name string
}

func (t *T) Print() string {
	t.Name = "test"
	return t.Name
}

var t *T

func main() {
	fmt.Printf("%#v\n", t)
	//t = new(T)<=>t = &T{} //ok
	*t = T{} //error:*t此时为nil,根本没有分配内存过(即没指向变量),谈何赋值
	fmt.Printf("%#v\n", t)
	fmt.Println(t.Print())
}
```

## 指针

unsafe模块的文档中提到几条转换规则，理解了以后就很容易做指针运算了：

A pointer value of any type can be converted to a Pointer.
A Pointer can be converted to a pointer value of any type.
A uintptr can be converted to a Pointer.
A Pointer can be converted to a uintptr.

unsafe.Pointer类似C的void *,即通用指针类型,可表示任意类型值的指针,用于转换不同类型指针，但它不可以参与指针运算;uintptr是golang的内置类型，是能存储指针的uint整型,是变量的首地址,可用于指针运算.

```go
func main() {
	d := struct {
		s string
		x int
	}{"abc", 100}

	p := uintptr(unsafe.Pointer(&d)) // *struct -> Pointer -> uintptr

	p += unsafe.Offsetof(d.x) // uintptr + offset
	p2 := unsafe.Pointer(p)   // uintptr -> Pointer

	px := (*int)(p2) // Pointer -> *int
	*px = 200        // d.x = 200
	fmt.Printf("%#v\n", d)
}
```

### defer

**defer在函数返回前执行;注册多个 defer,按 FILO 次序执行,哪怕函数或某个延迟调用用发生生错误,这些调用用依旧
会被执行**
```go
func main() {
	fmt.Println(f()) // 返回： 15
}

func f() (i int) {
	defer func() {
		i *= 5
	}()
	return 3
}
```
```go
func main() {
	fmt.Println("beg ------------")
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}
	fmt.Println("end ------------")
}
/*
beg ------------
end ------------
4 3 2 1 0 
*/
```

- [defer表达式中变量的值在defer表达式被定义时就已经明确](http://www.xiaozhou.net/something-about-defer-2014-05-25.html),即被推迟函数的实参在defer注册时就会求值， 而不是在调用执行时才求值,这样不仅无需担心变量值在函数执行时被改变.
```
func a() {
    i := 0
    defer fmt.Println(i) //i==0
    i++
    return
}
```
```
for i := 0; i < 5; i++ {
	defer fmt.Printf("%d ", i) //4 3 2 1 0
}

for i := 0; i < 5; i++ {
	defer func() {
		fmt.Printf("%d ", i) // 闭包,引用 // 5 5 5 5 5
	}()
}
```
```
func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}

func b() {
	defer un(trace("b"))
	fmt.Println("in b")
	a()
}

func main() {
	b()
}
//---output:
entering: b
in b
entering: a
in a
leaving: a
leaving: b
```
- [*Golang中defer、return、返回值之间执行顺序](https://xiequan.info/golang%E4%B8%ADdefer%E3%80%81return%E3%80%81%E8%BF%94%E5%9B%9E%E5%80%BC%E4%B9%8B%E9%97%B4%E6%89%A7%E8%A1%8C%E9%A1%BA%E5%BA%8F/),即return最先给返回值赋值；接着defer开始执行一些收尾工作；最后return指令携带返回值退出函数.
```
func main() {
	fmt.Println("a return:", a()) // 打印结果为 a return: 0
}

func a() int {
	var i int
	defer func() {
		i++
		fmt.Println("a defer2:", i) // 打印结果为 a defer2: 2
	}()
	defer func() {
		i++
		fmt.Println("a defer1:", i) // 打印结果为 a defer1: 1
	}()
	return i
}
/*
a()int 函数的返回值没有被提前声名，其值来自于其他变量的赋值，而defer中修改的也是其他变量（其实该defer根本无法直接访问到返回值），因此函数退出时返回值并没有被修改。
其实此时的`return i`可以理解为"x:=i;return x;//x是隐藏的实际返回值",后续的defer只影响了`i`,而非`x`.
*/
```

### init函数

- 每个源文件都可以定义一个或多个init
- 编译器不保证多个init**执行行次序**
- init在单一线程被调用,仅执行一次
- init在包所有全局变量初始化后执行
- 在所有init执行后才执行 main.main
- 无法不能在其他函数中调用init

### 闭包

闭包是由函数及其相关引用环境组合而成的实体,即闭包=函数+引用环境.

闭包只是在形式和表现上像函数，但实际上不是函数。函数是一些可执行的代码，这些代码在函数被定义后就确定了，不会在执行时发生变化，所以一个函数只有一个实例。**闭包在运行时可以有多个实例**，不同的引用环境和相同的函数组合可以产生不同的实例。所谓引用环境是指在程序执行中的某个点所有处于活跃状态的约束所组成的集合。其中的约束是指一个变量的名字和其所代表的对象之间的联系。

闭包函数出现的条件：

1. **被嵌套的函数引用到非本函数的外部变量**，而且这外部变量不是“全局变量”
1. 嵌套的函数被独立了出来(被父函数返回或赋值 变成了独立的个体)，而被引用的变量所在的父函数已结束

```go
//fn[i]为闭包函数,for循环中定义的局部变量i对于此闭包而言就等于外部的变量，内部可以访问外部的变量
//第二个for循环打印时,第一个for循环定义的那个局部变量i已经递增至3且已不再改变
func main() {
	var fn [3]func()

	for i := 0; i < len(fn); i++ {
		fn[i] = func() {
			fmt.Println(i)
		}
	}

	for _, f := range fn {
		f()
	}
}
/*---
3
3
3
---*/
//fn[i]为闭包函数,第一个for循环中使用的变量mian.i对于此闭包而言就等于外部的变量
//第二个for循环打印时,mian.i被重置,再次开始递增.
//如果将第二个for中的i改为局部变量的话,将打印3个3,原因是闭包函数引用的是mian.i,而不是for里的局部变量i,此时mian.i已是3.
func main() {
	var fn [3]func()
	var i int

	for i = 0; i < len(fn); i++ {
		fn[i] = func() {
			fmt.Println(i)
		}
	}

	for i = 0; i < len(fn); i++ {
		fn[i]()
	}
}
/*---
0
1
2
---*/
//fn[i]为闭包函数,查看make_fn()的参数定义,属于值传递,举例fn[0]=func(0){fmt.Println(0)}
func main() {
	var fn [3]func(int)

	for i := 0; i < len(fn); i++ {
		fn[i] = make_fn()
	}

	for i, f := range fn {
		f(i)
	}
}

func make_fn() func(i int) {
	return func(i int) {
		fmt.Println(i)
	}
}
/*---
0
1
2
---*/
```

```go
func main() {
	// 
	for i := 0; i < 3; i++ {
		defer fmt.Println(i)                    // 2,1,0
		defer func() { fmt.Println(i) }()       // 3,3,3
		defer func(i int) { fmt.Println(i) }(i) // 2,1,0
		defer print(&i)                         // 3,3,3
		go fmt.Println(i)                       // 0,1,2 # 顺序未知
		go func() { fmt.Println(i) }()          // 3,3,3 # 顺序未知
	}
}

func print(pi *int) { fmt.Println(*pi) }
```

```go
package main

import (
	"fmt"
)

type A struct {
	Name string
}

func (a *A) GetName() {
	fmt.Println(a.Name)
	fmt.Printf("%p\n", a)
}

func GetName() {
	fmt.Println("NoStruct")
}

type B struct{}

func (b *B) Call(fn func()) {
	fn() // 匿名函数会保留现场
}

func (b *B) Call2(fn func(*A), a *A) {
	fn(a)
}

func main() {
	a := &A{
		Name: "sdf",
	}
	b := new(B)

	fmt.Printf("%T\n", a.GetName)     // method value
	fmt.Printf("%T\n", (*A).GetName)  // method expresion
	b.Call(GetName)
	b.Call(a.GetName)
	b.Call2((*A).GetName, a)

	/*output:
	func()
	func(*main.A)
	NoStruct
	sdf
	0x1040c108
	sdf
	0x1040c108
	*/
}
```

## 常量

常量只有在最终被赋值给一个变量的时候才可以会出现溢出的情况,因此下面的语句是合法的:

    const t=(1<<100)>>97 //t=8

### 转换

#### unsafe.Pointer转Slice

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type T struct {
	Int  int
	Slice []int
	Map   map[string]int
}

func main() {
	t := new(T)
	t.Int = 12
	t.Slice = make([]int, 3, 6)
	t.Map = map[string]int{"A": 1, "B": 2}
	fmt.Println(t)

	fmt.Printf("%p\n", &t.Slice)

	k := unsafe.Pointer(uintptr(unsafe.Pointer(t)) + uintptr(unsafe.Sizeof(int(0))))
	fmt.Printf("%#v\n", k)

	p := (*struct {
		array uintptr
		len   int
		cap   int
	})(k)

	fmt.Printf("output: %+v\n", p) //检查k的地址是否正确(看len)
	fmt.Printf("%x\n", p.array)

	s := *(*[]int)(k) //Pointer转slice
	fmt.Println(s)

	//---其他获取slice底层结构的方式

	hdr := (*reflect.SliceHeader)(k)

	// Create a header for comparison via pointer manipulation
	address := (uintptr)(k)
	lenAddr := address + unsafe.Sizeof(address)
	capAddr := lenAddr + unsafe.Sizeof(int(0))

	unsafeHdr := reflect.SliceHeader{
		Data: *(*uintptr)(unsafe.Pointer(address)),
		Len:  *(*int)(unsafe.Pointer(lenAddr)),
		Cap:  *(*int)(unsafe.Pointer(capAddr)),
	}

	fmt.Printf("Real Header:  %#v\n", *hdr)
	fmt.Printf("UnsafeHeader: %#v\n", unsafeHdr)
}
```
```go
package main

import (
	"fmt"
	"unsafe"
)

type (
	Router struct {
		Level int
		*Route
	}

	Route struct {
		Path string
	}
)

func main() {
	route := new(Route)
	route.Path = "/hello"
	r := Router{Level: 1, Route: route}

	t := new(int)
	p := unsafe.Pointer(uintptr(unsafe.Pointer(&r)) + uintptr(unsafe.Sizeof(*t)))

	m := *(**Route)(p) //p是指向*Route的指针
	fmt.Printf("%#v\n", m)

	m0 := *(*Route)(m)
	fmt.Printf("%#v\n", m0)

	p1 := unsafe.Pointer(uintptr(unsafe.Pointer(&r)) + 0)//将实际地址转换为指向该地址的指针
	m1 := *(*int)(p1) //p1是指向Router.Level的指针
	fmt.Println(m1)
}
```

### 方法

每个类型都有与之关联的方法集,这会影响到接口实现规则:
- 类型 T 方法集包含全部 receiver T 方法
- 类型 *T 方法集包含全部 receiver T + *T 方法
- 如类型 S 包含匿名字段 T,则 S 方法集包含 T 方法
- 如类型 S 包含匿名字段 *T,则 S 方法集包含 T + *T 方法
- 不管嵌入 T 或 *T,*S 方法集总是包含 T + *T 方法

用实例 value 和 pointer 调用方法 (含匿名字段) 不受方法集约束,编译器总是查找全部
方法,并自动转换 receiver 实参.


```go
package main

import (
	"fmt"
)

type Tc interface {
	V()
}
type myType struct {
}

func (m *myType) V() {
	fmt.Println("V()")
}
func main() {
	//var i Tc = myType{}
	// i.V() //报错,类型 T 的可调用方法集不包含接受者为 *T 的方法,参考:http://se77en.cc/2014/05/05/methods-interfaces-and-embedded-types-in-golang/

	// 而这里不报错
	var w myType = myType{}
	w.V()//=>(&w).V(),编译器自动推导
}
```

### 作用域

```go
func main() {
	DoTheThing()
}

func DoTheThing() (err error) {
	fmt.Printf("%p", &err)

	if true {
		// 两个err作用域不同,会重新定义
		result, err := tryTheThing()
		fmt.Printf("%p", &err)
		if err != nil {
			err = errors.New("b")
		}

		fmt.Println(result)
	}

	return err
}

func tryTheThing() (string, error) {
	return "", errors.New("b")
}
```

```go
func main() {
	DoTheThing()
}

func DoTheThing() (err error) {
	fmt.Printf("%p", &err)

	// 相同作用域,已定义过的不再重新定义
	result, err := tryTheThing()
	fmt.Printf("%p", &err)
	if err != nil {
		err = errors.New("b")
	}

	fmt.Println(result)

	return err
}

func tryTheThing() (string, error) {
	return "", errors.New("b")
}
```

## 其他

```go
func main() {
    var a int =1
    var v interface{} =&a

    t:=4
    p:=v.(*int)
    p=&t
    //_=p

    //v.(*int)=&t // cannot assign to v.(*int)
    fmt.Println(a)
}
```

疑问:

1. `_=p`注释掉时报`p declared and not used`, p不是在`p=&t`中使用了吗?
1. `v.(*int)=&t`为什么报错?

解答:

1. 这里只是给 p 赋了值，不算使用，真正的算使用要在 = 右侧，作为参数传入方法里面,或`*`取值操作即`*p`.
2. `v.(*int)`这是类型断言，不能赋值,即无变量可用于操作(也就是没有变量可存储&t的值),可更正为`*v.(*int)=t`,此时`t`的值将存入变量`a`.


```go
// 接口和指针
package main

import (
	"fmt"
	"reflect"
)

type MyErr struct{}

func (m MyErr) Error() string {
	return "123"
}

func main() {
	v := new(error) //已分配接口存储位置
	*v = MyErr{}
	fmt.Printf("%#v\n", v)
	fmt.Println(reflect.TypeOf(v), (*v).Error())

	//x := (*error)(nil) //未分配接口存储位置
	//*x =&MyErr{} //因为没有存储位置而报错

	z := (*error)(nil)
	var e error
	e = MyErr{}
	z = &e
	fmt.Printf("%#v\n", z)
	fmt.Println(reflect.TypeOf(z), (*z).Error())
}
```
