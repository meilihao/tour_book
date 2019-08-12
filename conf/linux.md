### Tools/工具

#### 制作u盘启动盘(linux下)

1. Rufus, 推荐, 也支持windows
1 . [live-usb-install](http://sourceforge.net/projects/liveusbinstall/files/?source=navbar)
1 . [UNetbootin](http://sourceforge.net/projects/unetbootin/files/UNetbootin/),推荐

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

- 归档文件管理器解压zip中文乱码

	通过unzip命令解压，指定字符集`unzip -O CP936 xxx.zip` (用GBK, GB18030也可以)

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

### 安全

- 执行当前可执行程序加`./`的原因：

	主要是安全原因，因为在linux中执行程序时，会先搜索当前目录然后是系统目录，所以如果当前目录中有与系统可执行程序重名的程序，比如cp，她就会优先执行当前目录中的cp，但是如果当前目录的cp是木马，就会威胁到系统安全，所以这是Linux的一种安全策略，所以默认并没有把当前目录加到环境变量PATH中去.

### 优化

#### swappiness

在linux里面，swappiness的值的大小对如何使用swap分区是有着很大的联系的。swappiness=0的时候表示最大限度使用物理内存，然后才是 swap空间，swappiness＝100的时候表示积极的使用swap分区，并且把内存上的数据及时的搬运到swap空间里面。两个极端，对于ubuntu的默认设置，这个值等于60，建议修改为10。具体这样做：

1. 查看你的系统里面的swappiness

       $ cat /proc/sys/vm/swappiness

2. 修改swappiness值为10

       $ sudo sysctl vm.swappiness=10

 这只是临时性的修改，在你重启系统后会恢复默认的60，所以，还要做一步：

       $ sudo gedit /etc/sysctl.conf

 在这个文档的最后加上这样一行:

       vm.swappiness=10

改成10后感觉开关机,运行都慢了许多,调为50即可.

swap分区:
- `<=4G`: 内存的2倍
- `>4G&&<=16G`: 内存大小
- `>16G`: 不设置swap

> db server建议使用大内存且关闭swap
> 永久禁用swap: 修改`/etcd/fstab`
> 临时禁用swap: `swapoff -a`

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

# 编辑器
## Ubuntu修改默认编辑器
`sudo update-alternatives --config editor`

# 快捷键
1. `super + shift + tab`: 重启窗口管理器
1. `sudo systemctl restart lightdm`: 注销用户