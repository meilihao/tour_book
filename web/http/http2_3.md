# http2
参考:
- [HTTP/2 in GO(五)--大结局](https://www.tuicool.com/articles/jMreQbN)
- [HTTP/2 资料汇总](https://imququ.com/post/http2-resource.html)

> HTTP/2可以不配合HTTPS来实现和使用, 但当前所有其实现都是基于tls的.

HTTP 2.0 通过头压缩、分帧、二进制编码、多路复用等技术提升性能.

HTTP/2解决了HTTP的队头拥塞（head of line blocking）问题, 客户端必须等待一个请求完成才能发送下一个请求的日子过去了.

随着丢包率的增加，HTTP/2的表现越来越差. 在2%的丢包率（一个很差的网络质量）中，测试结果表明HTTP/1用户的性能更好，因为HTTP/1一般有六个TCP连接，哪怕其中一个连接阻塞了，其他没有丢包的连接仍然可以继续传输. 
在限定的条件下，在TCP下解决这个问题相当困难, 即http2未解决TCP队头阻塞, 但http3解决了.

## http3
ref:
- [检查是否已启用h3](https://http3check.net)
- [HTTP/3 来了：QUIC 协议在 OPPO 的应用](https://my.oschina.net/u/4273516/blog/8597013)
- [nginx http3 demo](https://quic.nginx.org/)
- [Current results of QUIC interop test runner](https://interop.seemann.io/)
- [**闷棍暴打面试官HTTP3-上篇**](https://zhuanlan.zhihu.com/p/580021385)
- [5.5.2 面向HTTP 3.0时代的高性能网络协议栈](https://developer.aliyun.com/article/1229379)
- [甩掉TCP协议的HTTP/3，真的很牛吗？](https://dbaplus.cn/news-160-5556-1.html)

    0-RTT存在前向安全问题

    QUIC更耗CPU和内存, 至少需要增加20%的服务器成本

    开源协议栈的选择有很多，比较知名的有：Google的quiche，cloudflare的quiche，lsquic，亚马逊的s2n-quic

    丢包就需要重传，传统的TCP对于重传包有“二义性”，无法精确地计算RTT。而QUIC采用单调递增的Packet Number方式，解决了重传包的“二义性”问题，RTT计算准确.

HTTP/3不存在明文的不安全版本.

### test status
date: 2023/4/27

server: salvo 0.38, 使用http3时会同时开启http2支持
client: curl(自编译支持http3), [quiche](https://github.com/cloudflare/quiche)

测试情况:
- curl/quiche 访问http3报错
- curl(os自带) 访问http2响应正常当, salvo有出现错误日志

### [构建curl with http3](https://curl.se/docs/http3.html)
**注意**: 对于在 Linux 上为 x86_64 架构构建的 OpenSSL 3.0.0 或更高版本, 需要将所有出现的`/lib`替换为`/lib64`, 下文已修改.

```bash
--- args:
-- somewhere1 : /usr/local
-- somewhere2 : /usr/local
-- somewhere3 : /usr/local
--- Build (patched) OpenSSL
$ git clone --depth 1 -b openssl-3.0.8+quic https://github.com/quictls/openssl
$ cd openssl
$ ./config enable-tls1_3 --prefix=<somewhere1>
$ make
$ sudo make install
--- Build nghttp3
$ cd ..
$ git clone -b v0.10.0 https://github.com/ngtcp2/nghttp3
$ cd nghttp3
$ autoreconf -fi
$ ./configure --prefix=<somewhere2> --enable-lib-only
$ make
$ sudo make install
--- Build ngtcp2
$ cd ..
$ git clone -b v0.13.1 https://github.com/ngtcp2/ngtcp2
$ cd ngtcp2
$ autoreconf -fi
$ ./configure PKG_CONFIG_PATH=<somewhere1>/lib64/pkgconfig:<somewhere2>/lib64/pkgconfig LDFLAGS="-Wl,-rpath,<somewhere1>/lib64" --prefix=<somewhere3> --enable-lib-only
$ make
$ sudo make install
--- Build curl
$ cd ..
$ git clone https://github.com/curl/curl
$ cd curl
$ autoreconf -fi
$ PKG_CONFIG_PATH=<somewhere1>/lib64/pkgconfig LDFLAGS="-Wl,-rpath,<somewhere1>/lib64" ./configure --with-openssl=<somewhere1> --with-nghttp3=<somewhere2> --with-ngtcp2=<somewhere3>
$ make
$ sudo make install
$ /usr/local/bin/curl --http3-only https://nghttp2.org:4433
```

可能的问题:
- nghttp3#`Libtool library used but 'LIBTOOL' is undefined` : 安装libtool
- curl#`--with-ngtcp2 was specified but could not find ngtcp2 pkg-config file`: 指定PKG_CONFIG_PATH

## quic
ref:
- [quic-go now supports QUIC version 1 (RFC 9000) and version 2 (RFC 9369)](https://github.com/quic-go/quic-go/releases/tag/v0.37.0)

QUIC没有非加密的版本.

跟 TLS 类似，QUIC 的 0RTT 握手，是建立在已经同一个服务器建立过连接的基础上，所以如果是纯的第一次连接，仍然需要一个 RTT.

## tengine http3
ref:
- [tengine-ingress/images/tengine/rootfs/build.sh](https://github.com/alibaba/tengine-ingress/blob/2a8038560c0b26f948fb40d592de427b00100e8a/images/tengine/rootfs/build.sh)
- [tengine/modules/ngx_http_xquic_module](https://github.com/alibaba/tengine/tree/master/modules/ngx_http_xquic_module)
- [Module ngx_http_v3_module](https://nginx.org/en/docs/http/ngx_http_v3_module.html)
- [使用 nginx-http3 部署 HTTP/3 服务](https://www.thebyte.com.cn/http/nginx-quic.html)
- [xquic:a client and server implementation of QUIC and HTTP/3 as specified by the IETF. Currently supported QUIC versions are v1 and draft-29](https://github.com/alibaba/xquic)

env:
- chrome 114
- firefox 115

这里按照`https://github.com/alibaba/tengine/tree/master/modules/ngx_http_xquic_module#readme`来编译.

能编译成功, 启动tengine有监听2443端口(`ss -anlp |grep 2443`).

http3.conf:
```conf
worker_processes  1;
user root;
events {
    worker_connections  1024;
}

error_log  logs/error.log;
error_log  logs/error.log  notice;
error_log  logs/error.log  info;
xquic_log   "pipe:rollback /usr/local/tengine/logs/tengine-xquic.log baknum=10 maxsize=1G interval=1d adjust=600" debug;

http {
    xquic_ssl_certificate        /usr/local/tengine/ssl/www.pem;
    xquic_ssl_certificate_key    /usr/local/tengine/ssl/www-key.pem;

    ssl_certificate        /usr/local/tengine/ssl/www.pem;
    ssl_certificate_key    /usr/local/tengine/ssl/www-key.pem;

    server {
        listen 80 default_server reuseport backlog=4096;
        listen 443 default_server reuseport backlog=4096 ssl http2;
        listen 2443 default_server reuseport backlog=4096 xquic;

        ssl_protocols TLSv1.3;
        ssl_early_data on;

        add_header Alt-Svc 'h3=":2443"; ma=86400, h3-29=":2443"; ma=86400';

        location / {
            root   html;
            index  index.html index.htm;
        }
    }
}
```

> firefox/chrome测试发现: 访问`http://127.0.0.1`时Alt-Svc没生效, 一直使用h1.1, 推测需要先重定向到h2.

改为:
```
worker_processes  1;
user root;
events {
    worker_connections  1024;
}

error_log  logs/error.log;
error_log  logs/error.log  notice;
error_log  logs/error.log  info;
xquic_log   "pipe:rollback /usr/local/tengine/logs/tengine-xquic.log baknum=10 maxsize=1G interval=1d adjust=600" debug;

http {
    xquic_ssl_certificate        /usr/local/tengine/ssl/www.pem;
    xquic_ssl_certificate_key    /usr/local/tengine/ssl/www-key.pem;

    ssl_certificate        /usr/local/tengine/ssl/www.pem;
    ssl_certificate_key    /usr/local/tengine/ssl/www-key.pem;

    server {
        listen 80 default_server reuseport backlog=4096;
        server_name _;
        return 301 https://$host$request_uri;
    }

    server {
        listen 443 default_server reuseport backlog=4096 ssl http2;
        listen 2443 default_server reuseport backlog=4096 xquic;

        ssl_protocols TLSv1.3;
        ssl_early_data on;

        add_header Alt-Svc 'h3-29=":2443"; ma=2592000, h3=":2443"; ma=2592000' always;

        location / {
            root   html;
            index  index.html index.htm;
        }
    }
}
```

可能的问题:
- xquic-1.6.0:

    - `error: dangling pointer to ‘peer_tp’ may be used [-Werror=dangling-pointer=]`

        参考`tengine-ingress/images/tengine/rootfs/build.sh`, 在`xquic-1.6.0/CMakeLists.txt`的`MESSAGE("SSL_INC_PATH= ${SSL_INC_PATH}")`下追加`set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-dangling-pointer")`
    - cmake启用`DXQC_ENABLE_TESTING`报`fatal error: event2/event.h: No such file or directory`

        `apt install libevent-dev`
    - cmake启用`DXQC_ENABLE_TESTING`报`error: ignoring return value of ‘fscanf’ declared with attribute ‘warn_unused_result’ [-Werror=unused-result]`

        参考`Werror=dangling-pointer=`追加`set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-unused-result")`
    - cmake启用`DXQC_ENABLE_TESTING`报`error: ‘__builtin_strncpy’ output may be truncated copying 255 bytes from a string of length 255 [-Werror=stringop-truncation]`

        参考`Werror=dangling-pointer=`追加`set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wno-stringop-truncation")`
- tengine

    权限问题: xquic_ssl_certificate/xquic_ssl_certificate_key在root下, tengine读取会报`|xquic|lib[error] |xqc_create_server_ssl_ctx|SSL_CTX_use_PrivateKey_file| error info:error:0200100D:system library:fopen:Permission denied||`

    tengine 配置启用`user root;`
- chrome访问`https://localhost`报`ERR_SSL_KEY_USAGE_INCOMPATIBLE`

    [直接使用了ca.cert作为ssl_certificate, 它的用途不包括digitalSignature, keyEncipherment](https://www.nextdoorwith.info/wp/se/x509-certificate-error-reason-and-how-to-resolve/), chrome检查了证书用途导致报错.

    解决方法: 使用ca证书创建server证书
- 访问`https://localhost`
    
    - chrome访问一直是h2(已开启`Experimental QUIC protocol`, 即`--enable-quic`), 可能还需要手动设置[`--quic-version=h3-29`](https://www.taurusxin.com/http3-enabled/), 修改google-chrome.desktop, 追加`--quic-version=h3-29`/`--quic-version=h3`无效果.
    - firefox需要`移除例外`再访问才使用h3, 否则使用h2. 协议变化还可能与缓存有关, 因为首次请求是h2的. 取消`禁用缓存`, 此时缓存的是h2响应, 后续请求都变成h2. 同时确认当前使用h3, 过几分钟后刷新页面, 请求变h2, 也可能不变, 当前发现变h2概率更大些.


其他:
- Alt-Svc(Alternative-Service, 备选服务) : 该头部列举了当前站点备选的访问方式列表, 一般用于在提供`QUIC`等新兴协议支持的同时, 实现向下兼容.

    经验证访问`https://localhost`(by f12)实际使用了h3.

    `ngx_http_xquic_module#readme`有说明xquic 依赖 ngx_http_v2_module, 直接访问https://localhost:2443会报错, 设置`add_header Alt-Svc 'h3-29=":2443"; ma=86400, h3=":2443"; ma=86400';`后访问`https://localhost`(by f12)实际使用了h3.

    > Alt-Svc + ssl_certificate=ca.pem, firefox访问`https://localhost`一直显示使用h2, 而`chrome报ERR_SSL_KEY_USAGE_INCOMPATIBLE`.

    > [对于 HTTP/1 来说，Alt-Svc 头部必须依附于首次响应，只有从第二个请求开始浏览器才会使用替代服务地址；而在 HTTP/2 中，ALTSVC 帧可以独立发送，浏览器从首次请求开始就能用上新地址](https://imququ.com/post/http-alt-svc.html)

## FAQ
### [http3为什么不基于UDP使用SCTP](https://http3-explained.haxx.se/zh/why-quic/why-tcpudp)
SCTP是一个支持数据流的可靠的传输层协议，而且在WebRTC上已有基于UDP的对它的实现.

这看上去很好，但与QUIC相比还不够好，它：
- 没有解决数据流的队头阻塞问题
- 连接建立时需要决定数据流的数量
- 没有稳固的TLS/安全性支持
- 建立连接时候需要4次握手，而QUIC一次都不用（0-RTT）
- QUIC是类TCP的字节流，而SCTP是信息流（message-based）
- **QUIC连接支持IP地址迁移，SCTP不行**, 这对移动设备的支持很有优势.
- [协议僵化阻碍着sctp等新协议的部署](https://http3-explained.haxx.se/zh/why-quic/why-ossification)

若要了解更多SCTP与QUIC的差异，请参阅[A Comparison between SCTP and QUIC](https://http3-explained.haxx.se/zh/why-quic/why-tcpudp).

### quic实现在user space
目前已知的所有QUIC实现都位于用户空间，这使它能得到更快速的迭代（相较于内核空间中的实现）.

### [验证站点是否启用HTTP3](https://http3check.net/)

htt3 supported:
- https://cloudflare-quic.com