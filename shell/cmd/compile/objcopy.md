# objcopy(object copy)
用来分析和修改任意类型的 ELF 目标文件,还可以修改 ELF 节,或将 ELF 节复制
到 ELF 二进制中(或从 ELF 二进制中复制 ELF 节).

## 选项
- -I bfdname: 指定输入文件的BFD标准格式名:elf32-little, elf32-big等
- -O bfdname: 指定输出文件的BFD标准格式名
- -F bfdname: 指定输入输入文件的BFD标准格式名, 目标文件格式, 只用于在目标和源之间传输数据, 不转换
- -j sectionname: 只将由sectionname指定的节复制到输出文件, 可以多次指定
- -R sectionname: 从输出文件中去除由sectionname指定的节, 可以多次指定
- -S : 不从源文件复制符号信息和重定位信息
- -g: 不从原文件复制调试信息和相关的段. 对使用`-g`编译生成的可执行文件执行后, 生成的结果几乎和不使用`-g`进行编译生成的可执行文件一样
- -G symbolname: 只保留symbolname为全局变量, 让其他变量都是文件局部变量, 这样外部不可见
- -L symbolname: 将变量symbolname变为文件局部变量, 可以多次指定
- -x: 不从源文件中复制非全局变量

## 示例
```
objcopy -only-section=.data <infile> <outfile> // 将.data 节从一个 ELF 目标文件复制到另一个文件中
objcopy -I binary -O elf32-i386 -B i386 image.jpg image.o # 将二进制文件(图片)做成elf文件的一个段
objdump -ht image.o # 查看image.o的段内容
```