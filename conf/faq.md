# 常见问题

## 上网

### chrome离线下载
`https://www.google.cn/chrome/?standalone=1&platform=win64`, 不加platform时即下载当前os对应的版本.

platform:
- win64/win : win是windows 32位
- mac
- linux

> [测试版](https://www.google.cn/intl/zh-CN/chrome/beta/?hl=zh-CN&standalone=1), [开发者版](https://www.google.cn/intl/zh-CN/chrome/dev/?hl=zh-CN&standalone=1)

> [firefox离线包](https://www.mozilla.org/zh-CN/firefox/all)

> [Microsoft Edge 离线安装包](https://www.microsoft.com/zh-cn/edge/business-pages/download)

### apt安装google chrome
```sh
$ wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
$ sudo sh -c 'echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
$ sudo apt update
$ sudo apt-cache search chrome
$ sudo apt install google-chrome-stable
```

> 参考: https://www.ubuntuupdates.org/ppa/google_chrome?dist=stable

https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm
https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

### [apt安装firefox](https://www.omgubuntu.co.uk/2022/04/how-to-install-firefox-deb-apt-ubuntu-22-04)
ref:
- [Firefox Developer Edition and Beta: Try out Mozilla’s .deb package!](https://hacks.mozilla.org/2023/11/firefox-developer-edition-and-beta-try-out-mozillas-deb-package/)

```bash
sudo snap remove firefox
sudo add-apt-repository ppa:mozillateam/ppa

echo '
Package: *
Pin: release o=LP-PPA-mozillateam
Pin-Priority: 1001
' | sudo tee /etc/apt/preferences.d/mozilla-firefox # 改变Firefox软件包的优先级，以确保PPA/deb/apt版本的Firefox总是首选的

echo 'Unattended-Upgrade::Allowed-Origins:: "LP-PPA-mozillateam:${distro_codename}";' | sudo tee /etc/apt/apt.conf.d/51unattended-upgrades-firefox # 希望将来的Firefox升级自动安装

sudo apt update
sudo apt install firefox
```

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
3、将"ro recovery"改为"rw single init=/bin/bash". 有些系统没有bash可改为`init=sh`
4、按ctrl+X或者F10启动，进入root shell
5、执行"passwd 用户名"
6、修改完成后按ctrl + alt + del重启电脑
```

单用户默认启用网络: `systemctl start NetworkManager`

单用户环境:
- lvm:
    - lvm盘已识别, 见`/dev/mapper`, mount即可

## android清除密码
前提: 已root

1. 下载adb工具, 比如adb1.0.32.zip
2. 执行`adb shell`, 强制手机打开usb调试模式
3. 执行`rm password.key`, 再`reboot`

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

## linux登录后应用自启动
```sh
$ pwd
/home/chen/.config/autostart
$ cat Zoiper5.desktop
[Desktop Entry]
Encoding=UTF-8
Name=Zoiper5
Comment=VoIP Softphone
Exec=/home/chen/opt/Zoiper5/zoiper
Terminal=false
Icon=
Type=Application
$ cat alarm-clock-applet.desktop
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

[XDG Autostart 规范](https://specifications.freedesktop.org/autostart-spec/autostart-spec-latest.html):
针对所有用户: `/etc/xdg/autostart/*.desktop`
针对某用户: `/home/<user>/.config/autostart/*.desktop`

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

## linux 搜狗输入法 禁用半角切换
打开Fcitx配置界面 -> 全局配置, 选中左下角的`显示高级选项`, 重新定义`切换全角`的快捷键即可.

### 查找域名对应的ip
`http://IPAddress.com`

> [for github被屏蔽](https://zhuanlan.zhihu.com/p/65154116)

### linux 字符界面输入出现`]`等乱码
用`shift + backspace`来删除

### 强制命令输出是英文
`LANG="POSIX" ls -l`

### 进入单用户模式, `systemctl start NetworkManager`报`System has not been booted with systemd as init system (PID 1). Can't operate.`
将进入单用户模式的init参数换成`init=/sbin/init`即可.

### 字符界面下启用wifi
1. 启动NetworkManager: `systemctl start NetworkManager`
1. 执行nmtui, 选择第二项`Active a connection`
1. 选中wifi输入密码即可

### chrome使用代理
```bash
google-chrome-stable --proxy-server="socks5://127.0.0.1:20170" // 需梯子
```

### chrome无法同步, 登录账号后提示：`无法同步到“xxx@gmail.com” Request canceled`, 然后chrome账号退出
和 SwitchyOmega 有关, 解决方法:
1. 同步时直接将其直接切换到proxy模式
1. SwitchyOmega中添加规则`*.googleapis.com`(**推荐**)


### `do-release-upgrade -d`升级下载新版包时意外中断(即开始安装新版包前)后, 重新执行时报`请在升级前安装您的发行版所有可用更新`
将`/etc/apt/sources.list`(新版源)替换回`/etc/apt/sources.list.distUpgrade`(旧版源), 重新执行`do-release-upgrade`

### 用`do-release-upgrade`将ubuntukylin 20.04升级到22.04问题
1. 部分应用没有标题栏, 即使有标题栏的应用的其最大化按钮也失效; 点击窗口无法获取输入焦点; 点击窗口后无法前置

    在`apt dist-upgrade`时发现`kylin-nm ukui-control-center`没有升级(kylin-nm依赖ukui-control-center), `apt upgrade ukui-control-center`时提示缺`libkylin-chkname1`和`ukui-biometric-manager`, 其中`ukui-biometric-manager`没有, `libkylin-chkname1`安装成功, 重启后这些问题消失.
1. 无法设置背景

    log显示`ukui-control-ce`有段错误

根源: 在ubuntukylin 22.04未发布情况(即源还未更新, 部分依赖不完整甚至错误)下进行了升级, 比如[ukui-biometric-manager_1.0.3-1_amd64.deb](https://archive.ubuntukylin.com/ubuntukylin/pool/main/u/ukui-biometric-manager/ukui-biometric-manager_1.0.3-1_amd64.deb)是依赖libopencv4.2, 而实际上ubuntu 22.04的源里libopencv已是`4.5`

> 预计ubuntukylin 22.04在4.22发布, 升级日期是4.20

部分mirrors没有xxx-partner源, 比如清华源, aliyun mirror有jammy-partner但未与官方及时同步. 观察到华为源`https://mirrors.huaweicloud.com/ubuntukylin`比较及时.

### 安装搜狗linux输入法后只有繁体
修改"设置-常用-默认状态"也没用.

解决: 输入法工具栏右击再点击"简繁切换"即可.

### [ubuntu 22.04安装nvidia驱动](https://zh-cn.linuxcapable.com/%E5%A6%82%E4%BD%95%E5%9C%A8-ubuntu-22-04-lts-%E4%B8%8A%E5%AE%89%E8%A3%85-nvidia-%E9%A9%B1%E5%8A%A8%E7%A8%8B%E5%BA%8F/)
```bash
# ubuntu-drivers devices
# apt install nvidia-driver-515
# reboot
# nvidia-smi
```

或使用`软件更新器`->`设置`->`附加驱动`

注意: 第一代笔记本用`软件更新器`安装nvidia-driver-515后无法进入图形界面, ctrl+alt+f1切换到terminal, 再安装nvidia-driver-510重启后恢复正常.

### sogou linux不能输入中文(已在`拼音`输入模式), 英文正常
`rm -rf ~/.config/sogou*`, 再重启fcitx输入法即可, 推测是配置损坏.

### 恢复chrome的使用自定义证书的https拦截
打开对应的网站, 点击url左侧的"不安全"按钮, 点击"开启警告".

### linux rdp
`apt install rdesktop && rdesktop 192.168.xxx.xxx`

报错:
- `Failed to initialize NLA, do you have correct Kerberos TGT initialized`

    在系统属性-远程中开启：【允许远程协助连接这台计算机】+【允许远程连接到此计算机】，如果勾选了【仅运行运行使用网络级别身份验证的远程桌面单位计算机连接】, 那么 rdesktop 就报该错


其他软件:
- xfreerdp
- GNOME Connections

    `flatpak install flathub org.gnome.Connections`
- krdc

    `flatpak install flathub org.kde.krdc`
- Remmina

    `flatpak install flathub org.remmina.Remmina`

### 忘记密钥环密码, chrome密码管理器变空
1. 删除原先密钥环

    ```bash
    $ cd ~/.local/share/keyrings/
    $ cp login.keyring login.keyring.bak
    seahorse # 通过seahorse删除左侧"登录"项
    $ sudo reboot
    ```
1. 打开chrome, 登入google账号等待同步即可

### linux图形界面重启后黑屏
重启桌面显示管理器, 比如`systemctl restart  lightdm.service`

### 点向日葵图标无法启动向日葵
env: ubuntu 22.04

解决方法: 先启动向日葵服务`systemctl status runsunloginclient.service`, 再点击图标启动即可.

### todesk无法启动
env:ubuntu 24.04

`todeskd.service`报`[pulseaudio] backend-native.c: org.bluez.ProfileManager1.RegisterProfile() failed: org.bluez.Error.NotPermitted: UUID already registered`.

给`/usr/lib/systemd/user/pulseaudio.service`添加`--system`参数并`systemctl --user daemon-reload`+`systemctl --user restart pulseaudio.service`后错误消失, 但todesk依旧core dump.

[按照官方技术支持](https://uc.todesk.com/serviceSupport/workOrderDetail?id=79318), 查看`lscpu |grep -i avx2`, 当前cpu没有avx2, 让安装了低版本的[v4.1.0](https://dl.todesk.com/linux/todesk_4.1.0_amd64.deb), 能启动, 连接远程可能报错:
1. 连接卡在1%: 取消重试
1. `未知错误(12202)`: 12202是不支持低版本的主控和其他设备同时远程一个被控的提示

官方客服回复: 高版本需要avx2.

#### fedora 40
安装[todesk-v4.7.2.0-x86_64.rpm](https://dl.todesk.com/linux/todesk-v4.7.2.0-x86_64.rpm)后, 点击图标无法启动, 且在terminal启动报`/opt/todesk/bin/ToDesk: error while loading shared libraries: libappindicator3.so.1: cannot open shared object file: No such file or directory`.

解决: 参考[Fedora/CentOS/RedHat](https://www.todesk.com/linux.html)安装文档, 执行`dnf install libappindicator-gtk3`再启动即可. debian安装`libappindicator3-1`

### 停止gvfs
```bash
systemctl --user stop gvfs-daemon
systemctl --user mask gvfs-daemon
```

### 输入法切换
env: fedora 40 + kde + fcitx5, 自带的`输入法选择器`无法切换输入法

```bash
$ imsettings-info # 当前输入法信息
$ imsettings-list # 支持的输入法
$ imsettings-switch fcitx5 # 切换输入法
```

### firefox + fcitx5无法输入中文
env: fedora 40 + kde + fcitx5

使用chrome

### fedora 39/40禁止自动挂起
服务器用的fedora 40(使用了默认的gnome桌面)，ssh总是隔一段时间就自动挂起了, 并且提示：
```
Broadcast message from chen@fedora (Sun 2024-09-01 13:37:31 CST):

The system will suspend now!
```

解决方案
1. 推荐
```bash
vim /etc/systemd/logind.conf
IdleAction=ignore
```
1. 按电源键恢复
1. gnome里配置电源策略

### kde的baloo文件提取器总是崩溃
解决方法: 
- `sudo balooctl disable`(已测试)
- 清理文件(未测试)

    ```bash
    # rm -rf ~/.config/baloo*
    # rm -rf ~/.local/share/baloo/*
    ```

> `balooctl status`

### 主目录下自动生成乱码名称文件夹，内含/crashpad/helpers.bin
ref:
- [中文系统，主文件夹下名字乱码文件夹超来越来，求教如何解决](https://forum.ubuntu.org.cn/viewtopic.php?t=494933)

env:
- ubuntu 24.04

    google后发现debian 12也有该现象

嫌疑最大的软件: 微信LINUX, 钉钉LINUX

### kde的fcitx5-configtool 中显示`无法通过 DBus 连接到 Fcitx，Fcitx 是否正在运行？`
ref:
- [设置 Fcitx 5](https://fcitx-im.org/wiki/Setup_Fcitx_5/zh-cn#.2Fetc.2Fprofile)

调试输出见fcitx5-diagnose

### chrome登入google账号一会就退出, 终端报`[3809650:3809676:0925/140540.278349:ERROR:google_apis/gcm/engine/registration_request.cc:291] Registration response error message: DEPRECATED_ENDPOINT`
chrome已最新并重装, 问题依旧

将ZeroOmega代理切换到proxy, 再登入

### chrome不同步数据
进`chrome://sync-internals/`点击`Trigger GetUpdates`

解决:
1. 将ZeroOmega切换到proxy, 并重启chrome, **关键步骤**
2. 重新触发`Trigger GetUpdates`