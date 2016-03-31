# whatis

## 描述

查询命令的简要说明,等同于使用`man -f`命令.
whatis会显示命令所在的具体的文档类别,方便在man page中查找.

## 参数

- -w : 使用正则匹配

## 例
```sh
$whatis -w "loca*"
```
