# file

## 描述

显示给定文件的类型

## 选项
- -b : 输出精简格式(即有文件类型但没有文件名)
- -L : 直接显示符号链接所指向文件的类型
- -z : 显示压缩文件的信息
- -i : 如果文件不是常规文件, 则不进一步对文件类型进行分类

## 例

    # file /bin/ls # 判断当前系统的位数
    # file t.elf # 输出包含"not stripped"表示有符号表
    # file -i xxx # 查看文件编码
    # file zsha2 # 查看开发语言, 这里是golang
    zsha2: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=9h4qNkFzYZaEW35kATVl/JI23wkfKIQqWQnWPQmI_/A2rlxTCOrY5CK4DB9Ypw/3t5KM7fcCCtzwU2Btlpl, stripped # 目标文件是`LSB relocatable`, 因为目标文件是可重定位文件, 即其中的符号尚未定位, 可用nm验证

## FAQ
### 创建稀疏文件
```bash
# dd if=/dev/zero of=sparse_file bs=1M seek=100 count=0
# truncate -s 100M sparse_file
# fallocate -o 500M -l 1G testfile # 从500MB的位置开始预留1GB的空间, 大小是1.5G
```

fallocate: 为文件预分配物理空间, 但不初始化它们, 不是生成空洞文件
truncate: 生成的是空洞文件, 并不占用实际的磁盘空间