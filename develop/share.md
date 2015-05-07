## 分享

参考:[分享到微信微博QQ空间JS代码大全](http://www.xunhuweb.com/blog/416)

### 微博分享

[微博分享按钮文档](http://open.weibo.com/sharebutton?siteid=372287067)

根据文档,构建网页,再提取相关html代码:

分享所需的url为`http://service.weibo.com/share/share.php?url=http%3A%2F%2F127.0.0.1%3A8080%2F&appkey=&language=zh_cn&title=%E5%94%A4%E9%86%92%E4%B8%96%E7%95%8C%E7%9A%84%E5%A2%99&pic=http://gopher.qiniudn.com/ad/jpush-3.jpg&searchPic=true&ralateUid=XXX`,以`onclick="window.open(url)"`方式打开分享窗口.

参数说明:

- [url] 要分享的url,在分享页面会以`http://t.cn/XXX`短域名的形式显示
- [appkey] [AppKey的作用仅作来源显示用](http://open.weibo.com/wiki/FAQ)
- language 语言设置
- title 要分享的文字
- [pic] 预设要分享的图片
- [searchPic] 抓图服务,true时会抓取参数url网页中的图片追加为要分享的图片
- [ralateUid] 关联账号,在分享文字后追加分享来源(`（分享自 @XXX）`)
- [content] 设置页面编码gb2312|utf-8

### qq空间分享(可支持转播到腾讯微博)

[Qzone分享组件文档](http://connect.qq.com/intro/share/),其会自动在参数url上抓图及提取summary并显示在分享页面上.

分享所需的url为`http://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url=http%3A%2F%2F127.0.0.1%3A8080%2F&showcount=1&desc=1&summary=2&title=3&site=4&pics=http%3A%2F%2Fgopher.qiniudn.com%2Fad%2Fjpush-3.jpg&otype=share`

参数说明:

- [url] 要分享的url
- [showcount] 是否显示分享总数
- desc 要分享的文字
- summary 要分享的摘要
- title 要分享的标题
- site 分享来源
- pics 要分享的图片,多图分享如何分隔未知(已尝试`|`,`||`,`,`均不可行)
- otype 未知

### 微信

分享方式:先把分享url生成二维码,微信扫二维码再分享

实现方式:

- 第三方分享插件
- 自己生成二维码或调用其他接口生成二维码
