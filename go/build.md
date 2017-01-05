# Golang 编译

## 条件编译

### 编译标签(`build tag`)：

```
// +build { GOOS }, { GOOS }, { !GOOS }
```
```
// +build (linux AND 386) OR (darwin AND (NOT cgo))
```
```
// +build linux darwin
// +build 386
```

1. 以`// +build`开头,一个源文件里可以有多个编译标签，多个编译标签之间是逻辑"与"的关系
1. 支持 GOOS 与 GOARCH 并可以具有多个值，用`,`分割， 例如：`// +build linux, darwin, freebsd`
1. 支持 不等条件 `!` ， 例如：`\\ +build !windows`,即不在windows环境下时，均可编译此文件。
1. 支持 与或非 逻辑， AND OR NOT
1. 条件编译需要前后空一行，否则无法识别

### 文件后缀：

```
xxx_{ GOOS }.go xxx_{ GOOS }_{ GOARCH }.go
```

1. 支持 GOOS ，例如： curl_windows.go
1. 支持 GOARCH， 例如： curl_386.go
1. 支持 上述两种叠加，但不可调换顺序 xxx_{ GOOS }_{ GOARCH }.go ，例如： curl_windows_amd64.go

### 如何选择

这两者可以叠加使用，但注意不要出现冗余，如：curl_windows.go 里面写`// +build windows`则重复了.
如果编译的文件是一一对应关系的话，使用文件后缀更简单些，如对每个 GOOS 生成一个文件.
如果有复杂条件的话，可以使用标签编译方式。如：
```
curl_windows.go 对应 windows 平台。
curl_others.go 里面写 \\ +build !windows 对应 非windows 平台
```

## 交叉编译

Golang 1.5及以上修改`GOOS和GOARCH`即可.如：

```
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test
$ file test
```

[GOOS 与 GOARCH 支持的参数](https://golang.org/doc/install/source#environment)

## 编译参数

### go get

-u - 强制使用网络去更新包和它的依赖包。当不加 -u 时，如果本地以存在，则不再从远程加载( 更新 )
-d - 只下载不安装
-v - 与 go build 参数含义相同
-x - 与 go build 参数含义相同

特性：
使用 go get 时，会自动切换到与当前 go version 一样的分支，如：当前 go version 是 1.5 则会自动查询 tags/go1 or branch/go1

支持的版本控制系统：
Mercurial
Git
Subversion
Bazaar

### go build

1. 不包含 main() 的话，执行后，不会产生任何文件.如果想要生成 package 或 可执行文件 需要使用`go install`
1. 包含 main() 的话，执行后，会在当前目录生成 package 或 可执行性文件
1. 以 `_.` 开头的文件会被忽略，例如：`_xxx.go 或者 .xxx.go`
1. 支持上面描述的 条件编译.

常用参数：

-o - 指定编译的文件名，可以带上路径
-a - 强行对所有涉及到的代码包（包含标准库中的代码包）进行重新构建，即使它们已经是最新的了
-n - 打印编译期间所用到的其它命令，但是并不真正执行它们
-v - 打印出那些被编译的代码包的名字
-x - 打印编译期间所用到的其它命令
-work - 打印出编译时生成的临时工作目录的路径，并在编译结束时保留它。在默认情况下，编译结束时会删除该目录

### go install

支持绝大多数的 go builid 参数，在 go build 执行的基础上，即：如果定义了 GOBIN 的话，会在此目录下生成可执行文件。( 前提是需要包含 main() 方法 )， 而 go build 只会在当前文件夹下生成 可执行性文件

### go run

包含了两个动作，编译 + 运行。
与 go build | go install 一样，支持它们的参数。

### go test

-c - 生成用于运行测试的可执行文件，但不执行它。
-i - 安装/重新安装运行测试所需的依赖包但不编译和运行测试代码。

### 有用的参数

-i - 安装相应的包。编译 +go install
-race - 开启编译的时候自动检测数据竞争的情况，目前只支持64位的机器
