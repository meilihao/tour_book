# unsafe
参考:
- [深度解密Go语言之unsafe](http://www.cnblogs.com/qcrao-2018/p/10964692.html)

为了安全考虑, Golang 的指针是类型安全的, 相比 C 的指针有很多限制, 好处就是可以享受指针带来的便利，又避免了指针的危险性; 坏处就是缺少C指针的灵活.

unsafe包的主要使用场景:
1. 与操作系统以及非Go编写的代码的通信
1. 高效类型转换

限制:
1. 不能对指针做算术运算
    ```go
    v := 1
    p := &v

    p++ // error
    ```
1. 不同类型的指针不能相互转换
    ```go
    v := 1
    var f *float32

    f = &a // error
    ```
1. 不同类型的指针不能使用 == 或 != 比较
    只有在两个指针类型相同或者可以相互转换的情况下，才可以对两者进行比较. 另外指针可以通过 == 和 != 直接和 nil 作比较
1. 不同类型的指针变量不能相互赋值

packege unsafe可避开这些限制, 可直接操作内存. 在某些情况下，它会使代码更高效，当然也更危险. 因为它是不安全的，官方并不建议使用.

避开限制举例:
- unsafe 包提供的非类型安全的指针 unsafe.Pointer

    unsafe.Pointer的特别之处在于，它可以绕过 Go 语言类型系统的检查，与任意的指针类型互相转换。也就是说，如果两种类型具有相同的内存结构（layout）, 因此可以将unsafe.Pointer当做桥梁，让这两种类型的指针相互转换，从而实现同一份内存拥有两种不同的解读方式. 比如说，`[]byte和string`其实内部的存储结构都是一样的.
- 绕过 Go 语言的类型系统，操作一个结构体的未导出成员

> unsafe 包用于 Go 编译器，在编译阶段使用.



## 为什么有 unsafe
Go 语言类型系统是为了安全和效率设计的，有时安全会导致效率低下. 有了 unsafe 包就可以利用它绕过类型系统的低效. 阅读 Go 源码，会发现有大量使用 unsafe 包的例子.

## unsafe 实现原理
```go
type ArbitraryType int // Go的任意类型
type Pointer *ArbitraryType // 指向任意类型的指针, 类似c里的`void*`
```

unsafe 包还有其他三个函数：
- func Sizeof(x ArbitraryType) uintptr
    返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小
- func Offsetof(x ArbitraryType) uintptr
    返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员
- func Alignof(x ArbitraryType) uintptr
    获取变量的对齐值，除 int、uintptr 这些依赖CPU位数的类型，基本类型的对齐值都是固定的. 结构体的对齐值取他的成员对齐值的最大值

> 这三个函数执行的结果和操作系统、编译器相关，所以是不可移植的

unsafe 包提供了 2 点重要的能力：
- 任何类型的指针和 unsafe.Pointer 可以相互转换
- uintptr 类型和 unsafe.Pointer 可以相互转换

unsafe.Pointer 不能直接进行数学运算，但可以把它转换成 uintptr，对 uintptr 类型进行数学运算，再还原成 unsafe.Pointer.

> uintptr 并没有指针的语义, 意味着 uintptr 所指向的对象会被 gc 回收; 而 unsafe.Pointer 有指针语义，可以保护它所指向的对象在"有用"的时候不会被gc回收.
> 在 /usr/local/go/src/cmd/compile/internal/gc/unsafe.go 中可以看到编译期间 Go 对 unsafe 包中函数的处理

### [unsafe.Pointer的安全使用模式](https://pkg.go.dev/unsafe)
go1.24在unsafe的文档中定义了6条安全使用模式
1. `*T1 -> unsafe.Pointer -> *T2`

    本质就是内存块的重解释:将原本解释为T1类型的内存重新解释为T2类型, **它不等价于go的类型转换**

    转换后类型T2的对齐系数不能比转换前类型T1的对齐系数更严格,即Alignof(T1) >= Alignof(T2). 在x86平台(属于复杂指令集)上, 对这种未对齐的地址进行指针解引用并不会造成严重后果(可能会对性能有少许影响),但是在一些对内存地址对齐比较严格的平台(如SPARC、ARM等)上,对未对齐内存地址进行指针解引用可能会出现“总线错误”等无法恢复的异常情况并导致程序崩溃
2. `unsafe.Poiner -> uintptr`且不再转回Poiner
3. 模拟指针运算

    注意事项:
    1. 不要越界
    2. unsafe.Pointer -> uintptr -> unsafe.Pointer的转换**要在一个表达式中**, go编译器会保证两次转换期间对象的有效性, 否则对象可能被gc

    uintptr仅是一个整型值,它无法起到对象引用的效果, 无法阻止GC回收内存对象
4. 调用syscall.Syscall系列函数时指针类型到uintptr类型参数的转换, 同3转换操作需在一个表达式中
5. 将reflect.Value.Pointer或reflect.Value.UnsafeAddr转换为指针, 同3转换操作需在一个表达式中
6. reflect.SliceHeader和reflect.StringHeader必须通过模式1构建

Go 1.14编译器在-race和-msan命令行选型开启的情况下,会执行-d=checkptr检查, 即对unsafe.Pointer进行模式1,2两项合规性检查

## unsafe 使用
### 获取slice header的信息
slice header 的结构体定义：
```go
// runtime/slice.go
type slice struct {
    array unsafe.Pointer // 底层数组的指针
    len   int // 长度 
    cap   int // 容量
}
```

调用 make 函数新建一个 slice，底层调用的是 makeslice 函数, 返回的是 unsafe.Pointer, 但实际使用中是 slice 结构体???：
```go
func makeslice(et *_type, len, cap int) unsafe.Pointer
```

演示:
```go
func main() {
    s := make([]int, 9, 20)
    p := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
    fmt.Println(p)
    
    var Len = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))
    fmt.Println(Len, len(s)) // 9 9

    var Cap = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
    fmt.Println(Cap, cap(s)) // 20 20
}
```

### 获取map len
map的定义:
```go
// runtime/map.go
type hmap struct {
    count     int
    flags     uint8
    B         uint8
    noverflow uint16
    hash0     uint32

    buckets    unsafe.Pointer
    oldbuckets unsafe.Pointer
    nevacuate  uintptr

    extra *mapextra
}
```

和 slice 不同的是，makemap 函数返回的是 hmap 的**指针**：
```go
func makemap(t *maptype, hint int, h *hmap) *hmap
```

演示:
```go
	mp := make(map[string]int)
	mp["qcrao"] = 100
	mp["stefno"] = 18

	count := **(**int)(unsafe.Pointer(&mp))
	fmt.Println(count, len(mp)) // 2 2
```

### Offsetof 获取成员偏移量
对于一个结构体，通过 offset 函数可以获取结构体成员的偏移量，进而获取成员的地址，读写该地址的内存，就可以达到改变成员值的目的.

> 结构体会被分配一块连续的内存，结构体的地址也代表了第一个成员的地址.

演示:
```go
package main

import (
	"fmt"
	"unsafe"
)

type Programmer struct {
	name     string
	age      int
	language string
}

func main() {
	p := Programmer{"stefno", 1, "go"}
	fmt.Println(p)

	name := (*string)(unsafe.Pointer(&p))
	*name = "qcrao"

	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.language)))
	*lang = "Golang"

	fmt.Println(p)

	lang2 := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Sizeof(int(0)) + unsafe.Sizeof(string(""))))
	*lang2 = "Golang2"

	fmt.Println(p)
}
```

### string 和 slice 的相互转换
实现字符串和 bytes 切片之间的转换，要求是 zero-copy.

这里需要了解 slice 和 string 的底层数据结构：
```go
type StringHeader struct {
    Data uintptr
    Len  int
}

type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```
上面是reflect包下的结构体, 只需要共享底层 []byte 数组就可以实现 zero-copy 。

通过构造 slice header 和 string header，来完成 string 和 byte slice 之间的转换:
```go
func string2bytes(s string) []byte {
    stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

    bh := reflect.SliceHeader{
        Data: stringHeader.Data,
        Len:  stringHeader.Len,
        Cap:  stringHeader.Len,
    }

    return *(*[]byte)(unsafe.Pointer(&bh))
}

func bytes2string(b []byte) string {
    sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

    sh := reflect.StringHeader{
        Data: sliceHeader.Data,
        Len:  sliceHeader.Len,
    }

    return *(*string)(unsafe.Pointer(&sh))
}
```

