## package 和func 定义使用变量的区别

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

## 权威导入路径(import paths)

我们经常使用托管在公共代码托管服务中的代码，诸如`github.com`，这意味着包导入路径包含托管服务名，比如`github.com/rsc/pdf`.一些场景下为了不破坏用户代码，我们用`rsc.io/pdf`，屏蔽底层具体哪家托管服务，比如`rso.io/pdf`的背后可能是`github.com`也可能是`bitbucket`.但这样会引入一个问题，那就是不经意间我们为一个包生成了两个合法的导入路径.如果一个程序中使用了这两个合法路径，一旦某个路径没有被识别出有更新，或者将包迁移到另外一个不同的托管公共服务下去时，使用旧导入路径包的程序就会报错.

Go 1.4引入一个规范包导入的注释，用于标识这个包的权威导入路径.如果使用的导入的路径不是权威路径，go命令会拒绝编译.
语法很简单,给非权威导入路径的package的任一go文件添加权威导入路径即可.这里是给github.com/rsc/pdf的任一go文件添加:

```go
package pdf // import "rsc.io/pdf"
```

那么那些尝试使用github.com/rsc/pdf导入路径的程序将会被go编译器拒绝编译.

这个权威导入路径检查是在编译期进行的，而不是下载阶段.
