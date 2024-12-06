# test
包内测试: 将测试代码放在与被测包同名的包中的测试方法
包外测试: 将测试代码放在名为被测包包名+"_test"的包中的测试方法

包内测试这种方法本质上是一种白盒测试方法。由于测试代码与被测包源码在同一包名下,测试代码可以访问该包下的所有符号,无论是导出符号还是未导出符号;并且由于包的内部实现逻辑对测试代码是透
明的,包内测试可以更为直接地构造测试数据和实施测试逻辑,可以很容易地达到较高的测试覆盖率.

包内测试缺点:
1. 测试代码自身需要经常性的维护

	包内测试的白盒测试本质意味着它是一种面向实现的测试。测试代码的测试数据构造和测试逻辑通常与被测包的特定数据结构设计和函数/方法的具体实现逻辑是紧耦合的, 需要同步配合调整.
1. 硬伤:包循环引用

与包内测试本质是面向实现的白盒测试不同,包外测试的本质是一种面向接口的黑盒测试. 包外测试由于将测试代码放入独立的包中,它更适合编写偏向集成测试的用例.

包外测试缺点: 测试盲区. 由于测试代码与被测试目标并不在同一包名下,测试代码仅有权访问被测包
的导出符号,并且仅能通过导出API这一有限的“窗口”并结合构造特定数据来验证被测包行为。在这样的约束下,很容易出现对被测试包的测试覆盖不足的情况. 解决版本就是export_test.go, 该文件中的代码位于被测包名下,但它既不会被包含在正式产品代码中(因为位于_test.go文件中),又不包含任何测试代码仅用于将被测包的内部符号在测试阶段暴露给包外测试代码.

倾向于优先选择包外测试:
1. 优先保证被测试包导出API的正确性
1. 可从用户角度验证导出API的有效性
1. 保持测试代码的健壮性,尽可能地降低对测试代码维护的投入
1. 不失灵活!可通过export_test.go来导出需要的内部符号,满足窥探包内实现逻辑的需求

## FAQ

### `go test`时提示"undefined: XXX"

将其测试的源码文件也作为参数即可

```shell
# 举例, cd golang源码目录/src/database/sql/driver
$ go test types_test.go  # 将提示:"undefined: ValueConverter"
$ go test types_test.go  types.go driver.go # ok
```

### `go test`时未输出调试用的fmt.Println和t.Log信息

`go test`默认不打印调试信息,但加上参数`-v`即可打印,其表示在测试运行结束后打印出所有在测试运行过程中被记录的日志.此时,fmt.Println和t.Log区别是:fmt是测试运行时打印,t是测试结束后打印.


	go test -v xxx_test.go

### `go test -v -run=XXX`,预先设置的调试信息未输出，直接`PASS ok`

XXX函数名字和*_test.go里的待测试函数的名称不一致.


### `go test -v -bench=BenchmarkSingle -benchtime 200x`如何确定b.N
其实BenchmarkSingle会被执行2次, 第一次b.N=1确定基准, 第二次b.N=200执行真正测试.

### go test build
`go test -c -v  -timeout 30s -run ^TestVersion$ fstack/pkg/plugins/system`

只有当个test case时, 可直接运行; 否则需要指定test case: `./xxx.test -test.run TestYYY`

### go test 禁用缓存
每当执行 go test 时, 如果功能代码和测试代码没有变动, 则在下一次执行时, 会直接读取缓存中的测试结果，并通过 (cached) 进行标记.

要禁用测试缓存, 可以通过`-count=1`标志来实现.

> 其他方法: `go clean -testcache`: expires all test results

### 查找包内/包外测试文件
```bash
go list -f={{.TestGoFiles}} .  # 包内
go list -f={{.XTestGoFiles}} . # 包外
```