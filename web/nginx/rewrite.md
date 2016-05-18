## Rewrite

使用nginx提供的全局变量或自己设置的变量，结合正则表达式和标志位实现url重写以及重定向.
rewrite只能放在`server{},location{},if{}`中，并且只能对域名后边的**除去传递的参数外的字符串起作用**.

### 执行顺序:

1. 执行server块的rewrite指令
2. 执行location匹配
3. 执行选定的location中的rewrite指令

如果其中某步URI被重写，则重新循环执行1-3，直到找到真实存在的文件.
循环超过10次，则返回500 Internal Server Error错误

### flag标志位

- last : 停止执行当前这一轮的ngx_http_rewrite_module指令集,新的uri会被nginx的location处理,但后续的rewrite指令将会被忽略
- break : 停止执行当前这一轮的ngx_http_rewrite_module指令集
- redirect : 返回302临时重定向，地址栏会显示跳转后的地址
- permanent : 返回301永久重定向，地址栏会显示跳转后的地址

> 因为301和302不能简单的只返回状态码，还必须有重定向的URL，这就是return指令无法返回301,302的原因了.

这里 last 和 break 区别有点难以理解：

- last一般写在server和if中，而break一般使用在location中
- last不终止重写后的url匹配，即新的url会再从server走一遍匹配流程，而break终止重写后的匹配
