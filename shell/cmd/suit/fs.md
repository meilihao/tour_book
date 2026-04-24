# fs
### 内存盘
```sh
$ mkdir -p /data/tmpfs
```
```
# gedit /etc/fstab,加入以下内容:
# tmpfs
tmpfs /data/tmpfs tmpfs rw,nodev,nosuid,size=256m 0 0
```

> 不使用`~/tmpfs`的原因: 避免部分软件因为权限问题而无法访问, 比如virt-manager访问存储报`Permission Denied`

> 如果已经是挂载了, 可以使用remount进行扩容, 保证原有的数据不会丢失: `mount -t tmpfs -o size=4G [-o remount] tmpfs /data/tmpfs`

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

## FAQ
### 随桌面环境自启动
对于支持 xdg 标准的桌面环境，例如 gnome，kde，xfce，lxde， 可以将文件 安装目录/share/applications/fcitx.desktop 建立符号链接或者复制到 ~/.config/autostart/ 或者 /etc/xdg/autostart（/usr/local/etc/xdg/autostart/） 目录里.

### no space left on device

1. 检查磁盘空间(`df -h`)
2. 检查inode(`df -i`)
3. 检查`/proc/sys/fs/inotify/max_user_watches`,inotify达到上限(??求查询inotify使用的句柄数)
```
$ sudo sysctl fs.inotify.max_user_watches=8192 # 临时修改
$ vim /etc/sysctl.conf # 添加max_user_watches=8192，然后sysctl -p生效,永久生效
```

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

### EXT4-fs error (device sdb4) ext4_find_entry:1436 inode #2 comm pvestatd reading directory lblock 0
[怀疑是sata接口接触不良](https://m.newsmth.net/article/KernelTech/75125?p=1)

### mknod: cannot set permissions of 'console': Operation not supported
在docker container(已使用`--privileged`)的chroot环境中执行`mknod -m 640 console c 5 1`报错, 推测是chroot环境权限受限导致.

解决方法: 在chroot的外层即container中执行mknod即可.

### 更新正在运行进程的可执行文件时报`text file busy`的原因及解决方法
线上替换一个正在运行进程的文件时（包括二进制、动态库、需要读取的资源文件等）, 应避免使用`cp/scp操作`, 而需要使用`mv/rsync`作为替代.

原因：cp是将源文件截断然后写入新内容, 也就是说正在打开这个文件的进程可以立刻感知到修改. 修改文件内容很可能导致程序逻辑错误甚至崩溃. 而mv则是标记”删除“老文件，然后放一个新的同名文件过去, 也就是说老文件和新文件其实是两个不同文件（inode不同），只是名字一样而已. 正在打开老文件的进程不会受到影响. 如果进程使用了mmap打开某文件（比如载入so），如果目标文件被使用cp覆盖并且长度变小, 那么读取差额部分的地址时（在新文件中其实已经不存在了），会导致SIGBUS信号, 进而使进程崩溃.

至于可执行文件本身, 倒是不怕cp导致崩溃. 因为cp时会报”text file busy", 压根cp不了. 这时候也应该使用mv类操作, 替换完成后重启进程, 此时执行的就是新的可执行文件了.

### du和df的统计结果为什么不一样
参考:
- [详细分析du和df的统计结果为什么不一样](https://www.tuicool.com/articles/UjaYvu2)

du是通过stat命令来统计每个文件(包括子目录)的空间占用总和. 因为会对每个涉及到的文件使用stat命令，所以速度较慢.

df是读取每个分区的superblock来获取空闲数据块、已使用数据块，从而计算出空闲空间和已使用空间，因此df统计的速度极快(superblock才占用1024字节).

### 文件监控限制
[fs.inotify.max_user_watches=524288](https://code.visualstudio.com/docs/setup/linux#_visual-studio-code-is-unable-to-watch-for-file-changes-in-this-large-workspace-error-enospc)

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