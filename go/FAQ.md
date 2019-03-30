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