# 简单的http协议

用于 HTTP 协议交互的信息被称为 HTTP 报文.请求端（客户端）的 HTTP 报文叫做请求报文，响应端（服务器端）的叫做响应报文.HTTP 报文本身是由多行（用 CR+LF 作换行符）数据构成的字符串文本.

HTTP 报文大致可分为报文首部和报文主体两块.两者由最初出现的空行（CR+LF）来划分.通常，并不一定要有报文主体.报文首部是服务器或客户端需处理的请求或响应的内容及属性.报文主体是应被发送的数据.

请求报文是由请求行(请求方法、请求 URI、协议版本)、可选的请求首部字段和内容实体构成的.

响应报文基本上由响应行(协议版本、状态码（表示请求成功或失败的数字代码）、原因短语(用以解释状态码))、可选的响应首部字段以及实体主体构成.

首部字段是指包含表示请求和响应的各种条件和属性的各类首部，一般有 4 种首部，分别是：通用首部、请求首部、响应首部和实体首部。

> 报文主体和实体主体的差异：
>
> 报文（message）
>
> 是 HTTP 通信中的基本单位，通过 HTTP 通信传输。
>
> 实体（entity）
>
> 作为请求或响应的有效载荷数据（补充项）被传输，其内容由实体首部和实体主体组成。
>
> HTTP 报文的主体用于传输请求或响应的实体主体。
>
> 通常，报文主体等于实体主体。只有当传输中进行编码操作时，实体主体的内容发生变化，才导致它和报文主体产生差异。

http是**无状态**协议,但**Cookie技术可保持状态**.

## http方法

```
GET      请求获取Request-URI所标识的资源
POST     在Request-URI所标识的资源后附加新的数据
HEAD     请求获取由Request-URI所标识的资源的响应消息报头(用以确认URI的有效性及资源更新的日期时间和通信状态等)
PUT      请求服务器存储一个资源，并用Request-URI作为其标识(**本身不带验证机制,需自实现**)
DELETE   请求服务器删除Request-URI所标识的资源(同PUT**本身不带验证机制**)
TRACE    请求服务器回送收到的请求信息，主要用于测试或诊断
CONNECT  用隧道协议连接代理
OPTIONS  请求查询服务器的性能，或者查询与资源相关的选项和需求(通常是查询支持的方法)
```
### CONNECT

CONNECT方法要求在与代理服务器通信时建立隧道，实现用隧道协议进行 TCP 通信.主要使用 SSL（Secure Sockets Layer，安全套接层）和 TLS（Transport Layer Security，传输层安全）协议把通信内容加密后经网络隧道传输

CONNECT 方法的格式:`CONNECT 代理服务器名:端口号 HTTP版本`.

### TRACE

追踪路径.

发送请求时，在 Max-Forwards 首部字段中填入数值，每经过一个服务器端就将该数字减 1，当数值刚好减到 0 时，就停止继续传输，最后接收到请求的服务器端则返回状态码 200 OK 的响应.

客户端通过 TRACE 方法可以查询发送出去的请求是怎样被加工修改 / 篡改的。这是因为，请求想要连接到源目标服务器可能会通过代理中转，TRACE 方法就是用来确认连接过程中发生的一系列操作.

但是，TRACE 方法本来就不怎么常用，再加上它容易引发 XST（Cross-Site Tracing，跨站追踪）攻击，通常就更不会用到了.

## 持久连接

HTTP/1.1 默认支持持久连接（HTTP Persistent Connections，也称为 HTTP keep-alive 或 HTTP connection reuse）的方法.持久连接的特点是，只要任意一端没有明确提出断开连接，则保持 TCP 连接状态.

持久连接旨在建立 1 次 TCP 连接后进行多次请求和响应的交互.它的好处在于减少了 TCP 连接的重复建立和断开所造成的额外开销，减轻了服务器端的负载,提升 Web 页面的响应速度.

## 管线化

持久连接使得多数请求以管线化（pipelining）方式发送成为可能.从前发送请求后需等待并收到响应，才能发送下一个请求.管线化后就能够做到同时并行发送多个请求，而传送过程中不需先等待服务端的回应,即不用等待响应亦可直接发送下一个请求.

## Cookie

Cookie 技术通过在请求和响应报文中写入 Cookie 信息来控制客户端的状态.

Cookie 会根据从服务器端发送的响应报文内的一个叫做 Set-Cookie 的首部字段信息，通知客户端保存 Cookie.当下次客户端再往该服务器发送请求时，客户端会自动在请求报文中加入 Cookie 值.

## http编码

用于提升传输速率.

常用的内容编码有以下几种：

gzip（GNU zip）

compress（UNIX 系统的标准压缩）

deflate（zlib）

identity（不进行编码）

## 分割发送的分块传输编码

