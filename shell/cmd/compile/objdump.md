# objdump
一个对代码进行快速反编译的工具, 适合反编译未被加固、精简(stripped)或者以任何方式混淆的普通二进制文件.

缺点:
1. 依赖ELF节头(没有时不能工作), 且不会进行控制流分析.

## 示例
```
objdump –D <elf_object>  // 查看 ELF 文件中所有节的数据或代码
objdump –d <elf_object>  // 只查看 ELF 文件中的程序代码
objdump –tT <elf_object> // 查看所有符号
```