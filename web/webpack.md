# webpack

参考:

- [webpack中文](https://webpack.vuefe.cn/index/)
- [webpack](https://zy108830.gitbooks.io/webpack-doc/content/index.html)

如何调用webpack

选择一个目录下有webpack.config.js文件的文件夹，然后运行下面的命令:

- webpack : 开发环境下编译
- webpack -p : 产品编译及压缩
- webpack --watch : 开发环境下持续的监听文件变动来进行实时编译(非常快!)
- webpack -d : 提供source maps,方便调试代码

## 参数

- --display-error-details : 查看错误,比如依赖缺失的情况
- --progress : 显示打包进度
- --colors : 输出内容加颜色
- --config : 指定webpack配置文件的位置

## loader

- css-loader

遍历 CSS 文件，然后找到`url()`表达式然后处理它们

- style-loader

把原来的 CSS 代码插入页面中的一个 style 标签中

- file-loader

```
# npm install file-loader --save-dev
	参考格式：
	{
		test: /\.(eot|woff|svg|ttf|woff2|gif)(\?|$)/,
		loader: 'file-loader?name=[hash].[ext]'
	}
```

- url-loader

把需要转换的图片路径变成 BASE64 字符串, limit参数用于转化图片的最大大小限制.

```
# npm install url-loader --save-dev
  参考格式：
  {
  test: /\.(png|jpg)$/,
  loader: 'url?limit=1200&name=[hash].[ext]'
  }
```

## config

### entry

打包入口

### output

- path : 打包的目录
- `__dirname` : 当前目录
- filename : 打包后输出的文件名

### module

#### loader

Webpack 本身只能处理 JavaScript 模块，如果要处理其他类型的文件，就需要使用 loader 进行转换。
Loader 可以理解为是模块和资源的转换器，它本身是一个函数，接受源文件作为参数，返回转换的结果。

Loader 可以通过管道方式链式调用，每个 loader 可以把资源转换成任意格式并传递给下一个 loader ，但是最后一个 loader 必须返回 JavaScript.

- loader: 'style!css!sass' : 需要注意loader的顺序，意为先使用sass加载器处理，解析为普通的css文件，再处理css文件，最后处理样式，类似于pipe的概念.
- loader: 'url-loader?limit=8192' : 内联的base64的图片地址，图片要小于8k，直接的url的地址则不解析

### resolve

- extensions ['', '.js', '.json', '.coffee'] : require文件时省略文件的扩展名,现在你require文件的时候可以直接使用require('file')，不用使用require('file.coffee')

## 其他

使webpack支持es6需要用到babel，它可以帮助我们将es6语法转化为es5的格式，可以在这里测试。
首先，需要安装babel

```
$ npm install babel-loader babel-core
```

安装babel-preset

```
$ npm install babel-preset-es2015
```

然后，配置webpack.config.js，如下：

```
module: {
  loaders: [
    { test: /\.js$/, exclude: /node_modules/, loader: "babel-loader" }
  ]
}
```

最后还需要加上babel的配置文件，在项目的根目录下创建文件.babelrc

```
{ "presets": ["es2015"] }
```
这样，我们的webpack就支持通过es6的export和import实现模块化了.

[ES7配置](http://es6.ruanyifeng.com/#docs/intro)

## webpack-dev-server

webpack还为我们提供了webpack-dev-server,它是一个小型的基于node的express服务器。利用webpack-dev-server我们可以轻松地启动服务。且webpack-dev-server默认采用的是watch模式，也就是说它会自动监测文件的变化，并在页面做出实时更新，我们不需要每次都重新编译.其默认使用8080端口.

```
$ npm install webpack-dev-server -g
```

### 参数

- --devtool eval : 为代码创建源地址,当有任何报错的时候可以让你更加精确地定位到文件和行号

### 原理

webpack-dev-server通过sockjs实现实时变化监测，当文件变化时，服务器就会向客户端发送通知，客户端重新渲染。
每次文件变化都会触发webpack-dev-server对文件进行重新编译，但是这个编译后的文件并不会每次都保存到我们的dist目录下，而是存放在内存中,默认情况下，这个文件的路径为当前路径。
也就是说，每次文件变化后，内存中的bundle.js做出了实时的变化，而output中配置的文件其实并没有变。
也就是说，如果我们在index.html中使用./bundle.js(` <script type="text/javascript" src="./dist/bundle.js"></script>`)则能够实现页面内容的实时更新。
那么，如何配置可以使其支持我们当前的写法./dist/bunlde.js呢？就需要用到publicPath这个属性。
设置后就可以通过./dist/bundle.js路径访问到内存中的文件，在当前路径下已经存在同名文件的情况下，webpack-dev-server会优先使用内存中的文件.

### 自动刷新

webpack-dev-server提供了两种自动刷新的模式

1. iframe模式
1. inline模式

这两种模式都支持Hot Module Replacement（热加载），所谓热加载是指当文件发生变化后，内存中的bundle文件会收到通知，同时更新页面中变化的部分，而非重新加载整个页面。

- iframe模式

```
$ webpack-dev-server
```

启动后,页面顶部会显示一条黑色背景的横幅,实际的页面其实是在iframe内显示.
这种方式有一点需要注意：浏览器地址栏的url地址不会受页面跳转的影响，将一直保持为`http://localhost:8080/webpack-dev-server`.

- inline模式,**推荐**

```
$ webpack-dev-server --inline
```
或
```
// 在webpack.config.js中加入
devServer: { inline: true }
```

这样我们就可以通过`http://localhost:8080/<path>`来访问我们的文件了,比如这样http://localhost:8080/index.html来访问index.html，且页面跳转回反映在浏览器的地址栏中。

### Hot Module Replacement（热加载）

```
$ webpack-dev-server --hot
```

打开页面可以在Console控制台看到启动内容，说明热加载配置成功。其中HMS表示热加载模块，WDS表示webpack-dev-server.
