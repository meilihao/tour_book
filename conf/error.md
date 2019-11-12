## NotFound

### No such file or directory
> [关于usr/bin/ld: cannot find -lxxx问题总结](http://eminzhang.blog.51cto.com/5292425/1285705)
- gli/gli.hpp : apt-get install libgli-dev
- glm/glm.hpp : apt-get install libglm-dev
- assimp : apt-get install libassimp-dev
- apt-get install libx11-xcb-dev
- Could NOT find X11_XCB : apt-get install libx11-dev
- Could not find X11 : apt-get install libx11-dev
- Xlib:  extension "NV-GLX" missing on display : apt install mesa-vulkan-drivers # https://bbs.deepin.org/forum.php?mod=viewthread&tid=143398&page=1#pid363502
- /usr/bin/ld: cannot find -lpng

### zlib header files not found
sudo apt-get install zlib1g-dev

### OpenSSL header files not found
sudo apt-get install libssl-dev

### curses not found
sudo apt-get install libncurses5-dev

### libevent not found
sudo apt-get install libevent-dev

### mbed TLS libraries not found
sudo apt-get install libmbedtls-dev

### The Sodium crypto library libraries not found
sudo apt-get install libsodium-dev

### The c-ares library libraries not found
sudo apt-get install libc-ares-dev

### Couldn't find libev
sudo apt-get install libev-dev

### Cannot find asciidoc in PATH

直接`apt install asciidoc`安装需要下载几十个依赖不可取.
因此先安装xmlto,再用`apt -f install`补全依赖,最后安装asciidoc.xxx.deb

- [deepin 15.3](http://packages.deepin.com/deepin/pool/main)
- [ubuntu 16.04](packages.ubuntu.com)

### no space left on device

1. 检查磁盘空间(`df -h`)
2. 检查inode(`df -i`)
3. 检查`/proc/sys/fs/inotify/max_user_watches`,inotify达到上限(??求查询inotify使用的句柄数)
```
$ sudo sysctl fs.inotify.max_user_watches=8192 # 临时修改
$ vim /etc/sysctl.conf # 添加max_user_watches=8192，然后sysctl -p生效,永久生效
```

### vscode : Run VS Code with admin privileges so the changes can be applied.

```shell
sudo code --user-data-dir=/home/chen/.vscode .
```

### /usr/lib/x86_64-linux-gnu/libmirprotobuf.so.3: undefined symbol: _ZNK6google8protobuf11MessageLiteXXX

卸载旧版libprotobuf-lite,删除/usr/local/lib/libprotobuf-lite.so*,重新安装相应版本即可(sudo apt --reinstall install libprotobuf-lite10).

### SELinux is preventing systemd from read access on the file xxx.service

使用sudo sealert -a /var/log/audit/audit.log查看具体日志，里面有解决方案.

> [参考](http://www.tuicool.com/articles/myYv6v)

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
1. 再新建一个新的info文件夹,更新缓存信息,恢复info文件夹的内容 : `sudo mkdir /var/lib/dpkg/info && sudo apt-get update, apt-get -f install`
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
$ sudo apt-get install rsyslog
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

### vnc4server启动时默认绑定localhost
因为vncserver没有使用TLSVnc, 不安全启动时默认绑定到localhost.

解决方法: 在`/etc/vnc.conf`中追加`$localhost = "no";`, 重启系统再重新运行`vncserver`即可.

`/etc/vnc.conf`的配置项`$geometry`支持修改分辨率.

ps: vnc推荐使用vnc4server.

### vnc viewer登录后灰屏/没有进入桌面
检查`$HOME/.vnc/xstartup`的配置.

xfce4的配置, 高分辨率会糊(来回切换分辨率就能解决或等会自行恢复), **推荐**:
```
#!/bin/sh

# Uncomment the following two lines for normal desktop:
unset SESSION_MANAGER
unset DBUS_SESSION_BUS_ADDRESS
# exec /etc/X11/xinit/xinitrc

[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
exec startxfce4
```

第二种xfce4配置, 有时成功有时失败:
```text
#!/bin/sh

# Uncomment the following two lines for normal desktop:
# unset SESSION_MANAGER
# exec /etc/X11/xinit/xinitrc

[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
xsetroot -solid grey
vncconfig -iconic &
x-terminal-emulator -geometry 80x24+10+10 -ls -title "$VNCDESKTOP Desktop" &
x-window-manager &
x-session-manager &
# desktop config
xfdesktop &
xfce4-panel &
xfsettingsd &
xfconfd &
xfce4-session &
xfwm4 &
```

ps: 如何找到上面的`desktop config`: 正常登录到系统, 看看它启动了哪些桌面环境相关的进程, 再结合网上资料, 补充完整即可.

这是网上ubuntu 19.04 + gnome的xstartup:
```
#!/bin/sh
[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
vncconfig -iconic &
dbus-launch --exit-with-session gnome-session & # 导致sougou输入法无法启动
```

其他:
```
#!/bin/sh

# Uncomment the following two lines for normal desktop:
# unset SESSION_MANAGER
exec /etc/X11/xinit/xinitrc

[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
```

ps: 直接删除xstartup也可进入桌面.

### exo-helper-1: not found
env: ubuntu 19.04 + xfce4

解决方法:
```
sudo apt install libexo-1-0
```

### 搜狗输入法无法运行在ubuntu 19.04 gnome/xfce4下
使用fcitx的其他中文输入法

ps: ubuntu 19.04 xfce4 用ibus也无法输入中文.