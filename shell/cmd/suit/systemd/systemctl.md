# systemctl

## 描述

Systemctl是一个systemd工具，主要负责控制systemd系统和服务管理器。

Systemd是一个系统管理守护进程、工具和库的集合，用于取代System V初始进程。Systemd的功能是用于集中管理和配置类UNIX系统。

在Linux生态系统中，Systemd被部署到了大多数的标准Linux发行版中，只有为数不多的几个发行版尚未部署。Systemd通常是所有其它守护进程的父进程，但并非总是如此。

>本机OS: ubuntu 15.04 amd64

### Systemd初体验和Systemctl基础

1. 首先检查你的系统中是否安装有systemd并确定当前安装的版本

       # systemd --version
       systemd 219
       +PAM +AUDIT +SELINUX +IMA +APPARMOR +SMACK +SYSVINIT +UTMP +LIBCRYPTSETUP +GCRYPT -GNUTLS +ACL +XZ -LZ4 -SECCOMP +BLKID -ELFUTILS +KMOD -ID

 上例中很清楚地表明，我们安装了219版本的systemd。

1. 检查systemd和systemctl的二进制文件和库文件的安装位置

       # whereis systemd
       systemd: /usr/lib/systemd /bin/systemd /etc/systemd /lib/systemd /usr/share/systemd /usr/share/man/man1/systemd.1.gz
       # whereis systemctl
       systemctl: /bin/systemctl /usr/share/man/man1/systemctl.1.gz

1. 检查systemd是否运行

 ```shell
# ps -eaf | grep [s]ystemd
root       273     1  0 08:59 ?        00:00:00 /lib/systemd/systemd-journald
root       301     1  0 08:59 ?        00:00:00 /lib/systemd/systemd-udevd
systemd+   622     1  0 08:59 ?        00:00:00 /lib/systemd/systemd-timesyncd
root       789     1  0 08:59 ?        00:00:00 /lib/systemd/systemd-logind
message+   815     1  0 08:59 ?        00:00:00 /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation
chen      2151     1  0 09:00 ?        00:00:00 /lib/systemd/systemd --user
```

 注意：systemd是作为父进程（PID=1）运行的。在上面带（-e）参数的ps命令输出中，选择所有进程，（-a）选择除会话前导外的所有进程，并使用（-f）参数输出完整格式列表（即 -eaf）。

 也请注意上例中后随的方括号和例子中剩余部分。方括号表达式是grep的字符类表达式的一部分。

1. 分析systemd启动进程

       # systemd-analyze
       Startup finished in 2.791s (kernel) + 45.161s (userspace) = 47.953s

1. 分析启动时各个进程花费的时间

 ```shell
# systemd-analyze blame
  15.399s gpu-manager.service
  14.272s plymouth-quit-wait.service
  13.446s mysql.service
  10.750s NetworkManager-wait-online.service
  7.665s NetworkManager.service
  7.590s systemd-udev-settle.service
  ....
```
1. 分析启动时的关键链
       # systemd-analyze critical-chain xxx.service
 重要：Systemctl接受服务（.service），挂载点（.mount），套接口（.socket）和设备（.device）作为单元

1. 启动时序图: systemd-analyze plot

       关键路径用红色高亮显示

       **注意**: 仅显示启动成功的, 需要留意某些使用了Restart服务, 比如第一次启动失败, 之后启动成功.
1. systemd-analyze dot | dot -Tsvg > /tmp/test.svg

       生成依赖图
1. systemd-analyze condition xxx

       测试 systemd 单元文件中的条件和断言.
1. systemd-analyze verify [/etc/systemd/system/backup.service]

       检查语法是否正确
1. systemd-analyze cat-config systemd/system/display-manager.service

       输出配置文件

