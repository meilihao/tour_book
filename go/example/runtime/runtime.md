## Caller

`func Caller(skip int) (pc uintptr, file string, line int, ok bool)`

Caller报告当前go程调用栈所执行的函数的文件和行号信息。实参skip为上溯的栈帧数(1表示上一级调用者)，0表示Caller的调用者（Caller本身所在的调用栈）。（由于历史原因，skip的意思在Caller和Callers中并不相同。）函数的返回值为调用栈标识符、文件名、该调用在文件中的行号。如果无法获得信息，ok会被设为false。

```go
where := func() {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d", file, line)
}
where()
```