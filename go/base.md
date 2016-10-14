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
	s = append(s, sc...)
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

### defer

**defer在函数返回前执行，即可以修改函数返回值**
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
- [defer表达式中变量的值在defer表达式被定义时就已经明确](http://www.xiaozhou.net/something-about-defer-2014-05-25.html),即被推迟函数的实参在defer执行时就会求值， 而不是在调用执行时才求值.这样不仅无需担心变量值在函数执行时被改变， 同时还意味着单个已推迟的调用可推迟多个函数的执行.
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
- [Golang中defer、return、返回值之间执行顺序](https://xiequan.info/golang%E4%B8%ADdefer%E3%80%81return%E3%80%81%E8%BF%94%E5%9B%9E%E5%80%BC%E4%B9%8B%E9%97%B4%E6%89%A7%E8%A1%8C%E9%A1%BA%E5%BA%8F/),即return最先给返回值赋值；接着defer开始执行一些收尾工作；最后return指令携带返回值退出函数.
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

init函数会被自动调用，在main函数之前执行，且不能在其他函数中调用(会报错该函数未定义)

所有init函数都会被自动调用，调用顺序如下：

- 同一个go文件的init函数调用顺序是 从上到下的
- 同一个package中按go源文件名字符串比较 从小到大顺序调用各文件中的init函数
- 不同的package，如果不相互依赖的，按照main包中 先import的后调用的顺序调用其包中的init函数
- 如果package存在依赖，则先调用最早被依赖的package中的init函数

### 闭包

闭包是由函数及其相关引用环境组合而成的实体,即闭包=函数+引用环境.

闭包只是在形式和表现上像函数，但实际上不是函数。函数是一些可执行的代码，这些代码在函数被定义后就确定了，不会在执行时发生变化，所以一个函数只有一个实例。**闭包在运行时可以有多个实例**，不同的引用环境和相同的函数组合可以产生不同的实例。所谓引用环境是指在程序执行中的某个点所有处于活跃状态的约束所组成的集合。其中的约束是指一个变量的名字和其所代表的对象之间的联系。

闭包函数出现的条件：

1. **被嵌套的函数引用到非本函数的外部变量**，而且这外部变量不是“全局变量”
1. 嵌套的函数被独立了出来(被父函数返回或赋值 变成了独立的个体)，而被引用的变量所在的父函数已结束

```go
//fn[i]为闭包函数,for循环中定义的局部变量i对于此闭包而言就等于外部的变量，内部可以访问外部的变量
//第二个for循环打印时,第一个for循环定义的那个局部变量i已经递增至3
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
//fn[i]为闭包函数,第一个for循环中定义的变量mian.i对于此闭包而言就等于外部的变量
//第二个for循环打印时,mian.i是递增的.
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

## 其他

```go
func main() {
    var a int =1
    var v interface{} =&a

    t:=4
    p:=v.(*int)
    p=&t
    //_=p

    //v.(*int)=&t
    fmt.Println(a)
}
```

疑问:

1. `_=p`注释掉时报`p declared and not used`, p不是在`p=&t`中使用了吗?
1. `v.(*int)=&t`为什么报错?

解答:

1. 这里只是给 p 赋了值，不算使用，真正的算使用要在 = 右侧，作为参数传入方法里面,或`*`取值操作即`*p`.
2. `v.(*int)`这是类型断言，不能赋值,即无变量可用于操作(也就是没有变量可存储&t的值),可更正为`*v.(*int)=t`,此时`t`的值将存入变量`a`.

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
	Name  int
	Slice []int
	Map   map[string]int
}

func main() {
	t := new(T)
	t.Name = 12
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
```go
//<<Go语言编程>>.许式伟.p74
type Integer int

func (a Integer) Less(b Integer) bool {
	return a < b
}

func (a *Integer) Add(b Integer) {
	*a += b
}

type LessAdder interface {
	Less(b Integer) bool
	Add(b Integer)
}

func main() {
	var a Integer = 1
	var b LessAdder = &a //正确,编译器自动生成会func (a *Integer) Less(b Integer) bool，其实直接按照接口赋值规则理解更方便
	var b LessAdder = a //错误,go无法自动生成func (a Integer) Add(b Integer),原因如下
}
```
```go
type Integer int
func (a *Integer) Add(b Integer) {
    *a += b
}

func (a Integer) Bdd(b Integer) {
    (&a).Add(b)
}

type LessAdder interface {
    Add(b Integer)
}
func main() {
    var a,b,c Integer = 1,1,2
    p:=&a

	p.Add(c)
	fmt.Println(a)

	b.Bdd(c)
	fmt.Println(b) //b=1，和实际期望3不符。
}
```

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
