# fcitx
ref:
- [Linux中文输入法安装指南](https://mp.weixin.qq.com/s/eaaxp0J4VxAxqT-xrTQA3Q)

```
# apt install fcitx5 fcitx5-chinese-addons fcitx5-configtool
```

配置:
- kde: 系统设置 -> 自动启动 -> `添加`-`应用程序`, 选`工具`-`Fcitix 5`

## FAQ
### kde的fcitx5-configtool 中显示`无法通过 DBus 连接到 Fcitx，Fcitx 是否正在运行？`
ref:
- [设置 Fcitx 5](https://fcitx-im.org/wiki/Setup_Fcitx_5/zh-cn#.2Fetc.2Fprofile)

调试输出见fcitx5-diagnose

### firefox + fcitx5无法输入中文
env: fedora 40 + kde + fcitx5

使用chrome

### qoder无法输入中文
ref:
- [Using Fcitx 5 on Wayland](https://fcitx-im.org/wiki/Using_Fcitx_5_on_Wayland#KDE_Plasma)
- [fcix5输入法在部分软件中无法输入中文](https://forum.archlinuxcn.org/t/topic/14626)

找到启动菜单中的qoder快捷图标, 右键选择`编辑应用程序`, 命令行参数从`%F`改为`--enable-features=UseOzonePlatform --ozone-platform=wayland --enable-wayland-ime %F`

### 微信wechat无法输入中文
同上编辑微信快捷图标, 环境变量加`QT_IM_MODULE=fcitx XMODIFIERS=@im=fcitx`

### wps office无法输入中文
同上编辑`WPS Office`快捷图标, 环境变量加`QT_IM_MODULE=fcitx XMODIFIERS=@im=fcitx`