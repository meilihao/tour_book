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

## interface{}

### 接口断言
```
package main

import "fmt"

func main() {
	var t T1 = 10
	var i1 I1 = t
	i1.Speak()

	var i2 I2
	i2 = i1
	i2.Speak()

	i1 = i2.(I1) // i2接口里的实际值是否实现了(I1)接口
	i1.Step()
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
