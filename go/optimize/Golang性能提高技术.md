# Golang 性能提高技术
参考:
- [Golang 性能提高技术----基础编码原则](https://www.jianshu.com/p/0dafe1059fdc)

基础编码原则:
- 高级设计 ：指的是程序整体的设计，采用适当的算法和数据结构
- 基本编码原则 ：从指令的角度考虑，开发中应如何编码，才能减少执行的指令
- 低级优化 ：针对现代处理器，如何让cpu的流水线尽量饱合

测试文件g_test.go:
```go
package test

import (
	"testing"
	"runtime/debug"
)

func init() {
    debug.SetGCPercent(-1)
}

var testData = CreateTestData()

func GetValue(index int) int {
   return testData[index]
}

//获取测试数据长度
func GetDataLen() int{
   return len(testData)
}

//创建一个6000000大小的整型切片, 选用比较大的测试数据，是为了减少运行中其它因素的干扰影响
func CreateTestData()[]int  {
   data := make([]int,6000000)
   for index,_ := range data{
      data[index] = index % 128
   }
   return data
}

func toSum1(result *int)  {
    for i:=0;i< GetDataLen();i++{ // GetDataLen()没必要在循环里重复调用
        *result += GetValue(i)
    }
}

// toSum1 -> toSum2 : 消除连续的函数调用
func toSum2(result *int)  {
    dataLength := GetDataLen()
    for i:=0;i< dataLength;i++{
        *result += GetValue(i) // 减少不必要的函数调用过程, 直接访问slice即可
    }
}

func GetData() []int {
    return testData
}

// toSum2 -> toSum3 : 减少函数的调用过程，而对切片内容的访问还是需要的，所以看到的效果不是很明显
func toSum3(result *int)  {
    data := GetData()
    dataLength := len(data)
    for i:=0;i< dataLength;i++{
        *result += data[i]
    }
}

/* toSum3的`*result += data[i]`的汇编步骤：
    ...
    0x0053 00083 (g_test.go:47)	MOVQ	"".dataLength+32(SP), AX ; AX = dataLength
	0x0058 00088 (g_test.go:47)	CMPQ	"".i+24(SP), AX          ; CMPQ i dataLength
	0x005d 00093 (g_test.go:47)	JLT	97                           ; i < dataLength时,jump 97 ; instructions["JLT"] = x86.AJLT  /* less than (signed) (SF != OF) */
	0x005f 00095 (g_test.go:47)	JMP	160
	0x0061 00097 (g_test.go:48)	PCDATA	$2, $2
	0x0061 00097 (g_test.go:48)	MOVQ	"".result+80(SP), AX     ; AX = result的指针地址
	0x0066 00102 (g_test.go:48)	TESTB	AL, (AX)
	0x0068 00104 (g_test.go:48)	PCDATA	$2, $0
	0x0068 00104 (g_test.go:48)	MOVQ	(AX), AX			 	 ; AX = *result 			; 1. 通过result地址值，从内存取出内容放在寄存器AX中
	0x006b 00107 (g_test.go:48)	MOVQ	"".i+24(SP), CX          ; CX= i
	0x0070 00112 (g_test.go:48)	PCDATA	$2, $1
	0x0070 00112 (g_test.go:48)	MOVQ	"".data+40(SP), DX       ; DX = data
	0x0075 00117 (g_test.go:48)	CMPQ	"".data+48(SP), CX		 ; len(data) i
	0x007a 00122 (g_test.go:48)	JHI	126							 ; len(data) > i时,jump 126 ; instructions["JHI"] = x86.AJHI  /* higher (unsigned) (CF = 0 && ZF = 0) */
	0x007c 00124 (g_test.go:48)	JMP	170
    0x007e 00126 (g_test.go:48)	PCDATA	$2, $3
	0x007e 00126 (g_test.go:48)	MOVQ	"".result+80(SP), BX     ; BX = result的指针地址	  
	0x0083 00131 (g_test.go:48)	TESTB	AL, (BX)
	0x0085 00133 (g_test.go:48)	PCDATA	$2, $4
	0x0085 00133 (g_test.go:48)	MOVQ	(DX)(CX*8), CX           ; CX = data[i]				; 2. 再通过切片数组的首地址获取第i个元素到寄存器CX中
	0x0089 00137 (g_test.go:48)	ADDQ	CX, AX                   ; AX+=CX					; 3. 两者相加
	0x008c 00140 (g_test.go:48)	PCDATA	$2, $0
	0x008c 00140 (g_test.go:48)	MOVQ	AX, (BX)                 ; *result = AX				; 4. 将寄存AX写回 result 指向的内存地址
*/

// toSum3 -> toSum4 : 消除不必要的存储器引用
func toSum4(result *int)  {
    k := *result
    data := GetData()
    dataLength := len(data)
    for i:=0;i< dataLength;i++{
        k += data[i]
    }
    *result = k
}

// 循环展开是通过程序的变换，通过增加每次迭代计算的元素， 减少循环的迭代次数.
// 减少迭代次数，意味着可以节省掉一半的条件判断指令执行(i<datalength)，其次每次迭代可以减少几个周期的延迟，因为在后续的指令执行需要等待控制语句更新程序计数器才能往下继续执行，在计数器没有更新时，是不能将计算结果更新到寄存器或内存中的
func toSum5(result *int)  {
	k := *result
	data := GetData()
	dataLength := len(data)
	for i:=1;i< dataLength;i+=2{
		k += data[i] + data[i - 1] // 此处实际运行中可能是先执行k+=data[i]，等计算结果写入到寄存器中才能执行 k+=data[i-1]. 对于加法指令可能并不会有太多的周期延迟，但是如果是针对乘除指令就会比较明显, 可尝试提高并行性
	}
	if dataLength % 2 == 1{
		k += data[dataLength-1]
	}

	*result = k
}

// 提高并行性
// 对比toSum4() 和 toSum6()的性能对比又有了明显的提升. 在这里只是对程序进行2次展开2次并行的处理，如果想让流水线更多饱合，那还可以进行更多的展开和并行处理;
// 但不是越多越好，要考虑到**寄存器的有限**，如果累加值超过剩余寄存器的数量，增加多余的内存读写操作反而得不偿失，具体展开和并行次数还得根据程序所在的机器运行的情况决定
func toSum6(result *int)  {
	k1 := 0
	k2 := 0
	data := GetData()
	dataLength := len(data)

	for i:=1;i< dataLength;i+=2{
		k1 += data[i]
		k2 += data[i - 1]
	}

	//如果是传入的数量是奇数，则单独对最后一个数进行累加
	if dataLength % 2 == 1{
		k1 += data[dataLength-1]
	}

	*result = k1 + k2
}

func BenchmarkData1(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	// b.N = 1
	
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum1(sum)
	}

	b.StopTimer()
}

func BenchmarkData2(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum2(sum)
	}

	b.StopTimer()
}

func BenchmarkData3(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum3(sum)
	}

	b.StopTimer()
}

func BenchmarkData4(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum4(sum)
	}

	b.StopTimer()
}

func BenchmarkData5(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum5(sum)
	}

	b.StopTimer()
}

func BenchmarkData6(b *testing.B)  {
    var sum *int = new(int)
	*sum = 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		toSum6(sum)
	}

	b.StopTimer()
}
```

```sh
$ go test -c g_test.go // 生成测试可执行文件
$ ./test.test -test.bench=.
goos: linux
goarch: amd64
BenchmarkData1-12    	     200	   7632973 ns/op
BenchmarkData2-12    	     200	   7657932 ns/op
BenchmarkData3-12    	     200	   7881545 ns/op
BenchmarkData4-12    	     500	   2747947 ns/op // 优化到toSum4就已经足够了
BenchmarkData5-12    	     500	   2728008 ns/op
BenchmarkData6-12    	     500	   2763390 ns/op
PASS
ok  	command-line-arguments	11.959s
$
$ go test -c -gcflags "-N -l" g_test.go
$ ./test.test -test.bench=. # =`go test -c -gcflags "-N -l" g_test.go`
goos: linux
goarch: amd64
BenchmarkData1-12    	     100	  21786558 ns/op
BenchmarkData2-12    	     100	  17309308 ns/op
BenchmarkData3-12    	     100	  12404992 ns/op
BenchmarkData4-12    	     100	  11053417 ns/op
BenchmarkData5-12    	     200	   5976350 ns/op
BenchmarkData6-12    	     200	   6539464 ns/op
PASS
ok  	command-line-arguments	10.119s
```