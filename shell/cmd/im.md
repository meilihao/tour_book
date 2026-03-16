# im
ref:
- [Linux 输入法](https://mp.weixin.qq.com/s/dhY7tA83OW17Jmlcy4L-HA)

桌面系统使用im-config包中的im-config命令配置默认输入法框架. 启动时通过/etc/X11/Xsession.d/70im-config_launch脚本读取下面两个文件（定义在 /usr/share/im-config/xinputrc.common中）的输入法变量来决定启动哪个输入法框架服务的，但是启动后，会被fcitx.desktop自启覆盖:
- $HOME/.xinputrc
- /etc/X11/xinit/xinputrc

服务器主要是使用imsettings包中 imsettings-switch命令行工具配置默认输入法框架. 当然也有图形化的im-chooser，需要手动安装. 启动时通过imsettings-start.desktop在自启动时根据下面两个文件的输入法变量实现来决定启动哪个输入法框架服务的:
- $HOME/.config/imsettings/xinputrc
- /etc/X11/xinit/xinputrc

## example
```
# 检查im-config配置（某些系统使用）
$ im-config -l # 列出当前配置
$ im-config -m # 显示当前激活的框架
# 手动切换框架，保存到 /etc/X11/xinit/xinputrc
$ im-config # 设置框架（会弹出图形界面）
$ fcitx-remote
$ fcitx5-configtool # 配置后如果没问题即可使用fcitx5
$ fcitx5-diagnose
```

fcitx配置:
```bash
$ vim ~/.bashrc
export XMODIFIERS=@im=fcitx # 告诉X11应用，通过XIM协议去找名为"fcitx"的服务器
export GTK_IM_MODULE=fcitx # 告诉GTK应用，使用内部的fcitx模块去连接
export QT_IM_MODULE=fcitx # 告诉Qt应用，使用内部的fcitx模块去连接
export SDL_IM_MODULE=fcitx # SDL2
export XIM=fcitx # SDL1
```

逻辑：应用 -> GTK/QT的 fcitx 模块 -> D-Bus -> Fcitx 输入法服务