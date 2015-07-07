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

**defer在函数返回后执行，可以修改函数返回值**

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