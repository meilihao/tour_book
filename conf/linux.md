### Tools/工具

#### 制作u盘启动盘(linux下)

1 . [live-usb-install](http://sourceforge.net/projects/liveusbinstall/files/?source=navbar)
2 . [UNetbootin](http://sourceforge.net/projects/unetbootin/files/UNetbootin/),推荐

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

#### cpu

	#cpu信息
	cat /proc/cpuinfo
    #cpu当前运行模式:32/64
    getconf LONG_BIT

#### 内存

	cat /proc/meminfo

#### linux发行版本

	cat /etc/issue

#### 获取有关硬件方面的信息

**DMI (Desktop Management Interface)**就是帮助收集电脑系统信息的管理系统，DMI信息的收集必须在严格遵照SMBIOS规范的前提下进行。

[Dmidecode命令详解](http://www.ha97.com/4120.html)

### FAQ

- fedora 22 dnf update时报`google-chrome-stable-xxx.i386.rpm 的公钥没有安装`

	```shell
$ cd /etc/yum.repos.d
$ sudo vim google-chrome.repo
# 在文件末尾追加“gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub”，再重新更新即可
# 或者设置gpgcheck=0，关闭检查也可(不推荐)
```

- fedora 22 运行goagent报`请安装python-vte`

	```shell
$ sudo dnf install vte
```

- fedora22下运行goagent（python goagent-gtk.py）出错`ImportError: No module named OpenSSL`

	```shell
$ sudo dnf install pyOpenSSL
# 安装再运行后即可在桌面的左下角的任务栏里看到goagent运行的图标
```

- fedora22下运行goagent（python goagent-gtk.py）出错`Load Crypto.Cipher.ARC4 Failed, Use Pure Python Instead`

	```shell
$ sudo dnf install pycrypto
```

- fedora22下运行python uploader.py出错`urllib2.URLError: <urlopen error [SSL: CERTIFICATE_VERIFY_FAILED] certificate verify failed (_ssl.c:581)>`

	```shell
Goagent上传是可以使用自身goagent服务或VPN，因为上传过程也是需要”梯子“才可以。默认开启了goagent服务，且现有GAE工作不正常的时候就会报这个错。
解决方法：临时关闭goagent服务；再找个VPN，通过VPN（http://www.i-vpn.net/free-vpn/）完成上传
```