# [wails](https://v3.wails.io/zh-cn/quick-start/installation/)
ref:
- [wails v3 —— 用 Go 和 Vue 做桌面应用](https://www.fengfengzhidao.com/p/wails-v3)

```bash
# sudo dnf install gtk4-devel  webkitgtk6.0-devel
# go install github.com/wailsapp/wails/v3/cmd/wails3@latest
# wails3 init -n myapp -t vue
# cd myapp && wails3 dev
```

其他命令
```bash
# PACKAGE_MANAGER=pnpm wails3 dev # 使用pnpm安装依赖
```

## wails dev
`wails dev`=`wails3 task dev`:
1. 执行build/linux/Taskfile.yml的`build:native`

    1. 执行build/Taskfile.yml的`build:frontend`, 打包frontend项目
1. 执行build/linux/Taskfile.yml的`run`

## FAQ
### 用vp create的项目代替wails的frontend项目
1. 用vp create创建项目, 项目的目录名还是要为frontend
1. package.json添加

    ```json
    {
        "scripts": {
            "build:dev": "tsc -b && vp build --minify false --mode development"
        }
    }
    ```
1. 修改build/Taskfile.yml的
    - `generate:bindings`

        - `wails3 generate bindings -f '{{.BUILD_FLAGS}}' -clean=true` -> `wails3 generate bindings -f '{{.BUILD_FLAGS}}' -clean=true -ts`
    - `install:frontend:deps:npm`

        `npm install` -> `vp install`
    - `frontend:dev:npm`
        `npm run dev -- --port {{.VITE_PORT}} --strictPort` -> `npm run dev -- --host 0.0.0.0 --port {{.VITE_PORT}} --strictPort`
1. 修改vite.config.ts, 追加

    ```ts
    {
        plugins: [react()],
        server: {
            port: 9245,
        },
    }
    ```


问题:
1. wails3 dev启动后白屏, 终端上报`proxy error error="dial tcp4 127.0.0.1:9245: connect: connect refused"`, 但http://localhost:9245能访问成功

    应该是执行build/Taskfile.yml的`frontend:dev:npm`->`vp dev`时, `ss -anltp|grep 9245`返回`[::1]:9245`, 说明9245端口绑定到ipv6, 导致ipv4的127.0.0.1无法访问

    wails3 build是正常的

    解决: `frontend:dev:npm`追加`--host 0.0.0.0`
1. wails3 dev报`xxx implicitly has an 'any' type`

    生成的binding默认是js, 切换到ts: build/Taskfile.yml的`generate:bindings`, 追加参数`-ts`
1. [wails3 dev报`fatal error: non-Go code set up signal handler without SA_ONSTACK flag`](https://github.com/wailsapp/wails/blob/6329e9d2bec94a6916b9a700667b78a9569a3d9b/website/docs/guides/linux.mdx#panic-recovery--signal-handling-issues)



    