
# 构建

## 条件编译

### 编译标签(`build tag`)：

```
// +build { GOOS }, { GOOS }, { !GOOS }
```
```
// +build (linux AND 386) OR (darwin AND (NOT cgo))
```
```
// +build linux darwin
// +build 386
```

1. 以`// +build`开头,一个源文件里可以有多个编译标签，多个编译标签之间是逻辑"与"的关系
1. 支持 GOOS 与 GOARCH 并可以具有多个值，用`,`分割， 例如：`// +build linux, darwin, freebsd`
1. 支持 不等条件 `!` ， 例如：`\\ +build !windows`,即不在windows环境下时，均可编译此文件。
1. 支持 与或非 逻辑， AND OR NOT
1. 条件编译需要前后空一行，否则无法识别

### 文件后缀：

```
xxx_{ GOOS }.go xxx_{ GOOS }_{ GOARCH }.go
```

1. 支持 GOOS ，例如： curl_windows.go
1. 支持 GOARCH， 例如： curl_386.go
1. 支持 上述两种叠加，但不可调换顺序 xxx_{ GOOS }_{ GOARCH }.go ，例如： curl_windows_amd64.go

### 如何选择

这两者可以叠加使用，但注意不要出现冗余，如：curl_windows.go 里面写`// +build windows`则重复了.
如果编译的文件是一一对应关系的话，使用文件后缀更简单些，如对每个 GOOS 生成一个文件.
如果有复杂条件的话，可以使用标签编译方式。如：
```
curl_windows.go 对应 windows 平台。
curl_others.go 里面写 \\ +build !windows 对应 非windows 平台
```

## 交叉编译

Golang 1.5及以上修改`GOOS和GOARCH`即可.如：

```
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test
$ file test
```

[GOOS 与 GOARCH 支持的参数](https://golang.org/doc/install/source#environment)