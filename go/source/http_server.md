## Go Web 开发学

参考:
- [Golang Http Server源码阅读](http://www.cnblogs.com/yjf512/archive/2012/08/22/2650873.html)
- [Go源码分析——http.ListenAndServe()是如何工作的](http://blog.csdn.net/gophers/article/details/37815009)
- [Go如何使得Web工作](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.3.md)
- [Go的http包详解](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.4.md)

用Go写一个最简单的 Http Server
```go
package main
import (
    "fmt"
    "net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello Go.")
}
func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe(":8080", nil)
}
```

编译运行这个程序之后，当我们用浏览器请求 localhost:8080 就可以看到 “Hello Go.” 了。

ListenAndServe 这个函数的作用是监听8080端口，并且在收到请求之后用第二个参数指定的处理器来处理.这里我们的第二个参数是nil，这种情况下会使用一个Go提供的默认的处理器(http.DefaultServeMux)来处理请求.那go到底如何处理请求呢?

我们将注意力转到 [ListenAndServe](http://localhost:6060/src/net/http/server.go) 函数，看看它是如何工作的:
```go
func ListenAndServe(addr string, handler Handler) error {
    // 创建一个Server结构体，调用该结构体的ListenAndServer方法
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}

type Server struct {
    Addr           string        // TCP address to listen on, ":http" if empty
    Handler        Handler       // handler to invoke, http.DefaultServeMux if nil
    ...
}
```
从这个函数的源码中就可以看出，调用http.ListenAndServe之后真正起作用的是Server结构体的LisntenAndServe方法，给http.ListenAndServe传递的参数只是用来创建一个Server结构体实例.

接下来继续看看Server.ListenAndServe的代码:
```go
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http" // 等价于":80"
	}
	ln, err := net.Listen("tcp", addr)  // 创建了一个TCP Listener,用于接收客户端的连接请求
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```
可以看到，Server.ListenAndServe中创建了一个服务器Listener，然后在返回时把它传给了Server.Serve()方法并调用了该方法.

继续分析Server.Serve:
```go
func (srv *Server) Serve(l net.Listener) error {
	...
  // 这里是服务器的主循环，通过传进来的listener接收来自客户端的请求并建立连接，  
  // 然后为每一个连接创建goroutine执行c.serve()，这个c.serve就是具体的服务处理函数
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew) // before Serve can return
		go c.serve() // 为每一个建立的连接创建goroutine来进行服务
	}
}
```

再看conn.serve具体的处理代码:
```go
func (c *conn) serve() {
	...
	for {
		w, err := c.readRequest()  // 读取客户端的请求
		...
		serverHandler{c.server}.ServeHTTP(w, w.req)  //对请求进行处理
		...
	}
}

type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}
```

到这里就可以发现,http.ListenAndServe的第二个参数为nil时，那么它会自动以http.DefaulServeMux作为默认的处理器来处理请求.

那http.DefaulServeMux又是如何处理请求呢？
```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r) // 返回对应设置路由的处理Handler
	h.ServeHTTP(w, r)
}

func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path { // 返回规范的url路径,清除.和..
			_, pattern = mux.handler(r.Host, p) // 查找具体url路径对应的Handler
			url := *r.URL
			url.Path = p
			return RedirectHandler(url.String(), StatusMovedPermanently), pattern // 跳转请求
		}
	}
	return mux.handler(r.Host, r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path) // 查找Handler
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), "" // 没有找到,使用NotFoundHandler,返回 404 错误给浏览器
	}
	return
}

// Does path match pattern?
// 前缀匹配, 因此pathMatch("/","/a")和pathMatch("/a","/a")都会返回true
func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
// match 会遍历路由信息字典，找到所有匹配该路径最长的那个.因为使用前缀匹配,因此如果找不到但存在"/"时会返回"/"的Handler.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[path]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
根据上面的代码,http.DefaultServeMux先根据注册在其中url路由路径查找到具体的Handler,再调用Handler的ServeHTTP方法进行具体处理.
那http.DefaultServeMux里的Handler是哪来的呢?
答案就在 HandleFunc 这个函数，它向默认的处理器注册了指定的url由哪个函数来处理.先来看看Go默认提供的处理器是如何工作的:
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }

// HandlerFunc type是一个适配器，通过类型转换让我们可以将普通的函数作为Handler使用
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

//  HandlerFunc(handler) 这里是把 hello 函数包装成了 HandlerFunc 类型，而 HandlerFunc 有 ServeHTTP 方法，并且 ServeHTTP 函数实际上就是调用我们的 hello 函数.
// 这样我们的 hello 函数就被包装成了 Handler.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}

// 接下来的 Handle 函数就是创建 muxEntry 保存到 map 里面，之后收到 url 请求就会根据这个映射关系来调用我们的处理函数。
// 到这里，示例程序里面 HandleFunc 的流程就走完了，通过多次调用 HandleFunc，我们就可以创建多个 url 和 处理函数的关联，示例程序里面我们只创建了 “/“ 和 hello 函数的关联.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if mux.m[pattern].explicit {
		panic("http: multiple registrations for " + pattern)
	}
	mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}
	if pattern[0] != '/' { //注册的url是否带hosts
		mux.hosts = true
	}
	// Helpful behavior:
	// If pattern is /tree/, insert an implicit permanent redirect for /tree.
	// It can be overridden by an explicit registration.
  // 只注册/tree/时(即不存在/tree),会隐性注册/tree,并将/tree跳转到/tree/;但当明确注册/tree时,其会覆盖隐性url(靠explicit字段实现).
	n := len(pattern)
	if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {
		// If pattern contains a host name, strip it and use remaining
		// path for redirect.
		path := pattern
		if pattern[0] != '/' {
			// In pattern, at least the last character is a '/', so
			// strings.Index can't be -1.
			path = pattern[strings.Index(pattern, "/"):] //舍弃host name
		}
		url := &url.URL{Path: path}
		mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(url.String(), StatusMovedPermanently), pattern: pattern}
	}
}

type ServeMux struct {
	mu    sync.RWMutex         // 读写锁，由于请求涉及到并发处理，因此这里需要一个锁机制
	m     map[string]muxEntry  // 路由规则，一个string对应一个mux实体，这里的string就是url()注册的路由表达式)
	hosts bool                 // 是否在任意的url中带有host信息
}
type muxEntry struct {
	explicit bool    // 是否已注册该url
	h        Handler // 路由表达式对应哪个handler
	pattern  string  // 注册的url路径(即匹配的url)
}
// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }
// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()
```
我们发现http.DefaultServeMux是ServeMux的实例(Mux = Multiplexer,多路复用器),是一个全局变量，由 NewServeMux 函数创建.
ServeMux内部维护了一个map，map value 是 muxEntry 结构体，每个 muxEntry 就保存一个 url 和 Handler 之间的关联.ServeMux 还有一个读写锁，用来在并行环境下维护 muxEntry 的数据完整.

这样一个go的web程序就跑起来了.

### 总结
通过对http包的分析之后，现在梳理一下整个的代码执行过程:

首先调用Http.HandleFunc:
1. 调用了DefaultServeMux的HandleFunc
1. 调用了DefaultServeMux的Handle
1. 往DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则

其次调用http.ListenAndServe(":8080", nil)
1. 实例化Server
1. 调用Server的ListenAndServe()
1. 调用net.Listen("tcp", addr)监听端口
1. 启动一个for循环，在循环体中Accept请求
1. 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()
1. 读取每个请求的内容w, err := c.readRequest()
1. 判断handler是否为空，如果没有设置handler（这个例子就没有设置handler），handler就设置为DefaultServeMux
1. 调用handler的ServeHTTP()

按照上面的例子，下面就进入到DefaultServeMux.ServeHTTP()

1. 根据request选择handler
  A 判断是否有路由能满足这个request（循环遍历ServerMux的muxEntry）
  B 如果有路由满足，返回这个路由的handler
  C 如果没有路由满足，返回NotFoundHandler
1. 执行handler.ServeHTTP(),输出结果

## 进阶
如下代码所示，我们自己实现了一个简易的路由器
```go
package main

import (
    "fmt"
    "net/http"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        sayhelloName(w, r)
        return
    }
    http.NotFound(w, r)
    return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello myroute!")
}

func main() {
    mux := &MyMux{}
    http.ListenAndServe(":9090", mux)
}
```
