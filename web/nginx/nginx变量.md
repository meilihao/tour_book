# nginx config
参考:
- [nginx server 如何写，如何写 nginx 配置](https://deepzz.com/post/how-to-write-nginx-server.html)

## 请求头

- $host	请求Header中的Host
- $http_user_agent 请求Header中的User-Agent
- $http_referer	请求Header中的Referer,告诉server从哪个页面链接而来
- $http_cookie 请求Header中的cookie
- $http_via 请求Header中的Via,客户端可能使用的代理
- $http_x_forwarded_for	client访问server的路径
- $http_xxx 获取header中的其他信息,名称要小写,名称中的`-`用`_`代替

## 响应头

- $sent_http_content_type HTTP header中的Content-Type
- $sent_http_content_length HTTP header中的Content-Length,表明http body的长度
- $sent_http_location HTTP header中的Location,重定向client到一个新URI地址
- $sent_http_last_modified HTTP header中的Last-Modified
- $sent_http_keep_alive HTTP header中的Keep-Alive
- $sent_http_transfer_encoding HTTP header中的Transfer-Encoding
- $sent_http_xxx 获取header中的其他信息,名称要小写,名称中的`-`用`_`代替

## nginx产生的变量

- $arg_xxx 请求行中的指定参数,xxx替换为具体的参数名
- $args 请求行中的所有参数，同$query_string
- $binary_remote_addr 二进制的ip地址(ipv4,四字节;ipv6,16字节)
- $body_bytes_sent 响应body中发送的字节数
- $bytes_sent 发送给客户端的总字节数
- $connection 连接的序列号
- $connection_requests 当前通过一个连接获得的请求数量
- $content_length 相当于HTTP Header中的Content-Length
- $content_type 相当于HTTP Header中的Content-Type
- $cookie_xxx 指定cookie的数据,xxx替换为具体的cookie name
- $document_root 对当前请求返回root或alias指令的指定值
- $document_uri 同$uri
- $host 请求行中的host name,或请求头中的"Host"字段,或nginx匹配到的server_name的值(原始请求没有提供host信息)
- $hostname host name
- $http_xxx 任意请求头的属性值；属性名使用小写并且用`_`替代`-`
- $https 如果连接使用SSL，返回“on”，否则返回空字符串
- $is_args 如果请求行带有参数，返回“?”，否则返回空字符串
- $limit_rate 允许设置此值来限制连接的传输速率
- $msec 当前时间，单位是秒，精度是毫秒
- $nginx_version nginx版本
- $pid worker进程的PID
- $pipe 如果请求是通过HTTP流水线(pipelined)发送，pipe值为`p`，否则为`.`
- $proxy_protocol_addr PROXY协议头中的客户端地址或"".前提条件,PROXY协议必须在使用前通过proxy_protocol参数启用.
- $query_string 同$args
- $realpath_root 对当前请求返回root或alias指令的指定值的绝对路径,其中的符号链接都会解析成真实路径.
- $remote_addr client地址
- $remote_port client端口
- $remote_user 为基本用户认证提供的用户名
- $request 完整的原始请求行
- $request_body 请求正文.在proxy_pass指令,fastcgi_pass指令,uwsgi_pass指令和scgi_pass指令的处理过程中， 这个变量可用.
- $request_body_file 请求body的临时文件名.处理完成时，该临时文件将被删除.如果希望总是将请求正文写入文件，需要开启     client_body_in_file_only.如果在被代理的请求或FastCGI请求中传递临时文件名，就应该禁止传递请求正文本身.
可使用proxy_pass_request_body off指令和fastcgi_pass_request_body off指令分别禁止在代理和FastCGI中传递请求正文.
- $request_completion 请求完成时返回“OK”，否则返回空字符串
- $request_filename
file path for the current request, based on the root or alias directives, and the request URI
- $request_id 唯一的请求标识符
- $request_length 请求的总长度(包括request line, header, 和 request body)
- $request_method HTTP方法，通常为"GET"或者"POST"
- $request_time 请求处理时间,单位是秒,精度是毫秒.请求处理时间由从客户端接收到第一个字节开始计时, 用于log_format, **推荐使用**
- $upstream_response_time Nginx 接受完 client 的请求后，再和 upstream server 请求的过程，这个指标才能真正反应 upstream server 的响应情况，虽然通常`upstream_response_time<=request_time`, 用于log_format, **推荐使用**
- $request_uri 完整的原始请求行（带参数）
- $scheme 请求协议类型，为"http"或"https"
- $sent_http_xxx 任意响应头的属性值；属性名使用小写并且用`_`替代`-`
- $server_addr 接受请求的服务器地址.为计算这个值，通常需要进行一次系统调用.为了避免系统调用，必须指定listen指令 的地址，并且使用bind参数
- $server_name 接受请求的虚拟主机的主机名
- $server_port 接受请求的虚拟主机的端口
- $server_protocol 请求协议，通常为"HTTP/1.0","HTTP/1.1"或"HTTP/2.0"
- $status 响应的status
- $tcpinfo_rtt, $tcpinfo_rttvar, $tcpinfo_snd_cwnd, $tcpinfo_rcv_space 客户端TCP连接的信息，在支持套接字选项TCP_INFO的系统中可用
- $time_iso8601 ISO8601标准格式下的本地时间
- $time_local 通用日志格式中使用的本地时间
- $uri 当前请求规范化以后的URI.变量$uri的值可能随请求的处理过程而改变.比如，当进行内部跳转时，或者使用默认页文件.

> 规范化，见[location.md](location.md)
