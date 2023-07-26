# stat

## 描述

查看文件状态, 比如时间戳

文件时间戳分类:
- atime : access time, 最后访问时间. 读取文件内容时，就会更新该时间
- mtime : modification time，最后修改时间. 当文件的内容(而不是文件属性)更改时会更新该时间
- ctime : change/status time, 最后状态时间. 当文件的状态(比如更改了属性/文件内容/位置)改变时会更新该时间

其他:
- IO Block: 文件块大小
- Blocks: 使用的底层块数(单位是512B from `man stat`)

  du获取的大小是 Blocks * 512
- IO Block: 文件系统的块大小
- device: 文件所在设备的设备号, 以十六进制和十进制表示. non-device 文件的情况下, 不会打印.

  它是文件所在设备的主设备号和次设备号组合成一个值. 804h 是主要设备 8, 次要设备 4.
- Context: selinux上下文

## 格式

  stat [OPTION]... FILE...

## 选项
- f : 获取文件所属fs的相关信息

## 例
```bash
# stat -c "%U:%G" /etc
root:root
```
