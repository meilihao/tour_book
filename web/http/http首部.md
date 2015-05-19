# http首部

HTTP 协议的请求和响应报文中必定包含 HTTP 首部.首部内容为客户端和服务器分别处理请求和响应提供所需要的信息.

## HTTP 首部字段结构

HTTP 首部字段是由首部字段名和字段值构成的，中间用冒号“:” 分隔,例如`首部字段名: 字段值`.

## 4 种 HTTP 首部字段类型

HTTP 首部字段根据实际用途被分为以下 4 种类型:

### 通用首部字段（General Header Fields）

请求报文和响应报文两方都会使用的首部.

首部字段名	说明
Cache-Control	控制缓存的行为
Connection	逐跳首部、连接的管理
Date	创建报文的日期时间
Pragma	报文指令
Trailer	报文末端的首部一览
Transfer-Encoding	指定报文主体的传输编码方式
Upgrade	升级为其他协议
Via	代理服务器的相关信息
Warning	错误通知

### 请求首部字段（Request Header Fields）

从客户端向服务器端发送请求报文时使用的首部.补充了请求的附加内容、客户端信息、响应内容相关优先级等信息.

首部字段名	说明
Accept	用户代理可处理的媒体类型
Accept-Charset	优先的字符集
Accept-Encoding	优先的内容编码
Accept-Language	优先的语言（自然语言）
Authorization	Web认证信息
Expect	期待服务器的特定行为
From	用户的电子邮箱地址
Host	请求资源所在服务器
If-Match	比较实体标记（ETag）
If-Modified-Since	比较资源的更新时间
If-None-Match	比较实体标记（与 If-Match 相反）
If-Range	资源未更新时发送实体 Byte 的范围请求
If-Unmodified-Since	比较资源的更新时间（与If-Modified-Since相反）
Max-Forwards	最大传输逐跳数
Proxy-Authorization	代理服务器要求客户端的认证信息
Range	实体的字节范围请求
Referer	对请求中 URI 的原始获取方
TE	传输编码的优先级
User-Agent	HTTP 客户端程序的信息

### 响应首部字段（Response Header Fields）

从服务器端向客户端返回响应报文时使用的首部.补充了响应的附加内容，也会要求客户端附加额外的内容信息.

首部字段名	说明
Accept-Ranges	是否接受字节范围请求
Age	推算资源创建经过时间
ETag	资源的匹配信息
Location	令客户端重定向至指定URI
Proxy-Authenticate	代理服务器对客户端的认证信息
Retry-After	对再次发起请求的时机要求
Server	HTTP服务器的安装信息
Vary	代理服务器缓存的管理信息
WWW-Authenticate	服务器对客户端的认证信息

### 实体首部字段（Entity Header Fields）

针对请求报文和响应报文的实体部分使用的首部.补充了资源内容更新时间等与实体有关的信息.

首部字段名	说明
Allow	资源可支持的HTTP方法
Content-Encoding	实体主体适用的编码方式
Content-Language	实体主体的自然语言
Content-Length	实体主体的大小（单位：字节）
Content-Location	替代对应资源的URI
Content-MD5	实体主体的报文摘要
Content-Range	实体主体的位置范围
Content-Type	实体主体的媒体类型
Expires	实体主体过期的日期时间
Last-Modified	资源的最后修改日期时间

## End-to-end 首部和 Hop-by-hop 首部

HTTP 首部字段将定义成缓存代理和非缓存代理的行为，分成 2 种类型:

### 端到端首部（End-to-end Header）

分在此类别中的首部会转发给请求 / 响应对应的最终接收目标，且必须保存在由缓存生成的响应中，另外规定它**必须被转发**.

### 逐跳首部（Hop-by-hop Header）

分在此类别中的首部只对单次转发有效，会**因通过缓存或代理而不再转发**.HTTP/1.1 和之后版本中，如果要使用 hop-by-hop 首部，需提供 Connection 首部字段.

下面列举了 HTTP/1.1 中的逐跳首部字段,除这 8 个首部字段之外，其他所有字段都属于端到端首部.

- Connection
- Keep-Alive
- Proxy-Authenticate
- Proxy-Authorization
- Trailer
- TE
- Transfer-Encoding
- Upgrade

## HTTP/1.1 通用首部字段

### Cache-Control

通过指定首部字段 Cache-Control 的指令，就能操作缓存的工作机制.

#### 缓存请求指令

指令	参数	说明
no-cache	无	强制向源服务器再次验证
no-store	无	不缓存请求或响应的任何内容
max-age = [ 秒]	必需	响应的最大Age值
max-stale( = [ 秒])	可省略	接收已过期的响应
min-fresh = [ 秒]	必需	期望在指定时间内的响应仍有效
no-transform	无	代理不可更改媒体类型
only-if-cached	无	从缓存获取资源
cache-extension	-	新指令标记（token）

#### 缓存响应指令

指令	参数	说明
public	无	可向任意方提供响应的缓存
private	可省略	仅向特定用户返回响应
no-cache	可省略	缓存前必须先确认其有效性
no-store	无	不缓存请求或响应的任何内容
no-transform	无	代理不可更改媒体类型
must-revalidate	无	可缓存但必须再向源服务器进行确认
proxy-revalidate	无	要求中间缓存服务器对缓存的响应有效性再进行确认
max-age = [ 秒]	必需	响应的最大Age值
s-maxage = [ 秒]	必需	公共缓存服务器响应的最大Age值
cache-extension	-	新指令标记（token）

#### no-cache

使用 no-cache 指令的目的是为了防止从缓存中返回过期的资源.

客户端发送的请求中如果包含 no-cache 指令，则表示客户端将不会接收缓存过的响应.于是，“中间”的缓存服务器必须把客户端请求转发给源服务器。

