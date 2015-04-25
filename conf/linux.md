### Tools/工具

#### 制作u盘启动盘(linux下)

1 . [live-usb-install](http://sourceforge.net/projects/liveusbinstall/files/?source=navbar)
2 . [UNetbootin](http://sourceforge.net/projects/unetbootin/files/UNetbootin/),推荐

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
