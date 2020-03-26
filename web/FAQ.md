# FAQ
### 无法下载
response:
```
content-type:application/octet-stream
content-disposition:attachment;filename="a.txt"
```

resp返回的header name变成小写导致.

### 格式化Curl返回的Json字符
- `curl https://news-at.zhihu.com/api/4/news/latest | jq .` # 基于[jq](https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64), **推荐**
- `curl https://news-at.zhihu.com/api/4/news/latest | python -m json.tool`
- `curl https://news-at.zhihu.com/api/4/news/latest -s | json` # 基于`npm install -g json`