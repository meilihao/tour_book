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

###Trailer

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

## 请求首部字段

### Accept

客户端能够处理的媒体类型及媒体类型的相对优先级.可使用 type/subtype 这种形式，一次指定多种媒体类型.

星号"*"用于按范围将类型分组，用"*/*"指示可接受全部类型；用"type/*"指示可接受 type类型的所有子类型.

若想要给显示的媒体类型增加优先级，则使用 q= 来额外表示权重值 1，用分号（;）进行分隔.权重值 q 的范围是 0~1（可精确到小数点后 3 位），且 1 为最大值。不指定权重 q 值时，默认权重为 q=1.0.当服务器提供多种内容时，将会首先返回权重值最高的媒体类型.

`Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`可理解为`Accept: (text/html,application/xhtml+xml,application/xml;q=0.9),(*/*;q=0.8)`

常见几个媒体类型:

- 文本文件

    text/html, text/plain, text/css ...
    application/xhtml+xml, application/xml ...

- 图片文件

    image/jpeg, image/gif, image/png ...

- 视频文件

    video/mpeg, video/quicktime ...

- 应用程序使用的二进制文件

    application/octet-stream, application/zip ...

### Accept-Charset

客户端支持的字符集及字符集的相对优先顺序,可一次性指定多种字符集.与首部字段 Accept 相同的是可用权重 q 值来表示相对优先级.

该首部字段应用于内容协商机制的服务器驱动协商.

### Accept-Encoding

支持的内容编码及内容编码的优先级顺序,可一次性指定多种内容编码.

常见编码:

- gzip

由文件压缩程序 gzip（GNU zip）生成的编码格式（RFC1952），采用 Lempel-Ziv 算法（LZ77）及 32 位循环冗余校验（Cyclic Redundancy Check，通称 CRC）。

- compress

由 UNIX 文件压缩程序 compress 生成的编码格式，采用 Lempel-Ziv-Welch 算法（LZW）。

- deflate

组合使用 zlib 格式（RFC1950）及由 deflate 压缩算法（RFC1951）生成的编码格式。

- identity

不执行压缩或不会变化的默认编码格式

采用权重 q 值来表示相对优先级，这点与首部字段 Accept 相同.另外，也可使用星号（*）作为通配符，指定任意的编码格式.

### Accept-Language

告知服务器用户代理能够处理的自然语言集（指中文或英文等），以及自然语言集的相对优先级,可一次指定多种自然语言集.

和 Accept 首部字段一样，按权重值 q 来表示相对优先级.

### Authorization

用来告知服务器，用户代理的认证信息（证书值）.通常，想要通过服务器认证的用户代理会在接收到返回的 401 状态码响应后，把首部字段 Authorization 加入请求中.共用缓存在接收到含有 Authorization 首部字段的请求时的操作处理会略有差异.

### Expect

客户端使用首部字段 Expect 来告知服务器，期望出现的某种特定行为.因服务器无法理解客户端的期望作出回应而发生错误时，会返回状态码 417 Expectation Failed.

### From

首部字段 From 用来告知服务器使用用户代理的用户的电子邮件地址.通常，其使用目的就是为了显示搜索引擎等用户代理的负责人的电子邮件联系方式。使用代理时，应尽可能包含 From 首部字段（但可能会因代理不同，将电子邮件地址记录在 User-Agent 首部字段内）.

### Host

请求的资源所处的互联网主机名和端口号.Host 首部字段在 HTTP/1.1 规范内是唯一一个必须被包含在请求内的首部字段.

首部字段 Host 和以单台服务器分配多个域名的虚拟主机的工作机制有很密切的关联，这是首部字段 Host 必须存在的意义.

### If-Match

形如 If-xxx 这种样式的请求首部字段，都可称为条件请求.服务器接收到附带条件的请求后，只有判断指定条件为真时，才会执行请求.

首部字段 If-Match，属附带条件之一，它会告知服务器匹配资源所用的实体标记（ETag）值,只有当 If-Match 的字段值跟 ETag 值匹配一致时，服务器才会接受请求;反之，则返回状态码 412 Precondition Failed 的响应.也可使用星号（*）指定 If-Match 的字段值,此时服务器将会忽略 ETag 的值，只要资源存在就处理请求.

### If-Modified-Since

首部字段 If-Modified-Since，属附带条件之一，如果在 If-Modified-Since 字段指定的日期时间后，资源发生了更新，服务器会接受请求.而在指定 If-Modified-Since 字段值的日期时间之后，如果请求的资源都没有过更新，则返回状态码 304 Not Modified 的响应.

