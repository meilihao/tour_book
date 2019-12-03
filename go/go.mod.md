# go mod
Go1.11和Go1.12只能在$GOPATH以外的目录中使用Go Modules.

> 目前go mod的包及相关数据均缓存在 $GOPATH/pkg/mod和 $GOPATH/pkg/sum 下，未来或将移至 $GOCACHE/mod 和$GOCACHE/sum 下(可能会在当 $GOPATH 被淘汰后)
> 可以使用`go clean -modcache`清理所有已缓存的go mod数据
> go mod导致`$GOPATH`的弱化甚至淘汰及`$GOBIN(go install)`的提升

## mod升级
```
$ go get github.com/objcoding/testmod@v1.0.1
```
或
```
$ go mod edit -require="github.com/objcoding/testmod@v1.0.1" // 主动修改 go.md 文件中依赖的版本号
$ go mod tidy // 对版本进行更新，这是一条神一样的命令，它会自动清理掉不需要的依赖项，同时可以将依赖项更新到当前版本
```

## replace
从 Go 1.11 版本开始，新增支持了 go modules 用于解决包依赖管理问题. 该工具提供了 replace，就是为了解决包的别名问题，也能替我们解决 golang.org/x 无法下载的的问题.

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

如果在运行go mod vendor时，提示`Get https://sum.golang.org/lookup/xxxxxx: dial tcp 216.58.200.49:443: i/o timeout`，则是因为Go 1.13设置了默认的GOSUMDB=sum.golang.org用于验证包的有效性，而这个网站是被墙了, 可以通过命令关闭：`go env -w GOSUMDB=off`.


## GOPRIVATE
控制哪些私有仓库和依赖(公司内部仓库)不通过 proxy 来拉取，直接走本地

```sh
# 设置不走 proxy 的私有仓库，多个用逗号相隔
go env -w GOPRIVATE=*.example.com
```

## FAQ
### go: cannot determine module path for source directory
在 $GOPATH 之外使用 go modules, 如果是现有项目的话可以直接 go mod init, 现有项目会根据 git remote 自动识别 module 名, 但是新项目的话就会报`go: cannot determine module path for source directory`, 此时需要带上 module 名即可.

### go mod使用gitlab私有仓库作为项目的依赖包
```sh
$ git config --global url."git@code.aliyun.com:xxx_backend/saas.git".insteadOf "https://code.aliyun.com/xxx_backend/saas.git"
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

更新时直接使用`go get -u code.aliyun.com/xxx_backend/saas`会报错, 可将replace里的版本, 比如这里的`v0.0.0-20190617102944-e1b0da75851a`替换为`latest`, `go build`时会自动更新到最新版.

注意, go1.13运行`go build`时要使用:
```bash
$ env GONOPROXY="code.aliyun.com" GONOSUMDB="code.aliyun.com" go build
```

> GONOPROXY,GONOSUMDB有多项时需用`,`分隔

### athens deploy
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