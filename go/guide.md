## 版本
- [golang发展历史](https://studygolang.com/topics/6369)

## install
go 环境变量配置
```bash
//etc/profile
#golang
export GOROOT=/opt/go
export GOPATH=/home/chen/git/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
>### go的环境变量
>- GOROOT：⽤于设定Go语⾔的根⽬录；标准命令的可执⾏⽂件在`$GOROOT/bin`中，标准库的源码⽂件在`$GOROOT/src`中，标准库的归档⽂件在`$GOROOT/pkg`中。
- GOPATH：⽤于设定⼯作区⽬录；可以包含⼀个或多个⼯作区⽬录的路径，每个⼯作区⽬录都应有src⼦⽬录。

>#### 环境变量说明
>- 当Go语⾔发现我们用import 语句导⼊了⼀个代码包时，会到以下⽬录查找该代码包的归档⽂件：
 1. `$GOROOT/pkg`⽬录
 2. `$GOPATH`包含的所有⼯作区⽬录的`pkg`⼦⽬录
- `$GOPATH[i]/src` ⽬录中的库源码⽂件总会被`go install`命令安装到`GOPATH[i]/pkg`中
- `$GOPATH[i]/src` ⽬录中的命令源码⽂件会被`go install`命令安装到 $GOBIN ⽬录中

## tools

```shell
gocode: go get -u -v github.com/nsf/gocode
godef: go get -u -v github.com/rogpeppe/godef
golint: go get -u -v github.com/golang/lint/golint
go-find-references: go get -u -v github.com/lukehoban/go-find-references
go-outline: go get -u -v github.com/lukehoban/go-outline
goreturns: go get -u -v github.com/sqs/goreturns
gorename: go get -u -v golang.org/x/tools/cmd/gorename
gooracle: go get -u -v golang.org/x/tools/cmd/oracle
go debug: go get -u -v github.com/derekparker/delve/cmd/dlv
goimports: go get -u -v golang.org/x/tools/cmd/goimports
```

### sublime插件

gosublime
```
{
	"env":{
		"GOROOT":"/opt/go",
		"GOPATH":"/home/chen/git/go"
	},
	"fmt_cmd": ["goimports"]
}
```
godef
```
[
    { "keys": ["ctrl+d"], "command": "godef" },
    { "keys": ["ctrl+f"], "command": "jump_forward"},
    { "keys": ["ctrl+g"], "command": "jump_back"}
]
```
