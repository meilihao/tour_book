# 位运算

## 判断操作系统位数
```go
func main() {
    bit := 32 << (^uint(0) >> 63)
    println(bit)
}
```

对于32位系统：

    ^unit(0)：232 − 1
    (232 − 1) >> 63，得到0
    32 << 0，等于 32

对于64位系统：

    ^unit(0)：264 − 1
    (264 − 1) >> 63，得到1
    32 << 1，等于 64

## 对齐
```go
const msgAlignTo = 4

func msgAlign(len int) int {
	return (len + nlmsgAlignTo - 1) & ^(nlmsgAlignTo - 1)
}
```

## 奇偶
```go
if i & 1 == 1 {
	// is 奇数
}
```