1. 列出所有可用单元
 ```shell
# systemctl list-unit-files [--type=service --state=enabled]
UNIT FILE                                  STATE
proc-sys-fs-binfmt_misc.automount          static
dev-hugepages.mount                        static
dev-mqueue.mount                           static
proc-sys-fs-binfmt_misc.mount              static
sys-fs-fuse-connections.mount              static
sys-kernel-config.mount                    static
sys-kernel-debug.mount                     static
tmp.mount                                  disabled
cups.path                                  enabled
systemd-ask-password-conso
.....
```
1. 列出所有运行中单元
```shell
# systemctl list-units # 同systemctl
# systemctl list-units --type service     ← 同上，只是以service为单位
```
1. 列出所有失败单元
       # systemctl --failed

1. 检查某个单元（如 cron.service）是否开机时启动
       # systemctl is-enabled crond.service

1. 检查某个单元或服务是否运行
       # systemctl status firewalld.service

1. 刷新systemd的配置
        # systemctl daemon-reload # service连续启动失败多次后, 再start会直接报错, 应先用daemon-reload重置其状态, 再start
        # systemctl reload apache2 # reload one service

1. 列出所有slice

       `systemctl -t slice`

### 使用Systemctl控制并管理服务
1. 列出所有服务（包括启用的和禁用的）
       # systemctl list-unit-files --type=service

1. Linux中如何启动、重启、停止、重载服务以及检查服务（如 httpd.service）状态
```shell
# systemctl edit linstor-satellite # 不建议直接编辑systemd service文件, 它通过覆盖配置实现, 类似ini配置覆盖的形式
# systemctl start httpd.service
# systemctl restart httpd.service
# systemctl stop httpd.service
# systemctl reload httpd.service # 重新载入httpd配置而不中断服务
# systemctl condrestart httpd.service # condrestart会检查服务是否已运行, 如果已运行则重启; 否则忽略
# systemctl status [httpd.service]
# ls /etc/systemd/system/*.wants/httpd.service 查看httpd服务在各个运行级别下的启用和禁用情况
```
注意：当我们使用systemctl的start，restart，stop和reload命令时，我们不会从终端获取到任何输出内容，只有status命令可以打印输出。

1. 如何激活服务并在启动时启用或禁用服务（即系统启动时自动启动服务）
        # systemctl is-active httpd.service
        # systemctl enable --now httpd.service # `--now`即同时start service
        # systemctl disable httpd.service
1. 如何屏蔽（让它不能启动）或显示服务（如 httpd.service）
        # systemctl mask httpd.service
        # systemctl unmask httpd.service
1. 使用systemctl命令杀死服务
        # systemctl kill httpd
        # systemctl status httpd

### 使用Systemctl控制并管理挂载点
1. 列出所有系统挂载点
       # systemctl list-unit-files --type=mount

1. 挂载、卸载、重新挂载、重载系统挂载点并检查系统中挂载点状态
```shell
# systemctl start tmp.mount
# systemctl stop tmp.mount
# systemctl restart tmp.mount
# systemctl reload tmp.mount
# systemctl status tmp.mount
```
1. 在启动时激活、启用或禁用挂载点（系统启动时自动挂载）
```shell
# systemctl is-active tmp.mount
# systemctl enable tmp.mount
# systemctl disable  tmp.mount
```
1. 在Linux中屏蔽（让它不能启用）或可见挂载点
       # systemctl mask tmp.mount
       # systemctl unmask tmp.mount

### 使用Systemctl控制并管理套接口

1. 列出所有可用系统套接口
       # systemctl list-unit-files --type=socket
1. 在Linux中启动、重启、停止、重载套接口并检查其状态
```shell
# systemctl start cups.socket
# systemctl restart cups.socket
# systemctl stop cups.socket
# systemctl reload cups.socket
# systemctl status cups.socket
```
1. 在启动时激活套接口，并启用或禁用它（系统启动时自启动）
```shell
# systemctl is-active cups.socket
# systemctl enable cups.socket
# systemctl disable cups.socket
```
1. 屏蔽（使它不能启动）或显示套接口
       # systemctl mask cups.socket
       # systemctl unmask cups.socket

### 服务的CPU利用率（分配额）
1. 获取当前某个服务的CPU分配额（如httpd）
       # systemctl show -p CPUShares httpd.service
       CPUShares=1024
