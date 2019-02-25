# FAQ

## 组件的data是函数
如果传给组件的 data 是一个原始对象的话，则在建立多个组件实例时它们就会共用这个 data 对象，修改其中一个组件实例的数据就会影响到
其他组件实例的数据，这显然不是我们所期望的.

## 访问404
vue单页应用使用vue-router会有两种配置，即history模式和hash模式，但是hash模式其实会有很多限制，最主要的一点，url地址太丑了，所以我们在生产环境中也希望用到history模式.

但在HTML5 History 模式下, 使用nginx等代理服务器时，会遇到404的问题, 这是因为vue项目编译出来的dist中，并没有真正的本地资源提供给nginx，正确的做法是需要转交给vue-router来做前端路由.

```nginx
...
location / {
    root /root/dist;
    try_files $uri $uri/ /index.html =404; # 原理 404时直接返回index.html,再由ue-router来做前端路由
}
...
```

## vscode格式化vue文件时会自动把单引号转为双引号
vscode使用[vetur](https://github.com/vuejs/vetur)格式化代码, 而其默认使用`"vetur.format.defaultFormatter.js": "prettier"`格式化代码导致[这个问题](https://github.com/vuejs/vetur/issues/807).

解决方法,使用其他格式化插件或安装prettier并配置其参数. 我选择vscode自带的js格式化插件,只要修改vscode的用户设置即可:
```json
"vetur.format.defaultFormatter.js": "vscode-typescript",
"javascript.format.insertSpaceBeforeFunctionParenthesis": true # 使用vscode-typescript时js function的小括号前的空格被省略, 该行配置使其保留.
```

> [vetur格式化](https://vuejs.github.io/vetur/formatting.html#formatters)

## [Vue2 dist 目录下各个文件的区别](https://github.com/vuejs/vue/tree/dev/dist)
按照构建方式分, 可以分成 完整构建(包含独立构建和运行时构建) 和 运行时构建
按照规范分, 可以分成 UMD, CommonJS 和 ES Module.

完整构建 和 运行时构建的区别就是: 完整构建可以使用template选项; 而运行时构建只能用render选项.
*.dev.js 和 *.prod.js区别: 未压缩版的用于开发模, 包含了完整的警告和调试模式 因此其在开发时错误的参考信息和定位调试是非常便利的; 压缩版的用于生产模式.

- vue.esm．＊ ： 基于 ES Module 的构建, 可以用于 Webpack-2 和 rollup 之类打包工具, **推荐**
- vue.* : 基于 UMD 的构建, 可以用于 CDN 引用
- vue.common.* :　基于 CommonJS 的构建，可以用于 Webpack-1 和 Browserify 之类打包工具, **不推荐**

## 生命周期
- beforeMount和mounted不会在ssr中调用, 因为这两个方法需要操作DOM, 而ssr环境没有DOM; beforeCreated和created则可以在ssr中调用.
- beforeCreated及以前不能操作vm.$data(比如调用ajax/fetch获取数据), 因为此时reactivity还未就绪, 最早操作数据需要到created.
- render方法是在beforeMount和mounted之间调用(`*.vue文件`中的`template`会由`vue-loader`编译成render方法).

## watch, computed 和 methods 区别
watch仅**监控到指定值变化**才会执行, 且组件初始化时不执行(除非设置`immediate`和`handler`); `deep`属性允许监控object内的变化.
computed用于组合数据且**带缓存**, 仅依赖值变化才会被调用.
methods不带缓存, vm.$data中任何值变化都会被重新调用.

> **computed里不建议设置set方法.**
> 不要在computed和watch里修改依赖值, 否则会导致循环

## 方法
- `errorCaptured()`用于收集错误, 包括子组件.

## 指令
**dom的属性默认都是字符串, 除非使用v-bind表示该属性才能反映出其原始类型(比如整数)**, 比如[传递静态或动态 Prop](https://cn.vuejs.org/v2/guide/components-props.html).

- `v-cloak`: 仅页面直接使用vue.js库时需要, webpack打包时就不需要了.

## 组件
### props
- 属性`type`: [Boolean,String,Array,Object,Function...](https://cn.vuejs.org/v2/guide/components-props.html), 或者他们的组合, 比如`type: [String, Number]`
- 属性`required`与`default`是互斥的, 二选一即可
- **default返回对象时应与组件data一样返回function**.
- 属性`validator`支持更复杂的校验.

### name
父组件的`name`属性可在子组件用`this.$parent.$options.name`访问到(即**子组件可通过`this.$parent`访问/修改父组件, 但不推荐, 因为会导致逻辑混乱**).

### provide 和 inject
父/祖组件提供`provide`并再在子组件里使用`inject`, 即可在子组件里访问到父/祖组件, 但父/祖组件传递的provide内容不是响应式的.

## vue-router
### props属性
在route里启用`props`则会将route params(也可自行定义参数, 具体见api)作为props传入对应的route component, 这样就便于解耦,增强复用性: 无需在组件里调用`this.$route`而耦合vue-router(**推荐**).

## vuex
原因: 多个组件共享状态

异步请求数据用`action`.
modules声明`namespaced: true`, 用于区分同名变量或函数.