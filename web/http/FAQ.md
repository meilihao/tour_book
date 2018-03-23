### url

#### 参数内容被截断

参数内容带**分号**,将其放入url中传值,发现该参数分号后内容被截断(包括分号).解决方法:将其放入Post请求的Form Data中传输.

#### chrome 无法取消301跳转
设置`DevTool - setting - Disable cache`无效
解决: chrome - 设置 - 清除浏览数据, 清除"浏览记录"和"缓存的图片和文件"