# 常见问题

## 上网

### fonts.googleapis.com被屏蔽导致网站加载变慢

Google的字体(fonts.googleapis.com)服务被屏蔽，导致很多网站打开都极慢.

```shell
# 通过修改hosts文件解决,以linux为例
# 编辑/etc/hosts
# 方法1: 将谷歌字体服务的链接替换成[科大LUG](https://lug.ustc.edu.cn/wiki/mirrors/help/revproxy)
fonts.googleapis.com         fonts.lug.ustc.edu.cn
ajax.googleapis.com          ajax.lug.ustc.edu.cn
themes.googleusercontent.com google-themes.lug.ustc.edu.cn
storage.googleapis.com       storage-googleapis.lug.ustc.edu.cn
fonts.gstatic.com            fonts-gstatic.lug.ustc.edu.cn
gerrit.googlesource.com      gerrit-googlesource.lug.ustc.edu.cn
secure.gravatar.com          gravatar.lug.ustc.edu.cn
# 方法2: 直接屏蔽,缺点是看不到Google字体的真正效果
127.0.0.1       fonts.googleapis.com
```

类似:
- [ReplaceGoogleCDN](https://github.com/justjavac/ReplaceGoogleCDN)

## linux 忘记密码
- dedpin 15.4.1
```
1、首先开机选择"Advanced options for *****"这一行按回车
2、然后选中最后是"（recovery mode）"这一行按"E"进入编辑页面
3、将"ro recovery"改为"rw single init=/bin/bash"
4、按ctrl+X或者F10启动，进入root shell
5、执行"passwd 用户名"
6、修改完成后按ctrl + alt + del重启电脑
```

## fish添加环境变量
```sh
$ vim .config/fish/conf.d/golang.fish
```
添加:
```text
set -x GOROOT /opt/go
set -x GOPATH /home/xjm/git/go
set -x PATH {$PATH} {$GOROOT}/bin {$GOPATH}/bin
```

## apt安装google chrome
```sh
$ wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
$ sudo sh -c 'echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
$ sudo apt update
$ sudo apt-cache search chrome
$ sudo apt install google-chrome-stable
```

> 参考: https://www.ubuntuupdates.org/ppa/google_chrome?dist=stable

## linux登录后应用自启动
```sh
$ ~/.c/autostart pwd
/home/chen/.config/autostart
$ ~/.c/autostart cat Zoiper5.desktop
[Desktop Entry]
Encoding=UTF-8
Name=Zoiper5
Comment=VoIP Softphone
Exec=/home/chen/opt/Zoiper5/zoiper
Terminal=false
Icon=
Type=Application
$ ~/.c/autostart cat alarm-clock-applet.desktop
[Desktop Entry]
Name=Alarm Clock
Name[zh_CN]=闹钟
Comment=Wake up in the morning
Comment[zh_CN]=早晨唤醒
Icon=alarm-clock
Exec=alarm-clock-applet --hidden
Terminal=false
Type=Application
Categories=GNOME;GTK;Utility;
X-Ubuntu-Gettext-Domain=alarm-clock
```