If-Modified-Since 用于确认代理或客户端拥有的本地资源的有效性.获取资源的更新日期时间，可通过确认首部字段 Last-Modified 来确定.

### If-None-Match

首部字段 If-None-Match 属于附带条件之一.它和首部字段 If-Match 作用相反.用于指定 If-None-Match 字段值的实体标记（ETag）值与请求资源的 ETag 不一致时，它就告知服务器处理该请求.

在 GET 或 HEAD 方法中使用首部字段 If-None-Match 可获取最新的资源.因此，这与使用首部字段 If-Modified-Since 时有些类似.

### If-Range

首部字段 If-Range 属于附带条件之一.它告知服务器若指定的 If-Range 字段值（ETag 值或者时间）和请求资源的 ETag 值或时间相一致时，则作为范围请求处理.反之，则返回全体资源.

### If-Unmodified-Since

首部字段 If-Unmodified-Since 和首部字段 If-Modified-Since 的作用相反.它的作用的是告知服务器，指定的请求资源只有在字段值内指定的日期时间之后，未发生更新的情况下，才能处理请求.如果在指定日期时间后发生了更新，则以状态码 412 Precondition Failed 作为响应返回.

### Max-Forwards

通过 TRACE 方法或 OPTIONS 方法，发送包含首部字段 Max-Forwards 的请求时，该字段以十进制整数形式指定可经过的服务器最大数目.服务器在往下一个服务器转发请求之前，Max-Forwards 的值减 1 后重新赋值.当服务器接收到 Max-Forwards 值为 0 的请求时，则不再进行转发，而是直接返回响应.

### Proxy-Authorization

接收到从**代理服务器**发来的认证质询时，客户端会发送包含首部字段 Proxy-Authorization 的请求，以告知服务器认证所需要的信息.

### Range

对于只需获取部分资源的范围请求，即可告知服务器资源的指定范围.

接收到附带 Range 首部字段请求的服务器，会在处理请求之后返回状态码为 206 Partial Content 的响应.无法处理该范围请求时，则会返回状态码 200 OK 的响应及全部资源.

### Referer

首部字段 Referer 会告知服务器请求的原始资源的 URI,不推荐使用,可能导致安全问题.

### TE

首部字段 TE 会告知服务器客户端能够处理响应的｀传输编码｀方式及相对优先级．它和首部字段 Accept-Encoding（内容编码） 的功能很相像，但是用于传输编码．

ＴＥ还可以指定伴随 trailer 字段的分块传输编码的方式，应用后者时，只需把 trailers 赋值给该字段值．

### User-Agent

将创建请求的浏览器和用户代理名称等信息传达给服务器.

由网络爬虫发起请求时，有可能会在字段内添加爬虫作者的电子邮件地址.此外，如果请求经过代理，那么中间也很可能被添加上代理服务器的名称.

## 响应首部字段

### Accept-Ranges

用来告知客户端服务器是否能处理范围请求，以指定获取服务器端某个部分的资源.可指定的字段值有两种，可处理范围请求时指定其为 bytes，反之则指定其为 none.

### Age

告知客户端，源服务器在多久前创建了响应。字段值的单位为秒.

若创建该响应的服务器是缓存服务器，Age 值是指缓存后的响应再次发起认证到认证完成的时间值.代理创建响应时必须加上首部字段 Age.

### ETag

告知客户端实体标识,它是一种可将资源以字符串形式做唯一性标识的方式.服务器会为每份资源分配对应的 ETag 值，当资源更新时，ETag 值也需要更新.生成 ETag 值时，并没有统一的算法规则，而仅仅是由服务器来分配.

google中英文首页的 URI 是相同的，所以仅凭 URI 指定缓存的资源是相当困难的．若在下载过程中出现连接中断、再连接的情况，都会依照 ETag 值来指定资源．

强 ETag 值和弱 Tag 值：

－　强 ETag 值，不论实体发生多么细微的变化都会改变其值
－　弱 ETag 值只用于提示资源是否相同．只有资源发生了根本改变，产生差异时才会改变 ETag 值，需在字段值最开始处附加｀W/｀

### Location

将响应接收方引导至某个与请求 URI 位置不同的资源.基本上，该字段会配合 3xx ：Redirection 的响应，提供重定向的 URI.几乎所有的浏览器在接收到包含首部字段 Location 的响应后，都会强制性地尝试对已提示的重定向资源的访问.

### Proxy-Authenticate

把由代理服务器所要求的认证信息发送给客户端.

它与客户端和服务器之间的 HTTP 访问认证的行为相似，不同之处在于其认证行为是在客户端与代理之间进行的.而客户端与服务器之间进行认证时，首部字段 WWW-Authorization 有着相同的作用.

