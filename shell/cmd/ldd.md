# ldd

## 描述

查看程序或者库文件所依赖的共享库列表

> ldd不是一个可执行程序，而只是一个shell脚本,它显示可执行模块的dependency的工作原理，其实质是通过ld-linux.so（elf动态库的装载器）来实现的.

## 选项
- -v : 输出detail

## 补充

一个可执行文件链接了哪些动态库呢或在遇到“error while loading shared libraries”时，我们难免会对此产生好奇.

查看该信息的方法是通过ldd，如

```shell
$ ldd chrome
  linux-vdso.so.1 =>  (0x00007fff52dff000) # 程序需要依赖什么库 => 系统提供的与程序需要的库所对应的库 (库加载的开始地址)
  libX11.so.6 => /usr/lib/libX11.so.6 (0x00007f0caebe4000)
  libdl.so.2 => /lib/libdl.so.2 (0x00007f0cae9e0000)
  libXrender.so.1 => /usr/lib/libXrender.so.1 (0x00007f0cae7d6000)
  libXss.so.1 => /usr/lib/libXss.so.1 (0x00007f0cae5d3000)
  libXext.so.6 => /usr/lib/libXext.so.6 (0x00007f0cae3c1000)
  librt.so.1 => /lib/librt.so.1 (0x00007f0cae1b9000)
  ....（略）
```

要想看系统还没安装的动态库，可以借用grep：

```shell
$ldd chrome | grep 'not found'
  libnss3.so.1d => not found
  libnssutil3.so.1d => not found
  libsmime3.so.1d => not found
  libplc4.so.0d => not found
  libnspr4.so.0d => not found
```

## FAQ
1. libxxx => not found
解决方法:
1. 对应的so库是否存在,不存在则安装
1. 已存在对应的so库, 则使用`sudo ldconfig -v`刷新so缓存

### 执行guestfish报`libselinux.so.1: no version information available (required by xxx)`
guestfish依赖的libselinux版本不对

> readelf -d xxx.so // SONAME包含主版本

### ldd xxx.so报`no version information available`
情况:
1. 依赖的so不存在, 需要安装
2. 依赖的so存在若干版本, 但使用了错误的版本

	比如: 调用的qt库版本与可执行程序编译时用的qt库版本不一致

### 查看so依赖的GLIBC版本
`strings /usr/lib/x86_64-linux-gnu/libgbm.so.1 | grep ^GLIBC`
