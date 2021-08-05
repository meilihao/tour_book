# virt-manager
## 构建
前提:
1. qemu
1. libvirt
1. 源码中的`INSTALL.md`

```
# -- https://github.com/archlinux/svntogit-community/blob/packages/virt-manager/trunk/PKGBUILD
wget https://releases.pagure.org/virt-manager/virt-manager-2.2.1.tar.gz
tar -xvzf virt-manager-2.2.1.tar.gz
cd virt-manager-2.2.1
apt install intltool
# python3 setup.py configure --default-hvs qemu,lxc
python3 setup.py install [--force]
```

## FAQ
参考:
- [How to compile virt-manager on Debian or Ubuntu](https://www.xmodulo.com/compile-virt-manager-debian-ubuntu.html)

### unable to execute 'intltool-update': No such file or directory
`apt install intltool`

### virt-manager运行报`Namespace LibvirtGlib not available`
参考gobject-introspection的[Namespaces are represented on disk by type libraries (.typelib files).](https://valadoc.org/gobject-introspection-1.0/GI.Repository.html), 应该是没有对应的`.typelib`文件, 它原本应由`apt install gir1.2-libvirt-glib-1.0`提供.

> GObject Introspection（简称 GI）用于产生与解析 C 程序库 API 元信息，以便于动态语言（或托管语言）绑定基于 C + GObject 的程序库, 具体可见[GObject Introspection 的作用与意义](http://garfileo.is-programmer.com/2012/2/20/gobject-introspection-introduction.32200.html).

> Typelibs将从环境变量GI_TYPELIB_PATH和`/usr/lib/girepository-1. 0/`中的路径加载.

安装libvirt-glib, 参考[libvirt-glib/trunk/PKGBUILD](https://github.com/archlinux/svntogit-community/blob/packages/libvirt-glib/trunk/PKGBUILD)
```bash
# 参考源码的INSTALL
apt install python3-gi libgirepository1.0-dev valac
wget https://libvirt.org/sources/glib/libvirt-glib-1.0.0.tar.gz
tar -xf libvirt-glib-1.0.0.tar.gz
cd libvirt-glib-1.0.0
./configure --enable-introspection --enable-vala
make
make install
cp /usr/local/lib/girepository-1.0/* /usr/lib/aarch64-linux/gnu/girepository-1.0 # 或设置变量GI_TYPELIB_PATH追加/usr/local/lib/girepository-1.0
```

### virt-manager运行报`No module named 'libvirt'`
未安装libvirt的python绑定: python3-libvirt. 参考[libvirt-python/trunk/PKGBUILD](https://github.com/archlinux/svntogit-community/blob/packages/libvirt-python/trunk/PKGBUILD), 安装即可.

```bash
wget https://libvirt.org/sources/python/libvirt-python-6.0.0.tar.gz
tar -xf libvirt-python-6.0.0.tar.gz
cd libvirt-python-6.0.0
python3 setup.py clean
python3 setup.py install --optimize=1
```

### virt-manager运行报`No module named 'libxml2'`
ubuntu 20.04: `apt install python3-libxml2`
ubuntu 16.04没有python3-libxml2, 用`pip3 install libxml2-python3`

### virt-manager运行报`pygobject3 3.22.0 or later is required.`
`pip3 install PyGobject==3.36.1` # version from Ubuntu 20.04

可能会遇到`No package 'cairo' found`, 解决方法: `apt install libcairo2-dev`, 再执行`pip3 install pycairo==1.16.2`

### virt-manager运行报`gtk 3.22.0 or later is required.`
它从`2.1.0`开始gtk必须是`3.24`及以上, 将virt-manager降级到`2.0.0`, gtk只要`3.14`(Ubuntu 16.04.6使用gtk 3.18).

### virt-manager运行报`cannot import name Vte, introspection typelib not found`
`apt install gir1.2-vte-2.91`

### virt-manager无法新建vm, `virt-manager --debug`报`cannot import name 'vmmDetails'`
经核对, vmmDetails明显存在于`/usr/share/virt-manager/virtManager/details.py`中, 估计是上次使用了2.2.1安装, 为了解决gtk报错使用`python3 setup.py install --force`降级安装了2.0.0, 因为历史文件干扰导致, 使用`rm -rf /usr/share/virt-manager`删除再用`python3 setup.py install --force`安装即可.

### virt-manager调试
`virt-manager --debug`

### virt-manager打开新建虚拟机界面报错"Error: No active connection to install on"
调试日志报"Autostart connect error: Unable to connect to libvirt qemu:///system."

原因未知.

> 当前环境是虚拟机, kvm-ok验证/dev/kvm不存在.

### virt-manager新建连接报错: `Cannot recv data: ssh_askpass: exec(/usr/bin/ssh-askpass): No such file or directory`
在virt-manager所在机器执行`apt install ssh-askpass`