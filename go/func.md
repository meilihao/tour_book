## GoLang之各种函数的用法

### 大纲
本文总结GoLang中常用的几种函数用法，主要包括：

[0] 首先main是一个没有返回值的函数
[1] 普通函数
[2] 函数返回多个值
[3] 不定参函数
[4] 闭包函数
[5] 递归函数
[6] 类型方法, 类似C++中类的成员函数
[7] 接口和多态
[9] 错误处理, Defer接口
[10] 错误处理, Panic/Recover

### 测试代码

```go
package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func multi_ret(key string) (int, bool) {
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	var err bool
	var val int
	val, err = m[key]
	return val, err
}

func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func nextNum() func() int {
	i, j := 0, 1
	return func() int {
		var tmp = i + j
		i, j = j, tmp
		return tmp
	}
}

func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

// 长方形
type rect struct {
	width, height float64
}

func (r *rect) area() float64 {
	return r.width * r.height
}

func (r *rect) perimeter() float64 {
	return 2 * (r.width + r.height)
}

// 圆形
type circle struct {
	radius float64
}

func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c *circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// 接口
type shape interface {
	area() float64
	perimeter() float64
}

func interface_test() {
	r := rect{width: 2, height: 4}
	c := circle{radius: 4.3}

	// 通过指针实现
	s := []shape{&r, &c}

	for _, sh := range s {
		fmt.Println(sh)
		fmt.Println(sh.area())
		fmt.Println(sh.perimeter())
	}
}

type myError struct {
	arg    int
	errMsg string
}

// 实现error的Error()接口
func (e *myError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.errMsg)
}

func error_test(arg int) (int, error) {
	if arg < 0 {
		return -1, errors.New("Bad Arguments, negtive")
	} else if arg > 256 {
		return -1, &myError{arg, "Bad Arguments, too large"}
	}
	return arg * arg, nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("Open failed")
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		fmt.Println("Create failed")
		return
	}
	defer dst.Close()

	// 注意dst在前面
	return io.Copy(dst, src)
}

//////////////////////
func main() {

	// [0] 首先main是一个没有返回值的函数

	// [1] 普通函数
	fmt.Println(max(1, 100))

	// [2] 函数返回多个值
	v, e := multi_ret("one")
	fmt.Println(v, e)
	v, e = multi_ret("four")
	fmt.Println(v, e)
	// 典型的判断方法
	if v, e = multi_ret("five"); e {
		fmt.Println("OK")
	} else {
		fmt.Println("Error")
	}

	// [3] 不定参函数
	sum(1, 2)
	sum(2, 4, 5)
	nums := []int{1, 2, 3, 4, 5}
	sum(nums...)

	// [4] 闭包函数
	nextNumFunc := nextNum()
	for i := 0; i < 10; i++ {
		fmt.Println(nextNumFunc())
	}

	// [5] 递归函数
	fmt.Println(fact(4))

	// [6] 类型方法, 类似C++中类的成员函数
	r := rect{width: 10, height: 15}
	fmt.Println("area: ", r.area())
	fmt.Println("perimeter: ", r.perimeter())

	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perimeter: ", rp.perimeter())

	// [7] 接口和多态
	interface_test()

	// [8] 错误处理, Error接口
	for _, val := range []int{-1, 4, 1000} {
		if r, e := error_test(val); e != nil {
			fmt.Printf("failed: %d:%s\n", r, e)
		} else {
			fmt.Println("success: ", r, e)
		}
	}

	// [9] 错误处理, Defer接口
	if w, err := CopyFile("/data/home/gerryyang/dst_data.tmp", "/data/home/gerryyang/src_data.tmp"); err != nil {
		fmt.Println("CopyFile failed: ", e)
	} else {
		fmt.Println("CopyFile success: ", w)
	}

	// 你猜下面会打印什么内容
	fmt.Println("beg ------------")
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}
	fmt.Println("end ------------")

	// [10] 错误处理, Panic/Recover
	// 可参考相关资料, 此处省略

}

/*
output:

100
1 true
0 false
Error
[1 2] 3
[2 4 5] 11
[1 2 3 4 5] 15
1
2
3
5
8
13
21
34
55
89
24
area:  150
perimeter:  50
area:  150
perimeter:  50
&{2 4}
8
12
&{4.3}
58.088048164875275
27.01769682087222
failed: -1:Bad Arguments, negtive
success:  16 <nil>
failed: -1:1000 - Bad Arguments, too large
CopyFile success:  8
------------
------------
4 3 2 1 0
*/
```

## 函数作为变量

在go中函数也是一种变量，我们通过type定义这种变量的类型。拥有相同参数和相同返回值的函数属于同一种类型。

```go
package main

import (
	"fmt"
)

//声明函数类型 t1,t2
type t1 func(*int)
type t2 func(*int) int

//t的类型为t1
func method1(i int, t t1) {
	t(&i)
}

//声明func(*int)类型的函数
func test1(m *int) {
	fmt.Println(*m)
}

//t的类型为t2
func method2(i int, t t2) {
	//t的参数是一个指针类型的int,返回一个int类型
	r := t(&i)
	fmt.Println(r)
}

//声明(func(*int) int)类型的函数
func test2(n *int) int {
	return *n
}

func main() {
	//调用 test1为类型为t1的函数,函数当做值来传递
	method1(100, test1)
	//调用,test2为类型为t2的函数
	method2(10, test2)
	
	x:=t1(func(i *int){
		fmt.Println("this is ",*i)
	})
               
	y:=func(i *int){
		fmt.Println("this is ",*i)
	}
	
	method1(100,x)
	method1(101,y)
}
```
