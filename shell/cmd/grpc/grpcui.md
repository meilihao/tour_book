# grpcui
grpc调试工具grpcui/grpcurl, 推荐grpcui

```bash
# --- go get github.com/fullstorydev/grpcui
# go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
# grpcui -plaintext localhost:57081
# --- https://github.com/fullstorydev/grpcurl
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
# grpcurl -plaintext localhost:57081 list # 获取服务列表
driverdevice.RpcDevice
# grpcurl -plaintext localhost:57081 list driverdevice.RpcDevice
driverdevice.RpcDevice.GetDeviceConnectStatus
# grpcurl -plaintext localhost:57081 describe driverdevice.RpcDevice.GetDeviceConnectStatus
driverdevice.RpcDevice.GetDeviceConnectStatus is a method:
rpc GetDeviceConnectStatus ( .driverdevice.GetDeviceConnectStatusRequest ) returns ( .driverdevice.GetDeviceConnectStatusResponse );
```

> 类似的调试工具还有[github.com/fullstorydev/grpcurl](https://chai2010.cn/advanced-go-programming-book/ch4-rpc/ch4-08-grpcurl.html)

## FAQ
### `grpcui -insecure localhost:9203`报`Failed to compute set of methods to expose: server does not support the reflection API`
添加将grpc.Server注册到反射服务中即可:
```go
import (
    "google.golang.org/grpc/reflection"
)

func main() {
    s := grpc.NewServer()
    pb.RegisterYourOwnServer(s, &server{})

    // Register reflection service on gRPC server.
    reflection.Register(s)

    s.Serve(lis)
}
```

如果启动了gprc反射服务，那么就可以通过reflection包提供的反射服务查询gRPC服务或调用gRPC方法.