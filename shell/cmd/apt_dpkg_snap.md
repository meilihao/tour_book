# apt

## 描述

debian,ubuntu等发行版的包管理.

> apt删除一个包并不会删除**已修改的用户配置文件**, 以防用户意外删除了包. 如果想通过 apt 命令删除一个应用及其配置文件，请在之前删除过的应用程序上使用 purge 子命令.

## 格式

## 例
```bash
# apt search drbd-utils # 类似`apt-cache search`, 但更好用
# apt-cache madison pouch # 列出软件包的所有版本
# apt install pouch=1.0.0-0~ubuntu # 安装指定版本的软件包
# apt-get install --reinstall aptitude # 重新安装package
# apt-get install --only-upgrade samba # 仅更新单个package
# apt list -a cifs-utils # package all version
# apt-cache policy cifs-utils # package all version, 推荐
# rmadison cifs-utils # package all version, 推荐
# apt-cache depends -i samba # 查看依赖
# apt install --simulate samba # 仅模拟不安装
# apt install --download-only samba # 仅下载, 不安装. `--download-only`=`-d`
# apt list --installed # 查看已安装的package
# apt-cache show cpio # 查看软件依赖
# apt show cpio # 查看软件信息, 包括依赖
# apt purge package_name # remove命令卸载指定软件包，但是留下一些包文件. 如果想彻底卸载软件包，包括它的文件，使用purge替换remove
# apt build-dep xxx # 获取构建包xxx的依赖(包括源码), 前提是取消`/etc/apt/sources.list*`中相应deb-src源的注释
# apt download redis # 下载deb
# apt-mark hold package_name # 升级时锁定选定的软件包 
# apt-mark unhold package_name
# add-apt-repository ppa:jonathonf/vim # add repo
# add-apt-repository -r ppa:jonathonf/vim # remove repo
```

> apt-file也可用于查找文件; apt-rdepends生成依赖树

# dpkg

## example
```
$ sudo dpkg -i --force-bad-verify  acl_2.2.52-3_amd64.deb # 跳过签名验证, `--ignore-depends=<x1>,<x2>`忽略依赖
$ dpkg -I xxx.deb # 查看依赖
$ dpkg -r kylin-nm # 卸载包
$ dpkg -S file # 这个文件属于哪个已安装软件包
$ dpkg -L package # 列出软件包中的所有文件
$ dpkg -s package # 列出软件包中的描述信息
$ echo "<package_name> hold" | sudo dpkg --set-selections  ##锁定软件包
$ dpkg --get-selections | grep hold  ##显示锁定的软件包列表
$ echo "<package_name> install" | sudo dpkg --set-selections  ##解除对软件包的锁定
$ dpkg --info xxx.deb | grep Depends # 查看deb包的依赖
$ dpkg-deb -c xxx.deb # 查看deb内容
```

