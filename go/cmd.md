# cmd

参考: [go_command_tutorial](github.com/hyper0x/go_command_tutorial/blob/master/SUMMARY.md)

## go build
选项:
- -a : 强制重新构建所有包
- -x -v:让构建过程一目了然
- -race:让并发bug无处遁形
- -gcflags:传给编译器的标志选项集合


    -l:关闭内联
    -N:关闭代码优化
    -m:输出逃逸分析(决定哪些对象在栈上分配,哪些对象在堆上分配)的分析决策过程
    -S:输出汇编代码
- -ldflags:传给链接器的标志选项集合

    -X:设定包中string类型变量的值(仅支持string类型变量)
    -s:不生成符号表(symbol table)
    -w:不生成DWARF(Debugging With Attributed Record Formats)调试信息
- -tags:指定构建约束条件

主要用于编译代码.在包的编译过程中，若有必要，会同时编译与之相关联的包.go build编译链接后的可执行程序放在源程序目录.

注意:如果是普通包，当你执行go build之后，它不会产生任何文件.如果**你需要在$GOPATH/pkg下生成相应的文件，那就得执行go install**.

`go build -mod=vendor`可指定使用vendor.

### 参数设置
#### 调试
编译时,如果编译的结果需要gdb调试则使用参数-gcflags "-N -l",这样可以忽略Go内部做的一些优化，比如聚合变量和函数等优化.

    $ go build -gcflags "-N -l"
    $ gdb  ./httprouter
    (gdb) source /opt/go/src/runtime/runtime-gdb.py
    (gdb) ...
>进入gdb环境后需先运行`source /opt/go/src/runtime/runtime-gdb.py`命令以加载Go Runtime的支持，否则gdb中goroutine相关命令将无法运行,参考[官方文档](golang.org/doc/gdb).

#### 优化
如果编译的结果需要发布.则使用-ldflags "-s -w",可以去掉调试信息,减小大约一半的大小:
- `-s`: 去掉符号表和调试信息
- `-w`: 去掉DWARF调试信息


    go build -ldflags “-s -w”

## go install

与build命令相比，install命令在编译源码后还会将可执行文件或库文件安装到约定的目录下.

go install只会检查"参数指定的包所在的GOPATH"内的源码是否有更新，如果有则重新编译.对于依赖的其他GOPATH下的包，如果存在已经编译好的.a文件，则不会再检查源码是否有更新，不会重新编译([参考 : go install 的工作方式 ](http://blog.csdn.net/tiaotiaoyly/article/details/38517103)).

## go run

编译并运行Go程序.

### FAQ

#### 修改源码后,运行和预期不符

修改了依赖包代码之后使用go run XXX.go的情况(单一$GOPATH时,human包为依赖包,gomain目录为现有项目):

1. 如果将human包移到goman目录下，则无论怎么修改human.go，编译goman项目时都会得到预期结果.

1. 如果human包不在gomain目录下(即属于不同项目),当$GOPATH/pkg下没有相应的human.a时,修改human.go,编译goman项目时会得到预期结果;当$GOPATH/pkg下有相应的human.a时,修改human.go,编译goman项目时不会得到预期结果,需用`go install`更新`human.a`.

