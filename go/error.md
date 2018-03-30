### `mixture of field:value and value initializers`

```go
type Food struct {
	Id int
}

type Store struct {
	Id    int
	Foods []Food
}

func main() {
	// 给下面的[]Food添加上字段名"Foods"即可
	store := Store{Id: 0, []Food{Food{100}, Food{101}}}
	fmt.Println(store)
}
```

### `not enough arguments in call to %func%`

可能使用了静态类型去调用方法而不是其实例去调用方法

### `invalid character '{{X}}' looking for beginning of object key string

解析这样的json时碰到:
```json
[{Album: 1,Description: "cf"}]
```

json的key没有引号导致,更正为:
```json
[{"Album": 1,"Description": "cf"}]
```

### 自定义类型实现json.Marshaler
```go
type Filename struct {
	Md5           string
	Width, Height uint
	Ext           string
}
func (f Filename) MarshalJSON() ([]byte, error) {
	if f == zeroFilename {
                // 不能使用"return nil,nil",必须有值,否则报json: error calling MarshalJSON for type datatype.Filename: unexpected end of JSON input
		return []byte(`""`), nil
	}
        // return 字符串必须带双引号,否则报json: error calling MarshalJSON for type datatype.Filename: invalid character 'f' after top-level value
	return []byte(fmt.Sprintf(`"%s-%dx%d.%s"`, f.Md5, f.Width, f.Height, f.Ext)), nil
}
```

### sync.WaitGroup.Wait 报"... deadlock"
sync.WaitGroup应该使用指针传递而不是值传递.


### 类型判断
```go
type Student struct {
	Name string
}

func zhoujielun(v interface{}) {
	switch msg := v.(type) {
	case *Student, Student:
		fmt.Printf("%+v\n", msg.Name) // error, 有两个分支, 编译器无法知道具体类型,只能放入万能的interface{},因此这里的msg type是interface{}
	}
}
```

```go
type Student struct {
	Name string
}

func zhoujielun(v interface{}) {
	switch msg := v.(type) {
	case *Student:
		fmt.Printf("%+v\n", msg.Name) // ok, 只有一个分支, 编译器知道具体类型
	}
}
```

### map不可寻址
```go
type Student struct {
	name string
}

func main() {
	m := map[string]Student{"people": {"zhoujielun"}}
	m["people"].name = "wuyanzu" // error
}
```
map[string]Student 中储存的Student是数值拷贝，当要修改Student里面的Name时，就发生了不可寻址的错误.

更直接的理解就是:
`m["people"].name`的意图是修改元素的属性, 而`m["people"]`获取的元素的copy, 意图和实际不符.