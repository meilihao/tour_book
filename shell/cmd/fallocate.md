# fallocate

## 描述

为文件预分配物理空间, 其分配出的内容是空字符.

## 选项
- l : 文件大小, 默认是字节, 也可后跟k、m、g、t、p、e来指定单位，分别代表KB、MB、GB、TB、PB、EB.

## example
```sh
$ fallocate -l 100m new_file # 预分配一个100m大小的文件
$ od -tc new_file
```
