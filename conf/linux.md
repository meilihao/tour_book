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
解决方法：临时关闭goagent服务；再找个VPN，通过VPN（http://www.i-vpn.net/free-vpn/）完成上传(以测试可行)

 在uploader.py里加入（未测试）：
import ssl
ssl._create_default_https_context = ssl._create_unverified_context
```

### so/动态库

Linux中的.so文件类似于Windows中的DLL，是动态链接库，也有人译作共享库（因so的全称为Shared Object）。当多个程序使用同一个动态链接库时，既能节约可执行文件的大小，也能减少运行时的内存占用.

对于用户而言，经常遇到的问题是某些应用程序找不到其需要的.so文件：`error while loading shared libraries: ...: cannot open shared object file: No such file or directory`


#### so存放位置

Linux中绝大多数.so文件都存放在/lib、/usr/lib/（见Linux目录结构），对于64位和32位共存的系统，32位的动态库可能会放在/lib32、/usr/lib32，完整的动态库存放路径列表可通过/etc/ld.so.conf文件配置。（如果修改了配置，需要用 /sbin/ldconfig 命令更新缓存）

**应注意动态库搜寻路径并不包括当前文件夹，所以当即使可执行文件和其所需的so文件在同一文件夹，也会出现找不到so的问题**,如:`./chrome: error while loading shared libraries: libnss3.so.1d: cannot open shared object file: No such file or directory`

此时可用LD_LIBRARY_PATH环境变量做临时设置，如：

```shell
# 将当前目录路径加入LD_LIBRARY_PATH
export LD_LIBRARY_PATH=.:$LD_LIBRARY_PATH
```

也有些so文件是在程序执行时临时加载的（如插件），它们的路径就比较灵活，只要可执行文件能找到它就行了。

#### 程序链接的动态库

Linux下查看so动态链接库的常用命令:
- ldd  查看可执行文件链接(依赖)了哪些系统动态链接库,**推荐**
- nm用来列出目标文件的符号清单.
- ar命令可以用来创建、修改库，也可以从库中提出单个模块。
- objdump：显示目标文件中的详细信息(`objdump -d <command>`，可以查看这些工具究竟如何完成这项任务)
- readelf 显示关于 ELF 目标文件的信息(`readelf -d libffmpeg.so | grep NEEDED`)

#### 版本

动态库的版本总是个问题，如果编译时链接的库和执行时提供的不一样，难免会导致程序的执行发生诡异的错误。为解决此问题，Linux系列的做法是这样的：

首先，每个so文件有一个文件名，如libABC.so.x.y.z，这里ABC是库名称，x.y.z是文件的版本号，一般来说：[2]

    第一位x表示了兼容性，x不一样的so文件是不能兼容的。
    第二位y的变化表示可能引入了新的特性（Feature），但总的来讲还是兼容的。
    第三位z的变化一般表示仅是修正了Bug。
    并非所有.so文件都遵循此规则，但其应用确实很普遍。

在系统中，会存在一些符号链接：

```
libpam.so -> libpam.so.0.83.0
libpam.so.0 -> libpam.so.0.83.0
```

其中第一个主要在使用该库开发其它程序时使用，比如gcc想连接PAM库，直接连libpam.so就行了，不必在链接时给出它的具体版本。第二个则主要用在运行时，因为前面说了第一位版本一样的库是互相兼容的，所以程序运行时只要试图连接libpam.so.0就够了，不必在意其具体的版本。ldconfig可以自动生成这些链接。

那么编译程序时gcc在链接一个so文件（如libpam.so）时，如何知道该程序运行时链接哪个文件呢（上例中是libpam.so.0）？原来产生so文件时，可以指定一个soname，一般形如libABC.so.x。有人编译可执行文件时，如果链接了某个so，存在可执行文件里的.so文件名并不是其全名，而是这个soname。比如上例中，这个soname就是libpam.so.0。回头看一下上节ldd的结果，可以印证这一现象。

有时还会看到形如libABCn.so，即版本号出现在.so前面的库文件，如libXaw6.so。此类文件一般是为开发者着想，比如GTK+ 3已经推出，但很多开发者还是想用GTK+ 2开发软件。由于编译时只连接无版本号的.so文件，就只有把版本号放在.so前面了。

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

#### I/O调度器

- [如何更改 Linux I/O 调度器来调整性能 for ssd](https://linux.cn/article-8179-1.html)

修改后需更新grub: `sudo update-grub`.

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