# tools
## print struct
- [github.com/davecgh/go-spew， **推荐**](github.com/davecgh/go-spew)

	```go
	package main

	import (
	    "github.com/davecgh/go-spew/spew"
	)

	type Project struct {
	    Id      int64  `json:"project_id"`
	    Title   string `json:"title"`
	    Name    string `json:"name"`
	    Data    string `json:"data"`
	    Commits string `json:"commits"`
	}

	func main() {
	    o := Project{Name: "hello", Title: "world"}
	    spew.Dump(o)
	}
	```
- [Go数据结构完美打印](https://github.com/shivamMg/ppds)

## memmap 可视化数据结构工具
```go
package main

import (
	"bytes"
	"io/ioutil"

	"github.com/bradleyjkemp/memviz"
)

func main() {
	type T struct {
		Id     int
		Name   string
		Parent *T
	}

	var t, tP T

	tP.Id = 0
	tP.Name = "0"

	t.Id = 1
	t.Name = "1"
	t.Parent = &tP

	buf := &bytes.Buffer{}
	memviz.Map(buf, &t)
	ioutil.WriteFile("a.dot", buf.Bytes(), 0600)
}
```

生成图片:
```
dot -Tpng -o a.png a.dot
```

### gpt
- [github.com/diskfs/go-diskfs](https://github.com/diskfs/go-diskfs)
- [github.com/rekby/gpt](https://github.com/rekby/gpt)

## 压缩binary
- [使用 upx 压缩 go build 打包的可执行文件](https://abelsu7.top/2019/10/24/go-build-compress-using-upx/)

## 代码混淆
- [go代码混淆 - gobfuscate](https://www.bcskill.com/index.php/archives/1000.html)

## 调优
- [Go调优神器trace介绍](https://studygolang.com/articles/9693)

## IDE
### FAQ
#### vscode 智能提示突然消失
vscode的go插件启用了gopls, 但gopls总是崩溃.

解决方法, 关闭gopls:
```
"go.useLanguageServer": false
```

#### vscode 插件Go 安装工具
通过`ctrl+shift+p`打开面板, 输入`go:install/Update Tools`，回车后，选择所有插件(勾一下全选)，点击确认，进行安装即可

配置说明:
```json
"go.autocompleteUnimportedPackages": true, // 自动完成未导入的包
"go.inferGopath": true, // 如果遇到使用标准包可以出现代码提示，但是使用自己的包或者第三方库无法出现代码提示，可以查看一下该配置项
"http.proxy": "https://192.168.0.233:3142", // VSCode 的一些插件需要配置代理，才能够正常安装
//  "http.proxy": "192.168.0.233:3142",
"http.proxyStrictSSL": false,
"go.docsTool": "gogetdoc", // 如果引用的包使用了 ( . "aa.com/text") 那这个text包下的函数也无法跳转进去, 可将 "go.docsTool" 改为 gogetdoc，默认是 godoc
```

#### [vscode go依赖更新](https://github.com/golang/vscode-go/blob/master/docs/commands.md#go-installupdate-tools)
打开golang 项目, 右键选择`Go: Show All Commands...`(或Ctrl+Shift+P) -> 输入`Go: Install/Update Tools`, 在下拉中选中该命令 -> 选择全部插件, 点击输入框右侧的"ok"按钮即可.