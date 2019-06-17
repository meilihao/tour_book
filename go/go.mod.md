# go mod
Go1.11和Go1.12只能在$GOPATH以外的目录中使用Go Modules.

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
从 Go 1.11 版本开始，官方不仅支持了 go module 包依赖管理工具, 还新增了 GOPROXY 环境变量. 如果设置了该变量，下载源代码时将会通过这个环境变量设置的代理地址，而不再像以前那样直接从代码库下载.

不过，需要依赖于 go module, 可通过`export GO111MODULE=on`开启 MODULE.

更可喜的是，https://goproxy.io 这个开源项目帮我们实现好了我们想要的, 该项目允许开发者一键构建自己的 GOPROXY 代理服务. 它同时也提供了公用的代理服务 https://goproxy.io，我们只需设置该环境变量即可：
```sh
export GOPROXY=https://goproxy.io
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