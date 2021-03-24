# protoc
为proto文件生成代码的工具

## `File does not reside within any path specified using --proto_path (or -I)`
```proto3
syntax = "proto3";

package check;

import "google/protobuf/timestamp.proto";
import "gogoproto/gogo.proto";

message checkReq {
    string book = 1 [(gogoproto.jsontag) = "MyField2", (gogoproto.moretags) = 'customtag:"v"'];
    google.protobuf.Timestamp created_at = 2 [(gogoproto.stdtime) = true];
}

message checkResp {
    bool found = 1;
    int64 price = 2;
}

service checker {
    rpc check(checkReq) returns(checkResp);
}
```

`protoc -I=/home/chen/tmpfs/protobuf -I=/home/chen/tmpfs/protobuf/protobuf --gofast_out=plugins=grpc:. my.proto`报标题中的错误.


一旦使用了`-I`参数则必须将my.proto所在的path也要用`-I`指定, 比如这里是`-I=.`. 注意`-I`不支持`~`路径形式但支持`$GOPATH`.

## 将google.protobuf.Timestamp生成为gogo的`*types.Timestamp`
`protoc -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf  -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2 -I=. --gofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,plugins=grpc:. *.proto`

## protobuf 因为 XXX 插入数据库报错的情况
```go
// https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L2487
// generateInternalStructFields just adds the XXX_<something> fields to the message struct.
func (g *Generator) generateInternalStructFields(mc *msgCtx, topLevelFields []topLevelField) {
	if gogoproto.HasUnkeyed(g.file.FileDescriptorProto, mc.message.DescriptorProto) {
		g.P("XXX_NoUnkeyedLiteral\tstruct{} `json:\"-\"`") // prevent unkeyed struct literals
	}
	if len(mc.message.ExtensionRange) > 0 {
		if gogoproto.HasExtensionsMap(g.file.FileDescriptorProto, mc.message.DescriptorProto) {
			messageset := ""
			if opts := mc.message.Options; opts != nil && opts.GetMessageSetWireFormat() {
				messageset = "protobuf_messageset:\"1\" "
			}
			g.P(g.Pkg["proto"], ".XXX_InternalExtensions `", messageset, "json:\"-\"`")
		} else {
			g.P("XXX_extensions\t\t[]byte `protobuf:\"bytes,0,opt\" json:\"-\"`")
		}
	}
	if gogoproto.HasUnrecognized(g.file.FileDescriptorProto, mc.message.DescriptorProto) {
		g.P("XXX_unrecognized\t[]byte `json:\"-\"`")
	}
	if gogoproto.HasSizecache(g.file.FileDescriptorProto, mc.message.DescriptorProto) {
		g.P("XXX_sizecache\tint32 `json:\"-\"`")
	}
}
```

在`XXX_...`后追加`xorm:"-" or gorm:"-"`即可.