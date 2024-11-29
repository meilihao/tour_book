# interface
接口类型变量有两种内部表示, 即这两种表示分别用于不同接口类型的变量:
1. eface : 用于表示没有方法的空接口(empty interface)类型变量,即interface{}类型的变量
1. iface : 用于表示其余拥有方法的接口(interface)类型变量

```go
// $GOROOT/src/runtime/runtime2.go
type iface struct {
    tab *itab
    data unsafe.Pointer
}

type eface struct {
    _type *_type
    data unsafe.Pointer
}
```

这两种结构的共同点是都有两个指针字段,并且第二个指针字段的功用相同,都指向当前赋值给该接口类型变量的动态类型变量的值.

不同点在于eface所表示的空接口类型并无方法列表, 因此其第一个指针字段指向一个_type类型结构, 该结构为该接口类型变量的动态.

类型的信息:
```go
// $GOROOT/src/runtime/type.go
type _type struct {
    size uintptr
    ptrdata uintptr
    hash uint32
    tflag tflag
    align uint8
    fieldalign uint8
    kind uint8
    alg *typeAlg
    gcdata *byte
    str nameOff
    ptrToThis typeOff
```

而iface除了要存储动态类型信息之外,还要存储接口本身的信息(接口的类型信息、方法列表信息等)以及动态类型所实现的方法的信息, 因此iface的第一个字段指向一个itab类型结构:
```go
// $GOROOT/src/runtime/runtime2.go
type itab struct {
    inter *interfacetype
    _type *_type
    hash uint32
    _ [4]byte
    fun [1]uintptr
}

// $GOROOT/src/runtime/type.go
type interfacetype struct {
    typ _type
    pkgpath name
    mhdr []imethod
}
```

上面itab结构中的第一个字段inter指向的interfacetype结构存储着该接口类型自身的信息. interfacetype结构由类型信息(typ)、包路径名(pkgpath)和接口方法集合切片(mhdr)组成.

itab结构中的字段_type则存储着该接口类型变量的动态类型的信息,字段fun则是动态类型已实现的接口方法的调用地址数组.

虽然eface和iface的第一个字段有所差别,但tab和_type可统一看作动态类型的类型信息. Go语言中每种类型都有唯一的_type信息,无论是内置原生类型,还是自定义类型. Go运行时会为程序内的全部类型建立只读的共享_type信息表,因此拥有相同动态类型的同类接口类型变量的_type/tab信息是相同的. 而接口类型变量的data部分则指向一个动态分配的内存空间, 该内存空间存储的是赋值给接口类型变量的动态类型变量的值.

> Go提供了println预定义函数,可以用来输出eface或iface的两个指针字段的值. println在编译阶段会由编译器根据要输出的参数的类型将println替换为特定的函数,这些函数都定义在$GOROOT/src/runtime/print.go文件中

通过`go tool compile -S interface-internal-4.go > interface-internal-4.s`, 可在汇编中看到runtime是通过了convT2E和convT2I两个runtime包的函数(`$GOROOT/src/runtime/iface.go`)将动态类型变量赋值给接口类型变量的.

convT2E用于将任意类型转换为一个eface,convT2I用于将任意类型转换为一个iface。两个函数的实现逻辑相似,主要思路就是根据传入的类型信息(convT2E的_type和convT2I的tab._type)分配一块内存空
间,并将elem指向的数据复制到这块内存空间中,最后传入的类型信息作为返回值结构中的类型信息,返回值结构中的数据指针(data)指向新分配的那块内存空间.

传入它们的类型信息是依赖go编译器的. 。编译器知道每个要转换为接口类型变量(toType)的动
态类型变量的类型(fromType),会根据这一类型选择适当的convT2X函数, 见`$GOROOT/src/cmd/compile/internal/gc/walk.go`中的`func walkexpr(n *Node, init *Nodes) *Node {}`

装箱是一个有性能损耗的操作, 因此Go在不断对装箱操作进行优化,包括对常见类型(如整型、字符串、切片等)提供一系列快速转换函数:
```go
// $GOROOT/src/cmd/compile/internal/gc/builtin/runtime.go
// 实现在 $GOROOT/src/runtime/iface.go中
func convT16(val any) unsafe.Pointer //  val必须是一个uint16相关类型的参数
....
```