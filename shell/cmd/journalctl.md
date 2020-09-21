# journalctl

## 描述

journald是systemd独有的日志系统，替换了sysVinit中的syslog守护进程,使用journalctl用来读取日志.

1. 读取日志
       # journalctl
1. 查看所有引导日志
       # journalctl -b
1. 最近一次关机记录

    ```bash
    $ sudo journalctl -rb -1
    $ sudo vim /etc/systemd/system.conf
    #DefaultTimeoutStopSec=90s # 关机的默认等待时间通常设置为 90 秒. 启用该选项, 在这个时间之后，os会尝试强制停止服务
    ```

1. 实时显示系统日志（类似tail -f）
       # journalctl -f
1. 只显示指定的服务或可执行程序的日志
       # journalctl /usr/sbin/dnsmasq
1. 查看某个Unit的日志
       # journalctl -u nginx.service
       # journalctl -u nginx.service --since today
1. 合并显示多个 Unit 的日志
       # journalctl -u nginx.service -u php-fpm.service --since today
