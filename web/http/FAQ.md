### url

#### 参数内容被截断

参数内容带**分号**,将其放入url中传值,发现该参数分号后内容被截断(包括分号).解决方法:将其放入Post请求的Form Data中传输.

#### chrome 无法取消301跳转
设置`DevTool - setting - Disable cache`无效
解决: chrome - 设置 - 清除浏览数据, 清除"浏览记录"和"缓存的图片和文件"

### auth
#### Token和cookie区别

token机制是为了防止cookie被清除，另外cookie是会在所有域名请求都携带上，无意中增加了服务端的请求量，token只需要在有必要的时候携带。