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
QUIC没有非加密的版本.

跟 TLS 类似，QUIC 的 0RTT 握手，是建立在已经同一个服务器建立过连接的基础上，所以如果是纯的第一次连接，仍然需要一个 RTT.

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