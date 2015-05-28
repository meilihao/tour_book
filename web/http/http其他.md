## http瓶颈

- 一条连接上只可发送一个请求。
- 请求只能从客户端开始。客户端不可以接收除响应以外的指令。
- 请求 / 响应首部未经压缩就发送。首部信息越多延迟越大。
- 发送冗长的首部。每次互相发送相同的首部造成的浪费较多。
- 可任意选择数据压缩格式。非强制压缩发送。一条连接上只可发送一个请求。

## Ajax

Ajax（Asynchronous JavaScript and XML， 异 步 JavaScript 与 XML 技术）是一种有效利用 JavaScript 和 DOM（Document Object Model，文档对象模型）的操作，以达到局部 Web 页面替换加载的异步通信手段。和以前的同步通信相比，由于它只更新一部分页面，响应中传输的数据量会因此而减少，这一优点显而易见。

Ajax 的核心技术是名为 XMLHttpRequest 的 API，通过 JavaScript 脚本语言的调用就能和服务器进行 HTTP 通信。借由这种手段，就能从已加载完毕的 Web 页面上发起请求，只更新局部页面。

而利用 Ajax 实时地从服务器获取内容，有可能会导致大量请求产生.

## commet

Comet是一种用于web的推送技术，能使服务器实时地将更新的信息传送到客户端，而无须客户端发出请求.***Comet将被HTML5标准中的WebSocket取代**.

Comet 架构非常适合**事件驱动**的 Web 应用，以及对**交互性和实时性**要求很强的应用，如股票交易行情分析、聊天室和 Web 版在线游戏等。

基于 Comet 架构的 Web 应用使用客户端和服务器端之间的 HTTP 长连接来作为数据传输的通道。每当服务器端的数据因为外部的事件而发生改变时，服务器端就能够及时把相关的数据推送给客户端。通常来说，有两种实现长连接的策略：

- HTTP 流（HTTP Streaming）/iframe流
这种情况下，客户端打开一个单一的与服务器端的 HTTP 持久连接。服务器通过此连接把数据发送过来，客户端增量的处理它们。
具体实现是在页面中插入一个隐藏的iframe，利用其src属性在服务器和客户端之间创建一条长链接，服务器向iframe传输数据（通常是HTML，内有负责插入信息的javascript），来实时更新页面。

iframe流方式的优点是浏览器兼容好，Google公司在一些产品中使用了iframe流，如Google Talk。

- HTTP 长轮询（HTTP Long Polling）
这种情况下，由客户端向服务器端发出请求并打开一个连接。这个连接只有在收到服务器端的数据之后才会关闭。服务器端发送完数据之后，就立即关闭连接。客户端则马上再打开一个新的连接，等待下一次的数据。

## WebSocket

WebSocket，即 Web 浏览器与 Web 服务器之间全双工通信标准。

WebSocket建立在 HTTP 基础上的协议，因此连接的发起方仍是客户端.

WebSocket 协议的主要特点:

- 推送功能
支持由服务器向客户端推送数据的推送功能。这样，服务器可直接发送数据，而不必等待客户端的请求.

- 减少通信量
只要建立起 WebSocket 连接，就希望一直保持连接状态。和 HTTP 相比，不但每次连接时的总开销减少，而且由于WebSocket 的首部信息很小，通信量也相应减少了.

### WebSocket握手

- 请求
用 HTTP 的 Upgrade 首部字段，告知服务器通信协议发生改变.

- 响应
返回状态码 101 Switching Protocols 的响应.

WebSocket成功握手确立 WebSocket 连接之后，通信时不再使用 HTTP 的数据帧，而采用 WebSocket 独立的数据帧.

## JSON

JSON（JavaScript Object Notation）是一种以 JavaScript（ECMAScript）的对象表示法为基础的轻量级数据标记语言。能够处理的数据类型有 false/null/true/ 对象 / 数组 / 数字 / 字符串.

JSON 让数据更轻更纯粹，并且 JSON 的字符串形式可被 JavaScript 轻易地读入。因此常被用于开发web api.