### Retry-After

告知客户端应该在多久之后再次发送请求.主要配合状态码 503 Service Unavailable 响应，或 3xx Redirect 响应一起使用.字段值可以指定为具体的日期时间（Wed, 04 Jul 2012 06：34：24 GMT 等格式），也可以是创建响应后的秒数.

### Server

告知客户端当前服务器上安装的 HTTP 服务器应用程序的信息.不单单会标出服务器上的软件应用名称，还有可能包括版本号和安装时启用的可选项.但因安全问题,通常只包含server名称和版本.

### Vary

可对缓存进行控制.源服务器会向代理服务器传达关于本地缓存使用方法的命令，可以告诉缓存服务器依据客户端的请求头中的哪一项或者哪几项来进行缓存．

从代理服务器接收到源服务器返回包含 Vary 指定项的响应之后，若再要进行缓存，仅对请求中含有相同 Vary 指定首部字段的请求返回缓存.即使对相同资源发起请求，但由于 Vary 指定的首部字段不相同，因此必须要从源服务器重新获取资源.

比如我们有一份文档是多国语言版本的，但是资源的地址却是一个，"/doc/download/13"
客户端第一次请求这个地址，由于缓存服务器上没有这个文档的缓存，所以直接将请求转发到了原始服务器，原始服务器根据用户的accept-language请求头得到用户期望得到的语言，比如是｀fr, en;q=0.8｀．则原始服务器返回整个文档给缓存服务器，并且带上vray: accept-language，这样缓存服务器在以后收到同样accept-language为fr的请求时就不会再去原始服务器上拿原文档了，而是直接响应fr法文的文档给客户端，但是如果用户想要的是jp日文呢，则原始服务器会重复之前的动作，获取原文档，并且根据vary的头部accept-language来进行缓存该文档．

### WWW-Authenticate

用于 HTTP 访问认证.它会告知客户端适用于访问请求 URI 所指定资源的认证方案（Basic 或是 Digest）和带参数提示的质询（challenge）.状态码 401 Unauthorized 响应中，肯定带有首部字段 WWW-Authenticate.

## 实体首部

### Allow

用于通知客户端能够支持 Request-URI 指定资源的所有 HTTP 方法.当服务器接收到不支持的 HTTP 方法时，会以状态码 405 Method Not Allowed 作为响应返回.

### Content-Encoding

告知客户端服务器对实体的主体部分选用的内容编码方式.内容编码是指在不丢失实体信息的前提下所进行的压缩.

### Content-Language

告知客户端，实体主体使用的自然语言（指中文或英文等语言）.

### Content-Length

实体主体部分的大小（单位是字节）.**对实体主体进行内容编码传输时，不能再使用 Content-Length 首部字段**.

### Content-Location

给出与报文主体部分相对应的 URI.和首部字段 Location 不同，Content-Location 表示的是报文主体返回资源对应的 URI.

### Content-MD5

检查报文主体在传输过程中是否保持完整，以及确认传输到达.

采用这种方法，对内容上的偶发性改变是无从查证的，也无法检测出恶意篡改.

### Content-Range

针对范围请求，返回响应时使用的首部字段 Content-Range，能告知客户端作为响应返回的实体的哪个部分符合范围请求.字段值以字节为单位，表示当前发送部分及整个实体大小.

### Content-Type

Content-Type 说明了实体主体内对象的媒体类型.和首部字段 Accept 一样，字段值用 type/subtype 形式赋值,参数 charset常用utf-8.

### Expires

将资源失效的日期告知客户端.缓存服务器在接收到含有首部字段 Expires 的响应后，会以缓存来应答请求，在 Expires 字段值指定的时间之前，响应的副本会一直被保存.当超过指定的时间后，缓存服务器在请求发送过来时，会转向源服务器请求资源.

源服务器不希望缓存服务器对资源缓存时，最好在 Expires 字段内写入与首部字段 Date 相同的时间值.

但是，当首部字段 Cache-Control 有指定 max-age 指令时，比起首部字段 Expires，会优先处理 max-age 指令.

### Last-Modified

指明资源最终修改的时间。一般来说，这个值就是 Request-URI 指定资源被修改的时间．

## Cookie首部

Cookie 的工作机制是用户识别及状态管理．

### 为 Cookie 服务的首部字段

首部字段名	说明	首部类型
Set-Cookie	开始状态管理所使用的Cookie信息	响应首部字段
Cookie	服务器接收到的Cookie信息	请求首部字段


### Set-Cookie 字段的属性

