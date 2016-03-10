## map

### 不能赋值,提示"cannot assign to ..."
```go
type person struct {
       id   int
       name string
}
persons := make(map[int]person)
persons[0] = person{1, "a"}
persons[0].id=2 // 报错
```
map里面的元素是值拷贝,当前存入的person是不可寻址的,所以不能用persons[0].id这种方法赋值,解决方法:
1. 把element类型换成`*person`
2. 重新赋值
```go
p := person[0]
p.id = 2
person[0] = p
```
