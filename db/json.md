# json
参考:
- [MySQL常用Json函数](https://www.cnblogs.com/waterystone/p/5626098.html)
## mysql
```sql
select JSON_EXTRACT('{"a": 1, "b": 2, "c": {"d": 4}}','$.a') -- 操作对象
select JSON_EXTRACT('[{"id":1,"name":"1"},{"id":2,"name":"2"}]','$[*].id') -- 操作数组
```