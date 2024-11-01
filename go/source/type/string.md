# string
```go
// $GOROOT/src/runtime/string.go
type stringStruct struct {
    str unsafe.Pointer
    len int
}

// 实例化一个字符串对应的函数
func rawstring(size int) (s string, b []byte) {
    p := mallocgc(uintptr(size), nil, false)
    stringStructOf(&s).str = p
    stringStructOf(&s).len = size
    *(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}
}
```

rawstring调用后,新申请的内存区域还未被写入数据, 该slice就是供后续运行时层向其中写入数据("hello")用的。**写完数据后,该slice就可以被回收掉了**.