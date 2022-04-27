# vnc
参考:
- [How to Install and Configure VNC Server on Ubuntu](https://www.tecmint.com/install-and-configure-vnc-server-on-ubuntu/)
- [Ubuntu 20.04 上安装和配置 VNC](https://xie.infoq.cn/article/cf473dc0dea917b0b2a546ecd)

## 安装

### tigervnc(**推荐**)
```bash
# xfce4-goodies 提供了一些额外的插件和一些有用的工具，如 mousepad 编辑器
$ sudo apt install tigervnc-standalone-server tigervnc-xorg-extension tigervnc-viewer [xfce4-goodies]
```

> tigervnc有`/etc/vnc.conf`.

### ~~tightvncserver(不推荐, 2009年后已不维护)~~
env: xubuntu 20.04

> tightvncserver没有`/etc/vnc.conf`

```bash
$ sudo apt install tightvncserver xfce4-goodies
$ vncpasswd
$ cat > ~/.vnc/xstartup << "EOF"
#!/usr/bin/env bash
unset SESSION_MANAGER
unset DBUS_SESSION_BUS_ADDRESS
# 使用fcitx输入环境
export GTK_IM_MODULE="fcitx"
export QT_IM_MODULE="fcitx"
export XMODIFIERS="@im=fcitx"
exec startxfce4 # kubuntu使用/usr/share/xsessions/plasma.desktop里的`startplasma-x11`
EOF
$ chmod +x ~/.vnc/xstartup
$ vncserver -geometry 1852x900
```

> 进入桌面后需要在该session打开的terminal上执行fcitx-autostart才能使用输入法.

### systemd
```bash
$ cat << EOF | sudo tee /lib/systemd/system/vncserver@.service # 名称末尾的@符号将允许传入一个可以在服务配置中使用的参数用来指定管理服务时要使用的VNC显示端口. 个人使用时直接在配置VNC显示端口为1即可
[Unit]
Description=Systemd VNC server startup script for Ubuntu 20.04
After=syslog.target network.target

[Service]
Type=forking
User=chen
Group=chen
WorkingDirectory=/home/chen

PIDFile=/home/chen/.vnc/%H:%i.pid
ExecStartPre=-/usr/bin/vncserver -kill :%i > /dev/null 2>&1 # 固定VNC显示端口为1: ExecStartPre=-/usr/bin/vncserver -kill :1 > /dev/null 2>&1
ExecStart=/usr/bin/vncserver -depth 24 -geometry 1852x900 :%i
ExecStop=/usr/bin/vncserver -kill :%i

[Install]
WantedBy=multi-user.target
EOF
$ sudo systemctl daemon-reload
$ sudo systemctl enable vncserver # 仅用于固定VNC显示端口的情况
$ sudo systemctl start vncserver
# --- 或使用参数
$ sudo systemctl enable vncserver@1
$ sudo systemctl start vncserver@1
$ sudo systemctl status vncserver@1
```

## 连接vnc
```bash
$ sudo apt install xtightvncviewer
$ xtightvncviewer 192.168.0.245:5901
```

## 其他vnc命令
```bash
$ vncserver -list # 以打开的vnc display
$ vncserver -kill :<N> # N可通过`ps -ef|grep "vnc"`中的进程参数获取
$ vncserver -kill :* # 关闭全部vnc display
```

## FAQ
### xfce4 高分辨率会糊
解决方案:
- 来回切换分辨率
- 等会自行恢复, **推荐**

### vnc viewer连接灰/黑屏幕
1. 检查`~/.vnc/xstartup`是否有执行权限

	对应的session exec即`/usr/share/xsessions/*.desktop`的`Exec`属性
1. 运行vncserver, 查看vncserver输出的log, 通常是dbus-launch出了问题

	比如~/.vnc/ubuntu-18:1.log提示"/usr/local/bin/dbus-launch: /lib/x86_64-linux-gnu/libdbus-1.so.3: version `LIBDBUS_PRIVATE_1.10.6' not found (required by /usr/local/bin/dbus-launch)".

	ubuntu本身自带/usr/bin/dbus-launch, 直接运行`/usr/bin/dbus-launch`会有内容输出. 但之前开发某个项目时自编译了dbus, 后来删除了该项目, 但vncserver优先使用了`/usr/local/bin/dbus-launch`导致报错, 删除或重命名`/usr/local/bin/dbus-launch`即可.

### vnc viewer无法连接
tigervnc-standalone-server启动时默认绑定localhost. 因为vncserver没有使用TLSVnc, 不安全启动时默认绑定到localhost.

解决方法:
- 在`/etc/vnc.conf`中追加`$localhost = "no";`(**`;`不能省略**), 重启系统再重新运行`vncserver`即可.
- vncserver使用参数`-localhost no`

### vncserver修改分辨率
解决方案:
- `/etc/vnc.conf`的配置项`$geometry`支持修改分辨率, 比如`$geometry = "1850x900";`
- vncserver使用参数`-geometry 1852x900`