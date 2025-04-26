# defer
defer是Go语言提供的一种用于注册延迟调用的机制：让函数或语句可以在当前函数执行完毕后（包括通过return正常结束或者**panic导致的异常结束**）执行.

> 在defer函数定义时，对外部变量的引用是有两种方式的，分别是作为函数参数和作为闭包引用.

> defer将相关函数注册到其所在goroutine用于存放deferred函数的栈数据结构中, 这些deferred函数将在执行defer的函数退出前被按后进先出(LIFO)的顺序调度执行

defer关键字后面的表达式是在将deferred函数注册到deferred函数栈的时候进行求值的.

## defer + return
```go
return xxx
```

可拆解为:
```go
1. 返回值 = xxx
2. call defer
3. return
```

## defer +　recover
recover()函数只在defer的上下文中才有效（且只有通过在defer中用匿名函数调用才有效），直接调用的话，只会返回 nil

注意:
1. cgo中crash就无法恢复

## defer作用域
```go
func main() {
    {
        defer fmt.Println("defer runs")
        fmt.Println("block ends")
    }
    
    fmt.Println("main ends")
}

$ go run main.go
block ends
main ends
defer runs
```

defer 并不是在退出当前代码块的作用域时执行的，其只会在当前函数和方法返回之前才被调用.


## defer 的底层实现
参考:
- [defer的底层原理](https://golang.org/ref/spec#Defer_statements)
- [理解 Go 语言 defer 关键字的原理](https://draveness.me/golang-defer)

从汇编代码上看到, defer 的底层实现主要由两个函数构成：
```go
// deferprocStack是Go 1.13开始引入的
func deferproc(siz int32, fn *funcval) / func deferprocStack(d *_defer)
func deferreturn(arg0 uintptr)
```

每次defer语句执行的时候，会把函数“压栈”，函数参数会被拷贝下来；当**外层函数**（非代码块，如一个for循环）退出时，defer函数按照定义的逆序执行；如果defer执行的函数为nil, 那么会在最终调用函数的产生panic.

defer语句并不会马上执行，而是会进入一个栈，函数return前，会按先进后出的顺序执行.

## defer性能
参考:
- [深入理解 Go 语言 defer](https://zhuanlan.zhihu.com/p/63354092)

defer 的实现主要依靠编译器和运行时的协作, 相关的三种机制：

- 堆上分配 · 1.1 ~ 1.12

    - 编译期将 defer 关键字转换成 runtime.deferproc 并在调用 defer 关键字的函数返回之前插入 runtime.deferreturn
    - 运行时调用 runtime.deferproc 会将一个新的 runtime._defer 结构体追加到当前 Goroutine 的链表头
    - 运行时调用 runtime.deferreturn 会从 Goroutine 的链表中取出 runtime._defer 结构并依次执行

    堆上分配严重影响defer延迟调用的执行效率
- 栈上分配 · 1.13
    - 当该关键字在函数体中最多执行一次时，编译期间的 cmd/compile/internal/gc.state.call 会将结构体分配到栈上并调用 runtime.deferprocStack

        满足该条件的defer延迟调用（标准库中93%的延迟调用满足此条件）在栈上分配, 从而提高了defer延迟调用的执行效率

    在Go 1.14版本中,defer性能提升巨大,已经和不用defer的性能相差很小了.
- 开放编码 · 1.14 ~ 现在
    - 编译期间判断 defer 关键字、return 语句的个数确定是否开启开放编码优化
    - 通过 deferBits 和 cmd/compile/internal/gc.openDeferInfo 存储 defer 关键字的相关信息
    - 如果 defer 关键字的执行可以在编译期间确定，会在函数返回前直接插入相应的代码，否则会由运行时的 runtime.deferreturn 处理

相关源码在`/src/runtime/panic.go`