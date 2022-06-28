# go mod
为了确保一致性构建，Go 引入了 go.mod 文件来标记每个依赖包的版本，在构建过程中 go 命令会下载 go.mod 中的依赖包，下载的依赖包会缓存在本地，以便下次构建.

Go1.11和Go1.12只能在$GOPATH以外的目录中使用Go Modules.

> 目前go mod的包及相关数据均缓存在 $GOPATH/pkg/mod和 $GOPATH/pkg/sum 下，未来或将移至 $GOCACHE/mod 和$GOCACHE/sum 下(可能会在当 $GOPATH 被淘汰后)
> 可以使用`go clean -modcache`清理所有已缓存的go mod数据
> go mod导致`$GOPATH`的弱化甚至淘汰及`$GOBIN(go install)`的提升

**注意:** go get 在非go mod项目下运行默认不使用GOPROXY, 除非显式指定GO111MODULE=on

## mod升级
```
$ go get github.com/objcoding/testmod@v1.0.1
$ go get -u github.com/satori/go.uuid@master
```
或
```
$ go mod edit -require="github.com/objcoding/testmod@v1.0.1" // 主动修改 go.md 文件中依赖的版本号
$ go mod tidy // 对版本进行更新，这是一条神一样的命令，它会自动清理掉不需要的依赖项，同时可以将依赖项更新到当前版本
```

## replace
从 Go 1.11 版本开始，新增支持了 go modules 用于解决包依赖管理问题. 该工具提供了 replace，就是为了解决包的别名问题，也能替我们解决 golang.org/x 无法下载的的问题.
replace也支持有限的本地路径(以"/","./"或"../"开头的路径)

```
module kaiqi/saas

go 1.12

require (
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
    ...
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/net => github.com/golang/net v0.0.0-20190310074541-c10a0554eabf
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190509080239-a5b02f93d862
)
```

## GOPROXY
从 Go 1.11 版本开始，官方不仅支持了 go module 包依赖管理工具, 还新增了 GOPROXY 环境变量. 如果设置了该变量，下载源代码时将会**优先**通过这个环境变量设置的代理地址，而不再像以前那样直接从代码库下载.

不过，需要依赖于 go module, 可通过`export GO111MODULE=on`开启 MODULE.

