## NotFound

### zlib header files not found
sudo apt-get install zlib1g-dev

### OpenSSL header files not found
sudo apt-get install libssl-dev

### curses not found
sudo apt-get install libncurses5-dev

### libevent not found
sudo apt-get install libevent-dev

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
listen-address  127.0.0.1:6060 // 6060也就是你需要的http输出的端口
forward-socks5   /   127.0.0.1:1080  . // 1080也就是socks5输入的端口

其他类似软件: Polipo
