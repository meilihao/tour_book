# freeswitch
v1.8

[安装文档](https://freeswitch.org/confluence/display/FREESWITCH/Installation)

## error
> 安装缺失lib后需要清理`make clean && ./configure && make`

1. bootstrap: libtool not found.           You need libtool version 1.5.14 or newer to build FreeSWITCH from source
```
dpkg -L libtool

发现没有/usr/bin/libtool

dpkg -l libtool

libtool 是2.4.6-2版本的

在ubuntu只有libtoolize，修改bootstrap.sh，

libtool=${LIBTOOL:-`${LIBDIR}/apr/build/PrintPath glibtool libtool libtool22 libtool15 libtool14 libtoolize`}
```

1. mod_lua.cpp:37:10: fatal error: lua.h: 没有那个文件或目录
```
$ sudo apt install liblua5.3-dev 
$ cp /usr/include/lua5.3/*.h    src/mod/languages/mod_lua 
```

1. /usr/bin/ld: cannot find -llua
```
$ sudo find -name "liblua.so"
$ /usr/lib/x86_64-linux-gnu
$ sudo ln -sf liblua5.3.so liblua.so
```

1. You must install libopus-dev to build mod_opus
```
$ sudo apt install  libopus-dev
$ make clean && ./configure && make # 需要清理
```

1. You must install libsndfile-dev to build mod_sndfile
```
$ sudo apt install  libsndfile-dev
$ make clean && ./configure && make # 需要清理
```

## 中文语音下载地址
https://files.freeswitch.org/releases/sounds/
