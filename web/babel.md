# babel

Babel是一个编译器,可以把用最新标准编写的 JavaScript 代码向下编译转为ES5代码，从而在现有环境执行.

Babel6 更具模块化， 其将所有的转码功能以插件的形式提供.

## presets

preset是babel6插件的集合.比如常用的[ES2015 preset](https://babeljs.io/docs/plugins/preset-es2015/).

plugins和presets编译顺序:

- plugins优先于presets进行编译。
- plugins按照数组的index增序(从数组第一个到最后一个)进行编译
- presets按照数组的index倒序(从数组最后一个到第一个)进行编译