如果服务器返回的响应中包含 no-cache 指令，那么缓存服务器不能对资源进行缓存.源服务器以后也将不再对缓存服务器请求中提出的资源有效性进行确认，且禁止其对响应资源进行缓存操作.

Cache-Control: no-cache=Location
由服务器返回的响应中，若报文首部字段 Cache-Control 中对 no-cache 字段名具体指定参数值，那么客户端在接收到这个被指定参数值的首部字段对应的响应报文后，就不能使用缓.换言之，无参数值的首部字段可以使用缓存.**只能在响应指令中指定该参数**.

#### no-store

no-store 指令

该指令规定缓存不能在本地存储请求或响应的任一部分.通常其暗示请求（和对应的响应）或响应中包含机密信息.

事实上 no-cache 代表不缓存过期的资源，缓存会向源服务器进行有效期确认后处理资源;no-store 才是真正地不进行缓存.

#### s-maxage

指定缓存期限和认证的指令,当使用 s-maxage 指令后，则直接忽略对 Expires 首部字段及 max-age 指令的处理.s-maxage 指令的功能和 max-age 指令的相同，它们的不同点是 s-maxage 指令只适用于供多位用户使用的公共缓存服务器(即一般代理).也就是说，对于向同一用户重复返回响应的服务器来说，这个指令没有任何作用.

#### max-age

当客户端发送的请求中包含 max-age 指令时，如果判定缓存资源的缓存时间数值比指定时间的数值更小，那么客户端就接收缓存的资源.另外，当指定 max-age 值为 0，那么缓存服务器通常需要将请求转发给源服务器.

当服务器返回的响应中包含 max-age 指令时，缓存服务器将不对资源的有效性再作确认，而 max-age 数值代表资源保存为缓存的最长时间.

HTTP/1.1 版本的缓存服务器遇到同时存在 Expires 首部字段的情况时，会优先处理 max-age 指令，而忽略掉 Expires 首部字段.

#### min-fresh

min-fresh 指令要求缓存服务器返回至少还未过指定时间的缓存资源.

#### max-stale

使用 max-stale 可指示缓存资源，即使过期时间内也照常接收.

#### only-if-cached

客户端仅在缓存服务器本地缓存目标资源的情况下才会要求其返回(不确认资源有效性).若发生请求缓存服务器的本地缓存无响应，则返回状态码 504 Gateway Timeout.

#### must-revalidate

代理会向源服务器再次验证即将返回的响应缓存目前是否仍然有效,若代理无法连通源服务器再次获取有效资源的话，缓存必须给客户端一条 504（Gateway Timeout）状态码.且该指令会忽略请求的 max-stale 指令.

#### proxy-revalidate

要求所有的缓存服务器在接收到客户端带有该指令的请求返回响应之前，必须再次验证缓存的有效性.

#### no-transform

规定无论是在请求还是响应中，缓存都不能改变实体主体的媒体类型.这样做可防止缓存或代理压缩图片等类似操作.

### Connection

Connection 首部字段具备如下两个作用:

- 控制不再转发给下一个连接的首部字段
- 管理持久连接

#### 控制不再转发给下一连接的首部字段

在客户端发送请求和服务器返回响应内，使用 Connection 首部字段，可控制不再转发给下一个连接的首部字段（即 Hop-by-hop 首部）(即会删除connection首部和它指定的首部).

#### 管理持久连接

- `Connection: close`

HTTP/1.1 版本的默认连接都是持久连接。为此，客户端会在持久连接上连续发送请求.当服务器端想明确断开连接时，则指定 Connection 首部字段的值为 Close.

- `Connection: Keep-Alive`

**HTTP/1.1 之前的 HTTP 版本的默认连接都是非持久连接.**为此，如果想在旧版本的 HTTP 协议上维持持续连接，则需要指定 Connection 首部字段的值为 Keep-Alive.

### Date

首部字段 Date 表明创建 HTTP 报文的日期和时间.

HTTP/1.1 协议使用在 RFC1123 中规定的日期时间的格式:`Date: Tue, 03 Jul 2012 04:40:59 GMT`.

###　Trailer

首部字段 Trailer 会事先说明在报文主体后记录了哪些首部字段．该首部字段可应用在 HTTP/1.1 版本分块传输编码时．

### Transfer-Encoding

首部字段 Transfer-Encoding 规定了传输报文主体时采用的编码方式.

### Upgrade

用于检测 HTTP 协议及其他协议是否可使用更高的版本进行通信，其参数值可以用来指定一个完全不同的通信协议.

Upgrade 首部字段产生作用的 Upgrade 对象仅限于客户端和邻接服务器之间.因此，使用首部字段 Upgrade 时，还需要额外指定 Connection:Upgrade(以便在下一个连接转发前删除该信息).对于附有首部字段 Upgrade 的请求，服务器可用 101 Switching Protocols 状态码作为响应返回.

### Via

使用首部字段 Via 是为了追踪客户端与服务器之间的请求和响应报文的传输路径.

报文经过代理或网关时，会先在首部字段 Via 中附加该服务器的信息，然后再进行转发.这个做法和 traceroute 及电子邮件的 Received 首部的工作机制很类似.

首部字段 Via 不仅用于追踪报文的转发，还可避免请求回环的发生.所以必须在经过代理时附加该首部字段内容.

Via 首部是为了追踪传输路径，所以经常会和 TRACE 方法一起使用.比如，代理服务器接收到由 TRACE 方法发送过来的请求（其中 Max-Forwards: 0）时，代理服务器就不能再转发该请求了.这种情况下，代理服务器会将自身的信息附加到 Via 首部后，返回该请求的响应.

### Warning

通常会告知用户一些与缓存相关的问题的警告.