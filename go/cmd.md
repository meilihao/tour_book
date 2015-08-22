## go build

主要用于编译代码.在包的编译过程中，若有必要，会同时编译与之相关联的包.go build编译链接后的可执行程序放在源程序目录.

注意:如果是普通包，当你执行go build之后，它不会产生任何文件.如果**你需要在$GOPATH/pkg下生成相应的文件，那就得执行go install**.

### 参数设置
#### 调试
编译时,如果编译的结果需要gdb调试则使用参数-gcflags "-N -l",这样可以忽略Go内部做的一些优化，比如聚合变量和函数等优化.

    go build -gcflags "-N -l"
#### 优化
如果编译的结果需要发布.则使用-ldflags "-s -w",可以去掉调试信息,减小大约一半的大小:
- `-s`: 去掉符号信息。
- `-w`: 去掉DWARF调试信息。


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

