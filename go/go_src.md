# 源码理解

### struct中的接口理解

location:`database/sql/driver.Null`

```go
package main

import (
	"fmt"
	"strconv"
)

type Value interface{}

type ValueConverter interface {
	ConvertValue(v interface{}) (Value, error)
}

type Null struct {
	Converter ValueConverter
}

func (n Null) ConvertValue(v interface{}) (Value, error) {
	if v == nil {
		return nil, nil
	}
	// n.Converter.ConvertValue(v)==(n.Converter).ConvertValue(v)
	// 此例中,n.Converter==Bool为真 (同时比较接口中存储的值和类型)
	return n.Converter.ConvertValue(v)
}

type boolType struct{}

func (boolType) ConvertValue(src interface{}) (Value, error) {
	switch s := src.(type) {
	case bool:
		return s, nil
	case string:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("sql/driver: couldn't convert %q into type bool", s)
		}
		return b, nil
	}

	return nil, fmt.Errorf("sql/driver: couldn't convert %v (%T) into type bool", src, src)
}

func main() {
	var Bool boolType
	var n Null = Null{Converter: Bool}
	r, e := n.ConvertValue("true")
	fmt.Println(r, e)
	// 上面等价于下面
	r1, e1 := Bool.ConvertValue("true")
	fmt.Println(r1, e1)
}
```

## database 插入或查询
- 字段类型需要支持查询的scan时，需要实现sql.Scanner接口
- 字段类型需要支持插入时，需要实现driver.Valuer接口
