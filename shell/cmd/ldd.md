# ldd

## 描述

查看程序或者库文件所依赖的共享库列表

> ldd不是一个可执行程序，而只是一个shell脚本,它显示可执行模块的dependency的工作原理，其实质是通过ld-linux.so（elf动态库的装载器）来实现的.

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
