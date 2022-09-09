# w3m
字符浏览器

> lynx不支持https

## build
```bash
# dnf install g++
# ./configure --with-charset=GBK --with-charset=UTF-8 --with-gc=/usr/local/lib
# make && make install
```

> centos9构建失败: 构建时生成的functable.c为空, 因此执行到Makefile的`mktable 100 functable.tab > functable.c`时发生了core dumped.

## FAQ
### 乱码
- 在乱码页面, 按快捷键o进入设置界面, 将"Charset Setting"的"Display charset"设为"Unicode (UTF-8)", 再翻页到末尾点击"OK"按钮保存即可. 
- `w3m -I GBK  www.baidu.com`,解决中文乱码.

### `gc.h not found`
```bash
# wget https://www.hboehm.info/gc/gc_source/gc6.8.tar.gz
# tar -xf gc6.8.tar.gz && cd gc6.8
# ./configure
# make && make install
# vim /etc/ld.so.conf.d/gc.conf
/usr/local/lib
# ldconfig
```

### 编译报`while loading shared libraries: libgc.so.1: cannot open shared object file: No such file or directory`
```bash
# vim /etc/ld.so.conf.d/gc.conf
/usr/local/lib
# ldconfig
```