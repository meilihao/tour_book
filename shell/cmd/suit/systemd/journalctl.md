# journalctl
journald是systemd独有的日志系统，替换了sysVinit中的syslog守护进程,使用journalctl用来读取日志.

> 一些发行版已经配置了日志，以便将日志写入磁盘（/var/log/journal），而其他发行版将日志保存在内存中（/run/log/journal）.

# 选项
- -k : 内核日志
- -b : 启动日志
- -u : 指定服务
- -n : 指定条数
- -p : 指定类型

       - emerg 系统出现严重故障，比如内核崩溃
       - alert 应立即修复的故障，比如数据库损坏
       - crit 危险性较高的故障，比如硬盘损坏导致程序运行失败
       - err 危险性一般的故障，比如某个服务启动或运行失败
       - warning 警告信息，比如某个服务参数或功能出错
       - notice 不严重的一般故障，只是需要抽空处理的情况
       - info 通用性消息，用于提示一些有用的信息
       - debug 调试程序所产生的信息
       - none 没有优先级，不进行日志记录
- -f : 实时刷新（追踪日志）
- --since : 指定时间

       - "-1 hour" : 最近1小时

- --disk-usage : 占用空间
- -x/--catalog : 将错误消息与日志中的解释性文本放在一起输出
- -e/--pager-end: 从journalctl的输出末尾开始查看日志
- --boot <N>: 指定要查看哪个启动会话的日志

## example

1. 读取日志
       # journalctl
1. 查看所有引导日志
       # journalctl -b
1. 最近一次关机记录

    ```bash
    $ sudo journalctl -rb -1 # `-1` is from `journalctl --list-boots`
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

1. 按时间查询

       # journalctl --since "2018-03-26 20:20:00" [--until "2021-08-26 03:00"]
       # journalctl --since "12:00" --until "14:00"
       # journalctl --since "2020-07-01" --until "2020-08-01"

## 日志管理
Systemd 统一管理所有 Unit 的启动日志。日志的配置文件是`/etc/systemd/journald.conf`.

systemd-journald 服务收集到的日志默认保存在 /run/log 目录中，重启系统会丢掉以前的日志信息, 修改配置文件 /etc/systemd/journald.conf，把 Storage=auto 改为 Storage=persistent，并取消注释，然后`systemctl restart systemd-journald.service`即可实现持久化日志(`/var/log/journal`).

> [systemd loglevel](https://wiki.archlinux.org/title/Systemd/Journal)

```
# 查看所有日志（默认情况下 ，只保存本次启动的日志）
$ sudo journalctl

# 查看内核日志（不显示应用日志）
$ sudo journalctl -k

# 查看日志占据的磁盘空间
$ sudo journalctl --disk-usage

# 查看系统本次启动开始的日志
$ sudo journalctl -b
$ sudo journalctl -b -0

# 查看上一次启动的日志（需更改设置）
$ sudo journalctl -b -1

# 查看指定时间的日志
$ sudo journalctl --since "2019-12-23 09:00:00" > log.log
$ sudo journalctl --since "20 min ago"
$ sudo journalctl --since yesterday
$ sudo journalctl --since "2015-01-10" --until "2015-01-11 03:00"
$ sudo journalctl --since 09:00 --until "1 hour ago"

# 显示尾部的最新10行日志
$ sudo journalctl -n

# 显示尾部指定行数的日志
$ sudo journalctl -n 20

# 实时滚动显示最新日志
$ sudo journalctl -f

# 查看指定服务的日志
$ sudo journalctl /usr/lib/systemd/systemd

# 查看指定进程的日志
$ sudo journalctl _PID=1

# 查看某个路径的脚本的日志
$ sudo journalctl /usr/bin/bash

# 查看指定用户的日志
$ sudo journalctl _UID=33 --since today

# 查看某个 Unit 的日志
$ sudo journalctl [--boot] -u nginx.service # `--boot[=ID]`当前或某次boot时的log
$ sudo journalctl -u nginx.service --since today

# 实时滚动显示某个 Unit 的最新日志
$ sudo journalctl -u nginx.service -f

# 合并显示多个 Unit 的日志
$ journalctl -u nginx.service -u php-fpm.service --since today

