# ld
ref:
- [充分理解Linux GCC 链接生成的Map文件](https://zhuanlan.zhihu.com/p/502051758)

链接器命令, 将`.o`文件转成可执行文件

## 选项
- --dynamic-linker /lib/ld-linux.so.2 : 采用32-bit动态连接
- --dynamic-linker /lib64/ld-linux-x86-64.so.2 : 采用64-bit动态连接
- -e : 指定程序入口符号
- -m elf_i386 : 生成 32-bit 的程序
- -L : 将比如lib32-glibc的库加入库搜索路径
- -lc : 连接标准 C 语言库， 比如printf
- -N : 指定可读写的text和data(section)
- -Ttext 指定text段的加载地址

> lld可使用`-e xxx`指定程序入口; clang不加`-nostdlib`时会连接ld.so即反汇编时main函数前面的部分代码.

## example
```bash
# ld [-L/usr/lib64] -luring --verbose # test liburing, like `gcc -lhdf5 --verbose`
```

## FAQ
### os启动起点
由ld链接器决定

### `ld cannot find library but it exists`
```bash
# ninjia
...
/usr/bin/c++ ... -lcurl  -Wl,-Bstatic  -luring  -Wl,-Bdynamic  ../thirdParty/isa-l_crypto/.libs/libisal_crypto.a  bin/libhashtable.a  -lm  -lrt  -ldl  -lrdmacm  -libverbs  -lpthread ...
/usr/bin/ld: cannot find -luring
```

liburing is so, use `-lcurl  -Wl,-Bdynamic  -luring`