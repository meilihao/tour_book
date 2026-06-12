# tools

# 可视化
ref:
- [json](https://json4u.cn/editor)
- [数据与算法](https://staying.fun/)
- [Online Compiler, AI Tutor, and Visual Debugger for Python, Java, C, C++, and JavaScript](https://pythontutor.com/)

# 破解
ref:
- [hashcat 暴力破解详细教程](https://teamssix.com/210831-122611)

	步骤:
	1. 提前hash值
	2. 用hashcat破解

# convert

## pdf
ref:
- [OCRmyPDF : 将扫描的 PDF 文件转换为可搜索、可复制的文档](https://zhuanlan.zhihu.com/p/22188285184)

### xlsx/docx/pptx转pdf
- linux

	`libreoffice --invisible --language=zh-CN --convert-to pdf <in_path> --outdir <out_path>`
- windows

	`cmd /c soffice --headless --invisible --convert-to pdf <in_path> --outdir <out_path>`

# file
- [httpfs Go 编写的静态文件服务器](https://github.com/hellojukay/httpfs)
- [Surge：终端里的高速下载工具](https://github.com/surge-downloader/Surge)

# net
- [rathole 是一个安全、稳定、高性能的内网穿透工具，用 Rust 语言编写](https://github.com/rathole-org/rathole)
- [P2P STUN打洞: STUN Max(golang)](github.com/uk0/stun_max)

# download
## uget+aria2

[uget](http://ugetdm.com)是一款轻量级的自由开源的下载管理器类似迅雷，可运行Linux、windows和MAC系统上,支持队列下载和恢复下载和通过终端下载的功能.uget采用aria2作为后端，安装aria2插件后可与其进行交互。

`aria2` 是 Linux CLI 界面下的多线程下载工具，与 axel 类似，但比之更强大.它支持 HTTP/HTTPS, FTP, BitTorrent 和 Metalink 协议，支持多线程下断点续传.另外，这里有一个名为 aria2fe 的 aria2 前端 GUI 程序，直接执行里面编译好的二进制程序就可使用.

### 以centos 7 x64为例

1. 下载Fedora分类下的最新版[uget](http://sourceforge.net/projects/urlget/files/uget%20%28stable%29/1.10.4/uget-1.10.4-1.fc21.x86_64.rpm/download)和[aria2](http://sourceforge.net/projects/urlget/files/aria2-plugin/1.18.x/aria2-1.18.2-1.fc21.x86_64.rpm/download)
2. 安装uget和aria2
3. 配置
	```shell
    1. 在uget的`编辑`-`设置`-`插件`选项卡中勾选`启用 aria2 插件`
    2. uget主界面左侧的`分类`窗口-`Home`分类-右键`属性`-`新下载的默认设置1`选项卡-`每台服务器连接数`改为'16'(其它数字也可).
    ```

# 注意力
- [开源应用 Pomatez 使你保持专注](https://linux.cn/article-16268-1.html)
- [Ubuntu上安装番茄时钟](https://zhuanlan.zhihu.com/p/350023097)

    `pip3 install tomato-clock`
- [github.com/coolcode/tomato-clock-rs](https://github.com/coolcode/tomato-clock-rs)

	上面tomato-clock的rust实现
- [alarm clock](https://alarm-clock-applet.github.io/)

# 启动盘
1. [Ventoy](https://www.ventoy.net/cn/doc_start.html), 推荐

	```bash
	# cd ventoy-1.0.64
	# ./VentoyGUI.x86_64 [--qt5] # 推荐使用qt
	```

	操作步骤:
	1. 通过菜单`配置选项-清除Ventoy`清理U盘
	1. `配置选项-分区类型`选择GPT. 不能选择`安全启动支持`, 旧主板或部分主板不支持
	1. 点`安装`即可
	1. 用rsync拷贝iso(比cp快), 比如`rsync --progress -v ubuntukylin-20.04-pro-sp1-amd64.iso /media/chen/Ventoy`
	1. 执行`sync`并计算iso checksum

        可用`iostat -d -m -x 5`监控磁盘写入, 避免sync很长时还以为是sync卡住
	1. reboot并开始安装os

1. PXE(Preboot eXecute Environment,预启动执行环境)

	PXE可以让计算机通过网络来启动操作系统(前提是计算机上安装的网卡支持 PXE 技术),主要用于在无人机值守安装系统中引导客户端主机安装 Linux 操作系统.
	Kickstart 是一种无人值守的安装方式,其工作原理是预先把原本需要运维人员手工填写的参数保存成一个ks.cfg 文件,当安装过程中需要填写参数时则自动匹配 Kickstart 生成的文件即可.

	> system-config-kickstart 是一款图形化的 Kickstart 应答文件生成工具,可以根据自己的需求生成自定义的应答文件.

	其他批量安装系统工具:
	- [Cobbler](https://mp.weixin.qq.com/s/cFXIF5QqS_hewzTeHHuFhA)

		Cobbler 是基于 Kickstart 的更高级封装工具，集成了DHCP、TFTP、DNS和PXE环境

## 上网
### chrome离线下载
`https://www.google.cn/chrome/?standalone=1&platform=win64`, 不加platform时即下载当前os对应的版本.

platform:
- win64/win : win是windows 32位
- mac
- linux

> [测试版](https://www.google.cn/intl/zh-CN/chrome/beta/?hl=zh-CN&standalone=1), [开发者版](https://www.google.cn/intl/zh-CN/chrome/dev/?hl=zh-CN&standalone=1)

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

### [安装firefox](https://www.omgubuntu.co.uk/2022/04/how-to-install-firefox-deb-apt-ubuntu-22-04)
推荐: [firefox离线包](https://www.firefox.com/zh-CN/) + [从 Mozilla 安装 Firefox](https://support.mozilla.org/zh-CN/kb/install-firefox-linux)

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

## glxinfo
```
$ sudo apt-get install mesa-utils
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

### 恢复chrome的使用自定义证书的https拦截
打开对应的网站, 点击url左侧的"不安全"按钮, 点击"开启警告".

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

## FAQ
### qqmusic崩溃报`FATAL:gpu_data_manager_impl_private.cc(1034)] The display compositor is frequently crashing. Goodbye.`
追加`--no-sandbox`

### chrome登入google账号一会就退出, 终端报`[3809650:3809676:0925/140540.278349:ERROR:google_apis/gcm/engine/registration_request.cc:291] Registration response error message: DEPRECATED_ENDPOINT`
chrome已最新并重装, 问题依旧

将ZeroOmega代理切换到proxy, 再登入

### chrome不同步数据
进`chrome://sync-internals/`点击`Trigger GetUpdates`

解决:
1. 将ZeroOmega切换到proxy, 并重启chrome, **关键步骤**
2. 重新触发`Trigger GetUpdates`

遇到`google-chrome-stable --proxy-server="socks5://127.0.0.1:20170"`+上述步骤还是无法同步, 换成`google-chrome-stable --proxy-server="127.0.0.1:20171"`后恢复

### chrome使用代理
```bash
google-chrome-stable --proxy-server="socks5://127.0.0.1:20170" // 需梯子
```

### chrome无法同步, 登录账号后提示：`无法同步到“xxx@gmail.com” Request canceled`, 然后chrome账号退出
和 SwitchyOmega 有关, 解决方法:
1. 同步时直接将其切换到proxy模式
1. SwitchyOmega中添加规则`*.googleapis.com`(**推荐**)

### 钉钉无法登入
ref:
- [在 openSUSE-Leap-15.5-DVD-x86_64 中使用钉钉 dingtalk_7.0.40.30829_amd64](https://forum.suse.org.cn/t/topic/16484)
- [在 Linux "rpm" 系发行版上运行钉钉应用程序](https://tedding.dev/2023/06/16/188c2ae4970.html)

env:
- kubuntu 24.04
- com.alibabainc.dingtalk_7.5.20.40605_amd64.deb

报错原因:
- 找不到`xxx.so`
- `libpango undefined symbol: hb_ot_metrics_get_position`

patch:
```bash
# diff Elevator.sh Elevator.sh.bak
2d1
< export LD_LIBRARY_PATH=/snap/gnome-42-2204/176/usr/lib/x86_64-linux-gnu:/snap/gnome-42-2204/176/usr/lib/x86_64-linux-gnu/pulseaudio:$LD_LIBRARY_PATH
```