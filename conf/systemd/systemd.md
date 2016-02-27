## systemd 文件类型及存放路径
systemd配置文件被称为unit单元，根据类型不同，以不同的扩展名结尾:
- `.service` 系统服务；
- `.target` 一组系统服务；
- `.automount` 自动挂载点；
- `.device` 能被内核识别的设备；
- `.mount` 挂载点；
- `.path` 文件系统的文件或者目录；
- `.scope` 外部创建的进程；
- `.slice` 一组分层次管理的系统进程；
- `.snapshot` 系统服务状态管理；
- `.socket` 进程间通讯套接字；
- `.swap` 定义swap文件或者设备；
- `.timer`	定义定时器。

systemd单元文件放置位置:
- `/lib/systemd/system/` <=> `/usr/lib/systemd/system/` systemd默认单元文件安装目录,**推荐**使用`/lib/systemd/system/`
- `/run/systemd/system/` systemd运行时创建的服务脚本，这个目录优先于安装目录
- `/etc/systemd/system/` 系统管理员创建和管理的单元目录，优先级最高

## Unit文件
Systemd为每一个守护进程记录一个初始化结构文件，我们称之为一个Unit。Systemd系统取代了传统系统为每一个守护进程初始化一次脚本的做法。

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

下面列表，是按照顺序一个服务Unit（service unit）生命周期内的过程，以及作用：
- ExecStartPre: ExecStart之前运行的命令（指定在启动执行 ExecStart 的命令前的准备工作，**可以有多个**，所有命令会按照文件中书写的顺序**依次**被执行）.
- ExecStart:运行这个单元最主要的命令（这个参数是几乎每个 .service 文件都会有的，指定服务启动的主要命令，在每个配置文件中**只能使用一次**）.
- ExecStartPost: ExecStart 运行完成后要执行的命令（指定在启动执行 ExecStart 的命令后的收尾工作，也**可以有多个**）.
- ExecReload:当使用 systemctl重新加载服务所需执行的命令
- ExecStop: 通过执行systemctl 确认单元服务失败或者停止时执行的命令
- ExecStopPost: ExecStop 执行完成之后所执行的命令（指定在 ExecStop 命令执行后的收尾工作，也**可以有多个**）.
- RestartSec:在重启服务之前系统休眠时间（很有效的防止失败服务少于100ms重启一次;如果服务需要被重启，这个参数的值为服务被重启前的等待秒数）.

- Description： 可以在Systemd日志展示（可以通过journalctl和Systemdctl来检查日志文件）
- After=docker.service 和 Requires=docker.service 意思是只有docker.service启动后才可以被激活。当然也可以使用多个After来做限制，多个after使用空格分开即可.
- =- 是Systemd忽略错误的一个语法.在这种情况下，docker会发挥一个非0的退出码，如果想停止一个不存在的容器，这对我们来说并不是一个错误。
- WantedBy="multi-user.target" 当启动 multi-user.target时，Systemd将会将获取单元.

### 参考
- [CentOS7/RHEL7 systemd详解](http://xiaoli110.blog.51cto.com/1724/1629533)
- [Systemd 系列中文手册](http://www.jinbuguo.com/systemd/index.html)
