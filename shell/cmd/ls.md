# ls

## 描述

显示目录/文件的相关属性.

注意:
1. 输出的权限列表带后缀`+`表示启用了ACL权限.

## 参数
- -a : 显示所有文件包括隐藏文件
- -A : 同-a，但不列出“.”(表示当前目录)和“..”(表示当前目录的父目录)
- -c : 根据ctime排序
- -C : 按列输出, 即纵向排序
- -d : 仅显示目录
- --full-time : 已完整的时间格式输出, 比如`long-iso`
- -F : 输出结尾追加文件类型(`*:可执行文件;/,目录;@,符号链接; |, FIFO; =, socket`)
- -h : 已人类友好形式输出, 比如1KB, 2GB...
- -i : 显示inode信息
- -k : 以k字节的形式表示大小
- -l : 使用长格式列出文件及目录的信息
- -L : 不解析link
- -m : 横向输出文件名, 并以`,`作为分隔符
- -o : 显示信息不包括组信息
- -p : 给文件夹名称追加`/`
- -q : 用"?"代替不可打印的字符
- -Q : 把输出的文件名用双引号包裹
- -t : 根据mtime排序, 默认是以文件名排序
- -r : 使用反向排序
- -s : 在每个文件后输出其大小
- -S : 按照文件大小排序
- -R : 递归输出所有子目录
- --time={atime,ctime} : 输出指定时间, 默认是mtime
- --time-style=long-iso : 输出格式为`2022-07-27 10:20`
- -u : 根据atime排序
- -x : 按列输出, 即横向输出
- -Z : 查看文件的安全上下文值 for selinux. 见selinux

文件类型:
- 可执行文件 : *
- 目录 : /
- socket : =
- @ : 符号连接
- | : 管道

> `ls -l`的第二列是硬链接个数

## 格式

    ls [OPTION]... [FILE]...

## 通配符
- * : 任意字符
- ? : 单个任意字符
- [a-z] : 单个小写字母
- [A-Z] : 单个大写字母
- [a-Z] : 单个字母
- [0-9] : 单个数字
- [[:alpha:]] : 任意字母
- [[:upper:]] : 任意大写字母
- [[:lower:]] : 任意小写字母
- [[:digit:]] : 所有数字
- [[:alnum:]] : 任意字母加数字
- [[:punct:]] : 标点符号

## 例
```sh
# ls -Al --time-style=+"" #  不显示时间
# ls whateveryouwant | xargs -n 1 basename # 仅获取文件名
# ls  -ld  * # 输出不带"total ...", 仅罗列当前路径
# ls  -lL 2>dev/null | awk '{print $9}' # 仅显示文件名
# ls -F | grep '/$' |head -n 1 | 仅罗列目录
```

## FAQ
### `ls /`卡住
`/`下有挂载目录 from iscsi, iscsi服务端取消分配后, iscsi client没有umount

### CanonicalPath(规范路)
一个相对路径为`.././Java.txt`的文件, 那么:
- 它的绝对路径是 `/Users/androidyue/Documents/projects/PathSamples/.././Java.txt`
- 它的规范路径是 `/Users/androidyue/Documents/projects/Java.txt`