## DNS
- [扫盲 DNS 原理,兼谈"域名劫持"和"域名欺骗/域名污染"](2015_05_23_001.md)
- [DNS劫持和DNS污染的区别](2015_05_23_007.md)
- [全局精确流量调度新思路-HttpDNS服务详解](http://www.zmke.com/i/8705.html)
- [从理论到实践，全方位认识DNS（理论篇）](http://selfboot.cn/2015/11/05/dns_theory/)
- [从理论到实践，全方位认识DNS（实践篇）](http://selfboot.cn/2015/11/05/dns_theory/)
- [DNS 一站到家之记录类型](https://deepzz.com/post/dns-recording-type.html)
- [DNS 一站到家之常用工具](https://deepzz.com/post/dns-comman-tools.html)
- [*从理论到实践，全方位认识DNS](https://www.cloudxns.net/Index/index.html)

## 硬件
- [为什么寄存器比内存快？](2015_05_23_012.md)

## base
- [互联网协议入门（一）](http://www.ruanyifeng.com/blog/2012/05/internet_protocol_suite_part_i.html)
- [MVC，MVP 和 MVVM 的图示](http://www.ruanyifeng.com/blog/2015/02/mvcmvp_mvvm.html)
- [浏览器HTTP缓存原理分析](http://www.cnblogs.com/tzyy/p/4908165.html)
- [Web消息推送的技术及websocket](2016_09_11_001.md)
- [Http协议中的各种长度限制总结](https://sites.google.com/site/gzhpwiki/home/guo-cheng-shi-jian/http-xie-yi-zhong-de-ge-zhong-zhang-du-xian-zhi-zong-jie)

## web
- [已支持http2的软件](https://github.com/http2/http2-spec/wiki/Implementations)
- [前端工程精粹（一）：静态资源版本更新与缓存](http://www.infoq.com/cn/articles/front-end-engineering-and-performance-optimization-part1)
- [前端工程精粹（二）：静态资源管理与模板框架](http://www.infoq.com/cn/articles/front-end-engineering-and-performance-optimization-part2)
- [你必须了解的Session的本质](http://netsecurity.51cto.com/art/201402/428721.htm)
- [fouber/blog # 前端博客](github.com/fouber/blog)
- [75份开发者、设计师必备的速查表](http://info.9iphp.com/75-best-cheat-sheets-for-designers-and-programmers/)
- [中高级前端大厂面试秘籍，为你保驾护航金三银四，直通大厂](https://github.com/xd-tayde/blog/blob/master/interview-1.md)
- [HTTP2 in GO](https://www.tuicool.com/articles/jMreQbN)
- [Tailwind 重塑编写 CSS 的方式](https://www.tuicool.com/articles/BbYZFbj)
- [SFTP vs. FTPS](https://www.goanywhere.com/blog/2016/11/23/sftp-vs-ftps-the-key-differences)

    FTPS 和 SFTP 之间的一个主要区别是它们使用端口的方式:
    1. 对于所有 SFTP 通信，SFTP 只需要一个端口号，因此很容易保护
    1. FTPS 使用多个端口号. 命令通道的第一个端口用于身份验证和传递命令。但是，每次发出文件传输请求或目录列表请求时，都需要为数据通道打开另一个端口号. 因此它必须在防火墙中打开一系列端口以允许 FTPS 连接，这可能会带来更多的安全风险

    因此推荐使用sftp.

## css
 - [CSS动画简介](http://www.ruanyifeng.com/blog/2014/02/css_transition_and_animation.html)
 - [浏览器内建支持的响应式图像](http://zhuanlan.zhihu.com/FrontendMagazine/19858945)
 - [CSS参考手册](https://github.com/doyoe/css-handbook)
- [15个来自Codepen的令人赞叹的CSS动画例子](http://info.9iphp.com/15-inspiring-examples-of-css-animation-on-codepen/)
- [微博移动样式框架Marvel.css开发心得](http://uxfan.com/fe/css/2016/01/19/marvel.html)
- [实现canvas里图形的拖拽](https://github.com/smileyby/canvasDrag)
- [CSS网格（CSS Grid）布局入门](https://segmentfault.com/a/1190000010923572)
- [如何迁移 Sass 到 PostCSS](https://imweb.io/topic/5b422d444d378e703a4f4468)

## javascript
 - [React 入门实例教程](http://www.ruanyifeng.com/blog/2015/03/react.html)
 - [如何编写高扩展且可维护的 JavaScript系列一：”狂野西部“症](http://zhuanlan.zhihu.com/FrontendMagazine/19879147)
 - [如何编写高扩展且可维护的 JavaScript系列二：命名空间](http://zhuanlan.zhihu.com/FrontendMagazine/19880674)
 - [如何编写高扩展且可维护的 JavaScript系列三：模块](http://zhuanlan.zhihu.com/FrontendMagazine/19884662)
 - [如何编写高扩展且可维护的 JavaScript系列四：耦合](http://zhuanlan.zhihu.com/FrontendMagazine/19886391)
 - [新的 JavaScript 模块系统](http://zhuanlan.zhihu.com/FrontendMagazine/19850058)
 - [JavaScript闭包的底层运行机制](http://blog.leapoahead.com/2015/09/15/js-closure/)
 - [深入浅出ES6](http://www.infoq.com/cn/author/Jason-Orendorff)
 - [**近一万字的ES6语法知识点补充**](https://juejin.im/post/5c6234f16fb9a049a81fcca5)

## 静态资源
 - [静态资源（js以及css）发布对比](2015_05_26_001.md)
 - [图片服务架构演进](http://blog.aliyun.com/967)

## 认证
 - [使用 AngularJS & NodeJS 实现基于 token 的认证应用](http://zhuanlan.zhihu.com/FrontendMagazine/19920223)
 - [关于 Token，你应该知道的十件事](2015_05_29_001.md)
 - [身份验证中Cookies与 Tokens比较](2015_07_21_001.md)
 - [OAuth 2.0 笔记](https://blog.yorkxin.org)
 - [OAuth 2.0安全案例回顾](http://drops.wooyun.org/papers/598)
 - [JSON Web Token - 在Web应用间安全地传递信息](http://blog.leapoahead.com/2015/09/06/understanding-jwt/)
 - [八幅漫画理解使用JSON Web Token设计单点登录系统](http://blog.leapoahead.com/2015/09/07/user-authentication-with-jwt/)

## 经验
- [TJ Holowaychuk是怎样学习编程的？](http://zhuanlan.zhihu.com/FrontendMagazine/19572823)
- [将 Web 应用性能提高十倍的10条建议](https://linux.cn/article-7206-1.html)
- [最全前端资源汇集](github.com/nicejade/Front-end-tutorial)

## nginx
- [Nginx 性能优化整理](https://github.com/trimstray/nginx-quick-reference)
- [深入NGINX：我们如何设计它的性能和扩展性](https://linux.cn/article-5681-1.html)
- [Nginx下配置高性能，高安全性的https TLS服务](https://blog.helong.info/blog/2015/05/08/https-config-optimize-in-nginx/)
- [Nginx 服务器安装及配置文件详解](http://info.9iphp.com/nginx-install-and-config/)
- [Nginx 配置文件总结](http://info.9iphp.com/nginx-config-summary/)
- [NGINX缓存使用官方指南](https://linux.cn/article-5945-1-rel.html)
- [如何监控 NGINX（第一篇）](https://linux.cn/article-5970-1-rel.html)
- [如何收集 NGINX 指标（第二篇）](https://linux.cn/article-5985-1-rel.html)
- [Nginx+Keepalived实现站点高可用](https://linux.cn/article-5715-1-rel.html)
- [nginx 配置 location 总结及 rewrite 规则写法](https://linux.cn/article-5714-1-rel.html)
- [为最佳性能调优 Nginx](http://blog.jobbole.com/87531/)
- [agentzh 的 Nginx 教程](http://openresty.org/download/agentzh-nginx-tutorials-zhcn.html)
- [Nginx下流量拦截算法](http://homeway.me/2015/10/21/nginx-lua-traffic-limit-algorithm/)

- [NGINX源码阅读](http://www.cnblogs.com/ourroad/p/4863758.html)
- [Nginx源码分析之启动过程](http://www.rowkey.me/blog/2014/09/24/nginx-bootstrap/)
- [Nginx基本配置备忘](https://zhuanlan.zhihu.com/p/24524057)
- [Zstd](https://blog.gakaki.com/guide/zstd_chrome)

    [zstd-nginx-module](https://github.com/tokers/zstd-nginx-module)

    browser:
    - chrome: 118, 需要手动打开
    - fireforx: 126

## 架构
- [构建需求响应式亿级商品详情页](http://blog.jobbole.com/90359/)

## 交互
- [网站优化指南与用户体验五要素](http://www.shejidaren.com/web-user-experience.html?utm_source=tuicool)

## Rustful API
- [来自HeroKu的HTTP API 设计指南(中文版)](http://get.jobdeer.com/343.get)
- [理解RESTFul架构](http://blog.jobbole.com/88989/)
- [一个工作中用得到的RESTful API规范*](2016_05_27_001.md)
- [RESTful API 设计与工程实践](http://blog.m31271n.com/2017/03/02/RESTful-API-%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%B7%A5%E7%A8%8B%E5%AE%9E%E8%B7%B5/#资源排序)
- [是时候抛弃Postman了，试试 VS Code 自带神器插件](https://juejin.im/post/5e2067f7f265da3e405028fb)

## 负载均衡
- [负载均衡原理与技术实现](http://network.51cto.com/art/201509/492457_all.htm)
- [Cloudflare 放弃 Nginx，使用内部 Rust 编写的 Pingora](https://www.oschina.net/news/210473/cloudflare-pingora-nginx)
    
    - [Pingap 基于 Pingora 的反向代理软件](https://www.oschina.net/news/288724/pingap-0-2-0)

        `./pingap-linux-x86 --admin localhost:5555`

	- [pingsix](https://github.com/zhu327/pingsix)

## 优化
- [京东首页前端技术剖析与对比](http://www.barretlee.com/blog/2015/09/09/jd-architecture-analysis/)
- [利用现代浏览器所提供的强大 API 录制，回放并保存任意 web 界面中的用户操作](https://juejin.im/post/5c4fdf7c51882524fe52dc2b)

## 爬虫
- [爬虫技术浅析](http://drops.wooyun.org/tips/3915)
- [Python3 反爬虫原理与绕过实战](https://www.tuicool.com/articles/Zz2u2m3)

## 进阶
- [那些值得关注的 Web 开发者成长路线](http://www.oschina.net/news/75434/web-developer-growth-path)
- [Daily-Interview-Question : 壹题](https://github.com/Advanced-Frontend/Daily-Interview-Question)

## 安全
- [understanding-csrf](https://github.com/pillarjs/understanding-csrf/blob/master/README_zh.md)

## canvas
- [canvas图案拖动](http://www.hangge.com/blog/cache/detail_1057.html)

## web audio api
- [HTML5-Audio- The most useful filter nodes](http://masf-html5.blogspot.com/2016/04/html5-audio-most-useful-filter-nodes.html)
- [Web Audio API samples](https://webaudioapi.com/samples/)
- [standardized-audio-context](https://github.com/chrisguttandin/standardized-audio-context)

## [http3](https://github.com/quicwg)
- [quicwg](https://quicwg.org/)
- [httpwg](https://httpwg.org/specs/)
- [HTTP/3 explained](https://daniel.haxx.se/http3-explained/)
- [quic status](https://datatracker.ietf.org/wg/quic/about/)
- [HTTP/3 status](https://datatracker.ietf.org/doc/draft-ietf-quic-http/)
- [QUIC 协议在蚂蚁集团落地之综述](https://my.oschina.net/alimobile/blog/5296337)

    QUIC LB 组件：基于 NGINX 4层 UDP Stream 模块开发.
- [腾讯开源高性能轻量级跨平台 QUIC 协议库 - TQUIC](https://www.oschina.net/news/265538)
- [Async QUIC and HTTP/3 made easy: tokio-quiche is now open-source](https://blog.cloudflare.com/async-quic-and-http-3-made-easy-tokio-quiche-is-now-open-source/)

## 打包
- [Turbopack](https://www.oschina.net/news/214982/nextjs-13-released)

    Next.js 13 发布，推出快 700 倍的基于 Rust 的 Webpack 替代品

## ssr
- [Hydration](https://my.oschina.net/u/4090830/blog/10141047)

    "hydration"（水合）是指通过客户端 JavaScript 将静态 HTML 网页转化为动态网页的过程，以实现对 HTML 元素的事件处理。这个过程可以通过将事件处理程序附加到 HTML 元素上来完成

## https
- [HTTPS过程详解，tcpdump抓包一步一步分析](https://blog.csdn.net/simonchi/article/details/107563574)

## demo
- [用 Vue 全家桶纯手工搓了一个开源版「抖音」，高仿度接近 100%](https://www.oschina.net/news/286027)

## tools
- web ppt

    - [Slidev](https://sli.dev/)
    - [nodePPT](https://github.com/ksky521/nodePPT)