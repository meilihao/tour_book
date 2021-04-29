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
1. 官方说openssl1.1.1已默认支持tls1.3,但我用OpenSSL 1.1.1-pre9,OpenSSL 1.1.1-pre10,直接编译出来的nginx均不支持TLSv1.3,必须打上[相应的openssl-equal-preX_ciphers.patch](https://github.com/hakasenyang/openssl-patch)才可以. 这是[网上找到的解释: 在OpenSSL 1.1.1-preX中，除TLS1.3最终版本之外的所有Draft版本都已删除, 但是浏览器仅支持Draft版本](https://serverfault.com/questions/927601/tls1-3-not-working-on-nginx-1-15-2-with-openssl-1-1-1-pre9).
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

更好的验证工具:[testssl.sh, 需要相应的openssl版本支持](https://github.com/drwetter/testssl.sh)或MySSL.com,https://www.ssllabs.com/ssltest/.

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
ps: `2018.8.29`收到boringssl team的反馈邮件`TLS 1.3 ciphers in BoringSSL aren't configurable. The "ssl_ciphers" list only configures TLS 1.2.`.

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

### auth_basic_user_file
auth_basic_user_file格式:
```
username:crypted_password
```

密码生成:
```
openssl passwd [-crypt] xxx
```

### `expires`/缓存不起作用
```
location ^~ /static {
    root /var/www/files; # 静态文件目录
    expires 30d;
}
```

明明nginx使用了`expires`, 且返回的respone header里也有`Cache-Control: max-age=2592000`和`Expires: xxx`,但浏览器还是重新获取资源而不是使用缓存.

原因:
chrome开发者工具`Network-Disable cache(while DevTools is open)`选项被启用了, 而当时该工具又恰好开着, 结果就是该请求被缓存了.

> firefox的类似选项是`高级设置-禁用HTTP缓存(工具箱打开时)`

### 对后端某个节点的优雅下线
Nginx将请求代理给一个后端节点，这个请求耗时较长，在请求未处理完时后端恰好要做发布.这时在Nginx中先将此节点标记为不可用（在upstream中设置server的down属性）.
此时，只要请求连接还保持，Nginx并不会中断当前连接，但之后新的连接将不再使用这个节点.

这样在用Nginx的负载时，后端若需要做发布. 只需要将对就节点标记为不可用并留出一定的时间让忆有请求都响应完毕即可.

更严格一些，还应检测到后端节点的网络连接都已释放（那些EST、TIME_WAIT等连接都结束后）, 或使用signal获取应用还在服务的req count.

### tengine jemalloc gcc-9.1
jemalloc-5.2.0:
```sh
$ cd jemalloc-5.2.0
$ env CC=/usr/local/gcc-9.1/bin/gcc-9.1 ./configure
$ make
$ sudo make install
```

ngx_brotli:
```sh
$ git clone --depth=1 git@github.com:eustas/ngx_brotli.git
$ git submodule update --init
$ cd deps/brotli/
$ git pull origin master
```

> ngx_brotli使用git@github.com:google/ngx_brotli.git并更新依赖`deps/brotli`(`git pull origin master`)后, nginx报错: `Brotli library is missing from /xxx/ngx_brotli/deps/brotli directory`; 不更新`deps/brotli`时不报错. 怀疑是ngx_brotli/src与`deps/brotli`版本没对应的原因.
> 顺便再说一句: [google/ngx_brotli is dead](https://github.com/google/ngx_brotli/issues/62) and [eustas is primary contributor to google/brotli](https://github.com/NixOS/nixpkgs/pull/34943)

tengin-2.3.1:
```sh
$ cd tengine-2.3.1
$ ./configure --with-cc=/usr/local/gcc-9.1/bin/gcc-9.1 --with-jemalloc=/home/chen/jemalloc-5.2.0 --add-module=/home/chen/ngx_brotli
$ make
$ sudo make install
```

nginx.conf
```
# 在http块下增加以下指令, 重启nginx后,用浏览器或抓包查看css和图片里的请求头是否包含`Content-Encoding: br`.
brotli on; # 是否启用在on-the-fly方式压缩文件，启用后，将会在响应时对文件进行压缩并返回
brotli_types text/xml text/plain application/json text/css image/svg application/font-woff application/vnd.ms-fontobject application/vnd.apple.mpegurl application/javascript image/x-icon image/jpeg image/gif image/png; # 指定对哪些内容编码类型进行压缩. text/html内容是默认会被进行压缩,无需添加.
brotli_static on;     # 启用后将会检查是否存在带有br扩展的预先压缩过的文件。如果值为always，则总是使用压缩过的文件，而不判断浏览器是否支持
brotli_comp_level 6;  # 设置压缩质量等级(0~11)
brotli_buffers 16 10k;# 设置缓冲的数量和大小, 大小默认为一个内存页的大小，也就是4k或者8k
brotli_window 512k;   # 设置窗口大小
brotli_min_length 20; # 设置需要进行压缩的最小响应大小
```

### [nginx 上传大文件超时](https://my.oschina.net/ericquan8/blog/379265)
```conf
 location / {  
        proxy_pass     http://xxx;
        // for big upload
        client_max_body_size    1000m;
        client_body_timeout 15m;
        proxy_connect_timeout         15m;
	    proxy_read_timeout            15m;
    } 
```

### conf.d/sites-enabled/sites-available
> 新版nginx的目录设置与httpd类似.

conf.d
这是一个目录, 用于全局服务器配置,文件结尾一
定是.conf才可以生效(当然也可以通过修改nginx.conf来取消这个限制)

sites-enabled
这里面的配置文件其实就是sites-available里面的配置文件的软
连接,但是由于nginx.conf默认包含的是这个文件夹,所以在
sites-available里面建立了新的站点之后,还要建立个软连接到sites-enabled里面才行

sites-available
虚拟主机的目录，可在这里面可以创建多个虚拟主机

### nginx server_name directive is not allowed here
nginx.conf配置错误, 当前加载文件中的server配置在nginx的配置嵌套层次不正确, 它因在`http`配置项下.

### php-fpm
env: ubuntu 20.4, nginx 1.18

```bash
apt install php php-fpm
service php7.4-fpm status
```

```conf
server {

        listen       9100;
        server_name  bareos;
        root         /var/www/bareos-webui/public;

        location / {
                index index.php;
                try_files $uri $uri/ /index.php?$query_string;
        }

        location ~ .php$ {
                include snippets/fastcgi-php.conf;

                # php7-cgi alone:
                # pass the PHP
                # scripts to FastCGI server
                # listening on 127.0.0.1:9000
                #fastcgi_pass 127.0.0.1:9000;

                # php7-fpm:
                fastcgi_pass unix:/var/run/php/php-fpm.sock;

                # APPLICATION_ENV:  set to 'development' or 'production'
                #fastcgi_param APPLICATION_ENV development;
                fastcgi_param APPLICATION_ENV production;
        }

}
```

> php-fpm access log: `/etc/php/7.4/fpm/pool.d/xxx.conf`的`[xxx]`项的`access_log`

nginx+php原理:
1. nginx的worker进程直接管理每一个到nginx的网络请求

1. 对于php而言，由于在整个网络请求的过程中php是一个cgi程序的角色，所以采用名为php-fpm的进程管理程序来对这些被请求的php程序进行管理. php-fpm程序也如同nginx一样，需要监听端口，并且有master和worker进程. worker进程直接管理每一个php进程.

    1. 关于fastcgi：fastcgi是一种进程管理器，管理cgi进程. 当前有多种实现了fastcgi功能的进程管理器，php-fpm就是其中的一种. 再提一点，php-fpm作为一种fast-cgi进程管理服务，会监听端口(也可使用unix socks)，一般默认监听9000端口，并且是监听本机，也就是只接收来自本机的端口请求，可通过`ss -nlpt|grep php-fpm`查看

    1. 关于fastcgi的配置文件，目前fastcgi的配置文件一般放在nginx.conf同级目录下，配置文件形式，一般有两种：fastcgi.conf 和 fastcgi_params.

1. 当需要处理php请求时，nginx的worker进程会将请求移交给php-fpm的worker进程进行处理，也就是最开头所说的nginx调用了php，其实严格得讲是nginx间接调用php.

> [PHP Notice: compact(): Undefined variable: extras in src\Helper\HeadLink.php](https://github.com/zendframework/zend-view/pull/170/files)