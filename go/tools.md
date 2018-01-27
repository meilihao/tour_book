# tools

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