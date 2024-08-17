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
- --limit-rate: 限速, 比如1024k
- -L : 如果该url是跳转的,比如301,则返回新url的内容
- -o : 将输出写到该文件
- --trace <file>: 将追踪信息写入file,内容比`-v`详细
- --trace-time : 为`--trace`添加时间戳
- -X : 指定请求类型,比如POST
- -A, --user-agent : 指定ua
- -v : 显示一次http通信的整个过程
- -k : 不对服务器的证书进行检查
- --cacert : 检查服务端证书

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

## FAQ
### couchdb更新大视图报`done waiting for 100-continue`
ref:
- [curl Expect:100-continue](https://blog.csdn.net/fdipzone/article/details/42463727)

`curl -v -X PUT --data ...`返回该错并阻塞

原因: 如请求是POST/PUT, 且`--data`超过1024, curl并不会直接就发起POST请求, 而是会分两步:
1.发送一个请求，header中包含一个Expect:100-continue，询问Server是否愿意接受数据
2.接受到Server返回的100-continue回应后，才把数据POST到Server

比如couchdb更新大视图就会卡住, 解决方法: 添加`-H 'Expect:'`禁用该逻辑

补: 实际测试couchdb更新大视图即使分成2步也是没有问题. 部署14套环境, 2套出相同问题, 推测可能与防火墙有关, 待查明

### couchdb更新大视图报`curl -v -X PUT -H 'Expect:' --data ...`报`upload completely sent off: 40 out of 40 bytes`阻塞
