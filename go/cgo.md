# cgo

参考：
- [CGO编程](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/readme.html)
- [binnnliu/cgo-demo](https://github.com/binnnliu/cgo-demo)
- [draffensperger/go-interlang](https://github.com/draffensperger/go-interlang)
- [Go与C语言的互操作](https://tonybai.com/2012/09/26/interoperability-between-go-and-c/)
- [全面总结： Golang 调用 C/C++，例子式教程 - 方法2有问题：无法找到so](https://juejin.im/post/5a62f7cff265da3e4c07e0ab)
- [官方wiki/cgo](https://github.com/golang/go/wiki/cgo)

go build调用了名为cgo的工具, 它会识别和读取Go源文件中的C代码,并将其提取后交给外部的C编译器(clang或gcc)编译, 最后与Go源码编译后的目标文件链接成一个可执行程序.

内存管理:
1. 在C内部分配的内存,Go中的GC是无法感知到的,因此要记着在使用后手动释放

注意点:
1. import "C"不支持放在xx_test.go文件中

## 使用cgo的开销
1. cgo比调用Go函数多出一个甚至多个数量级
1. 增加线程数量暴涨的可能性

	C空间的函数导致执行这段代码的线程(M)挂起,这之后Go运行时调度代码只能创建新的线程以供其他没有绑定M的P上的goroutine使用. 典型场景是go代码通过cgo调用c中的sleep().
1. 失去跨平台交叉构建能力

	跨平台编译能力仅限于纯Go代码. 当Go编译器执行跨平台编译时,它会将CGO_ENABLED置为0,即关闭cgo
1. 其他开销

## cgo缺点
1. go使用gc, c手动管理内存, 两边切换容易产生bug, 以及带来很大心智负担
1. go的工具链在c前失效

	1. Go的竞态检测工具、性能剖析工具、测试覆盖率工具、模糊测试以及源码竞态分析工具等失效
	2. cgo难调试

## 在Go中使用C语言的类型
1. 原生类型

	1. 数值

		```go
		C.char
		C.schar (signed char)
		C.uchar (unsigned char)
		C.short
		C.ushort (unsigned short)
		C.int, C.uint (unsigned int)
		C.long
		C.ulong (unsigned long)
		C.longlong (long long)
		C.ulonglong (unsigned long long)
		C.float
		C.double
		```

		Go的数值类型与C中的数值类型不是一一对应的,因此在使用对方类型变量时少不了显式类型转换操作
	2. 指针

		原生数值类型的指针类型可按Go语法在类型前面加上星号`*`即可, 但`void*`比较特殊,在Go中用unsafe.Pointer表示它,这是因为任何类型的指针值都可以转换为unsafe.Pointer类型,而unsafe.Pointer类型也可以转换回任意类型的指针类型
	3. 字符串

		C语言中并不存在原生的字符串类型,在C中用带结尾'\0'的字符数组来表示字符串;而在Go中,string类型是语言的原生类型,因此这两种语言的互操作势必要进行字符串类型的转换

		`C.GoString()`相当于在Go世界重新分配一块内存对象,并复制了C的字符串的信息
	4. 数组

		C语言中的数组与Go语言中的数组差异较大,后者是原生的值类型,而前者与C中的指针在大部分场合可以随意转换。Go仅提供了C.GoBytes来将C中的char类型数组转换为Go中的[]byte切片类型

		其他类型的数组处理:
		```go
		func CArrayToGoArray(cArray unsafe.Pointer, elemSize uintptr, len int) (goArray []int32) {
			for i := 0; i < len; i++ {
				j := *(*int32)((unsafe.Pointer)(uintptr(cArray) + uintptr(i)*elemSize))
				goArray = append(goArray, j)
			}

			return
		}

		func main() {
			goArray := CArrayToGoArray(unsafe.Pointer(&C.cArray[0]), unsafe.Sizeof(C.cArray[0]), 7)
			fmt.Println(goArray)
		}
		```
2. 自定义类型

	1. 枚举
		```go
		// enum color {
		//    RED,
		//    BLUE,
		//    YELLOW
		// };
		import "C"
		import "fmt"

		func main() {
			var e, f, g C.enum_color = C.RED, C.BLUE, C.YELLOW
			fmt.Println(e, f, g) // 0 1 2
		}
		```
	2. struct

		以通过C.struct_xx来访问C中定义的结构体类型xx

		```go
		// #include <stdlib.h>
		//
		// struct employee {
		//     char *id;
		//     int  age;
		// };
		import "C"

		import (
			"fmt"
			"unsafe"
		)

		func main() {
			id := C.CString("1247")
			defer C.free(unsafe.Pointer(id))

			var p = C.struct_employee{
				id:  id,
				age: 21,
			}
			fmt.Printf("%#v\n", p)
		}
		```
	3. union

		```go
		// #include <stdio.h>
		// union bar {
		//        char   c;
		//        int    i;
		//        double d;
		// };
		import "C"
		import "fmt"

		func main() {
			var b *C.union_bar = new(C.union_bar)
			b[0] = 4
			fmt.Println(b)
		}
		```

		Go对待C的union类型与其他类型不同,Go将union类型看成[N]byte,其中N为union类型中最长字段的大小(圆整后的)

	4. 别名(typedef)

		在Go中访问C中使用typedef定义的别名类型时,其访问方式与原类型的访问方式相同

为了方便获得C世界中的类型的大小,Go提供了C.sizeof_T来获取C.T类型的大小。如果是结构体、枚举及联合体类型,我们需要在T前面
分别加上struct_、enum_和union_的前缀.

```go
// struct employee {
//     char *id;
//     int  age;
// };
import "C"

import (
	"fmt"
)

func main() {
	fmt.Printf("%#v\n", C.sizeof_int)             // 4
	fmt.Printf("%#v\n", C.sizeof_char)            // 1
	fmt.Printf("%#v\n", C.sizeof_struct_employee) // 16
}
```

## 在Go中链接外部C库
```go
package main
// #cgo CFLAGS: -I${SRCDIR}
// #cgo LDFLAGS: -L${SRCDIR} -lfoo
// #include <stdio.h>
// #include <stdlib.h>
// #include "foo.h"
import "C"
import "fmt"
func main() {
	fmt.Println(C.count)
	C.foo()
}
```

通过#cgo指示符告诉Go编译器在当前源码目录(${SRCDIR}会在编译过程中自动转换为当前源码所在目录的绝对路径)下查找头文件foo.h, 并链接当前源码目录下的libfoo共享库. C.count变量和C.foo函数的定义都在libfoo共享库中

Go支持多返回值,而C并不支持,因此当将C函数用在多返回值的Go调用中时,C的errno将作为函数返回值列表中最后那个error返回值返回.

### 使用cgo代码的静态构建
静态构建就是指构建后的应用运行所需的所有符号、指令和数据都包含在自身的二进制文件当中,没有任何对外部动态共享库的依赖。静态构建出的二进制文件由于包含所有符号、指令和数据,因而通
常要比非静态构建的应用大许多。**默认情况下,Go没有采用静态构建**.

根据Go官方文档($GOROOT/cmd/cgo/doc.go),Go链接器支持两种工作模式:内部链接(internal linking)和外部链接(external linking).

如果用户代码中仅仅使用了net、os/user等几个标准库中的依赖cgo的包,Go链接器默认使用内部链接,而无须启动外部链接器(如gcc、clang等)。不过Go链接器功能有限,仅仅将.o和预编译好的标准库的.a
写到最终二进制文件中。因此如果标准库中是在CGO_ENABLED=1的情况下编译的,那么编译出来的最终二进制文件依旧是动态链接的,即便在go build时传入-ldflags 'extldflags "-static"'也是如此,因为根本没有
用到外部链接器.

而外部链接机制则是Go链接器将所有生成的.o都写到一个.o文件中,再将其交给外部链接器(比如gcc或clang)去做最终的链接处理。如果此时在go build的命令行参数中传入-ldflags ‘extldflags “-static”’,那
么gcc/clang将会做静态链接,将.o中未定义(undefined)的符号都替换为真正的代码指令。可以通过-linkmode=external来强制Go链接器采用外部链接, go build可加`-x -v`验证.

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

> 使用 `go tool cgo -debug-gcc xxx.go` 可以得到中间代码和对象，在当前目录的 _obj 目录下

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

## cgo编译
参考:
- [完全静态编译一个Go程序](https://colobu.com/2018/07/20/totally-static-Go-builds/)

// use github.com/tecbot/gorocksdb librocksdb.a 且已安装了`apt install libc6-dev`
`CGO_LDFLAGS="-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" go build`

`CGO_CFLAGS="-I/usr/local/include" go build --ldflags '-extldflags "-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"'` // ??? 加了`-extldflags "-static"`即`-static`就报`Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking`, 不加就正常. 且此时编出的二进制仅rocksdb静态化, 其他的比如bz2, snappy还是so的形式.

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

### 返回内容需释放
```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* GetString() { // 返回 c string
	static const char* s = "0123456789";

	int len = 10;
	char* p = malloc(len);

	memcpy(p, s, len);

    return p;
}
*/
import "C"
import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
	"unsafe"

	_ "github.com/mkevac/debugcharts"
)

// 通过 top 命令的RES判断
func main() {
	go func() {
		// terminal: $ go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
		// web:
		// 1、http://localhost:8081/ui
		// 2、http://localhost:6060/debug/charts
		// 3、http://localhost:6060/debug/pprof
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	fmt.Println("pid: ", os.Getpid())

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("before, have", runtime.NumGoroutine(), "goroutines,",
		ms.Alloc, "bytes allocated", ms.HeapObjects, "heap object",
		"using mem", ms.Sys)

	for i := 0; i < 2; i++ {
		memTest()

		time.Sleep(3 * time.Second)
	}

	runtime.GC()
	fmt.Println("gc finish")
	runtime.ReadMemStats(&ms)
	fmt.Println("after gc, have", runtime.NumGoroutine(), "goroutines,",
		ms.Alloc, "bytes allocated", ms.HeapObjects, "heap object",
		"using mem", ms.Sys)

	time.Sleep(60 * time.Second)
}

func memTest() {
	fmt.Println("start")

	for i := 0; i < 10000000; i++ {
		str := C.GetString()
		_ = C.GoString(str)
		// 当cgo(GetString)调用返回后，s所占用的内存要被释放掉, 否则top's RES会持续增长导致内存泄露
		C.free(unsafe.Pointer(str))
	}

	fmt.Println("end")
}
```

### cgo 使用c struct
```go
// [Cgo中使用var声明C结构的变量是否需要释放内存？](https://segmentfault.com/q/1010000009568805)
// C.CString() 返回的 C 字符串是在堆上新创建的并且不受 GC 的管理，使用完后需要自行调用 C.free() 释放，否则会造成内存泄露，而且这种内存泄露用前文中介绍的 pprof 也定位不出来
// goroutine 通过 CGO 进入到 C 接口的执行阶段后，已经脱离了 golang 运行时的调度并且会独占线程，此时实际上变成了多线程同步的编程模型。如果 C 接口里有阻塞操作，这时候可能会导致所有线程都处于阻塞状态，其他 goroutine 没有机会得到调度，最终导致整个系统的性能大大较低。总的来说，只有在第三方库没有 golang 的实现并且实现起来成本比较高的情况下才需要考虑使用 CGO ，否则慎用.
// fmt.Println()会干扰内存泄露的排查, 需提前注释掉.
package main

/*
#include <stdlib.h>
#include <string.h>
typedef struct student
{
    int age;
    char name[1024];
}Student, *PStudent;

typedef struct student2
{
    int age;
    char* name;
}Student2, *PStudent2;

void set(Student *p) {
	p->age = 1;
	memset(p->name, 0, 1024);
	strcpy(p->name, "hello world!");

	return;
}

void set2(Student2 *p) {
	p->age = 1;
	p->name = (char*)malloc(sizeof(char)*1024);
	memset(p->name, 0, 1024);
	strcpy(p->name, "hello world!");

	return;
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
	"runtime"
)

func main() {
	n, c := 0, 0
	for {
		// s1()
		s2()

		n++

		if n%1000 == 0 {
			fmt.Println("--------c:", c)
			time.Sleep(time.Second)

			c++
		}

		// if c > 180 {
		// 	break
		// }
	}

	runtime.GC()
	
	fmt.Println("--------end:", n)

	select {}
}

func s1() {
	s := (*C.Student)(C.malloc(C.sizeof_Student))

	C.set(s)

	// fmt.Println(s.age)
	// fmt.Println(s.name)
	// fmt.Println(*(*int32)(unsafe.Pointer(s)))
	// fmt.Println(C.GoString((*C.char)(unsafe.Pointer(&((*s).name)))))

	C.free(unsafe.Pointer(s))
}

// valgrind --leak-check=full ./t
func s2() {
	s2 := (*C.Student2)(C.malloc(C.sizeof_Student2))

	C.set2(s2)

	fmt.Println(s2.age)
	fmt.Println(C.GoString(s2.name))

	C.free(unsafe.Pointer(s2.name)) // 释放c申请的内存, 否则会有memory leak
	C.free(unsafe.Pointer(s2))
}
```

#### Go 对 C 结构体的内存管理
Go 自动管理 Go 对象的内存（例如 Go 的 struct、字符串、数组等）. 但是, C 语言的内存管理由 C 语言的运行时控制，因此在 Go 中使用 cgo 时，必须手动管理 C 语言结构体的内存. 具体来说：
1. 直接声明的 C 结构体（如 C.struct_passwd）不需要手动释放内存，因为它们通常作为值在 Go 堆栈上分配。
1. 通过 malloc 分配的 C 内存（如分配 C 字符串或动态数组）需要手动释放。

C 结构体内存管理:
对于 C 结构体，在 Go 中使用 cgo 时，如果结构体的内存是在 C 端动态分配的（如使用 malloc），则需要确保在使用完毕后调用 free 来释放内存.

示例 1：直接声明 C 结构体
```go
/*
#include <pwd.h>
*/
import "C"

func main() {
    var pw C.struct_passwd
    // 使用 pw 进行操作，不需要释放内存，因为 pw 是在 Go 栈上分配的
}
```
在这种情况下，pw 变量是一个 Go 结构体类型的 C 结构体，存储在 Go 栈上，不需要手动释放.

示例 2：通过 malloc 分配 C 结构体
如果通过 malloc 或类似的 C 函数分配内存来存储 C 结构体（或数组），你必须手动释放它。例如：

```go
/*
#include <stdlib.h>
#include <pwd.h>

struct passwd* create_passwd() {
    return (struct passwd*)malloc(sizeof(struct passwd));
}

void free_passwd(struct passwd* pw) {
    free(pw);
}
*/
import "C"

func main() {
    // 动态分配 C 结构体
    pw := C.create_passwd()
    // 使用 pw 进行操作

    // 手动释放内存
    C.free_passwd(pw)
}
```
在这个示例中，C.create_passwd() 使用 malloc 分配内存来存储 struct passwd，因此需要使用 C.free_passwd(pw) 来释放它。

总结：
1. Go 中的 C 结构体：如果你在 Go 中直接声明 C 结构体（例如 var pw C.struct_passwd），那么这些结构体的内存是由 Go 管理的，不需要手动释放
1. 动态分配的 C 结构体：如果你使用 C 函数（例如 malloc）来动态分配 C 结构体的内存，那么你需要在不再使用时手动调用 free 释放内存

注意：
1. 内存泄漏：忘记释放通过 malloc 或类似函数分配的内存会导致内存泄漏, 确保在不再使用该内存时释放它
1. cgo 和垃圾回收：Go 的垃圾回收机制不会管理通过 malloc 分配的 C 内存，因此这些内存必须由你显式地释放

#### `*_Ctype_char` 的内存管理
在 cgo 中，C.CString 返回一个 *C.char 类型的指针（即 *_Ctype_char），该指针指向 C 语言的内存. 因为这个内存是通过 malloc 分配的（内部实现依赖于 C.CString），Go 的垃圾回收器（GC）不会管理它。因此，必须手动释放这个内存，否则会导致内存泄漏.

### go c换传数组
go->c:
```go
package main
/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int loop(int** list_data, int leng, char** data)
{
  int* m = (int*)list_data;
  int sum = 0;
  for(int i=0; i<leng; i++)
  {
    sum += m[i];
  }
  *data = "finised task"; // "finised task"是字符串常量, 分配在静态存储区, 不用C.free()
  return sum;
}
*/
import "C"
import (
    "unsafe"
    "fmt"
)
func GoSilence2CArray() {
    var ids = []int32{1, 2, 3, 5}
    var res *C.char
    length := C.int(len(ids))
    le := C.loop((**C.int)(unsafe.Pointer(&ids[0])), length, &res)
    fmt.Println(le)
    fmt.Println(C.GoString(res))
    fmt.Println(ids)
}
func main() {
    GoSilence2CArray()
}
```

```go
package main

/*
#include<stdio.h>

void slice(int *a){
	for(int i=0;i<4;i++){
		printf("%d\n",a[i]);
	}
}

*/
import "C"
import (
	"fmt"
)

func main() {

	intSlice := []C.int{108880, 18, 28, 83, 488} //使用cgo类型C.int
	fmt.Println(&intSlice[0])
	C.slice(&intSlice[0])
}
```

c->go:
```go
package main
/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

typedef struct{
   char* name;
}person;

person* get_person(int n){
   person* ret = (person*)malloc(sizeof(person) * n);
   for(int i=0;i<n;i++){
      ret[i].name="wu";
   }
   return ret;
}
*/
import "C"
import (
    "unsafe"
    "fmt"
)
func CArray2GoSilence() {
    size := 2
    person := C.get_person(C.int(size))
    person_array := (*[2]C.person)(unsafe.Pointer(person))
    var names []string
    for i := 0; i < size; i++ {
        name := C.GoString(person_array[i].name)
        names = append(names, name)
    }
    for _, name := range names {
        fmt.Println(name)
    }
    C.free(unsafe.Pointer(person))
}
func main() {
    CArray2GoSilence()
}
```

go->c, 多维数组:
```go
package main
/*
#include <stdio.h>
#include <string.h>

void fill_array(char *s) 
{
    strcpy(s, "cobbliu");
}

void fill_2d_array(char **arr, int columeSize) {                                                                                                                                                                  
    strcpy((char*)(arr+0*sizeof(char)*columeSize), "hello");
    strcpy((char*)(arr+1*sizeof(char)*columeSize/sizeof(char*)), "cgo");
}
*/
import "C"
import "fmt"
import "unsafe"

func main() {
        var dir [10]byte
		C.fill_array((*C.char)(unsafe.Pointer(&dir[0])))
		fmt.Println(dir)
        fmt.Println(string(dir[:])) // go屏蔽掉了多余的`\0`
        //var dirs [4][16]byte                                                                                                                                                                                    
        dirs := make([][]byte, 4)
        for i := 0; i < 4; i++ {
                dirs[i] = make([]byte, 16)
        }

        C.fill_2d_array((**C.char)(unsafe.Pointer(&dirs[0][0])), C.int(16))
        fmt.Println(dirs)                   
}
```

### cgo无法链接so
```bash
# CGO_CFLAGS="-I/usr/local/include/rocksdb" CGO_LDFLAGS="-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd" go get github.com/linxGnu/grocksdb
# CGO_CFLAGS="-I/usr/local/include/rocksdb" CGO_LDFLAGS="-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd" go build # 不加CGO_CFLAGS,CGO_LDFLAGS编译时会无法链接到librocksdb.so, 导致编译不报错, 但执行时崩溃报"非法指令".
```

### `go install -v github.com/go-delve/delve/cmd/dlv@latest`报`_cgo_export.c:3:10: fatal error: stdlib.h: No such file or directory`
`apt install g++`

### pkg
`// #cgo pkg-config: udev`=`// #cgo LDFLAGS: -ludev`

### 处理union & 清理CString
c实现set union的函数, go调用该函数

```go
// #include "sf_disk_lib_h"
/*
void set_aksk(sfdisklib_connect_params * connect_params, char *ak, char* sk) {
	connect_params->creds.ak_sk.access_key = ak;
	connect_params->creds.ak_sk.secret_key = sk;
}
*/

func NewParams(ak, sk string) (*C.sfdisklib_connect_params, []*C.char){
	cAk := C.CString(ak)
	cSk := C.CString(sk)

	var cParams C.sfdisklib_connect_params
	var cParams = []*C.char{cAk, cCk}
	C.set_aksk(&cParams, cAk, cSk)

	return &cParams, cParams
}

func freeCStrings(ls []*C.char) {
	for i:=range ls {
		C.free(unsafe.Pointer(ls[i]))
	}
}

func main() {
	cParams, toFree:=NewParams("ak", "sk")
	defer freeCStrings(toFree)
}
```

### go buffer给c
```go
goBuffer:=make([]byte, 4096)
cBuffer:=(*C.uint8)(unsafe.Pointer(&goBuffer[0]))
C.sfdisklib_write(dh, offset/512, batchRead/512, cBuffer)
```

### error: unknown type name 'xxx'
比如报error: unknown type name 'mode_t'

原因:
1. 没有导入头文件
2. 没有按正确的顺序导入头文件(需先依赖)

### use of cgo in test tdb_test.go not supported
只能在非`_test.go`中使用cgo

### cgo 宏
Go 无法直接调用 C 宏, 因此需要通过 C 函数封装宏的逻辑或用go代码实现相同逻辑

举例:
```go
package main

/*
#cgo CFLAGS: -I/path/to/your/c/include
#include <wbc.h> // 假设 WBC_ERROR_IS_OK 和 WBC_ERR_SUCCESS 定义在 wbc.h 中

// 封装 C 的宏定义
int isErrorOk(int x) {
    return WBC_ERROR_IS_OK(x);
}
*/
import "C"
import "fmt"

func main() {
    // 假设有一个错误码
    errCode := C.WBC_ERR_SUCCESS

    // 调用 C 函数检查错误码
    result := C.isErrorOk(errCode)

    // 输出结果
    if result == 1 {
        fmt.Println("Error is OK (WBC_ERR_SUCCESS)")
    } else {
        fmt.Println("Error is NOT OK")
    }
}
```

### [获取errno](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/ch2-04-func.html)
from `Go语言高级编程#2.4 函数调用`

### 内存泄漏
方法:
- `GO_CFLAGS="-O0 -g" go build -gcflags="all=-N -l" -ldflags '-extldflags="-fsanitize=address -fsanitize=leak"' -o xxx`
- `valgrind --leak-check=full --track-origins=yes ./your_go_program`
