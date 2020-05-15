# jq
将stdin的json文本格式化后输出.

```
# apt install jq
```

## example
```bash
$  curl http://localhost:5984/ussfed_tmm -s | jq  # `-s`作用: 忽略curl的process header, 否则jq的输出内容带该信息
```