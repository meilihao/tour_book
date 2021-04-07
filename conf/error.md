## NotFound

### No such file or directory
> [关于usr/bin/ld: cannot find -lxxx问题总结](http://eminzhang.blog.51cto.com/5292425/1285705)
- gli/gli.hpp : apt install libgli-dev
- glm/glm.hpp : apt install libglm-dev
- assimp : apt install libassimp-dev
- apt install libx11-xcb-dev
- Could NOT find X11_XCB : apt install libx11-dev
- Could not find X11 : apt install libx11-dev
- Xlib:  extension "NV-GLX" missing on display : apt install mesa-vulkan-drivers # https://bbs.deepin.org/forum.php?mod=viewthread&tid=143398&page=1#pid363502
- /usr/bin/ld: cannot find -lpng

### 'aclocal-1.15' is missing on your system
`sudo apt install automake`或系统存在更高版本的aclocal, 比如`aclocal-1.16`

### zlib header files not found
sudo apt install zlib1g-dev

## not find
1. libsatlas.so
```
sudo apt-get install libatlas-base-dev
```

### debuild not found
sudo apt install devscripts

### OpenSSL header files not found
sudo apt install libssl-dev

### curses not found
sudo apt install libncurses5-dev

### libevent not found
sudo apt install libevent-dev

### mbed TLS libraries not found
sudo apt install libmbedtls-dev

### The Sodium crypto library libraries not found
sudo apt install libsodium-dev

### The c-ares library libraries not found
sudo apt install libc-ares-dev

### virsh not found
apt install libvirt-bin # for ubuntu 14.04

### Couldn't find libev
sudo apt install libev-dev

### /usr/bin/ld: 找不到 -lpam
sudo apt install libpam0g-dev

### librocksdb.so: undefined reference to `ZSTD_freeCDict`
`nm -A librocksdb.so|grep ZSTD_freeCDict`可知, librocksdb.so需要链接ZSTD_freeCDict, 但linker找不到"ZSTD_freeCDict", 通常是rocksdb太新或libzstd-dev太旧, 版本不匹配导致.

### virsh unknown os type hvm
原因: qemu 虚拟机配置中的`<emulator>...</emulator>`路径不存在, 即未安装qemu kvm.
apt install qemu-kvm qemu-system-x86 # for ubuntu 16.04, 必须同时安装qemu-system-x86,否则会报"qemu-system-x86_64: not found"

同时用`lsmod |grep kvm`检查kernel是否已加载kvm, 不然用`modprobe kvm`加载.

> virsh create启动kvm虚拟机前提: cpu支持虚拟化, 检测方法`grep -E "^flags.*(vmx|svm)" /proc/cpuinfo`.

> 可用`kvm-ok`检查kvm环境是否ok.

### /usr/local/lib/libdbus-1.so.3: version 'LIBDBUS_PRIVATE_1.12.16' not found (required by dbus-uuidgen)
具体原因未知, 怀疑是系统内有多个版本三libdbus.so导致.

解决: `sudo apt install --reinstall libdbus-1-3`

> 不行的话, 先删除错误的libdbus-1.so.3, 再重装.

### Cannot find asciidoc in PATH

直接`apt install asciidoc`安装需要下载几十个依赖不可取.
因此先安装xmlto,再用`apt -f install`补全依赖,最后安装asciidoc.xxx.deb

- [deepin 15.3](http://packages.deepin.com/deepin/pool/main)
- [ubuntu 16.04](packages.ubuntu.com)

### 调试日志刷屏: `kernel: Incremented dev->dev_cur_ordered_id: xxx for SIMPLE`
临时: `echo 'file target_core_transport.c -p' > /sys/kernel/debug/dynamic_debug/control`
永久: 修改启动kernel的kernel arg, 追加`dynamic_debug.verbose=0`

### no space left on device

1. 检查磁盘空间(`df -h`)
2. 检查inode(`df -i`)
3. 检查`/proc/sys/fs/inotify/max_user_watches`,inotify达到上限(??求查询inotify使用的句柄数)
```
$ sudo sysctl fs.inotify.max_user_watches=8192 # 临时修改
$ vim /etc/sysctl.conf # 添加max_user_watches=8192，然后sysctl -p生效,永久生效
```

### makeinfo: command not found
`sudo apt install texinfo`

### NMI watchdog: BUG: soft lockup - CPU#2 stuck for 23s! [systemd-logind:893]
env: ubuntu 14:04, kernel:4.4.0

该[issue](https://bugzilla.redhat.com/show_bug.cgi?id=1154286#c18)里有说在4.4.3中已修复该问题, 但其`Comment 19`提示之前的修复未完全解决该问题.

### vscode : Run VS Code with admin privileges so the changes can be applied.

```shell
sudo code --user-data-dir=/home/chen/.vscode .
```

### /usr/lib/x86_64-linux-gnu/libmirprotobuf.so.3: undefined symbol: _ZNK6google8protobuf11MessageLiteXXX

卸载旧版libprotobuf-lite,删除/usr/local/lib/libprotobuf-lite.so*,重新安装相应版本即可(sudo apt --reinstall install libprotobuf-lite10).

### SELinux is preventing systemd from read access on the file xxx.service

使用sudo sealert -a /var/log/audit/audit.log查看具体日志，里面有解决方案.

> [参考](http://www.tuicool.com/articles/myYv6v)

### [`qt.qpa.plugin: Could not load the Qt platform plugin "xcb" in "" even though it was found`](https://github.com/visualfc/liteide/issues/1165)
```bash
# 开启QT_DEBUG_PLUGINS, 获取详细错误`Cannot load library /opt/liteide/plugins/platforms/libqxcb.so: (/opt/liteide/plugins/platforms/../../lib/libQt5XcbQpa.so.5: undefined symbol: _ZNK14QPlatformTheme14fileIconPixmapERK9QFileInfoRK6QSizeF6QFlagsINS_10IconOptionEE)`
QT_DEBUG_PLUGINS=1 ./liteide
```

解决方法: `export LD_LIBRARY_PATH="/PATH/TO/liteide/lib:$LD_LIBRARY_PATH" && ./liteide`, 建议加入`~/.bashrc`.

### mpt3sas_cm0: log_info(0x3112010a): originator(PL), code(0x12), sub_code(0x010a)
参考:
- [mpt2sas0: log_info(0x31120303) 问题分析与解决](https://blog.csdn.net/weixin_44648216/article/details/104070284)
- [hotplug_sata_drive.md](https://github.com/huataihuang/cloud-atlas-draft/blob/master/os/linux/kernel/storage/hotplug_sata_drive.md)
- [mpt2sas故障处理](https://huataihuang.gitbooks.io/cloud-atlas/content/storage/das/mpt2sas/troubleshooting/mpt2sas_offline_fail_disk.html)
- [使用命令行工具对LSI阵列卡进行高效管理](https://blog.51cto.com/1130739/1771506)

> mpt3sas is the driver for the SATA host bus adapter.

log_info是一个U32长度的变量, 后面是对它的解释, 定义在[`IOC LOGINFO defines`](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L20)

log_info逻辑在[_base_sas_log_info](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/scsi/mpt3sas/mpt3sas_base.c#L1230).

originator在:
```c
// https://elixir.bootlin.com/linux/v5.10.10/source/drivers/scsi/mpt3sas/mpt3sas_base.c#L1257
	switch (sas_loginfo.dw.originator) {
	case 0:
		originator_str = "IOP";
		break;
	case 1:
		originator_str = "PL";
		break;
	case 2:
		if (!ioc->hide_ir_msg)
			originator_str = "IR";
		else
			originator_str = "WarpDrive";
		break;
	}
```

[originator](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L31):
0x0 对应 IOP 意思为IO Processor
0x1 对应 PL 意思为 Protocol Layer
0x2 对应 IR 意思为 Intergrated RAID

`code(0x12)的`定义是[PL_LOGINFO_CODE_ABORT](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L137)

`sub_code(0x010a)`的定义是[PL_LOGINFO_SUB_CODE_<XXX>](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L144)

sub_code(0x010a)=`PL_LOGINFO_SUB_CODE_OPEN_FAILURE|PL_LOGINFO_SUB_CODE_OPEN_FAIL_BREAK`

other example:
```log
0x31110d01 == (
MPI_IOCLOGINFO_TYPE_SAS |
IOC_LOGINFO_ORIGINATOR_PL |
PL_LOGINFO_CODE_RESET |
PL_LOGINFO_SUB_CODE_SATA_LINK_DOWN |
PL_LOGINFO_SUB_CODE_OPEN_FAIL_NO_DEST_TIME_OUT
)
```

> /var/log/messages日志中不断出现hard reset，表明存储卡出现了异常即服务器RAID卡需要替换维修.

> 存储hba卡通常是来自LSI公司（LSI Corporation）,  一般地，支持RAID 5的卡，称其为阵列卡，都可以使用LSI官方提供的MegaCli, SAS2IRCU, SAS3IRCU等工具来管理; 而不支持RAID 5的卡，称其为SAS卡，使用lsiutil工具来管理. HP的服务器则使用其特有的hpacucli工具来管理.

### ext4变成read only, syslog报"ext4_journal_check_start detected aborted journal"
参考:
- [File system going in to read-only](https://access.redhat.com/solutions/35122)

原因: All paths to LUN were lost. If the kernel finds fatal corruption on the disk or if certain key IOs like journal writes start failing, it will remount the filesystem as read-only.

修复方法:
```bash
umount /dev/<dev> # 从文件系统中删除日志
tune2sf -O ^has_journal /dev/<dev>
e2fsck -f /dev/<dev>
tune2sf -j /dev/<dev>
```

### chrome 55 没有flash

`chrome://plugins`里的Adobe Flash Player显示: 

Location:  internal-not-yet-present // 即flash并没有下载

运行:
```
google-chrome-stable --proxy-server="socks5://127.0.0.1:1080" // 需梯子
```

再在`chrome://components/`下载`Adobe Flash Player`,重启即可.

> 其实就是在`~/.config/google-chrome/PepperFlash`下载了一个flash的版本(文件夹名是flash对应的版本号)和latest-component-updated-flash校验文件.
> 相应的命令行:`/usr/bin/google-chrome-stable %U --ppapi-flash-path=/home/chen/.config/google-chrome/PepperFlash/24.0.0.186/libpepflashplayer.so --ppapi-flash-version=24.0.0.186`

### socks5转http

```
apt install privoxy
vim /etc/privoxy/config
systemctl restart privoxy
```

config变动:
```
listen-address  127.0.0.1:6060 // 6060也就是你需要的http输出的端口
forward-socks5   /   127.0.0.1:1080  . // 1080也就是socks5输入的端口
```
```
//开启privoxy 服务就行
sudo  service  privoxy start 
// 设置http 和 https 全局代理
export http_proxy='http://localhost:8118'
export https_proxy='http://localhost:8118'
test : 
wget www.google.com  # 如果返回200 ，并且把google的首页下载下来了，那就是成功了
 ```

其他类似软件: Polipo

### apt update 由于没有公钥，无法验证下列签名： NO_PUBKEY A1715D88E1DF1F24
```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys A1715D88E1DF1F24
```

### curl: (60) SSL certificate problem: unable to get local issuer certificate
curl对https服务端证书的检查未通过,解决:
1. 加`-k`跳过服务端证书的检查
1. 时`--cacert`,检查服务端证书

### apt purge xxx,报"子进程 已安装 post-removal 脚本 返回错误状态 1"
一般来说是由于我们在安装/卸载的过程中突然中止, 所以导致的环境异常, 软件已经可能已安装/已卸载了, 但是系统的信息却没有更新:
- 软件的状态信息有误, 状态信息在`/var/lib/dpkg/status`, 删除相应记录即可.
- 软件的配置信息不全, 位于`/var/lib/dpkg/info/.*`, 删除相应的文件即可.

一般来说, 前面两种方法之一即可解决该问题, 如果还是没觉得, 可以重建info列表:
1. 首先将info文件夹更名备份 : `sudo mv /var/lib/dpkg/info /var/lib/dpkg/info_old`
1. 再新建一个新的info文件夹,更新缓存信息,恢复info文件夹的内容 : `sudo mkdir /var/lib/dpkg/info && sudo apt update, apt -f install`
1. 执行完上一步操作后会在新的info文件夹下生成一些文件，现用这些文件覆盖info_old文件夹的内容,`sudo mv /var/lib/dpkg/info/* /var/lib/dpkg/info_old`
1. 把自己新建的info文件夹删掉,再把以前的info文件夹重新改回名字 : `sudo rm -rf /var/lib/dpkg/info && sudo mv /var/lib/dpkg/info_old /var/lib/dpkg/info`

### wireshark : Lua: Error during loading: [string "/usr/share/wireshark/init.lua"]:44:dofile has been disabled due to runing Wireshark as superuser.
使用`sudo wireshark`启动时碰到该文件.

解决方法(两种):
- 编辑init.lua文件的倒数第二行,`sudo vim /etc/wireshark/init.lua`,改为`--dofile("console.lua")`
- 编辑init.lua,`sudo vim /etc/wireshark/init.lua`,直接禁用lua即`disable_lua = true`

### Failed at step USER spawning /bin/etcd: No such process
`systemctl start xxx.service`报`etcd.service: Failed at step USER spawning /bin/etcd: No such process`,原因etcd.sevice文件中使用的`User=etcd`还不存在
```
 sudo useradd etcd -M -s /sbin/nologin 
```

### vmware-hostd 占用 443
vmware workstation: Edit -> Preferences -> Shared VMs -> Disable "Enable virtual machine sharing and remote access".

### Ubuntu 无法修改 file-max
场景:
1. `ulimit -n 65535`报`ulimit: open files: 无法修改 limit 值: 不允许的操作`
1. `echo "fs.file-max = 65535" >> /etc/sysctl.conf`后再运行`sysctl -p`不生效
1. 向`/etc/security/limits.conf`添加
```
* soft nofile 65535
* hard nofile 65535
```
重启后不生效

有时候经过上面的更改后使用ulimit -n会看到默认值并没有改变，我在ubuntu 18.04中就遇到这种情况.

需要在`/etc/pam.d/common-session`中加入`session required pam_limits.so`，再使用su username登录当前用户，然后 就可以使用ulimit命令了. 原因可能是gnome terminal默认是none-login的，所以我们在配置文件中的修改并没有影响到当前的terminal.

### 百度网盘Linux启动停留在启动界面的解决办法
先`rm -rf ~/baidunetdisk`, 再重启即可.

### syslog.socket: Socket service syslog.service not loaded, refusing
启动freeswith.service报错:
```
$ sudo systemctl restart freeswitch.service
...
12月 25 13:25:04 chen-pc systemd[1]: syslog.socket: Socket service syslog.service not loaded, refusing.
12月 25 13:25:04 chen-pc systemd[1]: Failed to listen on Syslog Socket.
...
$ sudo systemctl start syslog.socket
...
12月 25 13:25:04 chen-pc systemd[1]: syslog.socket: Socket service syslog.service not loaded, refusing.
12月 25 13:25:04 chen-pc systemd[1]: Failed to listen on Syslog Socket.
...
```

原因是系统的log服务不见了, 解决方法:
```
$ sudo apt install rsyslog
$ sudo systemctl status rsyslog
```

rsyslog启动后, syslog.socket也会自行起来.

### sudo 找不到命令
在`~/.bashrc`里追加`alias sudo='sudo env PATH=$PATH'`即可.

### neovim乱码
参考:
 - [Nvim shows weird symbols (�[2 q) when changing modes](https://github.com/neovim/neovim/wiki/FAQ#nvim-shows-weird-symbols-2-q-when-changing-modes)
 - [[RDY] Fix incorrect DECSCUSR fixup codes #6997](https://github.com/neovim/neovim/pull/6997)

 问题出现在xterm及其兼容版本下, 解决方法:
 ```sh
 $ echo 'set guicursor=' > ~/.config/nvim/init.vim # 禁用guicursor
 ```

 > 查看当前使用的term: `echo $TERM`
 > 查看系统支持的term: `tree /usr/share/terminfo`

### Ubuntu `do-release-upgrade -d` 16.04 -> 18.04 报错:
```log
authenticate 'bionic.tar.gz' against 'bionic.tar.gz.gpg'
gpg exited 1
Debug information:


gpg: Signature made Mon 01 Jul 2019 04:13:12 PM CST using RSA key ID C0B21F32
gpg: /tmp/ubuntu-release-upgrader-yoyius5d/trustdb.gpg: trustdb created
gpg: BAD signature from "Ubuntu Archive Automatic Signing Key (2012) <ftpmaster@ubuntu.com>"

Authentication failed
Authenticating the upgrade failed. There may be a problem with the network or with the server.
```

解决方法:
[`sudo apt-key update`](https://ubuntuforums.org/showthread.php?t=2388040)

### 该软件包现在的状态极为不妥
`dpkg -i`安装报错, `apt -f install`也无法修改.

解决方法:
```
rm -rf /var/lib/dpkg/info/${error's package}*
dpkg --remove --force-remove-reinstreq ${error's package}
dpkg --purge --force-remove-reinstreq ${error's package}
```

### sougou输入法候选词乱码
很可能是: sougou当前是繁体状态, 而当前系统字符集不支持繁体, 切换到简体状态即可.

### 无法修正错误，因为您要求某些软件包保持现状，就是它们破坏了软件包间的依赖关系
```
apt install mariadb-server
正在读取软件包列表... 完成
正在分析软件包的依赖关系树       
正在读取状态信息... 完成       
有一些软件包无法被安装。如果您用的是 unstable 发行版，这也许是
因为系统无法达到您要求的状态造成的。该版本中可能会有一些您需要的软件
包尚未被创建或是它们已被从新到(Incoming)目录移出。
下列信息可能会对解决问题有所帮助：

下列软件包有未满足的依赖关系：
 mariadb-server : 依赖: mariadb-server-10.4 (>= 1:10.4.10+maria~xenial) 但是它将不会被安装
E: 无法修正错误，因为您要求某些软件包保持现状，就是它们破坏了软件包间的依赖关系。
```

使用了mariadb的[离线安装包](https://mariadb.com/downloads/), 其中mariadb-server-10.4_10.4.10+maria~xenial_arm64.deb明明存在, 还是报错.

原因: mariadb-server-10.4缺其他依赖, 但apt不会提示, 上面的apt文案有点误导的意思.

解决方法:
1. 联网
1. 运行`dpkg -i mariadb-server-10.4_10.4.10+maria~xenial_arm64.deb mysql-common_10.4.10+maria~xenial_all.deb`, 看看缺哪些依赖
1. `apt -f install`
1. `apt install mariadb-server`

ps: **推荐[使用在线安装mariadb](https://downloads.mariadb.org/)**

### exo-helper-1: not found
env: ubuntu 19.04 + xfce4

解决方法:
```
sudo apt install libexo-1-0
```

### fcitx 无输入法托盘
```
$ cd ~/.config
$ rm -rf fcitx*
$ fcitx # 不能使用sudo
```

### 随桌面环境自启动
对于支持 xdg 标准的桌面环境，例如 gnome，kde，xfce，lxde， 可以将文件 安装目录/share/applications/fcitx.desktop 建立符号链接或者复制到 ~/.config/autostart/ 或者 /etc/xdg/autostart（/usr/local/etc/xdg/autostart/） 目录里.

### 搜狗输入法卸载
1. 使用apt卸载
1. 删除配置
```sh
$ cd ~/.config
$ sudo rm -rf SogouPY*
$ sudo rm -rf sogou*
$ rm -rf fcitx*
```
1. 重启

### 搜狗输入法无法运行在ubuntu 19.04 gnome/xfce4下
使用fcitx的其他中文输入法

ps: ubuntu 19.04 xfce4 用ibus也无法输入中文.

### make sure that your system can connect to api.snapcraft.io
在`/lib/systemd/system/snapd.service`追加:
```text
[Service]
Environment=http_proxy=http://proxy:port
Environment=https_proxy=http://proxy:port
```

然后重新加载snapd服务，运行以下命令：
```sh
sudo systemctl daemon-reload
sudo systemctl restart snapd.service
```

方法二:
在/etc/environment(snapd会读取它，应用其中指定的配置信息)文件中加入以下两行：
```text
http_proxy=http://[服务器地址]:[端口号]
https_proxy=http://[服务器地址]:[端口号]
```
然后重启snapd服务，运行以下命令：
```sh
sudo systemctl restart snapd
```

### `fcitx -d`: getting session bus failed: //bin/dbus-launch terminated abnormally with the following error: Autolaunch error: X11 initialization failed.

执行程序的用户不对. 我是用sudo执行的程序，而session bus需要访问启动用户所在home目录的隐藏文件夹`/home/xxxxx/.dbus/`,
该目录在用户目录xxxx下. 因此，只需要去掉sudo用普通用户执行即可.

### qt.qpa.plugin: Could not load the Qt platform plugin "xcb" in "" even though it was found
`sudo wireshark`时碰到的错误:
```
Invalid MIT-MAGIC-COOKIE-1 keyqt.qpa.xcb: could not connect to display :1.0
qt.qpa.plugin: Could not load the Qt platform plugin "xcb" in "" even though it was found.
This application failed to start because no Qt platform plugin could be initialized. Reinstalling the application may fix this problem.

Available platform plugins are: eglfs, linuxfb, minimal, minimalegl, offscreen, vnc, xcb.
```

解决方法: 使用普通用户权限运行wireshark.

### Couldn’t run /usr/bin/dumpcap in child process: Permission denied
```sh
$ sudo usermod -a -G wireshark $USER
```

需注销或重启电脑.

### xshell 关闭 发送键盘输入
取消`工具`-`发送键输入到所有会话`

> `发送键盘输入`功能与`查看-撰写栏`功能类似.

### mktemp:failed to create file via template '/tmp/.colorlsXXX' : Read-only file system
dmesg提示`EXT4-fs error (device sda2): ext4_lookup:1601: inode #8288969: comm grep: deleted inode referenced: 272`.
文件系统变成只读, 通常是硬盘故障.

解决:
1. `fuser -m /dev/sda` # 查找占用
1. `umount` # 卸载
1. `smartctl --all /dev/sda` 查看硬盘信息
1. `badblocks -s -v  /dev/sda` 检查坏道
1. 尝试`mount -o remount,rw /`恢复可写

重启后使用`fsck -y /dev/sdb2`修复后, 再reboot.

### 麒麟系统网络配置
麒麟默认用Network Manager进行网络管理, 即nmcli, 配置信息在`/etc/NetworkManager/system-connections`下

配置完成后需`service network-manager restart`.

#### listen unix /run/zsysd.sock: bind: address already in use
```
# ss -anlp # 没看到在使用
# rm ${unix_socket_path} # 直接删除即可
```

#### firefox清除HSTS缓存
1. 首选项 - 隐私与安全 - 清除历史数据, 选择合适的时间段+数据tab的两个选项.

### gem install 报错`unable to convert xxx to US-ASCII for xxx, skipping`
```
# gem install ffi -v 1.9.10
# unable to convert "\xE2" to UTF-8 in conversion from ASCII-8BIT to UTF-8 to US-ASCII for lib/ffi/library.rb, skipping
# gem install fpm -v 1.3.3
# unable to convert U+91CE from UTF-8 to US-ASCII for lib/backports/tools/normalize.rb, skipping
# gem install ffi -v 1.9.10 --no-rdoc --no-ri # gem install生成文档时报错, 不生成文档即可.
```

### 'aclocal-1.15' is missing on your system
Ubuntu19.10使用的是aclocal-1.16.

[运行`./configure`之前运行`autoreconf -f -i`](https://stackoverflow.com/questions/33278928/how-to-overcome-aclocal-1-15-is-missing-on-your-system-warning/33279062)

### linux蓝牙无法连接(开关已打开)
```
$ sudo journalctl -f 
3月 21 11:44:56 chen-pc bluetoothd[1881]: Failed to set mode: Blocked through rfkill (0x12)
...
$ rfkill list                                                                                                                                         11:47:23
0: phy0: Wireless LAN
	Soft blocked: no
	Hard blocked: no
1: hci0: Bluetooth
	Soft blocked: yes # 命令deepin控制台的蓝牙开关已打开, 却还是显示blocked.
	Hard blocked: no
$ rfkill unblock 1
```

bluetoothctl命令可查看蓝牙状态, 比如`scan on`监听蓝牙设备的变化.
blueman是管理蓝牙的gui工具.

### sh: 0: getcwd() failed: No such file or directory
刚刚把某个目录给删除了，而命令还停在原来的目录上，因此出现了这种问题, 退回到上级还存在的目录即可.

### superblock last mount time ${last_mount_time, now_time} is in the future
主板时间与`superblock last mount time`冲突

### /bin/bash^M: bad interpreter
文件使用了windows换行.

方法:
- 用vim打开, 再`:set ff?`查看`fileformat`是否为`doc`, 即是否使用了windows的换行.
- 用`cat -A xxx`查看文档换行

解决方法: `vim` + `:set ff=unix` + `:wq`

### arm飞腾硬件时间变慢
```bash
# ntpdate -p 192.168.0.223 && hwclock -w && data # 同步ntp, 硬件, 系统时间
# -- 10m后查看, 发现硬件时间变慢
# date && hwclock
```

### 电脑无故重启可能与ACPI有关系(无硬件问题)
acpi与内核可能有不兼容的问题, kernel启动参数可追加[`acpi=ht`](https://wiki.ubuntu.org.cn/UbuntuHelp:BootOptions/zh#.E5.8F.82.E6.95.B0:_acpi.3Dht)

### deepin 20 按键盘出发的主板蜂鸣声
```bash
$ sudo rmmod pcspkr # 临时关闭
$ sudo deepin-editor /etc/modprobe.d/nobeep.conf # 永久关闭, 需重启
blacklist pcspkr
blacklist snd_pcsp
```

## fs
### EXT4-fs error (device sdb4) ext4_find_entry:1436 inode #2 comm pvestatd reading directory lblock 0
[怀疑是sata接口接触不良](https://m.newsmth.net/article/KernelTech/75125?p=1)

### mknod: cannot set permissions of 'console': Operation not supported
在docker container(已使用`--privileged`)的chroot环境中执行`mknod -m 640 console c 5 1`报错, 推测是chroot环境权限受限导致.

解决方法: 在chroot的外层即container中执行mknod即可.

## lib
### /lib/x86_64-linux-gnu/libm.so.6: version `GLIBC_2.27' not found (required by xxx)
缺少GLIBC_2.27, `2.27`是xxx需要的最高glibc version.

命令`strings /lib/x86_64-linux-gnu/libm.so.6 | grep GLIBC_`支持查找该so支持的glibc版本, 输出是个范围, 只要本机`ldd --version`的glibc在该区间即可.

解决方法有2:
1. 安装指定版本的glibc
1. 在目标机器上编译程序并运行即可

### /var/run/nscd/socket no such file or directory
实际原因是`/var/run/nscd`不存在.

检查`/lib/systemd/system/nscd.service`的`section "Service"`是否存在属性`RuntimeDirectory=nscd`.

### undefined reference to `dlopen'
类似的还有`dlclose, dlerror, dlsym`, 是某个地方引用了`#include <dlfcn.h>`所致.

解决方法: 编译选项里加`-ldl`, 即`gcc DBSim.c -o DBSim -ldl`

### inlining failed in call to always_inline 'uint32_t __crc32cd(uint32_t, uint64_t)': target specific option mismatch
使用了特定的指令集, 因此[gcc编译时需要特定的选项`-march=armv8-a+crc+crypto`](https://github.com/OSSystems/meta-browser/issues/258)

### [登录火狐浏览器账号后没有同步数据, 或同步后发现不一致，或火狐浏览器账号不存在](https://blog.csdn.net/qq_40157728/article/details/103994649)
火狐账号登录(国际版) 和 火狐通行证(中国版) 是两个完全不同的账号体系，数据不互通.

两种账号可切换: 选项-同步-切换至xxx服务(在页面下方), 推荐使用"切换至全球服务".

### 更新正在运行进程的可执行文件时报`text file busy`的原因及解决方法
线上替换一个正在运行进程的文件时（包括二进制、动态库、需要读取的资源文件等）, 应避免使用`cp/scp操作`, 而需要使用`mv/rsync`作为替代.

原因：cp是将源文件截断然后写入新内容, 也就是说正在打开这个文件的进程可以立刻感知到修改. 修改文件内容很可能导致程序逻辑错误甚至崩溃. 而mv则是标记”删除“老文件，然后放一个新的同名文件过去, 也就是说老文件和新文件其实是两个不同文件（inode不同），只是名字一样而已. 正在打开老文件的进程不会受到影响. 如果进程使用了mmap打开某文件（比如载入so），如果目标文件被使用cp覆盖并且长度变小, 那么读取差额部分的地址时（在新文件中其实已经不存在了），会导致SIGBUS信号, 进而使进程崩溃.

至于可执行文件本身, 倒是不怕cp导致崩溃. 因为cp时会报”text file busy", 压根cp不了. 这时候也应该使用mv类操作, 替换完成后重启进程, 此时执行的就是新的可执行文件了.

### du和df的统计结果为什么不一样
参考:
- [详细分析du和df的统计结果为什么不一样](https://www.tuicool.com/articles/UjaYvu2)

du是通过stat命令来统计每个文件(包括子目录)的空间占用总和. 因为会对每个涉及到的文件使用stat命令，所以速度较慢.

df是读取每个分区的superblock来获取空闲数据块、已使用数据块，从而计算出空闲空间和已使用空间，因此df统计的速度极快(superblock才占用1024字节).

### apt update报错, packages metadata的长度错误
经核对报错文件的官方repo大小是正确的即仅是本地文件大小错误, 因此应该是本地apt cache出现了问题, 清理一下再更新即可.
```bash
sudo apt clean all # 先清理apt cache即可
sudo apt update
```

### wine乱码
方法1: `env LC_ALL=zh_CN.UTF-8 wine xxx.exe`

方法2: 安装微软雅黑:
```bash
cp msyh.ttc msyhbd.ttc msyhl.ttc ~/.wine/drive_c/windows/Fonts # 准备雅黑字体

vim msyh_font.reg # 添加一下内容
REGEDIT4
[HKEY_LOCAL_MACHINE\Software\Microsoft\Windows NT\CurrentVersion\FontLink\SystemLink]
"Lucida Sans Unicode"="msyh.ttc"
"Microsoft Sans Serif"="msyh.ttc"
"MS Sans Serif"="msyh.ttc"
"Tahoma"="msyh.ttc"
"Tahoma Bold"="msyhbd.ttc"
"msyh"="msyh.ttc"
"Arial"="msyh.ttc"
"Arial Black"="msyh.ttc"
regedit msyh_font.reg

vim ~/.wine/system.reg # 查找关键字"FontSubstitutes"
"MS Shell Dlg"="SimSun" => "MS Shell Dlg"="msyh"
"MS Shell Dlg 2"="SimSun" => "MS Shell Dlg"="msyh"

winecfg # 在"应用程序"选项卡修改"windows版本"为"windows 10"
```

### wine执行exe崩溃报"fixme:actctx:parse_depend_manifests Could not find dependent assembly L"Microsoft.VC80.MFCLOC" (8.0.50608.0)"
```bash
sudo apt install winetricks
winetricks vcrun2005
```

如过还是报该错, 删除`~/.wine`后安装vcrun2005, [安装雅黑字体解决文字方块], 再重新安装软件

### 删除wine的快捷方式
`rm -rf /.local/share/applications/wine/Programs/vivo/vivo手机助手/vivoAPK安装器.desktop`

### dpkg-query: 错误: --listfiles 需要一个有效的软件包名。而 libldap-2.4-2 不是: 软件包名 'libldap-2.4-2' 含糊不清， 它有一个以上的安装实例
deepin wine应用依赖i386, i386和amd64共存导致`apt dist-upgrade`失败, 先卸载wine应用再`apt autoremove`, 最后重新更新即可.

### deepin更新并重启后发现无法进入图形界面，而是变成字符界面，且用正确的账号无法在字符界面登录
`apt reinstall lightdm dde-session-shell`

安装lightdm重启后, 发现登录背景是黑色, 且输入账号登录界面会崩溃, 但字符界面已能正常登录.

`sudo apt reinstall dde*`, 或用`apt search dde |egrep "^dde" |grep -v "dbgsym" |grep -v "dev"`精确筛选dde相关包.

安装dde相关包后能正常登录, 但控制中心无法打开, `sudo dde-control-center --show`直接崩溃, 通过执行`sudo QT_DEBUG_PLUGINS=1 dde-control-center --show`定位是`libdeepin-recovery-plugin.so`导致, 用`dpkg -S libdeepin-recovery-plugin.so`定位包, 因此卸载`deepin-recovery-plugin`即可.