## Go Web 开发学

参考:[Go Web 开发学习 (一)](http://blog.liwenqiu.me/2013/07/30/golang-default-server/),[Golang Http Server源码阅读](http://www.cnblogs.com/yjf512/archive/2012/08/22/2650873.html)

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

ListenAndServe 这个函数的作用是监听8080端口，并且在收到请求之后用第二个参数指定的处理器来处理.这里我们的第二个参数是nil，这种情况下会使用一个Go提供的默认的处理器来处理请求.默认的处理器怎么知道如何处理请求呢？答案就在 HandleFunc 这个函数，它向默认的处理器注册了指定的url由哪个函数来处理.先来看看Go默认提供的处理器是如何工作的.

打开 src/pkg/net/http/server.go, 我们从HandleFunc函数开始.

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

DefaultServeMux 就是Go提供的默认处理器(Mux = Multiplexer,多路复用器).它是一个全局变量，由 NewServeMux 函数创建

```go
var DefaultServeMux = NewServeMux()
func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool // whether any patterns contain hostnames//任何patterns是否包含hostnames,即是否以"/"开头
}
type muxEntry struct {
	explicit bool // 是否已注册
	h        Handler
	pattern  string //url
}
```

我们可以看到它初始化了一个 ServeMux 的结构体，ServeMux内部维护了一个map，map value 是 muxEntry 结构体，每个 muxEntry 就保存一个 url 和 Handler 之间的关联.ServeMux 还有一个读写锁，用来在并行环境下维护 muxEntry 的数据完整.

我们再来看看 Handler 是什么东西。

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

只要有函数 ServeHTTP(ResponseWriter, *Request) 的类型就可以当做一个 Hanlder。但是我们之前的那个示例程序里面的 hello 函数并不具备成为一个 Hanlder 的条件？

我们继续看 DefaultServeMux.HandleFunc 这个函数，它的第二个参数是我们示例程序里面的 hello 函数.

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

HandlerFunc(handler) 这里把我们的 hello 函数包装成了 HandlerFunc 类型，而 HandlerFunc 有 ServeHTTP 方法，并且 ServeHTTP 函数实际上就是调用我们的 hello 函数.这样我们的 hello 函数就被包装成了 Handler.

```go
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
	if pattern[0] != '/' {
		mux.hosts = true
	}
	// Helpful behavior:
	// If pattern is /tree/, insert an implicit permanent redirect for /tree.
	// It can be overridden by an explicit registration.
    // 只注册/tree/时,会隐性注册/tree,并将/tree跳转到/tree/;但当明确注册/tree时,其会覆盖隐性url(靠explicit字段实现).
	n := len(pattern)
	if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {
		// If pattern contains a host name, strip it and use remaining
		// path for redirect.
		path := pattern
		if pattern[0] != '/' {
			// In pattern, at least the last character is a '/', so
			// strings.Index can't be -1.
			path = pattern[strings.Index(pattern, "/"):]
		}
		mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(path, StatusMovedPermanently), pattern: pattern}
	}
}
```

接下来的 Handle 函数就是创建 muxEntry 保存到 map 里面，之后收到 url 请求就会根据这个映射关系来调用我们的处理函数。

到这里，我们示例程序里面 HandleFunc 的流程就走完了，通过多次调用 HandleFunc，我们就可以创建多个 url 和 处理函数的关联，示例程序里面我们只创建了 “/“ 和 hello 函数的关联.

我们再将注意力转到 ListenAndServe 函数，它的第二个参数是一个 Handler，当传nil给它的时候，默认会使用 DefaultServeMux 来当作 Handler.那么 ServeMux 也必须是一个 Handler，它应当具备 ServeHTTP 函数了？答案是当然.

```go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req) //开始交给Handler处理
}
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		w.Header().Set("Connection", "close")
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```

首先调用 ServeMux 的 Hanlder 函数来匹配 url 和 Handler

```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = mux.handler(r.Host, p)
			return RedirectHandler(p, StatusMovedPermanently), pattern
		}
	}
	return mux.handler(r.Host, r.URL.Path)
}
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}
// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
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
```

匹配到适合 url 的 Handler 之后就会调用 Handler 的 ServeHTTP 函数，如果没有匹配到合适的 url，就会返回一个 NotFoundHandler， 它会返回 404 错误给浏览器。

我们之前说过 hello 函数会被包装成一个 HandleFunc，并且在 ServeHTTP 函数里面调用我们的 hello 函数.这样当Go服务器收到请求之后，根据 url 找到合适的 Handler，调用 Handler 的 ServeHTTP，这样就调用到了我们的 hello 函数， hello 函数往 ResponseWriter 里面写入我们想要返回给浏览器的内容.

这样，示例程序里面一个请求，响应的流程就走完了.Go里面每个请求是在一个 goroutine 里面执行的，关于这部分，我们以后再研究,可参考:

- [Go源码分析——http.ListenAndServe()是如何工作的](http://blog.csdn.net/gophers/article/details/37815009)
- [](Go如何使得Web工作)(https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.3.md)
