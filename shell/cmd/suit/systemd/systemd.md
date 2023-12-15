# systemd

为系统的启动和管理提供一套完整的解决方案.

参考:
- [走进Linux之systemd启动过程](https://linux.cn/article-5457-1.html)
- [linux systemd工具学习](https://github.com/yifengyou/learn-systemd)
- [Systemd 入门教程：命令篇](http://www.ruanyifeng.com/blog/2016/03/systemd-tutorial-commands.html)
- [使用 systemd 来管理启动项](https://linux.cn/article-13402-1.html)
- [Rethinking PID 1](http://0pointer.de/blog/projects/systemd.html)
- [systemd for Administrators Part I ~ XI](http://0pointer.de/blog/projects/systemd-for-admins-1.html)
- [SysVinit vs Systemd Cheatsheet](https://fedoraproject.org/wiki/SysVinit_to_Systemd_Cheatsheet)
- [Understanding and administering systemd](https://docs.fedoraproject.org/en-US/quick-docs/understanding-and-administering-systemd/index.html)
- [Resource Management Guide - 资源管理](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/resource_management_guide/index)
- [Linux 启动性能分析](https://linux.cn/article-14462-1.html)
- [systemd攻略 - (3)如何利用systemd控制cgroup,实战](https://muahao.github.io/systemd-03/)

### systemd 与 System V init 的区别以及作用
|System V init 运行级别|systemd 目标名称|作用|
|0| runlevel0.target, poweroff.target |关机|
|1| runlevel1.target, rescue.target |一个基本的系统，包括挂载文件系统，但是只运行最基础的服务，以及一个主控制台上的用于救援的命令行解释器|
|2| runlevel2.target, multi-user.target |多用户，没有 NFS，但是运行其他所有的非 GUI 服务 |
|3| runlevel3.target, multi-user.target |多用户的文本界面: 运行所有的服务，但是只有命令行界面（CLI）|
|4| runlevel4.target, multi-user.target |未使用, 等同于级别 3 |
|5| runlevel5.target, graphical.target |带gui的multi-user.target|
|6| runlevel6.target, reboot.target |重启|
|S|emergency |emergency.target |单用户模式 —— 没有服务运行；文件系统没有挂载. 这是最基础级的操作模式, 只有一个运行在主控制台的用于紧急情况的命令行解释器，供用户和系统交互 |
|||halt.target|关机|

修改默认运行级别: `ln -sf /lib/systemd/system/multi-user.target /etc/systemd/system/default.target`, default.target 文件是指向真实的目标文件的符号链接.


# 配置

## systemd 文件类型及存放路径
systemd配置文件被称为unit单元，根据类型不同，以不同的扩展名结尾:
- `.service` 系统服务；
- `.target` 一组系统服务(多个 Unit 构成的一个组)；

	每个目标在配置文件中都描述了一组依赖关系. systemd 启动需要的依赖，即 Linux 主机运行在特定功能级别所需的服务.
- `.automount` 自动挂载点；
- `.device` 能被内核识别的硬件设备；
- `.mount` 文件系统的挂载点；
- `.path` 文件系统的文件或者目录；
- `.scope` 不是由 Systemd 启动的外部进程；
- `.slice` 一组分层次管理的进程组；
- `.snapshot` Systemd 快照，可以切回某个快照；
- `.socket` 进程间通信的 socket；
- `.swap` swap 文件；
- `.timer`  定时器

systemd单元文件放置位置:
- `/lib/systemd/system/` <=> `/usr/lib/systemd/system/` systemd默认单元文件安装目录,**推荐**使用`/lib/systemd/system/`
- `/run/systemd/system/` systemd运行时创建的服务脚本，这个目录优先于安装目录
- `/etc/systemd/system/` 系统管理员创建和管理的单元目录，优先级最高(大部分文件都是符号链接，指向systemd默认单元文件安装目录)

## Unit文件
Systemd为每一个守护进程记录一个初始化结构文件，我们称之为一个Unit。Systemd系统取代了传统系统为每一个守护进程初始化一次脚本的做法。

每个配置文件的状态，一共有四种:
- enabled：已建立启动链接
- disabled：没建立启动链接
- static：该配置文件没有[Install]部分, 不可以自己启动, 不过可能会被其它的 enabled 的服务来唤醒
- masked：这个 unit 无论如何都无法被启动！因为已经被强制注销. 可通过 systemctl unmask 改回原来的状态

下面我们演示创建一个Hello_world.service的简单Unit文件：
```
[Unit]
Description=HelloWorldApp
After=docker.service
Requires=docker.service
AssertFileIsExecutable=/usr/bin/docker # 断言是否可执行文件，不常用

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker rm busybox1
ExecStartPre=/usr/bin/docker pull busybox
ExecStart=/usr/bin/docker run --rm --name busybox1 busybox /bin/sh -c "while true; do echo Hello World; sleep 1; done"
ExecStop=/usr/bin/docker stop busybox1

[Install]
WantedBy=multi-user.target
```

### Var
- MAINPID : 服务的systemd变量, 它指向主应用程序的PID

	- `systemctl  show --property MainPID --value NetworkManager.service`

### `[Unit]`
[Unit]区块通常是配置文件的第一个区块，用来定义 Unit 的元数据，以及配置与其他 Unit 的关系.它的主要字段如下:
- Description：简短描述,systemd服务启停时会打印到Systemd日志（可以通过journalctl和Systemdctl来检查日志文件）
- Documentation：文档地址
- Requires：当前 Unit 依赖的其他 Unit(强依赖)，如果它们没有运行，当前 Unit 会启动失败
- Wants：与当前 Unit 配合的其他 Unit(弱依赖)，如果它们没有运行，当前 Unit 不会启动失败
- BindsTo：与Requires类似，它指定的 Unit 如果退出，会导致当前 Unit 停止运行
- Before：如果该字段指定的 Unit 也要启动，那么必须在当前 Unit 之后启动
- After：如果该字段指定的 Unit 也要启动，那么必须在当前 Unit 之前启动
- Conflicts：这里指定的 Unit 不能与当前 Unit 同时运行
- Condition...：当前 Unit 运行必须满足的条件，否则不会运行
- Assert...：当前 Unit 运行必须满足的条件，否则会报启动失败
- StartLimitIntervalSec=400. 用于设置时长， 默认值等于 DefaultStartLimitIntervalSec= 的值(默认为10秒)，设为 0 表示不作限制

	> Before systemd-230 it was called just StartLimitInterval.
- StartLimitBurst=3. 用于设置在一段给定的时长内，最多允许启动多少次.  默认值等于 DefaultStartLimitBurst= 的值(默认为5次)

> **After和Before字段只涉及启动顺序，不涉及依赖关系**; **Wants字段与Requires字段只涉及依赖关系，与启动顺序无关**

	Requires时, 依赖service停止, 当前服务也会停止.
> 自动重启逻辑的单元触碰到了启动频率限制，那么该单元将再也不会尝试自动重启. systemctl reset-failed 命令能够重置单元的启动频率计数器.

### `[Install]`
[Install]通常是配置文件的最后一个区块，用来定义如何安装这个配置文件即怎样做到开机启动. 它的主要字段如下:
- WantedBy：表示该服务所在的 Target,它的值是一个或多个 Target，当前 Unit 激活时（enable）符号链接会放入/etc/systemd/system目录下面以 Target 名 + .wants后缀构成的子目录中
- RequiredBy：它的值是一个或多个 Target，当前 Unit 激活时，符号链接会放入/etc/systemd/system目录下面以 Target 名 + .required后缀构成的子目录中
- Alias：当前 Unit 可用于启动的别名
- Also：当前 Unit 激活（enable）时，会被同时激活的其他 Unit

> 事实上，systemd在运行时并不使用此小节. 只有 systemctl 的 enable 与 disable 命令在启用/停用单元时才会使用此小节.

### `[Service]`
ref:
- [systemd.service 中文手册](http://www.jinbuguo.com/systemd/systemd.service.html)

[Service]区块用来 Service 的配置，只有 Service 类型的 Unit 才有这个区块。它的主要字段如下:
- Type：定义启动时的进程行为。它有以下几种值。
- Type=simple：默认值，执行ExecStart指定的命令，启动主进程
- Type=forking：以 fork 方式从父进程创建子进程，创建后父进程会立即退出
- Type=oneshot：一次性进程，Systemd 会等当前服务退出，再继续往下执行
- Type=dbus：当前服务通过D-Bus启动
- Type=notify：与 Type=simple 相同, 但当前服务启动完成(systemctl start会阻塞直至服务就绪)，会通知Systemd，再继续往下执行

	在默认情况下(也就是调用此命令的服务单元含有 NotifyAccess=none)， systemd 并不从此命令接受状态更新消息, 因此需要NotifyAccess=all
- Type=idle：若有其他任务执行完毕，当前服务才会运行
- ExecStart：运行这个单元最主要的命令（这个参数是几乎每个 .service 文件都会有的，指定服务启动的主要命令，在每个配置文件中**只能使用一次**）.
- ExecStartPre：ExecStart之前运行的命令（指定在启动执行 ExecStart 的命令前的准备工作，**可以有多个**，所有命令会按照文件中书写的顺序**依次**被执行）.
- ExecStartPost：ExecStart 运行完成后要执行的命令（指定在启动执行 ExecStart 的命令后的收尾工作，也**可以有多个**）.
- ExecReload：当使用 systemctl重新加载服务所需执行的命令
- ExecStop：通过执行systemctl 确认单元服务失败或者停止时执行的命令
- ExecStopPost：ExecStop 执行完成之后所执行的命令（指定在 ExecStop 命令执行后的收尾工作，也**可以有多个**）.
- RestartSec：在重启服务之前系统休眠时间（很有效的防止失败服务少于100ms重启一次;如果服务需要被重启，这个参数的值为服务被重启前的等待秒数）.
- Restart：定义何种情况 Systemd 会自动重启当前服务，可能的值包括always（总是重启）、on-success、on-failure、on-abnormal、on-abort、on-watchdog

    - no（默认值）：退出后不会重启
    - on-success：只有正常退出时（退出状态码为0），才会重启
    - on-failure：非正常退出时（退出状态码非0），包括被信号终止和超时，才会重启
    - on-abnormal：只有被信号终止和超时，才会重启
    - on-abort：只有在收到没有捕捉到的信号终止时，才会重启
    - on-watchdog：超时退出，才会重启
    - always：不管是什么退出原因，总是重启

    > 对于守护进程，推荐设为on-failure。对于那些允许发生错误退出的服务，可以设为on-abnormal

- TimeoutSec：定义 Systemd 停止当前服务之前等待的秒数
- Environment：指定环境变量

	```conf
	Environment=ETCD_CA_FILE=/path/to/CA.pem
	Environment=ETCD_CERT_FILE=/path/to/server.crt
	```
	或者
	```conf
	Environment="ETCD_CA_FILE=/path/to/CA.pem" "ETCD_CERT_FILE=/path/to/server.crt"
	```
- TimeoutStartSec: 设置该服务允许的最大启动时长。 如果守护进程未能在限定的时长内发出"启动完毕"的信号，那么该服务将被视为启动失败，并会被关闭. 如果一个 Type=notify 服务发送了 "EXTEND_TIMEOUT_USEC=…" 信号， 那么允许的启动时长将会在 TimeoutStartSec= 基础上继续延长指定的时长.
- RestartPreventExitStatus: 当符合某些退出状态时不要进行重启, 比如`123 SIGTERM`(123是return code)
- PrivateTmp : 私有tmp, 比如php-fpm配置它后会有一个`/tmp/systemd-private-xxx-php-fpm.service-yyy/tmp`目录
- PIDFile: systemd会在中间父进程退出后立即读取这个PID文件，读取成功后就认为该服务已经启动成功. 如果systemd读取PIDFile的时候，服务主进程可能还未将PID写入到PID文件中，这时systemd将出现问题. 更详细可参考[systemd service：Type=forking](https://www.junmajinlong.com/linux/systemd/service_2/)

	服务停止时该文件会被删除.

> [Unit 配置文件的完整字段清单](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)
> 所有的启动设置之前，都可以加上一个连词号（-），表示"抑制错误"，即发生错误的时候，不影响其他命令的执行。比如，EnvironmentFile=-/etc/sysconfig/sshd（注意等号后面的那个连词号），就表示即使/etc/sysconfig/sshd文件不存在，也不会抛出错误.

### 启动类型Type
Type字段定义启动类型。它可以设置的值如下:
- simple（默认值）：ExecStart字段启动的进程为主进程
- forking：ExecStart字段将以fork()方式启动，此时父进程将会退出，子进程将成为主进程

	Type=forking时的systemd提供了PIDFile指令(PIDFile只适合在Type=forking模式下使用，其它时候没必要使用，因为其它类型的Service主进程的PID都是确定的), systemd会从PIDFile指令所指定的PID文件中获取服务的主进程PID. 且PIDFile指令的值要和服务程序的PID文件路径保持一致.
- oneshot：类似于simple，但只执行一次，Systemd 会等它执行完，才启动其他服务
- dbus：类似于simple，但会等待 D-Bus 信号后启动
- notify：类似于simple，启动结束后会发出通知信号，然后 Systemd 再启动其他服务
- idle：类似于simple，但是要等到其他任务都执行完，才会启动该服务。一种使用场合是为让该服务的输出，不与其他服务的输出相混合

### 重启行为
KillMode字段定义 Systemd 如何停止服务,可以设置的值如下:
- control-group（默认值）：当前控制组里面的所有子进程，都会被杀掉
- process：只杀主进程
- mixed：主进程将收到 SIGTERM 信号，子进程收到 SIGKILL 信号
- none：没有进程会被杀掉，只是执行服务的 stop 命令。

Restart字段可以设置的值如下:
- no（默认值）：退出后不会重启
- on-success：只有正常退出时（退出状态码为0），才会重启
- on-failure：非正常退出时（退出状态码非0），包括被信号终止和超时，才会重启
- on-abnormal：只有被信号终止和超时，才会重启
- on-abort：只有在收到没有捕捉到的信号终止时，才会重启
- on-watchdog：超时退出，才会重启
- always：不管是什么退出原因，总是重启

## target
Target 就是一个 Unit 组，包含许多相关的 Unit 。启动某个 Target 的时候，Systemd 就会启动里面所有的 Unit.

Target 与 传统 RunLevel 的对应关系如下:

Traditional runlevel      New target name     Symbolically linked to...

Runlevel 0           |    runlevel0.target -> poweroff.target
Runlevel 1           |    runlevel1.target -> rescue.target
Runlevel 2           |    runlevel2.target -> multi-user.target
Runlevel 3           |    runlevel3.target -> multi-user.target
Runlevel 4           |    runlevel4.target -> multi-user.target
Runlevel 5           |    runlevel5.target -> graphical.target
Runlevel 6           |    runlevel6.target -> reboot.target

systemd与init进程的主要差别如下:
- 默认的 RunLevel（/etc/inittab -> /etc/systemd/system/default.target)
- 启动脚本的位置(/etc/init.d -> /lib/systemd/system和/etc/systemd/system)
- 配置文件的位置.init进程的配置文件是/etc/inittab，各种服务的配置文件存放在/etc/sysconfig;现在的配置文件主要存放在/lib/systemd目录，在/etc/systemd目录里面的修改可以覆盖原始设置.

## 启动优化
```sh
# systemd-cgls                            显示cgroup树
# systemd-analyze                         ← 查看系统引导用时
# systemd-analyze time                    ← 同上
# systemd-analyze blame                   ← 查看初始化任务所消耗的时间
# systemd-analyze plot > systemd.svg      ← 将启动过程输出为svg图
# systemd-cgtop                           ← 查看资源的消耗状态
```

## [socket](http://www.jinbuguo.com/systemd/systemd.socket.html)
每个socket单元都必须有一个与其匹配的服务单元(详见 systemd.service(5) 手册)， 以处理该套接字上的接入流量. 匹配的 .service 单元名称默认与对应的 .socket 单元相同， 但是也可以通过 Service= 选项(见下文)明确指定

### 参考
- [CentOS7/RHEL7 systemd详解](http://xiaoli110.blog.51cto.com/1724/1629533)
- [Systemd 系列中文手册](http://www.jinbuguo.com/systemd/index.html)

### log
system-journal服务通过监听`socket /dev/log`(`/dev/log -> /run/systemd/journal/dev-log=`)来获取日志并保存到内存里, 再间隙性写入到`/var/log/journal`目录中.

rsyslog 服务启动后监听`socket /run/systemd/journal/syslog`筛选分类, 并写入`/var/log/messages`文件中.

> systemd-journald是一种改进的日志管理服务，是 syslog 的补充，收集来自内核、启动过程早期阶段、标准输出、系统日志，守护进程启动和运行期间错误的信息
> 默认情况下，systemd 的日志保存在 /run/log/journal 中，系统重启就会清除，这是RHEL7的新特性. 通过新建 /var/log/journal 目录，日志会自动记录到这个目录中，并永久存储.

日志流转: 应用进程将日志通过/run/systemd/journal/dev-log发送到systemd， 然后systemd 再将日志通过/run/systemd/journal/syslog发送到rsyslogd, 具体如下:
```
[log management with systemd](https://unix.stackexchange.com/questions/205883/understand-logging-in-linux)
systemd has a single monolithic log management program, systemd-journald. This runs as a service managed by systemd.

It reads /dev/kmsg for kernel log data.
It reads /dev/log (a symbolic link to /run/systemd/journal/dev-log) for application log data from the GNU C library's syslog() function.
It listens on the AF_LOCAL stream socket at /run/systemd/journal/stdout for log data coming from systemd-managed services.
It listens on the AF_LOCAL datagram socket at /run/systemd/journal/socket for log data coming from programs that speak the systemd-specific journal protocol (i.e. sd_journal_sendv() et al.).
It mixes these all together.
It writes to a set of system-wide and per-user journal files, in /run/log/journal/ or /var/log/journal/.
If it can connect (as a client) to an AF_LOCAL datagram socket at /run/systemd/journal/syslog it writes journal data there, if forwarding to syslog is configured.
If configured, it writes journal data to the kernel buffer using the writable /dev/kmsg mechanism.
If configured, it writes journal data to terminals and the console device as well.
```

## 命令
### systemctl
systemctl是 Systemd 的主命令，用于管理系统
```
# 重启系统
$ sudo systemctl reboot

# 关闭系统，切断电源
$ sudo systemctl poweroff

# CPU停止工作
$ sudo systemctl halt

# 暂停系统
$ sudo systemctl suspend

# 让系统进入冬眠状态
$ sudo systemctl hibernate

# 让系统进入交互式休眠状态
$ sudo systemctl hybrid-sleep

# 启动进入救援状态（单用户状态）
$ sudo systemctl rescue
```

### systemd-cgtop
查看资源的消耗状态

### systemd-cgls
查看cgroup层次结构（以及组成单元的详细信息）. 可用于追踪父子进程关系(即使exec时设置了setsid)

或使用`systemctl status httpd.service`里显示的`CGroup`属性

`systemd-cgls memory`可顺便显示memory control信息

`systemctl show httpd`会显示ControlGroup信息

### systemd-analyze
systemd-analyze命令用于查看启动耗时
```
# 查看启动耗时
$ systemd-analyze                                                                                       

# 查看每个服务的启动耗时
$ systemd-analyze blame

# 绘制启动矢量图，得到各service启动顺序
$ systemd-analyze plot > boot.svg

# 显示内核和普通用户空间启动时所花的时间
$ systemd-analyze time

# 显示在所有系统单元中是否有语法错误
$ systemd-analyze verify

# 显示瀑布状的启动过程流
$ systemd-analyze critical-chain

# 显示指定服务的启动流
$ systemd-analyze critical-chain atd.service
```

### hostnamectl
hostnamectl命令用于查看当前主机的信息
```
# 显示当前主机的信息
$ hostnamectl

# 设置主机名。
$ sudo hostnamectl set-hostname rhel7
```

### localectl
localectl命令用于查看本地化设置
```
# 查看本地化设置
$ localectl

# 设置本地化参数。
$ sudo localectl set-locale LANG=en_GB.utf8
$ sudo localectl set-keymap en_GB
```

### timedatectl
timedatectl命令用于查看当前时区设置
```
# 查看当前时区设置
$ timedatectl

# 显示所有可用的时区
$ timedatectl list-timezones                                                                                   
# 设置当前时区
$ sudo timedatectl set-timezone America/New_York # Asia/Shanghai
$ sudo timedatectl set-time YYYY-MM-DD
$ sudo timedatectl set-time HH:MM:SS
```

### loginctl
loginctl命令用于查看当前登录的用户
```
# 列出当前session
$ loginctl list-sessions

# 列出当前登录用户
$ loginctl list-users

# 列出显示指定用户的信息
$ loginctl show-user ruanyf
```

### systemctl
```
###unit状态###
# 列出正在运行的 Unit
$ systemctl list-units

# 列出所有Unit，包括没有找到配置文件的或者启动失败的
$ systemctl list-units --all

# 列出所有没有运行的 Unit
$ systemctl list-units --all --state=inactive

# 列出所有加载失败的 Unit
$ systemctl list-units --failed

# 列出所有正在运行的、类型为 service 的 Unit
$ systemctl list-units --type=service

# 列出当前系统支持的所有等级
$ systemctl list-units --type=target

# 显示某个 Unit 是否正在运行
$ systemctl is-active application.service

# 显示某个 Unit 是否处于启动失败状态
$ systemctl is-failed application.service

# 显示某个 Unit 服务是否建立了启动链接
$ systemctl is-enabled application.service

###管理###
# 立即启动一个服务
$ sudo systemctl start apache.service

# 立即停止一个服务
$ sudo systemctl stop apache.service

# 重启一个服务
$ sudo systemctl restart apache.service

# 杀死一个服务的所有子进程
$ sudo systemctl kill apache.service

# 重新加载一个服务的配置文件
$ sudo systemctl reload apache.service

# 重载所有修改过的配置文件
$ sudo systemctl daemon-reload

# 显示某个 Unit 的所有底层参数
$ systemctl show httpd.service

# 显示某个 Unit 的指定属性的值
$ systemctl show -p CPUShares httpd.service

# 查看一个服务的状态
$ sudo systemctl status httpd

>上面的输出结果含义如下:
Loaded行：配置文件的位置，是否设为开机启动
Active行：systemd提供的状态
	- inactive : 未运行
	- active (running) : 有一个或多个进程正在执行
	- active (exited) : 仅执行一次就正常结束的服务, 但有后台进程正在继续执行
	- active (waiting) : 正在执行, 但正在等待其他事件才能继续执行.

Main PID行：主进程ID
Status行：由应用本身（这里是 httpd ）提供的软件当前状态
CGroup块：应用的所有子进程
日志块：应用的日志

# 设置某个 Unit 的指定属性
$ sudo systemctl set-property httpd.service CPUShares=500

###依赖关系###
# 列出一个 Unit 的所有依赖
$ systemctl list-dependencies nginx.service [--reverse] # `--reverse`会反向追踪是谁在使用这个 unit. output中green是活动的, red是非活动的.
# 有些依赖是 Target 类型（详见下文），默认不会展开显示。如果要展开 Target，就需要使用`--all`参数
$ systemctl list-dependencies --all nginx.service

###开机启动###
# 设置开机启动
$ sudo systemctl enable clamd@scan.service
# 撤销开机启动
$ sudo systemctl disable clamd@scan.service
# 列出开机启动项
$ sudo systemctl list-unit-files --type=service|grep enabled

###配置文件###
# 列出所有配置文件
$ systemctl list-unit-files

# 列出指定类型的配置文件, 及其启动与禁用情况
$ systemctl list-unit-files --type=service

# 查看配置文件的内容
$ systemctl cat atd.service

# 重新加载配置文件
$ sudo systemctl daemon-reload

###target###
# 查看当前系统的所有 Target
$ systemctl list-unit-files --type=target

# 查看一个 Target 包含的所有 Unit
$ systemctl list-dependencies multi-user.target

# 查看启动时的默认 Target
$ systemctl get-default

# 设置启动时的默认 Target
$ sudo systemctl set-default multi-user.target

# 切换 Target 时，默认不关闭前一个 Target 启动的进程，
# systemctl isolate 命令改变这种行为，
# 关闭前一个 Target 里面所有不属于后一个 Target 的进程
$ sudo systemctl isolate multi-user.target
```

## FAQ
### systemd debug
编辑`/etc/systemd/system.conf`, 将`LogLevel`设为`debug`并重启即可.

### 升级工具upgrader在停止system的服务时被kill
upgrader由该服务启动, 即使使用了setsid来脱离从父进程继承而来的已打开的终端、隶属进程组和隶属的会话, 但仍在该服务的cgroup里, 导致停止服务时, 升级工具会被kill.

解决方法: 通过systemd-run创建临时cgroup来解决: [`systemd-run [--remain-after-exit] --unit=my_system_upgrade --scope --slice=my_system_upgrade_slice -E  setsid start-the-upgrade > /tmp/some-logs.log 2>&1`](https://stackoverflow.com/questions/35200232/how-to-launch-a-process-outside-a-systemd-control-group), 里面无需再使用nohup.

> [transient cgroup with systemd-run](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/resource_management_guide/chap-using_control_groups#sec-Creating_Transient_Cgroups_with_systemd-run)

>  `--remain-after-exit`: 可以将命令执行的stdout存入systemd log. 与`--scope`互斥

> 在oracle linux 7.9 x64上systemd-run执行命令(执行bin文件, 其中通过bash执行了其他脚本, 且脚本逻辑正确. 但直接执行bin正常)时可能会莫名退出, 加`> /tmp/some-logs.log 2>&1`或`--remain-after-exit`可解决问题. 推测可能与`systemd-run`的`--wait`有关: [systemd-run 将会以异步模式在后台启动临时服务并在命令开始执行之后返回](http://www.jinbuguo.com/systemd/systemd-run.html)

### systemd 计时器
```config
[Unit]
Description=Backup
Requires=myBackup.service
[Timer]
OnBootSec=30min
OnUnitActiveSec=1d
[Install]
WantedBy=timers.target
```

### systemd service获取不到所需环境变量
ref:
- [Make systemd service inherit environment variables from /etc/profile.d](https://unix.stackexchange.com/questions/675521/make-systemd-service-inherit-environment-variables-from-etc-profile-d)

systemd service的环境变量与shell env无关, 因此systemd提供了Environment, EnvironmentFile.

见`man 5 systemd.exec`.

> `ExecStartPre=/usr/bin/bash xxx.sh`无用, 因为执行ExecStart前ExecStartPre已结束.

### Failed with result 'start-limit-hit'
是由于服务定义中的 Restart=on-failure 指令造成的, 服务总是崩溃重启导致启动次数达到了上限