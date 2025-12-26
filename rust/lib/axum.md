# axum
ref:
- [axum中文简述](https://luxiangdong.com/2023/03/01/axum1/)
- [Axum WEB开发](https://zhuanlan.zhihu.com/p/526613827)

examples:
- [axum.rs 专题之《漫游axum》](https://github.com/axumrs/roaming-axum)
- [axum/examples](https://github.com/tokio-rs/axum/tree/main/examples)

axum是一个专注于极简高效和模块化的web应用程序框架.

axum区别于其他框架的地方是没有自己的中间件系统，而是使用tower::Service, 这允许与使用hyper或tonic编写的应用程序共享中间件.

## handler
在axum中，handler是一个异步函数，它接受零个或多个extractors作为参数，并返回可以转换为response的内容.

## extractor
提取器是一种实现FromRequest或FromRequestParts的类型。提取器是您如何分离传入请求以获得处理程序所需的部分

## Response
任何实现IntoResponse的东西都可以从处理器返回

## 中间件
ref:
- [Axum笔记：中间件（middleware）](https://zhuanlan.zhihu.com/p/647873624)
- [Axum的middleware模块](https://yuxuetr.com/wiki/axum/axum-middleware)

使用:
1. 使用 Router::layer 和 Router::route_layer 对整个路由器进行中间件处理

	生效:
	1. Router::layer：对所有request生效，不管是否该请求是否有命中的route
	1. Router::route_layer：当且仅当request有命中的route时生效，其他不生效

	执行顺序:
	1. 当使用 Router::layer（或类似方法）添加中间件时，所有先前添加的路由都将被包裹在已有中间件中, 这导致中间件从底部到顶部执行, 即洋葱头模型
	1. ServiceBuilder 的工作原理是将所有层组合在一起，使它们从上到下依次运行
1. 通过 MethodRouter::layer 和 MethodRouter::route_layer 添加到方法路由器

	生效:
	1. MethodRouter::layer：除了只能绑定在一个特定的method_router，其他与 Router::layer 效果一致
	1. MethodRouter::route_layer：除了只能绑定在一个特定的method_router，其他与 Router::router_layer 效果一致，特别的，该方式在method_router的fallback中不触发
1. 对特定的处理程序使用 Handler::layer

	handler粒度的Layer，也是最小粒度的Layer

## 与处理器共享状态（Sharing state with handlers）
最常见的是:
1. 使用State提取器, **推荐且最惯用的**

	使用:
	1. 结构体定义：将您需要共享的所有状态封装到一个结构体中
	1. 初始化：在 Router 上调用 .with_state() 附加状态
	1. 提取：在 Handler 函数签名中使用 State<T> 提取状态

1. 使用请求扩展

	使用:
	1. 初始化：通过 layer 或中间件将数据注入到请求管道中
	1. 提取：在 Handler 函数签名中使用 Extension<T> 提取该数据

场景	推荐方法	理由
数据库连接池、全局配置、缓存	State<AppState>	最规范，将所有核心状态集中管理，并且 Axum 官方推荐。
中间件注入的认证信息	Extension<User>	数据是在请求处理过程中（例如在 Auth 中间件之后）生成的，需要传递给后续 Handler。