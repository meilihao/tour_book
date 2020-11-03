# dlv
正常go build/install出的go程序是完全优化过的, 强行使用调试器挂接调试时, 某些local变量/lamda表达式捕获的变量会直接进入寄存器, 无法使用调试器查看, 可在build或install时加入关闭编译器优化的参数`-gcflags "-N -l"`.

## FAQ
### [给被调试程序传入参数](https://github.com/go-delve/delve/issues/562)
```
dlv debug main.go -- -conf=conf.yaml
sudo dlv exec ./zrepl --  daemon --config=/home/jr/git/go/src/github.com/zrepl/zrepl/config/samples/push.yml
sudo dlv debug github.com/zrepl/zrepl --  daemon --config=/home/jr/git/go/src/github.com/zrepl/zrepl/config/samples/push.yml
```
> [`--`](https://stackoverflow.com/questions/39779238/using-delve-with-subcommand-and-flags)用于分隔dlv和自身的参数.

### 打断点
```
(dlv) b 93 # 当前文件的93行
(dlv) b /home/chen/git/go/src/golang.org/x/tools/godoc/server.go:276
(dlv) b handlerServer.ServeHTTP // func (h *handlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
(dlv) b (*Statement).AddVar:+6 // (*Statement).AddVar()的第六行
```

> 可用`sources`查看用到的源码文件, 便于打断点

### 远程调试
```go
# server
$ dlv attach $PID --headless --listen=:8181 // `dlv connect`退出后`dlv attach`才允许退出
# clinet
$ dlv connect 192.168.0.139:8181
```

**此时要用远程程序编译时的文件路径打断点而非本地路径**.

因为远程程序的编译路径可能和本地不一致, 就需要配置`~/.config/dlv/config.yml`, 比如:
```yaml
substitute-path:
  - {from: /root, to: /home/ubuntu/git} # 远程程序编译路径的前缀`/root`会被替换成本地的`/home/ubuntu/git`
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
