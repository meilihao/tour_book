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

