### `mixture of field:value and value initializers`

```go
type Food struct {
	Id     int
}

type Store struct {
    Id int
	Foods []Food
}

func main() {
    #给下面的[]Food添加上字段名"Foods"即可
	store := Store{Id:0,[]Food{Food{100}, Food{101}}}
	fmt.Println(store)
}
```

### `not enough arguments in call to %func%`

可能使用了静态类型去调用方法而不是其声明的变量去调用方法