# grpc

参考:
- [gRPC学习笔记](https://github.com/skyao/leaning-grpc)
- [proto3学习笔记和资料翻译](github.com/skyao/leaning-proto3)
- [gRPC 官方文档中文版](http://doc.oschina.net/grpc)
- [Protobuf3 语法指南](http://colobu.com/2017/03/16/Protobuf3-language-guide/)
- [Go使用gRPC与Protocol Buffers构建高性能API服务](http://www.tuicool.com/articles/ve6zYnf)
- [etcd: distributed key-value store with grpc/http2](https://blog.gopheracademy.com/advent-2015/etcd-distributed-key-value-store-with-grpc-http2/)
- [Protobuf 终极教程](https://colobu.com/2019/10/03/protobuf-ultimate-tutorial-in-go/)

## 搭建环境
ref:
- [Quick start](https://grpc.io/docs/languages/go/quickstart/)
- [Protocol Buffer Basics: Go](https://protobuf.dev/getting-started/gotutorial/)

1. 因grpc使用了ProtoBuffer作为IDL,根据[文档](https://github.com/golang/protobuf),需先安装[Protobuf编译器 c++版](https://github.com/google/protobuf/releases).

    **`github.com/golang/protobuf`已被`google.golang.org/protobuf(即https://github.com/protocolbuffers/protobuf-go)`取代**
2. 安装protoc的go语言插件:
    - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
    - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

### protoc

```
protoc --proto_path=IMPORT_PATH --go_out=plugins=grpc:OUT_DIR path/to/file.proto
```

- --proto_path=IMPORT_PATH IMPORT_PATH 指定一个目录用于查找`.proto文件`,缺省使用当前目录`.`
- --go_out 参数告诉 protoc 使用 protoc-gen-go 插件来生成 Go 代码, 主要用于生成 Protocol Buffers 消息的 Go 结构体代码. 只定义了数据结构（消息），而没有定义 gRPC 服务时，或者需要纯粹的 Protocol Buffers 消息操作时，使用该参数
- --go-grpc_out 参数告诉 protoc 使用 protoc-gen-go-grpc 插件来生成 Go 代码, 专门用于生成 gRPC 服务接口的 Go 代码. 比如定义了 gRPC 服务（service 关键字），并且需要为该服务生成客户端和服务端代码时，使用该参数.
- path/to/file.proto 可以提供一个或多个`.proto`文件作为输入. 多个`.proto`文件可以一次指定. 虽然文件被以当前目录的相对路径命名, 每个文件必须位于一个IMPORT_PATH路径下, 以便编译器可以检测到它的标准名字.

> `--go_out=plugins=grpc:OUT_DIR` 是在 Go Protobuf 和 gRPC 生态系统早期使用(`在旧版本的 protoc-gen-go v1.2 之前，对应 github.com/golang/protobuf 模块`)的一种方法，它允许一个插件同时处理消息和服务代码的生成. 但在当前的 Go gRPC 开发中，这种方式已被更模块化和清晰的两步生成方式取代，即分别使用 `--go_out 和 --go-grpc_out`

更多见[protobuf的官方文档](https://github.com/protocolbuffers/protobuf-go).

### 示例

- [sikang99/grpc-example](https://github.com/sikang99/grpc-example)
- [gRPC 官方文档中文版: Go](http://doc.oschina.net/grpc?t=60133)

### tls

- [官方文档/示例(github)]
- [grpc加密TLS初体验](http://studygolang.com/articles/3220)
- [etcd使用](https://github.com/coreos/etcd/blob/340df26883b8604b574bef32e3b31a396cdb3bad/etcdserver/api/v3rpc/grpc.go)

## 错误

### error while loading shared libraries: libprotoc.so.xx

根据[文档](https://github.com/google/protobuf/blob/master/src/README.md),protobuf的默认安装位置是`/usr/local`,
而`/usr/local/lib`不在Ubuntu系统默认的`LD_LIBRARY_PATH`里,参照文档configure时使用`./configure --prefix=/usr`.

## 负载均衡和服务发现
ref:
- [gRPC客户端的那些事儿](https://tonybai.com/2021/09/17/those-things-about-grpc-client/)
- [写给go开发者的gRPC教程-服务发现与负载均衡](https://zhuanlan.zhihu.com/p/641743763)
- [Kratos 源码分析：Warden 负载均衡算法之 P2C](https://pandaychen.github.io/2020/07/25/KRATOS-WARDEN-BALANCER-P2C-ANALYSIS/)

## FAQ
### 让gin和grpc-gateway共用端口
```go
    r := gin.Default()
	grpcServer := grpc.NewServer()
	pb.RegisterHelloHTTPServer(grpcServer, HelloHTTPService)

	r.Use(func(ctx *gin.Context) {
		// 判断协议是否为http/2
		// 判断是否是grpc
		if ctx.Request.ProtoMajor == 2 &&
			strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			// 按grpc方式来请求
			ctx.Status(http.StatusOK)
			grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			// 不要再往下请求了,防止继续链式调用拦截器
			ctx.Abort()
			return
		}
		// 当作普通api
		ctx.Next()
	})
```