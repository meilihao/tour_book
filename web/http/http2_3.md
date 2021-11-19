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
HTTP/3不存在明文的不安全版本.

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