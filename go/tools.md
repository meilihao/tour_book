# tools
## 性能
- [pyroscope](https://colobu.com/2022/01/27/pyroscope-a-continuous-profiling-platform/)

## others
- [csv解析](https://github.com/gocarina/gocsv)
- [skip utf8 bom](https://github.com/dimchansky/utfbom)
- [test](https://github.com/stretchr/testify)

## http client
- [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)
- [github.com/guonaihong/gout](https://github.com/guonaihong/gout)

## error
- [errgroup]()

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
- [goleak : 一个可以事前检测 Go 泄漏的工具](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247494572&idx=1&sn=f6281cd182e7bfb7f20cd3641cb93306)

## script
- [在 Linux 中使用 Go 作为脚本语言](https://studygolang.com/articles/12461)

    ```bash
    # go get github.com/erning/gorun
    # echo ':golang:E::go::/usr/local/bin/gorun:OC' | sudo tee /proc/sys/fs/binfmt_misc/register
    # echo '-1' |sudo tee /proc/sys/fs/binfmt_misc/golang # [删除golang](https://android.googlesource.com/kernel/x86_64/+/android-5.0.0_r0.12/Documentation/binfmt_misc.txt)
    ```
- [github.com/bitfield/script](https://github.com/bitfield/script)

## crypto
- [go-dongle](github.com/golang-module/dongle)
- [go-cryptobin](https://github.com/deatil/go-cryptobin)

## encoding
### xlsx
- [github.com/xuri/excelize/v2](https://github.com/qax-os/excelize)

	- 当且仅有一个sheet时, 其无法删除
### pdf
- [github.com/johnfercher/maroto](https://github.com/johnfercher/maroto)

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

#### VSCode 运行go test显示打印日志 & 不使用test cache
`文件->首选项->设置->工作区设置->在setting.json中编辑`, 在settings节点下添加`"go.testFlags": ["-v", "-count=1"]`, 保存即可

> [`"go.testFlags": ["-v"]`](https://github.com/Microsoft/vscode-go/issues/1377)

### 粘包
- [FixedHeaderReceiveFilter](https://github.com/zboyco/go-server/blob/master/socket.go)

	实现bufio.SplitFunc

	go内置了4个splitFunc,当然也支持自定义:

	type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
	splitFUnc的功能就是:根据两个参数返回下一次Scan需要前进几个字节(advance)，分割出来的数据(token)，以及错误(err)
	参数:
		-- data(字节切片): 	 缓冲区的有效数据
		-- atEOF(bool):		是否已经输入完成(scanner是否扫描到源的结尾了),若没有则duplicate扩容!
	返回值:
		-- advance(int):	Scan需要前进几个字节(advance)
		-- token(字节切片):	 分割出来的数据
		-- err(error):		错误