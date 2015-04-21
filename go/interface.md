### interface何时为nil

```golang
package main

import "fmt"

type Stringer interface {
	String() string
}
type String struct {
	data string
}

func (s *String) String() string {
	return s.data
}

func GetString() *String {
	return nil
}

func CheckString(s Stringer) bool {
	//interface底层实际上是一个结构体，包括两个成员，一个是指向数据的指针，一个包含了成员的类型信息,两者同时为nil时接口才是nil
	//此处接口s里存储的类型明显不为nil,而是*String
	return s == nil
}

func main() {
	fmt.Printf("%#v,%t\n", GetString(), GetString() == nil)
	println(CheckString(GetString()))
	var a interface{}

    //nil可看成是<nil>(nil)(空类型的空指针)
    //nil只能赋值给pointer,channel,func,interface,map,slice类型的变量,详见http://pkg.golang.org/pkg/builtin/#Type
	fmt.Printf("%#v,%t\n", a, a == nil)
}

// Output:
// (*main.String)(nil),true
// false
// <nil>,true

```
