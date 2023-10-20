# jq
将stdin的json文本格式化后输出.

```
# apt install jq
```

## 选项
- -r : 删除双引号

## example
```bash
$ jq -r .arch version.json
$  curl http://localhost:5984/ussfed_tmm -s | jq  # `-s`作用: 忽略curl的process header, 否则jq的输出内容带该信息
```

## FAQ
### [jq当项不存在时如何使用空串以代替null](https://www.jianshu.com/p/e7c8efe17d9d)
- 用替换操作运算符(Alternative operator): `jq -r '.c // ""' $JSON`
- 用if-then-else: `jq -r 'if .a == null then "" else .a end' $JSON` 或 `jq -r '.a | if . == null then "" else . end' $JSON`