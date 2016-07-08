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
- [pic] 预设要分享的图片(必须带扩展名)
- [searchPic] 抓图服务,true时会抓取参数url网页中的图片追加为要分享的图片
- [ralateUid] 关联账号,在分享文字后追加分享来源(`（分享自 @XXX）`)
- [content] 设置页面编码[gb2312|utf-8,推荐]

在title中分享url时,如果出现中文参数,分享后该参数会被识别成文本,除非该参数被重复编码3次及以上,不推荐.
因此建议(**推荐**):不论title中的url是否含非ASCII字符都先将title中的url转成[短链](http://open.weibo.com/wiki/Short_url/shorten)再分享,如果url含非ASCII字符时,先urlencode该参数,再生成短链.

[微博分享的卡片效果](http://travel.sohu.com/20160503/n447372910.shtml),该功能需要单独申请linkcard解析,申请该功能需要应用授权人数达到1000人以上.

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

### 豆瓣

[文档](http://www.douban.com/service/bookmarklet),在"推荐到豆瓣"按钮的html代码中.

分享所需的url为`http://www.douban.com/share/service?image=http://d.hiphotos.baidu.com/image/pic/item/b03533fa828ba61e43973e704234970a304e5970.jpg&href=image.baidu.com&name=title&text=content`

参数说明:

- [image] 要分享的图片
- href 分享的url
- title 要分享的标题
- content 分享内容