更可喜的是，[https://goproxy.cn](https://github.com/goproxy/goproxy.cn/blob/master/README.zh-CN.md) 这个开源项目帮我们实现好了我们想要的, 该项目允许开发者一键构建自己的 GOPROXY 代理服务. 它同时也提供了公用的代理服务 https://goproxy.cn，我们只需设置该环境变量即可：
```sh
export GOPROXY=https://goproxy.cn,direct
```

> direct 可以在一定程度上解决私有库的问题. 这个 GOPROXY 设定的工作原理是：当 go 在抓取目标模块时，若遇见了 404 错误，那么就回退到 direct 也就是直接去目标模块的源头（即指定的url处） 去抓取.

如果在运行go mod vendor时，提示`Get https://sum.golang.org/lookup/xxxxxx: dial tcp 216.58.200.49:443: i/o timeout`，则是因为Go 1.13设置了默认的GOSUMDB=sum.golang.org用于验证包的有效性，而这个网站是被墙了, 可以通过命令关闭：`go env -w GOSUMDB=off`.


## GOPRIVATE
控制哪些私有仓库和依赖(公司内部仓库)不通过 proxy 来拉取，直接走本地, 且跳过GOSUMDB的检查.

```sh
# 设置不走 proxy 的私有仓库，多个用逗号相隔
go env -w GOPRIVATE=*.example.com
```

## go.sum
考虑到go mod下载的依赖包有可能是被黑客恶意篡改的，以及缓存在本地的依赖包也有被篡改的可能，单单一个 go.mod 文件并不能保证一致性构建. 为了解决 Go module 的这一安全隐患，Go 开发团队在引入 go.mod 的同时也引入了 go.sum 文件，用于记录每个依赖包的哈希值，在构建时，如果本地的依赖包 hash 值与 go.sum 文件中记录得不一致，则会拒绝构建.

go.sum 文件中每行记录由 module 名、版本和哈希组成，并由空格分开：`<module> <version>[/go.mod] <hash>`, 比如
```txt
github.com/google/uuid v1.1.1 h1:Gkbcsh/GbpXz7lPftLA3P6TYMwjCLYm83jiFQZF/3gY=  
github.com/google/uuid v1.1.1/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

正常情况下，每个依赖包版本会包含两条记录，第一条记录为该依赖包版本整体（所有文件）的哈希值，第二条记录仅表示该依赖包版本中 go.mod 文件的哈希值，如果该依赖包版本没有 go.mod 文件，则只有第一条记录.

依赖包版本中任何一个文件（包括 go.mod）改动，都会改变其整体哈希值，此处再额外记录依赖包版本的 go.mod 文件主要用于计算依赖树时不必下载完整的依赖包版本，只根据 go.mod 即可计算依赖树.

> 每条记录中的哈希值前均有一个表示哈希算法的 h1:，表示后面的哈希值是由算法 SHA-256 计算出来的，自 Go module 从 v1.11 版本初次实验性引入，直至 v1.14 ，只有这一个算法.

go.sum 文件中记录的依赖包版本数量往往比 go.mod 文件中要多，这是因为二者记录的粒度不同导致的. go.mod 只需要记录直接依赖的依赖包版本，只在依赖包版本不包含 go.mod 文件时候才会记录间接依赖包版本，而 go.sum 则是要记录构建用到的所有依赖包版本.

## FAQ
### go: cannot determine module path for source directory
在 $GOPATH 之外使用 go modules时, 如果是现有项目的话可以直接 go mod init, 现有项目会根据 git remote 自动识别 module 名, 但是新项目的话就会报`go: cannot determine module path for source directory`, 此时需要带上 module 名即可.

### malformed module path "XXXX": missing dot in first path element
go1.13  mod 要求import 后面的path 第一个元素，符合域名规范，比如code.be.mingbai.com/tools/soa

即使是本项目下的其他包

如果无法使用域名，可以考虑使用replace+本地路径(`replace  code.be.mingbai.com/tools/soa  =>  ../../tools/soa`)，但不建议这样做.

### go mod使用gitlab私有仓库作为项目的依赖包
内网私有package难点:
1. 内网域名解析 : add dns server 或 /etc/hosts
1. go mod默认强制使用https : `go get -insecure`
1. GOPROXY和GOPRIVATE的配置, 特别是不能遗漏GOPROXY的direct
1. git repo 拉取时的权限问题

总结: 最佳实践是`replace + vendor`.

内网私有项目也可使用
- [goproxy, **推荐待测试**](https://mp.weixin.qq.com/s/Sxv5qb-v6OIhPptLWAHYUw)
- git submodule/subtree + go.mod replace(本地路径)来处理(有局限, 依赖多时不好管理, **不推荐**).

> go get是先通过https检查meta tag(含有vcs(版本控制系统) repo info)后再通过指定的vcs获取项目, 因此想通过`git config --global url."http://git.xxx.com".insteadOf "https://git.xxx.com"`以期待可直接使用`go get`是不能成功的.

> 为了安全, CI/CD会将私钥信息保存到其他地方, 在go mod处理依赖前引入即可, 保证镜像不包含私钥, 比如[jinkins + k8s会预先将私钥保存到k8s的secret中](https://medium.com/@ikolomiyets/go-modules-from-the-private-repository-in-the-pipeline-5921d3ba0e40).

**不使用replace的private repo(推荐, git repo不检查权限时用)**:
1. export GOPRIVATE="git.xxx.com" 
1. export GOPROXY="http://192.168.0.233:8801/repository/go-proxy-group/,direct" # [direct(不能省略)的作用](https://github.com/golang/go/issues/35861#issuecomment-558853990)
1. go get -insecure -u git.xxx.com/publicgomod/logrus-remote-hook # 在项目go.mod所在目录执行, 多个私有package就需要多次执行, 建议使用Makefile.
1. go build

**replace + vendor(推荐2)**
1. 先用go mod replace + 本地路径引入依赖并更新`.gitignore`
1. 再用go mod vendor生成vendor
1. 最后使用`go build -mod=vendor`即可.

> 建议: 引入的依赖的module name使用域名形式. 此时公网的第三方依赖可设置GOPROXY来来取.

> 该方法本质与`git submodule/subtree + go.mod replace`类似.

其他例子1:
```sh
$ git config --global url."git@code.aliyun.com:xxx_backend/saas.git".insteadOf "https://code.aliyun.com/xxx_backend/saas.git" # = git config --global url."git@code.aliyun.com:".insteadOf "https://code.aliyun.com/"
$ go get -u code.aliyun.com/xxx_backend/saas
```

go.mod:
```txt
...
require (
	...
	xxx/saas v0.0.0-00010101000000-000000000000 // go mod自动添加
)

// 有两种方法:
replace xxx/saas => code.aliyun.com/xxx_backend/saas v0.0.0-20190617102944-e1b0da75851a // 1. 使用私有仓库, 推荐
// replace xxx/saas => /home/chen/git/xxx/saas // 2. 使用本地package, 不推荐
```

更新时直接使用`go get -u code.aliyun.com/xxx_backend/saas`会报错, 可将replace里的版本, 比如这里的`v0.0.0-20190617102944-e1b0da75851a`替换为`latest`, `go build`时会自动替换成最新的`v0.0.0-20190617102944-e1b0da75851a`格式的版本.

注意, go1.13运行`go build`时要使用:
```bash
$ env GONOPROXY="code.aliyun.com" GONOSUMDB="code.aliyun.com" go build # 可先用`go build`试试, 不行再追加env GONOPROXY,GONOSUMDB
```

其他例子2:
1. go.mod追加`replace license_client => 192.168.0.226/OtherProject/License-Client latest`
1. git config --global url."git@192.168.0.226:OtherProject/License-Client.git".insteadOf "http://192.168.0.226/OtherProject/License-Client.git"
1. env GOPROXY="" go get -insecure -u 192.168.0.226/OtherProject/License-Client
1. `go build`

> 更新replace的version时要先将版本改为`latest`,在运行`go get`.

> GONOPROXY,GONOSUMDB有多项时需用`,`分隔

### athens deploy
推荐使用goproxy.io, 见下面的`代理`.

1. build
```sh
git clone https://github.com/gomods/athens
cd athens
make build-ver VERSION="0.7.0"
```

1. config
```
cd athens
cp config.dev.toml config.toml
touch FilterFile
echo "D" > FilterFile # 因为是直接重定向, 因此不用配置StorageType
# edit config.toml
FilterFile = "FilterFile"
GlobalEndpoint = "https://mirrors.aliyun.com/goproxy/"
```

1. run
```go
sudo ./athens -config_file config.toml
env GOPROXY=http://${athens_service_ip}:3000 go mod vendor # 使用
```

### Nexus Repository Manager 3 配置goproy需验证账号
Setting - Security - Anonymous, 启用匿名.

### $GOPROXY setting: cannot have comma
[使用go 1.13或以上版本](https://github.com/golang/go/issues/33725)
[在 master 分支的文档显示，只有在 GOPROXY 服务器返回 404 与 410 时，GOPROXY 才会使用逗号后面的下一个代理进行尝试](https://global.v2ex.com/t/566338#reply17)

### go mod vendor
vendor目录仅包含依赖到的代码, 未依赖的package会被忽略.

更新vendor时重新运行`go mod vendor`即可.

### 代理
- [Go Module 来了，企业私有代理你准备好了吗？](https://gocn.vip/topics/9831)
- [goproxy.cn - 为中国 Go 语言开发者量身打造的模块代理](https://segmentfault.com/a/1190000020293616)

### go mod项目嵌套
go mod项目嵌套时, go get需要指定具体的路径比如`go get -u golang.org/x/tools/gopls`,再将gopls拷入对应目录, 这里是拷入`$GOROOT/pkg/mod/golang.org/x/tools-v0.0.0-xxx-<xxx>`. 因为直接`go get -u golang.org/x/tools`时, tools目录不含gopls.

### 内网有goproxy, `go get xxx`报错
直接在启用go mod的项目中`import xxx`, 然后`go build`正常.

### 获取go mod支持的package versions
`GO111MODULE=on go list -m -json -versions go.etcd.io/etcd`

### /home/chen/git/go/pkg/mod/go.uber.org/zap@v1.16.0/level.go:24:2: missing go.sum entry for module providing package go.uber.org/atomic; try 'go mod tidy' to add it
需到go.mod所在目录执行`go mod tidy`

### go mod update管理
```bash
go get -u github.com/oligot/go-mod-upgrade
go-mod-upgrade
```
模块名称中的颜色有助于标识更新类型：

- 绿色进行较小的更新（minor update）
- 黄色，用于补丁更新（patch）
- 红色表示预发行更新（prerelease）

### go mod tidy报"module declares its path as: xxx\nbut was required as: xxx"
修正相关项目的go.mod命名后, 执行`go mod tidy`仍然报该错.

因为使用了GOPROXY, 其存在缓存的原因. 解决方法:
1. 先删除go.mod对应的rule
1. 执行`GOPROXY="direct" GOSUMDB="off" go get -u github.com/xxx/xxx`

### go install `xxx@latest`
参考:
- [Deprecation of 'go get' for installing executables](https://golang.google.cn/doc/go-get-install-deprecation)

`latest`是指主分支的lastest commit.

`go install github.com/meilihao/goi18n/v2@yaml`

### `go get github.com/meilihao/golib/v1`报`go get: github.com/meilihao/golib@v1.0.0: invalid version: go.mod has malformed module path "github.com/meilihao/golib/v1" at revision v1.0.0`
> env: go1.17
此时glib的go.mod使用`github.com/meilihao/golib/v1`, 将其替换为`github.com/meilihao/golib/v2`即可解决问题, 原因未知, 诡异.

### mod=`github.com/meilihao/golib/v2` file=`github.com/meilihao/golib/version.go` 无法使用ldflags传参
使用方法: `LDFLAGS="-X github.com/meilihao/golib/v2/golib.appName='xxx'`.

> 传参格式: `-X '{go mod名称}/包名.{变量名}'`

解决方法: 将version.go挪到version/version.go下, 使用`LDFLAGS="-X github.com/meilihao/golib/v2/version.appName='xxx'`即可.

### 更新go vendor
更新go.mod后重新执行`go mod tidy && go mod vendor`