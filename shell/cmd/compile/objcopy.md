# objcopy(object copy)
用来分析和修改任意类型的 ELF 目标文件,还可以修改 ELF 节,或将 ELF 节复制
到 ELF 二进制中(或从 ELF 二进制中复制 ELF 节).

## 示例
```
objcopy -only-section=.data <infile> <outfile> // 将.data 节从一个 ELF 目标文件复制到另一个文件中
objcopy -I binary -O elf32-i386 -B i386 image.jpg image.o # 将二进制文件(图片)做成elf文件的一个段
objdump -ht image.o # 查看image.o的段内容
```