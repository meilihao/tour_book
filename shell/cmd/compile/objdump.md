# objdump
一个对代码进行快速反编译的工具, 适合反编译未被加固、精简(stripped)或者以任何方式混淆的普通二进制文件.

缺点:
1. 依赖ELF节头(没有时不能工作), 且不会进行控制流分析.

## 选项
- -a: 显示文档头
- -b bfdname: 指定目标码格式位bfd标准格式名bfdname, 比如binary
- -C : 将C++符号名逆向解析
- -d : 将所有包含指令的section反汇编
- -D : 反汇编所有的程序段，包括：代码段、数据段、只读数据段以及一些系统段等等
- -EB/-EL: 指定目标文件的小端, 将影响反汇编出来的指令
- -f: 显示每个文件的整体header信息
- -g: 显示调试信息, 尽量解析保存在文件中的调试信息并以C语言的语法显示出来
- -h : 打印elf文件的关键的section的header信息, 完整信息可用`readelf -S xxx`获取
- -l : 反汇编代码中插入源代码的文件名和行号
- -i : 显示对于`-b`或`-m`可用的架构和目标格式列表
- -j section : 仅显示指定的section. 可以有多个`-j`参数来选择多个section
- -m machine : 指定反汇编目标文件的架构
- -r: 显示文件的重定位入口, 如果和`-d`或者`-D`联用, 重定位部分以反汇编后的格式显示
- -R : 显示文件的动态重定位入口, 仅对动态目标文件有意义, 比如共享库
- -s : 显示指定节的完整内容, 默认所有的非空节都会被显示
- -S : 将代码段反汇编的同时，将反汇编代码和源代码交替显示即尽可能反汇编出源码. 推荐编译时需要给出-g，即需要调试信息
- -t : 打印符号表
- -x : 获得比`-h`更详细的各section信息, 包括符号表, 重定位入口

## 示例
```
objdump -D <elf_object>  # 查看 ELF 文件中所有节的数据或代码
objdump -d <elf_object>  # 只查看 ELF 文件中的程序代码
objdump -dS <elf_object>  # 查看 ELF 文件中的程序代码及源码
objdump -T xxx.so #  查看动态符号表
objdump -t <elf_object> # 查看所有符号. -T 和 -t 选项在于 -T 只能查看动态符号，如库导出的函数和引用其他库的函数，而 -t 可以查看所有的符号，包括数据段的符号
objdump -r xxx.o # 查看重定位表

    输出分析:
    - OFFSET : 重定位入口在该段的偏移量
```

## FAQ
### 反汇编判断开发语言
1. 因为重载, c++的函数符号名带特定的前缀和后缀.

### 解读`objdump -h xxx`
- Size : section len
- File Offset : 段在elf文件中的偏移位置
- CONTENTS, ALLOC, LOAD, READONLY, DATA : 表示段的各种属性

    - CONTENTS : 该段在elf文件中真实存在. `.bss`没有CONTENTS表示其实际上在elf文件中不存在内容


size xxx命令可查看elf文件的text, data, bss段的长度, dec是这3个段长度和的10进制表示, hex是这3个段长度和的16进制表示.

### `traps: cdp[701] trap invalid opcode ip:7fb03d456f82 sp:7fff22b7b1b0 error:0 in librocksdb.so.9.1.2[7fb03d071000+659000]`
ref:
- [故障分析 | MongoDB 5.0 报错 Illegal instruction 解决](https://opensource.actionsky.com/20220124-mongodb/)
- [记录一次线上进程异常崩溃的排查过程](https://blog.csdn.net/u010231230/article/details/118543846)
- [error number](https://worthsen.blog.csdn.net/article/details/106896795)

    解释`error:0`的含义
- [[x86][linux]AVX512指令引起的进程crash](https://cloud.tencent.com/developer/article/1356935)

env:
- Intel(R) Xeon(R) Gold 6133 CPU, 正常运行环境, 且是librocksdb.so的构建环境
- Intel(R) Xeon(R) E-2224 CPU, cdp崩溃环境

so内存基址: 7fb03d071000, 长度659000; 崩溃ip(指令地址): 7fb03d456f82

相对地址为7fb03d456f82-7fb03d071000=3e5f82, 再结合`objdump -ld xxx.so > xxx.objdump`排查

未找到错误指令, 应该是构建时使用了E-2224不支持的指令集, 而构建时追加`PORTABLE=1`(以兼容更多cpu)则正常.

根据ref `[x86][linux]AVX512指令引起的进程crash`, xxx.objdump存在vmovdqa64指令(AVX512F支持的指令集), 而E-2224不支持AVX512F.