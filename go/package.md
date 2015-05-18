###package 和func 定义使用变量的区别

package 级别可先使用自定义类型而后再定义其类型(比如下方的myint,即类似于定义和使用是分别处于同package下的不同源码文件中),func级别必须先定义再使用
```go
package main

import "fmt"

type Result interface {
	test()
}

var x myint

type myint int

// "var _ XXX"形式是(从运行的角度来说没用，从编译的角度来说是)用于判断类型是否已定义或者是否已实现接口
var _ Result = myint(0)

func (myint) test() {
	fmt.Println("test")
}

func main() {
	var y myint
	y.test()
	fmt.Println(y)
}
```
