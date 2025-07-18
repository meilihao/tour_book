# 1Panel
## 部署
ref:
- [在线安装](https://1panel.cn/docs/installation/online_installation/)
- [installer](https://github.com/1Panel-dev/installer)

通过quick_start.sh获悉具体安装脚本是1panel-v2.0.4-linux-amd64.tar.gz的install.sh

## 开发
ref:
- [开发环境搭建](https://1panel.cn/docs/dev_manual/dev_manual/)

代码入口:
- web : `cmd/server/cmd/root.go`的`server.Start`

	```go
	...
	rootRouter := router.Routers() // 构建routes
	...
	```

## FAQ
### 1pctl源码
1pctl是对1panel的脚本封装

### tls cert操作
在`backend/app/service/website_*.go`, 比如:
- 创建ca: `(WebsiteCAService) Create()`