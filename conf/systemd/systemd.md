# systemd
-[Systemd 入门教程：命令篇](http://www.ruanyifeng.com/blog/2016/03/systemd-tutorial-commands.html)

## 配置

## systemd 文件类型及存放路径
systemd配置文件被称为unit单元，根据类型不同，以不同的扩展名结尾:
- `.service` 系统服务；
- `.target` 一组系统服务(多个 Unit 构成的一个组)；
- `.automount` 自动挂载点；
- `.device` 能被内核识别的硬件设备；
- `.mount` 文件系统的挂载点；
- `.path` 文件系统的文件或者目录；
- `.scope` 不是由 Systemd 启动的外部进程；
- `.slice` 一组分层次管理的进程组；
- `.snapshot` Systemd 快照，可以切回某个快照；
- `.socket` 进程间通信的 socket；
- `.swap` swap 文件；
- `.timer`	定时器。

systemd单元文件放置位置:
- `/lib/systemd/system/` <=> `/usr/lib/systemd/system/` systemd默认单元文件安装目录,**推荐**使用`/lib/systemd/system/`
- `/run/systemd/system/` systemd运行时创建的服务脚本，这个目录优先于安装目录
- `/etc/systemd/system/` 系统管理员创建和管理的单元目录，优先级最高(大部分文件都是符号链接，指向systemd默认单元文件安装目录)

## Unit文件
Systemd为每一个守护进程记录一个初始化结构文件，我们称之为一个Unit。Systemd系统取代了传统系统为每一个守护进程初始化一次脚本的做法。

每个配置文件的状态，一共有四种:
- enabled：已建立启动链接
- disabled：没建立启动链接
- static：该配置文件没有[Install]部分（无法执行），只能作为其他配置文件的依赖
- masked：该配置文件被禁止建立启动链接

下面我们演示创建一个Hello_world.service的简单Unit文件：
```
[Unit]
Description=HelloWorldApp
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker rm busybox1
ExecStartPre=/usr/bin/docker pull busybox
ExecStart=/usr/bin/docker run --rm --name busybox1 busybox /bin/sh -c "while true; do echo Hello World; sleep 1; done"
ExecStop=/usr/bin/docker stop busybox1

[Install]
WantedBy=multi-user.target
```

### `[Unit]`
[Unit]区块通常是配置文件的第一个区块，用来定义 Unit 的元数据，以及配置与其他 Unit 的关系.它的主要字段如下:
- Description：简短描述,可以在Systemd日志展示（可以通过journalctl和Systemdctl来检查日志文件）
- Documentation：文档地址
- Requires：当前 Unit 依赖的其他 Unit(强依赖)，如果它们没有运行，当前 Unit 会启动失败
- Wants：与当前 Unit 配合的其他 Unit(弱依赖)，如果它们没有运行，当前 Unit 不会启动失败
- BindsTo：与Requires类似，它指定的 Unit 如果退出，会导致当前 Unit 停止运行
- Before：如果该字段指定的 Unit 也要启动，那么必须在当前 Unit 之后启动
- After：如果该字段指定的 Unit 也要启动，那么必须在当前 Unit 之前启动
- Conflicts：这里指定的 Unit 不能与当前 Unit 同时运行
- Condition...：当前 Unit 运行必须满足的条件，否则不会运行
- Assert...：当前 Unit 运行必须满足的条件，否则会报启动失败

> 注意，After和Before字段只涉及启动顺序，不涉及依赖关系

### `[Install]`
[Install]通常是配置文件的最后一个区块，用来定义如何启动，以及是否开机启动。它的主要字段如下:
- WantedBy：表示该服务所在的 Target,它的值是一个或多个 Target，当前 Unit 激活时（enable）符号链接会放入/etc/systemd/system目录下面以 Target 名 + .wants后缀构成的子目录中
- RequiredBy：它的值是一个或多个 Target，当前 Unit 激活时，符号链接会放入/etc/systemd/system目录下面以 Target 名 + .required后缀构成的子目录中
- Alias：当前 Unit 可用于启动的别名
- Also：当前 Unit 激活（enable）时，会被同时激活的其他 Unit

### `[Service]`
[Service]区块用来 Service 的配置，只有 Service 类型的 Unit 才有这个区块。它的主要字段如下:
- Type：定义启动时的进程行为。它有以下几种值。
- Type=simple：默认值，执行ExecStart指定的命令，启动主进程
- Type=forking：以 fork 方式从父进程创建子进程，创建后父进程会立即退出
- Type=oneshot：一次性进程，Systemd 会等当前服务退出，再继续往下执行
- Type=dbus：当前服务通过D-Bus启动
- Type=notify：当前服务启动完毕，会通知Systemd，再继续往下执行
- Type=idle：若有其他任务执行完毕，当前服务才会运行
- ExecStart：运行这个单元最主要的命令（这个参数是几乎每个 .service 文件都会有的，指定服务启动的主要命令，在每个配置文件中**只能使用一次**）.
- ExecStartPre：ExecStart之前运行的命令（指定在启动执行 ExecStart 的命令前的准备工作，**可以有多个**，所有命令会按照文件中书写的顺序**依次**被执行）.
- ExecStartPost：ExecStart 运行完成后要执行的命令（指定在启动执行 ExecStart 的命令后的收尾工作，也**可以有多个**）.
- ExecReload：当使用 systemctl重新加载服务所需执行的命令
- ExecStop：通过执行systemctl 确认单元服务失败或者停止时执行的命令
- ExecStopPost：ExecStop 执行完成之后所执行的命令（指定在 ExecStop 命令执行后的收尾工作，也**可以有多个**）.
- RestartSec：在重启服务之前系统休眠时间（很有效的防止失败服务少于100ms重启一次;如果服务需要被重启，这个参数的值为服务被重启前的等待秒数）.
- Restart：定义何种情况 Systemd 会自动重启当前服务，可能的值包括always（总是重启）、on-success、on-failure、on-abnormal、on-abort、on-watchdog
- TimeoutSec：定义 Systemd 停止当前服务之前等待的秒数
- Environment：指定环境变量

> [Unit 配置文件的完整字段清单](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)
> 所有的启动设置之前，都可以加上一个连词号（-），表示"抑制错误"，即发生错误的时候，不影响其他命令的执行。比如，EnvironmentFile=-/etc/sysconfig/sshd（注意等号后面的那个连词号），就表示即使/etc/sysconfig/sshd文件不存在，也不会抛出错误.

### 启动类型Type
Type字段定义启动类型。它可以设置的值如下:
- simple（默认值）：ExecStart字段启动的进程为主进程
- forking：ExecStart字段将以fork()方式启动，此时父进程将会退出，子进程将成为主进程
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
# systemd-cgls                            - 启动的结构树
# systemd-analyze                         ← 查看系统引导用时
# systemd-analyze time                    ← 同上
# systemd-analyze blame                   ← 查看初始化任务所消耗的时间
# systemd-analyze plot > systemd.svg      ← 将启动过程输出为svg图
# systemd-cgtop                           ← 查看资源的消耗状态
```

### 参考
- [CentOS7/RHEL7 systemd详解](http://xiaoli110.blog.51cto.com/1724/1629533)
- [Systemd 系列中文手册](http://www.jinbuguo.com/systemd/index.html)
