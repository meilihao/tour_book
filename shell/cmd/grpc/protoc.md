# protoc
为proto文件生成代码的工具

## [grpc](https://www.grpc.io/docs/languages/go/quickstart/)
```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## FAQ
### gogo 插件区别
gofast: 速度优先，但此方式不支持其它 gogoprotobuf 的扩展选项.

gogofast、gogofaster、gogoslick: 更快的速度、会生成更多的代码:
- gogofast类似gofast，但是会引入 gogoprotobuf 库
- gogofaster类似gogofast，但是不会产生XXX_unrecognized类的指针字段，可以减少垃圾回收时间
- gogoslick类似gogofaster，但是会增加一些额外的string、gostring和equal method等

### 不修改代码注入tag
[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)

### 生成pb.go时报错: `File does not reside within any path specified using --proto_path (or -I)`
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

### protobuf 因为 XXX 插入数据库报错的情况
ref:
- [golang 解决 protobuf 因为 XXX 插入数据库报错的情况](https://www.hwholiday.com/2020/go_protoc_gen_go/)

**新版protobuf-go, 比如`google.golang.org/protobuf@v1.28.0`**已用其他字段取代了`XXX_`

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

### 将google.protobuf.Timestamp生成为gogo的`*types.Timestamp`
`protoc -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf  -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2 -I=. --gofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,plugins=grpc:. *.proto`

### 为`google.golang.org/protobuf/types/known/timestamppb.Timestamp`添加 Scanner 和 Valuer 接口
先生成vendor, 再在`$GOPATH/pkg/mod/google.golang.org/protobuf@v1.28.0/types/known/timestamppb/`下添加gorm.go
```go
// from https://github.com/hashicorp/boundary/blob/main/internal/db/timestamp/scanners.go
package timestamppb

import (
	"database/sql/driver"
	"errors"
	"math"
	"time"
)

var (
	NegativeInfinityTS = time.Date(math.MinInt32, time.January, 1, 0, 0, 0, 0, time.UTC)
	PositiveInfinityTS = time.Date(math.MaxInt32, time.December, 31, 23, 59, 59, 1e9-1, time.UTC)
)

// Scan implements sql.Scanner for protobuf Timestamp.
func (ts *Timestamp) Scan(value interface{}) error {
	switch t := value.(type) {
	case time.Time:
		ts = New(t) // google proto version
	case string:
		switch value {
		case "-infinity":
			ts = New(NegativeInfinityTS)
		case "infinity":
			ts = New(PositiveInfinityTS)
		}
	default:
		return errors.New("Not a protobuf Timestamp")
	}
	return nil
}

// Scan implements driver.Valuer for protobuf Timestamp.
func (ts *Timestamp) Value() (driver.Value, error) {
	if ts == nil {
		return nil, nil
	}
	return ts.AsTime(), nil
}

// GormDataType gorm common data type (required)
func (ts *Timestamp) GormDataType() string {
	return "timestamp"
}
```

直接在`$GOPATH/pkg/mod/google.golang.org/protobuf@v1.28.0/types/known/timestamppb/`下添加gorm.go不会被编译; 将gorm.go代码直接加入timestamp.pb.go, 编译时会报导入std pkg失败. 它们失败的原因未知, 可能是go mod的限制.