# net
## FAQ
### `address in use`
对于listener的监听Socket,Go默认设置了SO_REUSEADDR,这样当重启服务程序时, 不会因为`address in use`的错误而重启失败.