# websocket

WebSocket API是纯事件驱动的, 因此应用程序只要监听WebSocket对象上的事件，以便及时处理输入数据和连接状态的改变即可.

WebSocket编程遵循异步编程模式.

WebSocket对象调度4个不同的事件：
- open
- message : WebSocket API只输出完整的消息，而不是WebSocket帧.
- error : 错误还会导致WebSocket连接关闭
- close

## FAQ
### close时收到reason缺少头两个byte
```go
// wsutil.WriteServerMessage(conn, ws.OpClose, []byte("asdfasdf")) // 错误例子,ws client收到的CloseEvent.reason少两个字节
wsutil.WriteServerMessage(conn, ws.OpClose, []byte("\x01\x00asdfasdf"))
conn.Close()
```

构造断连请求的数据部分，前面两字节存放的是状态码,比如上例就是`256`,其体现在 ws client onclose事件CloseEvent的code属性上.

### close事件
close事件有3个有用的属性（property），可以用于错误处理和恢复：wasClean、code和error.

wasClean属性是一个布尔属性，表示连接是否顺利关闭.如果WebSocket的关闭是对来自服务器的一个close帧的响应，则该属性为true.如果连接是因为其他原因（例如，因为底层TCP连接关闭）关闭，则该属性为false.

code和reason属性表示服务器发送的关闭握手状态。这些属性和WebSocket.close()方法中的code和reason参数一致.