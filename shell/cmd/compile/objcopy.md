# objcopy(object copy)
用来分析和修改任意类型的 ELF 目标文件,还可以修改 ELF 节,或将 ELF 节复制
到 ELF 二进制中(或从 ELF 二进制中复制 ELF 节).

## 示例
```
objcopy –only-section=.data <infile> <outfile> // 将.data 节从一个 ELF 目标文件复制到另一个文件中
```