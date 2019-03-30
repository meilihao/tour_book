# vue cli

## FAQ
### vue cli3 在vue.config.js的configureWebpack配置DefinePlugin，不起作用
[因为默认配置里已经有一个 DefinePlugin 的实例了，请参考链接中提供的代码进行配置](https://github.com/vuejs/vue-cli/issues/3279)

### 如何配置url前缀
vue.config.js
```js
module.exports = {
  /**
   * You will need to set publicPath if you plan to deploy your site under a sub path,
   * for example GitHub Pages. If you plan to deploy your site to https://foo.github.io/bar/,
   * then publicPath should be set to "/bar/".
   * In most cases please use '/' !!!
   * Detail: https://cli.vuejs.org/config/#publicpath
   */
  publicPath: '/adminer/',
  ...
}
```