# dlv
[delve](https://github.com/go-delve/delve)紧跟Go语言版本演进, 是目前Go调试器的事实标准.

正常go build/install出的go程序是完全优化过的, 强行使用调试器挂接调试时, 某些local变量/lamda表达式捕获的变量会直接进入寄存器, 无法使用调试器查看, 可在build或install时加入关闭编译器优化的参数`-gcflags "-N -l"`.

> dlv支持调试go test

> gccgo是GCC编译器的前端, 比原生编译器慢, 同时支持的go版本也更低, 但具有更为强大的优化能力和支持更多的处理器架构, 也支持gdb调试

## FAQ
### terminal绑定给当前程序, 导致dlv exec时无法输入调试命令
```bash
dlv exec ./my-app --headless --listen=:2345 --api-version=2 --accept-multiclient
dlv connect :2345
```

### [给被调试程序传入参数](https://github.com/go-delve/delve/issues/562)
```
dlv debug main.go -- -conf=conf.yaml
sudo dlv exec ./zrepl --  daemon --config=/home/jr/git/go/src/github.com/zrepl/zrepl/config/samples/push.yml
sudo dlv debug github.com/zrepl/zrepl --  daemon --config=/home/jr/git/go/src/github.com/zrepl/zrepl/config/samples/push.yml
```
> [`--`](https://stackoverflow.com/questions/39779238/using-delve-with-subcommand-and-flags)用于分隔dlv和自身的参数.

### 命令
```
(dlv) b 93 # 当前文件的93行
(dlv) b /home/chen/git/go/src/golang.org/x/tools/godoc/server.go:276
(dlv) b handlerServer.ServeHTTP // func (h *handlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
(dlv) b (*Statement).AddVar:+6 // (*Statement).AddVar()的第六行
(dlv) clear 1
(dlv) bp
(dlv) list
(dlv) c : continue
(dlv) n : next
(dlv) s : step
(dlv) whatis a : 查看类型
(dlv) p a : 查看内容
(dlv) set a=4
(dlv) call foo.Foo(a, b)
(dlv) regs : 查看寄存器
(dlv) locals : 当前函数栈本地变量列表(包括变量的值)
(dlv) args : 当前函数栈参数和返回值列表(包括参数和返回值的值)
(dlv) examinemem(简写为x): 查看某一内存地址上的值
(dlv) bt : 查看函数调用栈信息
(dlv) up/down : 在函数调用栈的栈帧间进行跳转
(dlv) frame <N> : 栈帧跳转
(dlv) r : restart
(dlv) goroutines : `*`表示当前goroutine
(dlv) goroutine <N> : 切换goroutine
(dlv) threads
```

> 可用`sources`查看用到的源码文件, 便于打断点

### 调试coredump
```go
$ ulimit -c unlimited // 不限制core文件大小
$ go build main.go
$ GOTRACEBACK=crash ./main
$ dlv core ./main ./core
```

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

### `could not find statement at`
Go编译器对目标代码做了优化, 关闭优化即可`-gcflags=all="-N -l"`