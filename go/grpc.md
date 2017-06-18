# grpc

参考:
- [gRPC学习笔记](https://github.com/skyao/leaning-grpc)
- [proto3学习笔记和资料翻译](github.com/skyao/leaning-proto3)
- [gRPC 官方文档中文版](http://doc.oschina.net/grpc)
- [Protobuf3 语法指南](http://colobu.com/2017/03/16/Protobuf3-language-guide/)
- [Go使用gRPC与Protocol Buffers构建高性能API服务](http://www.tuicool.com/articles/ve6zYnf)
- [etcd: distributed key-value store with grpc/http2](https://blog.gopheracademy.com/advent-2015/etcd-distributed-key-value-store-with-grpc-http2/)

## 搭建环境

1. 因grpc使用了ProtoBuffer作为IDL,根据[文档](https://github.com/golang/protobuf),需先安装[Protobuf编译器 c++版](https://github.com/google/protobuf/releases).
2. 安装protoc的go语言插件,`go get -a github.com/golang/protobuf/protoc-gen-go`.

### protoc

```
protoc --proto_path=IMPORT_PATH --go_out=plugins=grpc:OUT_DIR path/to/file.proto
```

- --proto_path=IMPORT_PATH IMPORT_PATH 指定一个目录用于查找`.proto文件`,缺省使用当前目录`.`
- --go_out=plugins=grpc:OUT_DIR `plugins=grpc`表示使用grpc插件,多个插件用`,`分隔;OUT_DIR用于存放生成的go代码在,和插件间用`:`分隔.
- path/to/file.proto 可以提供一个或多个`.proto`文件作为输入. 多个`.proto`文件可以一次指定. 虽然文件被以当前目录的相对路径命名, 每个文件必须位于一个IMPORT_PATH路径下, 以便编译器可以检测到它的标准名字.

更多见[protobuf的官方文档](https://github.com/golang/protobuf).

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
