# vite+
ref:
- [guide](https://viteplus.dev/guide/)
- [给到夯！前端工具链新标杆 Vite Plus 初探](https://juejin.cn/post/7618920576733184026)
- [一文搞懂 Vite+ 中的vp env：优雅管理 Node.js 版本的新方式](https://juejin.cn/post/7619214301320380479)

注意:
1. v0.1.18: 还不支持创建vue项目

```bash
$ curl -fsSL https://vite.plus | bash
$ vp env list-remote # 获取所有可用的nodejs版本
$ vp env list # 已安装的nodejs版本
$ vp env uninstall v24.14.0 # 下载已安装的nodejs版本
$ vp env current # 当前使用的nodejs版本
$ vp env on # 启用托管模式
$ vp install # 选中pnpm
$ vp create vite:monorepo # 创建monorepo项目
$ vp create vite:application #添加向monorepo项目添加新项目
```

## FAQ
### vp create的`Vite+ Monorepo和Vite+ Application`区别
Vite+ Application是单个项目
Vite+ Monorepo是一组项目的组织方式

### vp dev和vp preview区别
dev = 开发模式, 支持热更新
preview = 构建后预览. 需先vp build, 再vp preview

### `Vite+ Monorepo`项目下执行`vp dev`无法访问, 但`vp run website#dev`或`cd <project>/apps/website && vp dev`能正常访问 
```bash
$ cat package.json
{
  ...,
  "scripts": {
    "ready": "vp check && vp run -r test && vp run -r build",
    "dev": "vp run website#dev",
    "prepare": "vp config"
  },
  ...
}
```
启动dev server时, 
`vp dev`使用`/home/xxx/.vite-plus/js_runtime/node/24.15.0/bin/node /data/tmpfs/voidzero/yyy/node_modules/.pnpm/vite-plus@0.1.18_@types+node@25.6.0_jiti@2.6.1_typescript@6.0.2_vite@8.0.8_@types+node@_6c3352028afd0d77232260e20359556f/node_modules/vite-plus/dist/bin.js dev`, 而`cd <project>/apps/website && vp dev`使用`/home/xxx/.vite-plus/js_runtime/node/24.15.0/bin/node /data/tmpfs/voidzero/yyy/node_modules/.pnpm/vite-plus@0.1.18_@types+node@25.6.0_@voidzero-dev+vite-plus-core@0.1.18_@types+node@25._7394d3fc3cf68e557d6eebcdb659f21e/node_modules/vite-plus/dist/bin.js dev`

`ss -anltp|grep 9245`看看9245端口绑定到ipv4还是ipv6, 从而确定用`http://localhost:9245`还是`http://127.0.0.1:9245`访问