# 查看指定优先级（及其以上级别）的日志，共有8级
# 0: emerg
# 1: alert
# 2: crit
# 3: err
# 4: warning
# 5: notice
# 6: info
# 7: debug
$ sudo journalctl -p err -b

# 日志默认分页输出，--no-pager 改为正常的标准输出
$ sudo journalctl --no-pager

# 以 JSON 格式（单行）输出
$ sudo journalctl -b -u nginx.service -o json

# 以 JSON 格式（多行）输出，可读性更好
$ sudo journalctl -b -u nginx.serviceqq
 -o json-pretty

# 指定日志文件占据的最大空间
$ sudo journalctl --vacuum-size=1G // /etc/systemd/journald.conf#SystemMaxUse=100M + systemctl restart systemd-journald

# 指定日志文件保存多久
$ sudo journalctl --vacuum-time=1years

# 显示journalctl日志的字段名, 保存位置
journalctl -f -o verbose

# 设置`/etc/machine-id`
systemd-machine-id-setup # = dbus-uuidgen --ensure=/etc/machine-id

# 清空systemd log
journalctl --flush --rotate # 要求日志守护进程轮换日志文件. 日志文件轮换的效果是所有当前**活动的日志文件**都被标记为已归档并重命名，以便将来永远不会写入它们, 然后在它们的位置创建新的（空的）日志文件. **注意: 活动日志文件不会被`--vacuum*=`命令删除.
journalctl --vacuum-time=1s # 使所有日志文件不包含早于 1s 的数据
journalctl --vacuum-size=400M # 清除所有存档的日志文件，并保留最后 400MB 的文件. 只对存档的文件有效
journalctl --vacuum-files=2 # 只有最后两个日志文件被保留，其他的都被删除. 只对存档的文件有效

# 检查日志文件的完整性
journalctl --verify
```

## FAQ
### 内核崩溃 (Kernel panics) systemd log丢失
参考:
- [systemd-journald missing crash logs](https://unix.stackexchange.com/questions/414871/systemd-journald-missing-crash-logs)

[查看 panic 信息](https://wiki.archlinux.org/index.php/General_troubleshooting_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)#%E6%95%85%E9%9A%9C%E6%81%A2%E5%A4%8D%E6%8E%A7%E5%88%B6%E5%8F%B0)

### 配置自动删除
/etc/systemd/journald.conf:
<table border="1" cellpadding="4">
<thead>
<tr>
<th>journald.conf 参数</th>
<th>描述</th>
<th>实例</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>SystemMaxUse</code></td>
<td style="white-space: normal;">指定日志在持久性存储中可使用的最大磁盘空间</td>
<td><code>SystemMaxUse=500M</code></td>
</tr>
<tr>
<td><code>SystemKeepFree</code></td>
<td style="white-space: normal;">指定在将日志条目添加到持久性存储时，日志应留出的空间量。</td>
<td><code>SystemKeepFree=100M</code></td>
</tr>
<tr>
<td><code>SystemMaxFileSize</code></td>
<td style="white-space: normal;">控制单个日志文件在被轮换之前在持久性存储中可以增长到多大。</td>
<td style="white-space: normal;"><code>SystemMaxFileSize=100M</code></td>
</tr>
<tr>
<td><code>RuntimeMaxUse</code></td>
<td style="white-space: normal;">指定在易失性存储中可以使用的最大磁盘空间（在&nbsp;<code>/run</code>&nbsp;文件系统内）。</td>
<td><code>RuntimeMaxUse=100M</code></td>
</tr>
<tr>
<td><code>RuntimeKeepFree</code></td>
<td style="white-space: normal;">指定将数据写入易失性存储（在&nbsp;<code>/run</code>&nbsp;文件系统内）时为其他用途预留的空间数量。</td>
<td><code>RuntimeMaxUse=100M</code></td>
</tr>
<tr>
<td><code>RuntimeMaxFileSize</code></td>
<td style="white-space: normal;">指定单个日志文件在被轮换之前在易失性存储（在&nbsp;<code>/run</code>&nbsp;文件系统内）所能占用的空间量。</td>
<td style="white-space: normal;"><code>RuntimeMaxFileSize=200M</code></td>
</tr>
</tbody>
</table>