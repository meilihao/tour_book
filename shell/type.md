# type

## 描述

用来显示指定命令的类型

一个命令的类型可以是如下之一:

- alias 别名
- keyword 关键字，Shell保留字
- function 函数，Shell函数
- builtin 内建命令，Shell内建命令
- file 文件，磁盘文件，外部命令
- unfound 没有找到

## 选项

- -a：显示所有可能的类型，比如有些命令如pwd是shell内建命令，也可以是外部命令.


