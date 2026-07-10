# protobuf
ref:
- [Modern Protobuf and gRPC](https://buf.build/)

    - [现代化的 Protobuf 构建工具 buf](https://blog.cong.moe/post/2022-05-18-buf-tool/)

## FAQ
### probuf定义`int64 timestamp`, http api调用时返回的json中timestamp是字符串类型
protobuf 的 int64 在 JSON 序列化时默认会被转成字符串，这是为了兼容 JavaScript 数字精度问题（JS 的 Number 无法安全表示所有 64 位整数）

解决方法:
1. 改用 int32
2. 使用 Timestamp 类型

### 在 protobuf 中传输 Go 的 json.RawMessage
方案	可读性	类型安全	性能	适用场景
bytes	❌ base64	❌ 手动解析	高	二进制/加密数据
string	✅ 原始 JSON	❌ 手动解析	高	**推荐**，灵活存储
Struct	✅ 嵌套对象	✅ 类型安全	中	动态结构，需访问字段
