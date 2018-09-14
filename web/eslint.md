# eslint

## vscode
- 安装ESLint扩展,eslint的配置文件名应为`.eslintrc.json`
- 配置
```json
    "eslint.validate": [
    "javascript",
    "javascriptreact",
    "html",
    "vue"
  ]
```

## vue
- 安装[eslint-config-vue](https://github.com/vuejs/eslint-config-vue)
- 配置
```json
{
    "extends":["standard", "vue"]
}
```
- 配置webpack2
```json
      {
        enforce: 'pre',
        test: /\.(js|vue)$/,
        exclude: /node_modules/,
        loader: 'eslint-loader'
      },
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      }
```

## FAQ
### iview 'vue/no-parsing-error parsing error invalid-first-character-of-tag-name'
有两种解决办法： 
- 1 MenuItem修改为：menu-item,**推荐**
- 2 在根目录下 .eslintrc.js 文件 rules 下添加：
```js
"vue/no-parsing-error": [2, { "x-invalid-end-tag": false }]
```

### eslint-loader不起作用
很可能是eslint配置里的`extends`规则没有涵盖到代码.