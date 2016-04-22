# curl

Curl is a **command line** tool for transferring data specified with URL syntax.

## 描述

转换url

## 选项

- --cookie : 指定cookie
- -b STRING/FILE : 使用STRING/FILE中的cookie
- -c, --cookie-jar FILE : 保存服务器返回的cookie到文件
- -d, --data : 指定post的data
- --data-urlencode <data> : 对data使用urlencode
- -e, --referer : 指定header的referer
- -f : 模拟表单提交
- -i : 返回内容带header
- -I : 仅返回header
- --header : 添加header
- -L : 如果该url是跳转的,比如301,则返回新url的内容
- -o : 将输出写到该文件
- --trace <file>: 将追踪信息写入file,内容比`-v`详细
- --trace-time : 为`--trace`添加时间戳
- -X : 指定请求类型,比如POST
- -A, --user-agent : 指定ua
- -v : 显示一次http通信的整个过程

## 例
```sh
$ curl www.sina.com # 仅html文档
$ curl -o [文件名] www.sina.com
$ curl example.com/form.cgi?data=xxx # get请求
$ curl -X POST --data-urlencode "date=April 1" example.com/form.cgi
$ curl --user-agent "[User Agent]" [URL]
$ curl --form upload=@localfilename --form username=chen [URL]
$ curl --cookie "user=root;pass=123456" example.com
$ curl -c cookies http://example.com
$ curl -b cookies http://example.com
$ curl --header "Content-Type:application/json" http://example.com
```
