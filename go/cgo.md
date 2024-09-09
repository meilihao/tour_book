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

### 处理union
c实现set union的函数, go调用该函数

### go buffer给c
```go
goBuffer:=make([]byte, 4096)
cBuffer:=(*C.uint8)(unsafe.Pointer(&goBuffer[0]))
```