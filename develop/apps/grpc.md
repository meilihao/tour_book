# grpc
ref:
- [史上最细gRPC(Go)入门教程(三)---gRPC流式传输--Streaming](https://blog.csdn.net/java_1996/article/details/113428990)

gRPC 中的 Service API 有如下4种类型:
1. UnaryAPI：普通一元方法
2. ClientStreaming：客户端推送流
3. ServerStreaming：服务端推送流
4. BidirectionalStreaming：双向推送流

Unary API 就是普通的 RPC 调用.

Stream 顾名思义就是一种流，可以源源不断的推送数据，很适合大数据传输，或者服务端和客户端长时间数据交互的场景. Stream API 和 Unary API 相比，因为省掉了中间每次建立连接的花费，所以效率上会提升一些.

demo:
```proto3
syntax = "proto3";

package streaming;

message StreamingRequest {
    string message = 1;
}

message StreamingResponse {
    string message = 1;
}

service Streaming {
    rpc Unary(StreamingRequest) returns (StreamingResponse) {}
    rpc ClientStreaming(stream StreamingRequest) returns (StreamingResponse) {}
    rpc ServerStreaming(StreamingRequest) returns (stream StreamingResponse) {}
    rpc BidirectionalStreaming(stream StreamingRequest) returns (stream StreamingResponse) {}
}
```