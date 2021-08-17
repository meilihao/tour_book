# sysvinit
参考:
- [SysVinit (简体中文)](https://wiki.archlinux.org/title/SysVinit_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))

sysvinit 就是system V 风格的init 系统, **推荐使用systemd, 这里仅是工作需要而记录.**


## FAQ
### `service supervisor restart`没效果
原因: 执行`start-stop-daemon --stop`后的7秒内, supervisord还在运行(因为部分被管理的应用还没停止). 又因为`DODTIME=5`, 即stop操作的5s后执行`start-stop-daemon --start`时就认为supervisord已运行而不进行操作.

参考[`man start-stop-daemon`](https://man7.org/linux/man-pages/man8/start-stop-daemon.8.html), 使用`start-stop-daemon --stop ... --retry=TERM/30/KILL/5`启用重试机制检查supervisord是否仍在执行.