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

### threads_pthread.c:(.text+0x16)：对‘pthread_atfork’未定义的引用
nginx: 1.13.5
在自编译指定openssl时碰到:
```sh
../openssl-master/.openssl/lib/libcrypto.a(threads_pthread.o)：In function ‘fork_once_func’：
threads_pthread.c:(.text+0x16)：undefined reference to 'pthread_atfork'
```

对该issue openssl提供了[解决方法](https://github.com/openssl/openssl/issues/3884#issuecomment-313857555).
出现原因是nginx的`./configure`生成的`obj/Makefile`有误.
修正方法:
- 找到`libcrypto.a`所在的行
- 仅保留一个`-lpthread`,删除其他多余的
- 将`-lpthread`移动到行尾的`\`前,注意保留其两边的空格.
- 保存文件,返回上级目录运行`make`即可.

## FAQ
### 添加模块
> [Compiling Third-Party Dynamic Modules for NGINX and NGINX Plus](https://www.nginx.com/blog/compiling-dynamic-modules-nginx-plus/)
```sh
$ cd nginx-1.13.4
$ ./configure --with-compat --add-dynamic-module=../echo-nginx-module-master
$ make modules
$ sudo cp objs/ngx_http_echo_module.so /etc/nginx/modules/
$ sudo vim /etc/nginx/nginx.conf # 在文件开始处加入"load_module modules/ngx_http_echo_module.so;"
```

### 查看nginx的版本及编译参数
```sh
# 官方mainline的编译参数
$ sudo nginx -V
nginx version: nginx/1.13.5
built by gcc 4.8.5 20150623 (Red Hat 4.8.5-11) (GCC) 
built with OpenSSL 1.0.1e-fips 11 Feb 2013
TLS SNI support enabled
configure arguments: --prefix=/etc/nginx --sbin-path=/usr/sbin/nginx --modules-path=/usr/lib64/nginx/modules --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --pid-path=/var/run/nginx.pid --lock-path=/var/run/nginx.lock --http-client-body-temp-path=/var/cache/nginx/client_temp --http-proxy-temp-path=/var/cache/nginx/proxy_temp --http-fastcgi-temp-path=/var/cache/nginx/fastcgi_temp --http-uwsgi-temp-path=/var/cache/nginx/uwsgi_temp --http-scgi-temp-path=/var/cache/nginx/scgi_temp --user=nginx --group=nginx --with-compat --with-file-aio --with-threads --with-http_addition_module --with-http_auth_request_module --with-http_dav_module --with-http_flv_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_mp4_module --with-http_random_index_module --with-http_realip_module --with-http_secure_link_module --with-http_slice_module --with-http_ssl_module --with-http_stub_status_module --with-http_sub_module --with-http_v2_module --with-mail --with-mail_ssl_module --with-stream --with-stream_realip_module --with-stream_ssl_module --with-stream_ssl_preread_module --with-cc-opt='-O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -fexceptions -fstack-protector-strong --param=ssp-buffer-size=4 -grecord-gcc-switches -m64 -mtune=generic -fPIC' --with-ld-opt='-Wl,-z,relro -Wl,-z,now -pie'
```

自编译追加参数:
- 指定openssl版本,启用tls1.3
```
# 等nginx编译完成后,到openssl里运行`make test`,检查一下tls1.3是否确实启用了.
--with-openssl=../openssl --with-openssl-opt='enable-tls1_3'
```
> 注意请使用**tls1.3-draft-18分支**而不是master分支或tls1.3-draft-19分支,因为到**2017.9.19**为止浏览器还不支持tls1.3的最新草案,
> 所以只启用tls1.3时用openssl s_client测试会成功,但真实浏览器请求却会失败报`ERR_SSL_VERSION_OR_CIPHER_MISMATCH`.

### ssl_ciphers选择
筛选命令(不包含tls1.3):
```sh
$ openssl ciphers -V 'ALL'|grep "1.2"|egrep -v "Kx=DH|Kx=PSK|Kx=ECDHEPSK|RSAPSK|Camellia"|egrep -v "Enc=AESGCM\(256\)|Enc=AESCCM\(256\)|Enc=AESCCM8"|grep -v "Mac=SHA384"|egrep -v "Enc=AES\(256\)"|column  -t 
```
```
TLS13-AES-128-GCM-SHA256 TLS13-CHACHA20-POLY1305-SHA256 TLS13-AES-128-CCM-SHA256 TLS13-AES-128-CCM-8-SHA256 TLS13-AES-256-GCM-SHA384 \ # tls1.3
ECDHE-ECDSA-AES128-GCM-SHA256 ECDHE-ECDSA-CHACHA20-POLY1305 ECDHE-ECDSA-AES128-CCM \ # ECDHE+ECDSA
ECDHE-RSA-AES128-GCM-SHA256 ECDHE-RSA-CHACHA20-POLY1305 \ # ECDHE+RSA+AEAD
ECDHE-ECDSA-AES128-SHA256 ECDHE-RSA-AES128-SHA256 # ECDHE+!AEAD
```