属性	说明
NAME=VALUE	赋予 Cookie 的名称和其值（必需项）
expires=DATE	Cookie 的有效期（若不明确指定则默认为浏览器关闭前为止）
path=PATH	将服务器上的文件目录作为Cookie的适用对象（若不指定则默认为文档所在的文件目录）
domain=域名	作为 Cookie 适用对象的域名 （若不指定则默认为创建 Cookie 的服务器的域名）
Secure	仅在 HTTPS 安全通信时才会发送 Cookie
HttpOnly	加以限制，使 Cookie 不能被 JavaScript 脚本访问

#### expires 属性

一旦 Cookie 从服务器端发送至客户端，服务器端就不存在可以显式删除 Cookie 的方法.但可通过覆盖**已过期**的 Cookie，实现对客户端 Cookie 的实质性删除操作

#### path 属性

Cookie 的 path 属性可用于限制指定 Cookie 的发送范围的文件目录.不过另有办法可避开这项限制.一般设为"/",以表示同一个站点的所有页面都可以访问这个Cookie.

#### domain 属性

通过 Cookie 的 domain 属性指定的域名可做到与结尾匹配一致.比如，当指定 example.com 后，除 example.com 以外，www.example.com 或 www2.example.com 等都可以发送 Cookie。

因此，除了针对具体指定的多个域名发送 Cookie 之外，不指定 domain 属性显得更安全.

#### secure 属性

Cookie 的 secure 属性用于限制 Web 页面仅在 HTTPS 安全连接时，才可以发送 Cookie．

发送 Cookie 时，指定 secure 属性的方法如下所示：

Set-Cookie: name=value; secure

以上例子仅当在 https://www.example.com/（HTTPS）安全连接的情况下才会进行 Cookie 的访问.也就是说，即使域名相同，http://www.example.com/（HTTP）也不会发生 Cookie被访问．

当省略 secure 属性时，不论 HTTP 还是 HTTPS，都会对 Cookie 进行回收．

#### HttpOnly 属性

Cookie 的 HttpOnly 属性是 Cookie 的扩展功能，它使 JavaScript 脚本无法获得 Cookie.其主要目的为防止跨站脚本攻击（Cross-site scripting，XSS）对 Cookie 的信息窃取.该扩展并非是为了防止 XSS 而开发的．

发送指定 HttpOnly 属性的 Cookie 的方法如下所示:

`Set-Cookie: name=value; HttpOnly`

通过上述设置，通常从 Web 页面内还可以对 Cookie 进行读取操作.但使用 JavaScript 的 document.cookie 就无法读取附加 HttpOnly 属性后的 Cookie 的内容了.因此，也就无法在 XSS 中利用 JavaScript 劫持 Cookie 了.

## 其他首部

HTTP 首部字段是可以自行扩展的.一些最为常用的首部字段:

- X-Frame-Options
- X-XSS-Protection
- DNT
- P3P

### X-Frame-Options

属于 HTTP 响应首部，用于控制网站内容在其他 Web 网站的 Frame 标签内的显示问题.其主要目的是为了防止点击劫持（clickjacking）攻击.

> 由于frame框架对网页可用性存在负面影响，因此在HTML 5中已经不再支持frame框架，但是支持iframe框架所以HTML 5中> > 废弃了frame框架的`frameset`、`frame`和`noframes`标签.

X-Frame-Options 有以下两个可指定的字段值。

- DENY ：拒绝
- SAMEORIGIN ：仅同源域名下的页面（Top-level-browsing-context）匹配时许可。（比如，当指定 http://hackr.jp/sample.html 页面为 SAMEORIGIN 时，那么 hackr.jp 上所有页面的 frame 都被允许可加载该页面，而 example.com 等其他域名的页面就不行了）

### X-XSS-Protection

属于 HTTP 响应首部，它是针对跨站脚本攻击（XSS）的一种对策，用于控制浏览器 XSS 防护机制的开关.

X-XSS-Protection 可指定的字段值如下:

- 0 ：将 XSS 过滤设置成无效状态
- 1 ：将 XSS 过滤设置成有效状态

### DNT

属于 HTTP 请求首部，其中 DNT 是 Do Not Track 的简称，意为拒绝个人信息被收集，是表示拒绝被精准广告追踪的一种方法.

DNT 可指定的字段值如下:

- 0 ：同意被追踪
- 1 ：拒绝被追踪

由于首部字段 DNT 的功能具备有效性，所以 Web 服务器需要对 DNT 做对应的支持.

### P3P

首部字段 P3P 属于 HTTP 相应首部，通过利用 P3P（The Platform for Privacy Preferences，在线隐私偏好平台）技术，可以让 Web 网站上的个人隐私变成一种仅供程序可理解的形式，以达到保护用户隐私的目的.