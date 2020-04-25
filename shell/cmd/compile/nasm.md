# nasm
## 选项
- -f : 指令指定要汇编的目标平台，可以使用`nasm -hf`查看支持的平台模式, linux x64是elf64

## FAQ
### ELF与BIN文件区别
- Gcc 编译出来的是ELF文件, 通常`gcc –o test test.c`生成的test文件就是ELF格式的，在linuxshell下输入`./test`就可以执行.
- Bin 文件是经过压缩的可执行文件，去掉ELF格式的东西, 是直接的内存映像的表示. 在系统没有加载操作系统的时候可以执行, 因此引导程序必定是bin.