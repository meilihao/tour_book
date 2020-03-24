# dd

## 描述

文件复制和转换工具

## 选项
- if : 输入来源
    - `/dev/zero` : 提供无限个空字符（NULL/0x00）
    - `/dev/urandom` : 提供不为空字符的随机字节数据流
- of : 输出到指定位置
- bs : 单次读取和写入的字节数
- count : 操作次数, 因此文件大小 = bs * count
- seek=BLOCKS : 跳过输出文件中指定大小的部分来创建大文件，但是并不实际写入, 因此速度很快. 同时因为不实际写入磁盘，所以在容量只有10G的硬盘上创建100G的此类文件也是可以的 

## example
```sh
$ dd if=/dev/zero of=test bs=1M count=1000 # 在当前目录下会生成一个1000M的test文件，文件内容为全0, 但是这样为实际写入硬盘，文件产生速度取决于硬盘读写速度
$ seq 1000000 | xargs -i dd if=/dev/zero of={}.dat bs=1024 count=1 # 随机生成1百万个1K的文件
$ dd if=/dev/cdrom of=centos.iso # 将光驱设备中的光盘制作成 iso 格式的镜像文件
```