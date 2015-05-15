# web基础

## 网络基础TCP/IP

### TCP/IP协议族分层(四层模型)

|分层|协议|作用|
|--------|--------|--------|
|应用层|http,ftp,dns|向应用提供服务|
|传输层|tcp,udp|提供处于网络连接中的计算机间的数据传输|
|网络层|ip,icmp,arp|处理在网络中传输的数据包|
|链路层|略|处理连接网络的硬件部分|

数据包在传输时,每经过一层,必定会被打上(或消去)该层所属的首部信息.

IP间通信依赖MAC地址.

ARP是一种用以解析地址的协议,可根据通信方的IP地址反查出对应的MAC地址.

路由选择是指选择通过互连网络从源节点向目的节点传输信息的机制，而且信息至少通过一个中间节点。路由选择工作在 OSI 参考模型的网络层.

TCP是面向连接的,可靠的数据流(或叫字节流)服务.其发起连接需[三次握手](http://baike.baidu.com/view/1003841.htm),会使用`SYN(synchronize)`和`ACK(acknowledgement)`标志,状态变化:

客户端:`SYN_SEND`(发送`syn[j]`包后)->`ESTABLISHED`(收到服务器`SYN[k]+ACK[j+1]`包,再向服务器发送`ACK[k+1]`后)

服务端:`SYN_RECV`(收到客户端的`syn[j]`,发送`自己的SYN[k]`+`ACK[j+1]`包后)->`ESTABLISHED`(收到客户端`ACK[K+1]`)

### DNS

它提供域名到IP地址之间的解析服务.

### URI和URL

统一资源标识符（Uniform Resource Identifier，或URI)是一个用于标识某一互联网资源名称的字符串.

`RFC3986：统一资源标识符（URI）通用语法`中列举了几种 URI 例子，如下所示:

```
ftp://ftp.is.co.za/rfc/rfc1808.txt
http://www.ietf.org/rfc/rfc2396.txt
ldap://[2001:db8::7]/c=GB?objectClass?one
mailto:John.Doe@example.com
news:comp.infosystems.www.servers.unix
tel:+1-816-555-1212
telnet://192.0.2.16:80/
urn:oasis:names:specification:docbook:dtd:xml:4.1.2
```

统一资源定位符是对可以从互联网上得到的资源的位置和访问方法的一种简洁的表示，是互联网上标准资源的地址.

URL是URI的子集.

URI抽象结构:

```html
[scheme:]scheme-specific-part[#fragment]

[scheme:][//authority][path][?query][#fragment]

其中authority为[user-info@]host[:port]
```