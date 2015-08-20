## install
go 环境变量配置
```bash
//etc/profile
#golang
export GOROOT=/opt/go
export GOPATH=/home/chen/git/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Go编码规范指南
参考: http://golanghome.com/post/550
### 格式化规范

go默认已经有了gofmt工具，但是我们强烈建议使用goimport工具，这个在gofmt的基础上增加了自动删除和引入包.

    go get golang.org/x/tools/cmd/goimports
不同的编辑器有不同的配置, sublime的配置教程：http://michaelwhatcott.com/gosublime-goimports/
LiteIDE默认已经支持了goimports，如果你的不支持请点击属性配置->golangfmt->勾选goimports
保存之前自动fmt你的代码。

### go vet

vet工具可以帮我们静态分析我们的源码存在的各种问题，例如多余的代码，提前return的逻辑，struct的tag是否符合标准等。

    go get golang.org/x/tools/cmd/vet
使用如下：

    go vet .

### import 规范

import在多行的情况下，goimports会自动帮你格式化，但是我们这里还是规范一下import的一些规范，如果你在一个文件里面引入了一个package，还是建议采用如下格式：
```go
import (
    "fmt"
)
```
如果你的包引入了三种类型的包，标准库包，程序内部包，第三方包，建议采用如下方式进行组织你的包：
```go
import (
    "encoding/json"
    "strings"

    "myproject/models"
    "myproject/controller"
    "myproject/utils"

    "github.com/astaxie/beego"
    "github.com/go-sql-driver/mysql"
)
```
有顺序的引入包，不同的类型采用空格分离，第一种实标准库，第二是项目包，第三是第三方包。

在项目中不要使用相对路径引入包：
```go
// 这是不好的导入
import "../net"

// 这是正确的做法
import "github.com/repo/proj/src/net"
```
### 注释规范

注释可以帮我们很好的完成文档的工作，写得好的注释可以方便我们以后的维护。详细的如何写注释可以
参考：http://golang.org/doc/effective_go.html#commentary