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

Go 1.13 defer有性能改进:
Go 1.13之前，所有的defer延迟调用都是记录在堆上的，这严重影响了defer延迟调用的执行效率. 从Go 1.13开始，满足某些条件的defer延迟调用（标准库中93%的延迟调用满足此条件）将被记录在栈上而不是堆上，从而提高了defer延迟调用的执行效率.
而在Go 1.14版本中,defer性能提升巨大,已经和不用defer的性能相差很小了.

使用了参考里的bench test, 发现go1.13比go12.5快25%左右.
```
$ ./main.test -test.bench=. // go12.5, use deferproc
goos: linux
goarch: amd64
BenchmarkCall-12         	100000000	        12.5 ns/op
BenchmarkDeferCall-12    	50000000	        38.3 ns/op
PASS
$ ./main.test13 -test.bench=. // go1.13beta1, use deferprocStack
goos: linux
goarch: amd64
BenchmarkCall-12         	120874596	         9.96 ns/op
BenchmarkDeferCall-12    	42044380	        28.7 ns/op
PASS
```

