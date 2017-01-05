# axel

参考:
- [goxel](github.com/WayneZhouChina/goxel)

## 原理

服务端支持Range请求

**下载时需注意文件是否变化,即检查`ETag`或`Last-Modified`;或使用`If-Range`**

## 具体

1. 获取资源信息 : 发送HEAD请求并读取resp.Header中的"Accept-Ranges"(不为`none`表示支持Range请求)和"Content-Length"(大小)
1. 发送带"Range"请求头的request, 成功时返回的resp的status是206.
1. 使用文件操作的Seek()来保存数据.
