# pm2
pm2 是 node 进程管理工具，可以利用它来简化很多 node应用管理的繁琐任务，如性能监控、自动重启、负载均衡等.

## cmds
```bash
# npm install -g pm2 # 全局安装pm2，依赖node和npm
# pm2 show <appname|id> # 查看详细状态信息
# pm2 list # 查看所有启动的进程列表
# pm2 start app.json # 启动进程
# pm2 stop <id|all> # 停止 指定/所有 进程
# pm2 restart <id|all> # 重启 指定/所有 进程
# pm2 delete <id|all> # 杀死 指定/所有 进程
```

## 配置
ref:
- [配置文件](https://pm2.fenxianglu.cn/docs/general/configuration-file)

PM2 文件默认放在`$HOME/.pm2`

## FAQ
### mode=fork
在 fork 模式下，PM2 会使用操作系统的 child_process.fork 或 child_process.spawn API 来启动应用程序作为一个独立的进程

### `pm2 start <id>`报`[PM2] Applying action restartProcessId on app [52](ids: [ '52' ])`
git pull重新构建go程序后, pm2 start报错.

pm2.json有变化(元信息变化, 内容没变), 需要`pm2 delete 52 && pm2 save` + `pm2 start 52 && pm2 save`