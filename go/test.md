### FAQ

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