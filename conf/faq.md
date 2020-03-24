# 常见问题

## 上网

### 迅雷99永远下不完
1. 把任务删除，但是不要删除本地文件
2. 点击右上角的小箭头--文件--选择导入未完成下载 把原先那个.td后缀的文件导入即可（或者可以直接把这个文件拖入迅雷）,选择`继续下载`即可, 稍后就能下载完成.

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

## linux 忘记/重置密码
- dedpin 15.4.1
```
1、首先开机选择"Advanced options for *****"这一行按回车
2、然后选中最后是"（recovery mode）"这一行按"E"进入编辑页面
3、将"ro recovery"改为"rw single init=/bin/bash"
4、按ctrl+X或者F10启动，进入root shell
5、执行"passwd 用户名"
6、修改完成后按ctrl + alt + del重启电脑
```

## [fish添加环境变量](https://github.com/fish-shell/fish-shell/issues/527)
将相应的fish配置文件放入`/home/chen/.config/fish/conf.d`即可, fish的`conf.d`类似于nginx的`conf.d`.

```sh
$ vim .config/fish/conf.d/golang.fish
```
添加:
```text
set -x GOROOT /opt/go
set -x GOPATH /home/chen/git/go
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

## [deepin apt 系统更新异常以及高版本软件降级保护](https://wiki.deepin.org/wiki/%E7%B3%BB%E7%BB%9F%E6%9B%B4%E6%96%B0%E5%BC%82%E5%B8%B8%E4%BB%A5%E5%8F%8A%E9%AB%98%E7%89%88%E6%9C%AC%E8%BD%AF%E4%BB%B6%E9%99%8D%E7%BA%A7%E4%BF%9D%E6%8A%A4)
调整位置: /etc/apt/preferences.d/deepin 文件.

如果想修改此方案，可以在同级目录下按照deepin文件格式编辑其他第三方源优先策略. Pin-Priority 值越大，优先级越高;如果不想使用此方案，可以删除/etc/apt/preferences.d/deepin 文件.

## ssh config 添加别名
```
Host speak scriptRepo
    HostName 192.168.11.80
    Port 22
    User root
    IdentityFile    /home/chen/.ssh/my_rsa
```
`scriptRepo`就是别名

## glxinfo
```
$ sudo apt-get install mesa-utils
```

## X11/Xlib.h not found
```
$ sudo apt install libx11-dev
```

## cmake用法
```
$ cd ${project} # 该目录需包含`CMakeLists.txt`
$ mkdir build && cd build
$ cmake ..
$ make -j
```
## 添加sudo
```sh
# visudo # 添加`%sudo	ALL=(ALL:ALL) ALL`, 即属于sudo用户组的用户均可使用sudo命令
```

`tidb ALL=(ALL) NOPASSWD: ALL`表示使用sudo命令时无需密码.

## proxy curl https 卡住
```bash
$ env https_proxy="https://127.0.0.1:1081" curl https://www.google.com -v # 卡住, 原因未知
$ env https_proxy="http://127.0.0.1:1081" curl https://www.google.com -v # ok
$ env https_proxy=127.0.0.1:1081 curl https://www.google.com -v # ok, 推荐
```

## 显示器分辨率被固定
```bash
$ xrandr # 获取 `xxx connected`项的名称`xxx`, 其实就是`xrandr --listmonitors`的列表项,我这里是`DP-1`
$ cvt 1920 1440 60 # 生成配置参数
$ xrandr --newmode "1920x1440_60.00"  233.50  1920 2064 2264 2608  1440 1443 1447 1493 -hsync +vsync
$ xrandr --addmode DP-1 "1920x1440_60.00"
$ xrandr --output DP-1 --mode "1920x1440_60.00"
```

## su root报错`无法执行 fish: 没有那个文件或目录`
```sh
$ cat /etc/passwd|grep root
root:x:0:0:root:/root:fish # root的shell的路径不对
```

## ssh
### ssh-add无法添加ed25519 key
```
$ ssh-add ~/.ssh/my_ed25519
Enter passphrase for /home/chen/.ssh/my_ed25519:
Bad passphrase, try again for /home/chen/.ssh/my_ed25519:
Could not add identity "/home/chen/.ssh/my_ed25519": communication with agent failed

$ ssh -V
OpenSSH_7.5p1 Debian-5, OpenSSL 1.0.2l  25 May 2017
$ echo $SSH_AUTH_SOCK
/run/user/1000/keyring/ssh
```

虽然ssh-add无法添加, 但`ssh xxx`还是可正常使用

### ssh自动断开
1. 服务端设置
    找到/etc/ssh/sshd_config, 并修改数值

    # 30表示30s给客户端发送一次心跳
    ClientAliveInterval 30
    # 此客户端没有返回心跳3次，则会断开连接
    ClientAliveCountMax 3
    # TCP保持连接不断开
    TCPKeepAlive yes
1. 客户端配置(推荐)
    如果是想让主机所有用户都生效，修改/etc/ssh/ssh_config
    如果只想让当前用户生效，则修改 ~/.ssh/config

    Host *
        ServerAliveInterval 30
        ServerAliveCountMax 3
        TCPKeepAlive yes

## linux 搜狗输入法 禁用半角切换
打开Fcitx配置界面 -> 全局配置, 选中左下角的`显示高级选项`, 重新定义`切换全角`的快捷键即可.