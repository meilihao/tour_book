# objdump
一个对代码进行快速反编译的工具, 适合反编译未被加固、精简(stripped)或者以任何方式混淆的普通二进制文件.

缺点:
1. 依赖ELF节头(没有时不能工作), 且不会进行控制流分析.

## 选项
- -C : 将C++符号名逆向解析
- -d : 将所有包含指令的section反汇编
- -D : 反汇编所有的程序段，包括：代码段、数据段、只读数据段以及一些系统段等等
- -h : 打印elf文件的关键的section的基本信息, 完整信息可用`readelf -S xxx`获取
- -l : 反汇编代码中插入源代码的文件名和行号
- -j section : 仅反汇编指定的section。可以有多个-j参数来选择多个section
- -R : 查看`.o`定义的symbol
- -s : 将所有section的内容以16进制输出
- -S : 将代码段反汇编的同时，将反汇编代码和源代码交替显示，编译时需要给出-g，即需要调试信息
- -t : 打印符号表
- -x : 获得比`-h`更详细的各section信息

## 示例
```
objdump -D <elf_object>  # 查看 ELF 文件中所有节的数据或代码
objdump -d <elf_object>  # 只查看 ELF 文件中的程序代码
objdump -dS <elf_object>  # 查看 ELF 文件中的程序代码及源码
objdump -T xxx.so #  查看动态符号表
objdump -t <elf_object> # 查看所有符号. -T 和 -t 选项在于 -T 只能查看动态符号，如库导出的函数和引用其他库的函数，而 -t 可以查看所有的符号，包括数据段的符号
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