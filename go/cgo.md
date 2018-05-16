# cgo

参考：
- [binnnliu/cgo-demo](https://github.com/binnnliu/cgo-demo)
- [draffensperger/go-interlang](https://github.com/draffensperger/go-interlang)
- [Go与C语言的互操作](https://tonybai.com/2012/09/26/interoperability-between-go-and-c/)
- [全面总结： Golang 调用 C/C++，例子式教程 - 方法2有问题：无法找到so](https://juejin.im/post/5a62f7cff265da3e4c07e0ab)
- [官方wiki/cgo](https://github.com/golang/go/wiki/cgo)

## 直接嵌入go
```
.
├── demo.c
└── demo.go

0 directories, 2 files
```

demo.c
```c
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

char* Demo(char* name) {
	int size = strlen("Hello ") + strlen(name) + 1;
	char* buf = (char *)malloc(size);
	memset(buf, 0, size);

	sprintf(buf, "Hello %s", name);
	return buf;
}

void main() {
	char *str = Demo("World");
	printf("%s\n", str); // printf第一个参数必须带‘\n’,否则无法输出字符串
	free(str);
}
```

demo.go
```go
package main

//#include <stdio.h>
//#include <string.h>
//#include <stdlib.h>
//
//
//
//char* Demo(char* name) {
//	int size = strlen("Hello ") + strlen(name) + 1;
//	char* buf = (char *)malloc(size);
//	memset(buf, 0, size);
//
//	sprintf(buf, "Hello %s", name);
//	return buf;
//}
import "C"
import "fmt"
import "unsafe"

func main() {
	name := C.CString("World")
	defer C.free(unsafe.Pointer(name))
	
	ret := C.Demo(name)  // name已是指针
	gret := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	fmt.Println(gret)
}
```

直接运行`go build`即可.

要点：
- 但凡要引用与 c/c++ 相关的内容，写到 go 文件的头部注释里面
- 嵌套的 c/c++ 代码必须符合其语法，不与 go 一样
- import "C" 这句话要紧随，注释后，不要换行，否则报错
- go 代码中调用 c/c++ 的格式是: C.xxx()，例如 C.Demo(name)

## 直接引用 c/c++ 文件
```
.
├── demo.c
├── demo.go
└── demo.h

0 directories, 3 files
```

demo.h
```c
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

char* Demo(char* name);
```

demo.c
```c
#include "demo.h"


char* Demo(char* name) {
	int size = strlen("Hello ") + strlen(name) + 1;
	char* buf = (char *)malloc(size);
	memset(buf, 0, size);

	sprintf(buf, "Hello %s", name);
	return buf;
}
```

demo.go
```go
package main

//#include <demo.h>
import "C"
import "fmt"
import "unsafe"


func main() {
	name := C.CString("World")
	defer C.free(unsafe.Pointer(name))
	
	ret := C.Demo(name)
	gret := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	fmt.Println(gret)
}
```

## 使用动态库so
沿用上面的demo.h,demo.c, 并用demo.c制作so.

```
$ gcc -fPIC -shared -o libdemo.so demo.c
```

demo.go
```go
package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -ldemo -Wl,-rpath=./
#include "demo.h"
*/
import "C"
import "fmt"
import "unsafe"

func main() {
	name := C.CString("World")
	defer C.free(unsafe.Pointer(name))
	
	ret := C.Demo(name)
	gret := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	fmt.Println(gret)
}
```

要点：
- CFLAGS: -I路径 这句话指明头文件所在路径，-I. 指明 当前项目根目录
- LDFLAGS: -L路径 -l名字 指明动态库的所在路径，-L. -ldemo，指明在 libdemo.so的位置，**`-l`后没有前缀`lib`**

运行`go build`正常，但执行报错`./demo: error while loading shared libraries: libdemo.so: cannot open shared object file: No such file or directory`，即`ldd demo` 报`libdemo.so => not found`，解决方法有两：
- 设置rpath， 编译时直接指定so的位置， `cgo LDFLAGS: -L. -ldemo -Wl,-rpath=./`, `./`在这里表示当前目录，直接用`.`会报错.
- 设置LD_LIBRARY_PATH，运行`LD_LIBRARY_PATH=. ./demo`