注意：各个服务的默认CPU分配份额=1024，你可以增加/减少某个进程的CPU分配份额。

1. 将某个服务（httpd.service）的CPU分配份额限制为2000 CPUShares/
       # systemctl set-property httpd.service CPUShares=2000
       # systemctl show -p CPUShares httpd.service
       CPUShares=2000
注意：当你为某个服务设置CPUShares，会自动创建一个以服务名命名的目录（如 httpd.service），里面包含了一个名为90-CPUShares.conf的文件，该文件含有CPUShare限制信息，你可以通过以下方式查看该文件：

       # vi /etc/systemd/system/httpd.service.d/90-CPUShares.conf

1. 检查某个服务的所有配置细节
       # systemctl show httpd

1. 分析某个服务（httpd）的关键链
       # systemd-analyze critical-chain httpd.service

1. 获取某个服务（httpd）的依赖性列表
       # systemctl list-dependencies httpd.service

1. 按等级列出控制组
       # systemd-cgls

1. 按CPU、内存、输入和输出列出控制组
       # systemd-cgtop

### 控制系统运行等级
1. 启动系统救援模式
       # systemctl rescue

1. 进入紧急模式
       # systemctl emergency

1. 列出当前使用的运行等级
       # systemctl get-default

1. 启动运行等级5，即图形模式
       # systemctl isolate runlevel5.target #或systemctl isolate graphical.target

1. 启动运行等级3，即多用户模式（命令行）
       # systemctl isolate runlevel3.target #或# systemctl isolate multiuser.target

1. 切换系统运行级别
       # systemctl isolate graphical.target
1. 设置多用户模式或图形模式为默认运行等级
       # systemctl set-default multi-user.target # runlevel3.target
       # systemctl set-default runlevel5.target
       Removed symlink /etc/systemd/system/default.target.
       Created symlink from /etc/systemd/system/default.target to /usr/lib/systemd/system/multi-user.target.
1. 重启、停止、挂起、休眠系统或使系统进入混合休眠
```shell
# systemctl reboot [--firmware-setup] # `--firmware-setup`与 reboot 一起使用时, 它会指示系统固件启动进入固件设置界面
# systemctl halt
# systemctl suspend # 待机/挂起
# systemctl hibernate # 休眠
# systemctl hybrid-sleep # 混合休眠
```

 对于不知运行等级为何物的人，参考`ls -al /lib/systemd/system/runlevel`, 说明如下:

 Runlevel 0 : 关闭系统(runlevel0.target, poweroff.target)
Runlevel 1 : 救援/维护模式(runlevel1.target, rescue.target)
Runlevel 3 : 多用户，无图形系统(runlevel3.target, multi-user.target)
Runlevel 2/4 : 多用户，无图形系统(runlevel2.target, runlevel4.target, multi-user.target)
Runlevel 5 : 多用户，图形化系统(runlevel5.target, graphical.target)
Runlevel 6 : 关闭并重启机器(runlevel6.target, reboot.target)

### 管理远程系统

1. systemctl命令通常都能被用来管理远程主机(用ssh通讯),只需将远程主机和用户名添加到systemctl命令后

       # systemctl status sshd -H root@1.2.3.4


### 时区设置
```sh
# timedatectl list-timezones # 查看当前所支持的时区信息
# timedatectl set-timezone ${zone} # 选择上述中的时区，然后设置
# timedatectl status # 查看当前时区设置的状态
```

### timer
计划的事件就是在特定时间需要激活的服务. systemd 管理一个名为定时器的工具，它类似 cron 的功能.

```sh
# systemctl list-timers
```

### sockets
`systemctl list-sockets`

## 参考

- [最简明扼要的 Systemd 教程，只需十分钟](https://linux.cn/article-6888-1.html)

## FAQ
### Failed at step USER spawning
启动service失败: `Failed at step USER spawning /home/chen/opt/xxx/yyy: No such process`, 已确认可执行程序存在且有执行权限.

解决方法:
```conf
[Service]
User=root #nobody
```

将`User=root #nobody`替换为`User=root`, 建议systemd配置文件的注释占用新行.