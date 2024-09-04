### Tools/工具

#### 制作u盘启动盘(linux下)

1. [Ventoy](https://www.oschina.net/p/ventoy)

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
1. Rufus, 推荐, 也支持windows
1 . [live-usb-install](http://sourceforge.net/projects/liveusbinstall/files/?source=navbar)
1 . [UNetbootin](https://unetbootin.org/)

#### 解压缩

- Fedora22归档管理器解压rar报错`Parsing filters is unsupported`

	```shell
wget http://www.rarsoft.com/rar/rarlinux-x64-5.2.1.tar.gz
tar zxvf rarlinux-x64-5.2.1.tar.gz
cd rar
make
make install
# 如果在32位系统上安装了64位的软件，运行unrar命令的时候，会提示如下错误:
# bash: /usr/local/bin/rar: cannot execute binary file
```

#### 开机启动项管理

- sysv-rc-conf : for ubuntu
- chkconfig : for redhat,centos,fedora

### 查询

#### 内核模块参数

```
#参数在parameters目录下
/sys/module/${modulename}/parameters/
```
#### 常见Linux的rpm和deb软件包查找地址

用于查找rpm包：
http://rpm.pbone.net/
[rpmfind.net](http://rpmfind.net/),推荐

Fedora的koji下的rpm包：
http://koji.fedoraproject.org/koji/

查找Ubuntu的deb包地址：
http://packages.ubuntu.com/,推荐

查找Debian的deb包地址：
https://www.debian.org/distrib/packages

#### 公钥认证

```shell
# 导入公钥
rpm --import /path/to/key_file
# 显示所有已经导入的gpg格式的公钥
rpm -qa gpg-*
# 显示某个密钥的详细信息
rpm -qi gpg-pubkey-NAME
# 检查包的公钥认证：安装过程中会自动执行
rpm -K /path/to/package_file
```

### 专业名词

- `fedora rawhide` ： 简单的说，Rawhide就是Fedora的滚动更新版，但这与 Gentoo、ArchLinux 等又不同，因为这个分支指向的是当前开发版（如同 FreeBSD 的 CURRENT 分支），所以极其不稳定。需要注意的是，这不是测试版，而是开发版.

### so

- `libwbclient.so.0`

```shell
# Fedora22,运行“系统设置”程序
$ gnome-control-center --overview
gnome-control-center: error while loading shared libraries: libwbclient.so.0: cannot open shared object file: No such file or directory
# dnf update后出现这个error，应该是/usr/lib/libwbclient.so.0不存在引起的.
# sudo dnf install libwbclient时却提示已安装，通过find命令找到其已存在于/usr/lib/samba/wbclient/libwbclient.so.0.12,作软连接即可解决。
$ sudo ln -s  /usr/lib/samba/wbclient/libwbclient.so.0.12 /usr/lib/libwbclient.so.0
```

### 美化

#### 安装雅黑字体

```shell
# 从网上下载字体或其他windows系统上获取
# 用"字体查看器"打开并安装字体,此时默认安装位置是'/home/$USER/.local/share/fonts/'
# 用工具"Ubuntu Tweak"-"调整"-"字体"中修改相关默认字体即可.
```

### 虚拟机

#### visualbox

- Ubuntu 14.04 `不能为虚拟电脑 Ubuntu 打开一个新任务.VT-x is disabled in the BIOS. (VERR_VMX_MSR_VMXON_DISABLED).`

> 解决方法：
> 1、移除当前不能打开的虚拟镜像
> 2、新建虚拟电脑,选择使用已有的虚拟硬盘文件(即刚才移除的)

### 硬件信息

#### 获取有关硬件方面的信息

**DMI (Desktop Management Interface)**就是帮助收集电脑系统信息的管理系统，DMI信息的收集必须在严格遵照SMBIOS规范的前提下进行。

[Dmidecode命令详解](http://www.ha97.com/4120.html)

### FAQ

### 优化
#### tcp bbr
要求 : linux kernel >=4.9

```
# /etc/sysctl.conf添加
net.ipv4.tcp_congestion_control=bbr
net.core.default_qdisc=fq
```

```
# 使新配置生效
sudo sysctl -p

# 结果都有 bbr , 则证明你的内核已开启bbr
sudo sysctl net.ipv4.tcp_available_congestion_control
sudo sysctl net.ipv4.tcp_congestion_control

# 看到有 tcp_bbr 模块即说明bbr已启动
lsmod | grep bbr
```

#### 文件监控限制
[fs.inotify.max_user_watches=524288](https://code.visualstudio.com/docs/setup/linux#_visual-studio-code-is-unable-to-watch-for-file-changes-in-this-large-workspace-error-enospc)

### FAQ

#### dmesg中出现`tpm_tis 00:05: A TPM error (6) occurred attempting to read a pcr value`

可能会导致开机时卡住几秒钟.

tpm资料,[TPM安全芯片](http://baike.baidu.com/view/687208.htm),[Trusted Platform Module](https://wiki.archlinux.org/index.php/Trusted_Platform_Module)

在BIOS里把TPM芯片禁用(推荐,thinkpad t430:F1->Security->Security Chip);或者禁用tpm系统模块(`echo blacklist tpm_tis > /etc/modprobe.d/tpm_tis.conf`),经测试,dmesg中还是有该错误信息.

#### `[drm:cpt_set_fifo_underrun_reporting] *ERROR* uncleared pch fifo underrun on pch transcoder A [drm:cpt_serr_int_handler] *ERROR* PCH transcoder A FIFO underrun`

导致开机卡住几秒钟,drm错误,待官方解决.

# 软件包管理工具

## dnf

- [27 个 Linux 下软件包管理工具 DNF 命令例子](https://linux.cn/article-5718-1.html)


# 快捷键
1. `super + shift + tab`: 重启窗口管理器
1. `sudo systemctl restart lightdm`: 注销用户

# deepin
## 升级内核
`sudo apt install linux-image-deepin-stable-amd64 linux-headers-deepin-stable-amd64`

### 启动时显示grub
```bash
$ sudo vim /etc/default/grub
GRUB_TIMEOUT=10
GRUB_TIMEOUT_STYLE=menu
...
$ sudo grub2-mkconfig -o "$(readlink -e /etc/grub2.cfg)" # for fedora/centos
$ sudo update-grub # for debian/ubuntu
```

### zram
ref:
- [Enable Zram on Linux For Better System Performance](https://fosspost.org/enable-zram-on-linux-better-system-performance)

debian:
```bash
$ sudo apt install systemd-zram-generator zram-tools # 推荐使用zram-tools(/etc/default/zramswap), 因为它可以设置为开机自启, 而systemd-zram-generator的systemd-zram-setup不可以
sudo vim /etc/systemd/zram-generator.conf
[zram0]
compression-algorithm = zstd
zram-size = min(ram/2,4096)
swap-priority = 100
$ sudo systemctl daemon-reload
$ sudo systemctl start systemd-zram-setup@zram0
```

opensuse:
```bash
$ sudo zypper install systemd-zram-service zram-generator
$ sudo vim /etc/systemd/zram-generator.conf
[zram0]
compression-algorithm = zstd
zram-size = min(ram/2,4096)
swap-priority = 100
$ sudo systemctl daemon-reload
$ sudo systemctl start systemd-zram-setup@zram0
$ sudo zramctl
$ sudo swapon --show
```

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

### kubuntu 24.04 安装com.tencent.wechat(**推荐**)
使用[铜豌豆](https://www.atzlinux.com/allpackages.htm)打包的com.tencent.wechat

```bash
apt -y install wget
wget -c -O atzlinux-v12-archive-keyring_lastest_all.deb https://www.atzlinux.com/atzlinux/pool/main/a/atzlinux-archive-keyring/atzlinux-v12-archive-keyring_lastest_all.deb
apt -y install ./atzlinux-v12-archive-keyring_lastest_all.deb
apt update
apt install com.tencent.wechat
```

**ps**: 需要同时安装[优麒麟的weinxin](https://www.ubuntukylin.com/applications/106-cn.html), 否则com.tencent.wechat能启动, 但手机扫码登入失败.

之前正常, 这周(24.9.2)开始突然无法扫描登入了, 重装com.tencent.wechat无效, 但重装优麒麟的weinxin后, com.tencent.wechat恢复正常.

### kubuntu 24.04 安装wechat from opencloudos(**还不可用**)
ref:
- [OpenCloudOS 支持 Linux 原生版微信，开启生态新篇章](https://www.cnblogs.com/OpenCloudOS/p/18252948)

	[rpm](https://mirrors.opencloudos.tech/opencloudos/9.2/extras/x86_64/os/Packages/wechat-beta_1.0.0.242_amd64.rpm)

将wechat-beta_1.0.0.242_amd64.rpm解压, 拷贝到相应目录, 根据启动报错修改wechat.desktop:
```bash
$ sudo vim /usr/share/applications/wechat.desktop
...
Exec=env LD_LIBRARY_PATH=/opt/wechat-beta:$LD_LIBRARY_PATH /usr/bin/wechat %U
...
```

修改后能打开wechat, 出现扫二维码登入界面, 此时点关闭会crash.

### 优麒麟weixin无法启动(**版本旧**)
根据`journalctl -f`的提示, 修改chrome-sandbox权限:
```bash
$ cd /opt/weixin
$ sudo chown root:root chrome-sandbox
$ sudo chmod 4755 chrome-sandbox
```