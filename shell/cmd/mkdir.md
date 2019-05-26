# mkdir

## 描述

创建目录

## 格式

  mkdir [OPTION] [directory]...

## 选项

- -m : 建立目录的同时设置目录的权限
- -p : 递归创建目录, 如果目录已存在则会跳过
- -v : 显示创建目录的过程

## 例
```sh
$ mkdir -m 700 /usr/local/test
$ mkdir -pv test/{a,b}/{1,2} # 同时创建多个多级目录
mkdir: 已创建目录 'test'
mkdir: 已创建目录 'test/a'
mkdir: 已创建目录 'test/a/1'
mkdir: 已创建目录 'test/a/2'
mkdir: 已创建目录 'test/b'
mkdir: 已创建目录 'test/b/1'
mkdir: 已创建目录 'test/b/2'
$ mkdir  -pv test2/dir{1..3} test3/{a..c}
mkdir: 已创建目录 'test2'
mkdir: 已创建目录 'test2/dir1'
mkdir: 已创建目录 'test2/dir2'
mkdir: 已创建目录 'test2/dir3'
mkdir: 已创建目录 'test3'
mkdir: 已创建目录 'test3/a'
mkdir: 已创建目录 'test3/b'
mkdir: 已创建目录 'test3/c'
```
