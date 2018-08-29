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
# 获取官方nginx的编译参数
$ sudo nginx -V
....
configure arguments: --prefix=/etc/nginx ...
...
#编译 nginx
$./configure --with-openssl=../openssl ...
```

编译过程可参考[Nginx 编译安装](https://www.dcc.cat/nginx/).

问题:
1. 官方说openssl1.1.1已默认支持tls1.3,但我用OpenSSL 1.1.1-pre9,OpenSSL 1.1.1-pre10,直接编译出来的nginx均不支持TLSv1.3,必须打上[相应的openssl-equal-preX_ciphers.patch](https://github.com/hakasenyang/openssl-patch)才可以.
1. `nginx: [emerg] SSL_CTX_set_cipher_list("TLS13-AES-128-GCM-SHA256 TLS13-CHACHA20-POLY1305-SHA256 TLS13-AES-256-GCM-SHA384") failed (SSL:)`的错误原因: [TLS 1.3 相关的 ciphers 的名字变化](https://www.v2ex.com/t/456622).


### 使用boringssl编译nginx
参考:
- [Nginx 启用 BoringSSL](https://sometimesnaive.org/article/64)

启用boringssl tls 1.3(git commit 2556f8ba60347356f078c753eed2cc65caf5e446,20180829):
```
vim boringssl/ssl/internal.h
# 将两处`tls13_variant_t tls13_variant = xxx;`的值设为`tls13_all`
```

编译boringssl:
```
# 安装编译所需依赖
# BoringSSL 需要 Golang 支持
apt install -y build-essential make cmake golang

# 把 BoringSSL 源码克隆下来
git clone --depth=1 https://boringssl.googlesource.com/boringssl

# 编译开始
cd boringssl
mkdir build && cd build
env CFLAGS=-fPIC CXXFLAGS=-fPIC cmake -DCMAKE_BUILD_TYPE=Release  ..
make
cd ..
mkdir -p .openssl/lib && cd .openssl && ln -s ../include .
cd ..
cp build/crypto/libcrypto.a build/ssl/libssl.a .openssl/lib
cd ..
```
> 其他path: [boringssl支持tls1.3的path, 已验证(git commit 23849f0)](https://github.com/S8Cloud/sslpatch/blob/master/BoringSSL-enable-TLS1.3.patch)

编译nginx 1.15.3:
```sh
./configure  --with-openssl=../boringssl ...
# 在 configure 后，要先 touch 一下，才能继续 make
touch ../boringssl/.openssl/include/openssl/ssl.h
make
sudo make install
```

效果:
```sh
~/tls/nginx-1.15.3/objs$ ./nginx -V
nginx version: nginx/1.15.3
built by gcc 5.4.0 20160609 (Ubuntu 5.4.0-6ubuntu1~16.04.10) 
built with OpenSSL 1.1.0 (compatible; BoringSSL) (running with BoringSSL)
...
```

验证:
```
$ env LD_LIBRARY_PATH=/home/chen/opt/openssl-OpenSSL_1_1_1-pre9 ./openssl s_client -connect openhello.net:443 -tls1_3
...
---
New, TLSv1.3, Cipher is TLS_AES_256_GCM_SHA384 // TLSv1.3 连接成功
Server public key is 256 bit
Secure Renegotiation IS NOT supported
Compression: NONE
Expansion: NONE
No ALPN negotiated
Early data was not sent
Verify return code: 20 (unable to get local issuer certificate)
---
...
```

更好的验证工具:[testssl.sh, 需要相应的openssl版本支持](https://github.com/drwetter/testssl.sh)

其他问题:
1.
```
$ sudo nginx -t
nginx: [warn] "ssl_stapling" ignored, not supported # 可参考[从无法开启 OCSP Stapling 说起](https://toutiao.io/posts/xs1d1d/preview)
...
```
2. 
```
nginx配置:
```
ssl_protocols TLSv1.3;
ssl_ciphers 'TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384';
```
报错:
```
SSL_CTX_set_cipher_list("TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384") failed (SSL: error:100000b1:SSL routines:OPENSSL_internal:NO_CIPHER_MATCH)
```

起先一直以为是boringssl编译出了问题,偶然拷贝了以前的配置发现ok. 经测试,只要添加几个其他TLSv1.x的cipher suite即可解决, 真是奇怪的问题.比如:
```
ssl_ciphers 'TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384 ECDHE-ECDSA-AES128-GCM-SHA256';
```

#### boringssl/.openssl/lib/libssl.a: error adding symbols: Bad value
```fish
$ env CFLAGS=-fPIC CXXFLAGS=-fPIC cmake -DCMAKE_BUILD_TYPE=Release  ..
```

### ssl_ciphers选择
筛选命令(只包含tls1.2):
```sh
$ openssl ciphers -V 'ALL'|grep "1.2"|egrep -v "Kx=DH|Kx=PSK|Kx=ECDHEPSK|RSAPSK|Camellia"|egrep -v "Enc=AESGCM\(256\)|Enc=AESCCM\(256\)|Enc=AESCCM8"|grep -v "Mac=SHA384"|egrep -v "Enc=AES\(256\)"|column  -t 
```
```
TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384 \ # tls1.3
ECDHE-ECDSA-AES128-GCM-SHA256 ECDHE-ECDSA-CHACHA20-POLY1305 \ # ECDHE+ECDSA+AEAD
ECDHE-RSA-AES128-GCM-SHA256 ECDHE-RSA-CHACHA20-POLY1305 \ # ECDHE+RSA+AEAD
ECDHE-ECDSA-AES128-SHA256 ECDHE-RSA-AES128-SHA256 # ECDHE+!AEAD
ECDHE-ECDSA-AES128-SHA ECDHE-RSA-AES128-SHA # TLSv1 for win7,旧Android
```

通过[ssllabs](https://www.ssllabs.com/ssltest/analyze.html)对比发现`ECDHE-ECDSA-*`和`ECDHE-RSA-*`支持的设备跨度是一样的,因此仅保留`ECDSA`即可:
```
TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384 \ # tls1.3
ECDHE-ECDSA-AES128-GCM-SHA256 ECDHE-ECDSA-CHACHA20-POLY1305 \ # ECDHE+ECDSA+AEAD
ECDHE-ECDSA-AES128-SHA256 # ECDHE+!AEAD
ECDHE-ECDSA-AES128-SHA # TLSv1 for win7,旧Android
```

> 在配置 CipherSuite 时，请务必参考权威文档，如：[CloudFlare 使用的配置](https://github.com/cloudflare/sslconfig/blob/master/conf);[Mozilla 的推荐配置](https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations)
> ssl_ecdh_curve选择: `ssl_ecdh_curve   X25519:P-256:P-384:P-224:P-521;`

### 日志优化
```sh
    # $http_host是http(s)的Header Host值,可能为空.
    log_format  main  '[$time_local] $remote_addr - $host - $ssl_protocol/$ssl_cipher $status $body_bytes_sent "$request" '
                      '"$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
```

### Nginx配置网站适配PC和手机
- [Nginx配置网站适配PC和手机](https://blog.csdn.net/xiao__gui/article/details/46680863)
- [detectmobilebrowsers](http://detectmobilebrowsers.com/)