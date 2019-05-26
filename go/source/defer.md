# defer
defer是Go语言提供的一种用于注册延迟调用的机制：让函数或语句可以在当前函数执行完毕后（包括通过return正常结束或者**panic导致的异常结束**）执行.

[defer的底层原理](https://golang.org/ref/spec#Defer_statements): 每次defer语句执行的时候，会把函数“压栈”，函数参数会被拷贝下来；当**外层函数**（非代码块，如一个for循环）退出时，defer函数按照定义的逆序执行；如果defer执行的函数为nil, 那么会在最终调用函数的产生panic.

defer语句并不会马上执行，而是会进入一个栈，函数return前，会按先进后出的顺序执行.

> 在defer函数定义时，对外部变量的引用是有两种方式的，分别是作为函数参数和作为闭包引用.

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