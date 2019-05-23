# slice
env: go version go1.12.5 linux/amd64

参考:
- [深度解密Go语言之slice](https://studygolang.com/articles/19469)

## 模型
[slice](runtime/slice.go):
```go
type slice struct {
	array unsafe.Pointer // 底层数组指针
	len   int // 长度
	cap   int // 容量
}
```

## 创建slice
创建 slice 的方式有以下几种：
1. 直接声明: var slice []int
1. new: slice := *new([]int)
1. 字面量 slice := []int{1,2,3,4,5}
1. make	slice := make([]int, 5, 10)
1. 从切片或数组"截取":  slice := array[1:5] 或 slice := sourceSlice[1:5]

通过汇编代码得知make的具体执行函数:
```go
var s0 []int // nil切片, 底层数组为nil, 不调用runtime.makeslice

s:=make([]int,0) // runtime.makeslice, 底层数组不为nil

fmt.Println(s)

var a1 = *(*[3]int)(unsafe.Pointer(&s0))
var a2 = *(*[3]int)(unsafe.Pointer(&s1))
fmt.Println(a1) // [0 0 0]
fmt.Println(a2) // [5714048 0 0]
```