在 HTTP 通信过程中，请求的编码实体资源尚未全部传输完成之前，浏览器无法显示请求页面.在传输大容量数据时，通过把数据分割成多块，能够让浏览器逐步显示页面.

这种把实体主体分块的功能称为分块传输编码（Chunked Transfer Coding）.

## MIME

MIME（Multipurpose Internet Mail Extensions，多用途因特网邮件扩展）机制，它允许邮件处理文本、图片、视频等多个不同类型的数据.在 MIME 扩展中会使用一种称为多部分对象集合（Multipart）的方法，来容纳多份不同类型的数据,通常在上传时使用.

Multipart包括：

- multipart/form-data

 在 Web 表单文件上传时使用.

 ```
Content-Type: multipart/form-data; boundary=AaB03x
　
--AaB03x
Content-Disposition: form-data; name="field1"
　
Joe Blow
--AaB03x
Content-Disposition: form-data; name="pics"; filename="file1.txt"
Content-Type: text/plain
　
...（file1.txt的数据）...
--AaB03x--
```

- multipart/byteranges

 状态码 206（Partial Content，部分内容）响应报文包含了多个范围的内容时使用.类似于 FlashGet 或者迅雷这类的 HTTP 下载工具都是使用此类响应实现断点续传或者将一个大文档分解为多个下载段同时下载.

```
HTTP/1.1 206 Partial Content
Date: Fri, 13 Jul 2012 02:45:26 GMT
Last-Modified: Fri, 31 Aug 2007 02:02:20 GMT
Content-Type: multipart/byteranges; boundary=THIS_STRING_SEPARATES


--THIS_STRING_SEPARATES
Content-Type: application/pdf
Content-Range: bytes 500-999/8000


...（范围指定的数据）...
--THIS_STRING_SEPARATES
Content-Type: application/pdf
Content-Range: bytes 7000-7999/8000


...（范围指定的数据）...
--THIS_STRING_SEPARATES--
```

在 HTTP 报文中使用多部分对象集合时，需要在首部字段里加上 Content-type.

使用 boundary 字符串来划分多部分对象集合指明的各类实体.在 boundary 字符串指定的各个实体的起始行之前插入“--”标记（例如：--AaB03x、--THIS_STRING_SEPARATES），而在多部分对象集合对应的字符串的最后插入“--”标记（例如：--AaB03x--、--THIS_STRING_SEPARATES--）作为结束.

多部分对象集合的每个部分类型中，都可以含有首部字段.另外，可以在某个部分中嵌套使用多部分对象集合.

## 获取部分内容的范围请求

想实现能从之前下载中断处恢复下载,就需要指定下载的实体范围.像这样，指定范围发送的请求叫做范围请求（Range Request）.

```
Get /logo.jpg HTTP/1.1
Host: XXX
Range: bytes =5001-10000
-------------------------
HTTP/1.1 206 Partial Content
Data: XXX
Content-Range: bytes =5001-10000/10000
Content-Length: 5000
Content-Type:image/jpeg
```

针对范围请求，响应会返回状态码为 206 Partial Content 的响应报文．另外，对于多重范围的范围请求，响应会在首部字段 Content-Type 标明 multipart/byteranges 后返回响应报文(见上面MIME中的例子)．

如果服务器端无法响应范围请求，则会返回｀状态码 200 OK｀和｀完整的实体内容｀．

## 内容协商返回最适合的内容

同一个 Web 网站有可能存在着多份相同内容的页面,比如google不同版本的首页.

当浏览器的默认语言为英语或中文，访问相同 URI 的 Web 页面时，则会显示对应的英语版或中文版的 Web 页面.这样的机制称为内容协商（Content Negotiation）.

内容协商机制是指客户端和服务器端就响应的资源内容进行交涉，然后提供给客户端最为适合的资源.内容协商会以响应资源的语言、字符集、编码方式等作为判断的基准.

包含在请求报文中的某些首部字段（如下）就是判断的基准:

- Accept
- Accept-Charset
- Accept-Encoding
- Accept-Language
- Content-Language

内容协商技术有以下 3 种类型:

- 服务器驱动协商（Server-driven Negotiation）

 由服务器端进行内容协商。以请求的首部字段为参考，在服务器端自动处理。但对用户来说，以浏览器发送的信息作为判定的依据，并不一定能筛选出最优内容。

- 客户端驱动协商（Agent-driven Negotiation）

 由客户端进行内容协商的方式。用户从浏览器显示的可选项列表中手动选择。还可以利用 JavaScript 脚本在 Web 页面上自动进行上述选择。比如按 OS 的类型或浏览器类型，自行切换成 PC 版页面或手机版页面。

- 透明协商（Transparent Negotiation）

 是服务器驱动和客户端驱动的结合体，是由服务器端和客户端各自进行内容协商的一种方法。

## 返回的HTTP状态码








