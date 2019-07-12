# dlv

## FAQ
### [给被调试程序传入参数](https://github.com/go-delve/delve/issues/562)
```
dlv debug main.go -- -conf=conf.yaml
```

### 打断点
```
(dlv) b /home/chen/git/go/src/golang.org/x/tools/godoc/server.go:276
(dlv) b handlerServer.ServeHTTP // func (h *handlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

### 函数调用中使用函数作为参数的调试方法
```go
info := h.GetPageInfo(abspath, relpath, mode, r.FormValue("GOOS"), r.FormValue("GOARCH"))
```

步骤:
1. 在`h.GetPageInfo`所在行设置断点, 再触发该断点
1. 输入`s`, dlv会进入`r.FormValue("GOOS")`, 一直回车, 待`r.FormValue`执行完成后dlv会回到`h.GetPageInfo`
1. 输入`s`, `r.FormValue("GOARCH")`与`r.FormValue("GOOS")`类似, 执行完成后dlv会回到`h.GetPageInfo`
1. 在输入`s`, dlv进入`h.GetPageInfo`

### print 显示不全即显示more
```
(dlv) config max-string-len 99999 // 显示最大输出长度
(dlv) config -list // 查看dlv配置
```