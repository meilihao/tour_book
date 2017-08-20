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

### fetch patch method return 400

chrome : 51.0.2704.84 (64-bit)
nginx : 1.10.1

```js
fetch('/topics/' + id, {
    method: 'patch', // 这里的"patch"应改为"PATCH"
    credentials: 'include',
    body: data
})
```

将method属性的值转化为大写即可.

原因: 此时chrome产生的请求行(`patch /topics/2016060000000007 HTTP/1.1`)的方法是小写,会被nginx拒绝;
而奇怪的是,如果是`method: 'post'`,chrome生成的请求行的方法又会自动转成大写.

### 添加模块
> [Compiling Third-Party Dynamic Modules for NGINX and NGINX Plus](https://www.nginx.com/blog/compiling-dynamic-modules-nginx-plus/)
```sh
$ cd nginx-1.13.4
$ ./configure --with-compat --add-dynamic-module=../echo-nginx-module-master
$ make modules
$ sudo cp objs/ngx_http_echo_module.so /etc/nginx/modules/
$ sudo vim /etc/nginx/nginx.conf # 在文件开始处加入"load_module modules/ngx_http_echo_module.so;"
```