# ls

## 描述

显示目录/文件的相关属性

## 参数
- -a : 显示所有文件包括隐藏文件
- -c : 根据ctime排序
- --full-time : 已完整的时间格式输出
- -F : 输出结尾追加文件类型(可执行文件)
- -h : 已人类友好形式输出, 比如1KB, 2GB...
- -i : 显示inode信息
- -l : 使用长格式列出文件及目录的信息
- -t : 根据mtime排序, 默认是以文件名排序
- -r : 使用反向排序
- -S : 按照文件大小排序
- -R : 递归输出所有子目录
- --time={atime,ctime} : 输出指定时间, 默认是mtime
- -u : 根据atime排序

文件类型:
- 可执行文件 : *
- 目录 : /
- socket : =
- @ : 符号连接
- | : 管道

> `ls -l`的第二列是硬链接个数

## 格式

    ls [OPTION]... [FILE]...

## 例
```sh
$ln -s /usr/mengqc/mub1 /usr/liu/abc # 在目录/usr/liu下建立一个符号链接文件abc，使它指向目录/usr/mengqc/mub1
```