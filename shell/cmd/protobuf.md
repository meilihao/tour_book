# protobuf
## FAQ
### probuf定义`int64 timestamp`, http api调用时返回的json中timestamp是字符串类型
protobuf 的 int64 在 JSON 序列化时默认会被转成字符串，这是为了兼容 JavaScript 数字精度问题（JS 的 Number 无法安全表示所有 64 位整数）

解决方法:
1. 改用 int32
2. 使用 Timestamp 类型