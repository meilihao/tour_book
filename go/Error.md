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