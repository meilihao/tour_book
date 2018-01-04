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

### goto
如果必须使用 goto，应当**只使用正序的标签**（标签位于 goto 语句之后），且**标签和 goto 语句之间不能定义新变量**，否则会导致编译失败.