# cgo

参考：
- [CGO编程](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/readme.html)
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
- import "C" 这句话要紧随，注释后，不要换行，否则报错. 其表示启用CGO特性, 而go build命令会在编译和链接阶段启动gcc编译器处理.
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

## FAQ
### so依赖
如果引入的so依赖其他so,那么先使用`ldd libai.so`查看并安装缺失的so(注意版本),否则`go build`时会报错,比如:
```
/usr/bin/x86_64-linux-gnu-ld: warning: libsrtp.so.2, needed by ./libai.so, not found (try using -rpath or -rpath-link)
...
/usr/bin/x86_64-linux-gnu-ld: warning: libzlog.so.1.2, needed by ./libai.so, not found (try using -rpath or -rpath-link)
./libai.so：对‘virtual_factory_set_recv_send_circbuf_cb’未定义的引用
./libai.so：对‘zlog_put_mdc’未定义的引用
...
./libai.so：对‘zlog’未定义的引用
./libai.so：对‘virtual_factory_set_record_op’未定义的引用
collect2: error: ld returned 1 exit status
```
当然最可靠和方便的方法是在构建引入so的电脑上运行`go build`,省事又省力.

### could not determine kind of name for C.CString
```
$ go build
../../git/go/src/pmanage/manager.go:234:12: could not determine kind of name for C.CString
../../git/go/src/pmanage/manager.go:407:3: could not determine kind of name for C.ai_process
../../git/go/src/pmanage/manager.go:125:10: could not determine kind of name for C.libai_init
../../git/go/src/pmanage/manager.go:241:2: could not determine kind of name for C.sync_process
../../git/go/src/pmanage/manager.go:237:2: could not determine kind of name for C.update_call_status
```

引入的"xxx.h"的某些定义函数缺少结尾的`;`.

### could not determine kind of name for C.free
添加头文件: `#include <stdlib.h>`