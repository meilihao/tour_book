# desktop
安装lxde(基于GTK2):
```sh
$ sudo apt install lxde lxde-common
```

安装xfce4:
```sh
$ sudo apt install xfce4/xubuntu-desktop # xubuntu-desktop需要更多空间, 也更接近XUbuntu
```

安装Ubuntukylin:
```sh
$ sudo apt install ubuntukylin-desktop
```

安装gnome:
```bash
apt install ubuntu-gnome-desktop
```

安装mate:
```bash
apt install ubuntu-mate-desktop
```

安装lxqt:
```bash
# apt install lxqt # 精简版
# apt install task-lxqt-desktop # 完整版
```

卸载lxqt:
```sh
$ sudo apt purge lxqt-* liblxqt*
```

查看系统当前支持的桌面环境 : `/usr/share/xsessions`
当前支持的X11 session在`/usr/share/xesssions`.
当前支持的Wayland session在`/usr/share/wayland-sessions`.
当前正在使用的显示管理器: `cat /etc/X11/default-display-manager`.
设置默认显示管理器: `sudo dpkg-reconfigure <sddm/gdm3/lightdm>`
进入桌面自动打开应用: `~/.config/autostart/*.desktop`

## 屏保xscreensaver

屏保位置为`/usr/lib/xscreensaver`.
启动屏保`xscreensaver-command -select 116`,屏保的order在`~/.xscreensaver`里,这里的`116`是`flurry`.

## terminix,推荐
- [terminix](https://github.com/gnunn1/terminix), 推荐使用配色方案`Monokai Dark`
- [Tilix](https://gnunn1.github.io/tilix-web/),　支持同步输入,保存布局

## FAQ
### ubuntu输入正确用户密码登陆时重新跳转到登陆界面即无法登陆

原因：用户home目录下的.Xauthority文件拥有者变成了root，从而以用户登陆的时候无法都取.Xauthority文件.

解决：删除home目录下的.Xauthority文件，再重启(chown修改文件属性不可行).
