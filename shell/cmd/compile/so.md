# so
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

## X11/Xlib.h not found
```
$ sudo apt install libx11-dev
```

## libwbclient.so.0
```shell
# Fedora22,运行“系统设置”程序
$ gnome-control-center --overview
gnome-control-center: error while loading shared libraries: libwbclient.so.0: cannot open shared object file: No such file or directory
# dnf update后出现这个error，应该是/usr/lib/libwbclient.so.0不存在引起的.
# sudo dnf install libwbclient时却提示已安装，通过find命令找到其已存在于/usr/lib/samba/wbclient/libwbclient.so.0.12,作软连接即可解决。
$ sudo ln -s  /usr/lib/samba/wbclient/libwbclient.so.0.12 /usr/lib/libwbclient.so.0
```

## FAQ
### /usr/local/lib/libdbus-1.so.3: version 'LIBDBUS_PRIVATE_1.12.16' not found (required by dbus-uuidgen)
具体原因未知, 怀疑是系统内有多个版本三libdbus.so导致.

解决: `sudo apt install --reinstall libdbus-1-3`

> 不行的话, 先删除错误的libdbus-1.so.3, 再重装.

### /lib/x86_64-linux-gnu/libm.so.6: version `GLIBC_2.27' not found (required by xxx)
缺少GLIBC_2.27, `2.27`是xxx需要的最高glibc version.

命令`strings /lib/x86_64-linux-gnu/libm.so.6 | grep GLIBC_`支持查找该so支持的glibc版本, 输出是个范围, 只要本机`ldd --version`的glibc在该区间即可.

解决方法有2:
1. 安装指定版本的glibc
1. 在目标机器上编译程序并运行即可

### undefined reference to `dlopen'
类似的还有`dlclose, dlerror, dlsym`, 是某个地方引用了`#include <dlfcn.h>`所致.

解决方法: 编译选项里加`-ldl`, 即`gcc DBSim.c -o DBSim -ldl`

### "pkg-config": executable file not found in $PATH
`apt-get install pkg-config`