按文件搜索package也可直接使用[debian package服务](https://www.debian.org/distrib/packages)

# snap
```bash
snap find solitaire # 查找指定程序
snap info solitaire
snap remove gnome-42-2204 && snap install gnome-42-2204 # 重装
snap refresh # 更新snap包
snap list --all
snap remove gtk-common-themes --revision=123 # 删除disabled的版本
snap set system refresh.retain=2 # 保留最近的两个版本, 这已是默认设置
bash -c 'rm -rf /var/lib/snapd/cache/*' # 清理cache

# --- 推荐使用clean-snap.sh
vim clean-snap.sh
#!/bin/bash
# Removes old revisions of snaps
# CLOSE ALL SNAPS BEFORE RUNNING THIS
set -eu
snap list --all | awk '/disabled/{print $1, $3}' |
    while read snapname revision; do
        snap remove "$snapname" --revision="$revision"
    done
```

## deb
参考:
- [Ubuntu下制作deb包的方法详解](https://blog.csdn.net/gatieme/article/details/52829907)
- [对一个deb包的解压、修改、重新打包全过程方法](https://blog.csdn.net/yygydjkthh/article/details/36695243), 重新打包后apt install遇到`Failed to fetch ....deb  Size mismatch`, 改用dpkg安装即可, 因为"Packages.gz"里的元信息是旧的.


deb包本身有三部分组成：
- 数据包，包含实际安装的程序数据，文件名为 data.tar.XXX

    data.tar.gz包含的是实际安装的程序数据，而在安装过程中，该包里的数据会被直接解压到根目录(即 / )，因此在打包之前需要根据文件所在位置设置好相应的文件/目录树
- 安装所需的信息及控制脚本包, 包含deb的安装说明，标识，脚本等，文件名为 control.tar.gz

    一般有 5 个文件：
    控制文件    描述
    control     用了记录软件标识，版本号，平台，依赖信息等数据
    preinst     在解包data.tar.gz前运行的脚本
    postinst    在解包数据(即安装)后运行的脚本
    prerm   卸载时, 在删除文件之前运行的脚本
    postrm  卸载时, 在删除文件之后运行的脚本

- 最后一个是deb文件的一些二进制数据，包括文件头等信息，一般看不到，在某些软件中打开可以看到

deb本身可以使用不同的压缩方式. tar格式并不是一种压缩格式，而是直接把分散的文件和目录集合在一起，并记录其权限等数据信息. 之前提到过的 data.tar.XXX，这里 XXX 就是经过压缩后的后缀名. deb默认使用的压缩格式为gzip格式，所以最常见的就是 data.tar.gz. 常有的压缩格式还有 bzip2 和 lzma，其中 lzma 压缩率最高，但压缩需要的 CPU 资源和时间都比较长.

### dpkg重新打包
```bash
#解压出包中的文件到extract目录下
$ dpkg -X ../openssh-client_6.1p1_i386.deb extract/
#解压出包的控制信息extract/DEBIAN/下：
$ dpkg -e ../openssh-client_6.1p1_i386.deb extract/DEBIAN/
$ dpkg-deb -b extract/ build/ # build存放打包好的deb
$ ll build/
-rw-r--r-- 1 ufo ufo 1020014  7月  3 20:20 openssh-client_6.1p1_i386.deb # 验证方法为：再次解开重新打包的deb文件，查看在etc/ssh/sshd_config文件是否已经被修改
$ dpkg-scanpackages . | gzip -9c > Packages.gz #  制作本地软件源
```

## dpkg-buildpackage
选项:
- -nc : doesn't call the clean target, 因此无需重新编译(但可能还是需要编译少量内容)
- -uc : don't sign the changes file
- -us : unsigned source package

## deb debug package
打deb包时, 通过将可执行文件的符号表通过剥离成独立的 dbg 包, 称为 debug package. 正常情况下 -dbg.deb 不会安装.

新版本（debhelper/9.20151219 or newer in Debian）的 debhelper 已经把 -dbg.deb 改为 -dbgsym.deb，详情请见[DebugPackage](https://wiki.debian.org/DebugPackage).

## 打包
- [BuildingTutorial](https://wiki.debian.org/BuildingTutorial#Building_the_modified_package)
- [Easy way to create a Debian package and local package repository](https://linuxconfig.org/easy-way-to-create-a-debian-package-and-local-package-repository)
- [apt source包打包](https://www.debian.org/doc/manuals/apt-howto/ch-sourcehandling.zh-cn.html)
- [Debian 新维护者手册](https://www.debian.org/doc/manuals/maint-guide/)

## FAQ
### dpkg-deb: error: archive '<file>.deb' has premature member 'data.tar.gz' before
dpkg的bug: [dpkg无法解析tar.xz格式-xz compressed control.tar files not supported](https://bugs.launchpad.net/ubuntu/+source/dpkg/+bug/1730627)

升级dpkg版本(>=1.17.5), 即`apt install dpkg`.

原因:
对于软件安装包的提供者而言, 一定是希望安装包具有更好的兼容性. 最好可以使用xz压缩data部分, 仍然用gzip打control部分. 旧版的dpkg-deb, 默认会把control和data分开用不同的格式打包, control默认始终使用gzip的格式打包. 而新版的dpkg-deb(1.19.0)之后都会使用相同的格式压缩control和data. 如果你指定了-Z xz , 那就都是xz. 

还好, dpkg-deb提供了一个参数：--no-uniform-compression加上这一句就可以了. 

默认是：--uniform-compression, 代表使用统一的格式进行压缩. 加上--no-uniform-compression后不再统一, control使用gz压缩. 

### apt install报`Size mismatch`
下载到的deb软件包信息和源信息列表Packages记录(Packages.gz)的数据不相符, 可用`dpkg -i`安装

### apt install 安装的deb的缓存位置
ubuntu中由apt-get获得的文件包保存在/var/cache/apt/archives

### 删除snap
ref:
- [Ubuntu 22.04 禁用（彻底移除）Snap](https://sysin.org/blog/ubuntu-remove-snap/)

```bash
for p in $(snap list | awk '{print $1}'); do
  sudo snap remove $p
done # 需要多次执行, 直至提示`No snaps are installed yet`

sudo systemctl stop snapd
sudo systemctl disable --now snapd.socket

for m in /snap/core/*; do
   sudo umount $m
done

sudo apt autoremove --purge snapd

rm -rf ~/snap
sudo rm -rf /snap
sudo rm -rf /var/snap
sudo rm -rf /var/lib/snapd
sudo rm -rf /var/cache/snapd

sudo sh -c "cat > /etc/apt/preferences.d/no-snapd.pref" << EOL
Package: snapd
Pin: release a=*
Pin-Priority: -10
EOL # 禁止 apt 安装 snapd

reboot # 否则文件管理器可能还会看到残留的snap loop设备
```

### 路径`debian/rules`
dpkg-buildpackage的构建目录结构

### 删除`dpkg -l`显示状态为的"rc"包
```bash
# --- 方法1, since ubuntu 20.04, 推荐
apt purge '~c'

# --- 方法2
dpkg -l |grep "^rc" |cut -d " " -f 3|sudo xargs dpkg --purge

# --- 方法3
aptitude search ~c # list the residual packages
sudo aptitude purge ~c # purge them
```

### apt/dpkg log
- `/var/log/apt/term.log`
- `/var/log/dpkg.log`

### apt `Couldn't find any package by glob 'qemu'`
`apt remove qemu*` 改为 `apt remove "qemu*"`

### `debian/rules`中的`dh --list`没有输出`systemd`
根据[dh源码](https://github.com/Debian/debhelper/blob/master/dh)结合locate查找"*.pm", 发现`dh --list`输出的是`/usr/share/perl5/Debian/Debhelper/Sequence`里的文件.

根据`dpkg -S /usr/share/perl5/Debian/Debhelper/Sequence/systemd.pm`, `systemd.pm`属于`debhelper/dh-systemd`.

> 高版本debhelper已包含`systemd.pm`. 低版本包含在`dh-systemd`.

### 通过deb-src构建deb
```bash
# vim /etc/apt/source.list # 添加deb-src源
deb-src http://pl.archive.ubuntu.com/ubuntu/ natty main restricted
# apt update
# apt build-dep ccache # 安装构建ccache所需的依赖
# apt-get -b source ccache/apt source --compile ccache # 获取ccache源码并构建 = apt source ccache && cd ccache* && dpkg-buildpackage [--no-sign]/dpkg-buildpackage -rfakeroot -b -uc -us/debuild -b -uc -us
# dpkg -i ccache*.deb
```

### dpkg-buildpackage报Warning "Compatibility levels before 9 are deprecated"
将项目`debian/compat`中的数字改为9即可

### 创建debian/changelog
`dch --create`, 并使用`dch -i`插入新changelog.

### 构建deb报`dch: error: fatal error occurred while parsing debian/changelog`
debian/changelog格式有错误, 解决方法有两种:
1. 找出错误
1. 使用`dch --create`构建一份新的changelog, 并将内容替换为旧debian/changelog的第一条记录.

### 构建deb报`debsign: gpg error occurred!  Aborting....`
debuild默认构建deb需要gpg签名, 通过`man debuild`可用`debuild -i -us -uc -b`构建deb

### 构建deb报`Now running lintian liburing_2.1-1_amd64.changes ...\nE: liburing-dev: debian-changelog-file-missing`
通过`man debuild`可用`debuild --no-lintian`禁止运行lintian, `--no-lintian`必须是第一个参数. **通过liburing构建deb观察到该报错在构建出deb之后, 因此推测它不影响deb构建.**

### 构建deb后没找到deb
1. deb在当前目录的上一层
1. 指定了构建目录, 比如liburing在`/tmp/release/<os>/liburing`

### 修改dpkg依赖
1. 安装前

    ```bash
    $ dpkg-deb -x krb5-config_2.3kord_all.deb krb-tmp
    $ dpkg-deb --control krb5-config_2.3kord_all.deb krb-tmp/DEBIAN
    $ vim krb-tmp/DEBIAN/control 
    $ dpkg -b krb-tmp llala.deb
    dpkg-deb: building package 'krb5-config' in 'llala.deb'.
    ```
2. 安装后

    `sudo vim /var/lib/dpkg/status`后修改Depends行即可

### 手动清理软件
```bash
# dpkg -L icaclient # 罗列文件
# rm -rf ... # 清理上述列出的文件
# vim /var/lib/dpkg/status # 清理掉icaclient
# sudo dpkg --configure -a
# 清理罗列的icaclient文件
```

### [Virtual Package(虚拟包)](https://www.debian.org/doc/manuals/debian-faq/pkg-basics.en.html#virtual)

### 依赖所需版本
`ukui-biometric-manager : 依赖: libopencv-core4.2 (>= 4.2.0+dfsg) 但无法安装它`, 其中`>= 4.2.0+dfsg`是指`>=4.2.0 && <4.3.0`, 其实名称`libopencv-core4.2`中的`4.2`就是提示, 只允许是`4.2.x`

### apt-key添加key报`apt-key is deprecated. Manage keyring files in trusted.gpg.d instead`
```bash
# apt-key list # 查看key
Warning: apt-key is deprecated. Manage keyring files in trusted.gpg.d instead (see apt-key(8)).
/etc/apt/trusted.gpg
--------------------
pub   rsa4096 2017-05-08 [SCEA]
      1EDD E2CD FC02 5D17 F6DA  9EC0 ADAE 6AD2 8A8F 901A
uid           [ 未知 ] Sublime HQ Pty Ltd <support@sublimetext.com>
sub   rsa4096 2017-05-08 [S]
...
# sudo apt-key del 1EDDE2CDFC025D17F6DA9EC0ADAE6AD28A8F901A
# wget -qO - https://download.sublimetext.com/sublimehq-pub.gpg |gpg --dearmor |sudo tee /usr/share/keyrings/sublimehq-pub.gpg
# vim /etc/apt/sources.list.d/sublime-text.list 
deb [arch=amd64 signed-by=/usr/share/keyrings/sublimehq-pub.gpg] https://download.sublimetext.com/ apt/stable/
```

### `apt install ./xxx.deb`报`couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)`
原因: xxx.deb无法被_apt用户访问导致, 即权限问题. 将xxx.deb移到`/tmp`再安装即可

### apt update NO_PUBKEY
```bash
gpg --keyserver keyserver.ubuntu.com --recv-key 9C949F2093F565FF
gpg -a --export  9C949F2093F565FF | sudo apt-key add -
```

### 离线部署
```bash
创建存放目录：mkdir -p /home/ubuntu/packs；
安装软件包dpkg-dev:apt-get install dpkg-dev
拷贝dep包至存放目录：sudo cp -r /var/cache/apt/archives/* /home/ubuntu/packs；
进入packs目录下，生成包的依赖信息：dpkg-scanpackages packs /dev/null |gzip > packs/Packages.gz
制作iso包：mkisofs -r -o /home/ubuntu/ubuntu-16.04.5-amd64.iso /home/ubunut/packs
```

### Ubuntu安装kernel
ref:
- [`/~kernel-ppa/mainline`](https://kernel.ubuntu.com/~kernel-ppa/mainline/)

方法1:
```bash
wget -c https://kernel.ubuntu.com/~kernel-ppa/mainline/v6.0/amd64/linux-headers-6.0.0-060000_6.0.0-060000.202210022231_all.deb
wget -c https://kernel.ubuntu.com/~kernel-ppa/mainline/v6.0/amd64/linux-headers-6.0.0-060000-generic_6.0.0-060000.202210022231_amd64.deb
wget -c https://kernel.ubuntu.com/~kernel-ppa/mainline/v6.0/amd64/linux-modules-6.0.0-060000-generic_6.0.0-060000.202210022231_amd64.deb
wget -c https://kernel.ubuntu.com/~kernel-ppa/mainline/v6.0/amd64/linux-image-unsigned-6.0.0-060000-generic_6.0.0-060000.202210022231_amd64.deb
sudo apt install ./linux-*.deb
sudo apt remove linux-headers-6.0.0* linux-modules-6.0.0* linux-image-unsigned-6.0.0*
```

方法2:
```bash
sudo add-apt-repository ppa:cappelikan/ppa
sudo apt update
sudo apt install mainline
mainline-gtk    
```

### [Key is stored in legacy trusted.gpg keyring](https://linux.cn/article-15565-1.html)
apt-key已废弃, 并改用`/etc/apt/trusted.gpg.d`存储GPG 密钥.

方法1:
```bash
$ sudo apt-key list
...
pub   rsa4096 2016-02-18 [SCEA]
      DB08 5A08 CA13 B8AC B917  E0F6 D938 EC0D 0386 51BD
uid           [ unknown] https://packagecloud.io/slacktechnologies/slack (https://packagecloud.io/docs#gpg_signing) <abhishek@itsfoss>
$ sudo apt-key export 038651BD | sudo gpg --dearmour -o /etc/apt/trusted.gpg.d/slack.gpg # 038651BD 是 pub的最后8个字符
```

方法2:
```bash
$ sudo cp /etc/apt/trusted.gpg /etc/apt/trusted.gpg.d # **推荐**
```

### snap调试执行
`SNAP_CONFINE_DEBUG=1 snap run firefox`

### snap list显示`broken`
查询系统日志, 进行错误处理, 比如重装相应snap package

### snap 操作报`run hook "connect-plug-host-hunspell": cannot perform operation: mount --rbind /dev /tmp/snap.rootfs_qopigF//dev: No such file or directory`

```bash
systemctl stop snapd
systemctl stop snapd.socket
reboot
sudo rm -rf /var/lib/snapd/cache/*
sudo rm -rf /tmp/snap.* # 如果还是不能删除, 那stop snapd.socket/snapd时还需要disable
systemctl start snapd.socket
systemctl start snapd # 查看日志补全缺失的package
```

### update kernel
1. debian 12
```bash
$ --- use apt repo: <bookworm>-backports
$ sudo apt-cache search linux-image
$ sudo apt install linux-headers-6.6.13+bpo-amd64 linux-image-6.6.13+bpo-amd64
```

### dpkg script位置
`/var/lib/dpkg/info/*.{postinst,postrm,preinst,prerm}`

### `apt install`报`unknown system group 'plocate' in statoverride file; the system group got removed`
plocate组不存在, 打开`/var/lib/dpkg/statoverride`, 删除plocate所在行, 再重现安装plocate即可

### apt dist-upgrade/upgrade区别
upgrade:系统将现有的Package升级,如果有相依性的问题, 而此相依性需要安装其它新的Package或影响到其它Package的相依性时,此Package就不会被升级,会保留下来. 
dist-upgrade:可以聪明的解决相依性的问题, 如果有相依性问题, 需要安装/移除新的Package,就会试着去安装/移除它. (所以通常这个会被认为是有点风险的升级) 

apt upgrade 和 apt dist-upgrade 本质上是没有什么不同的, 仅在处理依赖上有差异: dist-upgrade会识别出当依赖关系改变的情形并作出处理，而upgrade对此情形不处理

### apt更新包含phased upgrades
`apt -o APT::Get::Always-Include-Phased-Updates=true upgrade`, 也可将该配置写入apt配置文件(不推荐)