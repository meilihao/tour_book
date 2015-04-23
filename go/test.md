### FAQ

#### `go test`时提示"undefined: XXX"

将其测试的源码文件也作为参数即可

```shell
# 举例, cd golang源码目录/src/database/sql/driver
$ go test types_test.go  # 将提示:"undefined: ValueConverter"
$ go test types_test.go  types.go driver.go # ok
```