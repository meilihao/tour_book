# kratos

## 编码
ref:
- [从 Kratos v2 迁移到 v3](https://github.com/go-kratos/kratos/blob/main/docs/migration/v2-to-v3_zh.md)

v2 中 encoding/json 同时处理普通 Go JSON 值和 proto.Message，其中 protobuf 消息会使用 protobuf JSON 语义

## FAQ
### proto定义`bytes data_raw = 6`("DataRaw  []byte `protobuf:"bytes,6,opt,name=data_raw,json=dataRaw,proto3" json:"data_raw,omitempty"`"), DataRaw为nil时, kratos返回的json是`"data_raw":""`
原因: 
```go
// mod/github.com/go-kratos/kratos/v2@v2.9.2/encoding/json/json.go
func (codec) Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m) // 走到了这里
	default:
		return json.Marshal(m)
	}
}
```

解决:
```go
import "github.com/go-kratos/kratos/v2/encoding/json"

json.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: false, // 不输出未设置的字段
	}
```