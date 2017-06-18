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