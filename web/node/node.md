# node

node version: 8.5.0

## CommonJS,AMD 与 ES6 Modules实现区别
ES6 模块的设计思想，是尽量的静态化，使得编译时就能确定模块的依赖关系，以及输入和输出的变量. CommonJS和AMD模块，都只能在运行时确定这些东西.

**一个 CommonJS 的模块在没有被执行完之前，它的结构（API）是不可知的**,即使在它被执行完以后，它的结构也可以随时被其他代码修改.

而所有的 import 和 export 语句都会在代码执行之前被解析出来.

> Babel的做法实际上是将不被支持的import翻译成目前已被支持的require.

**推荐使用ESM**

## import
解析import时会执行该文件中的函数,比如vuejs的`src/core/instance/index.js`里的`initMixin`函数.

## Error
### `SyntaxError: Unexpected token import`或`SyntaxError: Unexpected token export`
[node8.5.0开始支持ESM](https://nodejs.org/api/esm.html),但需要使用`--experimental-modules`flag来开启,且文件的扩展名必须是`mjs`.

### self signed certificate in certificate chain 
```
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0' // 禁用检查tls证书
```