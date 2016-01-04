### `error: the HTTP rewrite module requires the PCRE library`

```
sudo apt-get install libpcre3 libpcre3-dev
sudo dnf -y install pcre-devel
```

### `error: SSL modules require the OpenSSL library`

```
sudo apt-get install libssl-dev
sudo dnf -y install  openssl-devel
```

### `nginx: [emerg] unknown log format "main" in /usr/local/nginx/conf/nginx.conf`

打开access_log选项后未启用日志格式main，启用log_format	main(nginx默认日志格式)即可.

### `An error occurred`

可能情况:

`proxy_pass`指定的地址不可用.

### error.log出现`failed (13: Permission denied)`

查看`ps aux | grep "nginx: worker process" | awk '{print $1}'`(即nginx.conf的user指令)显示的用户是否对该请求路径(绝对路径)有无访问权限.
我这里是因将WebRoot放在了主目录下导致用户nginx无权限访问的原因.
