# FAQ
### 无法下载/中文乱码
ref:
- [URL中的空格、加号究竟应该使用何种方式编码](https://segmentfault.com/a/1190000040919065)

response:
```
Content-Type:application/octet-stream
Content-Disposition:attachment; filename="a.txt"; filename*=utf-8''a.txt
```

> `filename*`是新标准. filename的文件名可用encodeURIComponent()编码, `filename*`使用指定编码或采用与filename相同的编码.

```go
// go的url.QueryEscape与js的encodeURIComponent是不一样的，主要体现在对空格的处理，此函数编码后的字符串可以被js的decodeURIComponent正确还原
func encodeURIComponent(str string) string {
    r := url.QueryEscape(str)
    r = strings.Replace(r, "+", "%20", -1)
    return r
}
```

resp返回的header name变成小写导致.

2010年 RFC 5987 发布，正式规定了 HTTP Header 中多语言编码的处理方式采用`parameter*=charset'lang'value`的格式，其中：
- charset 和 lang 不区分大小写
- lang 是用来标注字段的语言，以供读屏软件朗诵或根据语言特性进行特殊渲染，可以留空
- value 根据 RFC 3986 Section 2.1 使用**百分号编码**，并且规定浏览器至少应该支持 ASCII 和 UTF-8
- 当 parameter 和 parameter* 同时出现在 HTTP 头中时，浏览器应当使用后者

### 格式化Curl返回的Json字符
- `curl https://news-at.zhihu.com/api/4/news/latest | jq .` # 基于[jq](https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64), **推荐**
- `curl https://news-at.zhihu.com/api/4/news/latest | python -m json.tool`
- `curl https://news-at.zhihu.com/api/4/news/latest -s | json` # 基于`npm install -g json`